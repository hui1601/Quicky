package command

func NewTonePlayCommand(toneID byte) *Command {
	return NewCommand(0x3d, []byte{toneID})
}
