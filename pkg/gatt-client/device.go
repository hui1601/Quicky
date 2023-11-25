package gatt_client

import (
	"github.com/godbus/dbus/v5"
	"github.com/hui1601/Quicky/pkg/utils"
	"strings"
)

type Device struct {
	address string
	client  *Client
}

func NewDevice(client *Client, address string) *Device {
	return &Device{client: client, address: strings.ToUpper(address)}
}

func (d *Device) isDeviceConnected() bool {
	var connected bool
	busObj := d.client.conn.Object("org.bluez", dbus.ObjectPath(d.client.path+"/dev_"+strings.ReplaceAll(d.address, ":", "_")))
	err := busObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.bluez.Device1", "Connected").Store(&connected)
	if err != nil {
		return false
	}
	return connected
}

func (d *Device) IsQcyDevice() bool {
	// get manufacturer
	var manufacturer map[uint16]dbus.Variant
	busObj := d.client.conn.Object("org.bluez", dbus.ObjectPath(d.client.path+"/dev_"+utils.MacToPath(d.address)))
	err := busObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.bluez.Device1", "ManufacturerData").Store(&manufacturer)
	if err != nil {
		//println(err.Error())
		return false
	}
	if _, ok := manufacturer[0x521c]; ok {
		return true
	}
	return false
}

func (d *Device) Connect() error {
	// already connected
	if d.isDeviceConnected() {
		return nil
	}
	busObj := d.client.conn.Object("org.bluez", dbus.ObjectPath(d.client.path+"/dev_"+strings.ReplaceAll(d.address, ":", "_")))
	res := busObj.Call("org.bluez.Device1.Connect", 0)
	return res.Err
}

func (d *Device) Address() string {
	return d.address
}

func (d *Device) Disconnect() error {
	// already disconnected
	if !d.isDeviceConnected() {
		return nil
	}
	busObj := d.client.conn.Object("org.bluez", dbus.ObjectPath(d.client.path+"/dev_"+strings.ReplaceAll(d.address, ":", "_")))
	res := busObj.Call("org.bluez.Device1.Disconnect", 0)
	return res.Err
}

func (d *Device) IsConnected() bool {
	return d.isDeviceConnected()
}

func (d *Device) ServiceUUIDs() []string {
	var uuids []string
	busObj := d.client.conn.Object("org.bluez", dbus.ObjectPath(d.client.path+"/dev_"+strings.ReplaceAll(d.address, ":", "_")))
	err := busObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.bluez.Device1", "UUIDs").Store(&uuids)
	if err != nil {
		return nil
	}
	return uuids
}

func (d *Device) CharacteristicUUIDs() map[string]string {
	var uuids map[string]string
	var response map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	busObj := d.client.conn.Object("org.bluez", "/")
	err := busObj.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&response)
	if err != nil {
		return nil
	}
	uuids = make(map[string]string)
	for path := range response {
		//println(string(path))
		if strings.Contains(string(path), "char") && !strings.Contains(string(path), "desc") {
			// read characteristic
			var characteristic string
			busObj := d.client.conn.Object("org.bluez", path)
			err := busObj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.bluez.GattCharacteristic1", "UUID").Store(&characteristic)
			if err != nil {
				continue
			}
			uuids[characteristic] = string(path)
		}
	}
	return uuids
}
