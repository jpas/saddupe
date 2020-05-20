package packet

const SimpleButtonStatusID = InputPacketID | 0x3F

func init() {
	Unpackers[SimpleButtonStatusID] = unpackSimpleButtonStatus
}

type SimpleButtonStatus struct {
	Up, Down, Left, Right      bool
	SL, SR, LR, ZLR            bool
	Minus, Plus, Home, Capture bool
	LeftStick, RightStick      bool
	Hat                        HatDirection
}

type HatDirection byte

// We want HatCenter to be the default value, but the protocol actually uses 0x08 for so we add 1 to everything else
const (
	HatCenter    HatDirection = 0x00
	HatUp        HatDirection = 0x01
	HatUpRight   HatDirection = 0x02
	HatRightUp   HatDirection = 0x02
	HatRight     HatDirection = 0x03
	HatRightDown HatDirection = 0x04
	HatDownRight HatDirection = 0x04
	HatDown      HatDirection = 0x05
	HatDownLeft  HatDirection = 0x06
	HatLeftDown  HatDirection = 0x06
	HatLeft      HatDirection = 0x07
	HatLeftUp    HatDirection = 0x08
	HatUpLeft    HatDirection = 0x08
)

func (s SimpleButtonStatus) Pack() ([]byte, error) {
	var p [12]byte

	err := SimpleButtonStatusID.packInto(p[:2])
	if err != nil {
		return nil, nil
	}

	p[2] = byte(boolsToBits(s.Down, s.Right, s.Left, s.Up, s.SL, s.SR))
	p[3] = byte(boolsToBits(s.Minus, s.Plus, s.LeftStick, s.RightStick, s.Home, s.Capture, s.LR, s.ZLR))

	if s.Hat == HatCenter {
		p[4] = 0x08
	} else {
		p[4] = byte(s.Hat) - 1
	}

	filler := [...]byte{0x00, 0x80, 0x00, 0x80, 0x00, 0x80, 0x00, 0x80}
	copy(p[5:], filler[:])

	return p[:], nil
}

func unpackSimpleButtonStatus(p []byte) (*Packet, error) {
	return nil, nil
}
