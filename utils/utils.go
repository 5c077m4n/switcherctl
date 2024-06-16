// Package utils general helpers
package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"math"
	"strings"
	"switcherctl/consts"
	"time"
)

// WattsToAmps Convert power consumption from watts to electric current in amps.
func WattsToAmps(watts int) int {
	return int(math.Round(float64(watts) / float64(220)))
}

// CurrentTimeHexLE current time in little edian hex format
func CurrentTimeHexLE() string {
	now := time.Now()
	epochSeconds := now.Unix()

	timeLEBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeLEBuf, uint64(epochSeconds))

	return hex.EncodeToString(timeLEBuf)
}

// GenerateLoginPacketType1 generate a login packet to be sent to a device
func GenerateLoginPacketType1(deviceKey string) string {
	return fmt.Sprintf(consts.LoginPacketType1, CurrentTimeHexLE(), deviceKey)
}

var crc32q = crc32.MakeTable(0x1021)

// SignPacketWithCRCKey use CRC to sign a hex packet
func SignPacketWithCRCKey(hexPacket string) (string, error) {
	binPacketChecksum := crc32.Checksum([]byte(hexPacket), crc32q)

	bePacketChecksum := make([]byte, 8)
	binary.BigEndian.PutUint32(bePacketChecksum, binPacketChecksum)

	packetChecksumHex := hex.EncodeToString(bePacketChecksum)
	packetChecksumHexSlice := packetChecksumHex[6:8] + packetChecksumHex[4:6]

	binKey, err := hex.DecodeString(packetChecksumHexSlice + strings.Repeat("30", 32))
	if err != nil {
		return "", err
	}
	binKeyChecksum := crc32.Checksum(binKey, crc32q)

	beBinKeyChecksum := make([]byte, 8)
	binary.BigEndian.PutUint32(beBinKeyChecksum, binKeyChecksum)

	binKeyChecksumHex := hex.EncodeToString(beBinKeyChecksum)
	binKeyChecksumHexSlice := binKeyChecksumHex[6:8] + binKeyChecksumHex[4:6]

	return hexPacket + packetChecksumHexSlice + binKeyChecksumHexSlice, nil
}
