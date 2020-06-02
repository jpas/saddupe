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
	b[2] |= bit(0, s.Y.Held())
	b[2] |= bit(1, s.X.Held())
	b[2] |= bit(2, s.B.Held())
	b[2] |= bit(3, s.A.Held())
	b[2] |= bit(6, s.R.Held())
	b[2] |= bit(7, s.ZR.Held())

	b[3] = 0
	b[3] |= bit(0, s.Minus.Held())
	b[3] |= bit(1, s.Plus.Held())
	b[3] |= bit(2, s.RightStick.Held())
	b[3] |= bit(3, s.LeftStick.Held())
	b[3] |= bit(4, s.Home.Held())
	b[3] |= bit(5, s.Capture.Held())

	b[4] = 0
	b[4] |= bit(0, s.Down.Held())
	b[4] |= bit(1, s.Up.Held())
	b[4] |= bit(2, s.Right.Held())
	b[4] |= bit(3, s.Left.Held())
	b[4] |= bit(6, s.L.Held())
	b[4] |= bit(7, s.ZL.Held())

	switch s.Kind() {
	case state.SadRight:
		b[2] |= bit(4, s.SR.Held())
		b[2] |= bit(5, s.SL.Held())
		b[3] |= bit(6, s.ChargeGrip.Held())
	case state.SadLeft:
		b[3] |= bit(6, s.ChargeGrip.Held())
		b[4] |= bit(4, s.SR.Held())
		b[4] |= bit(5, s.SL.Held())
	}

	err := encodeStick(b[5:], &s.LeftStick, s.LeftStickCalibration())
	if err != nil {
		return err
	}

	err = encodeStick(b[8:], &s.RightStick, s.RightStickCalibration())
	if err != nil {
		return err
	}

	// TODO(jpas) Rumble report

	return nil
}

func decodeState(b []byte, s *state.State) error {
	return errors.New("not implemented")
}

func translateRange(t, a, b, c, d float64, round func(float64) float64) float64 {
	// [a, b] -> [c, d]
	if t < a {
		return c
	}
	if t > b {
		return d
	}
	v := c + (d-c)/(b-a)*(t-a)
	if round == nil {
		return v
	}
	return round(v)
}

func encodeAxis(t float64, axis *state.AxisCalibration) uint16 {
	var r float64
	switch {
	case t == 0:
		return axis.Center
	case t < 0:
		r = translateRange(
			t,
			-1, 0,
			float64(axis.Min), float64(axis.Center),
			math.Ceil,
		)
	case t > 0:
		r = translateRange(
			t,
			0, 1,
			float64(axis.Min), float64(axis.Center),
			math.Floor,
		)
	}
	return uint16(r) & 0x0fff
}

func decodeAxis(t uint16, axis *state.AxisCalibration) float64 {
	var r float64
	switch {
	case t < axis.Center:
		r = translateRange(
			float64(t),
			float64(axis.Min), float64(axis.Center),
			-1, 0,
			nil,
		)
	case t > axis.Center:
		r = translateRange(
			float64(t),
			float64(axis.Center), float64(axis.Max),
			0, 1,
			nil,
		)
	}
	return r
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
