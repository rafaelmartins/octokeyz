// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: GPL-2.0

package octokeyz

import (
	"errors"
	"sync"
)

// Modifier is an opaque structure that represent a modifier button.
//
// Modifier buttons allow the implementation of secondary functions for buttons,
// by checking if the modifier button is pressed or not in a button handler callback.
type Modifier struct {
	mtx     sync.Mutex
	pressed bool
}

// Handler is a ButtonHandler implementation for a modifier button. It should be added
// to the button that will be used as modifier.
func (m *Modifier) Handler(b *Button) error {
	if !m.mtx.TryLock() {
		return errors.New("modifier activated by more than one button")
	}
	defer m.mtx.Unlock()

	m.pressed = true
	b.WaitForRelease()
	m.pressed = false

	return nil
}

// Pressed returns true if the modifier button is pressed.
func (m *Modifier) Pressed() bool {
	return m.pressed
}
