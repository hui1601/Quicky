package command

func NewInEarTestCommand(isOn bool) *Command {
	var arg byte = 0x02
	if isOn {
		arg = 0x01
	}
	return NewCommand(0x06, []byte{arg})
}
