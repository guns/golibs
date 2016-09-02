package zero

import (
	"bytes"
	"errors"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

func sliceAddr(s unsafe.Pointer) uintptr {
	return ((*reflect.SliceHeader)(s)).Data
}

func TestCacheClear(t *testing.T) {
	src := append(make([]byte, 0, 512), "01234567"...)
	cache := NewCache(func() ([]byte, error) {
		return src, nil
	})

	err := cache.WithByteReader(func(r *bytes.Reader) {
		buf := make([]byte, len(src))
		n, e := r.Read(buf)
		if !reflect.DeepEqual(buf[:n], src) {
			t.Errorf("%v != %v", buf[:n], src)
		}
		if e != nil {
			t.Errorf("unexpected error: %v", e)
		}
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	addr := sliceAddr(unsafe.Pointer(&cache.bytes))
	cache.Clear()
	if !(addr != sliceAddr(unsafe.Pointer(&cache.bytes))) {
		t.Errorf("expected: addr != sliceAddr(unsafe.Pointer(&cache.bytes)), %v == %v; cache should not refer to same memory after reset", addr, sliceAddr(unsafe.Pointer(&cache.bytes)))
	}

	if cache.bytes != nil {
		t.Errorf("unexpected non-nil value: %v", cache.bytes)
	}
	if !reflect.DeepEqual(src, make([]byte, 8)) {
		t.Errorf("%v != %v", src, make([]byte, 8))
	}
}

func TestCacheIdempotencyAndIdentity(t *testing.T) {
	inits := 0

	cache := NewCache(func() ([]byte, error) {
		inits++
		return []byte("01234567")[:inits], nil
	})

	buf := make([]byte, 10)

	for i := 0; i < 10; i++ {
		err := cache.WithByteReader(func(r *bytes.Reader) {
			n, e := r.Read(buf)
			if !reflect.DeepEqual(buf[:n], []byte{'0'}) {
				t.Errorf("%v != %v; reading cache should be idempotent", buf[:n], []byte{'0'})
			}
			if e != nil {
				t.Errorf("unexpected error: %v", e)
			}
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	addr := sliceAddr(unsafe.Pointer(&cache.bytes))

	if inits != 1 {
		t.Errorf("%v != %v; cache init fn should only be fired once", inits, 1)
	}

	buf = make([]byte, 10)
	cache.Clear()

	for i := 0; i < 3; i++ {
		err := cache.WithByteReader(func(r *bytes.Reader) {
			_, e := r.Read(buf)
			if e != nil {
				t.Errorf("unexpected error: %v", e)
			}
		})
		if err == nil {
			t.Errorf("expected err to be an error, but got nil")
		}
	}

	if inits != 1 {
		t.Errorf("%v != %v; cache init fn should still only be fired once", inits, 1)
	}
	if !(addr != sliceAddr(unsafe.Pointer(&cache.bytes))) {
		t.Errorf("expected: addr != sliceAddr(unsafe.Pointer(&cache.bytes)), actual: %v == %v; cache should not refer to same memory after reset", addr, sliceAddr(unsafe.Pointer(&cache.bytes)))
	}
}

func TestCacheConcurrency(t *testing.T) {
	inits := 0
	src := []byte("README")

	cache := NewCache(func() ([]byte, error) {
		inits++
		return append([]byte{}, src...), nil
	})

	var complete int32
	vlen := 256
	vs := make([][]byte, vlen)
	ch := make(chan int, vlen)
	defer close(ch)

	for i := range vs {
		i := i // Capture the value of i
		go func() {
			err := cache.WithByteReader(func(r *bytes.Reader) {
				ch <- i
				time.Sleep(time.Duration(i) * time.Microsecond)
				vs[i] = make([]byte, len(src))
				_, e := r.Read(vs[i])
				atomic.AddInt32(&complete, 1)
				if e != nil {
					t.Errorf("unexpected error: %v", e)
				}
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	gocount := 0
	for range ch {
		gocount++
		if gocount == vlen {
			n := int(atomic.LoadInt32(&complete))
			if !(n < vlen) {
				t.Errorf("expected: n < vlen; actual: %v >= %v; should still have some running goroutines when we reset the Cache", n, vlen)
			}
			cache.Clear()
			break
		}
	}
	if int(complete) != vlen {
		t.Errorf("%v != %v; reset should have waited until all goroutines exited", int(complete), vlen)
	}
	if inits != 1 {
		t.Errorf("%v != %v; cache should only be initialized once", inits, 1)
	}
	for _, v := range vs {
		if !reflect.DeepEqual(v, src) {
			t.Errorf("%v != %v", v, src)
		}
	}
}

func TestCacheErrors(t *testing.T) {
	inits, n := 0, 256

	cache := NewCache(func() ([]byte, error) {
		inits++
		return nil, errors.New("cache init error")
	})

	ch := make(chan error, n)
	defer close(ch)

	for i := 0; i < n; i++ {
		go func() {
			err := cache.WithByteReader(func(r *bytes.Reader) {
				_, e := r.Read([]byte{})
				if e != nil {
					t.Errorf("unexpected error: %v", e)
				}
			})
			ch <- err
		}()
		// Unclosed: ch
	}

	for i := 0; i < n; i++ {
		if <-ch == nil {
			t.Errorf("expected <-ch to be an error, but got nil")
		}
	}

	if inits != 1 {
		t.Errorf("%v != %v", inits, 1)
	}

	if cache.err == nil {
		t.Errorf("expected cache.err to be an error, but got nil")
	}
	cache.Clear()
	if cache.err == nil {
		t.Errorf("Clear should not clear errors")
	}
}

func TestCacheClone(t *testing.T) {
	inits := 0

	cache := NewCache(func() ([]byte, error) {
		inits++
		return []byte{'a'}, errors.New("cache init error")
	})

	err := cache.WithByteReader(func(r *bytes.Reader) {
		_, e := r.Read([]byte{})
		if e != nil {
			t.Errorf("unexpected error: %v", e)
		}
	})
	if err == nil {
		t.Errorf("expected err to be an error, but got nil")
	}
	if err != cache.err {
		t.Errorf("%v != %v", err, cache.err)
	}
	if !reflect.DeepEqual(cache.bytes, []byte{'a'}) {
		t.Errorf("%v != %v", cache.bytes, []byte{'a'})
	}

	cache.Clear()
	*cache = *cache.Dup()
	if !reflect.DeepEqual(cache.bytes, []byte(nil)) {
		t.Errorf("%v != %v", cache.bytes, []byte(nil))
	}
	if cache.err != nil {
		t.Errorf("unexpected error: %v", cache.err)
	}

	err = cache.WithByteReader(func(r *bytes.Reader) {
		_, e := r.Read([]byte{})
		if e != nil {
			t.Errorf("unexpected error: %v", e)
		}
	})
	if err == nil {
		t.Errorf("expected err to be an error, but got nil")
	}
	if inits != 2 {
		t.Errorf("%v != %v", inits, 2)
	}
}
