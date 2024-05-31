// Package consts for storing util consts
package consts

// Type 1 devices: Heaters (v2, touch, v4, Heater), Plug
// Type 2 devices: Breeze, Runners

const (
	UDPPortType1    = 20_002
	UDPPortType1New = 10_002
	UDPPortType2    = 20_003
	UDPPortType2New = 10_003
	TCPPortType1    = 9_957
	TCPPortType2    = 10_000
)

const (
	DeviceCategoryWaterHeater = iota
	DeviceCategoryPowerPlug
	DeviceCategoryThermostat
	DeviceCategoryShutter
)

const (
	DeviceTypeMini = iota
	DeviceTypePowerPlug
	DeviceTypeTouch
	DeviceTypeV2ESP
	DeviceTypeV2QCA
	DeviceTypeV4
	DeviceTypeBreeze
	DeviceTypeRunner
	DeviceTypeRunnerMini
)

// DefaultIP the fallback IP
const (
	DefaultIP = "10.100.102.82"
)

var (
	DeviceCategoryToUDPPort = map[int]int{
		DeviceCategoryWaterHeater: UDPPortType1New,
		DeviceCategoryPowerPlug:   UDPPortType1,
		DeviceCategoryThermostat:  UDPPortType2,
		DeviceCategoryShutter:     UDPPortType2,
	}
	DeviceCategoryToTCPPort = map[int]int{
		DeviceCategoryWaterHeater: TCPPortType1,
		DeviceCategoryPowerPlug:   TCPPortType1,
		DeviceCategoryThermostat:  TCPPortType2,
		DeviceCategoryShutter:     TCPPortType2,
	}
	DeviceTypeToCategory = map[int]int{
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
