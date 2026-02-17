package device

import (
	"context"
	"errors"
	"sync"

	"github.com/hui1601/Quicky/internal/command"
	"github.com/hui1601/Quicky/internal/constant"
	"github.com/hui1601/Quicky/internal/response"
	"tinygo.org/x/bluetooth"
)

type Client struct {
	Adapter *bluetooth.Adapter
	MAC     bluetooth.MAC
	Device  bluetooth.Device

	commandChar bluetooth.DeviceCharacteristic
	notifyChar  bluetooth.DeviceCharacteristic
	eqChar      bluetooth.DeviceCharacteristic
	keyFuncChar bluetooth.DeviceCharacteristic
	batteryChar bluetooth.DeviceCharacteristic
	versionChar bluetooth.DeviceCharacteristic

	events    chan response.Event
	mu        sync.Mutex
	connected bool
}

func NewClient(mac string) (*Client, error) {
	deviceMAC, err := bluetooth.ParseMAC(mac)
	if err != nil {
		return nil, err
	}
	return &Client{
		Adapter: bluetooth.DefaultAdapter,
		MAC:     deviceMAC,
		events:  make(chan response.Event, 32),
	}, nil
}

func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return errors.New("already connected")
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if err := c.Adapter.Enable(); err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	device, err := c.Adapter.Connect(
		bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: c.MAC}},
		bluetooth.ConnectionParams{},
	)
	if err != nil {
		return err
	}
	c.Device = device

	if err := ctx.Err(); err != nil {
		_ = device.Disconnect()
		return err
	}

	if err := c.discoverCharacteristics(); err != nil {
		_ = device.Disconnect()
		return err
	}

	if err := c.subscribeNotifications(); err != nil {
		_ = device.Disconnect()
		return err
	}

	c.connected = true
	return nil
}

func (c *Client) discoverCharacteristics() error {
	services, err := c.Device.DiscoverServices([]bluetooth.UUID{constant.ServiceUUID})
	if err != nil {
		return err
	}
	if len(services) == 0 {
		return errors.New("QCY service not found")
	}

	charUUIDs := []bluetooth.UUID{
		constant.CommandUUID,
		constant.NotifyUUID,
		constant.EQUUID,
		constant.KeyFuncUUID,
		constant.BatteryUUID,
		constant.VersionUUID,
	}

	chars, err := services[0].DiscoverCharacteristics(charUUIDs)
	if err != nil {
		return err
	}
	if len(chars) < 6 {
		return errors.New("not all characteristics found")
	}

	c.commandChar = chars[0]
	c.notifyChar = chars[1]
	c.eqChar = chars[2]
	c.keyFuncChar = chars[3]
	c.batteryChar = chars[4]
	c.versionChar = chars[5]

	return nil
}

func (c *Client) subscribeNotifications() error {
	return c.notifyChar.EnableNotifications(func(buf []byte) {
		commands, err := command.ParsePacket(buf)
		if err != nil {
			ev := response.Event{
				Type:  response.EventUnknown,
				Raw:   buf,
				Error: err,
			}
			select {
			case c.events <- ev:
			default:
			}
			return
		}
		for _, cmd := range commands {
			ev := response.Dispatch(cmd.OperationCode, cmd.Parameters)
			select {
			case c.events <- ev:
			default:
			}
		}
	})
}

func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	err := c.Device.Disconnect()
	c.connected = false
	return err
}

func (c *Client) Events() <-chan response.Event {
	return c.events
}

func (c *Client) SendCommand(cmd *command.Command) error {
	c.mu.Lock()
	if !c.connected {
		c.mu.Unlock()
		return errors.New("not connected")
	}
	c.mu.Unlock()

	_, err := c.commandChar.WriteWithoutResponse(cmd.PackPacket())
	return err
}

func (c *Client) WriteEQ(data []byte) error {
	c.mu.Lock()
	if !c.connected {
		c.mu.Unlock()
		return errors.New("not connected")
	}
	c.mu.Unlock()

	_, err := c.eqChar.WriteWithoutResponse(data)
	return err
}

func (c *Client) WriteKeyFunction(data []byte) error {
	c.mu.Lock()
	if !c.connected {
		c.mu.Unlock()
		return errors.New("not connected")
	}
	c.mu.Unlock()

	_, err := c.keyFuncChar.WriteWithoutResponse(data)
	return err
}

func (c *Client) ReadBattery() (response.Battery, error) {
	c.mu.Lock()
	if !c.connected {
		c.mu.Unlock()
		return response.Battery{}, errors.New("not connected")
	}
	c.mu.Unlock()

	var buf [3]byte
	n, err := c.batteryChar.Read(buf[:])
	if err != nil {
		return response.Battery{}, err
	}
	return response.ParseBattery(buf[:n])
}

func (c *Client) ReadVersion() (response.Version, error) {
	c.mu.Lock()
	if !c.connected {
		c.mu.Unlock()
		return response.Version{}, errors.New("not connected")
	}
	c.mu.Unlock()

	var buf [6]byte
	n, err := c.versionChar.Read(buf[:])
	if err != nil {
		return response.Version{}, err
	}
	return response.ParseVersion(buf[:n])
}

func (c *Client) Connected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}
