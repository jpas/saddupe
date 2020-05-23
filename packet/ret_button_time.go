package packet

import (
	"encoding/binary"
)

type RetButtonTime struct {
	L    uint16
	R    uint16
	ZL   uint16
	ZR   uint16
	SL   uint16
	SR   uint16
	Home uint16
}

func init() {
	RegisterRet(&RetButtonTime{})
}

func (r *RetButtonTime) Op() OpCode {
	return 0x04
}

func (r *RetButtonTime) Ack() bool {
	return true
}

func (r *RetButtonTime) Type() byte {
	return 0x10
}

func (r *RetButtonTime) Encode() ([]byte, error) {
	var b [14]byte
	binary.LittleEndian.PutUint16(b[0:], r.L)
	binary.LittleEndian.PutUint16(b[2:], r.R)
	binary.LittleEndian.PutUint16(b[4:], r.ZL)
	binary.LittleEndian.PutUint16(b[6:], r.ZR)
	binary.LittleEndian.PutUint16(b[8:], r.SL)
	binary.LittleEndian.PutUint16(b[10:], r.SR)
	binary.LittleEndian.PutUint16(b[12:], r.Home)
	return b[:], nil
}

func (r *RetButtonTime) Decode(b []byte) error {
	return nil
}
