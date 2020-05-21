package packet

import "github.com/jpas/saddupe/hid"

const (
	SimpleButtonStatusID  PacketID = 0x3f
	SimpleButtonStatusLen int      = 12
)

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

func (s SimpleButtonStatus) Report() (*hid.Report, error) {
	p := [SimpleButtonStatusLen]byte{
		byte(SimpleButtonStatusID),
		boolsToByte(s.Down, s.Right, s.Left, s.Up, s.SL, s.SR, false, false),
		boolsToByte(s.Minus, s.Plus, s.LeftStick, s.RightStick, s.Home, s.Capture, s.LR, s.ZLR),
		byte(0x08 - (s.Hat & 0x0f)),
		0x00, 0x80,
		0x00, 0x80,
		0x00, 0x80,
		0x00, 0x80,
	}
	return hid.NewInputReport(p[:]), nil
}
