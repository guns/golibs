// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package zerocache

import (
	"bytes"
	"errors"
	"math/rand"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var testBytes = []byte("testBytes")
var errExpected = errors.New("errExpected")

func TestCacheInit(t *testing.T) {
	inits := int32(0)
	cache := New(func() ([]byte, error) {
		atomic.AddInt32(&inits, 1)
		time.Sleep(time.Duration(rand.Uint32()&0xff) * time.Nanosecond)
		return append([]byte{}, testBytes...), nil
	})

	wg := sync.WaitGroup{}
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			cache.Init()
			wg.Done()
		}()
	}

	wg.Wait()

	if inits != 1 {
		t.Errorf("%v != %v", inits, 1)
	}
	if !reflect.DeepEqual(cache.bytes, testBytes) {
		t.Errorf("%v != %v", cache.bytes, testBytes)
	}
	if cache.err != nil {
		t.Errorf("unexpected error: %v", cache.err)
	}
}

func TestCacheWithByteReaderCallback(t *testing.T) {
	cache := New(func() ([]byte, error) { return append([]byte{}, testBytes...), nil })
	errors := [1000]error{}
	writers := [1000]bytes.Buffer{}

	wg := sync.WaitGroup{}
	wg.Add(1000)

	start := time.Now()

	for i := 0; i < 1000; i++ {
		i := i
		go func() {
			errors[i] = cache.WithByteReader(func(r *bytes.Reader) {
				_, _ = r.WriteTo(&writers[i]) // errcheck: the writers are checked later
				time.Sleep(time.Millisecond)
			})
			wg.Done()
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	if !(elapsed < 100*time.Millisecond) {
		t.Errorf("expected: time.Sub(start) < 100*time.Millisecond, actual: %v", elapsed)
	}
	for i := range writers {
		if !reflect.DeepEqual(writers[i].Bytes(), testBytes) {
			t.Errorf("%v != %v", writers[i].Bytes(), testBytes)
			break
		}
		if errors[i] != nil {
			t.Errorf("unexpected error: %v", errors[i])
			break
		}
	}
}

func TestCacheWithByteReaderError(t *testing.T) {
	cache := New(func() ([]byte, error) { return nil, errExpected })

	didread := false
	err := cache.WithByteReader(func(_ *bytes.Reader) {
		didread = true
	})

	if err != errExpected {
		t.Errorf("%v != %v", err, errExpected)
	}
	if didread {
		t.Errorf("expected: !didread")
	}
}

func TestCacheClearAndReset(t *testing.T) {
	const cacheClear = 0
	const cacheReset = 1

	data := []struct {
		method, inits      int
		data, expectedData []byte
		err, expectedErr   error
	}{
		{
			cacheClear,
			1,
			append([]byte{}, testBytes...),
			make([]byte, len(testBytes)),
			nil,
			errReadAfterClear,
		},
		{
			cacheClear,
			1,
			nil,
			nil,
			errExpected,
			errExpected,
		},
		{
			cacheReset,
			2,
			append([]byte{}, testBytes...),
			make([]byte, len(testBytes)),
			nil,
			nil,
		},
		{
			cacheReset,
			2,
			nil,
			nil,
			errExpected,
			nil,
		},
	}

	for _, row := range data {
		inits := 0
		cache := New(func() ([]byte, error) {
			inits++
			return row.data, row.err
		})

		cache.Init()

		switch row.method {
		case cacheClear:
			cache.Clear()
		case cacheReset:
			cache.Reset()
		}

		if !reflect.DeepEqual(cache.bytes[:len(row.expectedData)], row.expectedData) {
			t.Errorf("%v != %v", cache.bytes[:len(row.expectedData)], row.expectedData)
		}
		if cache.err != row.expectedErr {
			t.Errorf("%v != %v", cache.err, row.expectedErr)
		}

		cache.Init()

		if inits != row.inits {
			t.Errorf("%v != %v", inits, row.inits)
		}
	}
}
