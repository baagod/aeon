package aeon

import (
    "fmt"
    "testing"
)

func TestFunc(_ *testing.T) {
    t := Parse("2026-01-24")
    t = t.StartWeekday(5, 18)

    fmt.Println(t)
}
