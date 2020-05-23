package packet

import "github.com/jpas/saddupe/state"

type CmdSetMode struct {
	Mode state.Mode
}

func init() {
	RegisterCmd(&CmdSetMode{})
}

func (c *CmdSetMode) Op() OpCode {
	return 0x03
}

func (c *CmdSetMode) Encode() ([]byte, error) {
	b := []byte{
		byte(c.Op()),
		byte(c.Mode),
	}
	return b, nil
}

func (c *CmdSetMode) Decode(b []byte) error {
	c.Mode = state.Mode(b[1])
	return nil
}
