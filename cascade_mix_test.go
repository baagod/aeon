package thru

import (
	"testing"
)

func TestMixSeriesDevilMatrix(t *testing.T) {
	// 基准时间: 2024-04-15 12:00:00
	base := Parse("2024-04-15 12:00:00")

	t.Run("At 系列 (Abs + Rel...)", func(t *testing.T) {
		// 1. StartAtYear(5, 1): 定位到本年代第 5 年 (2025)，然后加 1 个月
		// 2024-04-15 -> Abs Year(5) -> 2025-04-15 -> Rel Month(1) -> 2025-05-15 -> final/align -> 2025-05-01
		assert(t, base.StartAtYear(5, 1), "2025-05-01 00:00:00", "StartAtYear(5, 1)")

		// 2. EndAtMonth(6, 5): 定位到 6 月，然后加 5 天
		// 2024-06-15 -> 2024-06-20 -> final(End) -> 2024-06-20 23:59:59.999999999
		assert(t, base.EndAtMonth(6, 5), "2024-06-20 23:59:59.999999999", "EndAtMonth(6, 5)")

		// 3. EndInYear(-1, 6): 回跳 1 年，然后定位到 6 月末
		assert(t, base.EndInYear(-1, 6), "2023-06-30 23:59:59.999999999", "EndInYear(-1, 6)")

		// 4. EndInMonth(2, 20): 加 2 个月，然后定位到 20 号
		// 2024-04-15 -> 2024-06-15 -> 2024-06-20 -> final -> 2024-06-20 23:59:59.999999999
		assert(t, base.EndInMonth(2, 20), "2024-06-20 23:59:59.999999999", "EndInMonth(2, 20)")

		// 6. EndAtMonth(2, -5): 定位到 2 月 -> 2024-02-15 | Rel Day(-5) -> 2024-02-10 | End 对齐
		assert(t, base.EndAtMonth(2, -5), "2024-02-10 23:59:59.999999999", "EndAtMonth(2, -5)")

	})

	t.Run("混合模式下的 Overflow 标志", func(t *testing.T) {
		// 基准: 2024-01-31
		ref := Parse("2024-01-31 12:00:00")

		// StartAtMonth(Overflow, 2, 0): 定位到 2 月 (保护到 2/29)，然后加 0 天
		// 此时末级单位是 Day (0), align(Day) 不重置日期，故保持在 2/29
		assert(t, ref.StartAtMonth(Overflow, 2, 0), "2024-02-29 00:00:00", "StartAtMonth with Overflow")

		// StartInMonth(Overflow, 1, 1): 加 1 月 (1/31 -> 2/29)，然后定位到 1 号
		// 此时末级单位是 Month (1), align(Month) 会将日期重置为 1 号
		assert(t, ref.StartInMonth(Overflow, 1, 1), "2024-02-01 00:00:00", "StartInMonth with Overflow")
	})

	t.Run("纳秒精度混合操作 (Mix Nano)", func(t *testing.T) {
		base := Parse("2024-01-01 00:00:00")

		// 1. StartAtMilli(501, 1)
		// Abs Milli(501) -> 500ms
		// Rel Micro(1) -> 500ms + 1us = .500001
		// Start -> .500001000
		assert(t, base.StartAtMilli(501, 1), "2024-01-01 00:00:00.500001000", "StartAtMilli(501, 1)")
	})
}
