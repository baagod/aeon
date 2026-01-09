package thru

import (
	"database/sql/driver"
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
	ISOYearWeek
	Weekday
)

type from int

const (
	fromAbs     from = iota // Start/End (全绝对)
	fromRel                 // StartBy/EndBy (全相对)
	fromAt                  // StartAt/EndAt (定位后偏移: Abs + Rel..)
	fromIn                  // StartIn/EndIn (偏移后定位: Rel + Abs..)
	fromNoalign             // Add 全相对不对齐
)

const (
	// ISO 特殊标志位
	ISO = -1000000
	// Overflow 允许月份溢出标志位
	Overflow = -2000000
)

var (
	// DefaultWeekStartsAt 全局默认周起始日（默认为周一）
	DefaultWeekStartsAt = time.Monday
)

type Time struct {
	time         time.Time
	weekStartsAt time.Weekday
}

// ---- 创建时间 ----

func New(t time.Time) Time {
	return Time{time: t, weekStartsAt: DefaultWeekStartsAt}
}

func Now() Time {
	return Time{time: time.Now(), weekStartsAt: DefaultWeekStartsAt}
}

func Date[M ~int](
	year int, month M, day, hour,
	minute, sec, nsec int, loc *time.Location,
) Time {
	m := time.Month(month)
	return Time{
		time:         time.Date(year, m, day, hour, minute, sec, nsec, loc),
		weekStartsAt: DefaultWeekStartsAt,
	}
}

// WithWeekStartsAt 返回新实例，周起始日为 w。
func (t Time) WithWeekStartsAt(w time.Weekday) Time {
	return Time{time: t.time, weekStartsAt: w}
}

// Unix 返回给定时间戳的本地时间。secs 可以是秒、毫秒或纳秒级时间戳。
func Unix(secs int64) Time {
	if secs <= 9999999999 { // 10 位及以下，视为秒级时间戳
		return Time{time: time.Unix(secs, 0), weekStartsAt: DefaultWeekStartsAt}
	}
	return Time{time: time.Unix(0, secs), weekStartsAt: DefaultWeekStartsAt}
}

