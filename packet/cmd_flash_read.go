package packet

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

type CmdFlashRead struct {
	Addr uint32
	Len  int
}

const (
	CmdFlashReadOp  OpCode = 0x10
	FlashReadMaxLen int    = 0x1d
)

func init() {
	RegisterCmd(&CmdFlashRead{})
}

func (c *CmdFlashRead) Op() OpCode {
	return CmdFlashReadOp
}

func (c *CmdFlashRead) Encode() ([]byte, error) {
	var b [6]byte
	b[0] = byte(c.Op())
	binary.LittleEndian.PutUint32(b[1:], c.Addr)
	b[5] = byte(c.Len)
	return []byte{byte(c.Op())}, nil
}

func (c *CmdFlashRead) Decode(b []byte) error {
	c.Addr = binary.LittleEndian.Uint32(b[1:])
	c.Len = int(b[5])
	return nil
}

type RetFlashRead struct {
	Addr uint32
	Data []byte
}

func init() {
	RegisterRet(&RetFlashRead{})
}

func (r *RetFlashRead) Op() OpCode {
	return CmdFlashReadOp
}

func (r *RetFlashRead) Ack() bool {
	return true
}

func (r *RetFlashRead) Type() byte {
	return byte(CmdFlashReadOp)
}

func (r *RetFlashRead) Encode() ([]byte, error) {
	if len(r.Data) >= FlashReadMaxLen {
		return nil, errors.Errorf("data must be less than %d bytes", FlashReadMaxLen)
	}
	b := make([]byte, 7)
	if err := putRetHeader(b, r); err != nil {
		return nil, err
	}
	binary.LittleEndian.PutUint32(b[2:], r.Addr)
	b[6] = byte(len(r.Data))
	return append(b[:], r.Data...), nil
}

func (r *RetFlashRead) Decode(b []byte) error {
	r.Addr = binary.LittleEndian.Uint32(b[2:])
	r.Data = make([]byte, b[6])
	copy(r.Data, b[7:])
	return nil
}
