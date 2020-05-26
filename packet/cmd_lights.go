package packet

type CmdLightsSet struct {
}

const CmdLightsSetOp OpCode = 0x30

func init() {
	RegisterCmd(&CmdLightsSet{})
}

func (c *CmdLightsSet) Op() OpCode {
	return CmdLightsSetOp
}

func (c *CmdLightsSet) Encode() ([]byte, error) {
	b := []byte{
		byte(c.Op()),
	}
	return b, nil
}

func (c *CmdLightsSet) Decode(b []byte) error {
	return nil
}

func init() {
	RegisterRet(&RetAck{op: CmdLightsSetOp})
}
