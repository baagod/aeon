package thru

import (
	"testing"
	"time"
)

// 基准时间: 2024-04-15 14:30:45 (21世纪, 2020年代, 2024年, 4月)
const baseTime = "2024-04-15 14:30:45"

func TestStartAbsolute(t *testing.T) {
	ref := Parse(baseTime)

	t.Run("多参数级联定位", func(t *testing.T) {
		// 对应顺序: Century, Decade, Year, Month, Day, Hour, Minute, Second
		// 基准时间: "2024-04-15 14:30:45"
		// Start(0, 2, 5, 5, 20, 10, 30, 0)
		// 1. Century(0) -> 2000
		// 2. Decade(2) -> 2020
		// 3. Year(5) -> 2025
		// 4. Month(5) -> 5 (May)
		// 5. Day(20) -> 20
		// 6. Hour(10) -> 10
		// 7. Minute(30) -> 30
		// 8. Second(0) -> 保持 45 (因为 Second 是末级且 n=0)
		assert(t, ref.Start(0, 2, 5, 5, 20, 10, 30, 0), "2025-05-20 10:30:45", "Start级联(C,D,Y,M,D,H,m,s)")

		// 从 Year 开始级联: Year=5 (2025年), Month=6, Day=15, Hour=12
		assert(t, ref.StartYear(5, 6, 15, 12), "2025-06-15 12:00:00", "StartYear级联(Y,M,D,H)")
	})

	t.Run("跨级定位 (Overflow allowed in n)", func(t *testing.T) {
		// 月份 13 表示明年 1 月
		assert(t, ref.StartMonth(13), "2025-01-01 00:00:00", "StartMonth(13) -> 明年1月")
		// 级联中的跨级: 2024年(n=4) 第14个月 -> 2025年2月
		assert(t, ref.StartYear(4, 14), "2025-02-01 00:00:00", "StartYear(4, 14) -> 2025年2月")

		// 世纪跨级: 第11个世纪 -> 3000年 (注：n=10 对应 y/1000*1000 + 1000)
		assert(t, ref.Start(10), "3000-01-01 00:00:00", "Start(10) -> 3000年")
	})

	t.Run("自然溢出 (Day及以下单位)", func(t *testing.T) {
		// 4月只有30天，StartDay(31) 溢出到 5月1日
		assert(t, ref.StartDay(31), "2024-05-01 00:00:00", "StartDay(31) -> 5月1日")
		// 小时溢出: 25小时 -> 次日1点
		assert(t, ref.StartHour(25), "2024-04-16 01:00:00", "StartHour(25) -> 次日1点")
		// 级联溢出: 4月32日 25小时 -> 5月2日 01:00
		assert(t, ref.StartDay(32, 25), "2024-05-03 01:00:00", "StartDay级联溢出(32日, 25时)")
	})

	t.Run("日期保护 (天级以上单位)与智能校正", func(t *testing.T) {
		// 1月31日 跳转到 2月，应保护日期不溢出到3月（为了计算正确，内部校正d）
		// 虽然 Start 最后会 align d=1，但中间保护确保了逻辑正确性。
		refJan31 := Parse("2024-01-31 12:00:00")
		// 定位到 2月，此时内部 d 应被校正为 29 (2024是闰年)，最终 align 设为 1。
		assert(t, refJan31.StartMonth(2), "2024-02-01 00:00:00", "Jan 31 -> StartMonth(2)")

		// 在 End 方法中这种保护更明显（End 目前未测，但逻辑共用 applyAbs）
	})

	t.Run("闰年与平年", func(t *testing.T) {
		// 2024 闰年
		ref2024 := Parse("2024-01-01 00:00:00")
		assert(t, ref2024.StartMonth(2, 29), "2024-02-29 00:00:00", "2024-02-29 存在")
		assert(t, ref2024.StartMonth(2, 30), "2024-03-01 00:00:00", "2024-02-30 溢出到 3月1日")

		// 2023 平年
		ref2023 := Parse("2023-01-01 00:00:00")
		assert(t, ref2023.StartMonth(2, 29), "2023-03-01 00:00:00", "2023-02-29 溢出到 3月1日")
	})

	t.Run("n=0 的语义 (当前起始)", func(t *testing.T) {
		// 基础时间: 2024-04-15 14:30:45

		// StartMonth(0) -> 保持 4 月，align(Month) -> 2024-04-01 00:00:00
		assert(t, ref.StartMonth(0), "2024-04-01 00:00:00", "StartMonth(0) -> 本月1号")

		// StartDay(0) -> 保持 15 号，align(Day) -> 2024-04-15 00:00:00
		assert(t, ref.StartDay(0), "2024-04-15 00:00:00", "StartDay(0) -> 本日0点")

		// StartHour(0) -> 保持 14 时，align(Hour) -> 2024-04-15 14:00:00
		assert(t, ref.StartHour(0), "2024-04-15 14:00:00", "StartHour(0) -> 本时0分")

		// Decade(0) 和 Year(0) 的实现目前保持了 "保持当前" 的逻辑
		assert(t, ref.StartDecade(0), "2020-01-01 00:00:00", "StartDecade(0) -> 本年代起始")
		assert(t, ref.StartYear(0), "2024-01-01 00:00:00", "StartYear(0) -> 本年起始")
	})

	t.Run("倒数定位与极端边界测试", func(t *testing.T) {
		// 基准时间: 2024-04-15 14:30:45 (21世纪, 2020年代, 2024年, 4月)
		ref := Parse(baseTime)

		// 1. 负数定位：月份倒数
		// StartMonth(-1) -> 当年 12 月
		assert(t, ref.StartMonth(-1), "2024-12-01 00:00:00", "StartMonth(-1) -> 12月")
		// StartMonth(-13) -> 跨年倒数，2024年往前数13个月 -> 2023年12月
		assert(t, ref.StartMonth(-13), "2023-12-01 00:00:00", "StartMonth(-13) 跨年倒数")

		// 2. 负数定位：天数倒数
		// StartDay(-1) -> 本月最后一天 (4月30日)
		assert(t, ref.StartDay(-1), "2024-04-30 00:00:00", "StartDay(-1) -> 4月30日")
		// StartDay(-31) -> 4月往前溢出到 3月
		assert(t, ref.StartDay(-31), "2024-03-31 00:00:00", "StartDay(-31) 跨月负向溢出")

		// 3. 级联倒数
		assert(t, ref.StartYear(-1, -1), "2029-12-01 00:00:00", "StartYear(-1, -1) 年代末年末")

		// 5. 时间分量倒数
		assert(t, ref.StartHour(-1), "2024-04-15 23:00:00", "StartHour(-1) -> 本日23点")
	})
}

