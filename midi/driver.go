package midi

import (
	"fmt"

	"github.com/grandcat/zeroconf"
	"github.com/laenzlinger/go-midi-rtp/session"
)

// Driver represents the midi connections
type Driver struct {
	zeroConfServer *zeroconf.Server
	midiRtpSession *session.MIDINetworkSession
}

// Open the midi drivers
func Open() Driver {
	port := 5005
	bonjourName := "midibox"
	server, err := zeroconf.Register(bonjourName, "_apple-midi._udp", "local.", port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		panic(err)
	}
	session := session.Start(bonjourName, uint16(port))
	return Driver{
		zeroConfServer: server,
		midiRtpSession: session,
	}
}

// Close the midi drivers
func (md Driver) Close() {
	md.midiRtpSession.End()
	md.zeroConfServer.Shutdown()
}

// Note sends a note midi message
func (md Driver) Note(on bool, note byte) {
	if on {
		md.NoteOn(note)
	} else {
		md.NoteOff(note)
	}
}

// NoteOn sends a note on midi message
func (md Driver) NoteOn(note byte) {
	fmt.Println("note on: ", note)
	noteOn := []byte{0x96, note, 0x7f}
	md.midiRtpSession.SendMIDIPayload(noteOn)
}

// NoteOff sends a Note off midi message
func (md Driver) NoteOff(note byte) {
	fmt.Println("note off:", note)
	noteOff := []byte{0x86, note, 0x7f}
	md.midiRtpSession.SendMIDIPayload(noteOff)
}

// SendRawMessage sends raw midi payload
func (md Driver) SendRawMessage(m []byte) {
	md.midiRtpSession.SendMIDIPayload(m)
}

