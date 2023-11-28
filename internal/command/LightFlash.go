package command

func NewLightFlashCommand(isOn bool) *Command {
	// 0x00: off, 0x01: on?
	var arg byte = 0x00
	if isOn {
		arg = 0x01
	}
	return NewCommand(0x05, []byte{arg})
}
