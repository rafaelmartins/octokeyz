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

// String returns a string representation of a button identifier
func (b ButtonID) String() string {
	return fmt.Sprintf("BUTTON_%d", b)
}

// A b8 USB keypad contains 8 buttons
const (
	BUTTON_1 ButtonID = iota + 1
	BUTTON_2
	BUTTON_3
	BUTTON_4
	BUTTON_5
	BUTTON_6
	BUTTON_7
	BUTTON_8
)

// ButtonHandlerError represents the error returned by a button handler
// including the button identifier.
type ButtonHandlerError struct {
	ButtonID ButtonID
	Err      error
}

// Error returns a string representation of a button handler error.
func (b ButtonHandlerError) Error() string {
	return fmt.Sprintf("b8: %s: %s", b.ButtonID, b.Err)
}

// Button is an opaque structure that represents a b8 USB keypad button.
type Button struct {
	mtx      sync.Mutex
	id       ButtonID
	channel  chan bool
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
	return b.id.String()
}

func (b *Button) addHandler(h ButtonHandler) {
	if h == nil {
		return
	}

	b.mtx.Lock()
	b.handlers = append(b.handlers, h)
	b.mtx.Unlock()
}

func (b *Button) press(t time.Time, errCh chan error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.channel = make(chan bool)
	b.pressed = t
	b.released = time.Time{}
	b.duration = 0

	for _, h := range b.handlers {
		go func(btn *Button, hnd ButtonHandler) {
			if err := hnd(btn); err != nil {
				e := ButtonHandlerError{
					ButtonID: btn.id,
					Err:      err,
				}

				if errCh != nil {
					select {
					case errCh <- e:
					default:
					}
				} else {
					log.Printf("error: %s", e)
				}
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
	close(b.channel)
}

// WaitForRelease blocks a button handler until the button that was
// pressed to activated the handler is released. It returns the duration
// of the button press.
//
// This function should not be called outside a ButtonHandler. It may
// cause undefined behavior.
func (b *Button) WaitForRelease() time.Duration {
	<-b.channel
	return b.duration
}
