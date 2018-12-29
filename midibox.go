package main

import (
	"fmt"
	"net"
	"log"

	"periph.io/x/periph/conn/gpio"
    "periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

func watchPin(level chan gpio.Level, pinName string) {
	p := gpioreg.ByName(pinName)
	fmt.Printf("%s: %s\n", p, p.Function())

	if err := p.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
	
	for {
		p.WaitForEdge(-1)
		level <- p.Read()
	}
}

func main() {

	conn, err := net.Dial("udp", "127.0.0.1:5006")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	down := make(chan gpio.Level)
	go watchPin(down, "GPIO5")
	up := make(chan gpio.Level)
	go watchPin(up, "GPIO6")
	
	for i := 0; i < 10; i++ {
        select {
		case level1 := <-up:
			note(conn, level1, 0x3c)
        case level2 := <-down:
			note(conn, level2, 0x3d)
        }
    }

}

func note(conn net.Conn, level gpio.Level, note byte) {
	if level == gpio.Low {
		noteOn(conn, note)
	} else {
		noteOff(conn, note)
	}
}

func noteOn(conn net.Conn, note byte) {
	fmt.Println("note on: ", note)
	var noteOn = []byte{0xaa, 0x96, note, 0x7f}
	conn.Write(noteOn)
}

func noteOff(conn net.Conn, note byte) {
	fmt.Println("note off:", note)
	var noteOff = []byte{0xaa, 0x86, note, 0x7f}
	conn.Write(noteOff)
}
