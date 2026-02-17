Characteristic UUIDs found in QCY earphones.

'Std' means the UUID is defined in the Bluetooth standard.

> **Source**: Observed on QCY HT07, supplemented with BLE traffic analysis.

## QCY Main Service

| uuid                                   | description                       | mode       |
|----------------------------------------|-----------------------------------|------------|
| `0000a001-0000-1000-8000-00805f9b34fb` | Main BLE Service                  | service    |

## QCY Custom Characteristics

| uuid                                   | description                       | mode       |
|----------------------------------------|-----------------------------------|------------|
| `00000001-0000-1000-8000-00805f9b34fb` | Left Single Tap Key (V1)          | read/write |
| `00000002-0000-1000-8000-00805f9b34fb` | Right Single Tap Key (V1)         | read/write |
| `00000003-0000-1000-8000-00805f9b34fb` | Left Double Tap Key (V1)          | read/write |
| `00000004-0000-1000-8000-00805f9b34fb` | Right Double Tap Key (V1)         | read/write |
| `00000005-0000-1000-8000-00805f9b34fb` | Left Triple Tap Key (V1)          | read/write |
| `00000006-0000-1000-8000-00805f9b34fb` | Right Triple Tap Key (V1)         | read/write |
| `00000007-0000-1000-8000-00805f9b34fb` | Version                           | read       |
| `00000008-0000-1000-8000-00805f9b34fb` | Battery (Left/Right/Case)         | read       |
| `00000009-0000-1000-8000-00805f9b34fb` | Current Language                  | read       |
| `0000000a-0000-1000-8000-00805f9b34fb` | Reset (V1)                        | write      |
| `0000000b-0000-1000-8000-00805f9b34fb` | EQ (direct write, no 0xFF frame)  | read/write |
| `0000000c-0000-1000-8000-00805f9b34fb` | Send Time (V1)                    | write      |
| `0000000d-0000-1000-8000-00805f9b34fb` | Key Function V2 (direct write)    | read/write |
| `0000000e-0000-1000-8000-00805f9b34fb` | ZR Device Settings                | read/write |
| `0000000f-0000-1000-8000-00805f9b34fb` | In-Ear Check (JL devices)         | read       |
| `00001001-0000-1000-8000-00805f9b34fb` | Command Write (main protocol)     | write      |
| `00001002-0000-1000-8000-00805f9b34fb` | Settings Read / Notify            | read/notify|
| `00001003-0000-1000-8000-00805f9b34fb` | Unknown                           | unknown    |

## Bluetooth Standard Characteristics

| uuid                                   | description                       | mode       |
|----------------------------------------|-----------------------------------|------------|
| `00002902-0000-1000-8000-00805f9b34fb` | CCCD (Client Config Descriptor)   | config     |
| `00002a05-0000-1000-8000-00805f9b34fb` | Service Changed Std               | indicate   |
| `00002a19-0000-1000-8000-00805f9b34fb` | Battery Level Std                 | read       |
| `00002a24-0000-1000-8000-00805f9b34fb` | Model Number Std                  | read       |
| `00002a25-0000-1000-8000-00805f9b34fb` | Serial Number Std                 | read       |
| `00002a26-0000-1000-8000-00805f9b34fb` | Firmware Revision Std             | read       |
| `00002a27-0000-1000-8000-00805f9b34fb` | Hardware Revision Std             | read       |
| `00002a28-0000-1000-8000-00805f9b34fb` | Software Revision Std             | read       |
| `00002a29-0000-1000-8000-00805f9b34fb` | Manufacturer Name Std             | read       |
| `00002a50-0000-1000-8000-00805f9b34fb` | PnP ID Std                        | read       |

## Unknown / Vendor-Specific

| uuid                                   | description                       | mode       |
|----------------------------------------|-----------------------------------|------------|
| `00002001-0000-1000-8000-00805f9b34fb` | Unknown                           | unknown    |
| `00002002-0000-1000-8000-00805f9b34fb` | Unknown                           | unknown    |
| `0000ff01-0000-1000-8000-00805f9b34fb` | Unknown                           | unknown    |
| `0000ff02-0000-1000-8000-00805f9b34fb` | Unknown                           | unknown    |

## Notes

- **`00001001`** is the main command channel. All commands (except EQ and key functions) are sent here wrapped in `[0xFF, bodyLen, ...]` framing. See [protocol.md](protocol.md).
- **`00001002`** receives notifications with the same framing format. Subscribe to this for device state updates.
- **`0000000B`** (EQ) and **`0000000D`** (Key Function) use **direct writes** without the `0xFF` packet framing.
- **`00000007`** and **`00000008`** can be read directly for version and battery, or received via command 0x30/0x2F notifications on `00001002`.
- UUIDs `00000001`–`00000006` are V1 legacy key mapping characteristics (one per gesture). Modern firmware uses `0000000D` instead.
- UUID `0000000E` is used by ZR-chipset devices for device-specific settings.