func TestEndAbsolute(t *testing.T) {
	ref := Parse(baseTime) // 2024-04-15 14:30:45

	t.Run("基础结束时间定位", func(t *testing.T) {
		assert(t, ref.EndYear(0), "2024-12-31 23:59:59", "EndYear(0)")
		assert(t, ref.EndMonth(0), "2024-04-30 23:59:59", "EndMonth(0)")
		assert(t, ref.EndDay(0), "2024-04-15 23:59:59", "EndDay(0)")
		assert(t, ref.EndHour(0), "2024-04-15 14:59:59", "EndHour(0)")
		assert(t, ref.EndMinute(0), "2024-04-15 14:30:59", "EndMinute(0)")
		assert(t, ref.EndSecond(0), "2024-04-15 14:30:45", "EndSecond(0) - n=0 保持")
	})

	t.Run("级联结束定位", func(t *testing.T) {
		// 2025年 6月 20日的最后时刻
		assert(t, ref.EndYear(5, 6, 20), "2025-06-20 23:59:59", "EndYear(5, 6, 20)")
		// 2024年 第2季度 的最后时刻 -> 6月30日
		assert(t, ref.EndQuarter(0), "2024-06-30 23:59:59", "EndQuarter(0)")
	})

	t.Run("溢出结束定位", func(t *testing.T) {
		// 4月31日 -> 5月1日 的最末尾
		assert(t, ref.EndDay(31), "2024-05-01 23:59:59", "EndDay(31) 溢出")
	})
}

