package b8

import (
	"fmt"
	"io"
	"time"
	"unsafe"
)

const (
	evKey      = 1
	btnMacro   = 0x290
	numButtons = 8
)

type ButtonType uint8

const (
	BUTTON_1 ButtonType = iota
	BUTTON_2
	BUTTON_3
	BUTTON_4
	BUTTON_5
	BUTTON_6
	BUTTON_7
	BUTTON_8
)

type buttonStateType uint8

const (
	buttonUp buttonStateType = iota
	buttonDown
)

type Event struct {
	Time   time.Time
	Button ButtonType

	buttonState buttonStateType
	hold        time.Duration
	channel     chan time.Time
	done        bool
}

type kernelEvent struct {
	time_ struct {
		sec  int64
		usec int64
	}
	type_ uint16
	code  uint16
	value int32
}

const (
	kernelEventSize = int(unsafe.Sizeof(kernelEvent{}))
)

func newEvents(r io.Reader) ([]*Event, error) {
	buf := make([]byte, 64*kernelEventSize)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	if n%kernelEventSize != 0 {
		return nil, fmt.Errorf("b8: failed to read from evdev")
	}

	rv := []*Event{}

	for i := 0; i < n/kernelEventSize; i++ {
		ev := *(*kernelEvent)(unsafe.Pointer(&buf[i*kernelEventSize]))

		if ev.type_ != evKey {
			continue
		}

		rv = append(rv, &Event{
			Time:        time.Unix(ev.time_.sec, ev.time_.usec),
			Button:      ButtonType(int(ev.code) - btnMacro),
			buttonState: buttonStateType(ev.value),
		})
	}

	return rv, nil
}

func (e *Event) WaitForRelease() time.Duration {
	if e.done {
		return e.hold
	}

	c := <-e.channel
	e.hold = c.Sub(e.Time)
	return e.hold
}
