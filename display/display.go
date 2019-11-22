package display

import (
	"image"
	"image/color"
	"log"

	"github.com/pbnjay/pixfont"
	"github.com/disintegration/imaging"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
)

// Display wrappes the underlying device and offers operations to interact with the display
type Display struct {
	dev *ssd1306.Dev
}

const fontHeight = 10

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
		font.DrawString(img, 0, i*(fontHeight+1), line, color.White)
	}
	d.drawImage(img)
}

// DrawLargeText will draw double sized text
func (d Display) DrawLargeText(lines ...string) {
	var img = image.NewRGBA(image.Rect(0, 0, 64, 32))
	font := pixfont.DefaultFont
	font.SetVariableWidth(true)

	ystart := int((32 - fontHeight*len(lines)) / (len(lines) +1))
	for i, line := range lines {
		font.DrawString(img, 0, ystart + i*(fontHeight+1), line, color.White)
	}
	dst := imaging.Resize(img, 128, 64, imaging.Lanczos)
	d.dev.Draw(dst.Bounds(), dst, image.Point{})
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
