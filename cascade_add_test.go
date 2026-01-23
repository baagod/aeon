package aeon

import (
    "testing"
)

func TestAdd(t *testing.T) {
    // 基准时间：2024-01-31 12:00:00 (1月最后一天)
    base := New(2024, 1, 31, 12, 0, 0)

    t.Run("Pure Translation (High Fidelity)", func(t *testing.T) {
        // 1. AddMonth(): 默认加 1 个月 -> 2月29日 (不回滚到 2/1)
        assert(t, base.AddMonth(), "2024-02-29 12:00:00", "AddMonth() 默认值与平移")

        // 2. AddDay(1): 纯平移
        assert(t, base.AddDay(1), "2024-02-01 12:00:00", "AddDay(1)")
    })

    t.Run("Cascading Add", func(t *testing.T) {
        // AddYear(1, 2, 3) = +1年 +2月 +3天
        // 2024-01-31 -> 2025-01-31 -> 2025-03-31 -> 2025-04-03
        assert(t, base.AddYear(1, 2, 3), "2025-04-03 12:00:00", "AddYear(1, 2, 3) 级联")
    })

    t.Run("Bitmask in Add", func(t *testing.T) {
        // 虽然 Add 系列推荐纯数字，但级联引擎依然支持位掩码
        // 开启 Overflow: 1月31日 + 1月 -> 3月2日
        assert(t, base.AddMonth(Overflow, 1), "2024-03-02 12:00:00", "AddMonth with Overflow")
    })
}
