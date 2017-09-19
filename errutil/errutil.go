// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// Package errutil provides a simple functions for working with errors.
package errutil

import (
	"errors"
	"strings"
)

// Join joins error messages with a given separator. Nil arguments are
// ignored, and if all arguments are nil, Join returns nil.
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

// First returns the first non-nil error in errs.
func First(errs ...error) error {
	for i := range errs {
		if errs[i] != nil {
			return errs[i]
		}
	}
	return nil
}
