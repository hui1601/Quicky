package utils

import "github.com/muka/go-bluetooth/bluez/profile/device"

type DeviceType int

const (
	// DeviceTypeUnknown is an unknown device type.
	DeviceTypeUnknown DeviceType = iota
	// DeviceTypeController is a controller.
	DeviceTypeController
	// DeviceTypeEarPhone is an earphone.
	DeviceTypeEarPhone
)

// GetDeviceType determines the device type.
func GetDeviceType(device1 device.Device1) DeviceType {
	manufacturerData, err := device1.GetManufacturerData()
	if manufacturerData == nil || err != nil {
		return DeviceTypeUnknown
	}
	if _, ok := manufacturerData[0x521c]; ok {
		return DeviceTypeController
	}
	serviceData, err := device1.GetServiceData()
	if err != nil {
		return DeviceTypeUnknown
	}
	if _, ok := serviceData["00007034-0000-1000-8000-00805f9b34fb"]; ok {
		return DeviceTypeEarPhone
	}
	return DeviceTypeUnknown
}
