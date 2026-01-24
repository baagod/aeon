package aeon

import (
    "fmt"
    "testing"
)

func TestFunc(_ *testing.T) {
    t := Parse("2026-01-24")
    fmt.Println(t)
}
