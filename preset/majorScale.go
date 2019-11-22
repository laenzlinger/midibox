package preset

import (
	"fmt"

	"github.com/laenzlinger/midibox/display"
	"github.com/laenzlinger/midibox/keyboard"
	"github.com/laenzlinger/midibox/midi"
)

type majorScale struct {
	base     byte
	current  byte
	md       midi.Driver
	interval map[keyboard.JoystickDirection]byte
}

func (p majorScale) Name() string {
	return fmt.Sprintf("Major Scale")
}

func (p *majorScale) Init(md midi.Driver, display display.Display) {
	p.current = 0
	p.base = 0x3c
	p.md = md
	p.interval = map[keyboard.JoystickDirection]byte{
		keyboard.North:     0,
		keyboard.NorthEast: 2,
		keyboard.East:      4,
		keyboard.SouthEast: 5,
		keyboard.South:     7,
		keyboard.SouthWest: 9,
		keyboard.West:      11,
		keyboard.NorthWest: 12,
	}
}

func (p *majorScale) OnJoystick(j keyboard.Joystick) {
	if j.DirectionChanged {
		p.Shutdown()
		if j.Direction != keyboard.Center {
			p.current = p.base + p.interval[j.Direction]
			p.md.NoteOn(p.current)
		}
	}
}

func (p *majorScale) OnFootKey(f keyboard.FootKey) {
}

func (p *majorScale) OnUpDwon(u keyboard.UpDown) {
	if u == keyboard.Up && p.base <= 72 {
		p.base++
	} else if u == keyboard.Down && p.base >= 60 {
		p.base--
	}
}

func (p *majorScale) Shutdown() {
	if p.current > 0 {
		p.md.NoteOff(p.current)
		p.current = 0
	}
}
