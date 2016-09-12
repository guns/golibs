package shell

import (
	"os/exec"
	"testing"
)

// WARNING: UNSAFE! word argument is not shell escaped!
func unsafeShellSprintf(word string) (string, error) {
	bs, err := exec.Command("sh", "-c", "printf %s "+word).Output()
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func TestExpandUserDir(t *testing.T) {
	data := []string{"", "/", "~", "~/", "~/Mail", "~root", "~root/", "~root/Mail"}

	for _, row := range data {
		actual, err := ExpandUserDir(row)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected, err := unsafeShellSprintf(row)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if actual != expected {
			t.Errorf("%#v != %#v", actual, expected)
		}
	}

	s, err := ExpandUserDir("~nosuchuser/")
	if s != "" {
		t.Errorf("%v != %v", s, "")
	}
	if err == nil {
		t.Errorf("expected err to be an error, but got nil")
	}
}
