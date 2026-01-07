package thru

import (
	"testing"
)

func TestBySeriesDevilMatrix(t *testing.T) {
	t.Run("月末保护与级联位移 (2024)", func(t *testing.T) {
		// 基准: 2024-01-31 (1月最后一天)
		ref := Parse("2024-01-31 12:00:00")

		// 1. 默认 By 逻辑: StartByMonth(1)
		// Month + 1 -> 2/29 (保护) -> 2/1 (对齐生效)
		assertBy(t, ref.StartByMonth(1), "2024-02-01 00:00:00", "1月31日 + 1月 (保护+对齐)")

		// 2. 级联位移: StartByMonth(1, 1)
		// Month + 1 -> 2/29 (保护)
		// Day + 1 -> 3/1 (自然溢出)
		assertBy(t, ref.StartByMonth(1, 1), "2024-03-01 00:00:00", "1月31日 + 1月 + 1天")
	})

	t.Run("主权周与ISO周偏移 (2026)", func(t *testing.T) {
		// 基准: 2026-01-02 (周五) - 属于去年余波 (2025-W52)
		ref := Parse("2026-01-02 12:00:00")

		// 1. 保持当前主权周: StartByYearWeek(0) -> 2025-12-29
		assertBy(t, ref.StartByYearWeek(0), "2025-12-29 00:00:00", "2026-01-02当前主权周首")

		// 2. 下个主权周: StartByYearWeek(1) -> 2026-01-05
		assertBy(t, ref.StartByYearWeek(1), "2026-01-05 00:00:00", "2026-01-02下个主权周首")

		// 3. ISO 周强制周一: StartByYearWeek(ISO, 1)
		assertBy(t, ref.StartByYearWeek(ISO, 1), "2026-01-05 00:00:00", "2026-01-02下个ISO周首")
	})

	t.Run("动态账期末 (EndBy)", func(t *testing.T) {
		ref := Parse("2024-04-15 10:00:00")
		// 下个月的前一天结束: EndByMonth(1, -1)
		assertBy(t, ref.EndByMonth(1, -1), "2024-05-14 23:59:59", "下月账期前一天结束")
	})

	t.Run("巨量偏移压测", func(t *testing.T) {
		ref := Parse("2024-04-15 00:00:00")
		// 100万天之后 (校准后预期)
		assertBy(t, ref.StartByDay(1000000), "4762-03-13 00:00:00", "100万天后")
	})

	t.Run("负数跨年位移", func(t *testing.T) {
		ref := Parse("2024-01-01 12:00:00")
		assertBy(t, ref.StartByMonth(-1), "2023-12-01 00:00:00", "1月1日回退1月")
	})

	t.Run("显式强制保护 (Overflow Flag)", func(t *testing.T) {
		ref := Parse("2024-04-15 12:00:00")
		// 虽然 Day 本身不保护，但传了 Overflow 之后触发保护（截断到月末）：
		assertBy(t, ref.StartByDay(Overflow, 45), "2024-04-30 00:00:00", "显式触发保护的位移")
	})
}

func assertBy(t *testing.T, got Time, wantStr string, name string) {
	if got.String() != wantStr {
		t.Errorf("%s, got [%s], want [%s]", name, got.String(), wantStr)
	}
}
