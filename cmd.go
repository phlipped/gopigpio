package gopigpio

import (
	"io"
)
/*
typedef struct
{
   uint32_t cmd;
   uint32_t p1;
   uint32_t p2;
   union
   {
      uint32_t p3;
      uint32_t ext_len;
      uint32_t res;
   };
} cmdCmd_t;
*/

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

func sendCmdSimple(s io.ReadWriter, cmd Cmd) (Result, error) {
	_ = s
	_ = cmd

	return Result{}, nil
}
