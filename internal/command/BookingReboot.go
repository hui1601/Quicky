package command

func NewBookingRebootCommand(reserveTime int32, currentTime int32) *Command {
	return NewCommand(0x14, []byte{byte(reserveTime & 0xff), byte((reserveTime >> 8) & 0xff), byte(currentTime & 0xff), byte((currentTime >> 8) & 0xff)})
}
