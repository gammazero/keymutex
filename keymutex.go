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

// hashString32 performs a fnv-1a hash on the given string, withoug having to
// allocate to convert the string to bytes.
func hashString32(s string) uint32 {
	h := offset32
	for len(s) >= 8 {
		h = (h ^ uint32(s[0])) * prime32
		h = (h ^ uint32(s[1])) * prime32
		h = (h ^ uint32(s[2])) * prime32
		h = (h ^ uint32(s[3])) * prime32
		h = (h ^ uint32(s[4])) * prime32
		h = (h ^ uint32(s[5])) * prime32
		h = (h ^ uint32(s[6])) * prime32
		h = (h ^ uint32(s[7])) * prime32
		s = s[8:]
	}

	if len(s) >= 4 {
		h = (h ^ uint32(s[0])) * prime32
		h = (h ^ uint32(s[1])) * prime32
		h = (h ^ uint32(s[2])) * prime32
		h = (h ^ uint32(s[3])) * prime32
		s = s[4:]
	}

	if len(s) >= 2 {
		h = (h ^ uint32(s[0])) * prime32
		h = (h ^ uint32(s[1])) * prime32
		s = s[2:]
	}

	if len(s) > 0 {
		h = (h ^ uint32(s[0])) * prime32
	}

	return h
}
