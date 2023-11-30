package command

// MusicControlAction is not used in the app, but command code is present in the app.
type MusicControlAction int16

func NewMusicControlCommand(action MusicControlAction) *Command {
	return NewCommand(0x04, []byte{byte(action)})
}
