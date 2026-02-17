package constant

import "tinygo.org/x/bluetooth"

func init() {
	ServiceUUID, _ = bluetooth.ParseUUID("0000a001-0000-1000-8000-00805f9b34fb")
	CommandUUID, _ = bluetooth.ParseUUID("00001001-0000-1000-8000-00805f9b34fb")
	NotifyUUID, _ = bluetooth.ParseUUID("00001002-0000-1000-8000-00805f9b34fb")
	EQUUID, _ = bluetooth.ParseUUID("0000000b-0000-1000-8000-00805f9b34fb")
	KeyFuncUUID, _ = bluetooth.ParseUUID("0000000d-0000-1000-8000-00805f9b34fb")
	BatteryUUID, _ = bluetooth.ParseUUID("00000008-0000-1000-8000-00805f9b34fb")
	VersionUUID, _ = bluetooth.ParseUUID("00000007-0000-1000-8000-00805f9b34fb")
	CCCDUUID, _ = bluetooth.ParseUUID("00002902-0000-1000-8000-00805f9b34fb")
}

const (
	// QCYCompanyID is the BLE manufacturer data company ID for QCY devices.
	QCYCompanyID uint16 = 0x521c
)

var (
	// ServiceUUID is the QCY main BLE service UUID.
	ServiceUUID bluetooth.UUID
	// CommandUUID is the command write characteristic UUID.
	CommandUUID bluetooth.UUID
	// NotifyUUID is the notification/settings read characteristic UUID.
	NotifyUUID bluetooth.UUID
	// EQUUID is the EQ direct write characteristic UUID (no 0xFF framing).
	EQUUID bluetooth.UUID
	// KeyFuncUUID is the key function direct write characteristic UUID (no 0xFF framing).
	KeyFuncUUID bluetooth.UUID
	// BatteryUUID is the battery read characteristic UUID.
	BatteryUUID bluetooth.UUID
	// VersionUUID is the version read characteristic UUID.
	VersionUUID bluetooth.UUID
	// CCCDUUID is the Client Characteristic Configuration Descriptor UUID.
	CCCDUUID bluetooth.UUID
)
