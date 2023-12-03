package command

// NewLowLatencyCommand creates a new low latency command
// maybe not used?
func NewLowLatencyCommand(isOn bool) *Command {
	var arg byte = 0x02
	if isOn {
		arg = 0x01
	}
	return NewCommand(0x09, []byte{arg})
}
