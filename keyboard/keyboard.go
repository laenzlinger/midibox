package keyboard

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

// UpDown is sent when one of the up/down buttons is pressed.
//  -------------------
//           |        |
//   OLED    |     x  | up button
//   display |  x     | down button
//           |        |
//  -------------------
type UpDown uint8

const (
	// Down button
	Down UpDown = iota
	// Up button
	Up
	// Both buttons
	Both
)

func (upDown UpDown) String() string {
	names := [...]string{
		"down",
		"up",
		"both",
	}
	if upDown < Down || upDown > Both {
		return "Unknown"
	}
	return names[upDown]
}

// JoystickDirection represents the position of the joystick
//
//    NorthWest   North    NorthEast
//               --------
//  West        | Center |       East
//               --------
//    SouthWest   South   SouthEast
//
type JoystickDirection uint8

const (
	// Center direction
	Center JoystickDirection = iota
	// North direction
	North
	// NorthEast direction
	NorthEast
	// East direction
	East
	// SouthEast direction
	SouthEast
	// South direction
	South
	// SouthWest direction
	SouthWest
	// West direction
	West
	// NorthWest direction
	NorthWest
)

func (dir JoystickDirection) String() string {
	names := [...]string{
		"center",
		"north",
		"north-east",
		"east",
		"south-east",
		"south",
		"south-west",
		"west",
		"north-west",
	}
	if dir < Center || dir > NorthWest {
		return "Unknown"
	}
	return names[dir]
}

// Joystick is sent when a joystick action is detected.
type Joystick struct {
	Direction        JoystickDirection
	DirectionChanged bool
	Fire             bool
	FireChanged      bool
}

func (j Joystick) String() string {
	return fmt.Sprintf("%v fire:%v", j.Direction, j.Fire)
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

	var previous = 0
	var lastChanged = time.Now()
	for tickTime := range keyboardTicker.C {
		var result UpDown
		var current = 0
		if b.up.pressed() && b.down.pressed() {
			current = 1
			result = Both
		} else if b.down.pressed() {
			current = 2
			result = Down
		} else if b.up.pressed() {
			current = 3
			result = Up
		}
		// debounce
		if current != previous && time.Since(lastChanged) > 300*time.Millisecond {
			previous = current
			lastChanged = tickTime
			if current != 0 {
				upDown <- result
			}
		}
		// repeat
		if current != 0 && time.Since(lastChanged) > 1000*time.Millisecond {
			upDown <- result
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
	north buttonPin
	east  buttonPin
	south buttonPin
	west  buttonPin
	fire  buttonPin
}

var inactiveJoystick = Joystick{Direction: Center, Fire: false}

func watchJoystick(joystick chan<- Joystick, b joystickButtons) {

	keyboardTicker := time.NewTicker(50 * time.Millisecond)

	var previous = inactiveJoystick
	var lastChanged = time.Now()
	for tickTime := range keyboardTicker.C {
		var current = Joystick{
			Direction: defineJoystickDirection(b),
			Fire:      b.fire.pressed(),
		}
		directionChanged := current.Direction != previous.Direction
		fireChanged := current.Fire != previous.Fire
		debounced := time.Since(lastChanged) > 200*time.Millisecond
		if (directionChanged || fireChanged) && debounced {
			current.DirectionChanged = directionChanged
			current.FireChanged = fireChanged
			previous = current
			lastChanged = tickTime
			joystick <- current
		}
	}
}

func defineJoystickDirection(b joystickButtons) JoystickDirection {
	if b.north.pressed() && !b.east.pressed() && !b.west.pressed() {
		return North
	} else if b.east.pressed() && b.north.pressed() {
		return NorthEast
	} else if b.east.pressed() && !b.north.pressed() && !b.south.pressed() {
		return East
	} else if b.east.pressed() && b.south.pressed() {
		return SouthEast
	} else if b.south.pressed() && !b.west.pressed() && !b.east.pressed() {
		return South
	} else if b.south.pressed() && b.west.pressed() {
		return SouthWest
	} else if b.west.pressed() && !b.north.pressed() && !b.south.pressed() {
		return West
	} else if b.north.pressed() && b.west.pressed() {
		return NorthWest
	} else {
		return Center
	}
}

// OpenJoystick open a channel that sends Joystick events
func OpenJoystick() chan Joystick {

	buttons := joystickButtons{
		fire:  registerPin("GPIO4"),
		north: registerPin("GPIO17"),
		east:  registerPin("GPIO23"),
		south: registerPin("GPIO22"),
		west:  registerPin("GPIO27"),
	}

	joystick := make(chan Joystick)
	go watchJoystick(joystick, buttons)
	return joystick
}
