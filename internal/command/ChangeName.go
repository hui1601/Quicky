package command

func NewChangeNameCommand(name string) *Command {
	return NewCommand(0x18, []byte(name))
}
