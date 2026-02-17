package command

import "encoding/binary"

type EQBand struct {
	Freq     uint16
	Gain     int16
	Q        uint16
	BandType byte
}

func clampGain(gain int16) int16 {
	if gain > 1270 {
		return 1270
	}
	if gain < -1270 {
		return -1270
	}
	return gain
}

func NewEQV1Command(eqIndex byte, masterGain int16, bands []EQBand) *Command {
	params := []byte{eqIndex}
	mg := make([]byte, 2)
	binary.LittleEndian.PutUint16(mg, uint16(masterGain))
	params = append(params, mg...)
	for _, b := range bands {
		buf := make([]byte, 6)
		binary.LittleEndian.PutUint16(buf[0:2], b.Freq)
		binary.LittleEndian.PutUint16(buf[2:4], uint16(clampGain(b.Gain)))
		binary.LittleEndian.PutUint16(buf[4:6], b.Q)
		params = append(params, buf...)
	}
	return NewCommand(0x20, params)
}

func newEQV2Params(eqIndex byte, masterGain int16, bands []EQBand) []byte {
	params := []byte{eqIndex}
	mg := make([]byte, 2)
	binary.LittleEndian.PutUint16(mg, uint16(masterGain))
	params = append(params, mg...)
	for _, b := range bands {
		buf := make([]byte, 7)
		binary.LittleEndian.PutUint16(buf[0:2], b.Freq)
		binary.LittleEndian.PutUint16(buf[2:4], uint16(clampGain(b.Gain)))
		binary.LittleEndian.PutUint16(buf[4:6], b.Q)
		buf[6] = b.BandType
		params = append(params, buf...)
	}
	return params
}

func NewEQV2Command(eqIndex byte, masterGain int16, bands []EQBand) *Command {
	return NewCommand(0x22, newEQV2Params(eqIndex, masterGain, bands))
}

func NewEQLeftCommand(eqIndex byte, masterGain int16, bands []EQBand) *Command {
	return NewCommand(0x46, newEQV2Params(eqIndex, masterGain, bands))
}

func NewEQRightCommand(eqIndex byte, masterGain int16, bands []EQBand) *Command {
	return NewCommand(0x47, newEQV2Params(eqIndex, masterGain, bands))
}

func NewEQDirectData(eqType byte, data []byte) []byte {
	result := []byte{eqType}
	result = append(result, data...)
	return result
}
