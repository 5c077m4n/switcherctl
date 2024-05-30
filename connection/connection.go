// Package connection for creating a connection to a Switcher device
package connection

import (
	"fmt"
	"log"
	"net"
	"switcherctl/parse"
	"time"
)

// Connection the struct for the Switcher connection
type Connection struct{ serve *net.UDPConn }

// RemoteAddress get the remote address
func (c *Connection) RemoteAddress() net.Addr { return c.serve.RemoteAddr() }

// Write write a message to the remote server
func (c *Connection) Write(msg string) (int, error) { return c.serve.Write([]byte(msg)) }

// Read read the server's next message
func (c *Connection) Read() (*parse.DatagramParser, error) {
	messageBuffer := make([]byte, 1024)
	n, _, err := c.serve.ReadFromUDP(messageBuffer)
	if err != nil {
		return nil, err
	}

	data := parse.New(messageBuffer[0:n])
	return &data, nil
}

// Close close the connection
func (c *Connection) Close() error { return c.serve.Close() }

// TryNew try to create a new connection instance
func TryNew(ip string, port int) (*Connection, error) {
	deviceHostIP := fmt.Sprintf("%s:%d", ip, port)
	remoteAddr, err := net.ResolveUDPAddr("udp4", deviceHostIP)
	if err != nil {
		return nil, err
	}

	serve, err := net.DialUDP("udp4", nil, remoteAddr)
	if err != nil {
		return nil, err
	}
	if err := serve.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		log.Fatalln(err)
	}

	return &Connection{serve: serve}, nil
}
