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

func isPositive(key string, n int) Fn {
	return func() (bool, string, string) {
		return n > 0, key, "must be positive"
	}
}

func isEven(key string, n int) Fn {
	return func() (bool, string, string) {
		return (n&1 == 0), key, "must be even"
	}
}

func TestThat(t *testing.T) {
	data := []struct {
		in, out error
	}{
		{That(isPositive("x", 1)), nil},
		{That(isPositive("x", 1), isPositive("y", 0)), ErrorMap{"y": "must be positive"}},
		{That(isPositive("x", 1), isPositive("y", 0), isPositive("z", -1)), ErrorMap{"y": "must be positive", "z": "must be positive"}},
	}

	for _, row := range data {
		if !reflect.DeepEqual(row.in, row.out) {
			t.Errorf("%#v != %#v", row.in, row.out)
		}
	}
}

func TestAnd(t *testing.T) {
	data := []struct {
		in  Fn
		out error
	}{
		{And(), nil},
		{And(isPositive("x", 1)), nil},
		{And(isPositive("x", -1)), ErrorMap{"x": "must be positive"}},
		{And(isPositive("x", 2), isEven("x", 2)), nil},
		{And(isPositive("x", -2), isEven("x", -2)), ErrorMap{"x": "must be positive"}},
		{And(isPositive("x", 1), isEven("x", 1)), ErrorMap{"x": "must be even"}},
		{And(isPositive("x", -1), isEven("x", -1)), ErrorMap{"x": "must be positive"}},
	}

	for _, row := range data {
		if !reflect.DeepEqual(That(row.in), row.out) {
			t.Errorf("%#v != %#v", That(row.in), row.out)
		}
	}
}

func TestOr(t *testing.T) {
	data := []struct {
		in  Fn
		out error
	}{
		{Or(), nil},
		{Or(isPositive("x", 1)), nil},
		{Or(isPositive("x", -1)), ErrorMap{"x": "must be positive"}},
		{Or(isPositive("x", 2), isEven("x", 2)), nil},
		{Or(isPositive("x", -2), isEven("x", -2)), nil},
		{Or(isPositive("x", 1), isEven("x", 1)), nil},
		{Or(isPositive("x", -1), isEven("x", -1)), ErrorMap{"x": "must be even"}},
	}

	for _, row := range data {
		if !reflect.DeepEqual(That(row.in), row.out) {
			t.Errorf("%#v != %#v", That(row.in), row.out)
		}
	}
}
