package parse

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed working_v4_on_state_message.txt
var workingMessageV4OnState []byte

func TestDatagramMessage(t *testing.T) {
	parser := New(workingMessageV4OnState)
	assert.Equal(
		t,
		string(workingMessageV4OnState),
		parser.String(),
		"the two messages should be identical",
	)
}

func TestDatagramMessageNameParsing(t *testing.T) {
	parser := New(workingMessageV4OnState)
	assert.Equal(t, "Switcher_V4_B54A", parser.GetDeviceName())
}

func TestDatagramMessageKeyParsing(t *testing.T) {
	parser := New(workingMessageV4OnState)
	assert.Equal(t, "6c", parser.GetDeviceKey())
}

func TestDatagramMessageIPParsing(t *testing.T) {
	parser := New(workingMessageV4OnState)
	ip, err := parser.GetIPType1()
	if err != nil {
		assert.Errorf(t, err, "unexpected error")
	}

	assert.Equal(t, "75.194.149.194", ip.String())
}

func TestDatagramMessageMACParsing(t *testing.T) {
	parser := New(workingMessageV4OnState)
	mac, err := parser.GetDeviceMAC()
	if err != nil {
		assert.Errorf(t, err, "unexpected error")
	}

	assert.Equal(t, "89:c3:80:30:2a:c2", mac.String())
}

func TestDatagramMessageIsPoweredOnParsing(t *testing.T) {
	parser := New(workingMessageV4OnState)
	isOn := parser.IsPoweredOn()

	assert.True(t, isOn, "the device message should say powered on")
}
