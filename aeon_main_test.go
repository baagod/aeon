package aeon

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	x := Parse("2026-01-11 01:12:13.101")
	fmt.Println(x)
}
