package main

import (
	"github.com/laenzlinger/midibox/mode"
	"log"

	"github.com/laenzlinger/midibox/display"
	"github.com/laenzlinger/midibox/keyboard"

	"periph.io/x/periph/host"
)

func main() {

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	display := display.Open()
	defer display.Clear()

	upDown := keyboard.OpenUpDown()
	joystick := keyboard.OpenJoystick()

	
    active := mode.InitialMode(display)

	for i := 0; i < 50; i++ {
		select {
		case u := <-upDown:
			active = active.OnUpDwon(u)
		case j := <-joystick:
			active = active.OnJoystick(j)
		}
	}

	active.Exit()

}
