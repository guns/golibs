/*

Package check provides a simple validation framework. It is a distillation of
https://github.com/kat-co/vala by Katherine Cox-Buday.

For example:

	type form struct {
		host string
		port int
	}

	func portNumberIsValid(key string, port int) check.Fn {
		return func() (pass bool, key, msg string) {
			pass = port > 0 && port < 0x10000
			msg = "must be an integer between 0 and 65536"
			return
		}
	}

	func stringIsNotEmpty(key, s string) check.Fn {
		return func() (bool, string, string) {
			return len(s) > 0, key, "must not be blank"
		}
	}

	func validate(f form) error {
		return check.That(
			stringIsNotEmpty("host", f.host),
			portNumberIsValid("port", f.port),
		)
	}

	func APIHandler(w http.ResponseWriter, r *http.Request) {
		var f form
		// …
		if err := validate(f); err != nil {
			if m, ok := err.(check.ErrorMap); ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(m)
				return
			}
			// …
		}
		// …
	}

An API request failing validation would get the following JSON response:

	{
		"host": "must not be blank",
		"port": "must be an integer between 0 and 65536"
	}

*/
package check

import "strings"

// An ErrorMap stores failed error messages in a map for easy retrieval. Note
// that entries are strings, not errors.
type ErrorMap map[string]string

// A Fn is a validation function. Return values are:
//
//	pass:     the test result
//	errorKey: the map key for the errorMsg entry
//	errorMsg: an error message
//
// errorKey and errorMsg should always be returned so that composable check
// functions like Not() can be written.
type Fn func() (pass bool, errorKey, errorMsg string)

// Error returns a summary of the error message map. It also implements the
// error interface, allowing easy travel through a typical function stack. The
// original ErrorMap can be regained through a type assertion.
func (m ErrorMap) Error() string {
	if len(m) == 0 {
		return "validation passed"
	}

	errors := make([]string, len(m))
	i := 0
	for k, v := range m {
		errors[i] = k + " " + v
		i++
	}
	return "validation failed: " + strings.Join(errors, ", ")
}

// That runs Fns and returns nil if all Fns passed, and returns a non-nil
// error interface value with concrete type ErrorMap if any failed.
func That(fs ...Fn) error {
	var m ErrorMap
	for _, checker := range fs {
		if pass, key, msg := checker(); !pass {
			if len(m) == 0 {
				m = make(ErrorMap)
			}
			m[key] = msg
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

// Pipe creates a pipeline of Fns such that each Fn must pass before the
// next Fn is called. This pipeline Fn will return the values from the first
// failure, if any.
//
// e.g.
//	check.That(
//		check.Pipe(
//			lenWithinBounds("username", username, 4, 60),
//			isUnique("username", db, username),
//		),
//	)
//
func Pipe(fs ...Fn) Fn {
	return func() (pass bool, key, msg string) {
		for _, f := range fs {
			pass, key, msg = f()
			if !pass {
				return pass, key, msg
			}
		}
		return pass, key, msg
	}
}
