// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package zero

import (
	"bytes"
	"errors"
	"sync"
	"sync/atomic"
)

/*

A Cache is a synchronized read-only buffer that is initialized from a constant
function and can be zeroed and reset.

*/
type Cache struct {
	done   uint32
	mutex  sync.RWMutex
	bytes  []byte
	err    error
	initFn func() ([]byte, error)
}

// NewCache returns an object that caches bytes from initFn. If initFn returns
// an error, Init will record the error and subsequent reads of the underlying
// buffer will fail.
func NewCache(initFn func() ([]byte, error)) *Cache {
	return &Cache{initFn: initFn}
}

// Init initializes a Cache. This method is synchronized and idempotent.
func (cache *Cache) Init() {
	// cf. sync.once.Do()
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if atomic.CompareAndSwapUint32(&cache.done, 0, 1) {
		cache.bytes, cache.err = cache.initFn()
	}
}

// WithByteReader calls f with a *bytes.Reader on the data returned from the
// initialization function. If the cache is uninitialized, Init is called
// before f is executed. If there was an error during initialization, that
// error is returned without executing f.
//
// This method is synchronized with a read-lock and may be called concurrently
// from multiple goroutines. The immutability of the underlying buffer is only
// guaranteed during the lifetime of f.
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

var errReadAfterClear = errors.New("cannot read cleared zero.Cache")

// Clear zeroes and truncates the underlying buffer without resetting it.
// Cleared Caches cannot be read. This method is synchronized.
func (cache *Cache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	atomic.StoreUint32(&cache.done, 1)
	ClearBytes(cache.bytes)
	cache.bytes = cache.bytes[:0]
	if cache.err == nil {
		cache.err = errReadAfterClear
	}
}

// Reset clears and truncates the underlying buffer and forgets initialization
// errors. A Cache that has been Reset can be re-initialized with Init. This
// method is synchronized
func (cache *Cache) Reset() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	atomic.StoreUint32(&cache.done, 0)
	ClearBytes(cache.bytes)
	cache.bytes = cache.bytes[:0]
	cache.err = nil
}
