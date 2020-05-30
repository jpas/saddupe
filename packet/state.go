package packet

import (
	"math"

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

	b[1] = byte(s.Battery.Level) & 0xf0
	if s.Battery.Charging {
		b[1] |= 0x10
	}

	switch s.Kind() {
	case state.SadLeft, state.SadRight:
		b[1] |= 0x0e
	case state.Pro:
		// TODO(jpas) also for charge grip, but we don't have state for that...
		b[1] |= 0x00
	}

	b[2] = 0
	b[2] |= bit(0, s.Y.Pressed())
	b[2] |= bit(1, s.X.Pressed())
	b[2] |= bit(2, s.B.Pressed())
	b[2] |= bit(3, s.A.Pressed())
	b[2] |= bit(6, s.R.Pressed())
	b[2] |= bit(7, s.ZR.Pressed())

	b[3] = 0
	b[3] |= bit(0, s.Minus.Pressed())
	b[3] |= bit(1, s.Plus.Pressed())
	b[3] |= bit(2, s.RightStick.Pressed())
	b[3] |= bit(3, s.LeftStick.Pressed())
	b[3] |= bit(4, s.Home.Pressed())
	b[3] |= bit(5, s.Capture.Pressed())

	b[4] = 0
	b[4] |= bit(0, s.Down.Pressed())
	b[4] |= bit(1, s.Up.Pressed())
	b[4] |= bit(2, s.Right.Pressed())
	b[4] |= bit(3, s.Left.Pressed())
	b[4] |= bit(6, s.L.Pressed())
	b[4] |= bit(7, s.ZL.Pressed())

	switch s.Kind() {
	case state.SadRight:
		b[2] |= bit(4, s.SR.Pressed())
		b[2] |= bit(5, s.SL.Pressed())
		b[3] |= bit(6, s.ChargeGrip.Pressed())
	case state.SadLeft:
		b[3] |= bit(6, s.ChargeGrip.Pressed())
		b[4] |= bit(4, s.SR.Pressed())
		b[4] |= bit(5, s.SL.Pressed())
	}

	encodeStick(b[5:], &s.LeftStick, s.LeftStickCalibration())
	encodeStick(b[8:], &s.RightStick, s.RightStickCalibration())

	// TODO(jpas) Rumble report

	return nil
}

func decodeState(b []byte, s *state.State) error {
	return errors.New("not implemented")
}

func encodeAxis(t float64, axis *state.AxisCalibration) uint16 {
	// [a, b] -> [c, d]
	var a, b, c, d float64

	// We want to round towards the center
	var round func(float64) float64

	switch {
	case t == 0:
		return axis.Center
	case t < 0:
		a, b, c, d = -1, 0, float64(axis.Min), float64(axis.Center)
		round = math.Ceil
	case t > 0:
		a, b, c, d = 0, 1, float64(axis.Center), float64(axis.Max)
		round = math.Floor
	}

	r := round(c + (d-c)/(b-a)*(t-a))
	return uint16(r) & 0x0fff
}

func encodeStick(b []byte, s *state.Stick, c *state.StickCalibration) error {
	x0, y0 := s.Pos()
	x := encodeAxis(x0, &c.X)
	y := encodeAxis(y0, &c.Y)
	b[0] = byte(x)
	b[1] = byte(y<<4)&0xf0 | byte(x>>8)&0x0f
	b[2] = byte(y >> 4)
	return nil
}

func decodeStick(b []byte, s *state.Stick, c *state.StickCalibration) error {
	return errors.New("not implemented")
}
