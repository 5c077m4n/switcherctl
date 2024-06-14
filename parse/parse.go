// Package parse for parsing network messages
package parse

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"switcherctl/consts"
	"time"
)

const switcherMessagePrefix = "fef0"

// ErrInvalidMessage error for when the message is too short and so can't be parsed
var ErrInvalidMessage = errors.New("the received message is invalid (too short)")

type (
	// DatagramParser struct to parse incoming messages
	DatagramParser struct {
		msgHex string
		msg    []byte
	}
	// DatagramParsedJSON a JSON representation of the network datagram's content
	DatagramParsedJSON struct {
		Name             string `json:"name"`
		IP               string `json:"ip"`
		ID               string `json:"id"`
		Key              string `json:"key"`
		MAC              string `json:"mac"`
		TimeToShutdown   string `json:"timeToShutdown"`
		TimeRemaining    string `json:"remainingTime"`
		PowerOn          bool   `json:"powerOn"`
		PowerConsumption uint64 `json:"powerConsumption"`
	}
)

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
	rawName := string(parser.msg[42:74])
	name := strings.TrimRight(rawName, "\u0000")
	return name
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
func (parser *DatagramParser) GetDeviceMAC() (*net.HardwareAddr, error) {
	rawMAC := strings.ToUpper(parser.msgHex[160:172])

	mac, err := net.ParseMAC(strings.Join(
		[]string{rawMAC[0:2], rawMAC[2:4], rawMAC[4:6], rawMAC[6:8], rawMAC[8:10], rawMAC[10:12]},
		":",
	))
	if err != nil {
		return nil, err
	}
	return &mac, nil
}

// IsPoweredOn extract the device's state from the message
func (parser *DatagramParser) IsPoweredOn() bool {
	state := parser.msgHex[266:268]
	return state == "01"
}

// GetTimeToShutdown get the datetime of the auto power off
func (parser *DatagramParser) GetTimeToShutdown() (*time.Duration, error) {
	beAutoShutdown := parser.msgHex[310:318]
	leAutoShutdown := beAutoShutdown[6:8] + beAutoShutdown[4:6] + beAutoShutdown[2:4] + beAutoShutdown[0:2]
	autoShutdownSeconds, err := strconv.ParseUint(leAutoShutdown, 16, 32)
	if err != nil {
		return nil, err
	}

	endsIn := time.Duration(autoShutdownSeconds * uint64(time.Second))
	return &endsIn, nil
}

// GetRemainingTime get the remaining time to work end
func (parser *DatagramParser) GetRemainingTime() (*time.Duration, error) {
	beRemainingHex := parser.msgHex[294:302]
	leRemainingSeconds := beRemainingHex[6:8] + beRemainingHex[4:6] + beRemainingHex[2:4] + beRemainingHex[0:2]
	remainingSeconds, err := strconv.ParseUint(leRemainingSeconds, 16, 32)
	if err != nil {
		return nil, err
	}

	endsIn := time.Duration(remainingSeconds * uint64(time.Second))
	return &endsIn, nil
}

// GetPowerConsumption get the power consumption (in watts)
func (parser *DatagramParser) GetPowerConsumption() (uint64, error) {
	hexPowerConsumption := parser.msgHex[270:278]
	return strconv.ParseUint(
		hexPowerConsumption[2:4]+hexPowerConsumption[0:2],
		16,
		32,
	)
}

// MarshalJSON returns a JSON struct of the datagram packet
func (parser *DatagramParser) MarshalJSON() ([]byte, error) {
	ip, err := parser.GetIPType1()
	if err != nil {
		log.Fatalln(err)
	}
	autoShutdown, err := parser.GetTimeToShutdown()
	if err != nil {
		log.Fatalln(err)
	}
	remaining, err := parser.GetRemainingTime()
	if err != nil {
		log.Fatalln(err)
	}
	mac, err := parser.GetDeviceMAC()
	if err != nil {
		log.Fatalln(err)
	}
	powerConsumption, err := parser.GetPowerConsumption()
	if err != nil {
		log.Fatalln(err)
	}

	return json.Marshal(&DatagramParsedJSON{
		Name:             parser.GetDeviceName(),
		IP:               ip.String(),
		ID:               parser.GetDeviceID(),
		Key:              parser.GetDeviceKey(),
		MAC:              mac.String(),
		TimeToShutdown:   autoShutdown.String(),
		TimeRemaining:    remaining.String(),
		PowerOn:          parser.IsPoweredOn(),
		PowerConsumption: powerConsumption,
	})
}

// New create a DatagramParser instance
func New(msg []byte) DatagramParser {
	msgHex := hex.EncodeToString(msg)
	return DatagramParser{msg: msg, msgHex: msgHex}
}
