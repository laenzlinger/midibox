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
	var active bool


	display.DrawText("Select Preset", presets.Current().Name())

	for i := 0; i < 50; i++ {
		select {
		case u := <-upDown:
			if active {
				if u == keyboard.Both {
					active = false
					display.DrawText("Select Preset", presets.Current().Name())
					presets.Current().Shutdown()
				} else {
					presets.Current().OnUpDwon(u)
				}
			}
		case j := <-joystick:
			if active {
				presets.Current().OnJoystick(j)
			} else {
				if j.Direction == keyboard.North {
					presets.Previous()
					display.DrawText("Select Preset", presets.Current().Name())
				} else if j.Direction == keyboard.South {
					presets.Next()
					display.DrawText("Select Preset", presets.Current().Name())
				} else if j.Fire && j.FireChanged {
					display.DrawText("Active Preset", presets.Current().Name())
					presets.Current().Init(md)
					active = true
				}
			}
		}
	}

	presets.Current().Shutdown()

}
