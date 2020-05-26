package state

import "log"

type Flash struct {
	mem [0x10000]byte
}

func NewFlash() *Flash {
	var f Flash
	f.Reset()
	return &f
}

func (f *Flash) Read(b []byte, addr uint32, len int) error {
	log.Printf("flash read: %d (%d)", addr, len)
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
