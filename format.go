package aeon

import (
	"database/sql/driver"
	"time"
)

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

// --- 内置的格式化时间 ---

type formatDateTimeMilli string

func (formatDateTimeMilli) Layout() string { return DTMilli }

type DateTimeMilli = F[formatDateTimeMilli]