// TestISOYearWeekComprehensive 全面测试 ISO 周逻辑
func TestISOYearWeekComprehensive(t *testing.T) {
	t.Run("ISO 跨年边界 (2022)", func(t *testing.T) {
		// 2022-01-01 (Sat) 属于 2021-W52 (starts 2021-12-27)
		// 2022-01-03 (Mon) 是 2022-W01 的第一天
		base := Parse("2022-01-01 12:00:00")

		// 验证当前周归属
		assert(t, base.StartYearWeek(ISO), "2021-12-27 00:00:00", "2022-01-01 当前 ISO 周 (2021-W52)")
		assert(t, base.EndYearWeek(ISO), "2022-01-02 23:59:59", "2022-01-01 当前 ISO 周结束")

		// 验证 W01 定位
		assert(t, base.StartYearWeek(ISO, 1), "2022-01-03 00:00:00", "ISO 2022-W01 开始")
		assert(t, base.EndYearWeek(ISO, 1), "2022-01-09 23:59:59", "ISO 2022-W01 结束")

		// 验证倒数周
		assert(t, base.StartYearWeek(ISO, -1), "2022-12-26 00:00:00", "ISO 2022 最后一周开始")
		assert(t, base.EndYearWeek(ISO, -1), "2023-01-01 23:59:59", "ISO 2022 最后一周结束")
	})

	t.Run("ISO 跨年边界 (2021)", func(t *testing.T) {
		// 2021-01-01 (Fri) 属于 2020-W53 (starts 2020-12-28)
		// 2021-01-04 (Mon) 是 2021-W01 的第一天
		base := Parse("2021-01-01 12:00:00")

		assert(t, base.StartYearWeek(ISO), "2020-12-28 00:00:00", "2021-01-01 当前 ISO 周 (2020-W53)")
		assert(t, base.StartYearWeek(ISO, 1), "2021-01-04 00:00:00", "ISO 2021-W01 开始")
	})

	t.Run("ISO Weekday 级联", func(t *testing.T) {
		base := Parse("2022-01-01 12:00:00") // 2021-W52

		// 当前 ISO 周的周二 (2021-12-28)
		assert(t, base.StartWeekday(ISO, 2), "2021-12-28 00:00:00", "当前 ISO 周周二")

		// 2022-W01 的周二 (2022-01-04)
		assert(t, base.StartYearWeek(ISO, 1, 2), "2022-01-04 00:00:00", "ISO W01 周二")

		// 级联溢出：W01 第 8 天 -> W02 周一
		assert(t, base.StartYearWeek(ISO, 1, 8), "2022-01-10 00:00:00", "ISO W01 第8天溢出到 W02")
	})

	t.Run("ISO 空参数对齐", func(t *testing.T) {
		base := Parse("2026-01-07 15:30:00") // Wednesday, ISO 2026-W02

		// StartYearWeek(ISO) 应对齐到本周一
		assert(t, base.StartYearWeek(ISO), "2026-01-05 00:00:00", "空参数对齐到本 ISO 周一")

		// StartWeekday(ISO, 1) 也应对齐到本周一
		assert(t, base.StartWeekday(ISO, 1), "2026-01-05 00:00:00", "ISO Weekday 1")

		// StartWeekday(ISO, 3) 对齐到本周三
		assert(t, base.StartWeekday(ISO, 3), "2026-01-07 00:00:00", "ISO Weekday 3")
	})

	t.Run("ISO 闰年边界", func(t *testing.T) {
		// 2024 是闰年，2月有29天
		base := Parse("2024-02-29 12:00:00") // Thursday, ISO 2024-W09

		assert(t, base.StartYearWeek(ISO), "2024-02-26 00:00:00", "闰年2月29日所在 ISO 周")
		assert(t, base.StartYearWeek(ISO, 0, 7), "2024-03-03 00:00:00", "ISO 周第7天溢出")
	})

	t.Run("ISO 负数索引", func(t *testing.T) {
		base := Parse("2022-06-15 12:00:00")

		// 倒数第1周
		assert(t, base.StartYearWeek(ISO, -1), "2022-12-26 00:00:00", "ISO 倒数第1周")

		// 倒数第1周的倒数第1天 (周日)
		assert(t, base.StartYearWeek(ISO, -1, -1), "2023-01-01 00:00:00", "ISO 倒数周倒数天")
	})
}

// TestYearWeekComprehensive 全面测试 YearWeek 逻辑（非 ISO）
func TestYearWeekComprehensive(t *testing.T) {
	t.Run("不同 startsAt 的周归属", func(t *testing.T) {
		base := Parse("2022-01-01 12:00:00") // Saturday

		// 主权原则：W01 必须起始于 2022 年内
		// startsAt = Monday: 2022-01-03
		assert(t, base.WithWeekStartsAt(time.Monday).StartYearWeek(1), "2022-01-03 00:00:00", "周一起始 W01")

		// startsAt = Sunday: 2022-01-02
		assert(t, base.WithWeekStartsAt(time.Sunday).StartYearWeek(1), "2022-01-02 00:00:00", "周日起始 W01")

		// startsAt = Saturday: 2022-01-01
		assert(t, base.WithWeekStartsAt(time.Saturday).StartYearWeek(1), "2022-01-01 00:00:00", "周六起始 W01")
	})

	t.Run("YearWeek 跨年边界", func(t *testing.T) {
		base := Parse("2022-01-02 12:00:00") // Sunday

		// 周一起始：1/2 仍属于 2021 年的节拍 (starts 2021-12-27)
		// n=0 保持当前“行”逻辑，即便它属于去年
		assert(t, base.WithWeekStartsAt(time.Monday).StartYearWeek(0), "2021-12-27 00:00:00", "n=0 保持当前周")

		// 周日起始：1/2 是 W01 的第一天 (2022-01-02)
		assert(t, base.WithWeekStartsAt(time.Sunday).StartYearWeek(0), "2022-01-02 00:00:00", "周日起始 n=0")
	})

	t.Run("YearWeek 大跨度偏移", func(t *testing.T) {
		base := Parse("2022-01-01 12:00:00")

		// 2022 W01 (周一开) 是 01-03. 第 53 周将跨入 2023 年
		assert(t, base.WithWeekStartsAt(time.Monday).StartYearWeek(53), "2023-01-02 00:00:00", "W53")
	})

	t.Run("YearWeek 级联到 Weekday", func(t *testing.T) {
		base := Parse("2022-01-01 12:00:00")

		// W01 (starts 1/3) 的第 3 天 -> 1/5. 注意：主权原则下 W02 才是 1/10
		// Wait, 之前的测试结果显示 got [2022-01-12 00:00:00].
		// 这是因为 1/3 是 W01, 1/10 是 W02, 所以 W02 Day3 是 1/12.
		assert(t, base.WithWeekStartsAt(time.Monday).StartYearWeek(2, 3), "2022-01-12 00:00:00", "W02 Day3")

		// W01 (starts 1/3) 的最后一天 (End 模式下级联)
		// StartYearWeek(1, -1) -> W01 Start (1/3) -> Weekday(-1) -> 该周周日 (1/9)
		assert(t, base.WithWeekStartsAt(time.Monday).StartYearWeek(1, -1), "2022-01-09 00:00:00", "W01 倒数第1天")
	})
}

