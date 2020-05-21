package hid

import (
	"io"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type Device struct {
	control   io.ReadWriteCloser
	interrupt io.ReadWriteCloser
	MAC       [6]byte
}

func NewDevice(control, interrupt io.ReadWriteCloser, mac [6]byte) (*Device, error) {
	return &Device{control, interrupt, mac}, nil
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
	var b [1024]byte // hopefully this is big enough
	n, err := d.interrupt.Read(b[:])
	if err != nil {
		return nil, err
	}
	return NewReport(b[:n])
}

func (d Device) Write(r *Report) error {
	_, err := d.interrupt.Write(r.Bytes())
	return err
}
