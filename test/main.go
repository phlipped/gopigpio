package main

import (
	"fmt"
	"github.com/phlipped/gopigpio"
	"net"
	"time"
)

const (
	TEST_PIN = gopigpio.Pin(18)

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

	fmt.Printf("\nTesting notifications ...\n")
	h, err := gopigpio.NotifyOpen(s)
	if err != nil {
		panic(fmt.Sprintf("Failed to get notification handle: err is '%s'", err))
	}
	defer func() {
		if res, err := gopigpio.NotifyClose(s, h); err != nil {
			fmt.Printf("Failed to close notification handle: err is '%s', res is '%d'\n", err, res)
		}
	}()

	notifyFile, err := gopigpio.OpenNotificationFile(h)
	if err != nil {
		panic(fmt.Sprintf("Failed to open notification file for handle '%d': err is '%s'", h, err))
	}
	defer func() {
		if err := notifyFile.Close(); err != nil {
			fmt.Printf("Failed to close notification file: err is '%s'\n", err)
		}
	}()

	// Now start monitoring level changes on some Pin
	res, err := gopigpio.NotifyBegin(s, h, []gopigpio.Pin{TEST_PIN})
	if err != nil || res != 0 {
		panic(fmt.Sprintf("Failure when calling NotifyBegin: err is '%s', res is '%d'", err, res))
	}

	terminateChan := make(chan struct{})
	errorsChan := make(chan error)

	notificationChan := gopigpio.ReadNotificationsFromHandle(h, terminateChan, errorsChan)

	// Now cause some level changes on that Pin
	err = gopigpio.GpioSetMode(s, TEST_PIN, gopigpio.OUTPUT)
	if err != nil {
		panic(fmt.Sprintf("Failure when calling gopigpio.GpioSetMode: err is '%s'", err))

	}

	err = gopigpio.GpioWrite(s, TEST_PIN, gopigpio.GPIO_HIGH)
	if err != nil {
		panic(fmt.Sprintf("Failure when calling gopigpio.Write: err is '%s'", err))
	}

	// FIXME consider adding a delay in here
	err = gopigpio.GpioWrite(s, TEST_PIN, gopigpio.GPIO_LOW)
	if err != nil {
		panic(fmt.Sprintf("Failure when calling gopigpio.Write: err is '%s'", err))
	}

	// FIXME consider adding a delay in here
	err = gopigpio.GpioWrite(s, TEST_PIN, gopigpio.GPIO_HIGH)
	if err != nil {
		panic(fmt.Sprintf("Failure when calling gopigpio.Write: err is '%s'", err))
	}

	timeout := time.After(2 * time.Second)
	done := false
	for !done {
		select {
		case _ = <-timeout:
			fmt.Printf("Timeout hit - exiting\n")
			terminateChan <- struct{}{}
			timeout = nil
		case err := <-errorsChan:
			fmt.Printf("Got error from notification reader: '%s'\n", err)
		case notification, ok := <-notificationChan:
			if !ok {
				fmt.Printf("NotificationChan closed\n")
				done = true
				continue
			}
			fmt.Printf("Notification: %v\n", notification)
		}

	}
}