// cascade 级联时间
// Start/End (全绝对)
// StartBy/EndBy (全相对)
// StartAt/EndAt (锚定后偏移: Abs + Rel..)
// StartIn/EndIn (偏移后定位: Rel + Abs..)
func (t Time) cascade(f from, end bool, u Unit, args ...int) Time {
	y, month, d := t.time.Date()
	h, mm, sec := t.time.Clock()
	m, ns := int(month), t.time.Nanosecond()

	if len(args) == 0 {
		args = zeroArgs
	}

	var overflow bool
	if args[0] == ISO {
		if args = args[1:]; len(args) == 0 || u == Weekday {
			args = append([]int{0}, args...)
		}
		u = ISOYearWeek
	} else if args[0] == Overflow {
		overflow, args = true, args[1:]
	}

	p := u
	seq := sequence(u)
	w := t.Weekday()
	sw := t.weekStartsAt

	for i, n := range args {
		if i >= len(seq) {
			break
		}

		unit := seq[i]

		switch f {
		case fromNoalign:
			y, m, d, h, mm, sec, ns, w = applyOffset(end, overflow, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
		case fromAbs:
			y, m, d, h, mm, sec, ns, w = applyAbs(end, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
		case fromRel:
			y, m, d, h, mm, sec, ns, w = applyRel(end, overflow, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
		case fromAt:
			if i == 0 {
				y, m, d, h, mm, sec, ns, w = applyAbs(end, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
			} else {
				y, m, d, h, mm, sec, ns, w = applyRel(end, overflow, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
			}
		case fromIn:
			if i == 0 {
				y, m, d, h, mm, sec, ns, w = applyRel(end, overflow, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
			} else {
				y, m, d, h, mm, sec, ns, w = applyAbs(end, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
			}
		}

		p = unit
	}

	if f != fromNoalign {
		y, m, d, h, mm, sec, ns = align(end, p, y, m, d, h, mm, sec, ns)
	}

	return Time{
		time:         time.Date(y, time.Month(m), d, h, mm, sec, ns, t.Location()),
		weekStartsAt: t.weekStartsAt,
	}
}

// ---- 添加时间 ----

func (t Time) Add(d time.Duration) Time {
	return Time{time: t.time.Add(d), weekStartsAt: t.weekStartsAt}
}

func (t Time) AddCentury(n ...int) Time  { return t.cascade(fromNoalign, false, Century, n...) }
func (t Time) AddDecade(n ...int) Time   { return t.cascade(fromNoalign, false, Decade, n...) }
func (t Time) AddYear(n ...int) Time     { return t.cascade(fromNoalign, false, Year, n...) }
func (t Time) AddMonth(n ...int) Time    { return t.cascade(fromNoalign, false, Month, n...) }
func (t Time) AddDay(n ...int) Time      { return t.cascade(fromNoalign, false, Day, n...) }
func (t Time) AddHour(n ...int) Time     { return t.cascade(fromNoalign, false, Hour, n...) }
func (t Time) AddMinute(n ...int) Time   { return t.cascade(fromNoalign, false, Minute, n...) }
func (t Time) AddSecond(n ...int) Time   { return t.cascade(fromNoalign, false, Second, n...) }
func (t Time) AddMilli(n ...int) Time    { return t.cascade(fromNoalign, false, Millisecond, n...) }
func (t Time) AddMicro(n ...int) Time    { return t.cascade(fromNoalign, false, Microsecond, n...) }
func (t Time) AddNano(n ...int) Time     { return t.cascade(fromNoalign, false, Nanosecond, n...) }
func (t Time) AddQuarter(n ...int) Time  { return t.cascade(fromNoalign, false, Quarter, n...) }
func (t Time) AddWeek(n ...int) Time     { return t.cascade(fromNoalign, false, Week, n...) }
func (t Time) AddWeekday(n ...int) Time  { return t.cascade(fromNoalign, false, Weekday, n...) }
func (t Time) AddYearWeek(n ...int) Time { return t.cascade(fromNoalign, false, YearWeek, n...) }

// --- 全绝对定位级联 ---

func (t Time) Start(n ...int) Time         { return t.cascade(fromAbs, false, Century, n...) }
func (t Time) StartDecade(n ...int) Time   { return t.cascade(fromAbs, false, Decade, n...) }
func (t Time) StartYear(n ...int) Time     { return t.cascade(fromAbs, false, Year, n...) }
func (t Time) StartMonth(n ...int) Time    { return t.cascade(fromAbs, false, Month, n...) }
func (t Time) StartDay(n ...int) Time      { return t.cascade(fromAbs, false, Day, n...) }
func (t Time) StartHour(n ...int) Time     { return t.cascade(fromAbs, false, Hour, n...) }
func (t Time) StartMinute(n ...int) Time   { return t.cascade(fromAbs, false, Minute, n...) }
func (t Time) StartSecond(n ...int) Time   { return t.cascade(fromAbs, false, Second, n...) }
func (t Time) StartMilli(n ...int) Time    { return t.cascade(fromAbs, false, Millisecond, n...) }
func (t Time) StartMicro(n ...int) Time    { return t.cascade(fromAbs, false, Microsecond, n...) }
func (t Time) StartNano(n ...int) Time     { return t.cascade(fromAbs, false, Nanosecond, n...) }
func (t Time) StartQuarter(n ...int) Time  { return t.cascade(fromAbs, false, Quarter, n...) }
func (t Time) StartWeek(n ...int) Time     { return t.cascade(fromAbs, false, Week, n...) }
func (t Time) StartWeekday(n ...int) Time  { return t.cascade(fromAbs, false, Weekday, n...) }
func (t Time) StartYearWeek(n ...int) Time { return t.cascade(fromAbs, false, YearWeek, n...) }

func (t Time) End(n ...int) Time         { return t.cascade(fromAbs, true, Century, n...) }
func (t Time) EndDecade(n ...int) Time   { return t.cascade(fromAbs, true, Decade, n...) }
func (t Time) EndYear(n ...int) Time     { return t.cascade(fromAbs, true, Year, n...) }
func (t Time) EndMonth(n ...int) Time    { return t.cascade(fromAbs, true, Month, n...) }
func (t Time) EndDay(n ...int) Time      { return t.cascade(fromAbs, true, Day, n...) }
func (t Time) EndHour(n ...int) Time     { return t.cascade(fromAbs, true, Hour, n...) }
func (t Time) EndMinute(n ...int) Time   { return t.cascade(fromAbs, true, Minute, n...) }
func (t Time) EndSecond(n ...int) Time   { return t.cascade(fromAbs, true, Second, n...) }
func (t Time) EndMilli(n ...int) Time    { return t.cascade(fromAbs, true, Millisecond, n...) }
func (t Time) EndMicro(n ...int) Time    { return t.cascade(fromAbs, true, Microsecond, n...) }
func (t Time) EndNano(n ...int) Time     { return t.cascade(fromAbs, true, Nanosecond, n...) }
func (t Time) EndQuarter(n ...int) Time  { return t.cascade(fromAbs, true, Quarter, n...) }
func (t Time) EndWeek(n ...int) Time     { return t.cascade(fromAbs, true, Week, n...) }
func (t Time) EndWeekday(n ...int) Time  { return t.cascade(fromAbs, true, Weekday, n...) }
func (t Time) EndYearWeek(n ...int) Time { return t.cascade(fromAbs, true, YearWeek, n...) }

// --- 全相对定位级联 ---

func (t Time) StartBy(n ...int) Time         { return t.cascade(fromRel, false, Century, n...) }
func (t Time) StartByDecade(n ...int) Time   { return t.cascade(fromRel, false, Decade, n...) }
func (t Time) StartByYear(n ...int) Time     { return t.cascade(fromRel, false, Year, n...) }
func (t Time) StartByMonth(n ...int) Time    { return t.cascade(fromRel, false, Month, n...) }
func (t Time) StartByDay(n ...int) Time      { return t.cascade(fromRel, false, Day, n...) }
func (t Time) StartByHour(n ...int) Time     { return t.cascade(fromRel, false, Hour, n...) }
func (t Time) StartByMinute(n ...int) Time   { return t.cascade(fromRel, false, Minute, n...) }
func (t Time) StartBySecond(n ...int) Time   { return t.cascade(fromRel, false, Second, n...) }
func (t Time) StartByMilli(n ...int) Time    { return t.cascade(fromRel, false, Millisecond, n...) }
func (t Time) StartByMicro(n ...int) Time    { return t.cascade(fromRel, false, Microsecond, n...) }
func (t Time) StartByNano(n ...int) Time     { return t.cascade(fromRel, false, Nanosecond, n...) }
func (t Time) StartByQuarter(n ...int) Time  { return t.cascade(fromRel, false, Quarter, n...) }
func (t Time) StartByWeek(n ...int) Time     { return t.cascade(fromRel, false, Week, n...) }
func (t Time) StartByWeekday(n ...int) Time  { return t.cascade(fromRel, false, Weekday, n...) }
func (t Time) StartByYearWeek(n ...int) Time { return t.cascade(fromRel, false, YearWeek, n...) }

func (t Time) EndBy(n ...int) Time         { return t.cascade(fromRel, true, Century, n...) }
func (t Time) EndByDecade(n ...int) Time   { return t.cascade(fromRel, true, Decade, n...) }
func (t Time) EndByYear(n ...int) Time     { return t.cascade(fromRel, true, Year, n...) }
func (t Time) EndByMonth(n ...int) Time    { return t.cascade(fromRel, true, Month, n...) }
func (t Time) EndByDay(n ...int) Time      { return t.cascade(fromRel, true, Day, n...) }
func (t Time) EndByHour(n ...int) Time     { return t.cascade(fromRel, true, Hour, n...) }
func (t Time) EndByMinute(n ...int) Time   { return t.cascade(fromRel, true, Minute, n...) }
func (t Time) EndBySecond(n ...int) Time   { return t.cascade(fromRel, true, Second, n...) }
func (t Time) EndByMilli(n ...int) Time    { return t.cascade(fromRel, true, Millisecond, n...) }
func (t Time) EndByMicro(n ...int) Time    { return t.cascade(fromRel, true, Microsecond, n...) }
func (t Time) EndByNano(n ...int) Time     { return t.cascade(fromRel, true, Nanosecond, n...) }
func (t Time) EndByQuarter(n ...int) Time  { return t.cascade(fromRel, true, Quarter, n...) }
func (t Time) EndByWeek(n ...int) Time     { return t.cascade(fromRel, true, Week, n...) }
func (t Time) EndByWeekday(n ...int) Time  { return t.cascade(fromRel, true, Weekday, n...) }
func (t Time) EndByYearWeek(n ...int) Time { return t.cascade(fromRel, true, YearWeek, n...) }

// ---- 锚位（绝对）后偏移级联 ----

func (t Time) StartAt(n ...int) Time         { return t.cascade(fromAt, false, Century, n...) }
func (t Time) StartAtDecade(n ...int) Time   { return t.cascade(fromAt, false, Decade, n...) }
func (t Time) StartAtYear(n ...int) Time     { return t.cascade(fromAt, false, Year, n...) }
func (t Time) StartAtMonth(n ...int) Time    { return t.cascade(fromAt, false, Month, n...) }
func (t Time) StartAtDay(n ...int) Time      { return t.cascade(fromAt, false, Day, n...) }
func (t Time) StartAtHour(n ...int) Time     { return t.cascade(fromAt, false, Hour, n...) }
func (t Time) StartAtMinute(n ...int) Time   { return t.cascade(fromAt, false, Minute, n...) }
func (t Time) StartAtSecond(n ...int) Time   { return t.cascade(fromAt, false, Second, n...) }
func (t Time) StartAtMilli(n ...int) Time    { return t.cascade(fromAt, false, Millisecond, n...) }
func (t Time) StartAtMicro(n ...int) Time    { return t.cascade(fromAt, false, Microsecond, n...) }
func (t Time) StartAtNano(n ...int) Time     { return t.cascade(fromAt, false, Nanosecond, n...) }
func (t Time) StartAtQuarter(n ...int) Time  { return t.cascade(fromAt, false, Quarter, n...) }
func (t Time) StartAtWeek(n ...int) Time     { return t.cascade(fromAt, false, Week, n...) }
func (t Time) StartAtWeekday(n ...int) Time  { return t.cascade(fromAt, false, Weekday, n...) }
func (t Time) StartAtYearWeek(n ...int) Time { return t.cascade(fromAt, false, YearWeek, n...) }

func (t Time) EndAt(n ...int) Time         { return t.cascade(fromAt, true, Century, n...) }
func (t Time) EndAtDecade(n ...int) Time   { return t.cascade(fromAt, true, Decade, n...) }
func (t Time) EndAtYear(n ...int) Time     { return t.cascade(fromAt, true, Year, n...) }
func (t Time) EndAtMonth(n ...int) Time    { return t.cascade(fromAt, true, Month, n...) }
func (t Time) EndAtDay(n ...int) Time      { return t.cascade(fromAt, true, Day, n...) }
func (t Time) EndAtHour(n ...int) Time     { return t.cascade(fromAt, true, Hour, n...) }
func (t Time) EndAtMinute(n ...int) Time   { return t.cascade(fromAt, true, Minute, n...) }
func (t Time) EndAtSecond(n ...int) Time   { return t.cascade(fromAt, true, Second, n...) }
func (t Time) EndAtMilli(n ...int) Time    { return t.cascade(fromAt, true, Millisecond, n...) }
func (t Time) EndAtMicro(n ...int) Time    { return t.cascade(fromAt, true, Microsecond, n...) }
func (t Time) EndAtNano(n ...int) Time     { return t.cascade(fromAt, true, Nanosecond, n...) }
func (t Time) EndAtQuarter(n ...int) Time  { return t.cascade(fromAt, true, Quarter, n...) }
func (t Time) EndAtWeek(n ...int) Time     { return t.cascade(fromAt, true, Week, n...) }
func (t Time) EndAtWeekday(n ...int) Time  { return t.cascade(fromAt, true, Weekday, n...) }
func (t Time) EndAtYearWeek(n ...int) Time { return t.cascade(fromAt, true, YearWeek, n...) }

// ---- 偏移后锚位（绝对）级联 ----

func (t Time) StartIn(n ...int) Time         { return t.cascade(fromIn, false, Century, n...) }
func (t Time) StartInDecade(n ...int) Time   { return t.cascade(fromIn, false, Decade, n...) }
func (t Time) StartInYear(n ...int) Time     { return t.cascade(fromIn, false, Year, n...) }
func (t Time) StartInMonth(n ...int) Time    { return t.cascade(fromIn, false, Month, n...) }
func (t Time) StartInDay(n ...int) Time      { return t.cascade(fromIn, false, Day, n...) }
func (t Time) StartInHour(n ...int) Time     { return t.cascade(fromIn, false, Hour, n...) }
func (t Time) StartInMinute(n ...int) Time   { return t.cascade(fromIn, false, Minute, n...) }
func (t Time) StartInSecond(n ...int) Time   { return t.cascade(fromIn, false, Second, n...) }
func (t Time) StartInMilli(n ...int) Time    { return t.cascade(fromIn, false, Millisecond, n...) }
func (t Time) StartInMicro(n ...int) Time    { return t.cascade(fromIn, false, Microsecond, n...) }
func (t Time) StartInNano(n ...int) Time     { return t.cascade(fromIn, false, Nanosecond, n...) }
func (t Time) StartInQuarter(n ...int) Time  { return t.cascade(fromIn, false, Quarter, n...) }
func (t Time) StartInWeek(n ...int) Time     { return t.cascade(fromIn, false, Week, n...) }
func (t Time) StartInWeekday(n ...int) Time  { return t.cascade(fromIn, false, Weekday, n...) }
func (t Time) StartInYearWeek(n ...int) Time { return t.cascade(fromIn, false, YearWeek, n...) }

func (t Time) EndIn(n ...int) Time         { return t.cascade(fromIn, true, Century, n...) }
func (t Time) EndInDecade(n ...int) Time   { return t.cascade(fromIn, true, Decade, n...) }
func (t Time) EndInYear(n ...int) Time     { return t.cascade(fromIn, true, Year, n...) }
func (t Time) EndInMonth(n ...int) Time    { return t.cascade(fromIn, true, Month, n...) }
func (t Time) EndInDay(n ...int) Time      { return t.cascade(fromIn, true, Day, n...) }
func (t Time) EndInHour(n ...int) Time     { return t.cascade(fromIn, true, Hour, n...) }
func (t Time) EndInMinute(n ...int) Time   { return t.cascade(fromIn, true, Minute, n...) }
func (t Time) EndInSecond(n ...int) Time   { return t.cascade(fromIn, true, Second, n...) }
func (t Time) EndInMilli(n ...int) Time    { return t.cascade(fromIn, true, Millisecond, n...) }
func (t Time) EndInMicro(n ...int) Time    { return t.cascade(fromIn, true, Microsecond, n...) }
func (t Time) EndInNano(n ...int) Time     { return t.cascade(fromIn, true, Nanosecond, n...) }
func (t Time) EndInQuarter(n ...int) Time  { return t.cascade(fromIn, true, Quarter, n...) }
func (t Time) EndInWeek(n ...int) Time     { return t.cascade(fromIn, true, Week, n...) }
func (t Time) EndInWeekday(n ...int) Time  { return t.cascade(fromIn, true, Weekday, n...) }
func (t Time) EndInYearWeek(n ...int) Time { return t.cascade(fromIn, true, YearWeek, n...) }

// ---- 选择时间 ----

// Go 偏移 ±y 年并选择 m 月 d 日，如果 m, d 为负数，则从最后的月、日开始偏移。
func (t Time) Go(y int, md ...int) Time {
	year, month, day := t.Year()+y, t.Month(), t.Day()

	if i := len(md); i > 0 {
		if m := float64(md[0]); m > 0 {
			month = int(math.Min(m, 12))
		} else if m < 0 {
			month = int(math.Max(13+m, 1))
		}
		if i > 1 {
			day = md[1]
		}
	}

	maxDay := float64(DaysIn(year, month))
	if d := float64(day); d > 0 {
		day = int(math.Min(d, maxDay))
	} else {
		day = int(math.Max(maxDay+d+1, 1))
	}

	h, mm, sec := t.time.Clock()
	date := time.Date(year, time.Month(month), day, h, mm, sec, t.time.Nanosecond(), t.Location())
	return Time{time: date, weekStartsAt: t.weekStartsAt}
}

// GoYear 和 Go() 一样，但 y 指定为确切年份而非偏移。
func (t Time) GoYear(y int, md ...int) Time {
	return t.Go(y-t.Year(), md...)
}

func (t Time) GoMonth(m int, d ...int) Time {
	day := 0
	if len(d) > 0 {
		day = d[0]
	}
	return t.Go(0, m, day)
}

func (t Time) GoDay(d int) Time {
	return t.Go(0, 0, d)
}

// ---- 获取时间 ----

// Year 返回 t 的年份
func (t Time) Year() int {
	return t.time.Year()
}

// Month 返回 t 的月份
func (t Time) Month() int {
	return int(t.time.Month())
}

// Day 返回 t 的天数
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
//   - 1-9: 返回纳秒精度的小数部分，n 表示小数位数
func (t Time) Second(n ...int) int {
	if len(n) == 0 || n[0] == 0 {
		return t.time.Second()
	}
	divisor := int(math.Pow10(9 - clamp(n[0], 1, 9)))
	return t.time.Nanosecond() / divisor
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
	divisor := int64(math.Pow10(19 - precision))
	return t.time.UnixNano() / divisor
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
// 该函数将返回距离当前时间最近的那个 "跃点"，如果当前时间正好位于两个 "跃点" 中间，返回指向未来的那个 "跃点"。
//
// 示例：假设当前时间是 2021-07-21 14:35:29.650
//
//   - 舍入到秒：t.Round(time.Second)
//
//     根据定义，时间 14:35:29.650 距离下一个跃点 14:35:30.000 最近 (只需要 0.35ns)，
//     而距离上一个跃点 14:35:29.000 较远 (需要 65ns)，故返回下一个跃点 14:35:30.000。
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

// Truncate 与 Round 类似，但是 Truncate 永远返回指向过去时间的 "跃点"，不会进行四舍五入。
//
// 可视化理解：
//
//	----|----|----|----|----
//	   d1    t   d2   d3
//
// Truncate 将永远返回 d1，如果时间 t 正好处于某个 "跃点" 位置，返回 t。
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

// DiffIn 返回 t 和 u 的时间差。使用 unit 参数指定比较单位：
//   - "y": 年
//   - "M": 月
//   - "d": 日
//   - "h": 小时
//   - "m": 分钟
//   - "s": 秒
func (t Time) DiffIn(u Time, unit string) float64 {
	switch unit {
	case "y":
		tDays := float64(t.YearDay()) / float64(t.Days())
		uDays := float64(u.YearDay()) / float64(u.Days())
		return float64(t.Year()-u.Year()) + tDays - uDays
	case "M":
		// 计算 [初始月差] 并将结果转换为浮点数：初始月差 = (年份差 * 12) + 月份差
		months := float64((t.Year()-u.Year())*12 + t.Month() - u.Month())
		// 计算 [天差] 并将结果转换为浮点数：天差 = t.Day() - u.Day()
		days := float64(t.Day() - u.Day())
		// 如果天差为负数，表示未满一个完整的月，需要将月差减 1。
		if days < 0 {
			months--
		}
		// 计算 [总月份差]，包括小数部分：总月份差 = 初始月差 + 天差 / u.Days()
		// 计算小数部分使用：天差 / u.MonthDays()，得到天差所占 u 月天数的比值，
		// 再与 [初始月差] 相加得到 [总月份差]。
		return months + days/float64(u.MonthDays())
	case "d":
		return t.Sub(u).Hours() / 24
	case "h":
		return t.Sub(u).Hours()
	case "m":
		return t.Sub(u).Minutes()
	case "s":
		return t.Sub(u).Seconds()
	}
	return 0
}

func (t Time) DiffAbsIn(u Time, unit string) float64 {
	return math.Abs(t.DiffIn(u, unit))
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

// ---- 序列化时间 ----

// Scan 由 DB 转到 Go 时调用
func (t *Time) Scan(value any) error {
	if v, ok := value.(time.Time); ok {
		*t = Time{time: v, weekStartsAt: DefaultWeekStartsAt}
	}
	return nil
}

// Value 由 Go 转到 DB 时调用
func (t Time) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.time, nil
}

// MarshalJSON 将 t 转为 JSON 字符串时调用
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(DateTime) + `"`), nil
}

// UnmarshalJSON 将 JSON 字符串转为 t 时调用
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	*t, err = ParseE(string(b))
	return
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
