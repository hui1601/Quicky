package command

func NewAudioLangCommand(lang string) *Command {
	return NewCommand(0x19, []byte(lang))
}
