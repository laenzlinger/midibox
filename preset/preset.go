package preset

import (
	"github.com/laenzlinger/midibox/midi"	
	"github.com/laenzlinger/midibox/keyboard"
)

// Preset defines a certain behavior of the midibox
type Preset interface {
	Name() string
	Init(md midi.Driver)
    OnJoystick(j keyboard.Joystick)
	OnUpDwon(j keyboard.UpDown)
	Shutdown()
}
// AllPresets that are registered
func AllPresets() []Preset {
	return []Preset{
		&chromaticScale{},
		&majorScale{},
	}
}