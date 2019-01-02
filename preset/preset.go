package preset

import (
	"github.com/laenzlinger/midibox/keyboard"
	"github.com/laenzlinger/midibox/midi"
)

// Preset defines a certain behavior of the midibox.
type Preset interface {
	Name() string
	Init(md midi.Driver)
	OnJoystick(j keyboard.Joystick)
	OnUpDwon(j keyboard.UpDown)
	Shutdown()
}

// Presets is the coolection of presets defined in midibox.
type Presets struct {
	presets []Preset
	current int
}

// AllPresets that are registered.
func AllPresets() Presets {
	all := []Preset{
		&chromaticScale{},
		&majorScale{},
	}
	return Presets{
		presets: all,
		current: 0,
	}
}

// Current returns the currently selected preset.
func (p *Presets) Current() Preset {
	return p.presets[p.current]
}

// Next preset must be selected.
func (p *Presets) Next() {
	p.current++
	if p.current >= len(p.presets) {
		p.current = 0
	}
}

// Previous preset must be selected.
func (p *Presets) Previous() {
	p.current--
	if p.current < 0 {
		p.current = len(p.presets) - 1
	}
}
