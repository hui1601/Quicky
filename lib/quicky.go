package quicky

import (
	"context"
	"image/color"
	"time"

	"github.com/hui1601/Quicky/internal/command"
	"github.com/hui1601/Quicky/internal/device"
	"github.com/hui1601/Quicky/internal/response"
)

type Client struct {
	dev *device.Client
}

func New(mac string) (*Client, error) {
	dev, err := device.NewClient(mac)
	if err != nil {
		return nil, err
	}
	return &Client{dev: dev}, nil
}

func (c *Client) Connect(ctx context.Context) error {
	return c.dev.Connect(ctx)
}

func (c *Client) Disconnect() error {
	return c.dev.Disconnect()
}

func (c *Client) Connected() bool {
	return c.dev.Connected()
}

func (c *Client) Events() <-chan Event {
	ch := make(chan Event, 32)
	go func() {
		for ev := range c.dev.Events() {
			ch <- fromInternalEvent(ev)
		}
		close(ch)
	}()
	return ch
}

func (c *Client) send(cmd *command.Command) error {
	return c.dev.SendCommand(cmd)
}

func (c *Client) ResetDefault() error {
	return c.send(command.NewResetDefaultCommand())
}

func (c *Client) ClearPairing() error {
	return c.send(command.NewResetPairCommand())
}

func (c *Client) FactoryReset() error {
	return c.send(command.NewFactoryResetCommand())
}

type MusicAction = command.MusicControlAction

func (c *Client) MusicControl(action MusicAction) error {
	return c.send(command.NewMusicControlCommand(action))
}

func (c *Client) SetLightFlash(on bool) error {
	return c.send(command.NewLightFlashCommand(on))
}

func (c *Client) SetInEarDetection(on bool) error {
	return c.send(command.NewInEarTestCommand(on))
}

func (c *Client) SetNoiseValue(value byte) error {
	return c.send(command.NewNoiseValueCommand(value))
}

type NoiseCancelMode = command.NoiseCancelMode

const (
	NoiseCancelOff          = command.NoiseCancelOff
	NoiseCancelANC          = command.NoiseCancelANC
	NoiseCancelOutdoor      = command.NoiseCancelOutdoor
	NoiseCancelTransparency = command.NoiseCancelTransparency
)

func (c *Client) SetNoiseCancelMode(mode NoiseCancelMode) error {
	return c.send(command.NewNoiseCancelModeCommand(mode))
}

func (c *Client) SetANCSetting(mode, subScene, noiseValue byte) error {
	return c.send(command.NewANCSettingCommand(mode, subScene, noiseValue))
}

func (c *Client) SetVolume(left, right byte) error {
	return c.send(command.NewVolumeCommand(left, right))
}

func (c *Client) SetLowLatency(on bool) error {
	return c.send(command.NewLowLatencyCommand(on))
}

func (c *Client) SetMonitoring(value byte) error {
	return c.send(command.NewMonitoringCommand(value))
}

func (c *Client) SetTestMode(on bool) error {
	return c.send(command.NewTestModeCommand(on))
}

func (c *Client) SetSleepMode(on bool) error {
	return c.send(command.NewSleepModeCommand(on))
}

func (c *Client) StartEarTipFitTest() error {
	return c.send(command.NewEarTipFitStartCommand())
}

func (c *Client) StopEarTipFitTest() error {
	return c.send(command.NewEarTipFitStopCommand())
}

func (c *Client) SetLEDMode(on bool) error {
	return c.send(command.NewLEDModeCommand(on))
}

func (c *Client) SetPowerManager(reserveTime, currentTime int32) error {
	return c.send(command.NewBookingRebootCommand(reserveTime, currentTime))
}

func (c *Client) SetSoundBalance(value byte) error {
	return c.send(command.NewSoundBalanceCommand(int32(value)))
}

func (c *Client) SetName(name string) error {
	return c.send(command.NewChangeNameCommand(name))
}

func (c *Client) SetAudioLanguage(lang string) error {
	return c.send(command.NewAudioLangCommand(lang))
}

func (c *Client) SetToneVolume(volume byte) error {
	return c.send(command.NewToneVolumeCommand(volume))
}

