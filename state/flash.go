package state

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

type Flash [0x10000]byte

func NewFlash() *Flash {
	var f Flash
	f[0x6013] = 0xa0 // unknown value
	return &f
}

func (f *Flash) Serial() string {
	if f[0x6000] >= 0x80 {
		return ""
	}
	return string(f[0x6000:0x6010])
}

func (f *Flash) SetSerial(s string) {
	m := f[0x6000:]
	if s == "" {
		for i := 0; i < 16; i++ {
			m[i] = 0xff
		}
	}
	copy(m[:16], []byte(s))
}

type DeviceKind byte

const (
	SadLeft  DeviceKind = 0x01
	SadRight DeviceKind = 0x02
	Pro      DeviceKind = 0x03
)

func (f *Flash) Kind() DeviceKind {
	return DeviceKind(f[0x6012])
}

func (f *Flash) SetKind(t DeviceKind) {
	f[0x6012] = byte(t)
}

func (f *Flash) HasColor() bool {
	return f[0x601b] == 0x01
}

func (f *Flash) SetHasColor(b bool) {
	if b {
		f[0x601b] = 0x01
	} else {
		f[0x601b] = 0x00
	}
}

func (f *Flash) Read(dst []byte, addr uint32, n int) error {
	end := addr + uint32(n) - 1
	if uint32(len(f)) <= end {
		return errors.New("read out of bounds")
	}
	log.Printf("flash read: [0x%05x, 0x%05x]", addr, end)
	copy(dst[:n], f[addr:])
	log.Printf("flash read: %02x", dst[:n])
	return nil
}

func (f *Flash) Write(addr uint32, src []byte) error {
	end := addr + uint32(len(src)) - 1
	if uint32(len(f)) <= end {
		return errors.New("write out of bounds")
	}
	log.Printf("flash write: [0x%05x, 0x%05x]", addr, end)
	copy(f[addr:], src)
	return nil
}

type Color struct {
	R, G, B uint8
}

var colorRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomColor() Color {
	v := colorRand.Uint32()
	return Color{
		R: byte(v),
		G: byte(v >> 8),
		B: byte(v >> 16),
	}
}

type AxisCalibration struct {
	Min, Center, Max uint16
}

type StickCalibration struct {
	X, Y AxisCalibration
}

func (f *Flash) stickCalibration(addr int) *StickCalibration {
	m := f[addr:]
	return &StickCalibration{
		X: AxisCalibration{
			Min:    (uint16(m[7])<<8)&0xf00 | uint16(m[6]),
			Center: (uint16(m[4])<<8)&0xf00 | uint16(m[3]),
			Max:    (uint16(m[1])<<8)&0xf00 | uint16(m[0]),
		},
		Y: AxisCalibration{
			Min:    uint16(m[8])<<4 | uint16(m[7])>>4,
			Center: uint16(m[5])<<4 | uint16(m[4])>>4,
			Max:    uint16(m[2])<<4 | uint16(m[1])>>4,
		},
	}
}

func (f *Flash) putStickCalibration(s *StickCalibration, addr int) {
	m := f[addr:]
	m[0] = byte(s.X.Max)
	m[1] = byte(s.Y.Max<<4)&0xf0 | byte(s.X.Max>>8)&0x0f
	m[2] = byte(s.Y.Max >> 4)
	m[3] = byte(s.X.Center)
	m[4] = byte(s.Y.Center<<4)&0xf0 | byte(s.X.Center>>8)&0x0f
	m[5] = byte(s.Y.Center >> 4)
	m[6] = byte(s.X.Min)
	m[7] = byte(s.Y.Min<<4)&0xf0 | byte(s.X.Min>>8)&0x0f
	m[8] = byte(s.Y.Min >> 4)
}

const LeftStickCalibrationBase = 0x603d

func (f *Flash) LeftStickCalibration() *StickCalibration {
	return f.stickCalibration(LeftStickCalibrationBase)
}

func (f *Flash) SetLeftStickCalibration(s *StickCalibration) {
	f.putStickCalibration(s, LeftStickCalibrationBase)
}

const RightStickCalibrationBase = 0x6046

func (f *Flash) RightStickCalibration() *StickCalibration {
	return f.stickCalibration(RightStickCalibrationBase)
}

func (f *Flash) SetRightStickCalibration(s *StickCalibration) {
	f.putStickCalibration(s, RightStickCalibrationBase)
}

func (f *Flash) BodyColour() Color {
	return Color{f[0x6050], f[0x6051], f[0x6052]}
}

func (f *Flash) SetBodyColour(c Color) {
	f[0x6050] = c.R
	f[0x6051] = c.G
	f[0x6052] = c.B
}

func (f *Flash) ButtonColour() Color {
	return Color{f[0x6053], f[0x6054], f[0x6055]}
}

func (f *Flash) SetButtonColour(c Color) {
	f[0x6053] = c.R
	f[0x6054] = c.G
	f[0x6055] = c.B
}

func (f *Flash) LeftGripColour() Color {
	return Color{f[0x6056], f[0x6057], f[0x6058]}
}

func (f *Flash) SetLeftGripColour(c Color) {
	f[0x6056] = c.R
	f[0x6057] = c.G
	f[0x6058] = c.B
}

func (f *Flash) RightGripColour() Color {
	return Color{f[0x6059], f[0x605a], f[0x605b]}
}

func (f *Flash) SetRightGripColour(c Color) {
	f[0x6059] = c.R
	f[0x605a] = c.G
	f[0x605b] = c.B
}
