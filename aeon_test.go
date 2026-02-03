package aeon

import (
    "testing"
)

func TestFunc(_ *testing.T) {
    // t := Parse("2021-07-21 07:00:00")
}

func TestNear(t *testing.T) {
    base := Parse("2021-07-21 12:00:00")

    // 候选时间点
    t1 := Parse("2021-07-21 09:00:00") // 距离 -3h
    t2 := Parse("2021-07-21 10:00:00") // 距离 -2h
    t3 := Parse("2021-07-21 15:00:00") // 距离 +3h
    t4 := Parse("2021-07-21 18:00:00") // 距离 +6h

    t.Run("Nearest (<)", func(t *testing.T) {
        // 预期：t2 (10:00, 距离2h)
        assert(t, base.Near("<", t1, t2, t3, t4), "2021-07-21 10:00:00", "应该返回距离最近的 10:00")
    })

    t.Run("Furthest (>)", func(t *testing.T) {
        // 预期：t4 (18:00, 距离6h)
        assert(t, base.Near(">", t1, t2, t3, t4), "2021-07-21 18:00:00", "应该返回距离最远的 18:00")
    })

    t.Run("Equal Distance", func(t *testing.T) {
        // 相同距离，返回先出现的
        p1 := Parse("2021-07-21 11:00:00") // -1h
        p2 := Parse("2021-07-21 13:00:00") // +1h

        assert(t, base.Near("<", p1, p2), "2021-07-21 11:00:00", "距离相等时应返回第一个(11:00)")
        assert(t, base.Near("<", p2, p1), "2021-07-21 13:00:00", "距离相等时应返回第一个(13:00)")
    })

    t.Run("Edge Cases", func(t *testing.T) {
        assert(t, base.Near("<"), base.String(), "无候选应返回自身")
        assert(t, base.Near("", t1), base.String(), "无操作符应返回自身")
    })

    t.Run("Zero Distance (Exact Match)", func(t *testing.T) {
        // 列表中包含自身
        assert(t, base.Near("<", t1, base), base.String(), "存在自身时，最近的应该是自身(距离0)")
        // 注意：如果存在距离为0的点，找最远点不应该受影响，依然找距离最大的
        assert(t, base.Near(">", t1, base), t1.String(), "存在自身时，最远的应该是距离最大的点")
    })

    t.Run("Single Candidate", func(t *testing.T) {
        assert(t, base.Near("<", t1), t1.String(), "单元素应直接返回(Nearest)")
        assert(t, base.Near(">", t1), t1.String(), "单元素应直接返回(Furthest)")
    })

    t.Run("Invalid Operator", func(t *testing.T) {
        // 验证当前行为：未知操作符返回第一个元素 (因为不满足任何更新条件)
        assert(t, base.Near("unknown", t1, t2), t1.String(), "未知操作符应返回第一个候选者")
    })
}

func assert(t *testing.T, actual Time, expected string, name string) {
    t.Helper()
    if actual.String() != expected {
        t.Errorf("%s: got [%s], want [%s]", name, actual, expected)
    }
}

func assertZone(t *testing.T, actual Time, expectedOffset int, name string) {
    t.Helper()
    _, offset := actual.time.Zone()
    if offset != expectedOffset {
        t.Errorf("%s zone offset: got [%d], want [%d]", name, offset, expectedOffset)
    }
}
