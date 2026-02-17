package command

func NewLEDSwitchCommand(enable bool) *Command {
	arg := byte(0x02)
	if enable {
		arg = 0x01
	}
	return NewCommand(0x35, []byte{arg})
}
