package state

import (
	"log"
	"math/rand"
	"time"
)

type Flash struct {
	mem [0x10000]byte
}

func NewFlash() *Flash {
	var f Flash
	f.Reset()
	return &f
}

func (f *Flash) Read(b []byte, addr uint32, len int) error {
	log.Printf("flash read: [0x%05x, 0x%05x]", addr, addr+uint32(len)-1)
	m := f.mem[addr:]
	for i := 0; i < len; i++ {
		b[i] = m[i]
	}
	log.Printf("flash read: %02x", b[:len])
	return nil
}

func (f *Flash) Reset() {
	f.SetSerial("XCW10000000000")
}

func (f *Flash) Serial() string {
	if f.mem[0x6000] >= 0x80 {
		return ""
	}
	return string(f.mem[0x6000:0x6010])

}

func (f *Flash) SetSerial(s string) {
	m := f.mem[0x6000:]
	if s == "" {
		for i := 0; i < 16; i++ {
			m[i] = 0xff
		}
	}
	copy(m[:16], []byte(s))
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
	m := f.mem[addr:]
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
	m := f.mem[addr:]
	m[0] = byte(s.X.Max)
	m[1] = byte(s.Y.Max<<4)&0xf0 | byte(s.X.Max>>8)&0x0f
	m[2] = byte(s.Y.Max >> 4)
	m[3] = byte(s.X.Center)
	m[4] = byte(s.Y.Center<<4)&0xf0 | byte(s.X.Max>>8)&0x0f
	m[5] = byte(s.Y.Center >> 4)
	m[6] = byte(s.X.Min)
	m[7] = byte(s.Y.Min<<4)&0xf0 | byte(s.X.Max>>8)&0x0f
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
	return Color{f.mem[0x6050], f.mem[0x6051], f.mem[0x6052]}
}

func (f *Flash) SetBodyColour(c Color) {
	f.mem[0x6050] = c.R
	f.mem[0x6051] = c.G
	f.mem[0x6052] = c.B
}

func (f *Flash) ButtonColour() Color {
	return Color{f.mem[0x6053], f.mem[0x6054], f.mem[0x6055]}
}

func (f *Flash) SetButtonColour(c Color) {
	f.mem[0x6053] = c.R
	f.mem[0x6054] = c.G
	f.mem[0x6055] = c.B
}

func (f *Flash) LeftGripColour() Color {
	return Color{f.mem[0x6056], f.mem[0x6057], f.mem[0x6058]}
}

func (f *Flash) SetLeftGripColour(c Color) {
	f.mem[0x6056] = c.R
	f.mem[0x6057] = c.G
	f.mem[0x6058] = c.B
}

func (f *Flash) RightGripColour() Color {
	return Color{f.mem[0x6059], f.mem[0x605a], f.mem[0x605b]}
}

func (f *Flash) SetRightGripColour(c Color) {
	f.mem[0x6059] = c.R
	f.mem[0x605a] = c.G
	f.mem[0x605b] = c.B
}