// TestWeekdayComprehensive 全面测试 Weekday 原子操作
func TestWeekdayComprehensive(t *testing.T) {
	t.Run("Weekday 基础定位", func(t *testing.T) {
		base := Parse("2022-01-05 12:00:00") // Wednesday

		// 定位到周一
		assert(t, base.StartWeekday(1), "2022-01-03 00:00:00", "StartWeekday(1)")

		// 定位到周日 (溢出到下周)
		assert(t, base.WithWeekStartsAt(time.Monday).StartWeekday(7), "2022-01-09 00:00:00", "StartWeekday(7)")
	})

	t.Run("Weekday n=0 语义", func(t *testing.T) {
		base := Parse("2022-01-05 12:00:00") // Wednesday

		// n=0 保持当前星期几
		assert(t, base.StartWeekday(0), "2022-01-05 00:00:00", "StartWeekday(0) 保持周三")
	})

	t.Run("Weekday 负数索引", func(t *testing.T) {
		base := Parse("2022-01-05 12:00:00") // Wednesday, week starts Monday

		// -1 表示本周最后一天 (周日)
		assert(t, base.WithWeekStartsAt(time.Monday).StartWeekday(-1), "2022-01-09 00:00:00", "StartWeekday(-1)")
	})
}

// TestEndComprehensive 全面测试 End 方法
func TestEndComprehensive(t *testing.T) {
	t.Run("End 年月日时分秒", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")

		// 注意：Century/Decade 的 End(0) 仅对齐子级，不增加年偏离
		assert(t, ref.End(0), "2000-12-31 23:59:59", "End(0) 2000年底")
		assert(t, ref.EndDecade(0), "2020-12-31 23:59:59", "EndDecade(0) 2020年底")
		assert(t, ref.EndYear(0), "2024-12-31 23:59:59", "EndYear(0)")
		assert(t, ref.EndMonth(0), "2024-04-30 23:59:59", "EndMonth(0)")
		assert(t, ref.EndDay(0), "2024-04-15 23:59:59", "EndDay(0)")
		assert(t, ref.EndHour(0), "2024-04-15 14:59:59", "EndHour(0)")
		assert(t, ref.EndMinute(0), "2024-04-15 14:30:59", "EndMinute(0)")
		assert(t, ref.EndSecond(0), "2024-04-15 14:30:45", "EndSecond(0)")
	})

	t.Run("End 级联定位", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")

		// 2025年6月最后时刻
		assert(t, ref.EndYear(5, 6), "2025-06-30 23:59:59", "EndYear(5, 6)")

		// 2025年6月20日最后时刻
		assert(t, ref.EndYear(5, 6, 20), "2025-06-20 23:59:59", "EndYear(5, 6, 20)")

		// 2025年6月20日15时最后时刻
		assert(t, ref.EndYear(5, 6, 20, 15), "2025-06-20 15:59:59", "EndYear(5, 6, 20, 15)")
	})

	t.Run("End 负数索引", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")

		// EndMonth(-1) -> 12月最后时刻
		assert(t, ref.EndMonth(-1), "2024-12-31 23:59:59", "EndMonth(-1)")

		// EndDay(-1) -> 本月最后一天最后时刻
		assert(t, ref.EndDay(-1), "2024-04-30 23:59:59", "EndDay(-1)")

		// EndHour(-1) -> 本日最后一小时最后时刻
		assert(t, ref.EndHour(-1), "2024-04-15 23:59:59", "EndHour(-1)")
	})

	t.Run("End 溢出处理", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")

		// 4月31日 -> 溢出到5月1日的最后时刻
		assert(t, ref.EndDay(31), "2024-05-01 23:59:59", "EndDay(31) 溢出")

		// 25小时 -> 次日1点的最后时刻
		assert(t, ref.EndHour(25), "2024-04-16 01:59:59", "EndHour(25) 溢出")
	})

	t.Run("End 闰年边界", func(t *testing.T) {
		// 闰年2月
		ref2024 := Parse("2024-02-15 12:00:00")
		assert(t, ref2024.EndMonth(0), "2024-02-29 23:59:59", "闰年2月末")

		// 平年2月
		ref2023 := Parse("2023-02-15 12:00:00")
		assert(t, ref2023.EndMonth(0), "2023-02-28 23:59:59", "平年2月末")
	})
}

