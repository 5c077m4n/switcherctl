// Package connections for creating a listening connection to a Switcher device
package connections

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"switcherctl/consts"
	"switcherctl/utils"
	"time"
)

// BidirectionalConn the struct for the Switcher connection
type BidirectionalConn struct {
	conn      *net.UDPConn
	sessionID string
}

func (c *BidirectionalConn) login(ip net.IP, port int) error {
	conn, err := TryNewListener(ip, port)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Fatalln(closeErr)
		}
	}()

	data, err := conn.Read()
	if err != nil {
		return err
	}

	timestamp := utils.CurrentTimeHexLE()
	loginPackateHex := fmt.Sprintf(consts.LoginPacketType1Template, timestamp, data.GetDeviceID())

	loginPackate, err := hex.DecodeString(loginPackateHex)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(loginPackate)
	if err != nil {
		return err
	}

	responseBuf := make([]byte, 1024)
	n, err := c.conn.Read(responseBuf)
	if err != nil {
		return err
	}

	responseHex := hex.EncodeToString(responseBuf[:n])
	if len(responseHex) < 24 {
		return consts.ErrLoginFail
	}
	c.sessionID = string(responseHex[16:24])

	return consts.ErrNotImplemeted
}

// Close close the connection
func (c *BidirectionalConn) Close() error { return c.conn.Close() }

// TryNewBidirectionalConn try to create a new connection instance
func TryNewBidirectionalConn(ip net.IP, port int) (*BidirectionalConn, error) {
	localAddr := &net.UDPAddr{IP: ip, Port: port}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		return nil, err
	}

	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, err
	}

	bidirConn := &BidirectionalConn{conn: conn}
	if err := bidirConn.login(ip, port); err != nil {
		return nil, err
	}

	return bidirConn, nil
}
