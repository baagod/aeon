package thru

import (
	"testing"
	"time"
)

func TestTime_Second(t *testing.T) {
	tm := time.Date(2023, 1, 1, 12, 30, 45, 123456789, time.UTC)
	at := New(tm)

	tests := []struct {
		n    int
		want int
	}{
		{0, 45}, // 默认秒
		{1, 1},  // 1.2... -> 1
		{3, 123},
		{6, 123456},
		{9, 123456789},
	}

	for _, tt := range tests {
		got := at.Second(tt.n)
		if got != tt.want {
			t.Errorf("Second(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestTime_Unix(t *testing.T) {
	// 2023-01-01 12:30:45.123456789 UTC
	// Unix: 1672576245
	// UnixMilli: 1672576245123
	// UnixMicro: 1672576245123456
	// UnixNano: 1672576245123456789
	tm := time.Date(2023, 1, 1, 12, 30, 45, 123456789, time.UTC)
	at := New(tm)

	tests := []struct {
		n    int
		want int64
	}{
		{0, 1672576245},
		{3, 1672576245123},
		{6, 1672576245123456},
		{9, 1672576245123456789},
	}

	for _, tt := range tests {
		got := at.Unix(tt.n)
		if got != tt.want {
			t.Errorf("Unix(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func assert(t *testing.T, actual Time, expected string, msg string) {
	t.Helper()
	if actual.String() != expected {
		t.Errorf("%s, got [%s], want [%s]", msg, actual, expected)
	}
}
