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
		assert(t, ref.StartByMonth(1), "2024-02-01 00:00:00", "1月31日 + 1月 (保护+对齐)")

		// 2. 级联位移: StartByMonth(1, 1)
		// Month + 1 -> 2/29 (保护)
		// Day + 1 -> 3/1 (自然溢出)
		assert(t, ref.StartByMonth(1, 1), "2024-03-01 00:00:00", "1月31日 + 1月 + 1天")
	})

	t.Run("主权周与ISO周偏移 (2026)", func(t *testing.T) {
		// 基准: 2026-01-02 (周五) - 属于去年余波 (2025-W52)
		ref := Parse("2026-01-02 12:00:00")

		// 1. 保持当前主权周: StartByYearWeek(0) -> 2025-12-29
		assert(t, ref.StartByYearWeek(0), "2025-12-29 00:00:00", "2026-01-02当前主权周首")

		// 2. 下个主权周: StartByYearWeek(1) -> 2026-01-05
		assert(t, ref.StartByYearWeek(1), "2026-01-05 00:00:00", "2026-01-02下个主权周首")

		// 3. ISO 周强制周一: StartByYearWeek(ISO, 1)
		assert(t, ref.StartByYearWeek(ISO, 1), "2026-01-05 00:00:00", "2026-01-02下个ISO周首")
	})

	t.Run("动态账期末 (EndByCentury)", func(t *testing.T) {
		ref := Parse("2024-04-15 10:00:00")
		// 下个月的前一天结束: EndByMonth(1, -1)
		assert(t, ref.EndByMonth(1, -1), "2024-05-14 23:59:59.999999999", "下月账期前一天结束")
	})

	t.Run("巨量偏移压测", func(t *testing.T) {
		ref := Parse("2024-04-15 00:00:00")
		// 100万天之后 (校准后预期)
		assert(t, ref.StartByDay(1000000), "4762-03-13 00:00:00", "100万天后")
	})

	t.Run("负数跨年位移", func(t *testing.T) {
		ref := Parse("2024-01-01 12:00:00")
		assert(t, ref.StartByMonth(-1), "2023-12-01 00:00:00", "1月1日回退1月")
	})

	t.Run("允许溢出 (Overflow Flag)", func(t *testing.T) {
		ref := Parse("2024-04-15 12:00:00")
		// 传入 Overflow 之后允许自然溢出：
		assert(t, ref.StartByDay(Overflow, 45), "2024-05-30 00:00:00", "开启溢出的位移")
	})

	t.Run("Oracle 极端用例矩阵", func(t *testing.T) {
		// 1. 跨世纪的“量子坍缩”级联 (Century-Decade-Year Cascade)
		// 基准: 2024-02-29 (闰日)
		// StartByCentury(0,0,0,1,-1) -> Century(0)使y=2000 -> Month(1)使m=3 -> Day(-1)使d=28
		ref闰日 := Parse("2024-02-29 12:00:00")
		assert(t, ref闰日.StartByCentury(0, 0, 0, 1, -1), "2000-03-28 00:00:00", "跨世纪零位移坍缩")

		// 2. “溢出模式”全局压测 (Overflow Flag Natural Rollover)
		// 基准: 2024-01-31. Month(1) -> 03-02. Day(40) + Overflow -> 自然滑动到 4/11
		refJan31 := Parse("2024-01-31 12:00:00")
		assert(t, refJan31.StartByMonth(Overflow, 1, 40), "2024-04-11 00:00:00", "Overflow开启自然溢出")

		// 3. “时空倒流”的末端对齐 (EndByCentury with Multi-Negative Offsets)
		// 2024-05-15 (Q2). Quarter(-1) -> 1/15 -> final(end)使m=3 -> 3/15. Month(-1) -> 2/15. Day(-1) -> 2/14.
		refQ2 := Parse("2024-05-15 12:00:00")
		assert(t, refQ2.EndByQuarter(-1, -1, -1), "2024-02-14 23:59:59.999999999", "EndBy负向级联")

		// 4. “主权与自然的冲突” (ISO Week-Year Crossing)
		// 2026-01-01 (周四) 属于 ISO 2026-W01. StartByYearWeek(ISO, 0) -> 2025-12-29 (周一)
		refISO := Parse("2026-01-01 12:00:00")
		assert(t, refISO.StartByYearWeek(ISO, 0, 0), "2025-12-29 00:00:00", "ISO跨年周对齐")

		// 5. “世纪末的最后冲刺” (Extreme EndCentury-of-Century Cascade)
		// Century(0)使y=2099 -> Decade(-1)使y=2089 -> Year(-1)使y=2088 -> Month(-1)使y=2087, m=12
		refCentury := Parse("2024-01-01 12:00:00")
		assert(t, refCentury.EndByCentury(0, -1, -1, -1), "2087-12-31 23:59:59.999999999", "世纪末深度级联")
	})

	t.Run("验证大单位End模式的边界完整性", func(t *testing.T) {
		ref := Parse("2024-05-15 12:00:00")

		// 1. 验证：EndByQuarter(0) 是否能到达 Q2 的最后一天 (6月30日)
		assert(t, ref.EndByQuarter(0), "2024-06-30 23:59:59.999999999", "EndByQuarter(0) 应该到6月底")

		// 2. 验证：Decade 末尾
		gotDecade := ref.cascade(fromRel, true, Decade, 0)
		assert(t, gotDecade, "2029-12-31 23:59:59.999999999", "Decade(0)末尾应该到2029年底")

		// 3. 验证：Century 末尾
		gotCentury := ref.cascade(fromRel, true, Century, 0)
		assert(t, gotCentury, "2099-12-31 23:59:59.999999999", "Century(0)末尾应该到2099年底")
	})

	t.Run("纳秒精度相对位移 (Rel Nano)", func(t *testing.T) {
		// 基准: .000
		ref := Parse("2024-01-01 00:00:00")

		// 1. StartByMilli(1): +1ms -> .001. StartCentury(归零) -> .001000000
		assert(t, ref.StartByMilli(1), "2024-01-01 00:00:00.001", "StartByMilli(1)")

		// 2. EndByMicro(1): +1us -> .000001. EndCentury(置满) -> .000001999 (Micro只管后面纳秒置满)
		// Micro 对齐逻辑: ns = (ns/1e3)*1e3 + 999
		assert(t, ref.EndByMicro(1), "2024-01-01 00:00:00.000001999", "EndByMicro(1)")
	})
}
