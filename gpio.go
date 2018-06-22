package gopigpio

import (
	"fmt"
	"io"
)

type Pin uint32
type PinMode uint32

const (
	INPUT  PinMode = 0
	OUTPUT PinMode = 1
	ALT0   PinMode = 4
	ALT1   PinMode = 5
	ALT2   PinMode = 6
	ALT3   PinMode = 7
	ALT4   PinMode = 3
	ALT5   PinMode = 2
)

type PinVal uint32

const (
	GPIO_LOW  uint32 = 0
	GPIO_HIGH uint32 = 1

	GPIO_CLEAR uint32 = 0
	GPIO_SET   uint32 = 1

	GPIO_OFF uint32 = 0
	GPIO_ON  uint32 = 1
)

// Command IDs for GPIOs
const (
	GPIO_SET_MODE CmdID = 0
)

func GpioSetMode(p io.ReadWriter, gpio Pin, mode PinMode) error {
	cmd := Cmd{
		ID: GPIO_SET_MODE,
		P1: uint32(gpio),
		P2: uint32(mode),
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

func GpioWrite(p io.ReadWriter, gpio Pin, val PinVal) error {
	return nil
}
