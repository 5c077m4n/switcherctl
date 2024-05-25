package connection

import (
	"fmt"
	"net"
	"switcherctl/consts"
)

type Connection struct{ serve *net.UDPConn }

func (c *Connection) RemoteAddress() net.Addr     { return c.serve.RemoteAddr() }
func (c *Connection) Write(b []byte) (int, error) { return c.serve.Write(b) }
func (c *Connection) ReadFromUDP() (string, error) {
	b := make([]byte, 1024)
	n, _, err := c.serve.ReadFromUDP(b)
	if err != nil {
		return "", err
	}

	return string(b[0:n]), nil
}
func (c *Connection) Close() error { return c.serve.Close() }

func TryNew(ip string, port int) (*Connection, error) {
	deviceHostIP := consts.DefaultIP + ":" + fmt.Sprint(port)
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
