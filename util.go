package gopigpio

import (
	"io"
)

// Command IDs
const (
	TICK        CmdID = 16
	HW_REVISION CmdID = 17
	VERSION     CmdID = 26
)

// N.B. Pigpio will never return a failure for these functions, which means
// we will never interpret <res> as an error. Errors will only be returned
// for infrastructure errors (e.g. writing to / reading from <p>)
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
