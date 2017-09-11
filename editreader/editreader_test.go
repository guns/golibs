// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package editreader

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"unicode"
)

func TestEditReader(t *testing.T) {
	data := []struct {
		input  string
		rlen   int
		wlen   int
		editFn EditFn
		expect string
	}{
		// Simple inputs
		{input: "\n", expect: "\n"},
		{input: "foo bar baz\n", expect: "foo bar baz\n"},
		{input: "lorem ipsum dolor\n\nsit amet\n", expect: "lorem ipsum dolor\n\nsit amet\n"},
		{input: "\x00\x01\n\x02\n", expect: "\x00\x01\n\x02\n"},

		// Unflushed inputs
		{input: "", expect: ""},
		{input: "foo bar baz", expect: ""},
		{input: "foo\nbar baz", expect: "foo\n"},

		// Partial reads
		{input: "1234567890\n1234567890\n", wlen: 3, expect: "1234567890\n1234567890\n"},

		// Overflow
		{input: "1234567890\n", rlen: 3, expect: "1234567890\n"},
		{input: "1234567890", rlen: 3, expect: "123456789"},

		// Erase
		{input: "\b\n", expect: "\n"},
		{input: "12\x7f\b\x7f345\n", expect: "345\n"},
		{input: "1234567\b\x7f\b890\n", expect: "1234890\n"},

		// EraseWord
		{input: "\x17\n", expect: "\n"},
		{input: " \x17\n", expect: "\n"},
		{input: "foo\x17\n", expect: "\n"},
		{input: "foo \t\x17\n", expect: "\n"},
		{input: "foo bar\x17\n", expect: "foo \n"},
		{input: "foo bar \x17\n", expect: "foo \n"},
		{input: "foo bar\x00\x01\x02\x17\n", expect: "foo \n"},

		// Kill
		{input: "\x15\n", expect: "\n"},
		{input: "foo bar baz\x15\n", expect: "\n"},
		{input: "foo\nbar \x15baz\n", expect: "foo\nbaz\n"},

		// Close
		{input: "\x04", expect: ""},
		{input: "foo\x04", expect: "foo"},
		{input: "foo\x04\x04bar\n", expect: "foo"},

		// Custom EditFn
		{
			input: "one 2 three 4 five\n",
			editFn: func(i int, b byte) Op {
				if unicode.IsNumber(rune(b)) {
					return 0
				}
				return BasicLineEdit(i, b)
			},
			expect: "one  three  five\n",
		},
		{
			// This case tests handling of unread data after close
			input: "foo\b\x7f\x17\x15\x04bar\n\x07baz\n\x07",
			editFn: func(i int, b byte) Op {
				if b == 0x07 {
					return Flush | Close
				}
				return Append
			},
			expect: "foo\b\x7f\x17\x15\x04bar\n",
		},
	}

	for _, row := range data {
		for _, writeTo := range []bool{false, true} {
			if row.wlen == 0 {
				row.wlen = 80
			}

			r := New(strings.NewReader(row.input), row.rlen, true, row.editFn)
			out := make([]byte, 0, 256)
			buf := make([]byte, row.wlen)
			var n int
			var err error

			out, buf, err = readall(r, out, buf, writeTo)
			if err != io.EOF {
				t.Errorf("unexpected non-EOF error: %#v", err)
			}

			// Read once more; should get (0, EOF)
			n, err = r.Read(buf)
			if n != 0 {
				t.Errorf("%#v != %#v", n, 0)
			}
			if err != io.EOF {
				t.Errorf("%#v != %#v", err, io.EOF)
			}

			if len(r.rbuf) != 1 {
				t.Errorf("%#v != %#v", len(r.rbuf), 1)
			}

			if r.rbuf[0] != 0 {
				t.Errorf("should be clear: %#v", r.rbuf[0])
			}

			for i := range r.buf {
				if r.buf[i] != 0 {
					t.Errorf("should be clear: %#v", string(r.buf))
					break
				}
			}

			if string(out) != row.expect {
				t.Errorf("%#v != %#v", string(out), row.expect)
			}
		}
	}
}

func readall(r *T, out, buf []byte, writeTo bool) ([]byte, []byte, error) {
	var n int
	var err error

	if writeTo {
		bbuf := bytes.Buffer{}
		_, err = r.WriteTo(&bbuf)
		out = append(out, bbuf.Bytes()...)
		if err == nil {
			err = io.EOF
		}
	} else {
		for {
			n, err = r.Read(buf)
			if n > 0 {
				out = append(out, buf[:n]...)
			}
			if err != nil {
				break
			}
		}
	}

	return out, buf, err
}
