package check

import (
	"reflect"
	"testing"
)

func TestErrorMap(t *testing.T) {
	m := ErrorMap{}

	if m.Error() != "validation passed" {
		t.Errorf("%#v != %#v", m.Error(), "validation passed")
	}

	m["username"] = "must not be blank"

	s := "validation failed: username must not be blank"
	if m.Error() != s {
		t.Errorf("%#v != %#v", m.Error(), s)
	}

	m["password"] = "must be longer than 12 characters"

	s = "validation failed: username must not be blank, password must be longer than 12 characters"
	if m.Error() != s {
		t.Errorf("%#v != %#v", m.Error(), s)
	}

	e := func() error { return m }()
	switch e.(type) {
	case ErrorMap:
	default:
		t.Errorf("%v != %v", reflect.TypeOf(e), reflect.TypeOf(m))
	}
}

func TestThat(t *testing.T) {
	IsPositive := func(n int, key string) Fn {
		return func() (bool, string, string) {
			return n > 0, key, "must be positive"
		}
	}

	m := That(IsPositive(0, "x"), IsPositive(1, "y"), IsPositive(-1, "z"))
	emap := ErrorMap{"x": "must be positive", "z": "must be positive"}
	if !reflect.DeepEqual(m, emap) {
		t.Errorf("%#v != %#v", m, emap)
	}
}
