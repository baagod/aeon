package aeon

import (
	"testing"

	"github.com/dromara/carbon/v2"
)

func BenchmarkGoYear_Aeon(b *testing.B) {
	t := Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.GoYear(2025)
	}
}

func BenchmarkSetYear_Carbon(b *testing.B) {
	c := carbon.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.SetYear(2025)
	}
}

func BenchmarkNoOverflow_Aeon(b *testing.B) {
	t := NewDate(2025, 1, 31)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.GoMonth(Overflow, 2)
	}
}

func BenchmarkNoOverflow_Carbon(b *testing.B) {
	c := carbon.CreateFromDate(2025, 1, 31)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.SetMonthNoOverflow(2)
	}
}
