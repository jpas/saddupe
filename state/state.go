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

func NewState(t DeviceKind) *State {
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

func (s *State) ButtonByName(name string) (*Button, error) {
	var b *Button

	switch strings.ToLower(name) {
	case "y":
		b = &s.Y
	case "x":
		b = &s.X
	case "b":
		b = &s.B
	case "a":
		b = &s.A
	case "r":
		b = &s.R
	case "sr":
		b = &s.SR
	case "zr":
		b = &s.ZR
	case "l":
		b = &s.L
	case "sl":
		b = &s.SL
	case "zl":
		b = &s.ZL
	case "-", "minus":
		b = &s.Minus
	case "+", "plus":
		b = &s.Plus
	case "home":
		b = &s.Home
	case "capture":
		b = &s.Capture
	case "v", "down":
		b = &s.Down
	case "^", "up":
		b = &s.Up
	case ">", "right":
		b = &s.Right
	case "<", "left":
		b = &s.Left
	case "chargegrip":
		b = &s.ChargeGrip
	case "ls", "lstick", "leftstick":
		b = &s.LeftStick.Button
	case "rs", "rstick", "rightstick":
		b = &s.RightStick.Button
	default:
		return nil, errors.New("unknown button")
	}
	return b, nil
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
