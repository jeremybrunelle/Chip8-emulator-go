package emulator

import (
	"fmt"
	"os"
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

type Chip8 struct {
	memory      [4096]uint8   // 4kb of ram
	display     [32][64]uint8 //64 x 32 ram
	stack       [16]uint16    //program stack
	key         [16]uint8
	delay_timer uint8
	sound_timer uint8

	vx [16]uint8 //cpu v registers
	oc uint16
	pc uint16 // program counter, points to current instruction in memory
	iv uint16 // indexes memory
	sp int32  //stack pointer

	doDraw bool
}

func Init() Chip8 {
	emu := Chip8{
		pc:     0x200,
		sp:     0,
		doDraw: true,
	}

	for i := 0; i < len(font_set); i++ {
		emu.memory[i] = font_set[i]
	}

	for i := 0; i < len(emu.stack); i++ {
		emu.stack[i] = 0x0
	}

	return emu

}

func (c *Chip8) Buffer() [32][64]uint8 {
	return c.display
}

func (c *Chip8) Draw() bool {
	sd := c.doDraw
	c.doDraw = false
	return sd

}

func (c *Chip8) Push(inst uint16) {

	if c.stack[c.sp] == 0 {
		c.stack[c.sp] = inst
	} else {
		c.sp += 1
		c.stack[c.sp] = inst
	}
}

func (c *Chip8) Pop() uint16 {
	return c.stack[c.sp]
}

func (c *Chip8) Key(num uint8, down bool) {
	if down {
		c.key[num] = 1
	} else {
		c.key[num] = 0
	}
}

func (c *Chip8) Cycle() {

	//ideally this is the fetch stage
	c.oc = (uint16(c.memory[c.pc]) << 8) | uint16(c.memory[c.pc+1])
	c.pc += 2

	//ideally this is the decode stage
	switch c.oc & 0xF000 {
	case 0x0000:
		switch c.oc {
		case 0x00E0: //clear screen
			for i := 0; i < len(c.display); i++ {
				for j := 0; j < len(c.display[i]); j++ {
					c.display[i][j] = 0x0
				}
			}
			c.doDraw = true

		case 0x00EE: //Subroutines
			c.pc = c.Pop()
		}

	case 0x1000: //jump to 0x1NNN
		c.pc = c.oc & 0x0FFF

	case 0x2000:
		c.Push(c.pc)
		c.pc = c.oc & 0x0FFF

	case 0x3000:
		x := (c.oc & 0x0F00) >> 8
		nn := uint8(c.oc & 0x00FF)
		if c.vx[x] == nn {
			c.pc += 2
		}

	case 0x4000:
		x := (c.oc & 0x0F00) >> 8
		nn := uint8(c.oc & 0x00FF)
		if c.vx[x] != nn {
			c.pc += 2
		}
	case 0x5000:
		x := (c.oc & 0x0F00) >> 8
		y := (c.oc & 0x00F0) >> 4
		if c.vx[x] == c.vx[y] {
			c.pc += 2
		}
	case 0x6000: //Set
		x := (c.oc & 0x0F00) >> 8
		nn := uint8(c.oc & 0x00FF)
		c.vx[x] = nn

	case 0x7000: //add
		x := (c.oc & 0x0F00) >> 8
		nn := uint8(c.oc & 0x00FF)

		c.vx[x] = c.vx[x] + nn

	case 0x8000:
		x := (c.oc & 0x0F00) >> 8
		y := (c.oc & 0x00F0) >> 4
		switch c.oc & 0x000F {
		case 0x0000:

			c.vx[x] = c.vx[y]
		case 0x0001:
			c.vx[x] = c.vx[x] | c.vx[y]
		case 0x0002:
			c.vx[x] = c.vx[x] & c.vx[y]
		case 0x0003:
			c.vx[x] = c.vx[x] ^ c.vx[y]
		case 0x0004:
			if c.vx[x]+c.vx[y] > 255 {
				c.vx[0xF] = uint8(1)
			} else {
				c.vx[x] = c.vx[x] + c.vx[y]
				c.vx[0xF] = uint8(1)
			}
		case 0x0005:
			if c.vx[x] > c.vx[y] {
				c.vx[0xF] = uint8(1)
				c.vx[x] = c.vx[x] - c.vx[y]
			} else {
				c.vx[0xF] = uint8(0)
			}
		case 0x0006: //need to finish
			c.vx[0xF] = c.vx[(c.oc&0x0F00)>>8] & 0x1
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] >> 1
		case 0x0007:
			if c.vx[y] > c.vx[x] {
				c.vx[0xF] = uint8(1)
				c.vx[y] = c.vx[y] - c.vx[x]
			} else {
				c.vx[0xF] = uint8(0)
			}
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x00F0)>>4] - c.vx[(c.oc&0x0F00)>>8]

		}

	case 0x9000:
		x := (c.oc & 0x0F00) >> 8
		y := (c.oc & 0x00F0) >> 4
		if c.vx[x] != c.vx[y] {
			c.pc += 2
		}

	case 0xA000: //set index
		c.iv = uint16(c.oc & 0x0FFF)

	case 0xD000: //draw
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

	case 0xE000:
		switch c.oc & 0x00FF {
		case 0x009E:
			if c.key[c.vx[(c.oc&0x0F00)>>8]] == 1 {
				c.pc += 2
			}
		case c.oc & 0x00FF:
			if c.key[c.vx[(c.oc&0x0F00)>>8]] != 1 {
				c.pc += 2
			}

		}
	case 0xF000:
		switch c.oc & 0x00FF {
		case 0x0007:
			vx := (c.oc & 0x0F00) >> 8
			c.vx[vx] = c.delay_timer
		case 0x000A:
			c.pc -= 2
			key := (c.oc & 0x0F00) >> 8

			if c.key[key] == 1 {
				c.pc += 2
				c.vx[key] = uint8(key)
			}

		case 0x0015:
			vx := (c.oc & 0x0F00) >> 8
			c.delay_timer = c.vx[vx]
		case 0x0018:
			vx := (c.oc & 0x0F00) >> 8
			c.sound_timer = c.vx[vx]
		case 0x001E:
			if c.iv+uint16(c.vx[(c.oc&0x0F00)>>8]) > 0xFFF {
				c.vx[0xF] = 1
			} else {
				c.vx[0xF] = 0
			}
			c.iv += uint16(c.vx[(c.oc&0x0F00)>>8])
		case 0x0029:
			c.iv = uint16(c.vx[(c.oc&0xF00)>>8]) * 0x5

		case 0x0033:
			c.memory[c.iv] = c.vx[(c.oc&0x0F00)>>8] / 100
			c.memory[c.iv+1] = (c.vx[(c.oc&0x0F00)>>8] / 10) % 10
			c.memory[c.iv+2] = (c.vx[(c.oc&0x0F00)>>8] / 100) % 10

		case 0x0055:
			x := int((c.oc & 0x0F00) >> 8)
			for i := 0; i <= x+1; i++ {
				c.memory[c.iv+uint16(i)] = c.vx[i]
			}
			c.iv = ((c.oc & 0x0F00) >> 8) + 1
		case 0x0065:
			for i := 0; i < int((c.oc&0x0F00)>>8)+1; i++ {
				c.vx[i] = c.memory[c.iv+uint16(i)]
			}
			c.iv = ((c.oc & 0x0F00) >> 8) + 1

		}
	}
	if c.delay_timer > 0 {
		c.delay_timer = c.delay_timer - 1
	}

}

func (c *Chip8) Load(fileName string) error {

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}

	defer file.Close()
	fileInfo, FileStatErr := file.Stat()
	if FileStatErr != nil {
		return FileStatErr
	}

	if int64(len(c.memory)-512) < fileInfo.Size() {
		return fmt.Errorf("program size bigger than memory")
	}

	fileBuffer := make([]byte, fileInfo.Size())
	if _, ReadErr := file.Read(fileBuffer); ReadErr != nil {
		return ReadErr
	}

	for i := 0; i < len(fileBuffer); i++ {
		c.memory[i+512] = fileBuffer[i]
	}

	return nil
}
