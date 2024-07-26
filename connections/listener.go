// Package connections for creating a listening connection to a Switcher device
package connections

import (
	"errors"
	"net"
	"switcherctl/parse"
	"time"
)

// Listener the struct for the Switcher connection
type Listener struct {
	conn   *net.UDPConn
	remote *net.UDPAddr
}

// ErrWrongRemote wrong remote error
var ErrWrongRemote = errors.New("message did not originate from a Switcher device")

// Read read the server's next message
func (l *Listener) Read() (*parse.DatagramParser, error) {
	messageBuffer := make([]byte, 1024)
	n, remoteAddr, err := l.conn.ReadFromUDP(messageBuffer)
	if err != nil {
		return nil, errors.Join(ErrListenerRead, err)
	}
	if !remoteAddr.IP.Equal(l.remote.IP) {
		return nil, errors.Join(ErrListenerRead, ErrWrongRemote)
	}

	data := parse.New(messageBuffer[:n])
	if !data.IsSwitcher() {
		return nil, errors.Join(ErrListenerRead, ErrWrongRemote)
	}

	return &data, nil
}

// Close close the connection
func (l *Listener) Close() error { return l.conn.Close() }

// TryNewListener try to create a new connection instance
func TryNewListener(ip net.IP, port int) (*Listener, error) {
	localAddr := &net.UDPAddr{IP: net.IP{0, 0, 0, 0}, Port: port}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		return nil, err
	}

	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, err
	}

	return &Listener{conn: conn, remote: &net.UDPAddr{IP: ip, Port: port}}, nil
}
