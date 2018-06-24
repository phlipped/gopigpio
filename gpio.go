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
	GPIO_LOW   PinVal = 0
	GPIO_CLEAR PinVal = 0
	GPIO_OFF   PinVal = 0
	GPIO_HIGH  PinVal = 1
	GPIO_SET   PinVal = 1
	GPIO_ON    PinVal = 1
)

type PinPull uint32

const (
	GPIO_PULL_OFF PinPull = 0
	GPIO_PULL_UP PinPull = 1
	GPIO_PULL_HIGH PinPull = 1
	GPIO_PULL_DOWN PinPull = 2
	GPIO_PULL_LOW PinPull = 2
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
	_ = gpio
	_ = val
	return nil
}

func GpioSetPullUpDown(p io.ReadWriter, gpio Pin, pull PinPull) error {
	_ = gpio
	_ = pull
	return nil
}
