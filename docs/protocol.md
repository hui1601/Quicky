# Protocol
QCY uses GATT for device control.
Please refer to the [GATT service uuid](service.md).

## Connect
QCY earphones have three Bluetooth MAC addresses.

Control, left, right, in most cases you can see either side and the control MAC address.

The control MAC address is broadcast after connecting to either side and receiving the L2CAP Disconnect command. 

However, it also worked if I just connected to the device and then disconnected the device (that's a bad idea, but I couldn't find a Dbus API to disconnect the L2CAP connection in Bluez)

## Cmd

