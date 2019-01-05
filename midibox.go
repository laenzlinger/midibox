package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/laenzlinger/midibox/mode"

	"github.com/laenzlinger/midibox/display"
	"github.com/laenzlinger/midibox/keyboard"

	"periph.io/x/periph/host"
)

func main() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	display := display.Open()
	defer display.Clear()

	// TODO: refacter the API to let the client define the channels
	upDown := keyboard.OpenUpDown()
	joystick := keyboard.OpenJoystick()

	mode := mode.Initial(display)
	defer mode.Exit()

	run := true
	for run {
		select {
		case u := <-upDown:
			mode = mode.OnUpDwon(u)
		case j := <-joystick:
			mode = mode.OnJoystick(j)
		case <-sig:
			run = false
		}
	}

}
