package thru

import (
	"testing"
	"time"
)

func TestAbsSeriesDevilMatrix(t *testing.T) {
	// 基准时间: 2024-04-15 14:30:45 (21世纪, 2020年代, 2024年, 4月)
	base := Parse("2024-04-15 14:30:45")

	t.Run("核心纵向级联 (Century -> Second)", func(t *testing.T) {
		// 1. 全链路定位: Start(0, 2, 5, 5, 20, 10, 30, 0)
		// 2000(C) -> 2020(D) -> 2025(Y) -> 5月(M) -> 20日(D) -> 10:30 -> 45s(Second=0保持)
		assert(t, base.Start(0, 2, 5, 5, 20, 10, 30, 0), "2025-05-20 10:30:45", "全链路Start级联")

		// 2. 深度End级联: EndYear(5, 6, 20, 15)
		// 2025年(Y) -> 6月(M) -> 20日(D) -> 15时(H) -> 置满后续
		assert(t, base.EndYear(5, 6, 20, 15), "2025-06-20 15:59:59", "深度End级联")
	})

	t.Run("边界对齐与扩张 (n=0 语义)", func(t *testing.T) {
		// 1. Start 归零对齐
		assert(t, base.StartYear(0), "2024-01-01 00:00:00", "StartYear(0) 归零")
		assert(t, base.StartMonth(0), "2024-04-01 00:00:00", "StartMonth(0) 归零")
		assert(t, base.StartDecade(1), "2010-01-01 00:00:00", "StartDecade(1) 锚定本世纪第1个年代")

		// 2. End 置满扩张 (含 Decade/Century 修复验证)
		assert(t, base.EndYear(0), "2024-12-31 23:59:59", "EndYear(0) 置满")
		assert(t, base.EndDecade(0), "2029-12-31 23:59:59", "EndDecade(0) 扩张")
		assert(t, base.End(0), "2099-12-31 23:59:59", "EndCentury(0) 扩张")
		assert(t, base.EndQuarter(0), "2024-06-30 23:59:59", "EndQuarter(0) 扩张")
	})

	t.Run("日期保护与自然溢出", func(t *testing.T) {
		// 1. 月份保护: 1月31日定位到2月 -> 2月29日 (2024闰年)
		refJan31 := Parse("2024-01-31 12:00:00")
		assert(t, refJan31.StartMonth(2), "2024-02-01 00:00:00", "月份保护(Start)")
		assert(t, refJan31.EndMonth(2), "2024-02-29 23:59:59", "月份保护(End)")

		// 2. 位移单位自然溢出: 4月31日 -> 5月1日
		assert(t, base.StartDay(31), "2024-05-01 00:00:00", "Day自然溢出")
		assert(t, base.StartHour(25), "2024-04-16 01:00:00", "Hour自然溢出")
	})

	t.Run("倒数绝对索引 (Negative Index)", func(t *testing.T) {
		// 1. 月份倒数: -1 = 12月
		assert(t, base.StartMonth(-1), "2024-12-01 00:00:00", "月份倒数")
		// 2. 天数倒数: -1 = 本月最后一天
		assert(t, base.StartDay(-1), "2024-04-30 00:00:00", "天数倒数")
		// 3. 级联倒数: 年代末年末
		assert(t, base.StartYear(-1, -1), "2029-12-01 00:00:00", "级联倒数")
	})

	t.Run("周与主权/ISO 逻辑矩阵", func(t *testing.T) {
		// 1. ISO 跨年周: 2026-01-01 属于 2026-W01 (周一为 2025-12-29)
		refISO := Parse("2026-01-01 12:00:00")
		assert(t, refISO.StartYearWeek(ISO, 1), "2025-12-29 00:00:00", "ISO W01 定位")
		assert(t, refISO.StartYearWeek(ISO, 1, 2), "2025-12-30 00:00:00", "ISO W01 周二级联")

		// 2. 主权周 (Sovereignty): 2026 第一个周一 (1/5)
		assert(t, refISO.WithWeekStartsAt(time.Monday).StartYearWeek(1), "2026-01-05 00:00:00", "主权周W01")

		// 3. Weekday 绝对定位
		assert(t, base.WithWeekStartsAt(time.Monday).StartWeekday(1), "2024-04-15 00:00:00", "Weekday 1(Mon)")
		assert(t, base.WithWeekStartsAt(time.Monday).StartWeekday(7), "2024-04-21 00:00:00", "Weekday 7(Sun)")
	})

	t.Run("极端年份与安全边界", func(t *testing.T) {
		// 1. 千年边界级联: Century(-1) -> y=2900 -> Decade(-1) -> y=2990 -> Year(-1) -> y=2999
		assert(t, base.Start(-1, -1, -1, -1, -1), "2999-12-31 00:00:00", "千年边界")

		// 2. 安全性: Start(ISO) 无参数不崩溃
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Start(ISO) 崩溃: %v", r)
			}
		}()
		base.Start(ISO)
	})
}
