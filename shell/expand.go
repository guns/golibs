// Package shell provides conveniences found in a Unix shell.
package shell

import (
	"os/user"
	"strings"
)

// ExpandUserDir expands leading `~user` HOME expansions.
func ExpandUserDir(s string) (string, error) {
	if len(s) == 0 || s[0] != '~' {
		return s, nil
	}

	i := strings.IndexByte(s, '/')
	var u *user.User
	var err error

	switch i {
	case -1:
		u, err = user.Current() // "~"
		if err != nil {
			return "", err
		}
		return u.HomeDir, nil
	case 1:
		u, err = user.Current() // "~/"
	default:
		u, err = user.Lookup(s[1:i]) // "~user/"
	}

	if err != nil {
		return "", err
	}

	return u.HomeDir + s[i:], nil
}
