package parse

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	workingMessageV4OffStateHex = "fef0a500023c0200000000008c290100000000000000000051329c6600000000000000000000f0fe060053776974636865725f56345f423534410000000000000000000000000000000003170a64665200000000000000000000000000000253776974636865725f56345f4235344100000000000000000000000000000000020401001c000000000000003bb3a154010000000000000084030000201c00000101eadecab9"
	workingMessageV4OnStateHex  = "fef0a500023c020000000000418401000000000000000000708da26600000000000000000000f0fe060053776974636865725f56345f423534410000000000000000000000000000000003170a64665200000000000000000000000000000253776974636865725f56345f4235344100000000000000000000000000000000020401001c000100ce0800007426a25401000000770300000d000000201c00000101a192baab"
)

func unhexMessage(t *testing.T, msgHex string) []byte {
	t.Helper()

	msg, err := hex.DecodeString(msgHex)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	return msg
}

func TestDatagramOffMessage(t *testing.T) {
	msg := unhexMessage(t, workingMessageV4OffStateHex)
	parser := New(msg)

	results, err := parser.ToJSON()

	if assert.NoError(t, err) {
		assert.Equal(
			t,
			&DatagramParsedJSON{
				Name:             "Switcher_V4_B54A",
				IP:               "10.100.102.82",
				ID:               "000000",
				Key:              "06",
				MAC:              "00:00:00:00:00:00",
				TimeToShutdown:   "2h0m0s",
				TimeRemaining:    "0s",
				IsPoweredOn:      false,
				PowerConsumption: 0,
			},
			results,
		)
	}
}

func TestDatagramOnMessage(t *testing.T) {
	msg := unhexMessage(t, workingMessageV4OnStateHex)
	parser := New(msg)

	results, err := parser.ToJSON()

	if assert.NoError(t, err) {
		assert.Equal(
			t,
			&DatagramParsedJSON{
				Name:             "Switcher_V4_B54A",
				IP:               "10.100.102.82",
				ID:               "000000",
				Key:              "06",
				MAC:              "00:00:00:00:00:00",
				TimeToShutdown:   "2h0m0s",
				TimeRemaining:    "14m47s",
				IsPoweredOn:      true,
				PowerConsumption: 2254,
			},
			results,
		)
	}
}
