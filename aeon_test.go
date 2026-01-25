package aeon

import (
    "testing"
)

func assert(t *testing.T, actual Time, expected string, name string) {
    t.Helper()
    if actual.String() != expected {
        t.Errorf("%s: got [%s], want [%s]", name, actual, expected)
    }
}
