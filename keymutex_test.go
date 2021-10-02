package keymutex

import (
	"hash/fnv"
	"testing"
	"time"
)

const benchKey = "dsaflkjasdflkjasflkjlaksjfljasldflakjsfdkaljsdfaklsjdflaskdjflaksdflkasdjflasdkjflasdkjfasdlkfjasdlkfjlkjasdlkfjasdlkfjasdlkjflasdkjflkasdjflkasdjflkjasdlfkjasdlkjflkasdjflkjasdlkflkasdjklcvjnasklfjklasdjhello"

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
		sum3 := hashBytes32([]byte(k))
		if sum1 != sum3 {
			t.Errorf("invalid hash, expected %x but got %x", sum1, sum3)
		}
	}
}

func TestKeyMutex(t *testing.T) {
	km := New(0)
	km.Lock("hello")

	bytesKey := []byte("somebytes")
	km.LockBytes(bytesKey)

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

	doneBytes := make(chan struct{})
	go func() {
		km.LockBytes(bytesKey)
		close(doneBytes)
	}()

	to := time.After(time.Second)
	select {
	case <-doneHello:
		t.Fatal("hello unlocked")
	case <-doneWorld:
	case <-to:
		t.Fatal("world was not unlocked")
	}

	km.Unlock("world")

	select {
	case <-doneHello:
		t.Fatal("hello unlocked")
	case <-doneBytes:
		t.Fatal("bytes unlocked")
	case <-to:
	}

	km.Unlock("hello")
	km.UnlockBytes(bytesKey)

	to = time.After(time.Second)
	select {
	case <-doneHello:
	case <-to:
		t.Fatal("hello still locked")
	}

	to = time.After(time.Second)
	select {
	case <-doneBytes:
	case <-to:
		t.Fatal("bytes still locked")
	}
}

func BenchmarkKeyMutexString(b *testing.B) {
	km := New(1)

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		km.Lock(benchKey)
		km.Unlock(benchKey)
	}
}

func BenchmarkKeyMutexBytes(b *testing.B) {
	km := New(1)
	key := []byte(benchKey)

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		km.LockBytes(key)
		km.UnlockBytes(key)
	}
}
