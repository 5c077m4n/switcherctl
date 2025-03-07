// Package consts for storing util consts
package consts

import (
	"net"
)

// Type 1 devices: Heaters (v2, touch, v4, Heater), Plug
// Type 2 devices: Breeze, Runners

type (
	Port           uint
	DeviceCategory uint
	DeviceType     uint
)

// Port types and values
const (
	UDPPortType1    Port = 20_002
	UDPPortType1New Port = 10_002
	UDPPortType2    Port = 20_003
	UDPPortType2New Port = 10_003
	TCPPortType1    Port = 9_957
	TCPPortType2    Port = 10_000
)

// Device catgories
const (
	DeviceCategoryWaterHeater DeviceCategory = iota
	DeviceCategoryPowerPlug
	DeviceCategoryThermostat
	DeviceCategoryShutter
)

// Device types
const (
	DeviceTypeMini DeviceType = iota
	DeviceTypePowerPlug
	DeviceTypeTouch
	DeviceTypeV2ESP
	DeviceTypeV2QCA
	DeviceTypeV4
	DeviceTypeBreeze
	DeviceTypeRunner
	DeviceTypeRunnerMini
)

// Message lengths for different devices
const (
	MessageLengthDefault = 165
	MessageLengthBreeze  = 168
	MessageLengthRunner  = 159
)

// Communication packet templates
const (
	TemplatePacketLoginType1   = "fef052000232a10000000000340001000000000000000000%s00000000000000000000f0fe%s00000000000000000000000000000000000000000000000000000000000000000000000000"
	TemplatePacketGetSchedules = "fef0570002320102%s340001000000000000000000%s00000000000000000000f0fe%s00000000000000000000000000000000000000000000000000000000000000000000000000060000"
)

// DefaultIP the fallback IP
var DefaultIP = net.IP{10, 100, 102, 82}

// Devices join maps
var (
	DeviceCategoryToUDPPort = map[DeviceCategory]Port{
		DeviceCategoryWaterHeater: UDPPortType1New,
		DeviceCategoryPowerPlug:   UDPPortType1,
		DeviceCategoryThermostat:  UDPPortType2,
		DeviceCategoryShutter:     UDPPortType2,
	}
	DeviceCategoryToTCPPort = map[DeviceCategory]Port{
		DeviceCategoryWaterHeater: TCPPortType1,
		DeviceCategoryPowerPlug:   TCPPortType1,
		DeviceCategoryThermostat:  TCPPortType2,
		DeviceCategoryShutter:     TCPPortType2,
	}
	DeviceTypeToCategory = map[DeviceType]DeviceCategory{
		DeviceTypeMini:       DeviceCategoryWaterHeater,
		DeviceTypePowerPlug:  DeviceCategoryPowerPlug,
		DeviceTypeTouch:      DeviceCategoryWaterHeater,
		DeviceTypeV2ESP:      DeviceCategoryWaterHeater,
		DeviceTypeV2QCA:      DeviceCategoryWaterHeater,
		DeviceTypeV4:         DeviceCategoryWaterHeater,
		DeviceTypeBreeze:     DeviceCategoryThermostat,
		DeviceTypeRunner:     DeviceCategoryShutter,
		DeviceTypeRunnerMini: DeviceCategoryShutter,
	}
)
