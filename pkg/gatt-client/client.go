package gatt_client

import (
	"github.com/godbus/dbus/v5"
	"github.com/hui1601/Quicky/pkg/utils"
	"strings"
)

type Client struct {
	conn        *dbus.Conn
	path        string
	scanHandler func(address string)
}

func NewClientWithSystemBus() (*Client, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return NewClient(conn), nil
}

func NewClient(conn *dbus.Conn) *Client {
	// That's not a good idea to hardcode the path, but I don't know how to get it dynamically
	return &Client{conn: conn, path: "/org/bluez/hci0"}
}

func (c *Client) GetDeviceByAddress(address string) *Device {
	return NewDevice(c, strings.ToUpper(address))
}

func (c *Client) GetDeviceByPath(path dbus.ObjectPath) *Device {
	// change path to mac address(/org/bluez/hci0/dev_C4_AC_60_64_9B_67 -> C4:AC:60:64:9B:67)
	addr := utils.PathToMac(string(path))
	device := NewDevice(c, addr)
	if !device.isDeviceConnected() {
		return nil
	}
	return device
}
func (c *Client) IsDeviceAvailable(address string) bool {
	devices := c.KnownDevices()
	address = strings.ToUpper(address)
	for _, device := range devices {
		if device == address {
			return true
		}
	}
	return false
}

func (c *Client) AddScanHandler(handler func(address string)) {
	// InterfacesAdded signal
	err := c.conn.AddMatchSignal(dbus.WithMatchInterface("org.freedesktop.DBus.ObjectManager"), dbus.WithMatchMember("InterfacesAdded"))
	if err != nil {
		return
	}
	c.scanHandler = handler
	ch := make(chan *dbus.Signal, 10)
	c.conn.Signal(ch)
	go func() {
		for {
			signal := <-ch
			if signal.Name == "org.freedesktop.DBus.ObjectManager.InterfacesAdded" {
				// get address from path
				path := string(signal.Body[0].(dbus.ObjectPath))
				// remove prefix(dev_)
				if !strings.HasPrefix(path, c.path+"/dev_") {
					continue
				}
				address := strings.Replace(path, c.path+"/dev_", "", 1)
				if strings.Contains(address, "/") {
					continue
				}
				address = strings.ReplaceAll(address, "_", ":")
				if c.scanHandler != nil {
					c.scanHandler(address)
				}
			}
		}
	}()
}

func (c *Client) ResetScanFilter() {
	busObj := c.conn.Object("org.bluez", dbus.ObjectPath(c.path))
	filter := map[string]interface{}{
		"Transport":     "le",
		"UUIDs":         []string{},
		"DuplicateData": true,
	}
	busObj.Call("org.bluez.Adapter1.SetDiscoveryFilter", 0, filter)
}

func (c *Client) StartScan() {
	busObj := c.conn.Object("org.bluez", dbus.ObjectPath(c.path))
	busObj.Call("org.bluez.Adapter1.StartDiscovery", 0)
	if c.scanHandler == nil {
		return
	}
	devices := c.KnownDevices()
	for _, device := range devices {
		c.scanHandler(device)
	}
}

func (c *Client) StopScan() {
	busObj := c.conn.Object("org.bluez", dbus.ObjectPath(c.path))
	busObj.Call("org.bluez.Adapter1.StopDiscovery", 0)
}

func (c *Client) KnownDevices() []string {
	// org.freedesktop.DBus.ObjectManager
	busObj := c.conn.Object("org.bluez", "/")
	var devices map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	err := busObj.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&devices)
	if err != nil {
		return nil
	}
	var addresses []string
	for path := range devices {
		if strings.HasPrefix(string(path), c.path+"/dev_") {
			address := utils.PathToMac(string(path))
			if strings.Contains(address, "/") || address == "" {
				continue
			}
			addresses = append(addresses, address)
		}
	}
	return addresses
}
