package response

import (
	"encoding/binary"
	"fmt"
	"image/color"
)

type Battery struct {
	Left     BatteryInfo
	Right    BatteryInfo
	Case     BatteryInfo
}

type BatteryInfo struct {
	Level    byte
	Charging bool
}

func ParseBattery(params []byte) (Battery, error) {
	if len(params) < 3 {
		return Battery{}, fmt.Errorf("battery: need 3 bytes, got %d", len(params))
	}
	return Battery{
		Left:  parseBatteryByte(params[0]),
		Right: parseBatteryByte(params[1]),
		Case:  parseBatteryByte(params[2]),
	}, nil
}

func parseBatteryByte(b byte) BatteryInfo {
	return BatteryInfo{
		Level:    b & 0x7f,
		Charging: b&0x80 != 0,
	}
}

type Version struct {
	Left  string
	Right string
}

func ParseVersion(params []byte) (Version, error) {
	if len(params) == 3 {
		return Version{
			Left: fmt.Sprintf("%d.%d.%d", params[0], params[1], params[2]),
		}, nil
	}
	if len(params) == 6 {
		return Version{
			Left:  fmt.Sprintf("%d.%d.%d", params[0], params[1], params[2]),
			Right: fmt.Sprintf("%d.%d.%d", params[3], params[4], params[5]),
		}, nil
	}
	return Version{}, fmt.Errorf("version: need 3 or 6 bytes, got %d", len(params))
}

type ANCSetting struct {
	Mode       byte
	SubScene   byte
	NoiseValue byte
}

func ParseANCSetting(params []byte) (ANCSetting, error) {
	if len(params) < 3 {
		return ANCSetting{}, fmt.Errorf("anc setting: need 3 bytes, got %d", len(params))
	}
	return ANCSetting{
		Mode:       params[0],
		SubScene:   params[1],
		NoiseValue: params[2],
	}, nil
}

type PowerManager struct {
	PowerOffTime uint16
	CurrentTime  uint16
}

func ParsePowerManager(params []byte) (PowerManager, error) {
	if len(params) < 4 {
		return PowerManager{}, fmt.Errorf("power manager: need 4 bytes, got %d", len(params))
	}
	return PowerManager{
		PowerOffTime: binary.LittleEndian.Uint16(params[0:2]),
		CurrentTime:  binary.LittleEndian.Uint16(params[2:4]),
	}, nil
}

type Volume struct {
	Left  byte
	Right byte
	Max   byte
}

func ParseVolume(params []byte) (Volume, error) {
	if len(params) < 3 {
		return Volume{}, fmt.Errorf("volume: need 3 bytes, got %d", len(params))
	}
	return Volume{
		Left:  params[0],
		Right: params[1],
		Max:   params[2],
	}, nil
}

type ToneVolume struct {
	Volume    byte
	MaxVolume byte
}

func ParseToneVolume(params []byte) (ToneVolume, error) {
	if len(params) < 1 {
		return ToneVolume{}, fmt.Errorf("tone volume: need at least 1 byte, got %d", len(params))
	}
	tv := ToneVolume{Volume: params[0], MaxVolume: 100}
	if len(params) >= 2 {
		tv.MaxVolume = params[1]
	}
	return tv, nil
}

type WearingDetection struct {
	Enabled    bool
	MusicIndex byte
	ANCIndex   byte
	ToneEnable bool
	HasTone    bool
}

func ParseWearingDetection(params []byte) (WearingDetection, error) {
	if len(params) < 3 {
		return WearingDetection{}, fmt.Errorf("wearing detection: need at least 3 bytes, got %d", len(params))
	}
	wd := WearingDetection{
		Enabled:    params[0] == 0x01,
		MusicIndex: params[1],
		ANCIndex:   params[2],
	}
	if len(params) >= 4 {
		wd.HasTone = true
		wd.ToneEnable = params[3] == 0x01
	}
	return wd, nil
}

type EarTipFitResult struct {
	Status      byte
	LeftResult  byte
	RightResult byte
}

const (
	EarTipFitReady  byte = 0x00
	EarTipFitResult_ byte = 0x02
)

func ParseEarTipFit(params []byte) (EarTipFitResult, error) {
	if len(params) < 3 {
		return EarTipFitResult{}, fmt.Errorf("ear tip fit: need 3 bytes, got %d", len(params))
	}
	return EarTipFitResult{
		Status:      params[0],
		LeftResult:  params[1],
		RightResult: params[2],
	}, nil
}

type EQBand struct {
	Freq     uint16
	Gain     int16
	Q        uint16
	BandType byte
}

