package main

import (
	"fmt"
	"os"

	"chip8-emulator/emulator"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	CHIP8_WIDTH  int32 = 64 * 25
	CHIP8_HEIGHT int32 = 32 * 25
	windowTitle        = "CHIP-8"
)

func main() {
	if len(os.Args) < 2 {
		panic("wrong args")
	}

	filepath := os.Args[1]

	chip8 := emulator.Init()
	if err := chip8.Load(filepath); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	window, err := sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		CHIP8_WIDTH, CHIP8_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	defer window.Destroy()

	r, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	defer r.Destroy()
	running := true
	for running {
		chip8.Cycle()

		if chip8.Draw() {
			r.SetDrawColor(0, 0, 0, 0)
			r.Clear()

			vector := chip8.Buffer()
			for j := 0; j < len(vector); j++ {
				for i := 0; i < len(vector[j]); i++ {
					if vector[j][i] != 0 {
						r.SetDrawColor(255, 255, 255, 255)
					} else {
						r.SetDrawColor(0, 0, 0, 0)
					}
					r.FillRect(&sdl.Rect{
						Y: int32(j) * 25,
						X: int32(i) * 25,
						W: 25,
						H: 25,
					})
				}
			}
			r.Present()

		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
				os.Exit(0)

			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYUP {
					switch t.Keysym.Sym {
					case sdl.K_1:
						chip8.Key(0x1, false)
					case sdl.K_2:
						chip8.Key(0x2, false)
					case sdl.K_3:
						chip8.Key(0x3, false)
					case sdl.K_4:
						chip8.Key(0xC, false)
					case sdl.K_q:
						chip8.Key(0x4, false)
					case sdl.K_w:
						chip8.Key(0x5, false)
					case sdl.K_e:
						chip8.Key(0x6, false)
					case sdl.K_r:
						chip8.Key(0xD, false)
					case sdl.K_a:
						chip8.Key(0x7, false)
					case sdl.K_s:
						chip8.Key(0x8, false)
					case sdl.K_d:
						chip8.Key(0x9, false)
					case sdl.K_f:
						chip8.Key(0xE, false)
					case sdl.K_z:
						chip8.Key(0xA, false)
					case sdl.K_x:
						chip8.Key(0x0, false)
					case sdl.K_c:
						chip8.Key(0xB, false)
					case sdl.K_v:
						chip8.Key(0xF, false)

					}
				} else if t.Type == sdl.KEYDOWN {
					switch t.Keysym.Sym {
					case sdl.K_1:
						chip8.Key(0x1, true)
					case sdl.K_2:
						chip8.Key(0x2, true)
					case sdl.K_3:
						chip8.Key(0x3, true)
					case sdl.K_4:
						chip8.Key(0xC, true)
					case sdl.K_q:
						chip8.Key(0x4, true)
					case sdl.K_w:
						chip8.Key(0x5, true)
					case sdl.K_e:
						chip8.Key(0x6, true)
					case sdl.K_r:
						chip8.Key(0xD, true)
					case sdl.K_a:
						chip8.Key(0x7, true)
					case sdl.K_s:
						chip8.Key(0x8, true)
					case sdl.K_d:
						chip8.Key(0x9, true)
					case sdl.K_f:
						chip8.Key(0xE, true)
					case sdl.K_z:
						chip8.Key(0xA, true)
					case sdl.K_x:
						chip8.Key(0x0, true)
					case sdl.K_c:
						chip8.Key(0xB, true)
					case sdl.K_v:
						chip8.Key(0xF, true)
					}
				}

			}
		}
		sdl.Delay(1000 / 60)
	}

}
