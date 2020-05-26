package packet

type CmdSetLights struct {
}

const CmdSetLightsOp OpCode = 0x30

func init() {
	RegisterCmd(&CmdSetLights{})
}

func (c *CmdSetLights) Op() OpCode {
	return CmdSetLightsOp
}

func (c *CmdSetLights) Encode() ([]byte, error) {
	b := []byte{
		byte(c.Op()),
	}
	return b, nil
}

func (c *CmdSetLights) Decode(b []byte) error {
	return nil
}

func init() {
	RegisterRet(&RetAck{op: CmdSetLightsOp})
}
