package gopigpio

import (
	"fmt"
	"io"
)

const (
  INPUT uint32 = 0
  OUTPUT uint32 = 1
  ALT0 uint32 = 4
  ALT1 uint32 = 5
  ALT2 uint32 = 6
  ALT3 uint32 = 7
  ALT4 uint32 = 3
  ALT5 uint32 = 2
)

// Command IDs for GPIOs
const (
  GPIO_SET_MODE = 0
)

func GpioSetMode(p io.ReadWriter, gpio uint32, mode uint32) error {
  cmd := Cmd{
	  ID: GPIO_SET_MODE,
	  P1: gpio,
          P2: mode,
  }
  res, err := sendCmd(p, cmd)
  if err != nil {
	  return err
  }

  if res < 0 {
	  return fmt.Errorf("Error from GpioSetMode(gpio=%d, mode=%d): Error code %d (see pigpio documentation for meaning of error code)", gpio, mode, res)
  }

  return nil

}
