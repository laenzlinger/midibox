package preset

import (
	"bytes"
	"fmt"
	"time"

	"github.com/laenzlinger/midibox/display"
	"github.com/laenzlinger/midibox/keyboard"
	"github.com/laenzlinger/midibox/midi"
)

type mmcCommandID2 byte

const (
	stop         mmcCommandID2 = 0x01
	play         mmcCommandID2 = 0x02
	deferredPlay mmcCommandID2 = 0x03 // Play (play after no longer busy
	fastForward  mmcCommandID2 = 0x04
	rewind       mmcCommandID2 = 0x05
	recordStrobe mmcCommandID2 = 0x06 // AKA [[Punch in/out|Punch In]]
	recordExit   mmcCommandID2 = 0x07 // AKA [[Punch out (music)|Punch out]]
	recordPause  mmcCommandID2 = 0x08
	pause        mmcCommandID2 = 0x09 // pause playback
	eject        mmcCommandID2 = 0x0A // disengage media container from MMC device
	chase        mmcCommandID2 = 0x0B
	mMCReset     mmcCommandID2 = 0x0D // to default/startup state
	write        mmcCommandID2 = 0x40 // AKA Record Ready, AKA Arm Tracks parameters: <length1> 4F <length2> <track-bitmap-bytes>
	gotoCmnd     mmcCommandID2 = 0x44 // AKA Locate parameters: <length>=06 01 <hours> <minutes> <seconds> <frames> <subframes>
	shuttle      mmcCommandID2 = 0x47 // parameters: <length>=03 <sh> <sm> <sl> (MIDI Standard Speed codes)
)

// Send MMC transport messages
// see https://en.wikipedia.org/wiki/MIDI_Machine_Control
type transport struct {
	md      midi.Driver
	display display.Display
}

func (p transport) Name() string {
	return fmt.Sprintf("MMC Transport")
}

func (p *transport) Init(md midi.Driver, display display.Display) {
	p.md = md
	p.display = display
}

func (p *transport) OnFootKey(f keyboard.FootKey) {
	switch f {
	case keyboard.Two:
		go p.popupMessage("Stop")
		p.sendMMCMessage(stop)
	case keyboard.Three:
		go p.popupMessage("Play")
		p.sendMMCMessage(play)
	case keyboard.Four: 
		go p.popupMessage("Record")
		p.sendMMCMessage(recordStrobe)
	case keyboard.UP:
		go p.popupMessage("Rewind")
		p.sendMMCMessage(rewind)
	case keyboard.DOWN:
		go p.popupMessage("Forward")
		p.sendMMCMessage(fastForward)
	}
}

func (p *transport) OnJoystick(j keyboard.Joystick) {
}

func (p *transport) OnUpDwon(u keyboard.UpDown) {
}

func (p *transport) Shutdown() {
}

func (p *transport) stop() {
}

func (p *transport) popupMessage(m string) {
	p.display.DrawLargeText(m)
	time.Sleep(500 * time.Millisecond)
	p.display.Clear()
	p.display.DrawText("", p.Name())
}

func (p *transport) sendMMCMessage(command mmcCommandID2) {
	b := bytes.NewBuffer([]byte{})
	b.WriteByte(0xF0) // SYSEX start
	b.WriteByte(0x7F) // MMC
	b.WriteByte(0x7F) // Device Id (7F => all)
	b.WriteByte(0x06) // MMC command (0x07 is resposne)
	b.WriteByte(byte(command))
	b.WriteByte(0xF7) // SYSEX end
	p.md.SendRawMessage(b.Bytes())
}
