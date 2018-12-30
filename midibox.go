package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

// UpDown is sent when one of the up/down buttons is pressed.
//  -------------------
//           |        |
//   OLED    |     x  | up button
//   display |  x     | down button
//           |        |
//  -------------------
type UpDown bool

const (
	// Up button
	Up UpDown = true
	// Down button
	Down UpDown = false
)

func (upDown UpDown) String() string {
	if upDown {
		return "up"
	}
	return "down"
}

// JoystickDirection represents the position of the joystick
//         North
//        --------
// West  | Center | East
//        --------
//         South
type JoystickDirection uint8

const (
	// Center direction
	Center JoystickDirection = iota
	// North direction
	North
	// East direction
	East
	// South direction
	South
	// West direction
	West
)

func (dir JoystickDirection) String() string {
	names := [...]string{
		"Center",
		"North",
		"East",
		"South",
		"West",
	}
	if dir < Center || dir > West {
		return "Unknown"
	}
	return names[dir]
}

// Joystick is sent when a joystick action is detected.
type Joystick struct {
	direction JoystickDirection
	active    bool
}

func (j Joystick) String() string {
	return fmt.Sprintf("%v %v", j.direction, j.active)
}

type buttonPin struct {
	pin gpio.PinIn
}

func registerPin(name string) buttonPin {
	pin := gpioreg.ByName(name)
	if err := pin.In(gpio.PullUp, gpio.NoEdge); err != nil {
		log.Fatal(err)
	}
	return buttonPin{pin: pin}
}

func (pin buttonPin) pressed() bool {
	return pin.pin.Read() == gpio.Low
}

type upDownButtons struct {
	up   buttonPin
	down buttonPin
}

func watchUpDown(upDown chan<- UpDown, b upDownButtons) {
	keyboardTicker := time.NewTicker(100 * time.Millisecond)

	var active = false
	var changed = time.Now()
	for tickTime := range keyboardTicker.C {
		var value UpDown
		var buttonPressed = false
		if b.up.pressed() {
			buttonPressed = true
			value = Up
		} else if b.down.pressed() {
			buttonPressed = true
			value = Down
		}
		if buttonPressed {
			if active {
				if time.Since(changed) > 500*time.Millisecond {
					upDown <- value
				}
			} else {
				active = true
				changed = tickTime
				upDown <- value
			}
		} else {
			active = false
			changed = time.Now()
		}

	}
}

// OpenUpDown open a channel that sends UpDown events
func OpenUpDown() chan UpDown {

	buttons := upDownButtons{
		up:   registerPin("GPIO6"),
		down: registerPin("GPIO5"),
	}

	upDown := make(chan UpDown)
	go watchUpDown(upDown, buttons)
	return upDown
}

type joystickButtons struct {
	north  buttonPin
	east   buttonPin
	south  buttonPin
	west   buttonPin
	center buttonPin
}

func watchJoystick(joystick chan<- Joystick, b joystickButtons) {
	keyboardTicker := time.NewTicker(100 * time.Millisecond)

	var active = false
	var changed = time.Now()
	for tickTime := range keyboardTicker.C {
		var value Joystick
		var buttonPressed = false

		if b.north.pressed() {
			value.direction = North
			buttonPressed = true
		} else if b.east.pressed() {
			value.direction = East
			buttonPressed = true
		} else if b.south.pressed() {
			value.direction = South
			buttonPressed = true
		} else if b.west.pressed() {
			value.direction = West
			buttonPressed = true
		}
		if b.center.pressed() {
			value.active = true
			buttonPressed = true
		}

		if buttonPressed {
			if active {
				if time.Since(changed) > 500*time.Millisecond {
					joystick <- value
				}
			} else {
				active = true
				changed = tickTime
				joystick <- value
			}
		} else {
			active = false
			changed = time.Now()
		}

	}
}

// OpenJoystick open a channel that sends Joystick events
func OpenJoystick() chan Joystick {

	buttons := joystickButtons{
		center:  registerPin("GPIO4"),
		north:   registerPin("GPIO17"),
		east:  registerPin("GPIO23"),
		south:   registerPin("GPIO22"),
		west: registerPin("GPIO27"),
	}

	joystick := make(chan Joystick)
	go watchJoystick(joystick, buttons)
	return joystick
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

	for i := 0; i < 10; i++ {
		select {
		case upDown := <-OpenUpDown():
			fmt.Println("upDown:", upDown)
		case joystick := <-OpenJoystick():
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
