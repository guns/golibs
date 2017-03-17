package check

import (
	"reflect"
	"testing"
)

func TestErrorMapError(t *testing.T) {
	data := []struct {
		in  error
		out string
	}{
		{ErrorMap(nil), "validation passed"},
		{ErrorMap{}, "validation passed"},
		{ErrorMap{"alice": "must be alert"}, "validation failed: alice must be alert"},
		{ErrorMap{"alice": "must be alert", "bob": "must be bored"}, "validation failed: alice must be alert, bob must be bored"},
	}

	for _, row := range data {
		if len(row.in.Error()) != len(row.out) {
			t.Errorf("%#v != %#v", row.in.Error(), row.out)
		}
	}
}

func TestThat(t *testing.T) {
	IsPositive := func(n int, key string) Fn {
		return func() (bool, string, string) {
			return n > 0, key, "must be positive"
		}
	}

	data := []struct {
		in  error
		out error
	}{
		{That(IsPositive(1, "x")), nil},
		{That(IsPositive(1, "x"), IsPositive(0, "y")), ErrorMap{"y": "must be positive"}},
		{That(IsPositive(1, "x"), IsPositive(0, "y"), IsPositive(-1, "z")), ErrorMap{"y": "must be positive", "z": "must be positive"}},
	}

	for _, row := range data {
		if row.out == nil {
			if row.in != nil {
				t.Errorf("unexpected non-nil value: %#v", row.in)
			}
		} else {
			if !reflect.DeepEqual(row.in, row.out) {
				t.Errorf("%#v != %#v", row.in, row.out)
			}
		}
	}
}
