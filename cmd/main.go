package main

import (
	gattclient "github.com/hui1601/Quicky/pkg/gatt-client"
)

var client *gattclient.Client
var foundDevice chan *gattclient.Device

func main() {
	var err error
	var device *gattclient.Device
	device = gattclient.NewDevice(client, "")
	client, err = gattclient.NewClientWithSystemBus()
	if err != nil {
		panic(err)
	}
	println("client created")
	client.StopScan()
	println("scan stopped")
	knownDevices := client.KnownDevices()
	for _, v := range knownDevices {
		tmpDevice := client.GetDeviceByAddress(v)
		if tmpDevice.IsConnected() && tmpDevice.IsQcyDevice() {
			println("found connected device: " + v)
			device = tmpDevice
		}
	}
	if device.Address() == "" {
		foundDevice = make(chan *gattclient.Device)
		client.ResetScanFilter()
		println("scan filter reset")
		client.AddScanHandler(scanHandler)
		println("scan handler added")
		go client.StartScan()
		println("scan started\nWaiting for device...")
		// sleep until device is found
		device = <-foundDevice
		println("device found")
		client.StopScan()
		println("Connecting to device...")
		err = device.Connect()
		if err != nil {
			panic(err)
		}
		println("connected to device")
		if !device.IsConnected() {
			println("device not connected")
			return
		}
	}
	deviceInfo, err := device.EarphoneInfo()
	if err != nil {
		panic(err)
	}
	if deviceInfo.ClassMac != "" {
		println("found main device MAC: " + deviceInfo.ClassMac)
		err := gattclient.NewDevice(client, deviceInfo.ClassMac).Connect()
		if err != nil {
			println("failed to connect to main device\ntrying to connect to other device")
			if deviceInfo.OtherMac != "" {
				println("found other device MAC: " + deviceInfo.OtherMac)
				err := gattclient.NewDevice(client, deviceInfo.OtherMac).Connect()
				if err != nil {
					panic(err)
				}
			}
		}
	}
	println("device connected!")
	// get battery level
	info := device.GetBatteryLevel()
	println("left battery level: ", info.LeftBattery)
	println("right battery level: ", info.RightBattery)
	println("case battery level: ", info.CaseBattery)
	device.SetNoiseCanceling(gattclient.ModePassThrough, 5)
}

func scanHandler(address string) {
	//println(address)
	// check if device is QCY
	device := client.GetDeviceByAddress(address)
	if !device.IsQcyDevice() {
		// check service uuids including 00007034-0000-1000-8000-00805f9b34fb
		if device.IsConnected() {
			for _, v := range device.ServiceUUIDs() {
				if v == "00007034-0000-1000-8000-00805f9b34fb" {
					println("found QCY earphone… disconnect it to connect QCY control service")
					_ = device.Disconnect()
					return
				}
			}
		}
		return
	}
	println("found device: " + address)
	println("QCY device found")
	foundDevice <- device
}
