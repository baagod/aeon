package aeon

import (
    "testing"
)

func TestYearsNavigation(t *testing.T) {
    // 2021-02-02 13:14:15
    base := New(2021, 2, 2, 13, 14, 15)

    // --- Century ---
    assert(t, base.StartCentury(0), "2000-01-01 00:00:00", "StartCentury(0)")
    assert(t, base.StartCentury(1), "2100-01-01 00:00:00", "StartCentury(1)")
    assert(t, base.StartCentury(2), "2200-01-01 00:00:00", "StartCentury(2)")
    assert(t, base.StartCentury(-1), "2900-01-01 00:00:00", "StartCentury(-1)")

    assert(t, base.EndCentury(0), "2099-12-31 23:59:59.999999999", "EndCentury(0)")
    assert(t, base.EndCentury(1), "2199-12-31 23:59:59.999999999", "EndCentury(1)")
    assert(t, base.EndCentury(2), "2299-12-31 23:59:59.999999999", "EndCentury(2)")
    assert(t, base.EndCentury(-1), "2999-12-31 23:59:59.999999999", "EndCentury(-1)")

    assert(t, base.GoCentury(0), "2021-02-02 13:14:15", "GoCentury(0)")
    assert(t, base.GoCentury(1), "2121-02-02 13:14:15", "GoCentury(1)")
    assert(t, base.GoCentury(2), "2221-02-02 13:14:15", "GoCentury(2)")
    assert(t, base.GoCentury(-1), "2921-02-02 13:14:15", "GoCentury(-1)")

    // --- Decade ---
    assert(t, base.StartDecade(0), "2020-01-01 00:00:00", "StartDecade(0)")
    assert(t, base.StartDecade(1), "2010-01-01 00:00:00", "StartDecade(1)")
    assert(t, base.StartDecade(2), "2020-01-01 00:00:00", "StartDecade(2)")
    assert(t, base.StartDecade(-1), "2090-01-01 00:00:00", "StartDecade(-1)")

    assert(t, base.EndDecade(0), "2029-12-31 23:59:59.999999999", "EndDecade(0)")
    assert(t, base.EndDecade(1), "2019-12-31 23:59:59.999999999", "EndDecade(1)")
    assert(t, base.EndDecade(2), "2029-12-31 23:59:59.999999999", "EndDecade(2)")
    assert(t, base.EndDecade(-1), "2099-12-31 23:59:59.999999999", "EndDecade(-1)")

    assert(t, base.GoDecade(0), "2001-02-02 13:14:15", "GoDecade(0)")
    assert(t, base.GoDecade(1), "2011-02-02 13:14:15", "GoDecade(1)")
    assert(t, base.GoDecade(2), "2021-02-02 13:14:15", "GoDecade(2)")
    assert(t, base.GoDecade(-1), "2091-02-02 13:14:15", "GoDecade(-1)")

    // --- Year ---
    assert(t, base.StartYear(0), "2021-01-01 00:00:00", "StartYear(0)")
    assert(t, base.StartYear(1), "2021-01-01 00:00:00", "StartYear(1)")
    assert(t, base.StartYear(2), "2022-01-01 00:00:00", "StartYear(2)")
    assert(t, base.StartYear(-1), "2029-01-01 00:00:00", "StartYear(-1)")

    assert(t, base.EndYear(0), "2021-12-31 23:59:59.999999999", "EndYear(0)")
    assert(t, base.EndYear(1), "2021-12-31 23:59:59.999999999", "EndYear(1)")
    assert(t, base.EndYear(2), "2022-12-31 23:59:59.999999999", "EndYear(2)")
    assert(t, base.EndYear(-1), "2029-12-31 23:59:59.999999999", "EndYear(-1)")

    assert(t, base.GoYear(0), "2020-02-02 13:14:15", "GoYear(0)")
    assert(t, base.GoYear(1), "2021-02-02 13:14:15", "GoYear(1)")
    assert(t, base.GoYear(2), "2022-02-02 13:14:15", "GoYear(2)")
    assert(t, base.GoYear(-1), "2029-02-02 13:14:15", "GoYear(-1)")
}

