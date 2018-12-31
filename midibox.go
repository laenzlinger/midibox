package main

import (
	"time"
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

	keyboard := keyboard.OpenKeyboard()

	for i := 0; i < 20; i++ {
		select {
		case u := <-keyboard.UpDown:
			display.DrawText(fmt.Sprintf("%v", u))
		case j := <-keyboard.Joystick:
			go func() {
				note := 0x3c + byte(j.Direction)
				md.NoteOn(note)
				time.Sleep(300*time.Millisecond)
				md.NoteOff(note)
			}()
		}
	}

}
