package aeon

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const (
	ANSIC      = "Mon Jan _2 15:04:05 2006"
	UnixD      = "Mon Jan _2 15:04:05 MST 2006"
	RubyD      = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822     = "02 Jan 06 15:04 MST"
	RFC822Z    = "02 Jan 06 15:04 -0700"
	RFC850     = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123    = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z   = "Mon, 02 Jan 2006 15:04:05 -0700"
	RFC3339    = "2006-01-02T15:04:05Z07:00"
	RFC3339Ns  = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen    = "3:04PM"
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
	StampNs    = "Jan _2 15:04:05.999999999"

	// 核心布局 (DT/D 系列)
	DT       = "2006-01-02 15:04:05"
	DTMilli  = "2006-01-02 15:04:05.000"
	DTMicro  = "2006-01-02 15:04:05.000000"
	DTNano   = "2006-01-02 15:04:05.000000000"
	DTNs     = "2006-01-02 15:04:05.999999999"
	DateOnly = "2006-01-02"
	DMilli   = "2006-01-02.000"
	DMicro   = "2006-01-02.000000"
	DNano    = "2006-01-02.000000000"
	TimeOnly = "15:04:05"

	// 补强布局 (归一化复用版)
	DTFull         = "2006-01-02 15:04:05.999999999 -0700 MST"
	DTTZ           = "2006-01-02 15:04:05-07:00"
	DTTZShort      = "2006-01-02 15:04:05-07"
	DTISO          = "2006-01-02T15:04:05-07:00"
	DTPMMST        = "2006-01-02 15:04:05 PM MST"
	DTPMShortMST   = "2006-01-02 15:04:05PM MST"
	RFC3339Space   = "2006-01-02 15:04:05Z07:00"
	RFC3339NsSpace = "2006-01-02 15:04:05.999999999Z07:00"

	// 紧凑与特殊格式
	DTCompact     = "20060102150405"
	DCompact      = "20060102"
	TimeCompact   = "150405"
	DTVeryShort   = "2006-1-2 15:4:5"
	DTShort       = "2006-1-2 15:4"
	DOnlyShort    = "2006-1-2"
	TimeVeryShort = "15:4:5"
	TimeShort     = "15:4"
	MonthD        = "1-2"
	YearOnly      = "2006"

	// 从 Carbon 补强的高价值布局
	DTCompactTZ     = "20060102150405-07:00"
	DTCompactZ      = "20060102150405Z07:00"
	DTCompactMilli  = "20060102150405.000"
	DHourShort      = "2006-1-2 15"
	HourOnly        = "15"
	DMonth          = "2006-1"
	FormattedD      = "Jan 2, 2006"
	FormattedDayD   = "Mon, Jan 2, 2006"
	DayDateTime     = "Mon, Jan 2, 2006 3:04 PM"
	Cookie          = "Monday, 02-Jan-2006 15:04:05 MST"
	Http            = "Mon, 02 Jan 2006 15:04:05 GMT"
	RFC1036         = "Mon, 02 Jan 06 15:04:05 -0700"
	RFC7231         = "Mon, 02 Jan 2006 15:04:05 MST"
	TimeTZShort     = "15:04:05-07"
	DTFullVeryShort = "2006-1-2 15:4:5 -0700 MST"
	DTNsVeryShort   = "2006-1-2 15:4:5.999999999"

	// ISO8601 家族
	ISO8601       = "2006-01-02T15:04:05-07:00"
	ISO8601Ns     = "2006-01-02T15:04:05.999999999-07:00"
	ISO8601Zulu   = "2006-01-02T15:04:05Z"
	ISO8601ZuluNs = "2006-01-02T15:04:05.999999999Z"
)

var formats = []string{
	// 1. 核心级联 (得益于归一化与 .999 机制，三行即可覆盖 95% 场景)
	DT, DTNs, DateOnly,

	// 2. 标准协议与归一化变体
	RFC3339Space, RFC3339NsSpace, DTFull, DTTZ, DTTZShort,
	ISO8601, ISO8601Ns, ISO8601Zulu, ISO8601ZuluNs,

	// 3. 特殊与 PM 格式
	DTPMMST, DTPMShortMST, DCompact,
	DTCompact, DTCompactTZ, DTCompactZ, DTCompactMilli,

	// 4. 标准库命名格式
	RFC1123Z, RFC1123, RFC822Z, RFC822,
	ANSIC, UnixD, RubyD, RFC850, Kitchen,

	// 5. 时间戳样式
	StampNs, Stamp,

	// 6. 宽松与补全格式
	DTVeryShort, DTShort, DOnlyShort, DHourShort, TimeOnly,
	TimeVeryShort, TimeShort, TimeTZShort, TimeCompact,
	DMonth, MonthD, YearOnly, HourOnly,
}

var (
	buckets = map[int][]string{}
	Formats []string // 公开给用户添加自定义格式
)

func init() {
	for _, layout := range formats {
		// 判定变长因子：
		// 1. 纳秒系列 (.999)
		// 2. 文本月份/星期 (Jan, Mon)
		// 3. 时区 (MST)
		// 4. 宽松数字格式 (不带前导零的 1, 2, 3, 4, 5, 15)
		dynamic := strings.Contains(layout, ".999") ||
			strings.Contains(layout, "Jan") ||
			strings.Contains(layout, "Mon") ||
			strings.Contains(layout, "MST") ||
			strings.Contains(layout, "2006-1") ||
			strings.Contains(layout, "15:4") ||
			layout == "15" || layout == "1-2" || layout == "2006"

		if dynamic {
			Formats = append(Formats, layout)
			if !strings.Contains(layout, ".999") {
				buckets[len(layout)] = append(buckets[len(layout)], layout)
				continue
			}

			dotIdx := strings.Index(layout, ".")
			for i := dotIdx; i <= len(layout); i++ {
				buckets[i] = append(buckets[i], layout)
			}

			continue
		}

		buckets[len(layout)] = append(buckets[len(layout)], layout)
	}
}

// --- 格式化时间 ---

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(DT)+2)
	b = append(b, '"')
	b = t.time.AppendFormat(b, DT)
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
	return DT
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
		s = strings.ReplaceAll(s[:idx], ".", "-") + " " + s[idx+1:]
	} else if dotCount := strings.Count(s, "."); dotCount >= 2 {
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

	// 3. 精准分桶搜索 (O(1) 效率)
	if list, ok := buckets[len(s)]; ok {
		for _, layout := range list {
			if t, err := time.ParseInLocation(layout, s, l); err == nil {
				return Aeon(t), nil
			}
		}
	}

	// 4. 变长布局特快路径 (处理星期、月份名称等非固定长度布局)
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

// --- 内置的格式化时间 ---

type formatDateTimeMilli string

func (formatDateTimeMilli) Layout() string { return DTMilli }

type DateTimeMilli = F[formatDateTimeMilli]
