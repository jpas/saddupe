package state

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
}

func (b *Button) Press() {
	b.Pressed = true
}

func (b *Button) Release() {
	b.Pressed = false
}

type Stick struct {
	*Button
	X float64
	Y float64
}

type RumbleState struct{}

func NewState() *State {
	return &State{}
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
