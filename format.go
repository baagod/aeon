package thru

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	Layout      = "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	Stamp       = "Jan _2 15:04:05"
	StampMilli  = "Jan _2 15:04:05.000"
	StampMicro  = "Jan _2 15:04:05.000000"
	StampNano   = "Jan _2 15:04:05.000000000"
	DateTime    = "2006-01-02 15:04:05"
	DateOnly    = "2006-01-02"
	TimeOnly    = "15:04:05"
)

const (
	dateonly = `\d{4}(-\d{2}){2}`
	datetime = `(\d{2}:){2}\d{2}(\.\d{1,9})?`
	mst      = `[A-Z]{3,4}([+\-]\d{1,2})?`
	z0700    = `[+\-]\d{4}`
)

var patterns = map[string]*regexp.Regexp{
	"2006":                                re(`\d{4}`),
	"15:04":                               re(`\d{2}:\d{2}`),
	"3:04PM":                              re(`\d{1,2}:\d{2}[AP]M`),
	"2006-01":                             re(`\d{4}-\d{2}`),
	"15:04:05":                            re(datetime),
	"2006-01-02":                          re(dateonly),
	"2006-01-02 15":                       re(dateonly + ` \d{2}`),
	"Jan _2 15:04:05":                     re(`(?i)[a-z]{3} \d{1,2} ` + datetime),
	"2006-01-02 15:04":                    re(dateonly + ` \d{2}:\d{2}`),
	"2006-01-02 15:04:05":                 re(dateonly + " " + datetime),
	"02 Jan 06 15:04 MST":                 re(`\d{2} (?i)[a-z]{3} \d{2} \d{2}:\d{2} %s`, mst),
	"02 Jan 06 15:04 -0700":               re(`\d{2} (?i)[a-z]{3} \d{2} \d{2}:\d{2} %s`, z0700),
	"Mon Jan _2 15:04:05 2006":            re(`(?i)([a-z]{3} ){2}\d{1,2} %s \d{4}`, datetime),
	"01-02 03:04:05PM '06 -0700":          re(`\d{2}-\d{2} %s[AP]M '\d{2} %s`, datetime, z0700),
	"Mon Jan _2 15:04:05 MST 2006":        re(`(?i)([a-z]{3} ){2}\d{1,2} %s %s \d{4}`, datetime, mst),
	"Mon, 02 Jan 2006 15:04:05 MST":       re(`(?i)[a-z]{3}, \d{2} (?i)[a-z]{3} \d{4} %s %s`, datetime, mst),
	"Mon Jan 02 15:04:05 -0700 2006":      re(`(?i)([a-z]{3} ){2}\d{2} %s %s \d{4}`, z0700, datetime),
	"Monday, 02-Jan-06 15:04:05 MST":      re(`(?i)(Mon|Tues|Wednes|Thurs|Fri|Satur|Sun)day, \d{2}-(?i)[a-z]{3}-\d{2} %s %s`, datetime, mst),
	"Mon, 02 Jan 2006 15:04:05 -0700":     re(`(?i)[a-z]{3}, \d{2} (?i)[a-z]{3} \d{4} %s %s`, datetime, z0700),
	"2006-01-02T15:04:05.999999999Z07:00": re(`%sT%s(Z|[+\-]\d{2}:\d{2})`, dateonly, datetime),
}

// --- 格式化时间 ---

// MarshalJSON 将 t 转为 JSON 字符串时调用
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}

	b := make([]byte, 0, len(DateTime)+2)
	b = append(b, '"') // 添加开头双引号
	b = t.time.AppendFormat(b, DateTime)
	return append(b, '"'), nil // 添加结尾双引号
}

// UnmarshalJSON 将 JSON 字符串转为 t 时调用
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	*t, err = ParseE(string(b), t.Location())
	return
}

