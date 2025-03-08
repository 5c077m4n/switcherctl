// Package connections for creating a listening connection to a Switcher device
package connections

import (
	"errors"
	"fmt"
	"net"
	"switcherctl/consts"
	"switcherctl/parse"
	"time"
)

// Listener the struct for the Switcher connection
type Listener struct {
	conn   *net.UDPConn
	remote *net.UDPAddr
}

// Read read the server's next message
func (l *Listener) Read() (*parse.DatagramParser, error) {
	messageBuffer := make([]byte, 1024)
	n, remoteAddr, err := l.conn.ReadFromUDP(messageBuffer)
	if err != nil {
		return nil, errors.Join(ErrListenerRead, err)
	}
	if l.remote.IP != nil && !remoteAddr.IP.Equal(l.remote.IP) {
		return nil, errors.Join(
			ErrListenerRead,
			ErrWrongRemote,
			fmt.Errorf("unkown IP address: %v", remoteAddr.IP),
		)
	}

	data := parse.New(messageBuffer[:n])
	if !data.IsSwitcher() {
		return nil, errors.Join(ErrListenerRead, ErrWrongRemote)
	}

	return &data, nil
}

// Close the connection
func (l *Listener) Close() error {
	if err := l.conn.Close(); err != nil {
		return errors.Join(ErrListenerClose, err)
	}
	return nil
}

// TryNewListener try to create a new connection instance
func TryNewListener(ip net.IP, port consts.Port) (*Listener, error) {
	localAddr := &net.UDPAddr{Port: int(port)}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		return nil, errors.Join(ErrTryNewListener, err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, errors.Join(ErrTryNewListener, err)
	}

	return &Listener{conn: conn, remote: &net.UDPAddr{IP: ip, Port: int(port)}}, nil
}
