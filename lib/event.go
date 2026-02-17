package quicky

import "github.com/hui1601/Quicky/internal/response"

type EventType = response.EventType

const (
	EventUnknown          = response.EventUnknown
	EventResetDefault     = response.EventResetDefault
	EventClearPairing     = response.EventClearPairing
	EventFactoryReset     = response.EventFactoryReset
	EventMusicControl     = response.EventMusicControl
	EventLightFlash       = response.EventLightFlash
	EventInEarTest        = response.EventInEarTest
	EventNoiseValue       = response.EventNoiseValue
	EventVolume           = response.EventVolume
	EventLowLatency       = response.EventLowLatency
	EventMonitoring       = response.EventMonitoring
	EventNoiseCancelMode  = response.EventNoiseCancelMode
	EventTestMode         = response.EventTestMode
	EventSleepMode        = response.EventSleepMode
	EventEarTipFit        = response.EventEarTipFit
	EventLEDMode          = response.EventLEDMode
	EventPowerManager     = response.EventPowerManager
	EventSoundBalance     = response.EventSoundBalance
	EventANCSetting       = response.EventANCSetting
	EventRename           = response.EventRename
	EventAudioLang        = response.EventAudioLang
	EventToneVolume       = response.EventToneVolume
	EventTakePhoto        = response.EventTakePhoto
	EventStandby          = response.EventStandby
	EventEQV1             = response.EventEQV1
	EventEQV2             = response.EventEQV2
	EventLDAC             = response.EventLDAC
	EventAdaptiveEQ       = response.EventAdaptiveEQ
	EventANCResult        = response.EventANCResult
	EventANCWear          = response.EventANCWear
	EventKeyFunction      = response.EventKeyFunction
	EventWearingDetection = response.EventWearingDetection
	EventSpatialAudio     = response.EventSpatialAudio
	EventMusicMode        = response.EventMusicMode
	EventBattery          = response.EventBattery
	EventVersion          = response.EventVersion
	EventEnvAdaptation    = response.EventEnvAdaptation
	EventTWSEnable        = response.EventTWSEnable
	EventLEDSwitch        = response.EventLEDSwitch
	EventLEDEffect        = response.EventLEDEffect
	EventPlayMode         = response.EventPlayMode
	EventFocusMode        = response.EventFocusMode
	EventMusicStatus      = response.EventMusicStatus
	EventMusicInfo        = response.EventMusicInfo
	EventTonePlay         = response.EventTonePlay
	EventSyncTime         = response.EventSyncTime
	EventAlarm            = response.EventAlarm
	EventAI               = response.EventAI
	EventMaxEQCount       = response.EventMaxEQCount
	EventCustomEQTest     = response.EventCustomEQTest
	EventEQLeft           = response.EventEQLeft
	EventEQRight          = response.EventEQRight
	EventInEarSensitivity = response.EventInEarSensitivity
	EventGameConfig       = response.EventGameConfig
)

type Event struct {
	Type   EventType
	CmdID  byte
	Raw    []byte
	Parsed any
	Error  error
}

func fromInternalEvent(ev response.Event) Event {
	return Event{
		Type:   ev.Type,
		CmdID:  ev.CmdID,
		Raw:    ev.Raw,
		Parsed: ev.Parsed,
		Error:  ev.Error,
	}
}

type ANCSetting = response.ANCSetting
type PowerManager = response.PowerManager
type Volume = response.Volume
type ToneVolume = response.ToneVolume
type WearingDetection = response.WearingDetection
type EarTipFitResult = response.EarTipFitResult
type ResponseEQBand = response.EQBand
type EQParams = response.EQParams
type ResponseKeyMapping = response.KeyMapping
type LEDEffect = response.LEDEffect
type Alarm = response.Alarm
type MusicStatus = response.MusicStatus
type MusicInfoResult = response.MusicInfoResult
