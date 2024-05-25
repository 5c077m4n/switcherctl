// Package connection for creating a connection to a Switcher device
package connection

import (
	"fmt"
	"net"
)

// Connection the struct for the Switcher connection
type Connection struct{ serve *net.UDPConn }

// RemoteAddress get the remote address
func (c *Connection) RemoteAddress() net.Addr { return c.serve.RemoteAddr() }

// Write write a message to the remote server
func (c *Connection) Write(b []byte) (int, error) { return c.serve.Write(b) }

// ReadFromUDP read the server's next message
func (c *Connection) ReadFromUDP() (string, error) {
	b := make([]byte, 1024)
	n, _, err := c.serve.ReadFromUDP(b)
	if err != nil {
		return "", err
	}

	return string(b[0:n]), nil
}

// Close close the connection
func (c *Connection) Close() error { return c.serve.Close() }

// TryNew try to create a new connection instance
func TryNew(ip string, port int) (*Connection, error) {
	deviceHostIP := ip + ":" + fmt.Sprint(port)
	addr, err := net.ResolveUDPAddr("udp4", deviceHostIP)
	if err != nil {
		return nil, err
	}

	serve, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	return &Connection{serve: serve}, nil
}
