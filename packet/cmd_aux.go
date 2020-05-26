package packet

// TODO(jpas) requires implementation

type CmdAuxSetConfig struct{}

const CmdAuxSetConfigOp OpCode = 0x21

func init() {
	RegisterCmd(&CmdAuxSetConfig{})
}

func (c *CmdAuxSetConfig) Op() OpCode {
	return CmdAuxSetConfigOp
}

func (c *CmdAuxSetConfig) Encode() ([]byte, error) {
	b := []byte{byte(c.Op())}
	return b, nil
}

func (c *CmdAuxSetConfig) Decode(b []byte) error {
	return nil
}

type RetAuxSetConfig struct{}

func init() {
	RegisterRet(&RetAuxSetConfig{})
}

func (r *RetAuxSetConfig) Op() OpCode {
	return CmdAuxSetConfigOp
}

func (r *RetAuxSetConfig) Ack() bool {
	return true
}

func (r *RetAuxSetConfig) Type() byte {
	return 0x20
}

func (r *RetAuxSetConfig) Encode() ([]byte, error) {
	b := make([]byte, 2)
	if err := putRetHeader(b, r); err != nil {
		return nil, err
	}
	return b, nil
}

func (r *RetAuxSetConfig) Decode(b []byte) error {
	return nil
}
