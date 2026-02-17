package command

func NewTakePhotoCommand(action byte) *Command {
	return NewCommand(0x1e, []byte{action})
}
