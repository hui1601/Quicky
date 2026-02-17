package command

func NewStandbyCommand(state byte) *Command {
	return NewCommand(0x1f, []byte{state})
}
