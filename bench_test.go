package aeon

import (
    "testing"

    "github.com/dromara/carbon/v2"
)

// --- 1. 高维绝对级联 (Century -> Hour) ---
// 目标：2025年5月20日12点

func BenchmarkAeonCascadeAbs(b *testing.B) {
    t := Now()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = t.StartCentury(20, 2, 5, 5, 20, 12)
    }
}

func BenchmarkCarbonChainedAbs(b *testing.B) {
    c := carbon.Now()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = c.StartOfCentury().AddDecades(2).AddYears(5).SetMonth(5).SetDay(20).SetHour(12).StartOfHour()
    }
}

// --- 2. 季度流级联 (Quarter Sequence) ---

func BenchmarkAeonQuarterCascade(b *testing.B) {
    t := Now()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = t.StartQuarter(2, 2, 10)
    }
}

func BenchmarkCarbonQuarterManual(b *testing.B) {
    c := carbon.Now()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = c.StartOfQuarter().AddMonths(1).SetDay(10).StartOfDay()
    }
}

// --- 3. 黑色星期五 (寻找年底最后一个周五) ---

func BenchmarkAeonBlackFriday(b *testing.B) {
    t := Now()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = t.GoMonth(-1, -1).StartWeek(-1, 5)
    }
}

func BenchmarkCarbonBlackFriday(b *testing.B) {
    c := carbon.Now()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        target := c.EndOfYear()
        // Carbon v2 的周五判断 (0:周日, 5:周五)
        for target.DayOfWeek() != 5 {
            target = target.SubDay()
        }
        _ = target
    }
}

// --- 4. 级联路径 vs 链式路径 ---

func BenchmarkAeonMultiAdd(b *testing.B) {
    t := Now()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = t.ShYear(1, 2, 3, 4, 5, 6)
    }
}

func BenchmarkCarbonMultiAdd(b *testing.B) {
    c := carbon.Now()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = c.AddYears(1).AddMonths(2).AddDays(3).AddHours(4).AddMinutes(5).AddSeconds(6)
    }
}
