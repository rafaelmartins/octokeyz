package b8

import (
	"log"
	"sync"
	"time"
)

type HandlerFunc func(ev *Event) error

type handler struct {
	fn     HandlerFunc
	button ButtonType

	mtx    sync.Mutex
	events []*Event
}

func (h *handler) execute(ev *Event) {
	if h.fn != nil && ev != nil && ev.Button == h.button {
		switch ev.buttonState {
		case buttonDown:
			go func(ev Event) {
				e := &ev
				e.channel = make(chan time.Time, 1)

				h.mtx.Lock()
				h.events = append(h.events, e)
				h.mtx.Unlock()

				if err := h.fn(e); err != nil {
					log.Print("error: b8: handler: ", err)
				}

				if !e.done {
					<-e.channel
				}
			}(*ev)

		case buttonUp:
			h.mtx.Lock()
			tmp := []*Event{}
			for _, e := range h.events {
				if ev.Button == e.Button {
					e.channel <- ev.Time
				} else {
					tmp = append(tmp, e)
				}
			}
			h.events = tmp
			h.mtx.Unlock()
		}
	}
}
