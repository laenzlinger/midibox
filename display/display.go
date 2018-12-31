package display

import (
	"image"
	"image/color"
	"log"

	"github.com/pbnjay/pixfont"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
)

// Display wrappes the underlying device and offers operations to interact with the display
type Display struct {
	dev *ssd1306.Dev
}

// Open a new hander for the display device
func Open() Display {
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
	return Display{dev: dev}
}

// DrawText will show the given (max. 5) lines of text on the display
func (d Display) DrawText(lines ...string) {
	var img = newEmptyImage()
	font := pixfont.DefaultFont
	font.SetVariableWidth(true)
	for i, line := range lines {
		font.DrawString(img, 0, i*11, line, color.White)
	}
	d.drawImage(img)
}

// Clear the display
func (d Display) Clear() {
	d.drawImage(newEmptyImage())
}

func newEmptyImage() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, 128, 64))
}

func (d Display) drawImage(img *image.RGBA) {
	d.dev.Draw(img.Bounds(), img, image.Point{})
}
