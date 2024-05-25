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

type DatagramParser struct{ message []byte }

func (dp *DatagramParser) IsSwitcher() (bool, error) {
	msgHex := hex.EncodeToString(dp.message)

	decoded, err := hex.DecodeString(msgHex[0:4])
	if err != nil {
		return false, err
	}

	return string(decoded) == "fef0" &&
		(len(dp.message) == MessageLengthDefault || len(dp.message) == MessageLengthBreeze || len(dp.message) == MessageLengthRunner), nil
}

func (dp *DatagramParser) GetIPType1() (net.IP, error) {
	hexIP := hex.EncodeToString(dp.message)[152:160]
	ipAddress, err := strconv.ParseUint(hexIP[6:8]+hexIP[4:6]+hexIP[2:4]+hexIP[0:2], 16, 16)
	if err != nil {
		return nil, err
	}

	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, uint32(ipAddress))

	return ip, nil
}

func New(msg []byte) DatagramParser {
	return DatagramParser{message: msg}
}
