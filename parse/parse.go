// Package parse for parsing network messages
package parse

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"strconv"
)

const (
	MessageLengthDefault = 165
	MessageLengthBreeze  = 168
	MessageLengthRunner  = 159
)

type DatagramParser struct{ msg []byte }

func (dp *DatagramParser) String() string { return string(dp.msg) }

func (dp *DatagramParser) IsSwitcher() (bool, error) {
	msgHex := hex.EncodeToString(dp.msg)

	decoded, err := hex.DecodeString(msgHex[0:4])
	if err != nil {
		return false, err
	}

	return string(decoded) == "fef0" &&
		(len(dp.msg) == MessageLengthDefault || len(dp.msg) == MessageLengthBreeze || len(dp.msg) == MessageLengthRunner), nil
}

func (dp *DatagramParser) GetIPType1() (net.IP, error) {
	hexIP := hex.EncodeToString(dp.msg)[152:160]
	ipAddress, err := strconv.ParseUint(hexIP[6:8]+hexIP[4:6]+hexIP[2:4]+hexIP[0:2], 16, 16)
	if err != nil {
		return nil, err
	}

	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, uint32(ipAddress))

	return ip, nil
}

// New create a DatagramParser instance
func New(msg []byte) DatagramParser {
	return DatagramParser{msg: msg}
}
