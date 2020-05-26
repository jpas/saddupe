package packet

type CmdEnableSixaxis struct {
}

const CmdEnableSixaxisOp OpCode = 0x40

func init() {
	RegisterCmd(&CmdEnableSixaxis{})
}

func (c *CmdEnableSixaxis) Op() OpCode {
	return CmdEnableSixaxisOp
}

func (c *CmdEnableSixaxis) Encode() ([]byte, error) {
	b := []byte{
		byte(c.Op()),
	}
	return b, nil
}

func (c *CmdEnableSixaxis) Decode(b []byte) error {
	return nil
}

func init() {
	RegisterRet(&RetAck{op: CmdEnableSixaxisOp})
}
