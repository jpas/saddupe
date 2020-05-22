package state

type State struct {
	Tick     uint64
	Battery  BatteryState
	Charging bool
	HasGrip  bool
	Powered  bool

	Buttons    ButtonsState
	LeftStick  StickState
	RightStick StickState
	Vibration  VibrationState
}

type BatteryState struct{}
type StickState struct{}
type ButtonsState struct{}
type VibrationState struct{}

func (s *State) Encode() ([]byte, error) {
	b := [12]byte{}
	return b[:], nil
}

func (s *State) Decode([]byte) error {
	return nil
}

type BasicState struct {
	Up, Down, Left, Right      bool
	SL, SR, LR, ZLR            bool
	Minus, Plus, Home, Capture bool
	LeftStick, RightStick      bool
	Hat                        HatDirection
}

type HatDirection byte

// To ensure that HatCenter is the zero value we use 0x08 - hat instead.
const (
	HatCenter    HatDirection = 0x00
	HatUp        HatDirection = 0x08
	HatUpRight   HatDirection = 0x07
	HatRightUp   HatDirection = 0x07
	HatRight     HatDirection = 0x06
	HatRightDown HatDirection = 0x05
	HatDownRight HatDirection = 0x05
	HatDown      HatDirection = 0x04
	HatDownLeft  HatDirection = 0x03
	HatLeftDown  HatDirection = 0x03
	HatLeft      HatDirection = 0x02
	HatLeftUp    HatDirection = 0x01
	HatUpLeft    HatDirection = 0x01
)
