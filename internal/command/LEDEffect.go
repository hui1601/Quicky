package command

import "image/color"

func NewLEDEffectCommand(speed, brightness, effectIndex byte, colors []color.RGBA) *Command {
	params := []byte{speed, brightness, effectIndex}
	for _, c := range colors {
		params = append(params, c.R, c.G, c.B)
	}
	return NewCommand(0x36, params)
}
