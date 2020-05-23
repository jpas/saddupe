package packet

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

type RetFlashRead struct {
	Addr uint32
	Data []byte
}

func init() {
	RegisterRet(&RetFlashRead{})
}

func (r *RetFlashRead) Op() OpCode {
	return 0x10
}

func (r *RetFlashRead) Ack() bool {
	return true
}

func (r *RetFlashRead) Type() byte {
	return 0x10
}

func (r *RetFlashRead) Encode() ([]byte, error) {
	if len(r.Data) >= FlashReadMaxLen {
		return nil, errors.Errorf("data must be less than %d bytes", FlashReadMaxLen)
	}
	var b [6]byte
	binary.LittleEndian.PutUint32(b[0:4], r.Addr)
	b[4] = byte(len(r.Data))
	return append(b[:], r.Data...), nil
}

func (r *RetFlashRead) Decode(b []byte) error {
	return nil
}
