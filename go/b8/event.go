package b8

import (
	"fmt"
	"io"
	"time"
	"unsafe"
)

const (
	evKey    = 1
	btnMacro = 0x290
)

type event struct {
	etime   time.Time
	button  ButtonID
	pressed bool
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

func newEvents(r io.Reader) ([]*event, error) {
	buf := make([]byte, 64*kernelEventSize)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	if n%kernelEventSize != 0 {
		return nil, fmt.Errorf("b8: failed to read from evdev")
	}

	rv := []*event{}

	for i := 0; i < n/kernelEventSize; i++ {
		ev := *(*kernelEvent)(unsafe.Pointer(&buf[i*kernelEventSize]))

		if ev.type_ != evKey {
			continue
		}

		rv = append(rv, &event{
			etime:   time.Unix(ev.time_.sec, ev.time_.usec*1000),
			button:  ButtonID(int(ev.code) - btnMacro),
			pressed: ev.value > 0,
		})
	}

	return rv, nil
}
