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

// PinName returns the name of the GPIO pin
func (upDown UpDown) PinName() string {
	if upDown {
		return "GPIO6"
	}
	return "GPIO5"
}

// JoystickDirection represents the position of the joystick
//        North
//       --------
// West  | None | East
//       --------
//         South
type JoystickDirection uint8

const (
	// None not moved
	None JoystickDirection = iota
	// North direction movement
	North
	// East direction movement
	East
	// South direction movement
	South
	// West direction movement
	West
)

func (dir JoystickDirection) String() string {
	names := [...]string{
		"None",
		"North",
		"East",
		"South",
		"West",
	}
	if dir < None || dir > West {
		return "Unknown"
	}
	return names[dir]
}

// PinName returns the name of the GPIO pin
func (dir JoystickDirection) PinName() string {
	pinNames := [...]string{
		"GPIO4",
		"GPIO17",
		"GPIO23",
		"GPIO22",
		"GPIO27",
	}
	return pinNames[dir]
}

// Joystick is sent when a joystick action is detected.
type Joystick struct {
	Direction JoystickDirection
	Active    bool
}

func watchUpDown(upDown chan<- UpDown, keyboardTicker *time.Ticker) {

	up := gpioreg.ByName(Up.PinName())
	down := gpioreg.ByName(Down.PinName())

	if err := up.In(gpio.PullUp, gpio.NoEdge); err != nil {
		log.Fatal(err)
	}
	if err := down.In(gpio.PullUp, gpio.NoEdge); err != nil {
		log.Fatal(err)
	}

	var active = false
	var changed = time.Now()
	for tickTime := range keyboardTicker.C {
		var value UpDown
		var buttonPressed = false
		if up.Read() == gpio.Low {
			buttonPressed = true
			value = Up
		} else if down.Read() == gpio.Low {
			buttonPressed = true
			value = Down
		}
		if buttonPressed {
			if active {
				if time.Since(changed) > 500 * time.Millisecond {
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

func watchJoystick(joystick chan<- Joystick, value JoystickDirection) {
	p := gpioreg.ByName(value.PinName())

	if err := p.In(gpio.PullUp, gpio.FallingEdge); err != nil {
		log.Fatal(err)
	}

	for {
		p.WaitForEdge(-1)
		joystick <- Joystick{
			Direction: value,
			Active:    false,
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

	keyboardTicker := time.NewTicker(50 * time.Millisecond)
	upDown := make(chan UpDown)
	go watchUpDown(upDown, keyboardTicker)

	joystick := make(chan Joystick)
	// go watchJoystick(joystick, None)
	// go watchJoystick(joystick, North)
	// go watchJoystick(joystick, East)
	// go watchJoystick(joystick, South)
	// go watchJoystick(joystick, West)

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
