package aeon

import (
    "testing"
)

func TestFunc(t *testing.T) {
    base := Parse("2024-04-15 12:00:00")
    assert(t, base.EndAtMonth(6, 5), "2024-06-20 23:59:59.999999999", "EndAtMonth(6, 5)")
}

// 测试 Pick 函数
func TestPick(t *testing.T) {
    base := Parse("2021-07-21 12:00:00")

    // 候选时间点
    t1 := Parse("2021-07-21 09:00:00") // 距离 -3h
    t2 := Parse("2021-07-21 10:00:00") // 距离 -2h
    t3 := Parse("2021-07-21 15:00:00") // 距离 +3h
    t4 := Parse("2021-07-21 18:00:00") // 距离 +6h

    // 1. Max/Min 模式 (替代 Maxmin)
    t.Run("Max (>)", func(t *testing.T) {
        assert(t, Pick('>', t1, t2, t3), "2021-07-21 15:00:00", "应该返回最晚时间")
        assert(t, Pick('>', t3, t2, t1), "2021-07-21 15:00:00", "乱序输入应该返回最晚时间")
    })

    t.Run("Min (<)", func(t *testing.T) {
        assert(t, Pick('<', t1, t2, t3), "2021-07-21 09:00:00", "应该返回最早时间")
        assert(t, Pick('<', t3, t2, t1), "2021-07-21 09:00:00", "乱序输入应该返回最早时间")
    })

    // 2. Near/Far 模式 (替代 Near)
    // 注意：Pick 的 Near/Far 模式下，第一个参数是 Reference (Base)
    t.Run("Near (-)", func(t *testing.T) {
        // 原 base.Near('<', t1, t2...) -> Pick('-', base, t1, t2...)
        // 预期：t2 (10:00, 距离2h)
        assert(t, Pick('-', base, t1, t2, t3, t4), "2021-07-21 10:00:00", "Pick(-) 应该返回距离最近的 10:00")
    })

    t.Run("Far (+)", func(t *testing.T) {
        // 原 base.Near('>', t1, t2...) -> Pick('+', base, t1, t2...)
        // 预期：t4 (18:00, 距离6h)
        assert(t, Pick('+', base, t1, t2, t3, t4), "2021-07-21 18:00:00", "Pick(+) 应该返回距离最远的 18:00")
    })

    // 3. 边界情况
    t.Run("Equal Distance", func(t *testing.T) {
        p1 := Parse("2021-07-21 11:00:00") // -1h
        p2 := Parse("2021-07-21 13:00:00") // +1h
        // Pick('-', base, p1, p2)
        assert(t, Pick('-', base, p1, p2), "2021-07-21 11:00:00", "距离相等时应返回第一个(11:00)")
        assert(t, Pick('-', base, p2, p1), "2021-07-21 13:00:00", "距离相等时应返回第一个(13:00)")
    })

    t.Run("Edge Cases", func(t *testing.T) {
        // Max/Min 模式单元素
        assert(t, Pick('>', t1), t1.String(), "Pick(>) 单元素返回自身")

        // Near/Far 模式单元素 (无候选者，只有 Reference)
        // Pick('+', base) -> 应该返回 base (安全降级)
        assert(t, Pick('+', base), base.String(), "Pick(+) 仅有参考点时返回参考点")

        // 零输入
        assert(t, Pick('>'), Time{}.String(), "Pick(>) 无参数返回零值")

        // 无效操作符
        assert(t, Pick('u', t1, t2), t1.String(), "Pick(u) 返回第一个元素")
    })
}

func assert(t *testing.T, b Time, e, f string) {
    t.Helper()
    if b.String() != e {
        t.Errorf("%s: got [%s], want [%s]", f, b, e)
    }
}

func assertZone(t *testing.T, actual Time, expectedOffset int, name string) {
    t.Helper()
    _, offset := actual.time.Zone()
    if offset != expectedOffset {
        t.Errorf("%s zone offset: got [%d], want [%d]", name, offset, expectedOffset)
    }
}
