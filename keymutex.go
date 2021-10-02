package keymutex

import (
	"sync"
)

const (
	// FNV-1a
	offset32 = uint32(2166136261)
	prime32  = uint32(16777619)

	defaultLocks = 31
)

type KeyMutex struct {
	mutexes []sync.Mutex
}

// New returns a new instance of KeyMutex which hashes arbitrary keys to
// a fixed set of locks, specified by n.  Use the default value if n <= 0.
func New(n int) *KeyMutex {
	if n <= 0 {
		n = defaultLocks
	}
	return &KeyMutex{
		mutexes: make([]sync.Mutex, n),
	}
}

// Lock acquires the lock for the specified key.
func (km *KeyMutex) Lock(k string) {
	km.mutexes[hashString32(k)%uint32(len(km.mutexes))].Lock()
}

// Unlock releases the lock for the specified key.
func (km *KeyMutex) Unlock(k string) {
	km.mutexes[hashString32(k)%uint32(len(km.mutexes))].Unlock()
}

func (km *KeyMutex) LockBytes(k []byte) {
	km.mutexes[hashBytes32(k)%uint32(len(km.mutexes))].Lock()
}

func (km *KeyMutex) UnlockBytes(k []byte) {
	km.mutexes[hashBytes32(k)%uint32(len(km.mutexes))].Unlock()
}

// hashString32 performs a fnv-1a hash on the given string, withoug having to
// allocate to convert the string to bytes.
func hashString32(s string) uint32 {
	h := offset32
	for i := 0; i < len(s); i++ {
		h = (h ^ uint32(s[i])) * prime32
	}
	return h
}

// hashBytes32 performs a fnv-1a hash on the given byte slice
func hashBytes32(b []byte) uint32 {
	h := offset32
	for _, c := range b {
		h = (h ^ uint32(c)) * prime32
	}
	return h
}
