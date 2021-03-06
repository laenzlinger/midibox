package mode

import (
	"github.com/laenzlinger/midibox/display"
	"github.com/laenzlinger/midibox/keyboard"
	"github.com/laenzlinger/midibox/midi"
	"github.com/laenzlinger/midibox/preset"
)

var midiDriver midi.Driver

// Mode represents the current mode of operation of the midibox
type Mode interface {
	// OnJoystick reacts on Joystick input
	OnJoystick(j keyboard.Joystick) Mode
	// OnUpDown reacts on UpDown input
	OnUpDwon(u keyboard.UpDown) Mode
	// OnFootKey reacts on FootKey input
	OnFootKey(u keyboard.FootKey) Mode
	// Exit the mode
	Exit() Mode
}

// Initial returns the inital mode.
func Initial(display display.Display) Mode {
	midiDriver = midi.Open()
	m := selectPreset{}
	return m.Enter(display).SelectCurrentMode()
}

// Shutdown the current mode.
func Shutdown()  {
	midiDriver.Close()
}

// SelectPreset is the mode which allows to select the preset.
type selectPreset struct {
	display display.Display
	presets preset.Presets
}

func (m *selectPreset) Enter(display display.Display) *selectPreset {
	m.presets = preset.AllPresets()
	m.display = display
	m.display.DrawText("Select Preset", m.presets.Current().Name())
	return m
}

func (m *selectPreset) OnJoystick(j keyboard.Joystick) Mode {
	if j.Fire && j.FireChanged {
		return m.SelectCurrentMode()
	}
	if j.Direction == keyboard.North {
		m.presets.Previous()
	} else if j.Direction == keyboard.South {
		m.presets.Next()
	}  
	m.display.DrawText("Select Preset", m.presets.Current().Name())
	return m
}

func (m *selectPreset) OnUpDwon(u keyboard.UpDown) Mode {
	return m
}

func (m *selectPreset) OnFootKey(u keyboard.FootKey) Mode {
	return m
}

func (m *selectPreset) Exit() Mode {
	return m
}

func (m *selectPreset) SelectCurrentMode() Mode {
	next := playMode{}
	return next.Enter(m.display, m.presets.Current())
}

// playMode is the mode where the selected preset is played
type playMode struct {
	display display.Display
	preset  preset.Preset
}

func (m *playMode) Enter(display display.Display, preset preset.Preset) Mode {
	m.preset = preset
	m.display = display
	m.display.DrawText("Active Preset", m.preset.Name())
	m.preset.Init(midiDriver, display)
	return m
}

func (m *playMode) OnJoystick(j keyboard.Joystick) Mode {
	if j.Fire && j.FireChanged {
		return m.Exit()
	}
	m.preset.OnJoystick(j)
	return m
}

func (m *playMode) OnUpDwon(u keyboard.UpDown) Mode {
	m.preset.OnUpDwon(u)
	return m
}

func (m *playMode) OnFootKey(f keyboard.FootKey) Mode {
	m.preset.OnFootKey(f)
	return m
}

func (m *playMode) Exit() Mode {
	m.preset.Shutdown()
	next := selectPreset{}
	return next.Enter(m.display)
}
