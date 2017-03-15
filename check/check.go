/*

Package check provides a simple validation framework. It is a distillation of
https://github.com/kat-co/vala by Katherine Cox-Buday.

For example:

	type form struct {
		host string
		port int
	}

	func stringIsNotEmpty(s, key string) check.Fn {
		return func() (pass bool, key, msg string) {
			pass = len(s) > 0
			msg = key + " must not be blank"
			return
		}
	}

	func portNumberIsValid(port int, key string) check.Fn {
		return func() (pass bool, key, msg string) {
			pass = port > 0 && port < 0x10000
			msg = key + " must be an integer between 0 and 65536"
			return
		}
	}

	func validate(f form) error {
		return check.That(
			stringIsNotEmpty(f.host, "host"),
			portNumberIsValid(f.port, "port"),
		)
	}

	func APIHandler(w http.ResponseWriter, r *http.Request) {
		var f form
		// …
		if err := validate(f); err != nil {
			if m, ok := err.(check.ErrorMap); ok {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(m)
				return
			}
			// …
		}
		// …
	}

An API request failing validation would get the following JSON response:

	{
		"host": "host must not be blank",
		"port": "port must be an integer between 0 and 65536"
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
	for _, v := range m {
		errors[i] = v
		i++
	}
	return "validation failed: " + strings.Join(errors, ", ")
}

// That runs given Checkers and returns an ErrorMap.
func That(checkers ...Fn) ErrorMap {
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
