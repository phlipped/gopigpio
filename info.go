package gopigpio

import (
	"io"
)

const (
  HW_REVISION = 17
)

func HardwareRevision(p io.ReadWriter) (uint, error) {
  cmd := Cmd{
    ID: HW_REVISION,
    P1: 0,
    P2: 0,
  }

  result, err := sendCmd(p, cmd)
  if err != nil {
	return 0, err
  }

  return uint(result.res), nil
}