func TestWeeksNavigation(t *testing.T) {
    // 基准日期：2026-01-21(周三)，1月31天 = 4×7+3，序数周期：1-7, 8-14, 15-21, 22-28, 29-31(不完整)
    base := New(2026, 1, 21, 13, 14, 15)

    /* 2026-01
       一    二   三   四    五   六    日
       ---  ---  ---  ---  ---  ---  ---
       29   30   31   01   02   03   04  <-- {1月开始于周四}
       05   06   07   08   09   10   11
       12   13   14   15   16   17   18
       19   20  [21]  22   23   24   25  <-- [当前天: 周三]
       26   27   28   29   30  {31}  01  <-- {1最后一天: 周五}
       02  (03)  04   05   06   07   08  <-- (自然溢出落点)
    */

    t.Run("Week_Calendar", func(t *testing.T) {
        assert(t, base.GoWeek(0), "2026-01-21 13:14:15", "GoWeek(0)")
        assert(t, base.GoWeek(1), "2025-12-31 13:14:15", "GoWeek(1)")
        assert(t, base.GoWeek(2), "2026-01-07 13:14:15", "GoWeek(2)")
        assert(t, base.GoWeek(-1), "2026-01-28 13:14:15", "GoWeek(-1)")
        assert(t, base.GoWeek(-2), "2026-01-21 13:14:15", "GoWeek(-2)")

        assert(t, base.StartWeek(0), "2026-01-19 00:00:00", "StartWeek(0)")
        assert(t, base.StartWeek(1), "2025-12-29 00:00:00", "StartWeek(1)")
        assert(t, base.StartWeek(2), "2026-01-05 00:00:00", "StartWeek(2)")
        assert(t, base.StartWeek(-1), "2026-01-26 00:00:00", "StartWeek(-1)")
        assert(t, base.StartWeek(-2), "2026-01-19 00:00:00", "StartWeek(-2)")

        assert(t, base.EndWeek(0), "2026-01-25 23:59:59.999999999", "EndWeek(0)")
        assert(t, base.EndWeek(1), "2026-01-04 23:59:59.999999999", "EndWeek(1)")
        assert(t, base.EndWeek(2), "2026-01-11 23:59:59.999999999", "EndWeek(2)")
        assert(t, base.EndWeek(-1), "2026-02-01 23:59:59.999999999", "EndWeek(-1)")
        assert(t, base.EndWeek(-2), "2026-01-25 23:59:59.999999999", "EndWeek(-2)")
    })

    t.Run("Week_Full", func(t *testing.T) {
        assert(t, base.GoWeek(Full, 0), "2026-01-21 13:14:15", "GoWeek(Full, 0)")
        assert(t, base.GoWeek(Full, 1), "2026-01-07 13:14:15", "GoWeek(Full, 1)")
        assert(t, base.GoWeek(Full, 2), "2026-01-14 13:14:15", "GoWeek(Full, 2)")
        assert(t, base.GoWeek(Full, -1), "2026-01-28 13:14:15", "GoWeek(Full, -1)")
        assert(t, base.GoWeek(Full, -2), "2026-01-21 13:14:15", "GoWeek(Full, -2)")

        assert(t, base.StartWeek(Full, 0), "2026-01-19 00:00:00", "StartWeek(Full, 0)")
        assert(t, base.StartWeek(Full, 1), "2026-01-05 00:00:00", "StartWeek(Full, 1)")
        assert(t, base.StartWeek(Full, 2), "2026-01-12 00:00:00", "StartWeek(Full, 2)")
        assert(t, base.StartWeek(Full, -1), "2026-01-26 00:00:00", "StartWeek(Full, -1)")
        assert(t, base.StartWeek(Full, -2), "2026-01-19 00:00:00", "StartWeek(Full, -2)")

        assert(t, base.EndWeek(Full, 0), "2026-01-25 23:59:59.999999999", "EndWeek(Full, 0)")
        assert(t, base.EndWeek(Full, 1), "2026-01-11 23:59:59.999999999", "EndWeek(Full, 1)")
        assert(t, base.EndWeek(Full, 2), "2026-01-18 23:59:59.999999999", "EndWeek(Full, 2)")
        assert(t, base.EndWeek(Full, -1), "2026-02-01 23:59:59.999999999", "EndWeek(Full, -1)")
        assert(t, base.EndWeek(Full, -2), "2026-01-25 23:59:59.999999999", "EndWeek(Full, -2)")
    })

    t.Run("Week_ISO", func(t *testing.T) {
        assert(t, base.GoWeek(ISO, 0), "2026-01-21 13:14:15", "GoWeek(ISO, 0) 当前时间")
        assert(t, base.GoWeek(ISO, 1), "2025-12-31 13:14:15", "GoWeek(ISO, 1) 本年第 1 个 ISO 年周")
        assert(t, base.GoWeek(ISO, 2), "2026-01-07 13:14:15", "GoWeek(ISO, 2) 本年第 2 个 ISO 年周")
        assert(t, base.GoWeek(ISO, -1), "2026-12-30 13:14:15", "GoWeek(ISO, -1) 本年最后 1 个 ISO 年周")
        assert(t, base.GoWeek(ISO, -2), "2026-12-23 13:14:15", "GoWeek(ISO, -2) 本年倒数第 2 个 ISO 年周")

        assert(t, base.StartWeek(ISO, 0), "2026-01-19 00:00:00", "StartWeek(ISO, 0)")
        assert(t, base.StartWeek(ISO, 1), "2025-12-29 00:00:00", "StartWeek(ISO, 1)")
        assert(t, base.StartWeek(ISO, 2), "2026-01-05 00:00:00", "StartWeek(ISO, 2)")
        assert(t, base.StartWeek(ISO, -1), "2026-12-28 00:00:00", "StartWeek(ISO, -1)")
        assert(t, base.StartWeek(ISO, -2), "2026-12-21 00:00:00", "StartWeek(ISO, -2)")

        assert(t, base.EndWeek(ISO, 0), "2026-01-25 23:59:59.999999999", "EndWeek(ISO, 0)")
        assert(t, base.EndWeek(ISO, 1), "2026-01-04 23:59:59.999999999", "EndWeek(ISO, 1)")
        assert(t, base.EndWeek(ISO, 2), "2026-01-11 23:59:59.999999999", "EndWeek(ISO, 2)")
        assert(t, base.EndWeek(ISO, -1), "2027-01-03 23:59:59.999999999", "EndWeek(ISO, -1)")
        assert(t, base.EndWeek(ISO, -2), "2026-12-27 23:59:59.999999999", "EndWeek(ISO, -2)")
    })

    t.Run("Week_Ord", func(t *testing.T) {
        assert(t, base.GoWeek(Ord, 0), "2026-01-21 13:14:15", "GoWeek(Ord, 0) 当前时间")
        assert(t, base.GoWeek(Ord, 1), "2026-01-01 13:14:15", "GoWeek(Ord, 1) 本月 1 日")
        assert(t, base.GoWeek(Ord, 2), "2026-01-08 13:14:15", "GoWeek(Ord, 2) 本月第 2 周开始")
        assert(t, base.GoWeek(Ord, -1), "2026-01-31 13:14:15", "GoWeek(Ord, -1) 本月最后 1 天")
        assert(t, base.GoWeek(Ord, -2), "2026-01-24 13:14:15", "GoWeek(Ord, -2) 本月倒数第 2 周结束")

        assert(t, base.StartWeek(Ord, 0), "2026-01-19 00:00:00", "StartWeek(Ord, 0) 本周期起点")
        assert(t, base.StartWeek(Ord, 1), "2026-01-01 00:00:00", "StartWeek(Ord, 1) 本月 1 日")
        assert(t, base.StartWeek(Ord, 2), "2026-01-08 00:00:00", "StartWeek(Ord, 2) 本月第 2 周开始")
        assert(t, base.StartWeek(Ord, -1), "2026-01-25 00:00:00", "StartWeek(Ord, -1) 本月最后周开始")
        assert(t, base.StartWeek(Ord, -2), "2026-01-18 00:00:00", "StartWeek(Ord, -2) 本月倒数第 2 周开始")

        assert(t, base.EndWeek(Ord, 0), "2026-01-25 23:59:59.999999999", "EndWeek(Ord, 0) 本周结束")
        assert(t, base.EndWeek(Ord, 1), "2026-01-07 23:59:59.999999999", "EndWeek(Ord, 1) 本月第 1 周结束")
        assert(t, base.EndWeek(Ord, 2), "2026-01-14 23:59:59.999999999", "EndWeek(Ord, 2) 本月第 2 周结束")
        assert(t, base.EndWeek(Ord, -1), "2026-01-31 23:59:59.999999999", "EndWeek(Ord, -1) 本月最后 1 天")
        assert(t, base.EndWeek(Ord, -2), "2026-01-24 23:59:59.999999999", "EndWeek(Ord, -2) 本月倒数第 2 周结束")
    })

    t.Run("Week_Qtr", func(t *testing.T) {
        assert(t, base.GoWeek(Qtr, 0), "2026-01-21 13:14:15", "GoWeek(Qtr, 0) 当前时间")
        assert(t, base.GoWeek(Qtr, 1), "2026-01-07 13:14:15", "GoWeek(Qtr, 1) 本季度第 1 周")
        assert(t, base.GoWeek(Qtr, 2), "2026-01-14 13:14:15", "GoWeek(Qtr, 2) 本季度第 2 周")
        assert(t, base.GoWeek(Qtr, -1), "2026-04-01 13:14:15", "GoWeek(Qtr, -1) 本季度最后 1 周")
        assert(t, base.GoWeek(Qtr, -2), "2026-03-25 13:14:15", "GoWeek(Qtr, -2) 本季度倒数第 2 周")

        assert(t, base.StartWeek(Qtr, 0), "2026-01-19 00:00:00", "StartWeek(Qtr, 0) 本周期起点")
        assert(t, base.StartWeek(Qtr, 1), "2026-01-01 00:00:00", "StartWeek(Qtr, 1) 本季度第 1 周开始")
        assert(t, base.StartWeek(Qtr, 2), "2026-01-08 00:00:00", "StartWeek(Qtr, 2) 本季度第 2 周开始")
        assert(t, base.StartWeek(Qtr, -1), "2026-03-30 00:00:00", "StartWeek(Qtr, -1) 本季度最后 1 周开始")
        assert(t, base.StartWeek(Qtr, -2), "2026-03-23 00:00:00", "StartWeek(Qtr, -2) 本季度倒数第 2 周开始")

        assert(t, base.EndWeek(Qtr, 0), "2026-01-25 23:59:59.999999999", "EndWeek(Qtr, 0) 本周期终点")
        assert(t, base.EndWeek(Qtr, 1), "2026-01-07 23:59:59.999999999", "EndWeek(Qtr, 1) 本季度第 1 周结束")
        assert(t, base.EndWeek(Qtr, 2), "2026-01-14 23:59:59.999999999", "EndWeek(Qtr, 2) 本季度第 2 周结束")
        assert(t, base.EndWeek(Qtr, -1), "2026-04-05 23:59:59.999999999", "EndWeek(Qtr, -1) 本季度最后 1 周结束")
        assert(t, base.EndWeek(Qtr, -2), "2026-03-29 23:59:59.999999999", "EndWeek(Qtr, -2) 本季度倒数第 2 周结束")
    })

    t.Run("Week_Qtr_Ord", func(t *testing.T) {
        assert(t, base.GoWeek(Qtr|Ord, 0), "2026-01-21 13:14:15", "GoWeek(Qtr|Ord, 0) Current")
        assert(t, base.GoWeek(Qtr|Ord, 1), "2026-01-01 13:14:15", "GoWeek(Qtr|Ord, 1) -> 本季度第 1 个序数周")
        assert(t, base.GoWeek(Qtr|Ord, 2), "2026-01-08 13:14:15", "GoWeek(Qtr|Ord, 2) -> 本季度第 2 个序数周")
        assert(t, base.GoWeek(Qtr|Ord, -1), "2026-03-31 13:14:15", "GoWeek(Qtr|Ord, -1) -> 本季度最后 1 天")
        assert(t, base.GoWeek(Qtr|Ord, -2), "2026-03-24 13:14:15", "GoWeek(Qtr|Ord, -2) -> 从本季度末尾开始倒数的第 2 个序数周")

        assert(t, base.StartWeek(Qtr|Ord, 0), "2026-01-19 00:00:00", "StartWeek(Qtr|Ord, 0)")
        assert(t, base.StartWeek(Qtr|Ord, 1), "2026-01-01 00:00:00", "StartWeek(Qtr|Ord, 1)")
        assert(t, base.StartWeek(Qtr|Ord, 2), "2026-01-08 00:00:00", "StartWeek(Qtr|Ord, 2)")
        assert(t, base.StartWeek(Qtr|Ord, -1), "2026-03-25 00:00:00", "StartWeek(Qtr|Ord, -1)")
        assert(t, base.StartWeek(Qtr|Ord, -2), "2026-03-18 00:00:00", "StartWeek(Qtr|Ord, -2)")

        assert(t, base.EndWeek(Qtr|Ord, 0), "2026-01-25 23:59:59.999999999", "EndWeek(Qtr|Ord, 0)")
        assert(t, base.EndWeek(Qtr|Ord, 1), "2026-01-07 23:59:59.999999999", "EndWeek(Qtr|Ord, 1)")
        assert(t, base.EndWeek(Qtr|Ord, 2), "2026-01-14 23:59:59.999999999", "EndWeek(Qtr|Ord, 2)")
        assert(t, base.EndWeek(Qtr|Ord, -1), "2026-03-31 23:59:59.999999999", "EndWeek(Qtr|Ord, -1)")
        assert(t, base.EndWeek(Qtr|Ord, -2), "2026-03-24 23:59:59.999999999", "EndWeek(Qtr|Ord, -2)")
    })

    t.Run("Week_Ord_Day", func(t *testing.T) {
        assert(t, base.GoWeek(Ord, 0, 0), "2026-01-21 13:14:15", "base.GoWeek(Ord, 0, 0) 本周天")
        assert(t, base.GoWeek(Ord, 0, 1), "2026-01-19 13:14:15", "GoWeek(Ord, 0, 1) 本周一")
        assert(t, base.GoWeek(Ord, 0, 2), "2026-01-20 13:14:15", "GoWeek(Ord, 0, 2) 本周二")
        assert(t, base.GoWeek(Ord, 0, -1), "2026-01-25 13:14:15", "GoWeek(Ord, 0, -1) 本周日")
        assert(t, base.GoWeek(Ord, 0, -2), "2026-01-24 13:14:15", "GoWeek(Ord, 0, -2) 本周六")

        assert(t, base.GoWeek(Ord, 1, 0), "2026-01-01 13:14:15", "第1周期 保持周四 (锚点原地)")
        assert(t, base.GoWeek(Ord, 1, 1), "2026-01-05 13:14:15", "第1周期 找周一 (向后推进)")
        assert(t, base.GoWeek(Ord, 1, 2), "2026-01-06 13:14:15", "第1周期 找周二 (向后推进)")
        assert(t, base.GoWeek(Ord, 1, -1), "2026-01-04 13:14:15", "第1周期 找周日 (负数映射 -> 向后推进)")
        assert(t, base.GoWeek(Ord, 1, -2), "2026-01-03 13:14:15", "第1周期 找周六 (负数映射 -> 向后推进)")

        assert(t, base.GoWeek(Ord, 2, 0), "2026-01-08 13:14:15", "第2周期 保持周四 (锚点原地)")
        assert(t, base.GoWeek(Ord, 2, 1), "2026-01-12 13:14:15", "第2周期 找周一 (向后推进)")
        assert(t, base.GoWeek(Ord, 2, 2), "2026-01-13 13:14:15", "第2周期 找周二 (向后推进)")
        assert(t, base.GoWeek(Ord, 2, -1), "2026-01-11 13:14:15", "第2周期 找周日 (负数映射 -> 向后推进)")
        assert(t, base.GoWeek(Ord, 2, -2), "2026-01-10 13:14:15", "第2周期 找周六 (负数映射 -> 向后推进)")

        assert(t, base.GoWeek(Ord, -1, 0), "2026-01-31 13:14:15", "最后周期 保持周六 (锚点原地)")
        assert(t, base.GoWeek(Ord, -1, 1), "2026-01-26 13:14:15", "最后周期 找周一 (向前回溯)")
        assert(t, base.GoWeek(Ord, -1, 2), "2026-01-27 13:14:15", "最后周期 找周二 (向前回溯)")
        assert(t, base.GoWeek(Ord, -1, -1), "2026-01-25 13:14:15", "最后周期 找周日 (负数映射 -> 向前回溯)")
        assert(t, base.GoWeek(Ord, -1, -2), "2026-01-31 13:14:15", "最后周期 找周六 (负数映射 -> 命中锚点)")

        assert(t, base.GoWeek(Ord, -2, 0), "2026-01-24 13:14:15", "倒数第2周期 保持周六 (锚点原地)")
        assert(t, base.GoWeek(Ord, -2, 1), "2026-01-19 13:14:15", "倒数第2周期 找周一 (向前回溯)")
        assert(t, base.GoWeek(Ord, -2, 2), "2026-01-20 13:14:15", "倒数第2周期 找周二 (向前回溯)")
        assert(t, base.GoWeek(Ord, -2, -1), "2026-01-18 13:14:15", "倒数第2周期 找周日 (负数映射 -> 向前回溯)")
        assert(t, base.GoWeek(Ord, -2, -2), "2026-01-24 13:14:15", "倒数第2周期 找周六 (负数映射 -> 命中锚点)")
    })
}

