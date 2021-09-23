# keymutex
Acquire locks on arbitrary strings by hashing over a fixed set of locks.

KeyMutex uses a [FNV-1a](https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function) 32-bit hash of an input string to select a mutex from a list of locks.  The hash is computed in a way that does not cause allocation by converting string to bytes, as with the standard library.
