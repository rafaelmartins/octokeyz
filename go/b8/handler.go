package b8

import (
	"log"
)

type HandlerFunc func(ev *Event) error

type handler struct {
	fn          HandlerFunc
	button      ButtonType
	buttonState ButtonStateType
}

func (h *handler) execute(ev *Event) {
	if h.fn != nil && ev != nil && ev.Button == h.button && ev.ButtonState == h.buttonState {
		go func(ev *Event) {
			if err := h.fn(ev); err != nil {
				log.Print("handler error: ", err)
			}
		}(ev)
	}
}
