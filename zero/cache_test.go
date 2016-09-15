package zero

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
var testErr = errors.New("testError")

func TestCacheInit(t *testing.T) {
	inits := int32(0)
	cache := NewCache(func() ([]byte, error) {
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
	cache := NewCache(func() ([]byte, error) { return append([]byte{}, testBytes...), nil })
	errors := [1000]error{}
	writers := [1000]bytes.Buffer{}

	wg := sync.WaitGroup{}
	wg.Add(1000)

	start := time.Now()

	for i := 0; i < 1000; i++ {
		i := i
		go func() {
			errors[i] = cache.WithByteReader(func(r *bytes.Reader) {
				r.WriteTo(&writers[i])
				time.Sleep(time.Millisecond)
			})
			wg.Done()
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	if !(elapsed < 10*time.Millisecond) {
		t.Errorf("expected: time.Sub(start) < 10*time.Millisecond, actual: %v", elapsed)
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
	cache := NewCache(func() ([]byte, error) { return nil, testErr })

	didread := false
	err := cache.WithByteReader(func(_ *bytes.Reader) {
		didread = true
	})

	if err != testErr {
		t.Errorf("%v != %v", err, testErr)
	}
	if didread {
		t.Errorf("expected: !didread")
	}
}

const (
	cacheClear = 0
	cacheReset = 1
)

func TestCacheClearAndReset(t *testing.T) {
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
			cacheReadAfterClearError,
		},
		{
			cacheClear,
			1,
			nil,
			nil,
			testErr,
			testErr,
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
			testErr,
			nil,
		},
	}

	for _, row := range data {
		inits := 0
		cache := NewCache(func() ([]byte, error) {
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
