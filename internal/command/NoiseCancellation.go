// Package command
package command

type NoiseCancellationMode byte

const (
	NoiseCancellationModeOff NoiseCancellationMode = iota
	NoiseCancellationModeSilent
	NoiseCancellationModeWork
	NoiseCancellationModeNoisy
	NoseCancellationModePassThrough
)

func (n NoiseCancellationMode) String() string {
	switch n {
	case NoiseCancellationModeOff:
		return "Off"
	case NoiseCancellationModeSilent:
		return "Silent"
	case NoiseCancellationModeWork:
		return "Work"
	case NoiseCancellationModeNoisy:
		return "Noisy"
	case NoseCancellationModePassThrough:
		return "PassThrough"
	default:
		return "Unknown"
	}
}

func NewNoiseCancellationCommand(mode NoiseCancellationMode, level byte) *Command {
	// high 4 bits: mode
	// low 4 bits: level
	// => 1 byte
	// off => 0x0002
	arg := byte(0)
	// 0x02: silent, 0x03: work, 0x04: noisy, 0x0a: pass through
	switch mode {
	case NoiseCancellationModeSilent:
		arg = 0x02
	case NoiseCancellationModeWork:
		arg = 0x03
	case NoiseCancellationModeNoisy:
		arg = 0x04
	case NoseCancellationModePassThrough:
		arg = 0x0a
	default:
	}
	arg = arg << 4
	// 0x01: low, 0x02: medium, 0x03: high(0x00 might be off)
	// in pass through mode, 0x07 is the highest level(1-6 normal, 7 increase human voice)
	arg = arg | func(mode NoiseCancellationMode, level byte) byte {
		if mode == NoseCancellationModePassThrough {
			return level%6 + 1
		} else if mode == NoiseCancellationModeOff {
			return 0x2
		}
		return level%3 + 1
	}(mode, level)
	return NewCommand(0x0c, []byte{arg})
}
