package command

func NewSpatialAudioCommand(enable bool) *Command {
	arg := byte(0x02)
	if enable {
		arg = 0x01
	}
	return NewCommand(0x2d, []byte{arg})
}
