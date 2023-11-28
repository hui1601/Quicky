package command

func NewTestModeCommand(isOn bool) *Command {
	var arg byte = 0x02
	if isOn {
		arg = 0x01
	}
	return NewCommand(0x0d, []byte{arg})
}