func TestOrdWeekBoundary(t *testing.T) {
    // 2026-01-31 是周六 (6)
    endBase := New(2026, 1, 31, 13, 14, 15)
    // 2026-01-01 是周四 (4)
    startBase := New(2026, 1, 1, 13, 14, 15)

    t.Run("月末原地出发边界测试", func(t *testing.T) {
        assert(t, endBase.GoWeek(Ord, -1, 1), "2026-01-26 13:14:15", "最后一周 找周一 (向前回溯)")
        assert(t, endBase.GoWeek(Ord, -1, 6), "2026-01-31 13:14:15", "最后一周 找周六 (原地命中)")
        assert(t, endBase.GoWeek(Ord, -1, -1), "2026-01-25 13:14:15", "最后一周 找周日 (负数映射)")
    })

    t.Run("月初原地出发边界测试", func(t *testing.T) {
        assert(t, startBase.GoWeek(Ord, 1, 4), "2026-01-01 13:14:15", "第一周 找周四 (原地命中)")
        assert(t, startBase.GoWeek(Ord, 1, 1), "2026-01-05 13:14:15", "第一周 找周一 (向后推进)")
        assert(t, startBase.GoWeek(Ord, 1, -1), "2026-01-04 13:14:15", "第一周 找周日 (负数映射)")
    })
}
