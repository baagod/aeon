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
	assert(t, feb29.AddQuarter(1), "2024-04-29 15:04:05", "AddQuarter(1) - 非首月被强制对齐到季度首月偏移")

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
