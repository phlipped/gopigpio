package gopigpio

import (
	"io"
)

type Pigpio struct {
	io.ReadWriteCloser
}
