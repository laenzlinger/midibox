package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	conn, err := net.Dial("udp", "127.0.0.1:5006")
	defer conn.Close()
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	noteOn(conn)

	time.Sleep(1 * time.Second)

	noteOff(conn)

}

func noteOn(conn net.Conn) {
	var noteOn = []byte{0xaa, 0x96, 0x3c, 0x7f}
	conn.Write(noteOn)
}

func noteOff(conn net.Conn) {
	var noteOff = []byte{0xaa, 0x86, 0x3c, 0x7f}
	conn.Write(noteOff)
}
