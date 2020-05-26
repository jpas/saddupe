package packet

import (
	"fmt"

	"github.com/jpas/saddupe/hid"
	"github.com/jpas/saddupe/state"
	"github.com/pkg/errors"
)

type BasicStatePacket struct {
	State state.BasicState
}

func init() {
	RegisterPacket(&BasicStatePacket{})
}

func (p *BasicStatePacket) Header() hid.ReportHeader {
	return hid.InputReportHeader
}

func (p *BasicStatePacket) ID() PacketID {
	return 0x3f
}

func (p *BasicStatePacket) Encode() ([]byte, error) {
	s := p.State
	if 0x08 < s.Hat {
		return nil, errors.New("invalid hat direction")
	}

	b := [...]byte{
		byte(p.ID()),
		boolsToByte(s.Down, s.Right, s.Left, s.Up, s.SL, s.SR, false, false),
		boolsToByte(s.Minus, s.Plus, s.LeftStick, s.RightStick, s.Home, s.Capture, s.LR, s.ZLR),
		0x08 - byte(s.Hat),
		0x00, 0x80,
		0x00, 0x80,
		0x00, 0x80,
		0x00, 0x80,
	}
	return b[:], nil
}

func (p *BasicStatePacket) Decode(b []byte) error {
	fmt.Println(b)
	s := &p.State
	if PacketID(b[0]) != p.ID() {
		return errors.New("invalid id")
	}

	s.Down, s.Right, s.Left, s.Up, s.SL, s.SR, _, _ = byteToBools(b[1])
	s.Minus, s.Plus, s.LeftStick, s.RightStick, s.Home, s.Capture, s.LR, s.ZLR = byteToBools(b[2])

	if 0x08 < b[3] {
		return errors.New("invalid hat direction")
	}
	s.Hat = state.HatDirection(0x08 - b[3])

	return nil
}
