package packet

import "encoding/binary"

type CmdFlashRead struct {
	Addr uint32
	Len  int
}

const FlashReadMaxLen = 0x1d

func init() {
	RegisterCmd(&CmdFlashRead{})
}

func (c *CmdFlashRead) Op() OpCode {
	return 0x10
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
