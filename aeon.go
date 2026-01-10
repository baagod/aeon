package thru

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
	// DefaultWeekStartsAt 全局默认周起始日（默认为周一）
	DefaultWeekStartsAt = time.Monday
	// pow10 预定义的 10 的幂次方表，用于高精度计算
	pow10 = [...]int64{
		1, 10, 100, 1000, 10000, 100000, 1e6, 1e7, 1e8, 1e9,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
	}
)

type Time struct {
	time         time.Time
	weekStartsAt time.Weekday
}

// --- 创建时间 ---

func New(t ...time.Time) Time {
	if len(t) == 0 {
		t = []time.Time{{}}
	}
	return Time{time: t[0], weekStartsAt: DefaultWeekStartsAt}
}

func Now() Time {
	return Time{time: time.Now(), weekStartsAt: DefaultWeekStartsAt}
}

func Date[M ~int](
	year int, month M, day, hour,
	minute, sec, nsec int, loc *time.Location,
) Time {
	return Time{
		time:         time.Date(year, time.Month(month), day, hour, minute, sec, nsec, loc),
		weekStartsAt: DefaultWeekStartsAt,
	}
}

// WithWeekStartsAt 返回新实例，周起始日为 w。
func (t Time) WithWeekStartsAt(w time.Weekday) Time {
	return Time{time: t.time, weekStartsAt: w}
}

// Unix 返回给定时间戳的时间。secs 可以是秒、毫秒、微妙或纳秒级时间戳。
func Unix(secs int64) Time {
	v := secs
	if secs < 0 { // 处理公元前时间戳
		v = -secs
	}

	switch {
	case v <= 9999999999: // 10位：秒
		return New(time.Unix(secs, 0))
	case v <= 9999999999999: // 13位：毫秒
		return New(time.UnixMilli(secs))
	case v <= 9999999999999999: // 16位：微秒
		return New(time.UnixMicro(secs))
	default: // 19位及以上：纳秒
		return New(time.Unix(0, secs))
	}
}

// --- 获取时间 ---

// Year 返回本年
func (t Time) Year() int {
	return t.time.Year()
}

// Month 返回本月
func (t Time) Month() int {
	return int(t.time.Month())
}

// Day 返回本月的第几天
func (t Time) Day() int {
	return t.time.Day()
}

// Hour 返回小时，范围 [0, 23]
func (t Time) Hour() int {
	return t.time.Hour()
}

// Minute 返回分钟，范围 [0, 59]
func (t Time) Minute() int {
	return t.time.Minute()
}

// Second 返回时间的秒数或指定纳秒精度的小数部分
//
// 参数 n (可选) 指定返回的精度：
//   - 不提供或 0: 返回整秒数 (0-59)
//   - 1-9: 返回纳秒精度的小数部分，n 表示小数位数。
func (t Time) Second(n ...int) int {
	if len(n) == 0 || n[0] == 0 {
		return t.time.Second()
	}
	divisor := pow10[9-clamp(n[0], 1, 9)]
	return t.time.Nanosecond() / int(divisor)
}

// Date 返回 t 的年月日
func (t Time) Date() (int, int, int) {
	y, m, d := t.time.Date()
	return y, int(m), d
}

// Clock 返回一天中的小时、分钟和秒
func (t Time) Clock() (hour, minute, sec int) {
	return t.time.Clock()
}

// Weekday 返回星期几
func (t Time) Weekday() time.Weekday {
	return t.time.Weekday()
}

// YearDay 返回一年中的第几天，非闰年范围 [1,365]，闰年范围 [1,366]。
func (t Time) YearDay() int {
	return t.time.YearDay()
}

// Days 返回本年总天数
func (t Time) Days() int {
	return DaysIn(t.Year())
}

// MonthDays 返回本月总天数
func (t Time) MonthDays() int {
	return DaysIn(t.Year(), t.Month())
}

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
	return Time{time: t.time.UTC(), weekStartsAt: t.weekStartsAt}
}

// Local 返回本地时间
func (t Time) Local() Time {
	return Time{time: t.time.Local(), weekStartsAt: t.weekStartsAt}
}

// In 返回指定的 loc 时间
func (t Time) In(loc *time.Location) Time {
	return Time{time: t.time.In(loc), weekStartsAt: t.weekStartsAt}
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
	return Time{time: t.time.Round(d), weekStartsAt: t.weekStartsAt}
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
	return Time{time: t.time.Truncate(d), weekStartsAt: t.weekStartsAt}
}

// Time 返回 time.Time
func (t Time) Time() time.Time {
	return t.time
}

// Location 返回时区信息
func (t Time) Location() *time.Location {
	return t.time.Location()
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
		tDays := float64(t.YearDay()) / float64(t.Days()) // t 在本年的进度 (0~1)
		uDays := float64(u.YearDay()) / float64(u.Days()) // u 在本年的进度 (0~1)
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
		diff = months + days/float64(u.MonthDays())
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

// Before 返回 t 是否在 u 之前 (t < u)
func (t Time) Before(u Time) bool {
	return t.time.Before(u.time)
}

// After 返回 t 是否在 u 之后 (t > u)
func (t Time) After(u Time) bool {
	return t.time.After(u.time)
}

// Equal 返回 t == u
func (t Time) Equal(u Time) bool {
	return t.time.Equal(u.time)
}

// Compare 比较 t 和 u。
func (t Time) Compare(u Time) int {
	return t.time.Compare(u.time)
}

// Since 返回 t 到现在经过的持续时间
func Since(t Time) time.Duration {
	return time.Since(t.time)
}

// Until 返回现在到 t 经过的持续时间
func Until(t Time) time.Duration {
	return time.Until(t.time)
}

// ---- 其他 ----

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

// IsDST 返回时间是否夏令时
func (t Time) IsDST() bool {
	return t.time.IsDST()
}
