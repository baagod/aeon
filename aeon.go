package aeon

import (
    "math"
    "time"
)

var (
    // DefaultWeekStarts 全局默认周起始日（默认为周一）
    DefaultWeekStarts = time.Monday
    // DefaultTimeZone Parse() 使用的默认时区
    DefaultTimeZone = time.Local
)

type Time struct {
    time       time.Time
    weekStarts time.Weekday
}

// --- 创建时间 ---

func Aeon(t ...time.Time) Time {
    if len(t) == 0 {
        t = []time.Time{{}}
    }
    return Time{time: t[0], weekStarts: DefaultWeekStarts}
}

func Now(loc ...*time.Location) Time {
    l := DefaultTimeZone
    if len(loc) > 0 && loc[0] != nil {
        l = loc[0]
    }
    return Time{time: time.Now().In(l), weekStarts: DefaultWeekStarts}
}

func New(y, m, d, h, mm, s int, add ...any) Time {
    l, ns := DefaultTimeZone, 0

    for _, arg := range add {
        switch v := arg.(type) {
        case int:
            switch {
            case v == 0:
                ns = 0
            case v <= 999: // 1-3：推断为毫秒
                ns = v * 1_000_000
            case v <= 999_999: // 4-6：推断为微秒
                ns = v * 1000
            default: // 7 及以上：推断为纳秒
                ns = v
            }
        case time.Duration:
            ns = int(v.Nanoseconds())
        default: // string, *time.Location
            l = timeZone(v)
        }
    }

    return Time{
        time:       time.Date(y, time.Month(m), d, h, mm, s, ns, l),
        weekStarts: DefaultWeekStarts,
    }
}

// WithWeekStarts 返回新实例，周起始日为 w。
func (t Time) WithWeekStarts(w time.Weekday) Time {
    return Time{time: t.time, weekStarts: w}
}

// Unix 返回给定时间戳的时间。secs 可以是秒、毫秒、微妙或纳秒级时间戳。
func Unix(secs int64, loc ...*time.Location) Time {
    v, t := secs, time.Time{}
    if secs < 0 { // 处理公元前时间戳
        v = -secs
    }

    switch {
    case v <= 9999999999: // 10位：秒
        t = time.Unix(secs, 0)
    case v <= 9999999999999: // 13位：毫秒
        t = time.UnixMilli(secs)
    case v <= 9999999999999999: // 16位：微秒
        t = time.UnixMicro(secs)
    default: // 19位及以上：纳秒
        t = time.Unix(0, secs)
    }

    if len(loc) > 0 && loc[0] != nil {
        t = t.In(loc[0])
    }

    return Aeon(t)
}

// --- 获取时间 ---

func (t Time) Year() int                 { return t.time.Year() }
func (t Time) Month() int                { return int(t.time.Month()) }
func (t Time) Day() int                  { return t.time.Day() }
func (t Time) Hour() int                 { return t.time.Hour() }
func (t Time) Minute() int               { return t.time.Minute() }
func (t Time) Clock() (h, mm, s int)     { return t.time.Clock() }
func (t Time) Days() int                 { return DaysIn(t.Year(), t.Month()) }
func (t Time) Weekday() time.Weekday     { return t.time.Weekday() }
func (t Time) ISOWeek() (year, week int) { return t.time.ISOWeek() }

// YearDay 返回一年中的第几天，非闰年范围 [1,365]，闰年范围 [1,366]。
func (t Time) YearDay() int {
    return t.time.YearDay()
}

// YearDays 返回本年总天数
func (t Time) YearDays() int { return DaysIn(t.Year()) }

// Date 返回 t 的年月日
func (t Time) Date() (int, int, int) {
    y, m, d := t.time.Date()
    return y, int(m), d
}

// Second 返回时间的秒数 (0-59)
func (t Time) Second() int { return t.time.Second() }

// Milli 返回毫秒数 (0-999)
func (t Time) Milli() int { return t.time.Nanosecond() / 1e6 }

// Micro 返回微秒数 (0-999999)
func (t Time) Micro() int { return t.time.Nanosecond() / 1e3 }

// Nano 返回纳秒数 (0-999999999)
func (t Time) Nano() int { return t.time.Nanosecond() }

// Unix 返回时间戳，可选择指定精度。
//
// 参数 n (可选) 指定返回的时间戳精度：
//   - 不提供或 0: 秒级 (10位)
//   - 3: 毫秒级 (13位)
//   - 6: 微秒级 (16位)
//   - 9: 纳秒级 (19位)
//   - 其他值: 对应位数的时间戳
func (t Time) Unix(n ...int) int64 {
    if len(n) == 0 || n[0] == 0 {
        return t.time.Unix()
    }
    precision := clamp(n[0]+10, 1, 19)
    return t.time.UnixNano() / pow19[19-precision]
}

// UTC 返回 UTC 时间
func (t Time) UTC() Time {
    return Time{time: t.time.UTC(), weekStarts: t.weekStarts}
}

