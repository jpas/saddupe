package state

import (
	"math/cmplx"
)

type Stick struct {
	Button
	pos complex128
}

func (s *Stick) X() float64 {
	return real(s.pos)
}

func (s *Stick) Y() float64 {
	return imag(s.pos)
}

func (s *Stick) Pos() (float64, float64) {
	return s.X(), s.Y()
}

func (s *Stick) Set(x, y float64) {
	radius, angle := cmplx.Polar(complex(x, y))
	s.SetPolar(radius, angle)
}

func (s *Stick) Polar() (float64, float64) {
	return cmplx.Polar(s.pos)
}

func (s *Stick) SetPolar(radius, angle float64) {
	if radius > 1 {
		radius = 1
	}
	s.pos = cmplx.Rect(radius, angle)
}