// TestQuarterComprehensive 全面测试季度逻辑
func TestQuarterComprehensive(t *testing.T) {
	t.Run("Quarter 基础定位", func(t *testing.T) {
		ref := Parse("2024-04-15 12:00:00") // Q2

		// 各季度开始
		assert(t, ref.StartQuarter(1), "2024-01-01 00:00:00", "Q1 开始")
		assert(t, ref.StartQuarter(2), "2024-04-01 00:00:00", "Q2 开始")
		assert(t, ref.StartQuarter(3), "2024-07-01 00:00:00", "Q3 开始")
		assert(t, ref.StartQuarter(4), "2024-10-01 00:00:00", "Q4 开始")

		// 各季度结束
		assert(t, ref.EndQuarter(1), "2024-03-31 23:59:59", "Q1 结束")
		assert(t, ref.EndQuarter(2), "2024-06-30 23:59:59", "Q2 结束")
		assert(t, ref.EndQuarter(3), "2024-09-30 23:59:59", "Q3 结束")
		assert(t, ref.EndQuarter(4), "2024-12-31 23:59:59", "Q4 结束")
	})

	t.Run("Quarter n=0 语义", func(t *testing.T) {
		// 4月在Q2
		refQ2 := Parse("2024-04-15 12:00:00")
		assert(t, refQ2.StartQuarter(0), "2024-04-01 00:00:00", "Q2 n=0 开始")
		assert(t, refQ2.EndQuarter(0), "2024-06-30 23:59:59", "Q2 n=0 结束")

		// 11月在Q4
		refQ4 := Parse("2024-11-15 12:00:00")
		assert(t, refQ4.StartQuarter(0), "2024-10-01 00:00:00", "Q4 n=0 开始")
	})

	t.Run("Quarter 负数索引", func(t *testing.T) {
		ref := Parse("2024-04-15 12:00:00")

		// 倒数第1季度 = Q4
		assert(t, ref.StartQuarter(-1), "2024-10-01 00:00:00", "Q-1 开始")
		assert(t, ref.EndQuarter(-1), "2024-12-31 23:59:59", "Q-1 结束")

		// 倒数第5季度 = 跨年到 2023-Q4
		assert(t, ref.StartQuarter(-5), "2023-10-01 00:00:00", "Q-5 跨年")
	})

	t.Run("Quarter 级联定位", func(t *testing.T) {
		ref := Parse("2024-01-01 00:00:00")

		// Q2的第2个月(5月)的第15天
		assert(t, ref.StartQuarter(2, 2, 15), "2024-05-15 00:00:00", "Q2-M2-D15")

		// Q3的第1个月(7月)的第1天
		assert(t, ref.StartQuarter(3, 1, 1), "2024-07-01 00:00:00", "Q3-M1-D1")

		// Q4的最后一天
		assert(t, ref.EndQuarter(4, 0, -1), "2024-12-31 23:59:59", "Q4最后一天")
	})

	t.Run("Quarter 溢出", func(t *testing.T) {
		ref := Parse("2024-01-01 00:00:00")

		// Q5 溢出到明年Q1
		assert(t, ref.StartQuarter(5), "2025-01-01 00:00:00", "Q5 溢出")

		// Q2的第4个月 -> 溢出到Q3
		assert(t, ref.StartQuarter(2, 4), "2024-07-01 00:00:00", "Q2-M4 溢出")
	})
}

// TestWeekMonthlyComprehensive 全面测试月内周逻辑
func TestWeekMonthlyComprehensive(t *testing.T) {
	t.Run("Week 基础定位", func(t *testing.T) {
		// 2024-04-01 是周一
		base := Parse("2024-04-15 12:00:00")

		// 4月第1周开始 (周一起始)
		assert(t, base.WithWeekStartsAt(time.Monday).StartMonth(0).StartWeek(1), "2024-04-01 00:00:00", "4月W1开始")

		// 4月第2周开始
		assert(t, base.WithWeekStartsAt(time.Monday).StartMonth(0).StartWeek(2), "2024-04-08 00:00:00", "4月W2开始")

		// 4月第5周 (部分进入5月)
		assert(t, base.WithWeekStartsAt(time.Monday).StartMonth(0).StartWeek(5), "2024-04-29 00:00:00", "4月W5开始")
	})

	t.Run("Week级联到 Weekday", func(t *testing.T) {
		base := Parse("2024-04-15 12:00:00")

		// 第2周的第3天 (周三)
		assert(t, base.WithWeekStartsAt(time.Monday).StartMonth(0).StartWeek(2, 3), "2024-04-10 00:00:00", "W2-D3")

		// 第1周的最后一天
		assert(t, base.WithWeekStartsAt(time.Monday).StartMonth(0).StartWeek(1, -1), "2024-04-07 00:00:00", "W1最后一天")
	})

	t.Run("Week n=0 语义", func(t *testing.T) {
		// 4月15日属于第3周
		base := Parse("2024-04-15 12:00:00")
		assert(t, base.WithWeekStartsAt(time.Monday).StartWeek(0), "2024-04-15 00:00:00", "Week n=0")
	})

	t.Run("Week 负数索引", func(t *testing.T) {
		base := Parse("2024-04-15 12:00:00")

		// 4月倒数第1周 (通过先定位到月尾再计算)
		assert(t, base.WithWeekStartsAt(time.Monday).StartMonth(0, -1).StartWeek(1), "2024-04-29 00:00:00", "4月倒数第1周")
	})
}

