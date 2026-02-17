package command

func NewPlayModeCommand(mode byte) *Command {
	return NewCommand(0x37, []byte{mode})
}