// Scan 由 DB 转到 Go 时调用
func (t *Time) Scan(value any) (err error) {
	switch v := value.(type) {
	case time.Time:
		*t = New(v)
	case string:
		*t, err = ParseE(v, t.Location())
	default:
		*t = New(time.Time{})
	}
	return
}

// Value 由 Go 转到 DB 时调用
func (t Time) Value() (v driver.Value, err error) {
	if !t.IsZero() {
		v = t.time
	}
	return
}

// --- 自定义 JSON 格式 ---

// Formatter 定义 JSON 序列化的时间格式
type Formatter interface {
	Layout() string
}

// F 是支持自定义 JSON 格式的时间包装器。
//
// 用户需自定义格式类型并实现 Formatter 接口：
//
//	type MyDate struct{}
//	func (MyDate) Layout() string { return "2006-01-02" }
//
//	type User struct {
//	    Birthday aeon.F[MyDate] `json:"birthday"`
//	}
type F[T Formatter] struct {
	Time
	m T
}

func (f *F[T]) layout() string {
	if layout := f.m.Layout(); layout != "" {
		return layout
	}
	return DateTime
}

// MarshalJSON 将时间序列化为 JSON 字符串
func (f F[T]) MarshalJSON() ([]byte, error) {
	if f.IsZero() {
		return []byte("null"), nil
	}

	layout := f.layout()

	b := make([]byte, 0, len(layout)+2)
	b = append(b, '"') // 添加开头双引号
	b = f.time.AppendFormat(b, layout)

	return append(b, '"'), nil // 添加结尾双引号
}

// UnmarshalJSON 从 JSON 字符串反序列化时间
func (f *F[T]) UnmarshalJSON(b []byte) (err error) {
	f.Time, err = ParseE(string(b), f.Location())
	return
}

// Scan 实现 sql.Scanner 接口
func (f *F[T]) Scan(value any) (err error) {
	switch v := value.(type) {
	case time.Time:
		f.Time = New(v)
	case string:
		f.Time, err = ParseE(v, f.Location())
	default:
		f.Time = New(time.Time{})
	}
	return
}

// Value 实现 driver.Valuer 接口
func (f F[T]) Value() (v driver.Value, err error) {
	if !f.IsZero() {
		v = f.Time.Time()
	}
	return
}

// --- 解析时间 ---

// ParseE 解析 value 并返回它所表示的时间
func ParseE(s string, loc ...*time.Location) (Time, error) {
	s = strings.Trim(strings.TrimSpace(s), `"`)
	if s == "" || s == "null" {
		return Time{}, nil
	}

	var layout string
	s = strings.ReplaceAll(s, "/", "-")

	for k, v := range patterns {
		if v.MatchString(s) {
			layout = k
			break
		}
	}

	if loc == nil {
		loc = append(loc, time.Local)
	}

	pt, err := time.ParseInLocation(layout, s, loc[0])
	return Time{time: pt, weekStartsAt: DefaultWeekStartsAt}, err
}

// Parse 返回忽略错误的 ParseE()
func Parse(value string, loc ...*time.Location) Time {
	t, _ := ParseE(value, loc...)
	return t
}

func ParseByLayoutE(layout string, value string, loc ...*time.Location) (Time, error) {
	if loc == nil {
		loc = append(loc, time.Local)
	}
	pt, err := time.ParseInLocation(layout, value, loc[0])
	return Time{time: pt, weekStartsAt: DefaultWeekStartsAt}, err
}

func ParseByLayout(layout string, value string, loc ...*time.Location) Time {
	t, _ := ParseByLayoutE(layout, value, loc...)
	return t
}

func (t Time) String() string {
	if ns := t.time.Nanosecond(); ns == 0 {
		return t.time.Format(DateTime)
	}
	return t.time.Format("2006-01-02 15:04:05.000000000")
}

func (t Time) Format(layout string) string {
	return t.time.Format(layout)
}

func re(s string, a ...any) *regexp.Regexp {
	return regexp.MustCompile("^" + fmt.Sprintf(s, a...) + "$")
}
