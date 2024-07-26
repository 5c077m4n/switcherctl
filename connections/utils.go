package connections

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"strings"
	"time"
)

// currentTimeHexLE current time in little edian hex format
func currentTimeHexLE() string {
	now := time.Now()
	epochSeconds := now.Unix()

	timeLEBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeLEBuf, uint64(epochSeconds))

	return hex.EncodeToString(timeLEBuf)
}

var crc32q = crc32.MakeTable(0x1021)

// signPacketWithCRCKey use CRC to sign a hex packet
func signPacketWithCRCKey(hexPacket string) (string, error) {
	binPacketChecksum := crc32.Checksum([]byte(hexPacket), crc32q)

	bePacketChecksum := make([]byte, 8)
	binary.BigEndian.PutUint32(bePacketChecksum, binPacketChecksum)

	packetChecksumHex := hex.EncodeToString(bePacketChecksum)
	packetChecksumHexSlice := packetChecksumHex[6:8] + packetChecksumHex[4:6]

	binKey, err := hex.DecodeString(packetChecksumHexSlice + strings.Repeat("30", 32))
	if err != nil {
		return "", errors.Join(ErrSignPacket, err)
	}
	binKeyChecksum := crc32.Checksum(binKey, crc32q)

	beBinKeyChecksum := make([]byte, 8)
	binary.BigEndian.PutUint32(beBinKeyChecksum, binKeyChecksum)

	binKeyChecksumHex := hex.EncodeToString(beBinKeyChecksum)
	binKeyChecksumHexSlice := binKeyChecksumHex[6:8] + binKeyChecksumHex[4:6]

	return hexPacket + packetChecksumHexSlice + binKeyChecksumHexSlice, nil
}
