// Package parse for parsing network messages
package parse

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"net"
	"strconv"
	"switcherctl/consts"
)

const switcherMessagePrefix = "fef0"

// ErrInvalidMessage error for when the message is too short and so can't be parsed
var ErrInvalidMessage = errors.New("the received message is invalid (too short)")

// DatagramParser struct to parse incoming messages
type DatagramParser struct {
	msgHex string
	msg    []byte
}

func (parser *DatagramParser) String() string { return string(parser.msg) }

// IsSwitcher test if message originates from a Swticher device
func (parser *DatagramParser) IsSwitcher() bool {
	return parser.msgHex[:4] == switcherMessagePrefix &&
		(len(parser.msg) == consts.MessageLengthDefault ||
			len(parser.msg) == consts.MessageLengthBreeze ||
			len(parser.msg) == consts.MessageLengthRunner)
}

// GetIPType1 get the IP of the device from a message
func (parser *DatagramParser) GetIPType1() (net.IP, error) {
	if len(parser.msgHex) < 160 {
		return nil, ErrInvalidMessage
	}

	hexIP := parser.msgHex[152:160]
	ipAddress, err := strconv.ParseUint(hexIP[6:8]+hexIP[4:6]+hexIP[2:4]+hexIP[0:2], 16, 32)
	if err != nil {
		return nil, err
	}

	ip := net.IP{}
	binary.BigEndian.PutUint32(ip, uint32(ipAddress))

	return ip, nil
}

// GetDeviceID extract the device's ID from the message
func (parser *DatagramParser) GetDeviceID() string {
	return hex.EncodeToString(parser.msg[40:41])
}

// New create a DatagramParser instance
func New(msg []byte) DatagramParser {
	msgHex := hex.EncodeToString(msg)
	return DatagramParser{msg: msg, msgHex: msgHex}
}
