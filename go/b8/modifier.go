package b8

import (
	"errors"
	"sync"
)

type Modifier struct {
	mtx     sync.Mutex
	pressed bool
}

func (m *Modifier) Handler(b *Button) error {
	if !m.mtx.TryLock() {
		return errors.New("b8: modifier activated by more than one button")
	}
	defer m.mtx.Unlock()

	m.pressed = true
	b.WaitForRelease()
	m.pressed = false

	return nil
}

func (m *Modifier) Pressed() bool {
	return m.pressed
}