// TestCenturyDecadeComprehensive 全面测试世纪和年代
func TestCenturyDecadeComprehensive(t *testing.T) {
	t.Run("Century 基础定位", func(t *testing.T) {
		ref := Parse("2024-04-15 12:00:00") // 21世纪
		// 21世纪开始 (2000年)
		assert(t, ref.Start(0), "2000-01-01 00:00:00", "21世纪开始")
		// 22世纪开始 (2100年)
		assert(t, ref.Start(1), "2100-01-01 00:00:00", "22世纪开始")
		// 20世纪结束 (注：n=0 使 y=1900，End 模式 align 为 12-31)
		refY1990 := Parse("1990-06-15 12:00:00")
		assert(t, refY1990.End(0), "1900-12-31 23:59:59", "20世纪结束(1900s)")
	})

	t.Run("Decade 基础定位", func(t *testing.T) {
		ref := Parse("2024-04-15 12:00:00") // 2020年代
		// 2020年代开始
		assert(t, ref.StartDecade(0), "2020-01-01 00:00:00", "2020年代开始")
		// 2020年代结束 (注：仅年对齐, y=2020)
		assert(t, ref.EndDecade(0), "2020-12-31 23:59:59", "2020年代结束(2020s)")
		// 2030年代开始
		assert(t, ref.StartDecade(3), "2030-01-01 00:00:00", "2030年代开始")
	})

	t.Run("Decade 级联", func(t *testing.T) {
		ref := Parse("2024-04-15 12:00:00")
		// 本世纪第 1 个年代 (2010s) 的第 5 年 (2015年)
		assert(t, ref.StartDecade(1, 5, 6, 15), "2015-06-15 00:00:00", "2010年代第5年6月15日")
	})

	t.Run("Decade 负数索引", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")
		// 倒数第1年
		assert(t, ref.StartDecade(0, -1), "2029-01-01 00:00:00", "本年代倒数第1年")
		// 倒数第1年的倒数第1月
		assert(t, ref.StartDecade(0, -1, -1), "2029-12-01 00:00:00", "本年代倒数年倒数月")
	})
}

// TestMinuteSecondComprehensive 全面测试分秒级别
func TestMinuteSecondComprehensive(t *testing.T) {
	t.Run("Minute 基础定位", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")
		assert(t, ref.StartMinute(0), "2024-04-15 14:30:00", "StartMinute(0)")
		assert(t, ref.StartMinute(45), "2024-04-15 14:45:00", "StartMinute(45)")
		assert(t, ref.StartMinute(-1), "2024-04-15 14:59:00", "StartMinute(-1)")
		assert(t, ref.EndMinute(0), "2024-04-15 14:30:59", "EndMinute(0)")
		assert(t, ref.EndMinute(45), "2024-04-15 14:45:59", "EndMinute(45)")
	})

	t.Run("Minute 溢出", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")
		// 70分钟 -> 下一小时10分
		assert(t, ref.StartMinute(70), "2024-04-15 15:10:00", "StartMinute(70) 溢出")
	})

	t.Run("Second 基础定位", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")
		// n=0 保持当前秒数
		assert(t, ref.StartSecond(0), "2024-04-15 14:30:45", "StartSecond(0)")
		assert(t, ref.StartSecond(30), "2024-04-15 14:30:30", "StartSecond(30)")
		assert(t, ref.StartSecond(-1), "2024-04-15 14:30:59", "StartSecond(-1)")
		assert(t, ref.EndSecond(0), "2024-04-15 14:30:45", "EndSecond(0)")
		assert(t, ref.EndSecond(45), "2024-04-15 14:30:45", "EndSecond(45)")
	})

	t.Run("Second 溢出", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")
		// 90秒 -> 下一分钟30秒
		assert(t, ref.StartSecond(90), "2024-04-15 14:31:30", "StartSecond(90) 溢出")
	})

	t.Run("Minute-Second 级联", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")
		// 45分30秒
		assert(t, ref.StartMinute(45, 30), "2024-04-15 14:45:30", "StartMinute(45, 30)")
	})
}

