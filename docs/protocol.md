# Protocol
QCY uses GATT for device control.
Please refer to the [GATT service uuid](service.md).

## Connect
QCY earphones have three Bluetooth MAC addresses.

Control, left, right, in most cases you can see either side and the control MAC address.

The control MAC address is broadcast after connecting to either side and receiving the L2CAP Disconnect command. 

However, it also worked if I just connected to the device and then disconnected the device (that's a bad idea, but I couldn't find a Dbus API to disconnect the L2CAP connection in Bluez)

## Command

| 1byte | 1byte  | ...  |
|-------|--------|------|
| cmd   | length | data |

cmd is 1 byte, length is 1 byte, and data is length bytes.

command is sent with service uuid `00001001-0000-1000-8000-00805f9b34fb`.

### Command List
| cmd  | description     | data length     | data                                                      |
|------|-----------------|-----------------|-----------------------------------------------------------|
| 0x01 | Reset Default   | 0               |                                                           |
| 0x02 | Reset Pairing   | 0               |                                                           |
| 0x03 | Factory Reset   | 0               |                                                           |
| 0x04 | Music Control   | 1               | unknown byte                                              |
| 0x05 | Light Flash     | 1               | 0x00(OFF)/0x01(ON)                                        |
| 0x06 | In-Ear Test     | 1               | 0x02(OFF)/0x01(ON)                                        |
| 0x07 | Unknown         | -               | -                                                         |
| 0x08 | Unknown         | -               | -                                                         |
| 0x09 | Low Latency     | 1               | 0x02(OFF)/0x01(ON)                                        |
| 0x0a | Unknown         | -               | -                                                         |
| 0x0b | Unknown         | -               | -                                                         |
| 0x0c | Noise Cancel    | 1               | Mode(High/Low) See below                                  |
| 0x0d | Test Mode       | 1               | 0x02(OFF)/0x01(ON)                                        |
| 0x0e | Unknown         | -               | -                                                         |
| 0x0f | Unknown         | -               | -  1                                                      |
| 0x10 | Unknown         | -               | -                                                         |
| 0x11 | Unknown         | -               | -                                                         |
| 0x12 | Unknown         | -               | -                                                         |
| 0x13 | Unknown         | -               | -                                                         |
| 0x14 | Scedule Reboot  | 4               | Scheduled reboot time(High, Low), Current time(High, Low) |
| 0x15 | Unknown         | -               | -                                                         |
| 0x16 | Sound Balance   | 1               | 0x00(Left)/0x64(Right)                                    |
| 0x17 | Unknown         | -               | -                                                         |
| 0x18 | Rename Earphone | (String length) | String                                                    |
| 0x19 | Audio Language  | (String length) | String                                                    |

#### Noise Cancel Mode
| Mode | Description                           |
|------|---------------------------------------|
| 0x00 | OFF                                   |
| 0x1? | Unknown                               |
| 0x21 | Silent Environment ANC Mode(Level 1)  |
| 0x22 | Silent Environment ANC Mode(Level 2)  |
| 0x23 | Silent Environment ANC Mode(Level 3)  |
| 0x31 | Working Environment ANC Mode(Level 1) |
| 0x32 | Working Environment ANC Mode(Level 2) |
| 0x33 | Working Environment ANC Mode(Level 3) |
| 0x41 | Noise Environment ANC Mode(Level 1)   |
| 0x42 | Noise Environment ANC Mode(Level 2)   |
| 0x43 | Noise Environment ANC Mode(Level 3)   |
| 0xa1 | Transparency Mode(Level 1)            |
| 0xa2 | Transparency Mode(Level 2)            |
| 0xa3 | Transparency Mode(Level 3)            |
| 0xa4 | Transparency Mode(Level 4)            |
| 0xa5 | Transparency Mode(Level 5)            |
| 0xa6 | Transparency Mode(Level 6)            |
| 0xa7 | Transparency Mode(Increase Voice)     |