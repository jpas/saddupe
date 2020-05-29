package state

import (
	"errors"
	"strings"
)

type State struct {
	Tick uint64
	Mode Mode

	Battery Battery
	Powered bool

	HasGrip bool

	Y, X, B, A                             Button
	R, SR, ZR                              Button
	L, SL, ZL                              Button
	Minus, Plus, Home, Capture, ChargeGrip Button
	Down, Up, Right, Left                  Button

	LeftStick  Stick
	RightStick Stick

	Flash *Flash

	Rumble Rumble
}

func NewState() *State {
	return &State{Flash: NewFlash()}
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
	case "zr":
		b = &s.ZR
	case "l":
		b = &s.L
	case "zl":
		b = &s.ZL
	case "minus":
		b = &s.Minus
	case "plus":
		b = &s.Plus
	case "home":
		b = &s.Home
	case "capture":
		b = &s.Capture
	case "down":
		b = &s.Down
	case "up":
		b = &s.Up
	case "right":
		b = &s.Right
	case "left":
		b = &s.Left
	case "leftstick":
		b = &s.LeftStick.Button
	case "rightstick":
		b = &s.RightStick.Button
	default:
		return nil, errors.New("unknown button")
	}
	return b, nil
}

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
