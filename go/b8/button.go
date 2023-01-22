// Copyright 2022-2023 Rafael G.Martins. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package b8

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// ButtonHandler is a function prototype that helps defining a callback
// function to handle button events.
type ButtonHandler func(b *Button) error

// ButtonID represents the identifier of a button.
type ButtonID uint8

// A b8 USB keypad contains 8 buttons
const (
	BUTTON_1 ButtonID = iota
	BUTTON_2
	BUTTON_3
	BUTTON_4
	BUTTON_5
	BUTTON_6
	BUTTON_7
	BUTTON_8
)

// Button is an opaque structure that represents a b8 USB keypad button.
type Button struct {
	mtx      sync.Mutex
	id       ButtonID
	channel  chan time.Duration
	pressed  time.Time
	released time.Time
	duration time.Duration
	handlers []ButtonHandler
}

func newButtons() map[ButtonID]*Button {
	return map[ButtonID]*Button{
		BUTTON_1: {id: BUTTON_1},
		BUTTON_2: {id: BUTTON_2},
		BUTTON_3: {id: BUTTON_3},
		BUTTON_4: {id: BUTTON_4},
		BUTTON_5: {id: BUTTON_5},
		BUTTON_6: {id: BUTTON_6},
		BUTTON_7: {id: BUTTON_7},
		BUTTON_8: {id: BUTTON_8},
	}
}

// String returns a string representation of a button
func (b *Button) String() string {
	return fmt.Sprintf("BUTTON_%d", b.id+1)
}

func (b *Button) addHandler(h ButtonHandler) {
	if h == nil {
		return
	}

	b.mtx.Lock()
	b.handlers = append(b.handlers, h)
	b.mtx.Unlock()
}

func (b *Button) press(t time.Time) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	// currently pressed
	if !b.pressed.IsZero() && b.channel != nil {
		// best effort, just try to unlock any pending goroutine
		for range b.handlers {
			select {
			case b.channel <- 0:
			default:
			}
		}
	}

	b.channel = make(chan time.Duration, 1)
	b.pressed = t
	b.released = time.Time{}
	b.duration = 0

	for _, h := range b.handlers {
		go func(btn *Button, hnd ButtonHandler) {
			if err := hnd(btn); err != nil {
				log.Printf("error: b8: %s handler: %s", b, err)
			}
		}(b, h)
	}
}

func (b *Button) release(t time.Time) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	// currently released
	if !b.released.IsZero() {
		return
	}

	b.released = t
	b.duration = b.released.Sub(b.pressed)
	b.pressed = time.Time{}

	for range b.handlers {
		select {
		case b.channel <- b.duration:
		default:
		}
	}
}

// WaitForRelease blocks a button handler until the button that was
// pressed to activated the handler is released. It returns the duration
// of the button press.
//
// This function should not be called outside a ButtonHandler, as it may
// result in undefined behavior.
func (b *Button) WaitForRelease() time.Duration {
	if b.duration != 0 {
		return b.duration
	}

	return <-b.channel
}
