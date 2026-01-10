package thru

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	// 基准时间: 2024-01-30 15:04:05 (闰年)
	baseTime := time.Date(2024, 1, 30, 15, 4, 5, 0, time.Local)
	base := New(baseTime)

	// 1. 基础年/月/日 (正负测试)
	assert(t, base.AddYear(1), "2025-01-30 15:04:05", "AddYear(1)")
	assert(t, base.AddYear(-1), "2023-01-30 15:04:05", "AddYear(-1)")

	assert(t, base.AddMonth(1), "2024-02-29 15:04:05", "AddMonth(1) - 闰年2月截断")
	assert(t, base.AddMonth(-1), "2023-12-30 15:04:05", "AddMonth(-1) - 跨年回退")

	assert(t, base.AddDay(1), "2024-01-31 15:04:05", "AddDay(1)")
	assert(t, base.AddDay(-30), "2023-12-31 15:04:05", "AddDay(-30) - 跨年回退")

	// 2. 季度 (Quarter) - 验证特殊的对齐行为
	// 场景A: 1月30日 (Q1首月) -> 加1季度 -> 4月30日
	assert(t, base.AddQuarter(1), "2024-04-30 15:04:05", "AddQuarter(1) - Q1首月")

	// 场景B: 2月29日 (Q1中间) -> 加1季度
	// 逻辑: 2月先回退到Q1首月(1月)，再加3个月 -> 4月。结果应为 4月29日。
	feb29 := base.AddMonth(1) // 2024-02-29
	assert(t, feb29.AddQuarter(1), "2024-05-29 15:04:05", "AddQuarter(1) - 非首月被强制对齐到季度首月偏移")

	// 3. 多参数级联 (正向)
	// 2024-01-30 -> +1年(2025-01-30) -> +2月(2025-03-30) -> +3天(2025-04-02)
	assert(t, base.AddYear(1, 2, 3), "2025-04-02 15:04:05", "AddYear(1, 2, 3) - 多级正向")

	// 4. 多参数级联 (负向)
	// 2024-01-30 -> -1年(2023-01-30) -> -2月(2022-11-30) -> -3天(2022-11-27)
	assert(t, base.AddYear(-1, -2, -3), "2022-11-27 15:04:05", "AddYear(-1, -2, -3) - 多级负向")

	// 5. 季度多参数
	// 2024-01-30 -> +1季度(2024-04-30) -> +2月(2024-06-30)
	assert(t, base.AddQuarter(1, 2), "2024-06-30 15:04:05", "AddQuarter(1, 2) - 季度+月")

	// 6. 周与大参数
	assert(t, base.AddWeek(2), "2024-02-13 15:04:05", "AddWeek(2)")
	assert(t, base.AddWeek(-1), "2024-01-23 15:04:05", "AddWeek(-1)")

	// 7. 时分秒多级
	// 15:04:05 -> +1h(16:04:05) -> +1m(16:05:05) -> +1s(16:05:06)
	assert(t, base.AddHour(1, 1, 1), "2024-01-30 16:05:06", "AddHour(1, 1, 1)")
}

func TestAdd_Advanced(t *testing.T) {
	// 1. 纳秒进位测试
	// 2024-01-30 15:04:05.999999999 + 2ns -> 2024-01-30 15:04:06.000000001
	t.Run("Nano Rollover", func(t *testing.T) {
		base := ParseByLayout("2006-01-02 15:04:05.000000000", "2024-01-30 15:04:05.999999999")
		// 加 2 ns -> 应该进位到下一秒的 .000000001
		next := base.AddNano(2)
		assert(t, next, "2024-01-30 15:04:06.000000001", "AddNano 进位测试")
	})

	// 2. DST (夏令时) 墙上时间语义测试
	t.Run("DST Wall Clock Behavior", func(t *testing.T) {
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			t.Skip("Skipping DST test: America/Los_Angeles location not found")
		}

		// 2024-03-10 02:00:00 DST 开始，时间向前跳 1 小时 (02:00 -> 03:00)
		// 设定基准：2024-03-09 10:00:00 (PST)
		base := Date(2024, 3, 9, 10, 0, 0, 0, loc)

		// AddDay(1) -> 应该是 2024-03-10 10:00:00 (PDT)
		// 尽管实际上只过了 23 小时，但墙上时间保持 10:00
		dayAdded := base.AddDay(1)
		if h := dayAdded.Hour(); h != 10 {
			t.Errorf("AddDay(1) across DST: expected hour 10, got %d", h)
		}

		// 对比：标准库 Add(24h)
		// 24小时后应该是 2024-03-10 11:00:00 (PDT)
		stdAdd := base.Time().Add(24 * time.Hour)
		if h := stdAdd.Hour(); h != 11 {
			t.Errorf("StdLib Add(24h) across DST: expected hour 11, got %d", h)
		}

		// AddHour(24) -> 我们的 AddHour 是基于日历数字加法
		// 10 + 24 = 34. 34 % 24 = 10. +1天.
		// 所以应该是 2024-03-10 10:00:00 (PDT)
		hourAdded := base.AddHour(24)
		if h := hourAdded.Hour(); h != 10 {
			t.Errorf("AddHour(24) across DST: expected hour 10 (wall clock), got %d", h)
		}
	})

	// 3. 极值测试
	t.Run("Extreme Values", func(t *testing.T) {
		base := Now()
		// 加 1000 年
		future := base.AddYear(1000)
		if future.Year() != base.Year()+1000 {
			t.Errorf("AddYear(1000): expected year %d, got %d", base.Year()+1000, future.Year())
		}
	})
}