func (c *Client) TakePhoto(action byte) error {
	return c.send(command.NewTakePhotoCommand(action))
}

func (c *Client) SetStandby(state byte) error {
	return c.send(command.NewStandbyCommand(state))
}

func (c *Client) SetLDAC(on bool) error {
	return c.send(command.NewLDACCommand(on))
}

func (c *Client) SetAdaptiveEQ(on bool) error {
	return c.send(command.NewAdaptiveEQCommand(on))
}

func (c *Client) SetWearingDetection(enable bool, musicIndex, ancIndex byte) error {
	return c.send(command.NewWearingDetectionCommand(enable, musicIndex, ancIndex))
}

func (c *Client) SetWearingDetectionV2(enable bool, musicIndex, ancIndex byte, toneEnable bool) error {
	return c.send(command.NewWearingDetectionV2Command(enable, musicIndex, ancIndex, toneEnable))
}

func (c *Client) SetSpatialAudio(on bool) error {
	return c.send(command.NewSpatialAudioCommand(on))
}

func (c *Client) SetMusicMode(mode byte) error {
	return c.send(command.NewMusicModeCommand(mode))
}

func (c *Client) SetEnvAdaptation(state byte) error {
	return c.send(command.NewEnvAdaptationCommand(state))
}

func (c *Client) SetTWSEnable(on bool) error {
	return c.send(command.NewTWSEnableCommand(on))
}

func (c *Client) SetLEDSwitch(on bool) error {
	return c.send(command.NewLEDSwitchCommand(on))
}

func (c *Client) SetLEDEffect(speed, brightness, effectIndex byte, colors []color.RGBA) error {
	return c.send(command.NewLEDEffectCommand(speed, brightness, effectIndex, colors))
}

func (c *Client) SetPlayMode(mode byte) error {
	return c.send(command.NewPlayModeCommand(mode))
}

func (c *Client) SetFocusMode(on bool) error {
	return c.send(command.NewFocusModeCommand(on))
}

func (c *Client) SetMusicStatus(musicID uint32, isPlaying bool, playMode byte) error {
	return c.send(command.NewMusicStatusCommand(musicID, isPlaying, playMode))
}

type MusicFile struct {
	MusicID uint32
	Total   uint16
}

func (c *Client) SetMusicInfo(startToneID byte, files []MusicFile) error {
	internal := make([]command.MusicFile, len(files))
	for i, f := range files {
		internal[i] = command.MusicFile{MusicID: f.MusicID, Total: f.Total}
	}
	return c.send(command.NewMusicInfoCommand(startToneID, internal))
}

func (c *Client) PlayTone(toneID byte) error {
	return c.send(command.NewTonePlayCommand(toneID))
}

func (c *Client) SyncTime(t time.Time) error {
	return c.send(command.NewSyncTimeCommand(t))
}

func (c *Client) AddAlarm(alarmID, hour, minute, cycle byte) error {
	return c.send(command.NewAlarmAddCommand(alarmID, hour, minute, cycle))
}

func (c *Client) DeleteAlarm(alarmID byte) error {
	return c.send(command.NewAlarmDeleteCommand(alarmID))
}

func (c *Client) EditAlarm(alarmID byte, enable bool, hour, minute, cycle, index byte) error {
	return c.send(command.NewAlarmEditCommand(alarmID, enable, hour, minute, cycle, index))
}

func (c *Client) TriggerAI(action byte) error {
	return c.send(command.NewAICommand(action))
}

func (c *Client) SetCustomEQTest(state byte) error {
	return c.send(command.NewCustomEQTestCommand(state))
}

func (c *Client) SetInEarSensitivity(level byte) error {
	return c.send(command.NewInEarSensitivityCommand(level))
}

func (c *Client) SetGameConfig(config byte) error {
	return c.send(command.NewGameConfigCommand(config))
}

func (c *Client) RequestData(cmdID byte) error {
	return c.send(command.NewRequestDataCommand(cmdID))
}

type EQBand struct {
	Freq     uint16
	Gain     int16
	Q        uint16
	BandType byte
}

