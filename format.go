package aeon

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const (
	ANSIC         = "Mon Jan _2 15:04:05 2006"
	UnixDate      = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate      = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822        = "02 Jan 06 15:04 MST"
	RFC822Z       = "02 Jan 06 15:04 -0700"
	RFC850        = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123       = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z      = "Mon, 02 Jan 2006 15:04:05 -0700"
	RFC3339       = "2006-01-02T15:04:05Z07:00"
	RFC3339Ns     = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen       = "3:04PM"
	Stamp         = "Jan _2 15:04:05"
	StampMilli    = "Jan _2 15:04:05.000"
	StampMicro    = "Jan _2 15:04:05.000000"
	StampNano     = "Jan _2 15:04:05.000000000"
	DateTime      = "2006-01-02 15:04:05"
	DateTimeMilli = "2006-01-02 15:04:05.000"
	DateTimeMicro = "2006-01-02 15:04:05.000000"
	DateTimeNano  = "2006-01-02 15:04:05.000000000"
	DateTimeNs    = "2006-01-02 15:04:05.999999999"
	DateOnly      = "2006-01-02"
	TimeOnly      = "15:04:05"

	// 补强布局 (归一化复用版)
	DateTimeFull       = "2006-01-02 15:04:05.999999999 -0700 MST"
	DateTimeTZ         = "2006-01-02 15:04:05-07:00"
	DateTimeTZShort    = "2006-01-02 15:04:05-07"
	DateTimeISO        = "2006-01-02T15:04:05-07:00"
	DateTimePMMST      = "2006-01-02 15:04:05 PM MST"
	DateTimePMShortMST = "2006-01-02 15:04:05PM MST"
	RFC3339Space       = "2006-01-02 15:04:05Z07:00"
	RFC3339NsSpace     = "2006-01-02 15:04:05.999999999Z07:00"

	// 紧凑与特殊格式
	DateTimeCompact   = "20060102150405"
	DateCompact       = "20060102"
	DateTimeVeryShort = "2006-1-2 15:4:5"
	DateTimeShort     = "2006-1-2 15:4"
	DateOnlyShort     = "2006-1-2"
	TimeVeryShort     = "15:4:5"
	TimeShort         = "15:4"
	MonthDay          = "1-2"
	YearOnly          = "2006"

	// 从 Carbon 补强的高价值布局
	DateTimeCompactTZ     = "20060102150405-07:00"
	DateTimeCompactZ      = "20060102150405Z07:00"
	DateTimeCompactMilli  = "20060102150405.000"
	DateHourShort         = "2006-1-2 15"
	HourOnly              = "15"
	DateMonth             = "2006-1"
	FormattedDate         = "Jan 2, 2006"
	FormattedDayDate      = "Mon, Jan 2, 2006"
	DayDateTime           = "Mon, Jan 2, 2006 3:04 PM"
	Cookie                = "Monday, 02-Jan-2006 15:04:05 MST"
	Http                  = "Mon, 02 Jan 2006 15:04:05 GMT"
	RFC1036               = "Mon, 02 Jan 06 15:04:05 -0700"
	RFC7231               = "Mon, 02 Jan 2006 15:04:05 MST"
	TimeTZShort           = "15:04:05-07"
	DateTimeFullVeryShort = "2006-1-2 15:4:5 -0700 MST"
	DateTimeNsVeryShort   = "2006-1-2 15:4:5.999999999"

	// ISO8601 家族 (带偏移量)
	ISO8601      = "2006-01-02T15:04:05-07:00"
	ISO8601Milli = "2006-01-02T15:04:05.000-07:00"
	ISO8601Micro = "2006-01-02T15:04:05.000000-07:00"
	ISO8601Nano  = "2006-01-02T15:04:05.000000000-07:00"

	// ISO8601 Zulu 家族 (带 Z)
	ISO8601Zulu      = "2006-01-02T15:04:05Z"
	ISO8601ZuluMilli = "2006-01-02T15:04:05.000Z"
	ISO8601ZuluMicro = "2006-01-02T15:04:05.000000Z"
	ISO8601ZuluNano  = "2006-01-02T15:04:05.000000000Z"

	// 标准协议别名
	Atom = RFC3339
	W3C  = RFC3339
	RSS  = "Mon, 02 Jan 2006 15:04:05 -0700"

	// 精度对齐系列 (用于对齐输出)
	DateMilli   = "2006-01-02.000"
	DateMicro   = "2006-01-02.000000"
	DateNano    = "2006-01-02.000000000"
	TimeMilli   = "15:04:05.000"
	TimeMicro   = "15:04:05.000000"
	TimeNano    = "15:04:05.000000000"
	TimeCompact = "150405"
)

var Formats = []string{
	// 1. 现代 API 最常用格式 (已通过归一化消除 T/Space 冗余)
	DateTime, DateTimeMilli, DateTimeMicro, DateTimeNano, DateTimeNs,
	DateOnly, RFC3339Space, RFC3339NsSpace,
	DateTimeFull, DateTimeTZ, DateTimeTZShort,
	ISO8601, ISO8601Milli, ISO8601Micro, ISO8601Nano,
	ISO8601Zulu, ISO8601ZuluMilli, ISO8601ZuluMicro, ISO8601ZuluNano,

	// 2. 特殊与 PM 格式
	DateTimePMMST, DateTimePMShortMST,
	DateCompact, DateTimeCompact,
	DateTimeCompactTZ, DateTimeCompactZ, DateTimeCompactMilli,

	// 3. 命名与标准格式
	RFC1123Z, RFC1123, RFC822Z, RFC822, ANSIC, UnixDate, RubyDate, RFC850, Kitchen,
	Cookie, Http, RFC1036, RFC7231, RSS,

	// 4. 时间戳样式
	StampNano, StampMicro, StampMilli, Stamp,

	// 5. 人类友好与宽松格式
	DayDateTime, FormattedDate, FormattedDayDate,
	DateTimeFullVeryShort, DateTimeNsVeryShort,
	DateTimeVeryShort, DateTimeShort, DateOnlyShort, DateHourShort,
	TimeOnly, TimeVeryShort, TimeShort, TimeTZShort, TimeCompact,
	DateMonth, MonthDay, YearOnly, HourOnly,
}

