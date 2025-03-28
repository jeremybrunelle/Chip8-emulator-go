package gochip8

import (
	"errors"
)

var font_set = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type stack []uint16

func (s stack) Push(v uint16) stack {
	return append(s, v)
}

func (s stack) Pop() (uint16, error) {
	l := len(s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	result := s[l-1]
	s = s[:l-1]

	return result, nil
}

type Chip8 struct {
	memory      [4096]uint8   // 4kb of ram
	display     [64][32]uint8 //64 x 32 ram
	stack       [16]stack     //program stack
	key         [16]uint8
	delay_timer uint8
	sound_timer uint8

	vx [16]uint8 //cpu v registers
	oc uint16
	pc uint16 // program counter, points to current instruction in memory
	iv uint16 // indexes memory
	sp uint16 //stack pointer

	doDraw bool
}

func Init() Chip8 {
	emu := Chip8{
		pc:     0x200,
		doDraw: true,
	}

	for i := 0; i < len(font_set); i++ {
		emu.memory[i] = font_set[i]
	}

	return emu

}

func (c *Chip8) Draw() bool {
	sd := c.doDraw
	c.doDraw = false
	return sd

}

func (c *Chip8) Cycle() {

	//ideally this is the fetch stage
	c.oc = (uint16(c.memory[c.pc]) | uint16(c.memory[c.pc+1]))
	c.pc += 2

	//ideally this is the decode stage
	switch c.oc & 0xF000 {
	case 0x0000:
		switch c.oc {
		case 0x00E0: //clear screen
			for x := 0; x < 64; x++ {
				for y := 0; y < 32; y++ {
					c.display[x][y] = 0
				}
			}
		}
		break
	case 0x1000: //jump to 0x1NNN
		c.pc = c.oc & 0x0FFF
		break
	case 0x6000:
		x := (c.oc & 0x0F00) >> 8
		nn := uint8(c.oc & 0x00FF)
		c.vx[x] = nn
		break
	case 0x7000:
		x := (c.oc & 0x0F00) >> 8
		nn := uint8(c.oc & 0x00FF)

		c.vx[x] = c.vx[x] + nn
		break
	case 0xA000:
		c.iv = uint16(c.oc & 0x0FFF)
		break
	case 0xD000:
		n_pixels := c.oc & 0x000F
		x := c.vx[(c.oc&0x0F00)>>8]
		y := c.vx[(c.oc&0x00F0)>>4]
		c.vx[0xF] = 0
		var i uint16 = 0
		var j uint16 = 0
		for i = 0; i < n_pixels; i++ {
			pixel := c.memory[c.iv+i]
			for j = 0; j < 8; j++ {
				if (pixel & (0x80 >> j)) != 0 {
					if (c.display[(y + uint8(i))][x+uint8(j)]) == 1 {
						c.vx[0xF] = 1
					}
					c.display[(y + uint8(i))][x+uint8(j)] ^= 1
				}
			}

		}
		c.doDraw = true

		break

	}

}
