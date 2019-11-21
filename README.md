# Midibox

Playing with midi is a fun project to learn new technology. I created this small project to learn mainly about
the go language.

I am using a [Raspberry Pi Zero W](https://www.raspberrypi.org/products/raspberry-pi-zero-w/) and an
[Adafruit 128x64 OLED bonnet](https://www.adafruit.com/product/3531) to implement a small midi controller.

The reason for choosing this hardware is because it is so easily available and it offers a small screen and
some buttons to get started without any need to do any soldering.

The hardware is accessed by using the [periph](https://periph.io) hardware abstraction library.

This project is used to demonstrate the [go implementation of the RTP-MIDI protocol](https://github.com/laenzlinger/go-midi-rtp).

## Hardware

Coneection of the DigiTech Control SEVEN switches:

| Key    | Color   | Pin   | Chip   |
|--------|---------|-------|--------|
| 0      | türkis  |  12   | BCM 18 |
| 1      | violett |  18   | BCM 24 |
| 2      | violett |  22   | BCM 25 |
| 3      | türkis  |  32   | BCM 12 |
| 4      | white   |  36   | BCM 16 |
| UP     | blue    |  33   | BCM 13 |
| DOWN   | blue    |  37   | BCM 26 |

see [pinout](https://pinout.xyz/pinout/oled_bonnet)

Upper Row Ground: white:  (Pin 20)
Lower Row Ground: orange  (Pin 30)

## installation

``` bash
scp midibox.service pi@midibox
ssh pi@midibox sudo mv midibox.service /etc/systemd/system/midibox.service
ssh pi@midibox sudo systemctl enable midibox.service
make deploy
```
