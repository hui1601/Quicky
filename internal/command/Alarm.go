package command

type AlarmOperation byte

const (
	AlarmAdd    AlarmOperation = 0x01
	AlarmDelete AlarmOperation = 0x02
	AlarmEdit   AlarmOperation = 0x03
)

func NewAlarmAddCommand(alarmID, hour, minute, cycle byte) *Command {
	return NewCommand(0x3f, []byte{byte(AlarmAdd), alarmID, 0x01, hour, minute, cycle, 0x05})
}

func NewAlarmDeleteCommand(alarmID byte) *Command {
	return NewCommand(0x3f, []byte{byte(AlarmDelete), alarmID, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func NewAlarmEditCommand(alarmID byte, enable bool, hour, minute, cycle, index byte) *Command {
	en := byte(0x00)
	if enable {
		en = 0x01
	}
	return NewCommand(0x3f, []byte{byte(AlarmEdit), alarmID, en, hour, minute, cycle, index})
}
