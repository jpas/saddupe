package hid

import (
	"io"

	"github.com/pkg/errors"
)

type Channel struct {
	rwc io.ReadWriteCloser
	mtu int
}

func NewChannel(rwc io.ReadWriteCloser, mtu int) (*Channel, error) {
	var c Channel

	c.rwc = rwc

	if err := c.SetMTU(mtu); err != nil {
		return nil, errors.Wrap(err, "failed to set mtu")
	}

	return &Channel{rwc, mtu}, nil
}

const MinimumMTU = 672

func (c *Channel) SetMTU(mtu int) error {
	if mtu < MinimumMTU {
		return errors.Errorf("mtu must be at least %d", MinimumMTU)
	}
	c.mtu = mtu
	return nil
}

func (c Channel) Read() (*Report, error) {
	buf := make([]byte, c.mtu)
	n, err := c.rwc.Read(buf)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}
	return NewReport(buf[:n])
}

func (c Channel) Write(r *Report) error {
	if _, err := c.rwc.Write(r.Bytes()); err != nil {
		return errors.Wrap(err, "write failed")
	}
	return nil
}

func (c Channel) Close() error {
	if err := c.rwc.Close(); err != nil {
		return errors.Wrap(err, "close failed")
	}
	return nil
}
