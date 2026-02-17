package command

func NewEarTipFitCommand(left, right byte) *Command {
	return NewCommand(0x11, []byte{left, right})
}

func NewEarTipFitStartCommand() *Command {
	return NewEarTipFitCommand(0x01, 0x01)
}

func NewEarTipFitStopCommand() *Command {
	return NewEarTipFitCommand(0x02, 0x02)
}
