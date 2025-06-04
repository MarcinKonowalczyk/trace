// micro version of assert to copy/paste into other packages `internal`
package muert

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// Check that the predicate is true, otherwise it fail the test.
func That(t testing.TB, predicate bool, args ...any) { t.Helper(); assert(t, 2, predicate, args) }

// Check that two comparable values are equal, otherwise it fail the test.
func Equal[T comparable](t testing.TB, a T, b T) {
	t.Helper()
	assert(t, 2, a == b, []any{"expected '%v' (%T) == '%v' (%T)", a, a, b, b})
}

// Check that an error is nil.
func NoError(t testing.TB, err error, args ...any) { t.Helper(); assert(t, 2, err == nil, args) }

// Check that the error is not nil and contains the expected message.
func Error(t testing.TB, err error, expected string, args ...any) {
	t.Helper()
	errs := err.Error()
	assert(t, 2, strings.Contains(errs, expected), []any{
		"expected error to contain '%s', got '%s' (%T): %v",
		expected, errs, err, args,
	})
}

func getParentInfo(N int) (string, int) {
	parent, _, _, _ := runtime.Caller(1 + N)
	return runtime.FuncForPC(parent).FileLine(parent)
}

// convert 'args ...any' to the assertion message
// internal utility so we don't use variadics to make the calls a bit more consistent
func argsToMessage(args []any) string {
	var msg string = "assertion failed"
	if len(args) > 0 {
		switch a := args[0].(type) {
		case string:
			msg = fmt.Sprintf(a, args[1:]...)
		default:
			msg = fmt.Sprintf("%v", args)
		}
	}
	return msg
}

func assert(t testing.TB, N int, predicate bool, args []any) {
	t.Helper()
	if !predicate {
		file, line := getParentInfo(N)
		t.Errorf(argsToMessage(args)+" in %s:%d", file, line)
	}
}
