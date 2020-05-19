package main

import (
	"io"
	"log"

	"github.com/jpas/saddupe/l2"
	"github.com/pkg/errors"
)

type Dupe struct {
	ctrl io.Closer
	intr io.ReadWriteCloser
}

func NewDupe(console string) (*Dupe, error) {
	ctrl, err := l2.NewConn(console, 17)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open control channel")
	}

	intr, err := l2.NewConn(console, 19)
	if err != nil {
		ctrl.Close()
		return nil, errors.Wrap(err, "unable to open interrupt channel")
	}

	return &Dupe{io.Closer(ctrl), io.ReadWriteCloser(intr)}, nil
}

func (d *Dupe) Run() {
	var buf [1024]byte
	for {
		n, err := d.intr.Read(buf[:])
		if err != nil {
			panic(err)
		}
		log.Printf("recv: %v", buf[:n])
	}
}
