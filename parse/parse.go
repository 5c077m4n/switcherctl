// Package parse for parsing network messages
package parse

import (
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

// GetIPType1 get the IP of the device from the message
func (parser *DatagramParser) GetIPType1() (net.IP, error) {
	if len(parser.msgHex) < 160 {
		return nil, ErrInvalidMessage
	}

	beHexIP := parser.msgHex[152:160]

	ip := net.IP{}
	for i := 0; i <= 6; i += 2 {
		ipPart, err := strconv.ParseUint(beHexIP[i:i+2], 16, 8)
		if err != nil {
			return nil, err
		}
		ip = append(ip, uint8(ipPart))
	}

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
