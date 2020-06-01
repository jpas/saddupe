package state

import "time"

type Button struct {
	pressed bool
	start   time.Time
	millis  uint64
}

func (b *Button) Held() bool {
	return b.pressed
}

func (b *Button) SetHeld(p bool) {
	if p {
		b.Hold()
	} else {
		b.Release()
	}
}

func (b *Button) Hold() {
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

func (b *Button) Tap(t time.Duration) {
	b.Hold()
	time.Sleep(t)
	b.Release()
}
