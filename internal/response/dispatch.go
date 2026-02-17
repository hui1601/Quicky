package response

import "strings"

// EventType represents the type of notification received from the device.
type EventType byte

const (
	EventUnknown          EventType = 0x00
	EventResetDefault     EventType = 0x01
	EventClearPairing     EventType = 0x02
	EventFactoryReset     EventType = 0x03
	EventMusicControl     EventType = 0x04
	EventLightFlash       EventType = 0x05
	EventInEarTest        EventType = 0x06
	EventNoiseValue       EventType = 0x07
	EventVolume           EventType = 0x08
	EventLowLatency       EventType = 0x09
	EventMonitoring       EventType = 0x0A
	EventNoiseCancelMode  EventType = 0x0C
	EventTestMode         EventType = 0x0D
	EventSleepMode        EventType = 0x10
	EventEarTipFit        EventType = 0x11
	EventLEDMode          EventType = 0x12
	EventPowerManager     EventType = 0x14
	EventSoundBalance     EventType = 0x16
	EventANCSetting       EventType = 0x17
	EventRename           EventType = 0x18
	EventAudioLang        EventType = 0x19
	EventToneVolume       EventType = 0x1D
	EventTakePhoto        EventType = 0x1E
	EventStandby          EventType = 0x1F
	EventEQV1             EventType = 0x20
	EventEQV2             EventType = 0x22
	EventLDAC             EventType = 0x23
	EventAdaptiveEQ       EventType = 0x27
	EventANCResult        EventType = 0x28
	EventANCWear          EventType = 0x29
	EventKeyFunction      EventType = 0x2B
	EventWearingDetection EventType = 0x2C
	EventSpatialAudio     EventType = 0x2D
	EventMusicMode        EventType = 0x2E
	EventBattery          EventType = 0x2F
	EventVersion          EventType = 0x30
	EventEnvAdaptation    EventType = 0x32
	EventTWSEnable        EventType = 0x34
	EventLEDSwitch        EventType = 0x35
	EventLEDEffect        EventType = 0x36
	EventPlayMode         EventType = 0x37
	EventFocusMode        EventType = 0x39
	EventMusicStatus      EventType = 0x3A
	EventMusicInfo        EventType = 0x3B
	EventTonePlay         EventType = 0x3D
	EventSyncTime         EventType = 0x3E
	EventAlarm            EventType = 0x3F
	EventAI               EventType = 0x43
	EventMaxEQCount       EventType = 0x44
	EventCustomEQTest     EventType = 0x45
	EventEQLeft           EventType = 0x46
	EventEQRight          EventType = 0x47
	EventInEarSensitivity EventType = 0x48
	EventGameConfig       EventType = 0x4A
)

// Event represents a parsed notification from the device.
type Event struct {
	Type   EventType
	CmdID  byte
	Raw    []byte // raw parameter bytes (copy)
	Parsed any    // typed parsed value or nil
	Error  error  // parsing error, if any
}

// Dispatch parses a command ID and its parameters into a typed Event.
func Dispatch(cmdID byte, params []byte) Event {
	raw := make([]byte, len(params))
	copy(raw, params)

	ev := Event{
		Type:  EventType(cmdID),
		CmdID: cmdID,
		Raw:   raw,
	}

	switch cmdID {
	case 0x2F:
		v, err := ParseBattery(params)
		ev.Parsed, ev.Error = v, err
	case 0x30:
		v, err := ParseVersion(params)
		ev.Parsed, ev.Error = v, err
	case 0x17:
		v, err := ParseANCSetting(params)
		ev.Parsed, ev.Error = v, err
	case 0x14:
		v, err := ParsePowerManager(params)
		ev.Parsed, ev.Error = v, err
	case 0x08:
		v, err := ParseVolume(params)
		ev.Parsed, ev.Error = v, err
	case 0x1D:
		v, err := ParseToneVolume(params)
		ev.Parsed, ev.Error = v, err
	case 0x2C:
		v, err := ParseWearingDetection(params)
		ev.Parsed, ev.Error = v, err
	case 0x11:
		v, err := ParseEarTipFit(params)
		ev.Parsed, ev.Error = v, err
	case 0x20:
		v, err := ParseEQV1(params)
		ev.Parsed, ev.Error = v, err
	case 0x22:
		v, err := ParseEQV2(params)
		ev.Parsed, ev.Error = v, err
	case 0x46, 0x47: // EQ left/right channels use v2 format
		v, err := ParseEQV2(params)
		ev.Parsed, ev.Error = v, err
	case 0x2B:
		v, err := ParseKeyFunction(params)
		ev.Parsed, ev.Error = v, err
	case 0x36:
		v, err := ParseLEDEffect(params)
		ev.Parsed, ev.Error = v, err
	case 0x3F:
		v, err := ParseAlarmList(params)
		ev.Parsed, ev.Error = v, err
	case 0x3A:
		v, err := ParseMusicStatus(params)
		ev.Parsed, ev.Error = v, err
	case 0x3B:
		startID, files, err := ParseMusicInfo(params)
		ev.Parsed = MusicInfoResult{StartToneID: startID, Files: files}
		ev.Error = err

	case 0x18: // rename, NUL-terminated UTF-8
		ev.Parsed = strings.TrimRight(string(params), "\x00")
	case 0x19: // audio language
		ev.Parsed = strings.TrimRight(string(params), "\x00")

	case 0x28, 0x29: // ANC result/wear — variable-length raw
		ev.Parsed = raw

	case 0x01, 0x02, 0x03, // reset/clear/factory (ack)
		0x04, // music control
		0x05, // light flash
		0x06, // in-ear test
		0x07, // noise value
		0x09, // low latency
		0x0A, // monitoring
		0x0C, // noise cancel mode
		0x0D, // test mode
		0x10, // sleep mode
		0x12, // LED mode
		0x16, // sound balance
		0x1E, // take photo
		0x1F, // standby
		0x23, // LDAC
		0x27, // adaptive EQ
		0x2D, // spatial audio
		0x2E, // music mode
		0x32, // env adaptation
		0x34, // TWS enable
		0x35, // LED switch
		0x37, // play mode
		0x39, // focus mode
		0x3D, // tone play
		0x43, // AI
		0x44, // max EQ count
		0x45, // custom EQ test
		0x48, // in-ear sensitivity
		0x4A: // game config
		if len(params) >= 1 {
			ev.Parsed = params[0]
		}

	default:
		ev.Type = EventUnknown
	}

	return ev
}

// MusicInfoResult wraps the multi-return ParseMusicInfo result into a single struct.
type MusicInfoResult struct {
	StartToneID byte
	Files       []MusicFile
}
