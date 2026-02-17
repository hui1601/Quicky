package command

import "time"

func NewSyncTimeCommand(t time.Time) *Command {
	year := byte(t.Year() % 100)
	month := byte(t.Month())
	day := byte(t.Day())
	hour := byte(t.Hour())
	minute := byte(t.Minute())
	second := byte(t.Second())
	dow := byte(1 << (t.Weekday())) // Sunday=0 → bit 0
	return NewCommand(0x3e, []byte{year, month, day, hour, minute, second, dow})
}
