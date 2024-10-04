// Package connections for creating a listening connection to a Switcher device
package connections

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"switcherctl/consts"
	"time"
)

// BidirectionalConn the struct for the Switcher connection
type BidirectionalConn struct {
	addr           *net.UDPAddr
	conn           *net.UDPConn
	deviceID       string
	sessionID      string
	loginTimestamp string
}

func (c *BidirectionalConn) login() error {
	ip := c.addr.IP
	port := c.addr.Port

	conn, err := TryNewListener(ip, port)
	if err != nil {
		return errors.Join(ErrLoginFail, err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			panic(errors.Join(ErrLoginFail, closeErr))
		}
	}()

	data, err := conn.Read()
	if err != nil {
		return errors.Join(ErrLoginFail, err)
	}

	timestamp := currentTimeHexLE()
	loginPackateHex := fmt.Sprintf(consts.TemplatePacketLoginType1, timestamp, data.GetDeviceID())

	loginPackate, err := hex.DecodeString(loginPackateHex)
	if err != nil {
		return errors.Join(ErrLoginFail, err)
	}

	_, err = c.conn.Write(loginPackate)
	if err != nil {
		return errors.Join(ErrLoginFail, err)
	}

	responseBuf := make([]byte, 1024)
	n, err := c.conn.Read(responseBuf)
	if err != nil {
		return errors.Join(ErrLoginFail, err)
	}

	responseHex := hex.EncodeToString(responseBuf[:n])
	if len(responseHex) < 24 {
		return errors.Join(ErrLoginFail, ErrResponseTooShort)
	}

	c.sessionID = string(responseHex[16:24])
	c.loginTimestamp = timestamp

	return nil
}

// GetSchedules get device's work schedules
func (c *BidirectionalConn) GetSchedules() error {
	packet := fmt.Sprintf(
		consts.TemplatePacketGetSchedules,
		c.sessionID,
		c.loginTimestamp,
		c.deviceID,
	)
	signedPacket, err := signPacketWithCRCKey(packet)
	if err != nil {
		return errors.Join(ErrGetSchedules, err)
	}

	slog.Debug(
		"schedules",
		"packet", packet,
		"signed packet", signedPacket,
	)

	signedPacketRaw, err := hex.DecodeString(signedPacket)
	if err != nil {
		return errors.Join(ErrGetSchedules, err)
	}

	if _, err = c.conn.Write(signedPacketRaw); err != nil {
		return errors.Join(ErrGetSchedules, err)
	}

	responseBuf := make([]byte, 1024)
	n, err := c.conn.Read(responseBuf)
	if err != nil {
		return errors.Join(ErrGetSchedules, err)
	}

	responseHex := hex.EncodeToString(responseBuf[:n])
	slog.Debug(
		"schedules",
		"response", responseHex,
	)

	return consts.ErrNotImplemeted
}

// Close close the connection
func (c *BidirectionalConn) Close() error { return c.conn.Close() }

// TryNewBidirectionalConn try to create a new connection instance
func TryNewBidirectionalConn(ip net.IP, port int, deviceID string) (*BidirectionalConn, error) {
	addr := &net.UDPAddr{IP: ip, Port: port}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, errors.Join(ErrTryNewBidirectionalConn, err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, errors.Join(ErrTryNewBidirectionalConn, err)
	}

	bidirConn := &BidirectionalConn{addr: addr, conn: conn, deviceID: deviceID}
	if err := bidirConn.login(); err != nil {
		return nil, errors.Join(ErrTryNewBidirectionalConn, err)
	}

	return bidirConn, nil
}
