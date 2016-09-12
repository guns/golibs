package zero

import (
	"bytes"
	"errors"
	"sync"
	"sync/atomic"
)

/*

A Cache is a synchronized buffer that is initialized from a constant function
and can be zeroed.

*/
type Cache struct {
	// done is always written atomically, and is either read atomically or
	// read while holding a lock on mutex.
	done   uint32
	mutex  sync.RWMutex
	bytes  []byte
	err    error
	initFn func() ([]byte, error)
}

// NewCache creates a Cache that caches bytes from initFn.
func NewCache(initFn func() ([]byte, error)) *Cache {
	return &Cache{initFn: initFn}
}

// Init idempotently initializes a Cache.
func (cache *Cache) Init() {
	// cf. sync.once.Do()
	cache.mutex.Lock()
	if cache.done == 0 {
		cache.bytes, cache.err = cache.initFn()
		atomic.StoreUint32(&cache.done, 1)
	}
	cache.mutex.Unlock()
}

// WithByteReader calls f with a *bytes.Reader on the cache byte slice. If the
// cache is uninitialized, it will be atomically populated before f is called.
// The immutability of the underlying byte slice is only guaranteed during the
// lifetime of f.
//
// If an error is returned, it should be assumed that f() was never called.
func (cache *Cache) WithByteReader(f func(*bytes.Reader)) error {
	if atomic.LoadUint32(&cache.done) == 0 {
		cache.Init()
	}
	// Assertion:
	//	It is impossible to acquire this RLock without setting
	//	cache.done, cache.bytes, and cache.err
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	if cache.err != nil {
		return cache.err
	}
	f(bytes.NewReader(cache.bytes))
	return nil
}

// Clear locks, zeroes, truncates, and unlocks the cache. The initialization
// flag and error message are set to prevent reuse.
func (cache *Cache) Clear() {
	cache.mutex.Lock()
	atomic.StoreUint32(&cache.done, 1)
	if cache.err == nil {
		cache.err = errors.New("cannot read cleared zero.Cache")
	}
	ClearBytes(cache.bytes)
	cache.bytes = nil
	cache.mutex.Unlock()
}

// Dup returns a new uninitialized cache with the same initFn.
func (cache *Cache) Dup() *Cache {
	return &Cache{initFn: cache.initFn}
}
