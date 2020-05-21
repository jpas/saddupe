package hid

import (
	"io"

	"github.com/jpas/saddupe/hid/internal/l2"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type Device struct {
	control   *Channel
	interrupt *Channel
	MAC       [6]byte
}

func NewDevice(control, interrupt io.ReadWriteCloser, mac string) (*Device, error) {
	addr, err := l2.NewAddr(mac, 0)
	if err != nil {
		return nil, errors.Wrap(err, "invalid mac")
	}

	ctrl, err := NewChannel(control, MinimumMTU)
	if err != nil {
		return nil, errors.Wrap(err, "unable to wrap control channel")
	}

	intr, err := NewChannel(interrupt, MinimumMTU)
	if err != nil {
		return nil, errors.Wrap(err, "unable to wrap interrupt channel")
	}

	return &Device{ctrl, intr, addr.MAC}, nil
}

func (d Device) Close() error {
	var result error

	if err := d.interrupt.Close(); err != nil {
		err = errors.Wrap(err, "interrupt close failed")
		result = multierror.Append(result, err)
	}

	if err := d.control.Close(); err != nil {
		err = errors.Wrap(err, "control close failed")
		result = multierror.Append(result, err)
	}

	return result
}

// Ignore control channel for now, but might be needed later
func (d Device) Read() (*Report, error) {
	return d.interrupt.Read()
}

func (d Device) Write(r *Report) error {
	return d.interrupt.Write(r)
}
