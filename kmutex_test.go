package keymutex

import (
	"hash/fnv"
	"testing"
	"time"
)

func TestFnv1aString32(t *testing.T) {
	strs := []string{
		"jmZTUnAo2c90m;sdnI",
		"qqsP/<uZpXO+`^391",
		"=M9fw7f.qkCW",
		"wquo5HIPoO",
		"eM~`&#5sTnFsZ8",
		"gIVgLj-)5bic",
		"#M1s7ccv7(m*q",
		"5%KFAXeLDI7",
		"gPi6g]@c&^Ih3U",
		"wtwZ?{TYEC_Rl",
	}
	for _, k := range strs {
		h := fnv.New32a()
		h.Write([]byte(k))
		sum1 := h.Sum32()
		sum2 := hashString32(k)
		if sum1 != sum2 {
			t.Errorf("invalid hash, expected %x but got %x", sum1, sum2)
		}
	}
}

func TestKeyMutex(t *testing.T) {
	km := New(5)
	km.Lock("hello")

	doneHello := make(chan struct{})
	go func() {
		km.Lock("hello")
		close(doneHello)
	}()

	doneWorld := make(chan struct{})
	go func() {
		km.Lock("world")
		close(doneWorld)
	}()

	to := time.After(time.Second)
	select {
	case <-doneHello:
		t.Fatal("should not locked")
	case <-doneWorld:
	case <-to:
		t.Fatal("world was not unlocked")
	}

	km.Unlock("world")

	select {
	case <-doneHello:
		t.Fatal("should not locked")
	case <-to:
	}

	km.Unlock("hello")

	to = time.After(time.Second)
	select {
	case <-doneHello:
	case <-to:
		t.Fatal("hello not unlocked")
	}
}