type EQParams struct {
	EQType     byte
	MasterGain int16
	Bands      []EQBand
}

func ParseEQV1(params []byte) (EQParams, error) {
	if len(params) < 3 {
		return EQParams{}, fmt.Errorf("eq v1: need at least 3 bytes, got %d", len(params))
	}
	eq := EQParams{
		EQType:     params[0],
		MasterGain: int16(binary.LittleEndian.Uint16(params[1:3])),
	}
	data := params[3:]
	for len(data) >= 6 {
		eq.Bands = append(eq.Bands, EQBand{
			Freq: binary.LittleEndian.Uint16(data[0:2]),
			Gain: int16(binary.LittleEndian.Uint16(data[2:4])),
			Q:    binary.LittleEndian.Uint16(data[4:6]),
		})
		data = data[6:]
	}
	return eq, nil
}

func ParseEQV2(params []byte) (EQParams, error) {
	if len(params) < 3 {
		return EQParams{}, fmt.Errorf("eq v2: need at least 3 bytes, got %d", len(params))
	}
	eq := EQParams{
		EQType:     params[0],
		MasterGain: int16(binary.LittleEndian.Uint16(params[1:3])),
	}
	data := params[3:]
	for len(data) >= 7 {
		eq.Bands = append(eq.Bands, EQBand{
			Freq:     binary.LittleEndian.Uint16(data[0:2]),
			Gain:     int16(binary.LittleEndian.Uint16(data[2:4])),
			Q:        binary.LittleEndian.Uint16(data[4:6]),
			BandType: data[6],
		})
		data = data[7:]
	}
	return eq, nil
}

type KeyMapping struct {
	Key  byte
	Func byte
}

func ParseKeyFunction(params []byte) ([]KeyMapping, error) {
	if len(params)%2 != 0 {
		return nil, fmt.Errorf("key function: odd byte count %d", len(params))
	}
	mappings := make([]KeyMapping, 0, len(params)/2)
	for i := 0; i+1 < len(params); i += 2 {
		mappings = append(mappings, KeyMapping{Key: params[i], Func: params[i+1]})
	}
	return mappings, nil
}

type LEDEffect struct {
	Speed       byte
	Brightness  byte
	EffectIndex byte
	Colors      []color.RGBA
}

func ParseLEDEffect(params []byte) (LEDEffect, error) {
	if len(params) < 3 {
		return LEDEffect{}, fmt.Errorf("led effect: need at least 3 bytes, got %d", len(params))
	}
	effect := LEDEffect{
		Speed:       params[0],
		Brightness:  params[1],
		EffectIndex: params[2],
	}
	data := params[3:]
	for len(data) >= 3 {
		effect.Colors = append(effect.Colors, color.RGBA{R: data[0], G: data[1], B: data[2], A: 0xff})
		data = data[3:]
	}
	return effect, nil
}

type Alarm struct {
	AlarmID byte
	Enabled bool
	Hour    byte
	Minute  byte
	Cycle   byte
}

func ParseAlarmList(params []byte) ([]Alarm, error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("alarm list: need at least 1 byte, got %d", len(params))
	}
	data := params[1:]
	var alarms []Alarm
	for len(data) >= 6 {
		alarms = append(alarms, Alarm{
			AlarmID: data[0],
			Enabled: data[1] == 0x01,
			Hour:    data[2],
			Minute:  data[3],
			Cycle:   data[4],
		})
		data = data[6:]
	}
	return alarms, nil
}

type MusicStatus struct {
	MusicID   uint32
	IsPlaying bool
	PlayMode  byte
}

func ParseMusicStatus(params []byte) (MusicStatus, error) {
	if len(params) < 6 {
		return MusicStatus{}, fmt.Errorf("music status: need 6 bytes, got %d", len(params))
	}
	return MusicStatus{
		MusicID:   binary.LittleEndian.Uint32(params[0:4]),
		IsPlaying: params[4] == 0x01,
		PlayMode:  params[5],
	}, nil
}

type MusicFile struct {
	MusicID uint32
	Total   uint16
}

func ParseMusicInfo(params []byte) (byte, []MusicFile, error) {
	if len(params) < 1 {
		return 0, nil, fmt.Errorf("music info: need at least 1 byte, got %d", len(params))
	}
	startToneID := params[0]
	data := params[1:]
	var files []MusicFile
	for len(data) >= 6 {
		files = append(files, MusicFile{
			MusicID: binary.LittleEndian.Uint32(data[0:4]),
			Total:   binary.LittleEndian.Uint16(data[4:6]),
		})
		data = data[6:]
	}
	return startToneID, files, nil
}
