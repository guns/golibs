// Package check provides a simple validation framework. It is a distillation
// of https://github.com/kat-co/vala by Katherine Cox-Buday.
package check

import "strings"

// An ErrorMap stores failed error messages in a map for easy retrieval. Note
// that entries are strings, not errors.
type ErrorMap map[string]string

// A Checker is a validation function. A failure message should always be
// returned for composability.
type Checker func() (pass bool, errorKey, errorMsg string)

// Error returns a summary of the error message map. It also implements the
// error interface, allowing easy travel through a typical function stack. The
// original ErrorMap can be regained through a type assertion.
func (m ErrorMap) Error() string {
	if len(m) == 0 {
		return "validation passed"
	}

	errors := make([]string, len(m))
	i := 0
	for _, v := range m {
		errors[i] = v
		i++
	}
	return "validation failed: " + strings.Join(errors, ", ")
}

// That runs given Checkers and returns an ErrorMap.
func That(checkers ...Checker) ErrorMap {
	var m ErrorMap
	for _, checker := range checkers {
		if pass, key, msg := checker(); !pass {
			if len(m) == 0 {
				m = make(ErrorMap, len(checkers))
			}
			m[key] = msg
		}
	}
	return m
}