// --- 格式化时间 ---

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(DateTime)+2)
	b = append(b, '"')
	b = t.time.AppendFormat(b, DateTime)
	return append(b, '"'), nil
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	*t, err = ParseE(string(b), t.Location())
	return
}

func (t Time) MarshalText() ([]byte, error) {
	if t.IsZero() {
		return []byte(""), nil
	}
	return []byte(t.String()), nil
}

func (t *Time) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}

func (t *Time) Scan(value any) (err error) {
	switch v := value.(type) {
	case time.Time:
		*t = Aeon(v)
	case string:
		*t, err = ParseE(v, t.Location())
	default:
		*t = Aeon()
	}
	return
}

func (t Time) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.time, nil
}

// --- 自定义格式器 ---

type Formatter interface {
	Layout() string
}

type F[T Formatter] struct {
	Time
	m T
}

func (f *F[T]) layout() string {
	if l := f.m.Layout(); l != "" {
		return l
	}
	return DateTime
}

func (f F[T]) MarshalJSON() ([]byte, error) {
	if f.IsZero() {
		return []byte("null"), nil
	}
	layout := f.layout()
	b := make([]byte, 0, len(layout)+2)
	b = append(b, '"')
	b = f.time.AppendFormat(b, layout)
	return append(b, '"'), nil
}

func (f *F[T]) UnmarshalJSON(b []byte) (err error) {
	f.Time, err = ParseE(string(b), f.Location())
	return
}

func (f F[T]) MarshalText() ([]byte, error) {
	if f.IsZero() {
		return []byte(""), nil
	}
	return []byte(f.Format(f.layout())), nil
}

func (f *F[T]) UnmarshalText(data []byte) error {
	return f.UnmarshalJSON(data)
}

func (f *F[T]) Scan(value any) (err error) {
	switch v := value.(type) {
	case time.Time:
		f.Time = Aeon(v)
	case string:
		f.Time, err = ParseE(v, f.Location())
	default:
		f.Time = Aeon()
	}
	return
}

func (f F[T]) Value() (driver.Value, error) {
	if f.IsZero() {
		return nil, nil
	}
	return f.Time.Time(), nil
}

type milliFormat struct{}

func (milliFormat) Layout() string { return DateTimeMilli }

type MilliTime = F[milliFormat]

// --- 解析引擎 ---

func ParseE(s string, loc ...*time.Location) (Time, error) {
	s = strings.Trim(strings.TrimSpace(s), "\"")
	if s == "" || s == "null" {
		return Time{}, nil
	}

	// 1. 斜杠归一化
	s = strings.ReplaceAll(s, "/", "-")

	// 2. 分隔符与点号智能归一化
	if idx := strings.IndexAny(s, " T"); idx != -1 {
		// 有日期时间分隔符: 只替换分界点之前的点，且将分隔符统一为空格
		s = strings.ReplaceAll(s[:idx], ".", "-") + " " + s[idx+1:]
	} else if dotCount := strings.Count(s, "."); dotCount >= 2 {
		// 无分隔符且点数多
		if dotCount == 2 {
			s = strings.ReplaceAll(s, ".", "-")
		} else {
			lastDot := strings.LastIndex(s, ".")
			s = strings.ReplaceAll(s[:lastDot], ".", "-") + s[lastDot:]
		}
	}

	l := DefaultTimeZone
	if len(loc) > 0 && loc[0] != nil {
		l = loc[0]
	}

	for _, layout := range Formats {
		if t, err := time.ParseInLocation(layout, s, l); err == nil {
			return Aeon(t), nil
		}
	}

	return Time{}, fmt.Errorf("aeon 无法解析时间: %q", s)
}

func Parse(value string, loc ...*time.Location) Time {
	t, _ := ParseE(value, loc...)
	return t
}

func ParseByE(layout string, value string, loc ...*time.Location) (Time, error) {
	l := DefaultTimeZone
	if len(loc) > 0 && loc[0] != nil {
		l = loc[0]
	}
	pt, err := time.ParseInLocation(layout, value, l)
	return Time{time: pt, weekStarts: DefaultWeekStarts}, err
}

func ParseBy(layout string, value string, loc ...*time.Location) Time {
	t, _ := ParseByE(layout, value, loc...)
	return t
}

// --- 格式化 API ---

func (t Time) Format(layout string) string                 { return t.time.Format(layout) }
func (t Time) AppendFormat(b []byte, layout string) []byte { return t.time.AppendFormat(b, layout) }
func (t Time) String() string                              { return t.time.Format(DateTimeNs) }
func (t Time) ToString(f ...string) string {
	if len(f) > 0 {
		return t.time.Format(f[0])
	}
	return t.time.Format(DateTimeNs)
}
