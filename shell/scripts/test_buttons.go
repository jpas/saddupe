package main

import (
	"time"

	. "github.com/jpas/saddupe/shell/env"
)

func main() {
	State.B.Tap(time.Second / 20)
}
