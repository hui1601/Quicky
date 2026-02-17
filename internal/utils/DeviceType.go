package utils

import (
	"fmt"

	"github.com/hui1601/Quicky/internal/constant"
	"tinygo.org/x/bluetooth"
)

type DeviceType int

const (
	DeviceTypeUnknown DeviceType = iota
	DeviceTypeQCY
)

type AdvertisementInfo struct {
	VendorID        uint16
	ColorIndex      byte
	LeftBattery     byte
	RightBattery    byte
	BoxBattery      byte
	IsLeftCharging  bool
	IsRightCharging bool
	IsBoxCharging   bool
	ControlMAC      string
	OtherMAC        string
}

func GetDeviceType(advertise bluetooth.AdvertisementFields) DeviceType {
	manufacturerData := advertise.ManufacturerData
	if len(manufacturerData) == 0 {
		return DeviceTypeUnknown
	}
	for _, data := range manufacturerData {
		if data.CompanyID == constant.QCYCompanyID {
			return DeviceTypeQCY
		}
	}
	return DeviceTypeUnknown
}

// ParseManufacturerData parses QCY manufacturer data payload.
// Requires at least 20 bytes; 24 bytes for full MAC parsing.
func ParseManufacturerData(data []byte) (*AdvertisementInfo, error) {
	if len(data) < 20 {
		return nil, fmt.Errorf("manufacturer data too short: got %d bytes, need at least 20", len(data))
	}

	info := &AdvertisementInfo{}

	info.VendorID = uint16(data[0])<<8 | uint16(data[1])
	info.ColorIndex = (data[3] & 0x18) >> 1

	info.LeftBattery = data[5] & 0x7F
	info.IsLeftCharging = data[5]&0x80 != 0

	info.RightBattery = data[6] & 0x7F
	info.IsRightCharging = data[6]&0x80 != 0

	info.BoxBattery = data[7] & 0x7F
	info.IsBoxCharging = data[7]&0x80 != 0

	if len(data) >= 17 {
		info.ControlMAC = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
			data[12], data[11], data[13], data[16], data[15], data[14])
	}

	if len(data) >= 24 {
		info.OtherMAC = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
			data[19], data[18], data[20], data[23], data[22], data[21])
		if info.OtherMAC == "00:00:00:00:00:00" {
			info.OtherMAC = info.ControlMAC
		}
	}

	return info, nil
}
