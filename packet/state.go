package packet

import (
	"github.com/jpas/saddupe/hid"
	"github.com/jpas/saddupe/state"
	"github.com/pkg/errors"
)

type StatePacket struct {
	State state.State
}

func init() {
	RegisterPacket(&StatePacket{})
}

func (p *StatePacket) Header() hid.ReportHeader {
	return hid.InputReportHeader
}

func (p *StatePacket) ID() PacketID {
	return 0x30
}

func (p *StatePacket) Encode() ([]byte, error) {
	var b [50]byte

	b[0] = byte(p.ID())

	err := encodeState(b[1:], &p.State)
	if err != nil {
		return nil, errors.Wrap(err, "state encode faild")
	}
	return b[:], nil
}

func (p *StatePacket) Decode(b []byte) error {
	err := decodeState(b[1:], &p.State)
	if err != nil {
		return errors.Wrap(err, "status decode faild")
	}
	return nil
}

func encodeState(b []byte, s *state.State) error {
	b[0] = byte(s.Tick)

	// 0xe for joycon
	// 0x0 for procon or charge grip joycons
	// We are always a procon so just hard code it...
	b[1] = byte(s.Battery.Level) | 0x0
	if s.Battery.Charging {
		b[1] |= 0x10
	}

	b[2] |= bit(0, s.Y.Pressed())
	b[2] |= bit(1, s.X.Pressed())
	b[2] |= bit(2, s.B.Pressed())
	b[2] |= bit(3, s.A.Pressed())
	b[2] |= bit(6, s.R.Pressed())
	b[2] |= bit(7, s.ZR.Pressed())

	b[3] |= bit(0, s.Minus.Pressed())
	b[3] |= bit(1, s.Plus.Pressed())
	b[3] |= bit(2, s.RightStick.Pressed())
	b[3] |= bit(3, s.LeftStick.Pressed())
	b[3] |= bit(4, s.Home.Pressed())
	b[3] |= bit(5, s.Capture.Pressed())

	b[4] |= bit(0, s.Down.Pressed())
	b[4] |= bit(0, s.Up.Pressed())
	b[4] |= bit(0, s.Right.Pressed())
	b[4] |= bit(0, s.Left.Pressed())
	b[4] |= bit(0, s.L.Pressed())
	b[4] |= bit(0, s.ZL.Pressed())

	return nil
}

func decodeState(b []byte, s *state.State) error {
	return nil
}
