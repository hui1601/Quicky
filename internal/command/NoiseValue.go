package command

func NewNoiseValueCommand(value byte) *Command {
	return NewCommand(0x07, []byte{value})
}
