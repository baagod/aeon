package thru

import (
	"fmt"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	ref := Parse("2024-01-30 00:00:00.123456789")
	fmt.Println(ref.AddDecade(1).Format(time.DateTime + ".000000000"))
}

func assert(t *testing.T, actual Time, expected string, msg string) {
	t.Helper()
	if actual.String() != expected {
		t.Errorf("%s, got [%s], want [%s]", msg, actual, expected)
	}
}
