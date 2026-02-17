package command

type NoiseCancelMode byte

const (
	NoiseCancelOff          NoiseCancelMode = 0x00
	NoiseCancelANC          NoiseCancelMode = 0x01
	NoiseCancelOutdoor      NoiseCancelMode = 0x02
	NoiseCancelTransparency NoiseCancelMode = 0x03
)

func NewNoiseCancelModeCommand(mode NoiseCancelMode) *Command {
	return NewCommand(0x0c, []byte{byte(mode)})
}
