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
		// 2024-06-15 -> 2024-06-20 -> final(End) -> 2024-06-20 23:59:59
		assert(t, base.EndAtMonth(6, 5), "2024-06-20 23:59:59", "EndAtMonth(6, 5)")

		// 3. StartAtYear(0, 1): 保持当年 (2024)，加 1 个月
		// 2024-04-15 -> 2024-05-15 -> final -> 2024-05-01
		assert(t, base.StartAtYear(0, 1), "2024-05-01 00:00:00", "StartAtYear(0, 1)")
	})

	t.Run("In 系列 (Rel + Abs...)", func(t *testing.T) {
		// 1. StartInYear(1, 6): 加 1 年，然后定位到 6 月
		// 2024-04-15 -> 2025-04-15 -> 2025-06-15 -> final -> 2025-06-01
		assert(t, base.StartInYear(1, 6), "2025-06-01 00:00:00", "StartInYear(1, 6)")

		// 2. StartInYear(-1, 6): 回跳 1 年，然后定位到 6 月
		assert(t, base.StartInYear(-1, 6), "2023-06-01 00:00:00", "StartInYear(-1, 6)")

		// 3. EndInYear(-1, 6): 回跳 1 年，然后定位到 6 月末
		assert(t, base.EndInYear(-1, 6), "2023-06-30 23:59:59", "EndInYear(-1, 6)")

		// 4. EndInMonth(2, 20): 加 2 个月，然后定位到 20 号
		// 2024-04-15 -> 2024-06-15 -> 2024-06-20 -> final -> 2024-06-20 23:59:59
		assert(t, base.EndInMonth(2, 20), "2024-06-20 23:59:59", "EndInMonth(2, 20)")

		// 3. StartInDay(10, 12): 加 10 天 (4/25)，然后定位到 12 时
		// 2024-04-15 12:00 -> 2024-04-25 12:00 -> 2024-04-25 12:00 -> final -> 2024-04-25 12:00:00
		assert(t, base.StartInDay(10, 12), "2024-04-25 12:00:00", "StartInDay(10, 12)")
	})

	t.Run("魔鬼矩阵：混合模式深度穿透", func(t *testing.T) {
		// 1. At 系列：跨年借位穿透
		// 锁定 2024 -> Rel Month(-13) -> 2023-03 -> Align 1号
		assert(t, base.StartAtYear(0, -13), "2023-03-01 00:00:00", "StartAtYear(0, -13)")

		// 2. In 系列：跨月后的绝对锚定
		// 4/15 + 40天 -> 5/25 -> Abs Hour(2) 强制定位
		assert(t, base.StartInDay(40, 2), "2024-05-25 02:00:00", "StartInDay(40, 2)")

		// 3. 跨度混合：从世纪到月份的折返
		// 逻辑：At Century(20) -> 2000 + 2000 = 4000 | Rel Decade(9) -> 4090 | Rel Year(2) -> 4092 | Rel Month(-5) -> 4091-11
		assert(t, base.StartAt(20, 9, 2, -5), "4091-11-01 00:00:00", "StartAt(20, 9, 2, -5)")

		// 4. ISO 周逻辑的级联渗透
		// At ISO W01 -> 2024-01-01 (周一) | Rel Weekday(1) 相对位移 1 天 -> 2024-01-02
		assert(t, base.StartAtYearWeek(ISO, 1, 1), "2024-01-02 00:00:00", "StartAtYearWeek(ISO, 1, 1)")

		// 5. Overflow 极限削峰测试
		ref := Parse("2024-01-31 12:00:00")
		// At Month(2) -> 02-29 | Rel Day(31) + Overflow -> 拦截削峰至 02-29
		assert(t, ref.StartAtMonth(Overflow, 2, 31), "2024-02-29 00:00:00", "StartAtMonth(Overflow, 2, 31)")

		// 6. EndAt 边界逆向折返
		// At Month(2) -> 02-15 | Rel Day(-5) -> 02-10 | End 对齐
		assert(t, base.EndAtMonth(2, -5), "2024-02-10 23:59:59", "EndAtMonth(2, -5)")
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
}
