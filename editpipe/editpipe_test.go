package editpipe

import (
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"testing"
)

func TestEditPipe(t *testing.T) {
	data := []struct {
		input  string
		buflen int
		editFn EditFn
		expect string
	}{
		{input: "\n", expect: "\n"},
		{input: "\x00\x00\x00\n", expect: "\x00\x00\x00\n"},
		{input: "foo bar baz\n", expect: "foo bar baz\n"},
		{input: "foo bar baz", expect: ""},
		{input: "foo bar baz\x15\n", expect: "\n"},
		{input: "foo bar baz \x15\n", expect: "\n"},
		{input: "foo bar baz\x17\n", expect: "foo bar \n"},
		{input: "foo bar baz \x17\n", expect: "foo bar \n"},
		{input: "foo bar\tbaz\t\x01\t\x01 \x17\n", expect: "foo bar\t\n"},
		{input: "foo BAR \x17bar  \bbaz\n", expect: "foo bar baz\n"},
		{input: "\b\n", expect: "\n"},
		{input: "\x15\x17\n", expect: "\n"},
		{input: "\x17\x15\n", expect: "\n"},
		{input: "foo bar baz\b\x7f\b\n", expect: "foo bar \n"},
		{input: "\x04", expect: ""},
		{input: "foo bar baz\x04", expect: "foo bar baz"},
		{input: "foo bar baz\n\x04", expect: "foo bar baz\n"},
		{input: "12345678", buflen: 8, expect: ""},
		{input: "123456789", buflen: 8, expect: "12345678"},
		{input: "1234567890\x15\n", buflen: 8, expect: "12345678\n"},
		{input: "1\n23\n456\n7890\n", buflen: 8, expect: "1\n23\n456\n7890\n"},
		{
			input: "foo bar baz\n",
			editFn: func(i int, b byte) Op {
				if b == 'a' || b == 'o' || b == ' ' {
					return 0
				}
				return BasicLineEdit(i, b)
			},
			expect: "fbrbz\n",
		},
		{
			input: "foo\b\x7f\x17\x15\x04bar\n\x07baz\n\x07",
			editFn: func(i int, b byte) Op {
				if b == '\x07' {
					return Flush | Close
				}
				return Append
			},
			expect: "foo\b\x7f\x17\x15\x04bar\n",
		},
	}

	for _, row := range data {
		r, w := io.Pipe()
		p := New(strings.NewReader(row.input), w, true, row.buflen, row.editFn)
		wg := sync.WaitGroup{}
		var perr error

		wg.Add(1)
		go func() {
			perr = p.ProcessAll()
			wg.Done()
		}()

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		wg.Wait()

		if perr != nil {
			t.Errorf("%#v != %#v", perr, nil)
		}

		if len(p.rbuf) != 1 {
			t.Errorf("%#v != %#v", len(p.rbuf), 1)
		}

		if p.rbuf[0] != 0 {
			t.Errorf("%#v != %#v", p.rbuf[0], 0)
		}

		for i := range p.wbuf {
			if p.wbuf[i] != 0 {
				t.Error("should be clear: %#v", p.wbuf)
				break
			}
		}

		if string(out) != row.expect {
			t.Errorf("%#v != %#v", string(out), row.expect)
		}
	}
}
