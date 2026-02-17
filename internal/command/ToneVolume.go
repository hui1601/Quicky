package command

func NewToneVolumeCommand(volume byte) *Command {
	return NewCommand(0x1d, []byte{volume})
}
