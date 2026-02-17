# Protocol

QCY uses BLE GATT for device control. Please refer to the [GATT service UUID list](service.md).

> **Source**: Reverse-engineered from QCY earphone BLE traffic analysis.

## Connection

QCY earphones expose three Bluetooth MAC addresses: **Control**, **Left**, and **Right**.

In most cases you can see either side and the control MAC address.

The control MAC address is broadcast after connecting to either side and receiving the L2CAP Disconnect command.

However, it also works if you just connect to the device and then disconnect (not ideal, but a workaround when no direct L2CAP disconnect API is available).

## Packet Format

All commands sent via the main command characteristic (`00001001`) and all responses received via the notification characteristic (`00001002`) share the same framing:

```
[0xFF] [body_len] [cmd1] [param_len1] [params1...] [cmd2] [param_len2] [params2...] ...
```

| Offset | Size   | Description                                              |
|--------|--------|----------------------------------------------------------|
| 0      | 1 byte | Start of Frame: `0xFF` (signed: `-1`)                    |
| 1      | 1 byte | Body length: total packet length minus 2                 |
| 2+     | varies | One or more command blocks (see below)                   |

Each **command block** within the body:

| Offset | Size   | Description                                      |
|--------|--------|--------------------------------------------------|
| 0      | 1 byte | Command ID                                       |
| 1      | 1 byte | Parameter count (number of parameter bytes)      |
| 2+     | N bytes| Parameter data                                   |

Multiple commands can be packed into a single packet. The parser iterates through the body using `param_len` to find the next command block.

### Exceptions

The following are written **directly** to their own characteristics (no `0xFF` framing):

