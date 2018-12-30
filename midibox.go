package main

import (
	"fmt"
	"net"
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
    "periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

type UpDown bool

const (
	Up UpDown = true
	Down UpDown = false
)

type JoystickDirection uint8

const (
	None  JoystickDirection = iota 
	North
	East 
	South
	West
)

type Joystick struct {
	Direction JoystickDirection
	Active bool
}

func watchUpDown(upDown chan <- UpDown, pinName string, value UpDown) {
	p := gpioreg.ByName(pinName)

	if err := p.In(gpio.PullUp, gpio.FallingEdge); err != nil {
		log.Fatal(err)
	}
	
	for {
		p.WaitForEdge(-1)
		upDown <- value
		p.In(gpio.PullNoChange, gpio.NoEdge)
		time.Sleep(250 * time.Millisecond)
		p.In(gpio.PullNoChange, gpio.FallingEdge)
	}
}

func watchJoystick(joystick chan <- Joystick, pinName string, value JoystickDirection) {
	p := gpioreg.ByName(pinName)

	if err := p.In(gpio.PullUp, gpio.FallingEdge); err != nil {
		log.Fatal(err)
	}
	
	for {
		p.WaitForEdge(-1)
		joystick <- Joystick{
			Direction: value,
			Active: false,
		}
		p.In(gpio.PullNoChange, gpio.NoEdge)
		time.Sleep(250 * time.Millisecond)
		p.In(gpio.PullNoChange, gpio.FallingEdge)
	}
}


func main() {

	conn, err := net.Dial("udp", "127.0.0.1:5006")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	upDown := make(chan UpDown)
	go watchUpDown(upDown, "GPIO5", Down)
	go watchUpDown(upDown, "GPIO6", Up)
	

	joystick := make(chan Joystick)
	go watchJoystick(joystick, "GPIO4", None)
	go watchJoystick(joystick, "GPIO17", North)
	go watchJoystick(joystick, "GPIO23", East)
	go watchJoystick(joystick, "GPIO22", South)
	go watchJoystick(joystick, "GPIO27", West)


	for i := 0; i < 10; i++ {
        select {
		case upDown := <-upDown:
			fmt.Println("upDown:", upDown)
		case joystick := <-joystick:
			fmt.Println("joystick:", joystick)
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
