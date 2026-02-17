package command

import "encoding/binary"

func NewMusicStatusCommand(musicID uint32, isPlaying bool, playMode byte) *Command {
	params := make([]byte, 6)
	binary.LittleEndian.PutUint32(params[0:4], musicID)
	if isPlaying {
		params[4] = 0x01
	}
	params[5] = playMode
	return NewCommand(0x3a, params)
}
