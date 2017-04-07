// Copyright (c) 2015-2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

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
	var path string

	switch i {
	case -1: // "~" or "~user"
		if len(s) == 1 {
			u, err = user.Current()
		} else {
			u, err = user.Lookup(s[1:])
		}
	case 1: // "~/"
		u, err = user.Current()
		path = s[i:]
	default: // "~user/"
		u, err = user.Lookup(s[1:i])
		path = s[i:]
	}

	if err != nil {
		return "", err
	}

	return u.HomeDir + path, nil
}
