package command

func NewGameConfigCommand(config byte) *Command {
	return NewCommand(0x4a, []byte{config})
}
