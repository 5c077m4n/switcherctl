// Package connection for creating a connection to a Switcher device
package connection

import (
	"errors"
	"net"
	"switcherctl/consts"
	"switcherctl/parse"
	"time"
)

// ErrWrongRemote wrong remote error
var ErrWrongRemote = errors.New("message did not originate from a Switcher device")

// Connection the struct for the Switcher connection
type Connection struct {
	conn   *net.UDPConn
	remote *net.UDPAddr
}

// Read read the server's next message
func (c *Connection) Read() (*parse.DatagramParser, error) {
	messageBuffer := make([]byte, 1024)
	n, remoteAddr, err := c.conn.ReadFromUDP(messageBuffer)
	if err != nil {
		return nil, err
	}
	if !remoteAddr.IP.Equal(consts.DefaultIP) {
		return nil, ErrWrongRemote
	}

	data := parse.New(messageBuffer[:n])
	if !data.IsSwitcher() {
		return nil, ErrWrongRemote
	}

	return &data, nil
}

// Close close the connection
func (c *Connection) Close() error { return c.conn.Close() }

// TryNew try to create a new connection instance
func TryNew(ip net.IP, port int) (*Connection, error) {
	localAddr := &net.UDPAddr{IP: net.IP{0, 0, 0, 0}, Port: port, Zone: ""}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		return nil, err
	}

	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, err
	}

	return &Connection{
		conn:   conn,
		remote: &net.UDPAddr{IP: ip, Port: port, Zone: ""},
	}, nil
}
