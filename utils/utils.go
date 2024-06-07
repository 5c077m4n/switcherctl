// Package utils general helpers
package utils

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"time"
)

// WattsToAmps Convert power consumption to watts to electric current in amps.
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
