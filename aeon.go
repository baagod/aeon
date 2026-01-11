package aeon

import (
	"math"
	"time"
)

type Unit int

const (
	Century Unit = iota
	Decade
	Year
	Month
	Day
	Hour
	Minute
	Second
	Millisecond
	Microsecond
	Nanosecond
	Quarter
	Week
	YearWeek
	Weekday
)

var (
	// DefaultWeekStarts 全局默认周起始日（默认为周一）
	DefaultWeekStarts = time.Monday
	// DefaultTimeZone Parse() 使用的默认时区
	DefaultTimeZone = time.Local
	// pow10 预定义的 10 的幂次方表，用于高精度计算
	pow10 = [...]int64{
		1, 10, 100, 1000, 10000, 100000, 1e6, 1e7, 1e8, 1e9,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
	}
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

func New(y, m, d, h, mm, s int, others ...any) Time {
	loc, ns := DefaultTimeZone, 0

	for _, arg := range others {
		switch v := arg.(type) {
		case time.Duration:
			ns = int(v.Nanoseconds())
		case *time.Location:
			loc = v
		}
	}

	return Time{
		time:       time.Date(y, time.Month(m), d, h, mm, s, ns, loc),
		weekStarts: DefaultWeekStarts,
	}
}

func NewDate(y, m, d int, others ...any) Time {
	return New(y, m, d, 0, 0, 0, others...)
}

func NewHour(y, m, d, h int, others ...any) Time {
	return New(y, m, d, h, 0, 0, others...)
}

func NewMinute(y, m, d, h, mm int, others ...any) Time {
	return New(y, m, d, h, mm, 0, others...)
}

// WithWeekStarts 返回新实例，周起始日为 w。
func (t Time) WithWeekStarts(w time.Weekday) Time {
	return Time{time: t.time, weekStarts: w}
}

// Unix 返回给定时间戳的时间。secs 可以是秒、毫秒、微妙或纳秒级时间戳。
func Unix(secs int64, utc ...bool) Time {
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

	if len(utc) > 0 && utc[0] {
		t = t.UTC()
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
	return t.time.UnixNano() / pow10[19-precision]
}

// UTC 返回 UTC 时间
func (t Time) UTC() Time {
	return Time{time: t.time.UTC(), weekStarts: t.weekStarts}
}

// Local 返回本地时间
func (t Time) Local() Time {
	return Time{time: t.time.Local(), weekStarts: t.weekStarts}
}

// In 返回指定的 loc 时间
func (t Time) In(loc *time.Location) Time {
	return Time{time: t.time.In(loc), weekStarts: t.weekStarts}
}

// Round 返回距离当前时间最近的 "跃点"。
//
// 这个函数有点不好理解，让我来尝试解释一下：
// 想象在当前时间之外存在着另一个时间循环，时间轴每次移动 d，而每个 d 就是一个 "跃点"。
// 函数返回距离当前时间最近的那个 "跃点"，如果当前时间正好位于两个 "跃点" 中间，返回指向未来的那个 "跃点"。
//
// 示例：假设当前时间是 2021-07-21 14:35:29.650
//
//   - 舍入到秒：t.Round(time.Second)
//
//     根据定义，时间 14:35:29.650 距离下一个跃点 14:35:30.000 最近 (只需要 350ms)，
//     而距离上一个跃点 14:35:29.000 较远 (需要 650ms)，故返回下一个跃点 14:35:30.000。
//
//   - 舍入分钟：t.Round(time.Minute)
//
//     时间 14:35:29.650 距离上一个跃点 14:35:00 最近（只需要 29.650s），
//     而距离下一个跃点 14:36:00 较远 (需要 30.350s)，故返回上一个跃点 14:35:00.000。
//
//   - 舍入 15 分钟：t.Round(15 * time.Minute)
//
//     跃点：--- 14:00:00 --- 14:15:00 --- 14:30:00 -- t --- 14:45:00 ---
//
//     时间 14:35:29.650 处在 14:30 (上一个跃点) 和 14:45 (下一个跃点) 之间，
//     距离上一个跃点最近，故返回上一个跃点时间：14:30:00。
func (t Time) Round(d time.Duration) Time {
	return Time{time: t.time.Round(d), weekStarts: t.weekStarts}
}

// Truncate 返回最接近当前时间但不超过它的 "跃点"（向过去截断）。
//
// 与 Round 的区别：
//   - Round 会选择最近的跃点（可能向未来舍入）
//   - Truncate 永远向过去截断，不进行四舍五入
//
// 可视化理解：
//
//	时间轴: ---[d1]---- t ----[d2]----[d3]----
//	                   ↑
//	                当前时间 t
//
// 无论 t 距离 d1 还是 d2 更近，Truncate 都只会返回 d1（向过去）。
// 如果 t 正好落在某个跃点上（如 d2），则返回 t 本身。
//
// 示例：假设当前时间是 2021-07-21 14:35:29.650
//
//   - 截断到秒：t.Truncate(time.Second)
//
//     直接舍弃毫秒部分，返回 14:35:29.000（向过去）。
//
//   - 截断到分钟：t.Truncate(time.Minute)
//
//     舍弃秒和毫秒部分，返回 14:35:00.000（向过去）。
//
//   - 截断到 15 分钟：t.Truncate(15 * time.Minute)
//
//     跃点：--- 14:00:00 --- 14:15:00 --- 14:30:00 ---- t ---- 14:45:00 ---
//     ↑
//     当前时间 14:35:29.650
//
//     返回上一个跃点：14:30:00.000（向过去）。
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

// IsLeapYear 返回 t 是否长年
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

// IsSame 返回 t 与 target 在指定单位下是否相同（包含上级单位一致性）。
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
		l := t.cascade(fromAbs, false, u)
		r := target.cascade(fromAbs, false, u)
		return l.Eq(r)
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

// Near 在集合 times 中寻找距离 base 最近或最远的时间点。
//
// 参数 op：
//   - "<" : 最近 (Closest)
//   - ">" : 最远 (Farthest)
func Near(op string, base Time, times ...Time) Time {
	if len(times) == 0 {
		return base
	}

	res := times[0]
	diff := math.Abs(base.Sub(res).Seconds())

	for i := 1; i < len(times); i++ {
		secs := math.Abs(base.Sub(times[i]).Seconds())
		if op == ">" && secs > diff {
			diff, res = secs, times[i]
		} else if op == "<" && secs < diff {
			diff, res = secs, times[i]
		}
	}

	return res
}

// Maxmin 返回一组时间中的极值
//
// 参数 op：
//   - ">" : 最大值（最晚）
//   - "<" : 最小值（最早）
func Maxmin(op string, times ...Time) Time {
	if len(times) == 0 {
		return Time{}
	}

	res := times[0]
	for i := 1; i < len(times); i++ {
		if op == ">" && times[i].Gt(res) {
			res = times[i]
		} else if op == "<" && times[i].Lt(res) {
			res = times[i]
		}
	}

	return res
}

// Between 判断 t 是否在 (start, end) 区间内。
//
// 可选参数 bounds 用于控制边界包含性（默认为 "="）：
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

// Since 返回 t 到现在经过的持续时间
func Since(t Time) time.Duration {
	return time.Since(t.time)
}

// Until 返回现在到 t 经过的持续时间
func Until(t Time) time.Duration {
	return time.Until(t.time)
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
