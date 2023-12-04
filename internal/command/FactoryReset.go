package command

func NewFactoryResetCommand() *Command {
	return NewCommand(0x03, []byte{})
}
