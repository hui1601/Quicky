package command

func NewEnvAdaptationCommand(state byte) *Command {
	return NewCommand(0x32, []byte{state})
}
