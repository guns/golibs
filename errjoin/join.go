// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package errjoin provides a simple function that combines multiple errors.
package errjoin

import (
	"errors"
	"strings"
)

// Join joins errors with a given separator. Nil arguments are filtered, and
// if all arguments are nil, Join returns nil.
func Join(sep string, errs ...error) error {
	var errorStrings []string

	for i := range errs {
		if errs[i] != nil {
			errorStrings = append(errorStrings, errs[i].Error())
		}
	}

	if len(errorStrings) == 0 {
		return nil
	}

	return errors.New(strings.Join(errorStrings, sep))
}
