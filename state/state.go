package state

import (
	"errors"
	"strings"
)

type State struct {
	*Flash

	Y, X, B, A                             Button
	R, SR, ZR                              Button
	L, SL, ZL                              Button
	Minus, Plus, Home, Capture, ChargeGrip Button
	Down, Up, Right, Left                  Button

	LeftStick  Stick
	RightStick Stick

	Tick   uint64
	Mode   Mode
	Rumble Rumble

	Battery Battery
	Powered bool
}

func New(t DeviceKind) *State {
	s := &State{Flash: NewFlash()}

	s.SetSerial("")
	s.SetKind(t)

	s.SetHasColor(true)
	s.SetBodyColour(RandomColor())
	s.SetButtonColour(RandomColor())
	s.SetLeftGripColour(RandomColor())
	s.SetRightGripColour(RandomColor())

	axis := AxisCalibration{Min: 0, Center: 0x7ff, Max: 0xffe}
	stick := StickCalibration{X: axis, Y: axis}
	switch t {
	case SadLeft:
		s.SetLeftStickCalibration(&stick)
	case SadRight:
		s.SetRightStickCalibration(&stick)
	case Pro:
		s.SetLeftStickCalibration(&stick)
		s.SetRightStickCalibration(&stick)
	}

	s.Mode = FullMode
	s.Battery.Level = BatteryFull

	return s
}

func (s *State) Buttons() map[string]*Button {
	return map[string]*Button{
		"y":          &s.Y,
		"x":          &s.X,
		"b":          &s.B,
		"a":          &s.A,
		"r":          &s.R,
		"sr":         &s.SR,
		"zr":         &s.ZR,
		"l":          &s.L,
		"sl":         &s.SL,
		"zl":         &s.ZL,
		"minus":      &s.Minus,
		"plus":       &s.Plus,
		"home":       &s.Home,
		"capture":    &s.Capture,
		"down":       &s.Down,
		"up":         &s.Up,
		"right":      &s.Right,
		"left":       &s.Left,
		"chargegrip": &s.ChargeGrip,
		"leftstick":  &s.LeftStick.Button,
		"rightstick": &s.RightStick.Button,
	}
}

func (s *State) ButtonByName(name string) (*Button, error) {
	switch strings.ToLower(name) {
	case "-":
		name = "minus"
	case "+":
		name = "plus"
	case "v":
		name = "down"
	case "^":
		name = "up"
	case "<":
		name = "left"
	case ">":
		name = "right"
	case "ls", "lstick":
		name = "leftstick"
	case "rs", "rstick":
		name = "rightstick"
	}

	button, ok := s.Buttons()[name]
	if !ok {
		return nil, errors.New("unknown button")
	}
	return button, nil
}

func (s *State) StickByName(name string) (*Stick, error) {
	switch strings.ToLower(name) {
	case "l", "left":
		return &s.LeftStick, nil
	case "r", "right":
		return &s.RightStick, nil
	default:
		return nil, errors.New("unknown stick")
	}
}

type Mode byte

const (
	ActiveNFCIRMode0 Mode = 0x00
	ActiveNFCIRMode1 Mode = 0x01
	ActiveNFCIRMode2 Mode = 0x02
	ActiveNFCIRMode3 Mode = 0x03
	FullMode         Mode = 0x30
	NFCMode          Mode = 0x31
	BasicMode        Mode = 0x3f
)

type Rumble struct{}

type Battery struct {
	Level    BatteryLevel
	Charging bool
}

type BatteryLevel byte

const (
	BatteryFull     BatteryLevel = 0x80
	BatteryMedium   BatteryLevel = 0x60
	BatteryLow      BatteryLevel = 0x40
	BatteryCritical BatteryLevel = 0x20
	BatteryEmpty    BatteryLevel = 0x00
)
