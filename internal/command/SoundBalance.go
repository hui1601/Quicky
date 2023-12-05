package command

func NewSoundBalanceCommand(left int32) *Command {
	return NewCommand(0x16, []byte{byte(left & 0xff)})
}
