package packet

const (
	SimpleButtonStatusID   PacketID = InputPacketID | 0x3F
	SimpleButtonStatusSize int      = 12
)

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

func (s SimpleButtonStatus) Pack() ([]byte, error) {
	var p [SimpleButtonStatusSize]byte

	err := SimpleButtonStatusID.packInto(p[:2])
	if err != nil {
		return nil, nil
	}

	p[2] = boolsToByte(s.Down, s.Right, s.Left, s.Up, s.SL, s.SR, false, false)
	p[3] = boolsToByte(s.Minus, s.Plus, s.LeftStick, s.RightStick, s.Home, s.Capture, s.LR, s.ZLR)
	p[4] = byte((s.Hat - 0x08) & 0x0f)
	copy(p[5:], []byte{0x00, 0x80, 0x00, 0x80, 0x00, 0x80, 0x00, 0x80})

	return p[:], nil
}

func unpackSimpleButtonStatus(p []byte) (Packet, error) {
	var s SimpleButtonStatus

	if len(p) < SimpleButtonStatusSize {
		return nil, ErrTooShort
	}

	s.Down, s.Right, s.Left, s.Up, s.SL, s.SR, _, _ = byteToBools(p[2])
	s.Minus, s.Plus, s.LeftStick, s.RightStick, s.Home, s.Capture, s.LR, s.ZLR = byteToBools(p[3])
	s.Hat = HatDirection(0x08 - (p[4] & 0x0f))

	// We don't have to do anything with the padding!

	return &s, nil
}
