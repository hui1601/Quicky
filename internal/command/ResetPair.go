package command

func NewResetPairCommand() *Command {
	return NewCommand(0x02, []byte{})
}
