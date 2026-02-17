package command

func NewAICommand(action byte) *Command {
	return NewCommand(0x43, []byte{action})
}
