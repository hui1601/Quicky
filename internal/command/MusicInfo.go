package command

import "encoding/binary"

type MusicFile struct {
	MusicID uint32
	Total   uint16
}

func NewMusicInfoCommand(startToneID byte, files []MusicFile) *Command {
	params := []byte{startToneID}
	for _, f := range files {
		buf := make([]byte, 6)
		binary.LittleEndian.PutUint32(buf[0:4], f.MusicID)
		binary.LittleEndian.PutUint16(buf[4:6], f.Total)
		params = append(params, buf...)
	}
	return NewCommand(0x3b, params)
}