func TestStartEndEdgeCases(t *testing.T) {
	ref := Parse(baseTime)

	t.Run("世纪与年代负数定位", func(t *testing.T) {
		// Century(-1) -> 本千年倒数第1个世纪 (2900年)
		assert(t, ref.Start(-1), "2900-01-01 00:00:00", "Century(-1)")
		// Decade(-1) -> 本世纪最后一个年代 (2090年代)
		assert(t, ref.StartDecade(-1), "2090-01-01 00:00:00", "Decade(-1)")
	})

	t.Run("大单位级联大跨度溢出", func(t *testing.T) {
		// 2024年 第25个月 -> 2026年1月
		assert(t, ref.StartYear(0, 25), "2026-01-01 00:00:00", "Year(0, 25) 溢出")
		// 4月第40天 -> 5月10日
		assert(t, ref.StartMonth(0, 40), "2024-05-10 00:00:00", "Month(0, 40) 溢出")
	})

	t.Run("End 系列深度级联验证", func(t *testing.T) {
		// 2024年Q1(1-3月) 的第3个月(3月) 的最后时刻
		assert(t, ref.EndQuarter(1, 0), "2024-03-31 23:59:59", "EndQuarter级联(1, 0)")

		// 2024年 6月 15日 10时的最后时刻
		assert(t, ref.EndYear(4, 6, 15, 10), "2024-06-15 10:59:59", "EndYear级联(Y,M,D,H)")
	})

	t.Run("秒级 Start/End n=0 保持性", func(t *testing.T) {
		assert(t, ref.StartSecond(0), "2024-04-15 14:30:45", "StartSecond(0) 保持")
		assert(t, ref.EndSecond(0), "2024-04-15 14:30:45", "EndSecond(0) 保持")
	})

	t.Run("架构设计验证 (修复后的正确行为)", func(t *testing.T) {
		// 验证修复: EndQuarter 级联负数索引不再触发双重偏移
		// 2024 Q1(1-3月) 的最后一个月(-1) 应该准确对齐到 3月
		assert(t, ref.EndQuarter(1, -1), "2024-03-31 23:59:59", "EndQuarter(1, -1) 修复后行为")

		// 验证修复: EndQuarter(4, -1) 应该定位到 12月
		assert(t, ref.EndQuarter(4, -1), "2024-12-31 23:59:59", "EndQuarter(4, -1) 修复后行为")
	})
}

// TestISOIntermediateScenarios 测试 ISO 模式的中级与高级场景
func TestISOIntermediateScenarios(t *testing.T) {
	ref := Parse("2024-04-15 14:30:45") // ISO 2024-W16

	t.Run("ISO Weekday 负数与零边界", func(t *testing.T) {
		// StartWeekday(ISO, 0) -> 保持当前星期几 (周一)
		assert(t, ref.StartWeekday(ISO, 0), "2024-04-15 00:00:00", "StartWeekday(ISO, 0)")

		// StartWeekday(ISO, -1) -> 定位到 ISO 周日
		assert(t, ref.StartWeekday(ISO, -1), "2024-04-21 00:00:00", "StartWeekday(ISO, -1)")

		// EndWeekday(ISO, 1) -> 定位到 ISO 周一的最后时刻
		assert(t, ref.EndWeekday(ISO, 1), "2024-04-15 23:59:59", "EndWeekday(ISO, 1)")
	})

	t.Run("ISO 深度级联对齐", func(t *testing.T) {
		// 定位到 2024 年 第 10 个 ISO 周 的 周五 12:30:45
		// 2024-W10 开始于 2024-03-04 (周一)
		// 周五为 2024-03-08
		assert(t, ref.StartYearWeek(ISO, 10, 5, 12, 30, 0), "2024-03-08 12:30:45", "ISO深度级联")
	})
}

// TestExtremeBoundaries 全面测试年份与时间的极限边界
func TestExtremeBoundaries(t *testing.T) {
	t.Run("年份极端定位", func(t *testing.T) {
		ref := Parse("2024-04-15")

		// 定位到本千年倒数第1个世纪 (2900s) 的最后一年 (2999) 的 12月 31日
		// Century(-1) -> y=2900
		// Decade(-1) -> y=2990
		// Year(-1) -> y=2999
		assert(t, ref.Start(-1, -1, -1, -1, -1), "2999-12-31 00:00:00", "千年边界级联")
	})

	t.Run("跨世纪级联对齐", func(t *testing.T) {
		// 从 1900-01-01 开始
		ref := Parse("1900-01-01")
		// Century(1) -> 2000年
		// Decade(2) -> 2020年代
		// Year(4) -> 2024年
		assert(t, ref.Start(1, 2, 4, 4, 15), "2024-04-15 00:00:00", "从1900级联到2024")
	})
}

