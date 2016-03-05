package shell

import (
	"os"
	"testing"
)

func TestExpandUserDir(t *testing.T) {
	u, err := ExpandUserDir("")
	if !(len(u) == 0) {
		t.Errorf("expected: len(u) == 0, actual: %v", len(u))
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	u, err = ExpandUserDir("/")
	if u != "/" {
		t.Errorf("%v != %v", u, "/")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	u, err = ExpandUserDir("~")
	if u != os.ExpandEnv("${HOME}") {
		t.Errorf("%v != %v", u, os.ExpandEnv("${HOME}"))
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	u, err = ExpandUserDir("~/Mail")
	if u != os.ExpandEnv("${HOME}/Mail") {
		t.Errorf("%v != %v", u, os.ExpandEnv("${HOME}/Mail"))
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// FIXME: Not portable
	u, err = ExpandUserDir("~root/Mail")
	if u != os.ExpandEnv("/root/Mail") {
		t.Errorf("%v != %v", u, os.ExpandEnv("/root/Mail"))
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
