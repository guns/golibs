// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

package errutil

import (
	"errors"
	"reflect"
	"testing"
)

func TestJoin(t *testing.T) {
	data := []struct {
		in  []error
		out error
	}{
		{[]error{}, nil},
		{[]error{nil}, nil},
		{[]error{nil, nil}, nil},
		{[]error{errors.New("a"), nil}, errors.New("a")},
		{[]error{nil, errors.New("b")}, errors.New("b")},
		{[]error{errors.New("a"), nil, errors.New("c")}, errors.New("a; c")},
		{[]error{errors.New("a"), errors.New("b"), errors.New("c")}, errors.New("a; b; c")},
	}

	for _, row := range data {
		err := Join("; ", row.in...)
		if !reflect.DeepEqual(err, row.out) {
			t.Errorf("%#v != %#v", err, row.out)
		}
	}
}

func TestFirst(t *testing.T) {
	err1 := errors.New("1")
	err2 := errors.New("2")
	err3 := errors.New("3")

	data := []struct {
		in  []error
		out error
	}{
		{[]error{}, nil},
		{[]error{nil}, nil},
		{[]error{nil, nil}, nil},
		{[]error{err3, err2, err1}, err3},
		{[]error{nil, err2, err1}, err2},
		{[]error{nil, nil, err1}, err1},
	}

	for _, row := range data {
		err := First(row.in...)
		if err != row.out {
			t.Errorf("%#v != %#v", err, row.out)
		}
	}
}