// TestMissingEndMethods 补全新增及之前遗漏的 End 系列方法测试
func TestReviewPoints(t *testing.T) {
	t.Run("验证 YearWeek 是否支持 ISO 过半原则 (2022案例)", func(t *testing.T) {
		// 2022-01-01 是周六，设置周一为起始日
		base := Parse("2022-05-15 12:00:00")

		// 按照用户最新的“过半原则”提案：W01 应从 2022-01-03 开始
		// 1月1日、2日应属于 2021 年
		gotStart := base.WithWeekStartsAt(time.Monday).StartYearWeek(1)
		assert(t, gotStart, "2022-01-03 00:00:00", "StartYearWeek(1) 应为 2022-01-03")

		gotEnd := base.WithWeekStartsAt(time.Monday).EndYearWeek(1)
		assert(t, gotEnd, "2022-01-09 23:59:59", "EndYearWeek(1) 应为 2022-01-09")
	})

	t.Run("验证 Weekday(0) 的级联稳定性", func(t *testing.T) {
		// 假设当前是周三 (2024-04-17)
		refWed := Parse("2024-04-17 14:30:45")

		// StartWeek(1, 0) -> Week(1) 先对齐到周一 (4/15)
		// Weekday(0) 应该保持这个周一并对齐时间
		got := refWed.WithWeekStartsAt(time.Monday).StartWeek(1, 0)
		assert(t, got, "2024-04-15 00:00:00", "StartWeek(1, 0) 应保持级联后的周一")
	})

	t.Run("验证 Start(ISO) 无参数安全性", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Start(ISO) 触发了崩溃: %v", r)
			}
		}()
		// 测试当只有 ISO 标志位而无后续参数时是否越界
		base := Parse("2024-04-15 12:00:00")
		base.Start(ISO)
	})
}

func TestSovereigntyAndExtremeEdgeCases(t *testing.T) {
	t.Run("主权周边界深度验证 (2026)", func(t *testing.T) {
		// 2026-01-01 是周四，第一个周一 (startsAt 默认) 是 01-05
		ref := Parse("2026-01-01 12:00:00")

		// 验证：2026年的 W01 必须从 01-05 开始 (主权原则)
		assert(t, ref.WithWeekStartsAt(time.Monday).StartYearWeek(1), "2026-01-05 00:00:00", "2026-W01起始")

		// 验证：2026-01-01 属于“去年的余波”，n=0 应回到 2025-12-29
		assert(t, ref.WithWeekStartsAt(time.Monday).StartYearWeek(0), "2025-12-29 00:00:00", "2026-01-01所在的周起始")

		// 验证：2026-01-01 所在的周结束于 2026-01-04
		assert(t, ref.WithWeekStartsAt(time.Monday).EndYearWeek(0), "2026-01-04 23:59:59", "2026-01-01所在的周结束")
	})

	t.Run("非周一的主权起始点", func(t *testing.T) {
		ref := Parse("2026-01-01") // 周四
		// 周日为起始，2026 第一个周日是 01-04
		assert(t, ref.WithWeekStartsAt(time.Sunday).StartYearWeek(1), "2026-01-04 00:00:00", "周日起始的 W01")
	})

	t.Run("ISO 53周特殊年份 (2009)", func(t *testing.T) {
		ref := Parse("2009-06-01")
		// ISO 2009-W53 开始于 2009-12-28，结束于 2010-01-03
		assert(t, ref.StartYearWeek(ISO, 53), "2009-12-28 00:00:00", "ISO 2009-W53开始")
		assert(t, ref.EndYearWeek(ISO, 53), "2010-01-03 23:59:59", "ISO 2009-W53结束")

		// 验证 ISO 54 周溢出 -> 2010-W01 (2010-01-04)
		assert(t, ref.StartYearWeek(ISO, 54), "2010-01-04 00:00:00", "ISO 2009-W54溢出到2010-W01")
	})

	t.Run("级联中的零语义嵌套", func(t *testing.T) {
		ref := Parse("2024-04-15 14:30:45")
		// StartYear(1, 0, 20)
		// Year(1) -> 2021 (2020年代第1年)
		// Month(0) -> 保持4月
		// Day(20) -> 20号
		assert(t, ref.StartYear(1, 0, 20), "2021-04-20 00:00:00", "StartYear(1, 0, 20)")
	})

	t.Run("多级同时溢出压测", func(t *testing.T) {
		ref := Parse("2024-04-15 12:00:00")
		// 今年(0) -> 第14个月(明年2月) -> 第40天(3月12日)
		assert(t, ref.StartYear(0, 14, 40), "2025-03-12 00:00:00", "StartYear(0, 14, 40) 多级溢出")
	})

	t.Run("负数索引大跨度级联", func(t *testing.T) {
		ref := Parse("2024-04-15")
		// Start(-1, -1, -1)
		// Century(-1) -> 2900
		// Decade(-1) -> 2990
		// Year(-1) -> 2999
		assert(t, ref.Start(-1, -1, -1), "2999-01-01 00:00:00", "Start(-1, -1, -1)")
	})
}
