package main

import (
	"fmt"
	"log"
	"net"

	"github.com/laenzlinger/midibox/display"
	"github.com/laenzlinger/midibox/keyboard"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
)

func main() {

	conn, err := net.Dial("udp", "127.0.0.1:5006")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	display := display.Open()
	defer display.Clear()

	upDown := keyboard.OpenUpDown()
	joystick := keyboard.OpenJoystick()

	for i := 0; i < 20; i++ {
		select {
		case u := <-upDown:
			display.DrawText(fmt.Sprintf("%v", u))
		case j := <-joystick:
			display.DrawText(j.Direction.String(), fmt.Sprintf("%v", j.Fire))
		}
	}

}

func note(conn net.Conn, level gpio.Level, note byte) {
	if level == gpio.Low {
		noteOn(conn, note)
	} else {
		noteOff(conn, note)
	}
}

func noteOn(conn net.Conn, note byte) {
	fmt.Println("note on: ", note)
	var noteOn = []byte{0xaa, 0x96, note, 0x7f}
	conn.Write(noteOn)
}

func noteOff(conn net.Conn, note byte) {
	fmt.Println("note off:", note)
	var noteOff = []byte{0xaa, 0x86, note, 0x7f}
	conn.Write(noteOff)
}
