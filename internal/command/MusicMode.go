package command

func NewMusicModeCommand(mode byte) *Command {
	return NewCommand(0x2e, []byte{mode})
}
