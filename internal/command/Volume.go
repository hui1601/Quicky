package command

func NewVolumeCommand(left, right byte) *Command {
	return NewCommand(0x08, []byte{left, right, 0x00})
}
