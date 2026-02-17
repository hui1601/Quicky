package command

func NewRequestDataCommand(cmdID byte) *Command {
	return NewCommand(0xfe, []byte{cmdID})
}
