package thru

import (
	"fmt"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	ref, _ := time.Parse("2024-01-01 00:00:23.000000000", time.DateTime+".000000000")
	ref = ref.Add(time.Millisecond)
	fmt.Println(ref.Format(time.DateTime + ".000000000"))
	// assert(t, ref.StartMilli(1), "2024-01-01 00:00:00.001000000", "")
}

func assert(t *testing.T, actual Time, expected string, msg string) {
	t.Helper()
	if actual.String() != expected {
		t.Errorf("%s, got [%s], want [%s]", msg, actual, expected)
	}
}
