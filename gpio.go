package gopigpio

import (
	"io"
)

// GPIO Modes
type GpioMode uint32

const (
  INPUT GpioMode = 0
  OUTPUT GpioMode = 1
  ALT0 GpioMode = 4
  ALT1 GpioMode = 5
  ALT2 GpioMode = 6
  ALT3 GpioMode = 7
  ALT4 GpioMode = 3
  ALT5 GpioMode = 2
)

// Command IDs for GPIOs
const (
  GPIO_SET_MODE = 0
)

func GpioSetMode(p io.ReadWriter, gpio uint32, mode GpioMode) error {
  _ = gpio
  _ = mode

  // Build cmd struct
  cmd := Cmd{
	  ID: GPIO_SET_MODE,
	  P1: gpio,
          P2: uint32(mode),
  }
  // Send cmd
  result, err := sendCmd(p, cmd)

  // Get response
  // FIXME process result and err
  _ = result
  _ = err
  // Interpret response
  return nil
}
