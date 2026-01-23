package aeon

import (
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
    Quarter // 季度流
    Week    // 月周流
    Weekday
)

type path int

const (
    seAbs path = iota // Start/EndCentury (全绝对)
    seRel             // Start/EndByCentury (全相对)
    seAt              // StartAt/EndCentury (定位后偏移: Abs + Rel..)
    seIn              // Start/EndInCentury (偏移后定位: Rel + Abs..)
    goAbs
    goRel
    goAt
    goIn
)

const (
    // flagSign 是标志位的特征基座 (-1,073,741,824)，确保标志位处于 int32 的深水区
    flagSign = -1 << 30
    // flagThreshold 是标志位识别门槛 (-1,000,000,000)。
    // 任何小于此门槛的参数均被视为标志位包。
    flagThreshold = -1e9

    ISO  = flagSign | (1 << 0) // ISO 周标志
    Ord  = flagSign | (1 << 1) // Ord 周标志
    Full = flagSign | (1 << 2) // Full 周标志

    // Overflow 允许月份溢出标志
    Overflow = flagSign | (1 << 3)
    // ABS 绝对时间标志 (内部使用)
    ABS = flagSign | (1 << 4)
)

// Flag 承载级联操作的上下文配置
type Flag struct {
    isoWeek  bool // [ISO] 周标志 (遵循 ISO 周规则)
    fullWeek bool // [完整] 周标志 (从本月首周一开始)
    ordWeek  bool // [序数] 周标志 (从本月1日开始)
    overflow bool // 是否允许溢出
    abs      bool // 是否绝对年模式
    fill     bool // 是否置满时间
    goMode   bool // 是否跳转模式
}

