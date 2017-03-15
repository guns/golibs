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

	m["username"] = "username should not be blank"

	if m.Error() != "validation failed: username should not be blank" {
		t.Errorf("%#v != %#v", m.Error(), "validation failed: username should not be blank")
	}

	m["password"] = "password must be longer than 12 characters"

	n := len("validation failed: username should not be blank, password must be longer than 12 characters")
	if len(m.Error()) != n {
		t.Errorf("%#v != %#v", len(m.Error()), n)
	}

	e := func() error { return m }()
	switch e.(type) {
	case ErrorMap:
	default:
		t.Errorf("%v != %v", reflect.TypeOf(e), reflect.TypeOf(m))
	}
}

func TestThat(t *testing.T) {
	IsPositive := func(n int, key string) Checker {
		return func() (bool, string, string) {
			return n > 0, key, key + " must be positive"
		}
	}

	m := That(IsPositive(0, "x"), IsPositive(1, "y"), IsPositive(-1, "z"))
	emap := ErrorMap{"x": "x must be positive", "z": "z must be positive"}
	if !reflect.DeepEqual(m, emap) {
		t.Errorf("%#v != %#v", m, emap)
	}
}
