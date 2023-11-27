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
	ErrInvalidPacket := errors.New("invalid packet")
	if len(packet) < 3 {
		return nil, ErrInvalidPacket
	}
	if packet[0] != 0xff {
		return nil, ErrInvalidPacket
	}
	bodyLength := int(packet[1])
	if len(packet) != bodyLength+2 {
		return nil, ErrInvalidPacket
	}
	operationCode := int(packet[2])
	parameters := packet[3:]
	return &Command{
		OperationCode: operationCode,
		Parameters:    parameters,
	}, nil
}
