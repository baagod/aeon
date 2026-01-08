package thru

import (
	"testing"
)

func TestLog(t *testing.T) {
	ref := Parse("2024-05-15 12:00:00")
	assert(t, ref.StartDecade(1), "2010-01-01 00:00:00", "StartDecade()")
}

func assert(t *testing.T, actual Time, expected string, msg string) {
	t.Helper()
	if actual.String() != expected {
		t.Errorf("%s, got [%s], want [%s]", msg, actual, expected)
	}
}