// Local 返回本地时间
func (t Time) Local() Time {
    return Time{time: t.time.Local(), weekStarts: t.weekStarts}
}

// To 返回指定的 loc 时间
func (t Time) To(loc *time.Location) Time {
    return Time{time: t.time.In(loc), weekStarts: t.weekStarts}
}

// Round 返回距离当前时间最近的刻度点。
//
// 想象时间就是一把尺子，每隔 d 点设置刻度，返回与当前时间最近的刻度点时间。
//
// 示例：当前时间 t = 14:35:29，刻度 d = 15分钟
//
// 刻度尺：14:15 ──┬── 14:30 ── [t] ────┬── 14:45
//
// 此时 t (14:35) 距离刻度 14:30 更近，所以返回 14:30:00。
func (t Time) Round(d time.Duration) Time {
    return Time{time: t.time.Round(d), weekStarts: t.weekStarts}
}

// Truncate 与 Round 相同，但它永远截断在过去（而非未来）刻度的时间。
func (t Time) Truncate(d time.Duration) Time {
    return Time{time: t.time.Truncate(d), weekStarts: t.weekStarts}
}

// Time 返回 time.Time
func (t Time) Time() time.Time {
    return t.time
}

// Location 返回时区信息
func (t Time) Location() *time.Location {
    return t.time.Location()
}

// Zone 获取当前时区的名称和偏移量
func (t Time) Zone() (name string, offset int) {
    return t.time.Zone()
}

// ---- 比较时间 ----

// Diff 返回 t 和 u 的时间差。
//
// 参数：
//   - unit: 比较单位 ("y"年, "M"月, "d"日, "h"时, "m"分, "s"秒)
//   - abs: 可选，为 true 时返回绝对值
func (t Time) Diff(u Time, unit string, abs ...bool) float64 {
    var diff float64
    switch unit {
    case "y":
        // 年差 = 整数年差 + t 的年内进度 - u 的年内进度
        years := float64(t.Year() - u.Year())
        tDays := float64(t.YearDay()) / float64(t.YearDays()) // t 在本年的进度 (0~1)
        uDays := float64(u.YearDay()) / float64(u.YearDays()) // u 在本年的进度 (0~1)
        diff = years + tDays - uDays
    case "M":
        // 月差 = 整数月差 + 天数偏移比例
        months := float64((t.Year()-u.Year())*12 + t.Month() - u.Month())
        // 计算 [天差] 并将结果转换为浮点数：天差 = t.Day() - u.Day()
        days := float64(t.Day() - u.Day())
        // 如果天差为负数，表示未满一个完整的月，需要将月差减 1。
        if days < 0 {
            months-- // 天数不足一月，月差减 1
        }
        diff = months + days/float64(u.Days())
    case "d":
        diff = t.Sub(u).Hours() / 24
    case "h":
        diff = t.Sub(u).Hours()
    case "m":
        diff = t.Sub(u).Minutes()
    case "s":
        diff = t.Sub(u).Seconds()
    }

    if len(abs) > 0 && abs[0] {
        return math.Abs(diff)
    }

    return diff
}

// Sub 返回 t - u 的时间差
func (t Time) Sub(u Time) time.Duration {
    return t.time.Sub(u.time)
}

// Lt 返回 t 是否在 u 之前 (t < u)
func (t Time) Lt(u Time) bool {
    return t.time.Before(u.time)
}

// Gt 返回 t 是否在 u 之后 (t > u)
func (t Time) Gt(u Time) bool {
    return t.time.After(u.time)
}

// Eq 返回 t == u
func (t Time) Eq(u Time) bool {
    return t.time.Equal(u.time)
}

// Compare 比较 t 和 u。
//
// 如果 t == u，返回 0；t > u 返回 1；t < u 返回 -1。
func (t Time) Compare(u Time) int {
    return t.time.Compare(u.time)
}

// --- 其他 ---

// IsZero 返回 t 是否零时，即 0001-01-01 00:00:00 UTC。
func (t Time) IsZero() bool {
    return t.time.IsZero()
}

func (t Time) ZeroOr(u Time) Time {
    if t.IsZero() {
        return u
    }
    return t
}

// IsLeapYear 返回 t 是否闰年
func (t Time) IsLeapYear() bool {
    return IsLeapYear(t.Year())
}

// IsLongYear 返回 t 是否 ISO 长年
func (t Time) IsLongYear() bool {
    return IsLongYear(t.Year())
}

// IsAM 返回时间是否在上午 (00:00:00 ~ 11:59:59)
func (t Time) IsAM() bool { return t.Hour() < 12 }

// IsWeekend 返回时间是否在周末 (周六或周日)
func (t Time) IsWeekend() bool {
    w := t.Weekday()
    return w == time.Saturday || w == time.Sunday
}

// IsDST 返回 t 是否夏令时
func (t Time) IsDST() bool {
    return t.time.IsDST()
}

