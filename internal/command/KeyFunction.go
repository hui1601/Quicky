package command

type KeyID byte

const (
	KeyMusicLeftSingle  KeyID = 0x01
	KeyMusicRightSingle KeyID = 0x02
	KeyMusicLeftDouble  KeyID = 0x03
	KeyMusicRightDouble KeyID = 0x04
	KeyMusicLeftTriple  KeyID = 0x05
	KeyMusicRightTriple KeyID = 0x06
	KeyMusicLeftQuad    KeyID = 0x07
	KeyMusicRightQuad   KeyID = 0x08
	KeyMusicLeftLong    KeyID = 0x09
	KeyMusicRightLong   KeyID = 0x0a

	KeyVoiceLeftSingle  KeyID = 0x15
	KeyVoiceRightSingle KeyID = 0x16
	KeyVoiceLeftDouble  KeyID = 0x17
	KeyVoiceRightDouble KeyID = 0x18
	KeyVoiceLeftTriple  KeyID = 0x19
	KeyVoiceRightTriple KeyID = 0x1a
	KeyVoiceLeftQuad    KeyID = 0x1b
	KeyVoiceRightQuad   KeyID = 0x1c
	KeyVoiceLeftLong    KeyID = 0x1d
	KeyVoiceRightLong   KeyID = 0x1e
)

type FuncID byte

const (
	FuncNone           FuncID = 0x00
	FuncPlayPause      FuncID = 0x01
	FuncPrevious       FuncID = 0x02
	FuncNext           FuncID = 0x03
	FuncVoiceAssistant FuncID = 0x04
	FuncVolumeUp       FuncID = 0x05
	FuncVolumeDown     FuncID = 0x06
	FuncGameMode       FuncID = 0x07
	FuncAnswerCall     FuncID = 0x08
	FuncRejectCall     FuncID = 0x09
	FuncHoldCall       FuncID = 0x0a
	FuncRedial         FuncID = 0x0b
)

type KeyMapping struct {
	Key  KeyID
	Func FuncID
}

func NewKeyFunctionDirectData(mappings []KeyMapping) []byte {
	data := make([]byte, 0, len(mappings)*2)
	for _, m := range mappings {
		data = append(data, byte(m.Key), byte(m.Func))
	}
	return data
}
