// Package connections for creating a listening connection to a Switcher device
package connections

import (
	"errors"
	"net"
	"time"
)

// BidirectionalConn the struct for the Switcher connection
type BidirectionalConn struct{ conn *net.UDPConn }

func (c *BidirectionalConn) login() error {
	return errors.New("not implemented")
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
	if err := bidirConn.login(); err != nil {
		return nil, err
	}

	return bidirConn, nil
}
