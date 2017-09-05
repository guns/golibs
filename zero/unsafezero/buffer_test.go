// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package unsafezero

import (
	"bytes"
	"reflect"
	"testing"
)

func TestClearBuffer(t *testing.T) {
	buf := &bytes.Buffer{}

	// Fill the internal bootstrap array
	for i := 0; i < 128; i++ {
		err := buf.WriteByte(0xff)
		if err != nil {
			t.Fail()
		}
	}

	// As well as the runeBytes array
	if _, err := buf.WriteRune('❤'); err != nil {
		t.Fail()
	}

	if _, err := buf.Read(make([]byte, 8)); err != nil {
		t.Fail()
	}

	ClearBuffer(buf)

	if !reflect.DeepEqual(buf, bytes.NewBuffer([]byte{})) {
		t.Errorf("%#v != %#v", buf, bytes.NewBuffer([]byte{}))
	}
}
