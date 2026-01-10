package thru

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const (
	ANSIC        = "Mon Jan _2 15:04:05 2006"
	UnixDate     = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate     = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822       = "02 Jan 06 15:04 MST"
	RFC822Z      = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850       = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123      = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z     = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339      = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano  = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen      = "3:04PM"
	Stamp        = "Jan _2 15:04:05"
	StampMilli   = "Jan _2 15:04:05.000"
	StampMicro   = "Jan _2 15:04:05.000000"
	StampNano    = "Jan _2 15:04:05.000000000"
	DateTime     = "2006-01-02 15:04:05"
	DateTimeNano = "2006-01-02 15:04:05.999999999"
	DateOnly     = "2006-01-02"
	TimeOnly     = "15:04:05"
)

var Formats = []string{
	// 1. 现代 API 最常用格式 (优先匹配长格式)
	DateTime,     // "2006-01-02 15:04:05"
	DateTimeNano, // "2006-01-02 15:04:05.999999999"
	DateOnly,     // "2006-01-02"
	RFC3339,      // "2006-01-02T15:04:05Z07:00"
	RFC3339Nano,  // "2006-01-02T15:04:05.999999999Z07:00"
	"2006-01-02 15:04:05.999999999 -0700 MST",

	// 2. 带点号的变体 (点号不可替换)
	"2006.01.02 15:04:05.999999999",
	"2006.01.02 15:04:05",
	"2006.01.02",
	"2006.1.2 15:04:05",
	"2006.1.2",

	// 3. 标准库命名格式
	RFC1123Z,
	RFC1123,
	RFC822Z,
	RFC822,
	ANSIC,
	UnixDate,
	RubyDate,
	RFC850,
	Kitchen,

	// 4. 时间戳系列
	StampNano,
	StampMicro,
	StampMilli,
	Stamp,

	// 5. now 补全的紧凑型与宽松格式
	"20060102",
	"2006-1-2 15:4:5",
	"2006-1-2 15:4",
	"2006-1-2 15",
	"2006-1-2",
	TimeOnly, // "15:04:05"
	"15:4:5",
	"15:4",
	"15",
	"1-2",  // 月-日
	"2006", // 年
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

func (t Time) MarshalText() ([]byte, error) {
	if t.IsZero() {
		return []byte(""), nil // 文本格式零值通常为空
	}
	return []byte(t.String()), nil
}

func (t *Time) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}

// Scan 由 DB 转到 Go 时调用
func (t *Time) Scan(value any) (err error) {
	switch v := value.(type) {
	case time.Time:
		*t = New(v)
	case string:
		*t, err = ParseE(v, t.Location())
	default:
		*t = New()
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

// MarshalText 实现 encoding.TextMarshaler 接口
func (f F[T]) MarshalText() ([]byte, error) {
	if f.IsZero() {
		return []byte(""), nil
	}
	// 使用 F.layout() 而不是 Time.String()
	return []byte(f.Format(f.layout())), nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler 接口
func (f *F[T]) UnmarshalText(data []byte) error {
	return f.UnmarshalJSON(data)
}

// Scan 实现 sql.Scanner 接口
func (f *F[T]) Scan(value any) (err error) {
	switch v := value.(type) {
	case time.Time:
		f.Time = New(v)
	case string:
		f.Time, err = ParseE(v, f.Location())
	default:
		f.Time = New()
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

	// 统一分隔符优化性能与兼容性
	s = strings.ReplaceAll(s, "/", "-")

	l := time.Local
	if len(loc) > 0 && loc[0] != nil {
		l = loc[0]
	}

	for _, layout := range Formats {
		if t, err := time.ParseInLocation(layout, s, l); err == nil {
			return New(t), nil // 解析成功，立即返回
		}
	}

	return Time{}, fmt.Errorf("aeon 无法解析时间字符串: %q", s)
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

// -- 时间格式 ---

func (t Time) Format(layout string) string {
	return t.time.Format(layout)
}

func (t Time) AppendFormat(b []byte, layout string) []byte {
	return t.time.AppendFormat(b, layout)
}

func (t Time) String() string {
	return t.time.Format(DateTimeNano)
}
