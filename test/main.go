package main

import (
	"fmt"
	"github.com/phlipped/gopigpio"
	"net"
)

func main() {
	s, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	hwVersion, err := gopigpio.HardwareRevision(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(hwVersion)

	pigpioVersion, err := gopigpio.Version(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(pigpioVersion)

	tick, err := gopigpio.Tick(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(tick)
}
