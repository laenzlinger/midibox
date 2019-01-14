# Midibox

Playing with midi is a fun project to learn new technology. I created this small project to learn mainly about
the go language.

I am using a [Raspberry Pi Zero W](https://www.raspberrypi.org/products/raspberry-pi-zero-w/) and an
[Adafruit 128x64 OLED bonnet](https://www.adafruit.com/product/3531) to implement a small midi controller.

The reason for choosing this hardware is because it is so easily available and it offers a small screen and
some buttons to get started without any need to do any soldering.

The hardware is accessed by using the [periph](https://periph.io) hardware abstraction library.

This project is used to demonstrate the [go implementation of the RTP-MIDI protocol](https://github.com/laenzlinger/go-midi-rtp).