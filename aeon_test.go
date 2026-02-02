package aeon

import (
    "fmt"
    "testing"
    "time"
)

func TestFunc(_ *testing.T) {
    t := Parse("2021-07-21 14:35:29")
    t = t.Round(15 * time.Minute)
    fmt.Println(t)
}

func assert(t *testing.T, actual Time, expected string, name string) {
    t.Helper()
    if actual.String() != expected {
        t.Errorf("%s: got [%s], want [%s]", name, actual, expected)
    }
}

func assertZone(t *testing.T, actual Time, expectedOffset int, name string) {
    t.Helper()
    _, offset := actual.time.Zone()
    if offset != expectedOffset {
        t.Errorf("%s zone offset: got [%d], want [%d]", name, offset, expectedOffset)
    }
}
