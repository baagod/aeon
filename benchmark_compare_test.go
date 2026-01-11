package aeon

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/dromara/carbon/v2"
)

// --- 1. 创建性能 ---

func BenchmarkCreate_Aeon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewDate(2025, 1, 1)
	}
}

func BenchmarkCreate_Carbon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = carbon.CreateFromDate(2025, 1, 1)
	}
}

// --- 2. 深度级联对决 (世纪->年代->年->月->日) ---

func BenchmarkDeepCalc_Aeon(b *testing.B) {
	t := Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 🦬 单次级联穿透
		_ = t.GoDecade(2, 5, 5, 20)
	}
}

func BenchmarkDeepCalc_Carbon(b *testing.B) {
	c := carbon.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Carbon 必须链式调用
		_ = c.StartOfCentury().AddDecades(2).AddYears(5).SetMonth(5).SetDay(20)
	}
}

// --- 3. 相对位移对决 (年+月+日) ---

func BenchmarkAdd_Aeon(b *testing.B) {
	t := Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.AddYear(1, 2, 3) // 🦬 参数级联
	}
}

func BenchmarkAdd_Carbon(b *testing.B) {
	c := carbon.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.AddYears(1).AddMonths(2).AddDays(3)
	}
}

func BenchmarkAdd_Std(b *testing.B) {
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.AddDate(1, 2, 3)
	}
}

// --- 4. JSON 序列化 ---

func BenchmarkJSON_Aeon(b *testing.B) {
	t := Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(t)
	}
}

func BenchmarkJSON_Carbon(b *testing.B) {
	c := carbon.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(c)
	}
}
