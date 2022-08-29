package keymutex

import (
	"testing"
	"time"
)

func TestKeyRWMutex(t *testing.T) {
	km := NewRW(0)
	km.RLock("hello")

	doneRHello := make(chan struct{})
	go func() {
		km.RLock("hello")
		close(doneRHello)
	}()

	to := time.After(time.Second)
	select {
	case <-doneRHello:
	case <-to:
		t.Fatal("could not RLock hello was not unlocked")
	}

	doneHello := make(chan struct{})
	go func() {
		km.Lock("hello")
		close(doneHello)
	}()

	select {
	case <-doneHello:
		t.Fatal("hello unlocked")
	case <-to:
	}

	km.RUnlock("hello")
	km.RUnlock("hello")

	to = time.After(time.Second)
	select {
	case <-doneHello:
	case <-to:
		t.Fatal("could not Lock hello was not unlocked")
	}

	km.Unlock("hello")
}