| Data         | Characteristic UUID | Format                                    |
|--------------|--------------------|--------------------------------------------|
| EQ settings  | `0000000B`         | Raw EQ bytes (see [EQ](#eq-settings))      |
| Key functions| `0000000D`         | Raw key-function pairs (see [Key Functions](#key-functions)) |

---

## Command List

Commands are sent to `00001001-0000-1000-8000-00805f9b34fb`. Responses/notifications arrive on `00001002-0000-1000-8000-00805f9b34fb`.

### Overview

| Cmd  | Name               | Param Len | Description                        |
|------|--------------------|-----------|------------------------------------|
| 0x01 | Reset Default      | 0         | Reset settings to default          |
| 0x02 | Clear Pairing      | 0         | Clear Bluetooth pairing info       |
| 0x03 | Factory Reset      | 0         | Full factory reset                 |
| 0x04 | Music Control      | 1         | Play/pause/prev/next               |
| 0x05 | Light Flash        | 1         | LED flash on/off                   |
| 0x06 | In-Ear Detection   | 1         | Enable/disable in-ear detection    |
| 0x07 | Noise Value        | 1         | Set/read noise value               |
| 0x08 | Volume             | 3         | Set left/right volume              |
| 0x09 | Low Latency        | 1         | Game mode on/off                   |
| 0x0A | Monitoring         | 1         | Monitoring/sidetone value          |
| 0x0C | Noise Cancel Mode  | 1         | ANC/transparency/off mode          |
| 0x0D | Test Mode          | 1         | Enable/disable test mode           |
| 0x10 | Sleep Mode         | 1         | Enable/disable sleep mode          |
| 0x11 | Ear Tip Fit Test   | 2         | Start/stop ear tip fit test        |
| 0x12 | LED Mode           | 1         | LED indicator on/off               |
| 0x14 | Power Manager      | 4         | Auto power-off timer               |
| 0x16 | Sound Balance      | 1         | Left/right balance                 |
| 0x17 | ANC Setting        | 3         | ANC scene/level/noise setting      |
| 0x18 | Rename Device      | varies    | Set device name (UTF-8)            |
| 0x19 | Voice Language     | varies    | Set voice prompt language           |
| 0x1D | Tone Volume        | 1         | Set notification tone volume       |
| 0x1E | Take Photo         | 1         | Remote shutter trigger             |
| 0x1F | Standby            | 1         | Enter standby mode                 |
| 0x20 | EQ Params (v1)     | varies    | Read EQ parameters (6 bytes/band)  |
| 0x22 | EQ Params (v2)     | varies    | Read EQ parameters (7 bytes/band)  |
| 0x23 | LDAC               | 1         | LDAC codec toggle                  |
| 0x27 | Adaptive EQ        | 1         | Adaptive EQ toggle                 |
| 0x28 | ANC Result         | varies    | ANC calibration result             |
| 0x29 | ANC Wear           | 1-2       | ANC wearing state                  |
| 0x2B | Key Function       | varies    | Button action mapping              |
| 0x2C | Wearing Detection  | 3-4       | Wearing detection settings         |
| 0x2D | Spatial Audio      | 1         | Spatial audio toggle               |
| 0x2E | Music Mode         | 1         | Music playback mode                |
| 0x2F | Battery            | 3         | Battery levels (L/R/case)          |
| 0x30 | Version            | 3 or 6    | Firmware version                   |
| 0x32 | Env Adaptation     | 1         | Environmental adaptation           |
| 0x34 | TWS Enable         | 1         | TWS mode toggle                    |
| 0x35 | LED Switch         | 1         | LED switch                         |
| 0x36 | LED Effect         | varies    | LED color/effect settings          |
| 0x37 | Play Mode          | 1         | Playback mode                      |
| 0x39 | Focus Mode         | 1         | Focus mode toggle                  |
| 0x3A | Music Status       | 6         | Current music playback status      |
| 0x3B | Music Info         | varies    | Music file list                    |
| 0x3D | Tone Play          | 1         | Play a notification tone           |
| 0x3E | Sync Time          | 7         | Synchronize date/time              |
| 0x3F | Alarm              | 7         | Alarm clock management             |
| 0x43 | AI                 | 1         | AI assistant trigger               |
| 0x44 | Max EQ Count       | 1         | Query max custom EQ slots          |
| 0x45 | Custom EQ Test     | 1         | Custom EQ test mode                |
| 0x46 | EQ Left            | varies    | Left channel EQ (v2 format)        |
| 0x47 | EQ Right           | varies    | Right channel EQ (v2 format)       |
| 0x48 | In-Ear Sensitivity | 1         | In-ear detection sensitivity       |
| 0x4A | Game Config        | 1         | Game mode configuration            |
| 0xFE | Request Data       | 1         | Request current value of any cmd   |

---

## Command Details

### 0x01 — Reset Default
```
Send: [0x01, 0x00]
```
Resets all settings to factory defaults (without clearing pairing).

### 0x02 — Clear Pairing
```
Send: [0x02, 0x00]
```
Clears Bluetooth pairing information.

### 0x03 — Factory Reset
```
Send: [0x03, 0x00]
```
Full factory reset (clears pairing + resets all settings).

### 0x04 — Music Control
```
Send: [0x04, 0x01, action]
```
| Action | Description |
|--------|-------------|
| 0x01   | Play        |
| 0x02   | Pause       |
| 0x03   | Previous    |
| 0x04   | Next        |

### 0x05 — Light Flash
```
Send: [0x05, 0x01, state]
```
| State | Description    |
|-------|----------------|
| 0x00  | Off            |
| 0x01  | On (flashing)  |

### 0x06 — In-Ear Detection
```
Send: [0x06, 0x01, state]
```
| State | Description |
|-------|-------------|
| 0x01  | Enable      |
| 0x02  | Disable     |

**Response**: `[0x06, 0x01, state]` — `0x01` = enabled, other = disabled.

### 0x07 — Noise Value
```
Send: [0x07, 0x01, value]
Read: [0x07, 0x01, 0xFF]
```
Sets or reads the noise cancellation depth value (0–254). Sending `0xFF` requests the current value.

**Response**: `[0x07, 0x01, value]` — current noise value.

### 0x08 — Volume
```
Send: [0x08, 0x03, left, right, 0x00]
```
Sets left and right volume independently. Third byte is reserved (always 0x00).

**Response**: `[0x08, 0x03, leftVoice, rightVoice, maxVoice]` — current volumes and max volume (all unsigned).

### 0x09 — Low Latency (Game Mode)
```
Send: [0x09, 0x01, state]
```
| State | Description |
|-------|-------------|
| 0x01  | Enable      |
| 0x02  | Disable     |

**Response**: `[0x09, 0x01, state]` — `0x01` = enabled, other = disabled.

### 0x0C — Noise Cancel Mode

Simple mode (value <= 0xFF):
```
Send: [0x0C, 0x01, mode]
```

For advanced ANC scenes (value > 0xFF), the app unpacks the value and sends command 0x17 instead. See [ANC Setting](#0x17--anc-setting).

| Mode | Description   |
|------|---------------|
| 0x00 | Off           |
| 0x01 | ANC           |
| 0x02 | Outdoor       |
| 0x03 | Transparency  |

**Response**: `[0x0C, 0x01, mode]` — current mode.

### 0x0D — Test Mode
```
Send: [0x0D, 0x01, state]
```
| State | Description |
|-------|-------------|
| 0x01  | Enable      |
| 0x02  | Disable     |

**Response**: `[0x0D, 0x01, value]`

### 0x10 — Sleep Mode
```
Send: [0x10, 0x01, state]
```
| State | Description |
|-------|-------------|
| 0x01  | Enable      |
| 0x02  | Disable     |

**Response**: `[0x10, 0x01, state]`

### 0x11 — Ear Tip Fit Test (Compactness)
```
Send: [0x11, 0x02, left, right]
```
| Value | Description  |
|-------|--------------|
| 0x01  | Start test   |
| 0x02  | Stop test    |

**Response**: `[0x11, 0x03, status, leftResult, rightResult]`
- `status = 0x00`: Ready/testing
- `status = 0x02`: Result available
- `leftResult`/`rightResult`: fit test scores

### 0x12 — LED Mode
```
Send: [0x12, 0x01, state]
```
| State | Description |
|-------|-------------|
| 0x01  | Enable      |
| 0x02  | Disable     |

**Response**: `[0x12, 0x01, state]`

### 0x14 — Power Manager (Auto Power-Off)
```
Send: [0x14, 0x04, timeLo, timeHi, 0x00, 0x00]
```
- `timeLo`, `timeHi`: 16-bit little-endian auto power-off time (in minutes)
- Last 2 bytes reserved (0x00)

**Response**: `[0x14, 0x04, powerTimeLo, powerTimeHi, curTimeLo, curTimeHi]`
- Power-off time and current elapsed time, both 16-bit LE.

### 0x16 — Sound Balance
```
Send: [0x16, 0x01, value]
```
| Value | Description           |
|-------|-----------------------|
| 0x00  | Full left             |
| 0x32  | Center (default)      |
| 0x64  | Full right            |

**Response**: `[0x16, 0x01, value]`

### 0x17 — ANC Setting
```
Send: [0x17, 0x03, mode, subScene, noiseValue]
```
- `mode`: ANC mode (see noise cancel mode table)
- `subScene`: Sub-scene/level within the mode
- `noiseValue`: Noise depth value

The app uses a packed integer format for the combined noise mode:
```
packed = (mode << 16) | (subScene << 8) | noiseValue
```

When `mode=3, subScene=1, noiseValue=0`, the app remaps to `mode=3, subScene=2, noiseValue=0` (transparency special case).

#### ANC Scene Table

| Mode | SubScene | Description                            |
|------|----------|----------------------------------------|
| 0x00 | 0x00     | Off                                    |
| 0x02 | 0x01     | Silent Environment ANC (Level 1)       |
| 0x02 | 0x02     | Silent Environment ANC (Level 2)       |
| 0x02 | 0x03     | Silent Environment ANC (Level 3)       |
| 0x03 | 0x01     | Working Environment ANC (Level 1)      |
| 0x03 | 0x02     | Working Environment ANC (Level 2)      |
| 0x03 | 0x03     | Working Environment ANC (Level 3)      |
| 0x04 | 0x01     | Noisy Environment ANC (Level 1)        |
| 0x04 | 0x02     | Noisy Environment ANC (Level 2)        |
| 0x04 | 0x03     | Noisy Environment ANC (Level 3)        |
| 0x0A | 0x01–07  | Transparency Mode (Levels 1–7)         |

**Response**: `[0x17, 0x03, mode, subScene, noiseValue]`

### 0x18 — Rename Device
```
Send: [0x18, strLen, UTF-8 bytes...]
```

**Response**: `[0x18, strLen, UTF-8 bytes..., 0x00]` — NUL-terminated name string.

### 0x19 — Voice Language
```
Send: [0x19, strLen, UTF-8 bytes...]
```
Sets the voice prompt language (e.g., `"en"`, `"zh"`).

### 0x1D — Tone Volume
```
Send: [0x1D, 0x01, volume]
```

**Response**: `[0x1D, 0x01, toneVoice]` or `[0x1D, 0x02, toneVoice, maxVolume]`
- Default `maxVolume` is 100 if only 1 byte returned.

### 0x1E — Take Photo
```
Send: [0x1E, 0x01, action]
```
Remote camera shutter trigger.

### 0x1F — Standby
```
Send: [0x1F, 0x01, state]
```
Enter/exit standby mode.

### 0x20 — EQ Parameters (v1 Read)
```
Response: [0x20, paramLen, eqType, masterGainLo, masterGainHi, {freqLo, freqHi, gainLo, gainHi, qLo, qHi} x N]
```
- `eqType`: EQ preset index
- `masterGain`: 16-bit LE signed, divide by 100 for dB
- Per band (6 bytes each):
  - `freq`: 16-bit LE frequency in Hz
  - `gain`: 16-bit LE signed, divide by 100 for dB (clamped to +/-12.7 dB)
  - `q`: 16-bit LE, divide by 100 for Q factor

### 0x22 — EQ Parameters (v2 Read)
```
Response: [0x22, paramLen, eqType, masterGainLo, masterGainHi, {freqLo, freqHi, gainLo, gainHi, qLo, qHi, bandType} x N]
```
Same as v1 but with an extra `bandType` byte per band (7 bytes per band).

### 0x23 — LDAC
```
Send: [0x23, 0x01, state]
```
Toggle LDAC high-quality codec.

### 0x27 — Adaptive EQ
```
Send: [0x27, 0x01, state]
```
Toggle adaptive EQ.

### 0x28 — ANC Result
```
Response: [0x28, paramLen, ...]
```
ANC calibration result data.

### 0x29 — ANC Wear
```
Response: [0x29, 0x01, value]
```
If only 1 byte of params, the value is duplicated internally to 2 identical values.

### 0x2B — Key Function
```
Response: [0x2B, paramLen, keyId1, funId1, keyId2, funId2, ...]
```
Button action mapping as key-value pairs. See [Key Functions](#key-functions) for ID tables.

### 0x2C — Wearing Detection
```
Send: [0x2C, 0x03, enable, musicIndex, ancIndex]
Send: [0x2C, 0x03, enable, musicIndex, ancIndex, toneEnable]  (v2, 4th byte appended)
```
- `enable`: 0x01 = on, 0x02 = off
- `musicIndex`: Action when removed (e.g., pause music)
- `ancIndex`: Action when removed (e.g., change ANC mode)
- `toneEnable` (v2): Tone notification toggle

**Response**: `[0x2C, 0x03, isEnable, musicIndex, ancIndex]` or `[0x2C, 0x04, isEnable, musicIndex, ancIndex, toneEnable]`

### 0x2D — Spatial Audio
```
Send: [0x2D, 0x01, state]
```
Toggle spatial audio / head tracking.

### 0x2E — Music Mode
```
Send: [0x2E, 0x01, mode]
```

### 0x2F — Battery
```
Response: [0x2F, 0x03, left, right, box]
```
Each byte encodes:
- **Bit 7** (0x80): Charging flag (1 = charging)
- **Bits 0-6** (0x7F): Battery level (0-127, maps to percentage)

| Type  | Index |
|-------|-------|
| Left  | 0     |
| Right | 1     |
| Case  | 2     |

### 0x30 — Version
```
Response (3 bytes): [0x30, 0x03, major, minor, patch]
Response (6 bytes): [0x30, 0x06, Lmajor, Lminor, Lpatch, Rmajor, Rminor, Rpatch]
```
Formatted as `"major.minor.patch"`. 6-byte variant includes separate left and right firmware versions.

### 0x32 — Environmental Adaptation
```
Send: [0x32, 0x01, state]
```

### 0x34 — TWS Enable
```
Send: [0x34, 0x01, state]
```

### 0x35 — LED Switch
```
Send: [0x35, 0x01, state]
```

### 0x36 — LED Effect
```
Send: [0x36, paramLen, speed, brightness, effectIndex, R1, G1, B1, R2, G2, B2, ...]
```
- `paramLen`: 3 + (3 x numColors)
- `speed`: Animation speed
- `brightness`: LED brightness
- `effectIndex`: Effect pattern index
- Colors: RGB triplets (3 bytes each)

**Response**: Same format.

### 0x37 — Play Mode
```
Send: [0x37, 0x01, mode]
```

### 0x39 — Focus Mode
```
Send: [0x39, 0x01, state]
```

### 0x3A — Music Status
```
Send: [0x3A, 0x06, musicID_b0, musicID_b1, musicID_b2, musicID_b3, isPlaying, playMode]
```
- `musicID`: 32-bit little-endian music ID
- `isPlaying`: 0x01 = playing, 0x00 = paused
- `playMode`: Playback mode byte

**Response**: Same format.

### 0x3B — Music Info
```
Send: [0x3B, paramLen, startToneID, {musicID_b0, b1, b2, b3, totalLo, totalHi} x N]
```
- `startToneID`: Starting tone index (byte at offset 2)
- Per file (6 bytes):
  - `musicID`: 32-bit LE music file ID
  - `total`: 16-bit LE total count/duration

**Response**: Same format. `paramLen = 1 + (6 x N)`.

### 0x3D — Tone Play
```
Send: [0x3D, 0x01, toneID]
```
Play a specific notification tone by ID.

### 0x3E — Sync Time
```
Send: [0x3E, 0x07, year, month, day, hour, minute, second, dayOfWeekBitmask]
```
- `year`: Year mod 100 (e.g., 25 for 2025)
- `month`: 1-12
- `day`: 1-31
- `hour`: 0-23 (24-hour)
- `minute`: 0-59
- `second`: 0-59
- `dayOfWeekBitmask`: `1 << (calendar_day_of_week - 1)` where Sunday=1

### 0x3F — Alarm
```
Add:    [0x3F, 0x07, 0x01, alarmID, enable, hour, minute, cycle, 0x05]
Delete: [0x3F, 0x07, 0x02, alarmID, 0x00,   0x00, 0x00,   0x00, 0x00]
Edit:   [0x3F, 0x07, 0x03, alarmID, enable, hour, minute, cycle, index]
```
- `alarmID`: Unique alarm identifier
- `enable`: 0x01 = enabled, 0x00 = disabled
- `hour`/`minute`: Alarm time
- `cycle`: Day-of-week bitmask (bit 0 = Sunday, bit 6 = Saturday)
- Byte at index 2 is the operation type: 1=add, 2=delete, 3=edit

**Response**: `[0x3F, paramLen, count, {alarmID, enable, hour, minute, cycle, ?} x N]`
- Stride of 6 bytes per alarm, parsed from offset 1.

### 0x43 — AI
```
Send: [0x43, 0x01, action]
```
AI assistant trigger.

### 0x44 — Max EQ Count
```
Response: [0x44, 0x01, count]
```
Maximum number of custom EQ presets supported.

### 0x45 — Custom EQ Test
```
Send: [0x45, 0x01, state]
```

### 0x46 — EQ Left Channel
```
Send: [0x46, paramLen, eqIndex, masterGainLo, masterGainHi, {freqLo,freqHi,gainLo,gainHi,qLo,qHi,bandType} x N]
```
Same as v2 EQ format but for left channel only. Written via `00001001` with `0xFF` framing.

### 0x47 — EQ Right Channel
```
Send: [0x47, paramLen, eqIndex, masterGainLo, masterGainHi, {freqLo,freqHi,gainLo,gainHi,qLo,qHi,bandType} x N]
```
Same as v2 EQ format but for right channel only.

### 0x48 — In-Ear Sensitivity
```
Send: [0x48, 0x01, level]
```

### 0x4A — Game Config
```
Send: [0x4A, 0x01, config]
```

### 0xFE — Request Data
```
Send: [0xFE, 0x01, cmdId]
```
Request the current value/setting for any command ID. The device responds with the corresponding command's response format.

---

## EQ Settings

EQ data is written directly to characteristic `0000000B-0000-1000-8000-00805f9b34fb` **without** the `0xFF` packet framing.

### Simple EQ (via characteristic)
```
Write to 0000000B: [eqType, data...]
```

### Parametric EQ Write (v1, cmd 0x20)
```
Write to 00001001 (with 0xFF framing):
[0x20, paramLen, eqIndex, masterGainLo, masterGainHi, {freqLo, freqHi, gainLo, gainHi, qLo, qHi} x N]
```
- `paramLen = N x 6 + 3`
- `masterGain = float x 100`, 16-bit LE signed
- `gain = float x 100`, 16-bit LE signed, clamped to [-1270, 1270] (+/-12.7 dB)
- `q = float x 100`, 16-bit LE
- `freq`: raw Hz, 16-bit LE

### Parametric EQ Write (v2, cmd 0x22)
```
Write to 00001001 (with 0xFF framing):
[0x22, paramLen, eqIndex, masterGainLo, masterGainHi, {freqLo, freqHi, gainLo, gainHi, qLo, qHi, bandType} x N]
```
- `paramLen = N x 7 + 3`
- Same encoding as v1 with an extra `bandType` byte per band

### Per-Channel EQ (cmd 0x46/0x47)
Uses v2 format but with command IDs `0x46` (left) and `0x47` (right).

---

## Key Functions

Key function data is read from and written to characteristic `0000000D-0000-1000-8000-00805f9b34fb` **without** the `0xFF` packet framing.

### Format
```
[keyId1, funId1, keyId2, funId2, ...]
```
Key-value pairs, 2 bytes each.

### Key IDs

#### Music Mode (Normal)

| Key ID | Description        |
|--------|--------------------|
| 0x01   | Left single tap    |
| 0x02   | Right single tap   |
| 0x03   | Left double tap    |
| 0x04   | Right double tap   |
| 0x05   | Left triple tap    |
| 0x06   | Right triple tap   |
| 0x07   | Left quad tap      |
| 0x08   | Right quad tap     |
| 0x09   | Left long press    |
| 0x0A   | Right long press   |

#### Voice/Call Mode

| Key ID | Description        |
|--------|--------------------|
| 0x15   | Left single tap    |
| 0x16   | Right single tap   |
| 0x17   | Left double tap    |
| 0x18   | Right double tap   |
| 0x19   | Left triple tap    |
| 0x1A   | Right triple tap   |
| 0x1B   | Left quad tap      |
| 0x1C   | Right quad tap     |
| 0x1D   | Left long press    |
| 0x1E   | Right long press   |

### Function IDs

| Fun ID | Description       |
|--------|-------------------|
| 0x00   | None / Disabled   |
| 0x01   | Play / Pause      |
| 0x02   | Previous track    |
| 0x03   | Next track        |
| 0x04   | Voice assistant   |
| 0x05   | Volume up         |
| 0x06   | Volume down       |
| 0x07   | Game mode toggle  |
| 0x08   | Answer call       |
| 0x09   | Reject call       |
| 0x0A   | Hold call         |
| 0x0B   | Redial            |

---

## Noise Cancel Mode (Legacy)

For command 0x0C with simple single-byte values (models without ANC setting support):

| Mode | Description   |
|------|---------------|
| 0x00 | Off           |
| 0x01 | ANC           |
| 0x02 | Outdoor       |
| 0x03 | Transparency  |

For models with ANC setting support, the app uses command 0x17 with a packed integer:
```
packed = (mode << 16) | (subScene << 8) | noiseValue
```

See [ANC Setting](#0x17--anc-setting) for the scene table.

---

## Battery Read

Battery can be read from characteristic `00000008-0000-1000-8000-00805f9b34fb` or via command 0x2F notification.

**Format**: 3 bytes `[left, right, box]`

Each byte:
```
Bit 7 (0x80): 1 = charging, 0 = not charging
Bits 0-6 (0x7F): Battery level (0-127)
```

---

## Version Read

Version can be read from characteristic `00000007-0000-1000-8000-00805f9b34fb` or via command 0x30 notification.

**3 bytes**: `[major, minor, patch]` -> `"major.minor.patch"` (left earbud only)

**6 bytes**: `[Lmajor, Lminor, Lpatch, Rmajor, Rminor, Rpatch]` -> separate left + right versions

---

## Response Parsing

The device sends notifications on `00001002`. The data is parsed by stripping the `[0xFF, bodyLen]` header and iterating through command blocks.

### Default Response Types

| Param Count | Handling                                    |
|-------------|---------------------------------------------|
| 1 byte      | Parsed as a single value (SingleDataBean)   |
| N bytes     | Parsed as hex string (CmdHexDataBean)       |

Commands with specific parsers use their respective DataBean classes as documented in each command section above.

---

## Protocol Variants

The QCY app supports multiple chipset vendors with different protocol paths:

| Vendor  | Protocol Path                                | Notes                                |
|---------|----------------------------------------------|--------------------------------------|
| QCY     | `QCYHeadsetClient` via `00001001`/`00001002` | Main protocol documented above       |
| JL/Jieli| `JLDeviceImpl` via JL SDK                    | Jieli-specific BLE commands          |
| Airoha  | RACE/MMI commands, `AirohaPEQManager`        | Airoha-specific ANC and EQ handling  |
| ZR      | Via UUID `0000000E`                          | ZR device settings characteristic    |

This documentation covers the **QCY standard protocol** (used by most QCY earphones like HT07, T13, etc.).

---

## V1 Protocol (Legacy Characteristics)

Older models may use individual characteristics instead of the unified command channel:

| Characteristic UUID | V1 Purpose           | Cmd ID (internal) |
|--------------------|----------------------|-------------------|
| `00000001`         | Left single tap key  | -                 |
| `00000002`         | Right single tap key | -                 |
| `00000003`         | Left double tap key  | -                 |
| `00000004`         | Right double tap key | -                 |
| `00000005`         | Left triple tap key  | -                 |
| `00000006`         | Right triple tap key | -                 |
| `00000007`         | Version read         | 0x1004            |
| `00000008`         | Battery read         | 0x1001            |
| `00000009`         | Language read        | -                 |
| `0000000A`         | Reset                | 0x1005            |
| `0000000B`         | EQ read/write        | 0x1003 / 0x1013   |
| `0000000C`         | Send time            | -                 |
| `0000000D`         | Key function V2      | -                 |
| `0000000E`         | ZR device settings   | -                 |
| `0000000F`         | In-ear check (JL)    | -                 |

---

## Chinese Term Reference

Terms found in the protocol implementation:

| Chinese       | Pinyin    | English                         |
|---------------|-----------|----------------------------------|
| 入耳检测       | ruer      | In-ear detection                |
| 低延时         | diyanshi  | Low latency / Game mode         |
| 监听          | jianting  | Monitoring / Transparency mode   |
| 降噪          | noise     | Noise cancellation / ANC        |
| 佩戴          | peidai    | Wearing detection               |
| 声道平衡       | -         | Channel balance                 |
| 闪灯          | light     | LED flash                       |
| 配对          | pair      | Bluetooth pairing               |