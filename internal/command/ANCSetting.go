package command

func NewANCSettingCommand(mode, subScene, noiseValue byte) *Command {
	return NewCommand(0x17, []byte{mode, subScene, noiseValue})
}
