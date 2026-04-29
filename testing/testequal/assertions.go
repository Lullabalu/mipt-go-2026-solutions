//go:build !solution

package testequal

import "fmt"

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.

func equal(expected, actual any) bool {
	if expected == nil || actual == nil {
		if expected == nil && actual == nil {
			return true
		}
		return false
	}

	switch ex := expected.(type) {
	case int8:
		ac, ok := actual.(int8)
		return ok && ac == ex
	case int16:
		ac, ok := actual.(int16)
		return ok && ac == ex
	case int32:
		ac, ok := actual.(int32)
		return ok && ac == ex
	case int64:
		ac, ok := actual.(int64)
		return ok && ac == ex
	case int:
		ac, ok := actual.(int)
		return ok && ac == ex
	case uint8:
		ac, ok := actual.(uint8)
		return ok && ac == ex
	case uint16:
		ac, ok := actual.(uint16)
		return ok && ac == ex
	case uint32:
		ac, ok := actual.(uint32)
		return ok && ac == ex
	case uint64:
		ac, ok := actual.(uint64)
		return ok && ac == ex
	case uint:
		ac, ok := actual.(uint)
		return ok && ac == ex
	case string:
		ac, ok := actual.(string)
		return ok && ac == ex
	case map[string]string:
		ac, ok := actual.(map[string]string)
		if !ok || len(ex) != len(ac) {
			return false
		}
		if ac == nil || ex == nil {
			if (ac == nil) != (ex == nil) {
				return false
			}
			return true
		}

		for k, v := range ex {
			if v != ac[k] {
				return false
			}
		}

		for k, v := range ac {
			if v != ex[k] {
				return false
			}
		}
		return true
	case []int:
		ac, ok := actual.([]int)
		if !ok || len(ac) != len(ex) {
			return false
		}
		if ac == nil || ex == nil {
			if (ac == nil) != (ex == nil) {
				return false
			}
			return true
		}
		for i := range len(ex) {
			if ac[i] != ex[i] {
				return false
			}
		}
		return true
	case []byte:
		ac, ok := actual.([]byte)
		if !ok || len(ac) != len(ex) {
			return false
		}
		if ac == nil || ex == nil {
			if (ac == nil) != (ex == nil) {
				return false
			}
			return true
		}
		for i := range len(ex) {
			if ac[i] != ex[i] {
				return false
			}
		}
		return true
	}
	return false
}

func AssertEqual(t T, expected, actual any, msgAndArgs ...any) bool {
	t.Helper()
	if !equal(expected, actual) {
		msg := ""
		if len(msgAndArgs) > 0 {
			format, ok := msgAndArgs[0].(string)
			if ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs...)
			}
		}

		t.Errorf(
			"not equal: \n expected: %v \n actual  : %v \n message : %s", expected, actual, msg,
		)
		return false
	}

	return true
}

func AssertNotEqual(t T, expected, actual any, msgAndArgs ...any) bool {
	t.Helper()
	if equal(expected, actual) {
		msg := ""
		if len(msgAndArgs) > 0 {
			format, ok := msgAndArgs[0].(string)
			if ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs...)
			}
		}

		t.Errorf(
			"equal: \n expected: %v \n actual  : %v \n message : %s", expected, actual, msg,
		)
		return false
	}

	return true
}

func RequireEqual(t T, expected, actual any, msgAndArgs ...any) {
	t.Helper()
	if !equal(expected, actual) {
		msg := ""
		if len(msgAndArgs) > 0 {
			format, ok := msgAndArgs[0].(string)
			if ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs...)
			}
		}

		t.Errorf(
			"not equal: \n expected: %v \n actual  : %v \n message : %s", expected, actual, msg,
		)
		t.FailNow()
	}

}

func RequireNotEqual(t T, expected, actual any, msgAndArgs ...any) {
	t.Helper()
	if equal(expected, actual) {
		msg := ""
		if len(msgAndArgs) > 0 {
			format, ok := msgAndArgs[0].(string)
			if ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs...)
			}
		}

		t.Errorf(
			"equal: \n expected: %v \n actual  : %v \n message : %s", expected, actual, msg,
		)
		t.FailNow()
	}
}
