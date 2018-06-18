package gopigpio

type CmdID uint32

type Cmd struct {
  Cmd CmdID
  P1 uint32
  P2 uint32
}

type Result struct {
  cmd Cmd
  ext []byte
}

func (p Pigpio) sendCmdSimple(cmd Cmd) (Result, error) {
	// FIXME implement
	_ = p
	_ = cmd

	return Result{}, nil
}

func (p Pigpio) sendCmdExtended(cmd Cmd, ext []byte) (Result, error) {
	// FIXME implement
	_ = p
	_ = cmd
	_ = ext

	return Result{}, nil
}
