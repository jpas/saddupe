package packet

import "encoding/binary"

type CmdButtonTime struct{}

const CmdButtonTimeOp OpCode = 0x04

func init() {
	RegisterCmd(&CmdButtonTime{})
}

func (c *CmdButtonTime) Op() OpCode {
	return CmdButtonTimeOp
}

func (c *CmdButtonTime) Encode() ([]byte, error) {
	b := []byte{byte(c.Op())}
	return b, nil
}

func (c *CmdButtonTime) Decode(b []byte) error {
	return nil
}

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
	return CmdButtonTimeOp
}

func (r *RetButtonTime) Ack() bool {
	return true
}

func (r *RetButtonTime) Type() byte {
	return 0x10
}

func (r *RetButtonTime) Encode() ([]byte, error) {
	b := make([]byte, 16)
	putRetHeader(b, r)
	binary.LittleEndian.PutUint16(b[2:], r.L)
	binary.LittleEndian.PutUint16(b[4:], r.R)
	binary.LittleEndian.PutUint16(b[6:], r.ZL)
	binary.LittleEndian.PutUint16(b[8:], r.ZR)
	binary.LittleEndian.PutUint16(b[10:], r.SL)
	binary.LittleEndian.PutUint16(b[12:], r.SR)
	binary.LittleEndian.PutUint16(b[14:], r.Home)
	return b, nil
}

func (r *RetButtonTime) Decode(b []byte) error {
	return nil
}
