// Package parse for parsing network messages
package parse

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"switcherctl/consts"
	"time"
)

const switcherMessagePrefix = "fef0"

// ErrInvalidMessage error for when the message is too short and so can't be parsed
var ErrInvalidMessage = errors.New("the received message is invalid (too short)")

func secondsToISOTime(totalSeconds int) (*time.Duration, error) {
	minutes, seconds := totalSeconds/60, totalSeconds%60
	hours, minutes := minutes/60, minutes%60
	duration, err := time.ParseDuration(fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds))
	if err != nil {
		return nil, err
	}

	return &duration, nil
}

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

// GetDeviceName extract the device's name from the message
func (parser *DatagramParser) GetDeviceName() string {
	return string(parser.msg[42:74])
}

// GetDeviceID extract the device's ID from the message
func (parser *DatagramParser) GetDeviceID() string {
	return parser.msgHex[36:42]
}

// GetDeviceKey extract the device's key from the message
func (parser *DatagramParser) GetDeviceKey() string {
	return parser.msgHex[80:82]
}

// GetDeviceMAC extract the device's MAC address from the message
func (parser *DatagramParser) GetDeviceMAC() string {
	rawMAC := strings.ToUpper(parser.msgHex[160:172])
	return strings.Join(
		[]string{rawMAC[0:2], rawMAC[2:4], rawMAC[4:6], rawMAC[6:8], rawMAC[8:10], rawMAC[10:12]},
		":",
	)
}

// IsPoweredOn extract the device's state from the message
func (parser *DatagramParser) IsPoweredOn() bool {
	state := parser.msgHex[266:268]
	return state == "01"
}

// GetTimeToShutdown get the datetime of the auto power off
func (parser *DatagramParser) GetTimeToShutdown() (*time.Duration, error) {
	autoShutdown := parser.msgHex[310:318]
	leAutoShutdown := autoShutdown[6:8] + autoShutdown[4:6] + autoShutdown[2:4] + autoShutdown[0:2]
	autoShutdownSeconds, err := strconv.ParseUint(leAutoShutdown, 16, 32)
	if err != nil {
		return nil, err
	}

	endsIn := time.Duration(autoShutdownSeconds * uint64(time.Second))
	return &endsIn, nil
}

// GetRemainingTime get the remaining time to work end
func (parser *DatagramParser) GetRemainingTime() (*time.Duration, error) {
	remainingHex := parser.msgHex[294:302]
	leRemainingSeconds := remainingHex[6:8] + remainingHex[4:6] + remainingHex[2:4] + remainingHex[0:2]
	remainingSeconds, err := strconv.ParseUint(leRemainingSeconds, 16, 32)
	if err != nil {
		return nil, err
	}

	endsIn := time.Duration(remainingSeconds * uint64(time.Second))
	return &endsIn, nil
}

// New create a DatagramParser instance
func New(msg []byte) DatagramParser {
	msgHex := hex.EncodeToString(msg)
	return DatagramParser{msg: msg, msgHex: msgHex}
}
