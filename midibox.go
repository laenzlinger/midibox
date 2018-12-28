package main

import (
	"fmt"
	"net"
	"time"

	"github.com/warthog618/gpio"
)

func main() {

	conn, err := net.Dial("udp", "127.0.0.1:5006")
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	handler := func(pin *gpio.Pin) {
		if pin.Read() {
			noteOff(conn)
		} else {
			noteOn(conn)
		}
	}

	gpioErr := gpio.Open()
	defer gpio.Close()
	if gpioErr != nil {
		panic(gpioErr)
	}

	pin := gpio.NewPin(gpio.GPIO5)
	pin.Input()
	pin.PullUp()
	watchErr := pin.Watch(gpio.EdgeBoth, handler)
	defer pin.Unwatch()
	if watchErr != nil {
		panic(watchErr)
	}

	fmt.Println("setup done")
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
	}

}

func noteOn(conn net.Conn) {
	fmt.Println("note on")
	var noteOn = []byte{0xaa, 0x96, 0x3c, 0x7f}
	conn.Write(noteOn)
}

func noteOff(conn net.Conn) {
	fmt.Println("note off")
	var noteOff = []byte{0xaa, 0x86, 0x3c, 0x7f}
	conn.Write(noteOff)
}
