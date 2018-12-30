package main

import (
	"fmt"
	"log"
	"net"
	"image"
	"image/color"

	"github.com/laenzlinger/midibox/keyboard"
	"github.com/pbnjay/pixfont"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/host"

)

func main() {

	conn, err := net.Dial("udp", "127.0.0.1:5006")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	// Open a handle to a ssd1306 connected on the I²C bus:
	dev, err := ssd1306.NewI2C(bus, &ssd1306.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer drawText(dev, "");

	upDown := keyboard.OpenUpDown()
	joystick := keyboard.OpenJoystick()

	for i := 0; i < 20; i++ {
		select {
		case u := <-upDown:
			drawText(dev, fmt.Sprintf("%v", u));
		case j := <-joystick:
			drawText(dev, j.Direction.String(), fmt.Sprintf("%v", j.Fire));
		}
	}

}

func drawText(dev *ssd1306.Dev, lines ...string ) {
	var img = image.NewRGBA(image.Rect(0, 0, 128, 64))
	font := pixfont.DefaultFont
	font.SetVariableWidth(true)
	for i, line := range lines {
		font.DrawString(img, 0, i*11, line, color.White)
	}
	dev.Draw(img.Bounds(), img, image.Point{})
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
