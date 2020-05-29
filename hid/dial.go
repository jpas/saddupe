package hid

import (
	"github.com/jpas/saddupe/internal/l2"
	"github.com/pkg/errors"
)

func Dial(proto string, addr string) (*Device, error) {
	switch proto {
	case "bt":
		return BtDial(addr)
	default:
		return nil, errors.New("unsupported proto")
	}
}

func BtDial(mac string) (*Device, error) {
	control, err := l2.NewConn(mac, 17)
	if err != nil {
		return nil, errors.Wrap(err, "unable to control channel")
	}

	interrupt, err := l2.NewConn(mac, 19)
	if err != nil {
		return nil, errors.Wrap(err, "unable to interrupt channel")
	}

	localAddr := interrupt.LocalAddr().(*l2.Addr)
	otherAddr := interrupt.RemoteAddr().(*l2.Addr)

	return NewDevice(control, interrupt, localAddr.MAC, otherAddr.MAC)
}
