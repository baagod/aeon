package aeon

import "time"

type from int

const (
	fromAbs from = iota // Start/EndCentury (全绝对)
	fromRel             // Start/EndByCentury (全相对)
	fromAt              // StartAt/EndCentury (定位后偏移: Abs + Rel..)
	fromIn              // Start/EndInCentury (偏移后定位: Rel + Abs..)
	fromGoAbs
	fromGoRel
	fromGoAt
	fromGoIn
	fromOffset // Add 全相对不对齐
)

const (
	// ISO 标志位
	ISO = -1000000
	// Overflow 允许月份溢出标志
	Overflow = -2000000
	// ABS 绝对时间（仅对世纪、年份生效）标志
	ABS = -3000000
)

// Flag 承载级联操作的上下文配置
type Flag struct {
	iso      bool // 是否使用 ISO 周历
	overflow bool // 是否允许溢出
	abs      bool // 是否使用绝对值模式 (Thru Abs)
	fill     bool // 是否是对齐到结束 (End系列)
	goMode   bool // 是否是跳转模式 (不归零)
}

// cascade 级联时间
// Start/EndCentury (全绝对)
// Start/EndByCentury (全相对)
// StartAt/EndCentury (定位后偏移: Abs + Rel..)
// Start/EndInCentury (偏移后定位: Rel + Abs..)
// Add 全相对不对齐
func (t Time) cascade(f from, fill bool, u Unit, args ...int) Time {
	y, m, d := t.Date()
	h, mm, sec := t.Clock()
	ns := t.time.Nanosecond()

	c := Flag{
		fill:   fill,
		goMode: f >= fromGoAbs && f <= fromGoIn,
	}

	// 标志位解析循环
	var z int
Loop:
	for ; z < len(args); z++ {
		switch args[z] {
		case ISO:
			c.iso = true
		case Overflow:
			c.overflow = true
		case ABS:
			c.abs = true
		default:
			break Loop
		}
	}

	if args = args[z:]; len(args) == 0 {
		if f == fromOffset {
			args = oneArgs
		} else {
			args = zeroArgs
		}
	}

	p := u
	seq := u.seq()
	w := t.Weekday()
	sw := t.weekStarts

	for i, n := range args {
		if i >= len(seq) {
			break
		}

		unit := seq[i]

		switch f {
		case fromOffset:
			y, m, d, h, mm, sec, ns, w = applyOffset(c, unit, p, n, y, m, d, h, mm, sec, ns, w)
		case fromAbs, fromGoAbs:
			y, m, d, h, mm, sec, ns, w = applyAbs(c, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
		case fromRel, fromGoRel:
			y, m, d, h, mm, sec, ns, w = applyRel(c, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
		case fromAt, fromGoAt:
			if i == 0 {
				y, m, d, h, mm, sec, ns, w = applyAbs(c, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
			} else {
				y, m, d, h, mm, sec, ns, w = applyRel(c, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
			}
		case fromIn, fromGoIn:
			if i == 0 {
				y, m, d, h, mm, sec, ns, w = applyRel(c, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
			} else {
				y, m, d, h, mm, sec, ns, w = applyAbs(c, unit, p, n, y, m, d, h, mm, sec, ns, w, sw)
			}
		}

		p = unit
	}

	if !c.goMode && f != fromOffset {
		y, m, d, h, mm, sec, ns = align(c, p, y, m, d, h, mm, sec, ns)
	}

	return Time{
		time:       time.Date(y, time.Month(m), d, h, mm, sec, ns, t.Location()),
		weekStarts: t.weekStarts,
	}
}

// --- 添加时间 ---

func (t Time) Add(d time.Duration) Time {
	return Time{time: t.time.Add(d), weekStarts: t.weekStarts}
}

func (t Time) AddCentury(n ...int) Time  { return t.cascade(fromOffset, false, Century, n...) }
func (t Time) AddDecade(n ...int) Time   { return t.cascade(fromOffset, false, Decade, n...) }
func (t Time) AddYear(n ...int) Time     { return t.cascade(fromOffset, false, Year, n...) }
func (t Time) AddMonth(n ...int) Time    { return t.cascade(fromOffset, false, Month, n...) }
func (t Time) AddDay(n ...int) Time      { return t.cascade(fromOffset, false, Day, n...) }
func (t Time) AddHour(n ...int) Time     { return t.cascade(fromOffset, false, Hour, n...) }
func (t Time) AddMinute(n ...int) Time   { return t.cascade(fromOffset, false, Minute, n...) }
func (t Time) AddSecond(n ...int) Time   { return t.cascade(fromOffset, false, Second, n...) }
func (t Time) AddMilli(n ...int) Time    { return t.cascade(fromOffset, false, Millisecond, n...) }
func (t Time) AddMicro(n ...int) Time    { return t.cascade(fromOffset, false, Microsecond, n...) }
func (t Time) AddNano(n ...int) Time     { return t.cascade(fromOffset, false, Nanosecond, n...) }
func (t Time) AddQuarter(n ...int) Time  { return t.cascade(fromOffset, false, Quarter, n...) }
func (t Time) AddWeek(n ...int) Time     { return t.cascade(fromOffset, false, Week, n...) }
func (t Time) AddWeekday(n ...int) Time  { return t.cascade(fromOffset, false, Weekday, n...) }
func (t Time) AddYearWeek(n ...int) Time { return t.cascade(fromOffset, false, YearWeek, n...) }

// --- 全绝对定位级联 ---

func (t Time) Start(n ...int) Time {
	return t.cascade(fromAbs, false, Year, append([]int{ABS}, n...)...)
}

func (t Time) End(n ...int) Time {
	return t.cascade(fromAbs, true, Year, append([]int{ABS}, n...)...)
}

func (t Time) StartCentury(n ...int) Time  { return t.cascade(fromAbs, false, Century, n...) }
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

func (t Time) EndCentury(n ...int) Time  { return t.cascade(fromAbs, true, Century, n...) }
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

func (t Time) StartBy(n ...int) Time {
	return t.cascade(fromRel, false, Year, append([]int{ABS}, n...)...)
}

func (t Time) EndBy(n ...int) Time {
	return t.cascade(fromRel, true, Year, append([]int{ABS}, n...)...)
}

func (t Time) StartByCentury(n ...int) Time  { return t.cascade(fromRel, false, Century, n...) }
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

func (t Time) EndByCentury(n ...int) Time  { return t.cascade(fromRel, true, Century, n...) }
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

func (t Time) StartAt(n ...int) Time {
	return t.cascade(fromAt, false, Year, append([]int{ABS}, n...)...)
}

func (t Time) EndAt(n ...int) Time {
	return t.cascade(fromAt, true, Year, append([]int{ABS}, n...)...)
}

func (t Time) StartAtCentury(n ...int) Time  { return t.cascade(fromAt, false, Century, n...) }
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

func (t Time) EndAtCentury(n ...int) Time  { return t.cascade(fromAt, true, Century, n...) }
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

func (t Time) StartIn(n ...int) Time {
	return t.cascade(fromIn, false, Year, append([]int{ABS}, n...)...)
}

func (t Time) EndIn(n ...int) Time {
	return t.cascade(fromIn, true, Year, append([]int{ABS}, n...)...)
}

func (t Time) StartInCentury(n ...int) Time  { return t.cascade(fromIn, false, Century, n...) }
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

func (t Time) EndInCentury(n ...int) Time  { return t.cascade(fromIn, true, Century, n...) }
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

// --- Go 系列的绝对年方法 ---

func (t Time) Go(n ...int) Time {
	return t.cascade(fromGoAbs, false, Year, append([]int{ABS}, n...)...)
}

func (t Time) GoBy(n ...int) Time {
	return t.cascade(fromGoRel, false, Year, append([]int{ABS}, n...)...)
}

func (t Time) GoAt(n ...int) Time {
	return t.cascade(fromGoAt, false, Year, append([]int{ABS}, n...)...)
}

func (t Time) GoIn(n ...int) Time {
	return t.cascade(fromGoIn, false, Year, append([]int{ABS}, n...)...)
}

// --- Start 的保留精度版本 ---

func (t Time) GoCentury(n ...int) Time  { return t.cascade(fromGoAbs, false, Century, n...) }
func (t Time) GoDecade(n ...int) Time   { return t.cascade(fromGoAbs, false, Decade, n...) }
func (t Time) GoYear(n ...int) Time     { return t.cascade(fromGoAbs, false, Year, n...) }
func (t Time) GoMonth(n ...int) Time    { return t.cascade(fromGoAbs, false, Month, n...) }
func (t Time) GoDay(n ...int) Time      { return t.cascade(fromGoAbs, false, Day, n...) }
func (t Time) GoHour(n ...int) Time     { return t.cascade(fromGoAbs, false, Hour, n...) }
func (t Time) GoMinute(n ...int) Time   { return t.cascade(fromGoAbs, false, Minute, n...) }
func (t Time) GoSecond(n ...int) Time   { return t.cascade(fromGoAbs, false, Second, n...) }
func (t Time) GoMilli(n ...int) Time    { return t.cascade(fromGoAbs, false, Millisecond, n...) }
func (t Time) GoMicro(n ...int) Time    { return t.cascade(fromGoAbs, false, Microsecond, n...) }
func (t Time) GoNano(n ...int) Time     { return t.cascade(fromGoAbs, false, Nanosecond, n...) }
func (t Time) GoQuarter(n ...int) Time  { return t.cascade(fromGoAbs, false, Quarter, n...) }
func (t Time) GoWeek(n ...int) Time     { return t.cascade(fromGoAbs, false, Week, n...) }
func (t Time) GoWeekday(n ...int) Time  { return t.cascade(fromGoAbs, false, Weekday, n...) }
func (t Time) GoYearWeek(n ...int) Time { return t.cascade(fromGoAbs, false, YearWeek, n...) }

// --- StartBy 的保留精度版本 ---

func (t Time) GoByCentury(n ...int) Time  { return t.cascade(fromGoRel, false, Century, n...) }
func (t Time) GoByDecade(n ...int) Time   { return t.cascade(fromGoRel, false, Decade, n...) }
func (t Time) GoByYear(n ...int) Time     { return t.cascade(fromGoRel, false, Year, n...) }
func (t Time) GoByMonth(n ...int) Time    { return t.cascade(fromGoRel, false, Month, n...) }
func (t Time) GoByDay(n ...int) Time      { return t.cascade(fromGoRel, false, Day, n...) }
func (t Time) GoByHour(n ...int) Time     { return t.cascade(fromGoRel, false, Hour, n...) }
func (t Time) GoByMinute(n ...int) Time   { return t.cascade(fromGoRel, false, Minute, n...) }
func (t Time) GoBySecond(n ...int) Time   { return t.cascade(fromGoRel, false, Second, n...) }
func (t Time) GoByMilli(n ...int) Time    { return t.cascade(fromGoRel, false, Millisecond, n...) }
func (t Time) GoByMicro(n ...int) Time    { return t.cascade(fromGoRel, false, Microsecond, n...) }
func (t Time) GoByNano(n ...int) Time     { return t.cascade(fromGoRel, false, Nanosecond, n...) }
func (t Time) GoByQuarter(n ...int) Time  { return t.cascade(fromGoRel, false, Quarter, n...) }
func (t Time) GoByWeek(n ...int) Time     { return t.cascade(fromGoRel, false, Week, n...) }
func (t Time) GoByWeekday(n ...int) Time  { return t.cascade(fromGoRel, false, Weekday, n...) }
func (t Time) GoByYearWeek(n ...int) Time { return t.cascade(fromGoRel, false, YearWeek, n...) }

// --- StartAt 的保留精度版本 ---

func (t Time) GoAtCentury(n ...int) Time  { return t.cascade(fromGoAt, false, Century, n...) }
func (t Time) GoAtDecade(n ...int) Time   { return t.cascade(fromGoAt, false, Decade, n...) }
func (t Time) GoAtYear(n ...int) Time     { return t.cascade(fromGoAt, false, Year, n...) }
func (t Time) GoAtMonth(n ...int) Time    { return t.cascade(fromGoAt, false, Month, n...) }
func (t Time) GoAtDay(n ...int) Time      { return t.cascade(fromGoAt, false, Day, n...) }
func (t Time) GoAtHour(n ...int) Time     { return t.cascade(fromGoAt, false, Hour, n...) }
func (t Time) GoAtMinute(n ...int) Time   { return t.cascade(fromGoAt, false, Minute, n...) }
func (t Time) GoAtSecond(n ...int) Time   { return t.cascade(fromGoAt, false, Second, n...) }
func (t Time) GoAtMilli(n ...int) Time    { return t.cascade(fromGoAt, false, Millisecond, n...) }
func (t Time) GoAtMicro(n ...int) Time    { return t.cascade(fromGoAt, false, Microsecond, n...) }
func (t Time) GoAtNano(n ...int) Time     { return t.cascade(fromGoAt, false, Nanosecond, n...) }
func (t Time) GoAtQuarter(n ...int) Time  { return t.cascade(fromGoAt, false, Quarter, n...) }
func (t Time) GoAtWeek(n ...int) Time     { return t.cascade(fromGoAt, false, Week, n...) }
func (t Time) GoAtWeekday(n ...int) Time  { return t.cascade(fromGoAt, false, Weekday, n...) }
func (t Time) GoAtYearWeek(n ...int) Time { return t.cascade(fromGoAt, false, YearWeek, n...) }

// --- StartIn 的保留精度版本 ---

func (t Time) GoInCentury(n ...int) Time  { return t.cascade(fromGoIn, false, Century, n...) }
func (t Time) GoInDecade(n ...int) Time   { return t.cascade(fromGoIn, false, Decade, n...) }
func (t Time) GoInYear(n ...int) Time     { return t.cascade(fromGoIn, false, Year, n...) }
func (t Time) GoInMonth(n ...int) Time    { return t.cascade(fromGoIn, false, Month, n...) }
func (t Time) GoInDay(n ...int) Time      { return t.cascade(fromGoIn, false, Day, n...) }
func (t Time) GoInHour(n ...int) Time     { return t.cascade(fromGoIn, false, Hour, n...) }
func (t Time) GoInMinute(n ...int) Time   { return t.cascade(fromGoIn, false, Minute, n...) }
func (t Time) GoInSecond(n ...int) Time   { return t.cascade(fromGoIn, false, Second, n...) }
func (t Time) GoInMilli(n ...int) Time    { return t.cascade(fromGoIn, false, Millisecond, n...) }
func (t Time) GoInMicro(n ...int) Time    { return t.cascade(fromGoIn, false, Microsecond, n...) }
func (t Time) GoInNano(n ...int) Time     { return t.cascade(fromGoIn, false, Nanosecond, n...) }
func (t Time) GoInQuarter(n ...int) Time  { return t.cascade(fromGoIn, false, Quarter, n...) }
func (t Time) GoInWeek(n ...int) Time     { return t.cascade(fromGoIn, false, Week, n...) }
func (t Time) GoInWeekday(n ...int) Time  { return t.cascade(fromGoIn, false, Weekday, n...) }
func (t Time) GoInYearWeek(n ...int) Time { return t.cascade(fromGoIn, false, YearWeek, n...) }
