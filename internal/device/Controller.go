package device

import (
	"Quicky/internal/command"
	"Quicky/internal/constant"
)

func (c *Client) SendCommand(command *command.Command) error {
	packet := command.PackPacket()
	char, err := c.Device.GetCharByUUID(constant.CommandCharacteristicUUID)
	if err != nil {
		return err
	}
	err = char.WriteValue(packet, nil)
	if err != nil {
		return err
	}
	return nil
}
