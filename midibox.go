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
	footkey := keyboard.OpenFootKeys()

	m := mode.Initial(display)
	defer mode.Shutdown()

	run := true
	for run {
		select {
		case u := <-upDown:
			m = m.OnUpDwon(u)
		case j := <-joystick:
			m = m.OnJoystick(j)
		case f := <- footkey:
			m = m.OnFootKey(f)
		case <-sig:
			m.Exit()
			run = false
		}
	}
}
