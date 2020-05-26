package packet

import "github.com/jpas/saddupe/state"

type CmdModeSet struct {
	Mode state.Mode
}

const CmdModeSetOp OpCode = 0x03

func init() {
	RegisterCmd(&CmdModeSet{})
}

func (c *CmdModeSet) Op() OpCode {
	return CmdModeSetOp
}

func (c *CmdModeSet) Encode() ([]byte, error) {

	b := []byte{
		byte(c.Op()),
		byte(c.Mode),
	}
	return b, nil
}

func (c *CmdModeSet) Decode(b []byte) error {
	c.Mode = state.Mode(b[1])
	return nil
}

func init() {
	RegisterRet(&RetAck{op: CmdModeSetOp})
}