// IsSame 返回 t 与 target 是否在指定单位下相同 (包含上级单位一致性)
func (t Time) IsSame(u Unit, target Time) bool {
    switch u {
    case Century:
        return t.Year()/100 == target.Year()/100
    case Decade:
        return t.Year()/10 == target.Year()/10
    case Year:
        return t.Year() == target.Year()
    case Month:
        ty, tm, _ := t.Date()
        uy, um, _ := target.Date()
        return ty == uy && tm == um
    case Day:
        ty, tm, td := t.Date()
        uy, um, ud := target.Date()
        return ty == uy && tm == um && td == ud
    default:
        return a(t, seAbs, u).Eq(a(target, seAbs, u))
    }
}

// Until 返回 t 与当前时间 (Now) 的相对时长 (t - Now)。
//
// 返回值：
//   - 正数：t 在未来
//   - 负数：t 在过去
//   - 零：  t 是现在
func (t Time) Until() time.Duration {
    return time.Until(t.time)
}

// Near 返回在集合中距离 t 最近 ("<") 或 最远 (">") 的时间
//
// 类似于离散版的 Round。
func (t Time) Near(op string, times ...Time) Time {
    isGt, isLt := op == ">", op == "<"
    if len(times) == 0 || !isGt && !isLt {
        return t
    }

    // 初始化基准：默认第一个元素为当前最优解，并计算其与 t 的绝对距离
    res := times[0]
    best := t.Sub(res).Abs()

    for _, x := range times[1:] {
        d := t.Sub(x).Abs() // 计算当前元素与 t 的绝对距离
        if (isGt && d > best) || (isLt && d < best) {
            res, best = x, d
        }
    }

    return res
}

// Between 判断 t 是否在 (start, end) 区间内。
//
// 可选参数 bounds 用于控制边界包含性 (默认为 "=")：
//   - "=" : 包含边界
//   - "!" : 不包含边界
//   - "[" : 包含左边界
//   - "]" : 包含右边界
func (t Time) Between(start, end Time, bounds ...string) bool {
    b := "=" // 默认包含边界
    if len(bounds) > 0 {
        b = bounds[0]
    }
    switch b {
    case "!": // 全排除
        return t.Gt(start) && t.Lt(end)
    case "[": // 仅左含
        return (t.Gt(start) || t.Eq(start)) && t.Lt(end)
    case "]": // 仅右含
        return t.Gt(start) && (t.Lt(end) || t.Eq(end))
    default: // 全包含 (默认)
        return (t.Gt(start) || t.Eq(start)) && (t.Lt(end) || t.Eq(end))
    }
}

// --- 时间格式 ---

func (t Time) Format(layout string) string                 { return t.time.Format(layout) }
func (t Time) AppendFormat(b []byte, layout string) []byte { return t.time.AppendFormat(b, layout) }
func (t Time) String() string                              { return t.time.Format(DTNs) }

func (t Time) ToString(f ...string) string {
    if len(f) > 0 {
        return t.time.Format(f[0])
    }
    return t.time.Format(DTNs)
}

// --- Aeon 包方法 ---

// Maxmin 在一组时间中返回 最大 (">") 或 最小值 ("<")
//
//  - 如果 `op` 不是 ">" 或 "<"，返回集合中的第一个时间。
//  - 如果集合为空，返回零值 Time。
func Maxmin(op string, times ...Time) Time {
    if len(times) == 0 {
        return Time{}
    }

    res := times[0]
    isGt, isLt := op == ">", op == "<"

    if !isGt && !isLt {
        return res
    }

    for _, x := range times[1:] {
        if isGt && x.Gt(res) || isLt && x.Lt(res) {
            res = x
        }
    }

    return res
}

// IsLeapYear 返回 y 是否闰年
func IsLeapYear(y int) bool {
    return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}

// IsLongYear 返回当前年份是否为 ISO 8601 规定的 “长年”（包含 53 周）。
func IsLongYear(y int) bool {
    // 获取该年 1月1日 是周几 (0=Sun, 4=Thu)
    // 这里用我们内置的 weekday 函数，比 time.Date 快得多
    w := weekday(y, 1, 1)

    // ISO 8601 长年逻辑：
    // 1. 1月1日是周四
    // 2. 或者是闰年且1月1日是周三
    return w == time.Thursday || (w == time.Wednesday && IsLeapYear(y))
}

// DaysIn 返回 y 年 m 月最大天数，如果忽略 m 则返回 y 年总天数。
//
//   - 1, 3, 5, 7, 8, 10, 12 月有 31 天；4, 6, 9, 11 月有 30 天。
//   - 平年 2 月有 28 天，闰年 29 天。
func DaysIn(y int, m ...int) int {
    if len(m) > 0 {
        if m[0] == 2 && IsLeapYear(y) {
            return 29
        }
        return maxDays[m[0]]
    }

    if IsLeapYear(y) {
        return 366
    }

    return 365
}
