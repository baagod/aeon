package aeon

import (
	"fmt"
	"testing"
)

func TestExample(_ *testing.T) {
	t := Parse("2020-08-05 13:14:15.999999999")

	fmt.Println(t.AddDay())
}
