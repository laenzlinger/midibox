package preset

import (
	"fmt"

	"github.com/laenzlinger/midibox/keyboard"
	"github.com/laenzlinger/midibox/midi"
)

type transport struct {
	md midi.Driver
}

func (p transport) Name() string {
	return fmt.Sprintf("Transport")
}

func (p *transport) Init(md midi.Driver) {
	p.md = md
}

func (p *transport) OnFootKey(f keyboard.FootKey) {
	fmt.Println("Footkey: ", f)
}

func (p *transport) OnJoystick(j keyboard.Joystick) {
}

func (p *transport) OnUpDwon(u keyboard.UpDown) {
}

func (p *transport) Shutdown() {
}
