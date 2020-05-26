package packet

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

const (
	CmdFlashReadOp        OpCode = 0x10
	CmdFlashWriteOp       OpCode = 0x11
	CmdFlashSectorEraseOp OpCode = 0x12
	FlashOpMaxLen         int    = 0x1d
)

type CmdFlashRead struct {
	Addr uint32
	Len  int
}

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
	if len(r.Data) >= FlashOpMaxLen {
		return nil, errors.Errorf("data must be less than %d bytes", FlashOpMaxLen)
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

type CmdFlashWrite struct {
	Addr uint32
	Data []byte
}

func init() {
	RegisterCmd(&CmdFlashWrite{})
}

func (c *CmdFlashWrite) Op() OpCode {
	return CmdFlashWriteOp
}

func (c *CmdFlashWrite) Encode() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (c *CmdFlashWrite) Decode(b []byte) error {
	return errors.New("not implemented")
}

type RetFlashWrite struct {
	status byte
}

func init() {
	RegisterRet(&RetFlashWrite{})
}

func (r *RetFlashWrite) Op() OpCode {
	return CmdFlashWriteOp
}

func (r *RetFlashWrite) Ack() bool {
	return true
}

func (r *RetFlashWrite) Type() byte {
	return 0x00
}

func (c *RetFlashWrite) Encode() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (c *RetFlashWrite) Decode(b []byte) error {
	c.status = b[2]
	return nil
}

// Little-endian address, int8 size. Max size x1D data to write. Replies with x8011 ack and a uint8 status. x00 = success, x01 = write protected.
type CmdFlashSectorErase struct {
	Addr uint32
}

func init() {
	RegisterCmd(&CmdFlashSectorErase{})
}

func (c *CmdFlashSectorErase) Op() OpCode {
	return CmdFlashSectorEraseOp
}

func (c *CmdFlashSectorErase) Encode() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (c *CmdFlashSectorErase) Decode(b []byte) error {
	return errors.New("not implemented")
}

type RetFlashSectorErase struct {
	status byte
}

func init() {
	RegisterRet(&RetFlashSectorErase{})
}

func (r *RetFlashSectorErase) Op() OpCode {
	return CmdFlashSectorEraseOp
}

func (r *RetFlashSectorErase) Ack() bool {
	return true
}

func (r *RetFlashSectorErase) Type() byte {
	return 0x00
}

func (c *RetFlashSectorErase) Encode() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (c *RetFlashSectorErase) Decode(b []byte) error {
	c.status = b[2]
	return nil
}
