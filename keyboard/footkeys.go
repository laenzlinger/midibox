package keyboard

import (
	"time"
)

// FootKey is sent when one of the foot key buttons is pressed.
//  --------------------------------------
//  | DigiTech             CONTROL SEVEN |
//  |                                    |
//  |  display    [0]     [1]    [ up ]  |
//  |                                    |
//  |    [2]      [3]     [4]    [down]  |
//  --------------------------------------
//
// Implemented Behaviour:
// * Only one of the buttons can be active at a given time
// * Only pressing the button down edge is detected
// * when activating for more than one second => the event is triggered every 100ms again
type FootKey uint8

const (
	// Zero button was pressed
	Zero FootKey = iota
	// One button was pressed
	One
	// Two button was pressed
	Two
	// Three button was pressed
	Three
	// Four button was pressed
	Four
	// UP button was pressed
	UP
	// DOWN button was pressed
	DOWN
)

func (fk FootKey) String() string {
	names := [...]string{
		"zero",
		"one",
		"two",
		"three",
		"four",
		"up",
		"down",
	}
	if fk < Zero || fk > DOWN {
		return "Unknown"
	}
	return names[fk]
}

type footKeyButtons struct {
	zero  buttonPin
	one   buttonPin
	two   buttonPin
	three buttonPin
	four  buttonPin
	five  buttonPin
	up    buttonPin
	down  buttonPin
}

func watchFootKeys(footKey chan<- FootKey, b footKeyButtons) {
	keyboardTicker := time.NewTicker(100 * time.Millisecond)

	var previous = -1
	var lastChanged = time.Now()
	for tickTime := range keyboardTicker.C {
		var result FootKey
		var current = -1
		if b.zero.pressed() {
			current = 0
			result = Zero
		} else if b.one.pressed() {
			current = 1
			result = One
		} else if b.two.pressed() {
			current = 2
			result = Two
		} else if b.three.pressed() {
			current = 3
			result = Three
		} else if b.four.pressed() {
			current = 4
			result = Four
		} else if b.up.pressed() {
			current = 5
			result = UP
		} else if b.down.pressed() {
			current = 6
			result = DOWN
		}
		// debounce
		if current != previous && time.Since(lastChanged) > 300*time.Millisecond {
			previous = current
			lastChanged = tickTime
			if current >= 0 {
				footKey <- result
			}
		}
		// repeat
		if current >= 0 && time.Since(lastChanged) > 1000*time.Millisecond {
			footKey <- result
		}
	}
}

// OpenFootKeys open a channel that sends FootKey events
func OpenFootKeys() chan FootKey {

	buttons := footKeyButtons{
		zero:  registerPin("GPIO18"),
		one:   registerPin("GPIO24"),
		two:   registerPin("GPIO25"),
		three: registerPin("GPIO12"),
		four:  registerPin("GPIO16"),
		up:    registerPin("GPIO13"),
		down:  registerPin("GPIO26"),
	}

	footkey := make(chan FootKey)
	go watchFootKeys(footkey, buttons)
	return footkey
}
