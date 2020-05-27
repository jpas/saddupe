package state

import "time"

type State struct {
	Tick uint64
	Mode Mode

	Battery BatteryState
	Powered bool

	HasGrip bool

	Y, X, B, A                             Button
	R, SR, ZR                              Button
	L, SL, ZL                              Button
	Minus, Plus, Home, Capture, ChargeGrip Button
	Down, Up, Right, Left                  Button

	LeftStick  Stick
	RightStick Stick

	Flash Flash

	Rumble RumbleState
}

func NewState() *State {
	return &State{}
}

type Gamepad byte

type BatteryState struct {
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

type Button struct {
	Pressed bool
	start   time.Time
	millis  uint64
}

func (b *Button) Press() {
	if !b.Pressed {
		b.start = time.Now()
	}
	b.Pressed = true
}

func (b *Button) Release() {
	if b.Pressed {
		b.millis += uint64(time.Since(b.start).Milliseconds())
	}
	b.Pressed = false
}

func (b *Button) Milliseconds() uint64 {
	return b.millis
}

type Stick struct {
	Button
	X float64
	Y float64
}

type RumbleState struct{}

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
