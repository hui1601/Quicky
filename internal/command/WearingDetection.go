package command

func NewWearingDetectionCommand(enable bool, musicIndex, ancIndex byte) *Command {
	arg := byte(0x02)
	if enable {
		arg = 0x01
	}
	return NewCommand(0x2c, []byte{arg, musicIndex, ancIndex})
}

func NewWearingDetectionV2Command(enable bool, musicIndex, ancIndex byte, toneEnable bool) *Command {
	arg := byte(0x02)
	if enable {
		arg = 0x01
	}
	tone := byte(0x02)
	if toneEnable {
		tone = 0x01
	}
	return NewCommand(0x2c, []byte{arg, musicIndex, ancIndex, tone})
}
