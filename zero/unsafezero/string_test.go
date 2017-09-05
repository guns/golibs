// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package unsafezero

import (
	"math/rand"
	"testing"
)

func TestClearString(t *testing.T) {
	// Only dynamic strings can be mutated
	bs := make([]byte, 0x1000)
	for i := range bs {
		bs[i] = byte(rand.Uint32() & 0x7f)
	}

	str := string(bs)
	copy := str
	z := string(make([]byte, 0x1000))

	ClearString(str)

	if str != z {
		t.Errorf("%v != %v", str, z)
	}
	if copy != z {
		t.Errorf("%v != %v", copy, z)
	}
}
