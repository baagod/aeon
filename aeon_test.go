package thru

import (
	"testing"
)

func TestLog(t *testing.T) {
	// ref := Parse("2023-01-31")
	// fmt.Println(ref.JumpMonth(Overflow, 2)) // 2023-01-31 00:00:00
	//
	// tt := time.UnixMilli(1736416800000)
	// fmt.Println(tt.Format(time.DateTime))
	//
	// ref = Unix(1736416800000)
	// fmt.Println(ref)
}

func assert(t *testing.T, actual Time, expected string, msg string) {
	t.Helper()
	if actual.String() != expected {
		t.Errorf("%s, got [%s], want [%s]", msg, actual, expected)
	}
}