func toInternalBands(bands []EQBand) []command.EQBand {
	out := make([]command.EQBand, len(bands))
	for i, b := range bands {
		out[i] = command.EQBand{Freq: b.Freq, Gain: b.Gain, Q: b.Q, BandType: b.BandType}
	}
	return out
}

func (c *Client) SetEQV1(eqIndex byte, masterGain int16, bands []EQBand) error {
	return c.send(command.NewEQV1Command(eqIndex, masterGain, toInternalBands(bands)))
}

func (c *Client) SetEQV2(eqIndex byte, masterGain int16, bands []EQBand) error {
	return c.send(command.NewEQV2Command(eqIndex, masterGain, toInternalBands(bands)))
}

func (c *Client) SetEQLeft(eqIndex byte, masterGain int16, bands []EQBand) error {
	return c.send(command.NewEQLeftCommand(eqIndex, masterGain, toInternalBands(bands)))
}

func (c *Client) SetEQRight(eqIndex byte, masterGain int16, bands []EQBand) error {
	return c.send(command.NewEQRightCommand(eqIndex, masterGain, toInternalBands(bands)))
}

func (c *Client) WriteEQDirect(eqType byte, data []byte) error {
	return c.dev.WriteEQ(command.NewEQDirectData(eqType, data))
}

type KeyID = command.KeyID
type FuncID = command.FuncID

const (
	KeyMusicLeftSingle  = command.KeyMusicLeftSingle
	KeyMusicRightSingle = command.KeyMusicRightSingle
	KeyMusicLeftDouble  = command.KeyMusicLeftDouble
	KeyMusicRightDouble = command.KeyMusicRightDouble
	KeyMusicLeftTriple  = command.KeyMusicLeftTriple
	KeyMusicRightTriple = command.KeyMusicRightTriple
	KeyMusicLeftQuad    = command.KeyMusicLeftQuad
	KeyMusicRightQuad   = command.KeyMusicRightQuad
	KeyMusicLeftLong    = command.KeyMusicLeftLong
	KeyMusicRightLong   = command.KeyMusicRightLong

	KeyVoiceLeftSingle  = command.KeyVoiceLeftSingle
	KeyVoiceRightSingle = command.KeyVoiceRightSingle
	KeyVoiceLeftDouble  = command.KeyVoiceLeftDouble
	KeyVoiceRightDouble = command.KeyVoiceRightDouble
	KeyVoiceLeftTriple  = command.KeyVoiceLeftTriple
	KeyVoiceRightTriple = command.KeyVoiceRightTriple
	KeyVoiceLeftQuad    = command.KeyVoiceLeftQuad
	KeyVoiceRightQuad   = command.KeyVoiceRightQuad
	KeyVoiceLeftLong    = command.KeyVoiceLeftLong
	KeyVoiceRightLong   = command.KeyVoiceRightLong

	FuncNone           = command.FuncNone
	FuncPlayPause      = command.FuncPlayPause
	FuncPrevious       = command.FuncPrevious
	FuncNext           = command.FuncNext
	FuncVoiceAssistant = command.FuncVoiceAssistant
	FuncVolumeUp       = command.FuncVolumeUp
	FuncVolumeDown     = command.FuncVolumeDown
	FuncGameMode       = command.FuncGameMode
	FuncAnswerCall     = command.FuncAnswerCall
	FuncRejectCall     = command.FuncRejectCall
	FuncHoldCall       = command.FuncHoldCall
	FuncRedial         = command.FuncRedial
)

type KeyMapping struct {
	Key  KeyID
	Func FuncID
}

func (c *Client) WriteKeyFunctions(mappings []KeyMapping) error {
	internal := make([]command.KeyMapping, len(mappings))
	for i, m := range mappings {
		internal[i] = command.KeyMapping{Key: m.Key, Func: m.Func}
	}
	return c.dev.WriteKeyFunction(command.NewKeyFunctionDirectData(internal))
}

type Battery = response.Battery
type BatteryInfo = response.BatteryInfo
type Version = response.Version

func (c *Client) ReadBattery() (Battery, error) {
	return c.dev.ReadBattery()
}

func (c *Client) ReadVersion() (Version, error) {
	return c.dev.ReadVersion()
}
