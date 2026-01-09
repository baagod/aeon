package thru

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	ref := Parse("2023-01-31")
	fmt.Println(ref.JumpMonth(Overflow, 2)) // 2023-01-31 00:00:00
}

func assert(t *testing.T, actual Time, expected string, msg string) {
	t.Helper()
	if actual.String() != expected {
		t.Errorf("%s, got [%s], want [%s]", msg, actual, expected)
	}
}
