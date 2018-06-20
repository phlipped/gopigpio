package gopigpio

import (
	"io"
)

// Command IDs
const (
	TICK = 16
	HW_REVISION = 17
	VERSION = 26
)

func HardwareRevision(p io.ReadWriter) (uint32, error) {
  cmd := Cmd{
	  ID: HW_REVISION,
  }
  res, err := sendCmd(p, cmd)
  return uint32(res), err
}

func Version(p io.ReadWriter) (uint32, error) {
  cmd := Cmd{
	  ID: VERSION,
  }
  res, err := sendCmd(p, cmd)
  return uint32(res), err
}

func Tick(p io.ReadWriter) (uint32, error) {
  cmd := Cmd{
	  ID: TICK,
  }
  res, err := sendCmd(p, cmd)
  return uint32(res), err
}
