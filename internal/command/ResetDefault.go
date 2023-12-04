package command

func NewResetDefaultCommand() *Command {
	return NewCommand(0x01, []byte{})
}
