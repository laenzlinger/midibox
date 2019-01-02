package mode

import (
	"github.com/laenzlinger/midibox/keyboard"
)

// Mode represents the current mode of operation of the midibox
type Mode interface {
	// Enter the mode
	Enter() *Mode
	// OnJoystick reacts on Joystick input
	OnJoystick(j keyboard.Joystick) *Mode
	// OnUpDown reacts on UpDown input
	OnUpDwon(j keyboard.UpDown) *Mode
	// Exit the mode
	Exit() *Mode
}

