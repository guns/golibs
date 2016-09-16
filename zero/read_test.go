package zero

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadAll(t *testing.T) {
	data := []struct {
		initlen, initcap int
		init             string
		newlen, newcap   int
		readersrc        string
		err              error
	}{
		{
			17, 17, "Sally sells cshs ",
			33, 1024, "by the seashore.",
			nil,
		},
		{
			28, 64, "Lorem ipsum dolor sit amet, ",
			56, 64, "consectetur adipisicing elit",
			nil,
		},
		{
			0, 0, "",
			893, 1024, "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			nil,
		},
	}

	for _, row := range data {
		init := make([]byte, row.initlen, row.initcap)
		copy(init, row.init)

		initExpected := make([]byte, row.initlen, row.initcap)
		if row.initcap == row.newcap {
			copy(initExpected, row.init)
		}

		new, err := ReadAll(init, strings.NewReader(row.readersrc))
		if err != row.err {
			t.Errorf("%v != %v", err, row.err)
		}

		if !reflect.DeepEqual(init, initExpected) {
			t.Errorf("%v != %v", init, initExpected)
		}

		newExpected := make([]byte, row.newlen, row.newcap)
		copy(newExpected, row.init+row.readersrc)

		if !reflect.DeepEqual(new, newExpected) {
			t.Errorf("%v != %v", new, newExpected)
		}

		if cap(new) != row.newcap {
			t.Errorf("%v != %v", cap(new), row.newcap)
		}
	}
}
