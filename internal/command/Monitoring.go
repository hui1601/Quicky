package command

func NewMonitoringCommand(value byte) *Command {
	return NewCommand(0x0a, []byte{value})
}
