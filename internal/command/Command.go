package command

import "errors"

type Command struct {
	OperationCode int
	Parameters    []byte
}

func NewCommand(operationCode int, parameters []byte) *Command {
	return &Command{
		OperationCode: operationCode,
		Parameters:    parameters,
	}
}

func (c *Command) PackPacket() []byte {
	var packet []byte
	var body []byte
	body = append(body, byte(c.OperationCode))
	body = append(body, byte(len(c.Parameters)))
	body = append(body, c.Parameters...)
	// packet prefix
	packet = append(packet, 0xff)
	// body length
	packet = append(packet, byte(len(body)))
	// body
	packet = append(packet, body...)
	return packet
}

func ParsePacket(packet []byte) (*Command, error) {
	if len(packet) < 3 || packet[0] != 0xff || len(packet) != int(packet[1])+2 {
		return nil, errors.New("invalid packet")
	}
	return &Command{
		OperationCode: int(packet[2]),
		Parameters:    packet[4:],
	}, nil
}
