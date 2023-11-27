package gatt

import (
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

type Client struct {
	Adapter *adapter.Adapter1
	Device  *device.Device1
}

func NewClient(adapterID string, macAddr string) (*Client, error) {
	adapter1, err := adapter.GetAdapter(adapterID)
	if err != nil {
		return nil, err
	}

	device1, err := adapter1.GetDeviceByAddress(macAddr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Adapter: adapter1,
		Device:  device1,
	}, nil
}

func (c *Client) Connect() error {
	err := c.Device.Connect()
	if err != nil {
		return err
	}
	err = c.Device.Pair()
	if err != nil {
		return err
	}
	err = c.Device.SetTrusted(true)
	if err != nil {
		return err
	}
	return nil
}
