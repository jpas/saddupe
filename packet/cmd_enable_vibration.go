package packet

type CmdEnableVibration struct {
}

const CmdEnableVibrationOp OpCode = 0x48

func init() {
	RegisterCmd(&CmdEnableVibration{})
}

func (c *CmdEnableVibration) Op() OpCode {
	return CmdEnableVibrationOp
}

func (c *CmdEnableVibration) Encode() ([]byte, error) {
	b := []byte{
		byte(c.Op()),
	}
	return b, nil
}

func (c *CmdEnableVibration) Decode(b []byte) error {
	return nil
}

func init() {
	RegisterRet(&RetAck{op: CmdEnableVibrationOp})
}
