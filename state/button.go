package state

import "time"

type Button struct {
	pressed bool
	start   time.Time
	millis  uint64
}

func (b *Button) Pressed() bool {
	return b.pressed
}

func (b *Button) SetPressed(p bool) {
	if p {
		b.Press()
	} else {
		b.Release()
	}
}

func (b *Button) Press() {
	if !b.pressed {
		b.start = time.Now()
	}
	b.pressed = true
}

func (b *Button) Release() {
	if b.pressed {
		b.millis += uint64(time.Since(b.start).Milliseconds())
	}
	b.pressed = false
}

func (b *Button) Milliseconds() uint64 {
	return b.millis
}
