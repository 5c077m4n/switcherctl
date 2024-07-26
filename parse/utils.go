package parse

import "math"

// WattsToAmps Convert power consumption from watts to electric current in amps.
func WattsToAmps(watts int) int {
	return int(math.Round(float64(watts) / float64(220)))
}
