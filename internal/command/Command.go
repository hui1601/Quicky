package command

import "errors"

type Command struct {
	OperationCode byte
	Parameters    []byte
}

func NewCommand(opcode byte, parameters []byte) *Command {
	return &Command{
		OperationCode: opcode,
		Parameters:    parameters,
	}
}

func (c *Command) PackPacket() []byte {
	body := []byte{c.OperationCode, byte(len(c.Parameters))}
	body = append(body, c.Parameters...)
	packet := []byte{0xff, byte(len(body))}
	packet = append(packet, body...)
	return packet
}

func ParsePacket(packet []byte) ([]Command, error) {
	if len(packet) < 4 || packet[0] != 0xff {
		return nil, errors.New("invalid packet: too short or missing SOF")
	}
	bodyLen := int(packet[1])
	if bodyLen+2 != len(packet) {
		return nil, errors.New("invalid packet: body length mismatch")
	}

	var commands []Command
	offset := 2
	for offset < len(packet) {
		if offset+2 > len(packet) {
			return nil, errors.New("invalid packet: truncated command block")
		}
		cmdId := packet[offset]
		paramLen := int(packet[offset+1])
		offset += 2
		if offset+paramLen > len(packet) {
			return nil, errors.New("invalid packet: truncated parameters")
		}
		params := make([]byte, paramLen)
		copy(params, packet[offset:offset+paramLen])
		commands = append(commands, Command{OperationCode: cmdId, Parameters: params})
		offset += paramLen
	}

	if len(commands) == 0 {
		return nil, errors.New("invalid packet: no commands found")
	}
	return commands, nil
}
