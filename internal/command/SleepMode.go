package command

func NewSleepModeCommand(enable bool) *Command {
	arg := byte(0x02)
	if enable {
		arg = 0x01
	}
	return NewCommand(0x10, []byte{arg})
}
