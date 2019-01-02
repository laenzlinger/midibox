package main

import (
	"log"

	"github.com/laenzlinger/midibox/display"
	"github.com/laenzlinger/midibox/keyboard"
	"github.com/laenzlinger/midibox/midi"
	"github.com/laenzlinger/midibox/preset"

	"periph.io/x/periph/host"
)

func main() {

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	display := display.Open()
	defer display.Clear()

	md := midi.Open()
	defer md.Close()

	upDown := keyboard.OpenUpDown()
	joystick := keyboard.OpenJoystick()

	presets := preset.AllPresets()
	var current int
	var active bool

	display.DrawText("Select Preset", presets[current].Name())

	for i := 0; i < 50; i++ {
		select {
		case u := <-upDown:
			if active {
				if u == keyboard.Both {
					active = false
					display.DrawText("Select Preset", presets[current].Name())
					presets[current].Shutdown()
				} else {
					presets[current].OnUpDwon(u)
				}
			}
		case j := <-joystick:
			if active {
				presets[current].OnJoystick(j)
			} else {
				if j.Direction == keyboard.North {
					current--
					if current < 0 {
						current = len(presets) - 1
					}
					display.DrawText("Select Preset", presets[current].Name())
				} else if j.Direction == keyboard.South {
					current++
					if current >= len(presets) {
						current = 0
					}
					display.DrawText("Select Preset", presets[current].Name())
				} else if j.Fire && j.FireChanged {
					display.DrawText("Active Preset", presets[current].Name())
					presets[current].Init(md)
					active = true
				}
			}
		}
	}

	presets[current].Shutdown()

}
