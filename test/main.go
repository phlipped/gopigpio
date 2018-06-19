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
	result, err := gopigpio.HardwareRevision(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

}
