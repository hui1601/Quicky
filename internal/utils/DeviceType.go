package utils

type DeviceType int

const (
	// DeviceTypeUnknown is an unknown device type.
	DeviceTypeUnknown DeviceType = iota
	// DeviceTypeController is a controller.
	DeviceTypeController
	// DeviceTypeEarPhone is an earphone.
	DeviceTypeEarPhone
)
