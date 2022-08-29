package keymutex

import (
	"sync"
)

type KeyRWMutex struct {
	mutexes []sync.RWMutex
}

// NewRW returns a new instance of KeyRWMutex which hashes arbitrary keys to a
// fixed set of RWMutexes, specified by n. Use the default value if n <= 0.
func NewRW(n int) *KeyRWMutex {
	if n <= 0 {
		n = defaultLocks
	}
	return &KeyRWMutex{
		mutexes: make([]sync.RWMutex, n),
	}
}

// Lock acquires the lock for the specified key.
func (km *KeyRWMutex) Lock(k string) {
	km.mutexes[hashString32(k)%uint32(len(km.mutexes))].Lock()
}

// RLock acquires the shared lock for the specified key.
func (km *KeyRWMutex) RLock(k string) {
	km.mutexes[hashString32(k)%uint32(len(km.mutexes))].RLock()
}

// Unlock releases the lock for the specified key.
func (km *KeyRWMutex) Unlock(k string) {
	km.mutexes[hashString32(k)%uint32(len(km.mutexes))].Unlock()
}

// RUnlock releases the shared lock for the specified key.
func (km *KeyRWMutex) RUnlock(k string) {
	km.mutexes[hashString32(k)%uint32(len(km.mutexes))].RUnlock()
}

func (km *KeyRWMutex) LockBytes(k []byte) {
	km.mutexes[hashBytes32(k)%uint32(len(km.mutexes))].Lock()
}

func (km *KeyRWMutex) RLockBytes(k []byte) {
	km.mutexes[hashBytes32(k)%uint32(len(km.mutexes))].RLock()
}

func (km *KeyRWMutex) UnlockBytes(k []byte) {
	km.mutexes[hashBytes32(k)%uint32(len(km.mutexes))].Unlock()
}

func (km *KeyRWMutex) RUnlockBytes(k []byte) {
	km.mutexes[hashBytes32(k)%uint32(len(km.mutexes))].RUnlock()
}
