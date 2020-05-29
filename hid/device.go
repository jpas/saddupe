package hid

import (
	"io"

	"github.com/hashicorp/go-multierror"
	"github.com/jpas/saddupe/hw"
	"github.com/pkg/errors"
)

type Device struct {
	control   io.ReadWriteCloser
	interrupt io.ReadWriteCloser
	local     hw.MAC
	other     hw.MAC
}

func NewDevice(control, interrupt io.ReadWriteCloser, localAddr hw.MAC, otherAddr hw.MAC) (*Device, error) {
	return &Device{control, interrupt, localAddr, otherAddr}, nil
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
	if n == 0 {
		return nil, errors.New("empty report")
	}

	return NewReport(ReportHeader(b[0]), b[1:n])
}

func (d Device) Write(r *Report) error {
	_, err := d.interrupt.Write(r.Bytes())
	return err
}

func (d Device) LocalAddr() hw.MAC {
	return d.local
}

func (d Device) OtherAddr() hw.MAC {
	return d.other
}