// cascade 级联时间核心引擎
func cascade(t Time, f path, fill bool, u Unit, args ...int) Time {
    y, m, d := t.Date()
    h, mm, s := t.Clock()
    ns := t.time.Nanosecond()
    w := t.Weekday()
    sw := t.weekStarts

    // 🦬 级解析：提取首位参数的位掩码标志位
    c := Flag{fill: fill, goMode: f >= goAbs}

    if len(args) > 0 && args[0] < flagThreshold {
        mask := args[0]
        c.isoWeek = mask&ISO == ISO
        c.fullWeek = mask&Full == Full
        c.ordWeek = mask&Ord == Ord
        c.overflow = mask&Overflow == Overflow
        c.abs = mask&ABS == ABS
        args = args[1:]
    }

    if len(args) == 0 {
        if f == goRel {
            args = oneArgs
        } else {
            args = zeroArgs
        }
    }

    p, pN := u, args[0] // 父单元及其传值
    if len(args) == 1 { // 单参数路径
        y, m, d, h, mm, s, ns, w = apply(f, c, true, u, u, args[0], pN, y, m, d, h, mm, s, ns, w, sw)
    } else { // 级联循环
        seq := u.seq()
        if l := len(seq); len(args) > l {
            args = args[:l]
        }

        for i, n := range args {
            unit := seq[i]
            y, m, d, h, mm, s, ns, w = apply(f, c, i == 0, unit, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
            p, pN = unit, n
        }
    }

    if !c.goMode { // go 模式不对齐时间 (归零或置满)
        y, m, d, h, mm, s, ns = align(c, p, y, m, d, h, mm, s, ns)
    }

    return Time{
        time:       time.Date(y, time.Month(m), d, h, mm, s, ns, t.Location()),
        weekStarts: t.weekStarts,
    }
}

// a 归零时间
func a(t Time, p path, u Unit, args ...int) Time {
    return cascade(t, p, false, u, args...)
}

// z 置满时间
func z(t Time, p path, u Unit, args ...int) Time {
    return cascade(t, p, true, u, args...)
}

// --- 顶级导航方法（首个参数定位到真正的绝对年份）---

func (t Time) Start(n ...int) Time   { return a(t, seAbs, Year, append(absArgs, n...)...) }
func (t Time) StartAt(n ...int) Time { return a(t, seAt, Year, append(absArgs, n...)...) }

func (t Time) End(n ...int) Time   { return z(t, seAbs, Year, append(absArgs, n...)...) }
func (t Time) EndAt(n ...int) Time { return z(t, seAt, Year, append(absArgs, n...)...) }

func (t Time) Go(n ...int) Time { return a(t, goAbs, Year, append(absArgs, n...)...) }
func (t Time) At(n ...int) Time { return a(t, goAt, Year, append(absArgs, n...)...) }

// --- 全绝对定位级联 ---

func (t Time) StartCentury(n ...int) Time { return a(t, seAbs, Century, n...) }
func (t Time) StartDecade(n ...int) Time  { return a(t, seAbs, Decade, n...) }
func (t Time) StartYear(n ...int) Time    { return a(t, seAbs, Year, n...) }
func (t Time) StartMonth(n ...int) Time   { return a(t, seAbs, Month, n...) }
func (t Time) StartDay(n ...int) Time     { return a(t, seAbs, Day, n...) }
func (t Time) StartHour(n ...int) Time    { return a(t, seAbs, Hour, n...) }
func (t Time) StartMinute(n ...int) Time  { return a(t, seAbs, Minute, n...) }
func (t Time) StartSecond(n ...int) Time  { return a(t, seAbs, Second, n...) }
func (t Time) StartMilli(n ...int) Time   { return a(t, seAbs, Millisecond, n...) }
func (t Time) StartMicro(n ...int) Time   { return a(t, seAbs, Microsecond, n...) }
func (t Time) StartNano(n ...int) Time    { return a(t, seAbs, Nanosecond, n...) }
func (t Time) StartQuarter(n ...int) Time { return a(t, seAbs, Quarter, n...) }
func (t Time) StartWeek(n ...int) Time    { return a(t, seAbs, Week, n...) }
func (t Time) StartWeekday(n ...int) Time { return a(t, seAbs, Weekday, n...) }

func (t Time) EndCentury(n ...int) Time { return z(t, seAbs, Century, n...) }
func (t Time) EndDecade(n ...int) Time  { return z(t, seAbs, Decade, n...) }
func (t Time) EndYear(n ...int) Time    { return z(t, seAbs, Year, n...) }
func (t Time) EndMonth(n ...int) Time   { return z(t, seAbs, Month, n...) }
func (t Time) EndDay(n ...int) Time     { return z(t, seAbs, Day, n...) }
func (t Time) EndHour(n ...int) Time    { return z(t, seAbs, Hour, n...) }
func (t Time) EndMinute(n ...int) Time  { return z(t, seAbs, Minute, n...) }
func (t Time) EndSecond(n ...int) Time  { return z(t, seAbs, Second, n...) }
func (t Time) EndMilli(n ...int) Time   { return z(t, seAbs, Millisecond, n...) }
func (t Time) EndMicro(n ...int) Time   { return z(t, seAbs, Microsecond, n...) }
func (t Time) EndNano(n ...int) Time    { return z(t, seAbs, Nanosecond, n...) }
func (t Time) EndQuarter(n ...int) Time { return z(t, seAbs, Quarter, n...) }
func (t Time) EndWeek(n ...int) Time    { return z(t, seAbs, Week, n...) }
func (t Time) EndWeekday(n ...int) Time { return z(t, seAbs, Weekday, n...) }

// --- 全相对定位级联 ---

func (t Time) StartByCentury(n ...int) Time { return a(t, seRel, Century, n...) }
func (t Time) StartByDecade(n ...int) Time  { return a(t, seRel, Decade, n...) }
func (t Time) StartByYear(n ...int) Time    { return a(t, seRel, Year, n...) }
func (t Time) StartByMonth(n ...int) Time   { return a(t, seRel, Month, n...) }
func (t Time) StartByDay(n ...int) Time     { return a(t, seRel, Day, n...) }
func (t Time) StartByHour(n ...int) Time    { return a(t, seRel, Hour, n...) }
func (t Time) StartByMinute(n ...int) Time  { return a(t, seRel, Minute, n...) }
func (t Time) StartBySecond(n ...int) Time  { return a(t, seRel, Second, n...) }
func (t Time) StartByMilli(n ...int) Time   { return a(t, seRel, Millisecond, n...) }
func (t Time) StartByMicro(n ...int) Time   { return a(t, seRel, Microsecond, n...) }
func (t Time) StartByNano(n ...int) Time    { return a(t, seRel, Nanosecond, n...) }
func (t Time) StartByQuarter(n ...int) Time { return a(t, seRel, Quarter, n...) }
func (t Time) StartByWeek(n ...int) Time    { return a(t, seRel, Week, n...) }
func (t Time) StartByWeekday(n ...int) Time { return a(t, seRel, Weekday, n...) }

func (t Time) EndByCentury(n ...int) Time { return z(t, seRel, Century, n...) }
func (t Time) EndByDecade(n ...int) Time  { return z(t, seRel, Decade, n...) }
func (t Time) EndByYear(n ...int) Time    { return z(t, seRel, Year, n...) }
func (t Time) EndByMonth(n ...int) Time   { return z(t, seRel, Month, n...) }
func (t Time) EndByDay(n ...int) Time     { return z(t, seRel, Day, n...) }
func (t Time) EndByHour(n ...int) Time    { return z(t, seRel, Hour, n...) }
func (t Time) EndByMinute(n ...int) Time  { return z(t, seRel, Minute, n...) }
func (t Time) EndBySecond(n ...int) Time  { return z(t, seRel, Second, n...) }
func (t Time) EndByMilli(n ...int) Time   { return z(t, seRel, Millisecond, n...) }
func (t Time) EndByMicro(n ...int) Time   { return z(t, seRel, Microsecond, n...) }
func (t Time) EndByNano(n ...int) Time    { return z(t, seRel, Nanosecond, n...) }
func (t Time) EndByQuarter(n ...int) Time { return z(t, seRel, Quarter, n...) }
func (t Time) EndByWeek(n ...int) Time    { return z(t, seRel, Week, n...) }
func (t Time) EndByWeekday(n ...int) Time { return z(t, seRel, Weekday, n...) }

// ---- 锚位（绝对）后偏移级联 ----

func (t Time) StartAtCentury(n ...int) Time { return a(t, seAt, Century, n...) }
func (t Time) StartAtDecade(n ...int) Time  { return a(t, seAt, Decade, n...) }
func (t Time) StartAtYear(n ...int) Time    { return a(t, seAt, Year, n...) }
func (t Time) StartAtMonth(n ...int) Time   { return a(t, seAt, Month, n...) }
func (t Time) StartAtDay(n ...int) Time     { return a(t, seAt, Day, n...) }
func (t Time) StartAtHour(n ...int) Time    { return a(t, seAt, Hour, n...) }
func (t Time) StartAtMinute(n ...int) Time  { return a(t, seAt, Minute, n...) }
func (t Time) StartAtSecond(n ...int) Time  { return a(t, seAt, Second, n...) }
func (t Time) StartAtMilli(n ...int) Time   { return a(t, seAt, Millisecond, n...) }
func (t Time) StartAtMicro(n ...int) Time   { return a(t, seAt, Microsecond, n...) }
func (t Time) StartAtNano(n ...int) Time    { return a(t, seAt, Nanosecond, n...) }
func (t Time) StartAtQuarter(n ...int) Time { return a(t, seAt, Quarter, n...) }
func (t Time) StartAtWeek(n ...int) Time    { return a(t, seAt, Week, n...) }
func (t Time) StartAtWeekday(n ...int) Time { return a(t, seAt, Weekday, n...) }

func (t Time) EndAtCentury(n ...int) Time { return z(t, seAt, Century, n...) }
func (t Time) EndAtDecade(n ...int) Time  { return z(t, seAt, Decade, n...) }
func (t Time) EndAtYear(n ...int) Time    { return z(t, seAt, Year, n...) }
func (t Time) EndAtMonth(n ...int) Time   { return z(t, seAt, Month, n...) }
func (t Time) EndAtDay(n ...int) Time     { return z(t, seAt, Day, n...) }
func (t Time) EndAtHour(n ...int) Time    { return z(t, seAt, Hour, n...) }
func (t Time) EndAtMinute(n ...int) Time  { return z(t, seAt, Minute, n...) }
func (t Time) EndAtSecond(n ...int) Time  { return z(t, seAt, Second, n...) }
func (t Time) EndAtMilli(n ...int) Time   { return z(t, seAt, Millisecond, n...) }
func (t Time) EndAtMicro(n ...int) Time   { return z(t, seAt, Microsecond, n...) }
func (t Time) EndAtNano(n ...int) Time    { return z(t, seAt, Nanosecond, n...) }
func (t Time) EndAtQuarter(n ...int) Time { return z(t, seAt, Quarter, n...) }
func (t Time) EndAtWeek(n ...int) Time    { return z(t, seAt, Week, n...) }
func (t Time) EndAtWeekday(n ...int) Time { return z(t, seAt, Weekday, n...) }

// ---- 偏移后锚位（绝对）级联 ----

func (t Time) StartInCentury(n ...int) Time { return a(t, seIn, Century, n...) }
func (t Time) StartInDecade(n ...int) Time  { return a(t, seIn, Decade, n...) }
func (t Time) StartInYear(n ...int) Time    { return a(t, seIn, Year, n...) }
func (t Time) StartInMonth(n ...int) Time   { return a(t, seIn, Month, n...) }
func (t Time) StartInDay(n ...int) Time     { return a(t, seIn, Day, n...) }
func (t Time) StartInHour(n ...int) Time    { return a(t, seIn, Hour, n...) }
func (t Time) StartInMinute(n ...int) Time  { return a(t, seIn, Minute, n...) }
func (t Time) StartInSecond(n ...int) Time  { return a(t, seIn, Second, n...) }
func (t Time) StartInMilli(n ...int) Time   { return a(t, seIn, Millisecond, n...) }
func (t Time) StartInMicro(n ...int) Time   { return a(t, seIn, Microsecond, n...) }
func (t Time) StartInNano(n ...int) Time    { return a(t, seIn, Nanosecond, n...) }
func (t Time) StartInQuarter(n ...int) Time { return a(t, seIn, Quarter, n...) }
func (t Time) StartInWeek(n ...int) Time    { return a(t, seIn, Week, n...) }
func (t Time) StartInWeekday(n ...int) Time { return a(t, seIn, Weekday, n...) }

func (t Time) EndInCentury(n ...int) Time { return z(t, seIn, Century, n...) }
func (t Time) EndInDecade(n ...int) Time  { return z(t, seIn, Decade, n...) }
func (t Time) EndInYear(n ...int) Time    { return z(t, seIn, Year, n...) }
func (t Time) EndInMonth(n ...int) Time   { return z(t, seIn, Month, n...) }
func (t Time) EndInDay(n ...int) Time     { return z(t, seIn, Day, n...) }
func (t Time) EndInHour(n ...int) Time    { return z(t, seIn, Hour, n...) }
func (t Time) EndInMinute(n ...int) Time  { return z(t, seIn, Minute, n...) }
func (t Time) EndInSecond(n ...int) Time  { return z(t, seIn, Second, n...) }
func (t Time) EndInMilli(n ...int) Time   { return z(t, seIn, Millisecond, n...) }
func (t Time) EndInMicro(n ...int) Time   { return z(t, seIn, Microsecond, n...) }
func (t Time) EndInNano(n ...int) Time    { return z(t, seIn, Nanosecond, n...) }
func (t Time) EndInQuarter(n ...int) Time { return z(t, seIn, Quarter, n...) }
func (t Time) EndInWeek(n ...int) Time    { return z(t, seIn, Week, n...) }
func (t Time) EndInWeekday(n ...int) Time { return z(t, seIn, Weekday, n...) }

// --- Start 的保留精度版本 ---

func (t Time) GoCentury(n ...int) Time { return a(t, goAbs, Century, n...) }
func (t Time) GoDecade(n ...int) Time  { return a(t, goAbs, Decade, n...) }
func (t Time) GoYear(n ...int) Time    { return a(t, goAbs, Year, n...) }
func (t Time) GoMonth(n ...int) Time   { return a(t, goAbs, Month, n...) }
func (t Time) GoDay(n ...int) Time     { return a(t, goAbs, Day, n...) }
func (t Time) GoHour(n ...int) Time    { return a(t, goAbs, Hour, n...) }
func (t Time) GoMinute(n ...int) Time  { return a(t, goAbs, Minute, n...) }
func (t Time) GoSecond(n ...int) Time  { return a(t, goAbs, Second, n...) }
func (t Time) GoMilli(n ...int) Time   { return a(t, goAbs, Millisecond, n...) }
func (t Time) GoMicro(n ...int) Time   { return a(t, goAbs, Microsecond, n...) }
func (t Time) GoNano(n ...int) Time    { return a(t, goAbs, Nanosecond, n...) }
func (t Time) GoQuarter(n ...int) Time { return a(t, goAbs, Quarter, n...) }
func (t Time) GoWeek(n ...int) Time    { return a(t, goAbs, Week, n...) }
func (t Time) GoWeekday(n ...int) Time { return a(t, goAbs, Weekday, n...) }

// --- StartAt 的保留精度版本 ---

func (t Time) AtCentury(n ...int) Time { return a(t, goAt, Century, n...) }
func (t Time) AtDecade(n ...int) Time  { return a(t, goAt, Decade, n...) }
func (t Time) AtYear(n ...int) Time    { return a(t, goAt, Year, n...) }
func (t Time) AtMonth(n ...int) Time   { return a(t, goAt, Month, n...) }
func (t Time) AtDay(n ...int) Time     { return a(t, goAt, Day, n...) }
func (t Time) AtHour(n ...int) Time    { return a(t, goAt, Hour, n...) }
func (t Time) AtMinute(n ...int) Time  { return a(t, goAt, Minute, n...) }
func (t Time) AtSecond(n ...int) Time  { return a(t, goAt, Second, n...) }
func (t Time) AtMilli(n ...int) Time   { return a(t, goAt, Millisecond, n...) }
func (t Time) AtMicro(n ...int) Time   { return a(t, goAt, Microsecond, n...) }
func (t Time) AtNano(n ...int) Time    { return a(t, goAt, Nanosecond, n...) }
func (t Time) AtQuarter(n ...int) Time { return a(t, goAt, Quarter, n...) }
func (t Time) AtWeek(n ...int) Time    { return a(t, goAt, Week, n...) }
func (t Time) AtWeekday(n ...int) Time { return a(t, goAt, Weekday, n...) }

// --- StartIn 的保留精度版本 ---

func (t Time) InCentury(n ...int) Time { return a(t, goIn, Century, n...) }
func (t Time) InDecade(n ...int) Time  { return a(t, goIn, Decade, n...) }
func (t Time) InYear(n ...int) Time    { return a(t, goIn, Year, n...) }
func (t Time) InMonth(n ...int) Time   { return a(t, goIn, Month, n...) }
func (t Time) InDay(n ...int) Time     { return a(t, goIn, Day, n...) }
func (t Time) InHour(n ...int) Time    { return a(t, goIn, Hour, n...) }
func (t Time) InMinute(n ...int) Time  { return a(t, goIn, Minute, n...) }
func (t Time) InSecond(n ...int) Time  { return a(t, goIn, Second, n...) }
func (t Time) InMilli(n ...int) Time   { return a(t, goIn, Millisecond, n...) }
func (t Time) InMicro(n ...int) Time   { return a(t, goIn, Microsecond, n...) }
func (t Time) InNano(n ...int) Time    { return a(t, goIn, Nanosecond, n...) }
func (t Time) InQuarter(n ...int) Time { return a(t, goIn, Quarter, n...) }
func (t Time) InWeek(n ...int) Time    { return a(t, goIn, Week, n...) }

// --- 添加时间 ---

func (t Time) Add(d time.Duration) Time { return Time{time: t.time.Add(d), weekStarts: t.weekStarts} }
func (t Time) AddCentury(n ...int) Time { return a(t, goRel, Century, n...) }
func (t Time) AddDecade(n ...int) Time  { return a(t, goRel, Decade, n...) }
func (t Time) AddYear(n ...int) Time    { return a(t, goRel, Year, n...) }
func (t Time) AddMonth(n ...int) Time   { return a(t, goRel, Month, n...) }
func (t Time) AddDay(n ...int) Time     { return a(t, goRel, Day, n...) }
func (t Time) AddHour(n ...int) Time    { return a(t, goRel, Hour, n...) }
func (t Time) AddMinute(n ...int) Time  { return a(t, goRel, Minute, n...) }
func (t Time) AddSecond(n ...int) Time  { return a(t, goRel, Second, n...) }
func (t Time) AddMilli(n ...int) Time   { return a(t, goRel, Millisecond, n...) }
func (t Time) AddMicro(n ...int) Time   { return a(t, goRel, Microsecond, n...) }
func (t Time) AddNano(n ...int) Time    { return a(t, goRel, Nanosecond, n...) }
func (t Time) AddQuarter(n ...int) Time { return a(t, goRel, Quarter, n...) }
func (t Time) AddWeek(n ...int) Time    { return a(t, goRel, Week, n...) }
