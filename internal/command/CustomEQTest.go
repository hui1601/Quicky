package command

func NewCustomEQTestCommand(state byte) *Command {
	return NewCommand(0x45, []byte{state})
}
