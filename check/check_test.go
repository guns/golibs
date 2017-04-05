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

func isPositive(n int, key string) Fn {
	return func() (bool, string, string) {
		return n > 0, key, "must be positive"
	}
}

func isEven(n int, key string) Fn {
	return func() (bool, string, string) {
		return (n&1 == 0), key, "must be even"
	}
}

func TestThat(t *testing.T) {
	data := []struct {
		in, out error
	}{
		{That(isPositive(1, "x")), nil},
		{That(isPositive(1, "x"), isPositive(0, "y")), ErrorMap{"y": "must be positive"}},
		{That(isPositive(1, "x"), isPositive(0, "y"), isPositive(-1, "z")), ErrorMap{"y": "must be positive", "z": "must be positive"}},
	}

	for _, row := range data {
		if !reflect.DeepEqual(row.in, row.out) {
			t.Errorf("%#v != %#v", row.in, row.out)
		}
	}
}

func TestPipe(t *testing.T) {
	data := []struct {
		in, out error
	}{
		{That(Pipe(isPositive(2, "x"), isEven(2, "x"))), nil},
		{That(Pipe(isPositive(-2, "x"), isEven(-2, "x"))), ErrorMap{"x": "must be positive"}},
		{That(Pipe(isPositive(1, "x"), isEven(1, "x"))), ErrorMap{"x": "must be even"}},
		{That(Pipe(isPositive(-1, "x"), isEven(-1, "x"))), ErrorMap{"x": "must be positive"}},
	}

	for _, row := range data {
		if !reflect.DeepEqual(row.in, row.out) {
			t.Errorf("%#v != %#v", row.in, row.out)
		}
	}
}
