package preset

import ( 
	"fmt"
	"github.com/laenzlinger/midibox/keyboard"
     "github.com/laenzlinger/midibox/midi"
)

type chromaticScale struct {
	base byte
	current byte
	md midi.Driver
}

func (p chromaticScale) Name() string {
	return fmt.Sprintf("Major Scale: %d", p.base)
}

func (p *chromaticScale) Init(md midi.Driver) {
	p.current = 0
	p.md = md
}

func (p *chromaticScale) OnJoystick(j keyboard.Joystick) {
	if j.DirectionChanged {
		if (p.current > 0) {
			p.md.NoteOff(p.current)
		}
		if j.Direction != keyboard.Center {
			p.current = p.base + byte(j.Direction)
			p.md.NoteOn(p.current)
		}
	}
}

func (p *chromaticScale) OnUpDwon(u keyboard.UpDown) {

}

func (p *chromaticScale) Shutdown() {

}
