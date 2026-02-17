package command

func NewFocusModeCommand(enable bool) *Command {
	arg := byte(0x02)
	if enable {
		arg = 0x01
	}
	return NewCommand(0x39, []byte{arg})
}
