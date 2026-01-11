package aeon

import (
	"fmt"
	"testing"
	"time"

	"github.com/dromara/carbon/v2"
	_ "github.com/dromara/carbon/v2"
	_ "github.com/jinzhu/now"
)

func TestExample(_ *testing.T) {
	// t := Parse("20200805131415.999999999")
	//
	// fmt.Println(t)
	//
	carbon.NewCarbon().IsLongYear()
	Aeon(time.Now())
	t := Unix(-1, true)
	fmt.Println(t.ToString(DTFull)) // 1969-12-31 23:59:59 +0000 UTC
}
