package gatt_client

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/hui1601/Quicky/pkg/utils"
)

type EarphoneInfo struct {
	updated         bool
	ClassMac        string
	OtherMac        string
	Vendor          uint32
	IsBlack         bool
	IsLeftCharging  bool
	IsRightCharging bool
	IsCaseCharging  bool
	LeftBattery     uint8
	RightBattery    uint8
	CaseBattery     uint8
}

var earphoneInfo EarphoneInfo

func (d *Device) EarphoneInfo() (EarphoneInfo, error) {
	if earphoneInfo.updated {
		return earphoneInfo, nil
	}
	var manufacturer map[uint16]dbus.Variant
	busObj := d.client.conn.Object("org.bluez", dbus.ObjectPath(d.client.path+"/dev_"+utils.MacToPath(d.address)))
	err := busObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.bluez.Device1", "ManufacturerData").Store(&manufacturer)
	if err != nil {
		return EarphoneInfo{}, err
	}
	if _, ok := manufacturer[0x521c]; ok {
		earphoneInfo = EarphoneInfo{}
		b := manufacturer[0x521c].Value().([]byte)
		earphoneInfo.ClassMac = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", b[12], b[11], b[13], b[16], b[15], b[14])
		earphoneInfo.OtherMac = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", b[19], b[18], b[20], b[23], b[22], b[21])
		earphoneInfo.Vendor = (uint32(b[0]) << 8) + uint32(b[1])
		earphoneInfo.IsBlack = (b[3] & 2) == 2
		earphoneInfo.IsLeftCharging = b[5] == 0
		earphoneInfo.IsRightCharging = b[6] == 0
		// case charging is not supported by QCY HT07 APP
		earphoneInfo.IsCaseCharging = b[7] == 0
		earphoneInfo.LeftBattery = b[5] & 100
		earphoneInfo.RightBattery = b[6] & 100
		earphoneInfo.CaseBattery = b[7] & 100
		earphoneInfo.updated = true
		return earphoneInfo, nil
	}
	return EarphoneInfo{}, nil
}

func (d *Device) readCharacteristic(uuid string) ([]byte, error) {
	// get characteristic path
	characteristicPath, ok := d.CharacteristicUUIDs()[uuid]
	if !ok {
		return nil, nil
	}
	// read characteristic
	busObj := d.client.conn.Object("org.bluez", dbus.ObjectPath(characteristicPath))
	var value []byte
	err := busObj.Call("org.bluez.GattCharacteristic1.ReadValue", 0, map[string]dbus.Interface{}).Store(&value)
	return value, err
}

func (d *Device) writeCharacteristic(uuid string, value []byte) error {
	// get characteristic path
	characteristicPath, ok := d.CharacteristicUUIDs()[uuid]
	if !ok {
		return nil
	}
	// read characteristic
	busObj := d.client.conn.Object("org.bluez", dbus.ObjectPath(characteristicPath))
	arg := map[string]interface{}{
		"type": "command",
	}
	call := busObj.Call("org.bluez.GattCharacteristic1.WriteValue", 0, value, arg)
	<-call.Done
	err := call.Err
	return err
}

func (d *Device) writeCommand(cmd int, args []byte) error {
	var cmdBytes []byte
	cmdBytes = append(cmdBytes, 0xff)
	cmdBytes = append(cmdBytes, byte(len(args)+2))
	cmdBytes = append(cmdBytes, byte(cmd))
	cmdBytes = append(cmdBytes, byte(len(args)))
	cmdBytes = append(cmdBytes, args...)
	return d.writeCharacteristic("00001001-0000-1000-8000-00805f9b34fb", cmdBytes)
}

func (d *Device) GetBatteryLevel() EarphoneInfo {
	if !earphoneInfo.updated {
		earphoneInfo, _ = d.EarphoneInfo()
	}
	// refresh battery
	v, err := d.readCharacteristic("00000008-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return earphoneInfo
	}
	if len(v) < 3 {
		return earphoneInfo
	}
	earphoneInfo.LeftBattery = v[0] & 100
	earphoneInfo.RightBattery = v[1] & 100
	earphoneInfo.CaseBattery = v[2] & 100
	return earphoneInfo
}

// noise-canceling enum
// this value must not change
//
//goland:noinspection ALL
const (
	ModeSilent = iota
	ModeWork
	ModeNoise
	ModePassThrough
	ModeOff
)

// SetNoiseCanceling mode(0: silent, 1: work, 2: noise, 3: pass-through, 4: off), level(noise:0~2, pass:0~6(6 is voice enhance))
func (d *Device) SetNoiseCanceling(mode int, level int) {
	b := byte(0)
	if mode <= ModeNoise {
		b = (byte(mode)+4)<<4 + byte(level%3+1)
	} else if mode == ModePassThrough {
		if level != 6 {
			b = 0xa<<4 + byte(level%6+1)
		} else {
			b = 0xa<<4 + byte(0)
		}
	} else {
		b = 0x02
	}
	err := d.writeCommand(0x0C, []byte{b})
	if err != nil {
		panic(err)
	}
}
