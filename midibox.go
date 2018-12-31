package main

import (
	"fmt"
	"log"

	"github.com/laenzlinger/midibox/display"
	"github.com/laenzlinger/midibox/keyboard"
	"github.com/laenzlinger/midibox/midi"

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

	var note byte
	for i := 0; i < 20; i++ {
		select {
		case u := <-upDown:
			display.DrawText(fmt.Sprintf("%v", u))
		case j := <-joystick:
			if j.DirectionChanged {
				md.NoteOff(note)
				if j.Direction != keyboard.Center {
					note = 0x3c + byte(j.Direction)
					md.NoteOn(note)
				}
			}
		}
	}

}
