package midi

import (
	"fmt"
	"log"
	"net"
)

// Driver represents the midi connections
type Driver struct {
	conn net.Conn
}

// Open the midi drivers
func Open() Driver {
	conn, err := net.Dial("udp", "127.0.0.1:5006")
	if err != nil {
		log.Fatal(err)
	}
	return Driver{ conn: conn }
}

// Close the midi drivers
func (md Driver) Close()  {
	md.conn.Close()
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
	var noteOn = []byte{0xaa, 0x96, note, 0x7f}
	md.conn.Write(noteOn)
}

// NoteOff sends a Note off midi message
func (md Driver) NoteOff(note byte) {
	fmt.Println("note off:", note)
	var noteOff = []byte{0xaa, 0x86, note, 0x7f}
	md.conn.Write(noteOff)
}
