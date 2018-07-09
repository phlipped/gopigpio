package gopigpio

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/npat-efault/poller"
)

const (
	NOTIFICATION_FILE_PATH_PATTERN = "/dev/pigpio%d"

)

const (
	NOTIFY_OPEN CmdID = 18
	NOTIFY_BEGIN CmdID = 19
	NOTIFY_CLOSE CmdID = 21
)

const (
	NOTIFY_FLAGS_WDOG = 1 << 5
	NOTIFY_FLAGS_ALIVE = 1 << 6
	NOTIFY_FLAGS_EVENT = 1 << 7
)

type NotifyHandle uint32

type Notification struct {
	SeqNo uint16
	WatchDogFlag bool
	AliveFlag bool
	EventFlag bool
	Tick uint32
	Levels []PinVal
}

func NotificationFromBytes(buf []byte) Notification {
	n := Notification{}
	n.SeqNo = binary.LittleEndian.Uint16(buf[0:2])

	flags := binary.LittleEndian.Uint16(buf[2:4])
	n.WatchDogFlag = (flags & NOTIFY_FLAGS_WDOG) != 0
	n.AliveFlag = (flags & NOTIFY_FLAGS_ALIVE) != 0
	n.EventFlag = (flags & NOTIFY_FLAGS_EVENT) != 0

	n.Tick = binary.LittleEndian.Uint32(buf[4:8])

	levelsBuf := binary.LittleEndian.Uint32(buf[8:12])
	for i := uint(0); i < 32; i++ {
		level := PinVal((levelsBuf>>i) & 1)
		n.Levels = append(n.Levels, level)
	}

	return n
}


func OpenNotificationFile(h NotifyHandle) (*poller.FD, error) {
	filename := fmt.Sprintf(NOTIFICATION_FILE_PATH_PATTERN, h)
	return poller.Open(filename, poller.O_RO)
}

func ReadNotificationsFromHandle(h NotifyHandle, terminate <-chan struct{}, errors chan<- error) (<-chan Notification) {
	f, err := OpenNotificationFile(h)
	if err != nil {
		errors<- err
	}

	// goroutine that Waits for a terminate signal, and closes f when it sees that notification
	go func() {
		_ = <-terminate
		if err := f.Close(); err != nil {
			errors<- err
		}
	}()

	notificationChan := make(chan Notification, 5) // Fixme make it a bigger buffer - may as well, right?

	go func() {
		for {
			// Keeps reading notifications from f - f will be closed by a separate goroutine
			notificationBuffer := make([]byte, 12)
			n, err := io.ReadFull(f, notificationBuffer)
			if err != nil {
				if n != 0 {
					// FIXME make custom error that indicates how many bytes were read so the notification handle can be cleaned.
					errors<- err
				}
				close(notificationChan)
				return
			}
			// Build a notification object and send it out the output channel
			notificationChan<- NotificationFromBytes(notificationBuffer)
		}
	}()

	return notificationChan
}


func NotifyOpen(p io.ReadWriter) (NotifyHandle, error) {
	cmd := Cmd{
		ID: NOTIFY_OPEN,
	}

	res, err := sendCmd(p, cmd)
	h := NotifyHandle(res)
	if err != nil {
		return h, err
	}
	if res < 0 {
		return h, fmt.Errorf("Error while executing NotifyOpen command: Error Code was %d", res)
	}
	return h, nil
}

func NotifyBegin(p io.ReadWriter, h NotifyHandle, pins []Pin) (int32, error) {
	var bits uint32
	for _, pin := range pins {
		bits |= 1 << pin
	}
	cmd := Cmd{
		ID: NOTIFY_BEGIN,
		P1: uint32(h),
		P2: bits,
	}

	res, err := sendCmd(p, cmd)
	if err != nil {
		return res, err
	}
	if res < 0 {
		return res, fmt.Errorf("Error while executing NotifyBegin command: Error Code was %d", res)
	}
	return res, nil
}

func NotifyClose(p io.ReadWriter, h NotifyHandle) (int32, error) {
	cmd := Cmd{
		ID: NOTIFY_CLOSE,
		P1: uint32(h),
	}

	res, err := sendCmd(p, cmd)
	if err != nil {
		return res, err
	}
	if res < 0 {
		return res, fmt.Errorf("Error while executing NotifyOpen command: Error Code was %d", res)
	}
	return res, nil
}
