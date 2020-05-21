package hid

import (
	"github.com/jpas/saddupe/hid/internal/l2"

	"github.com/pkg/errors"
)

func Dial(proto string, addr string) (*Device, error) {
	switch proto {
	case "bt":
		return BtDeviceDial(addr)
	default:
		return nil, errors.New("unsupported proto")
	}
}

func BtDeviceDial(mac string) (*Device, error) {
	control, err := l2.NewConn(mac, 17)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open control channel")
	}

	interrupt, err := l2.NewConn(mac, 19)
	if err != nil {
		control.Close()
		return nil, errors.Wrap(err, "unable to open interrupt channel")
	}

	return NewDevice(control, interrupt, mac)
}
