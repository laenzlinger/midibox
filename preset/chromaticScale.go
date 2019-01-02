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
	return fmt.Sprintf("Chrom. Scale")
}

func (p *chromaticScale) Init(md midi.Driver) {
	p.current = 0
	p.base = 0x3c
	p.md = md
}

func (p *chromaticScale) OnJoystick(j keyboard.Joystick) {
	if j.DirectionChanged {
		p.Shutdown()
		if j.Direction != keyboard.Center {
			p.current = p.base + byte(j.Direction)
			p.md.NoteOn(p.current)
		}
	}
}

func (p *chromaticScale) OnUpDwon(u keyboard.UpDown) {
	if (u == keyboard.Up && p.base <= 72) {
		p.base++
	} else if (u == keyboard.Down && p.base >= 60) {
		p.base--
	}
}

func (p *chromaticScale) Shutdown() {
	if (p.current > 0) {
		p.md.NoteOff(p.current)
		p.current = 0
	}
}
