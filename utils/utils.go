package utils

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"time"
)

// Convert power consumption to watts to electric current in amps.
func WattsToAmps(watts int) int {
	return int(math.Round(float64(watts) / float64(220)))
}

func CurrentTimeLE() string {
	now := time.Now()
	epochSeconds := now.Unix()

	timeLEBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeLEBuf, uint64(epochSeconds))

	timeHexBuf := make([]byte, 8)
	hex.Encode(timeHexBuf, timeLEBuf)

	return string(timeHexBuf)
}
