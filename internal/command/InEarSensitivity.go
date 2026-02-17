package command

func NewInEarSensitivityCommand(level byte) *Command {
	return NewCommand(0x48, []byte{level})
}
