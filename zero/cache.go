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

// Init initializes a Cache.
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
// If an error is returned, it should be assumed that f() was never called.
func (cache *Cache) WithByteReader(f func(*bytes.Reader)) error {
	if atomic.LoadUint32(&cache.done) == 0 {
		cache.Init()
	}
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
		cache.err = errors.New("cannot read cleared Cache")
	}
	ClearBytes(cache.bytes)
	cache.bytes = nil
	cache.mutex.Unlock()
}

// Dup returns a new uninitialized cache with the same initFn.
func (cache *Cache) Dup() *Cache {
	return &Cache{initFn: cache.initFn}
}
