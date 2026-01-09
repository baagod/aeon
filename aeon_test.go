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

func TestTime_Diff(t *testing.T) {
	// 基准时间：2020-08-05 13:14:15
	base := Date(2020, 8, 5, 13, 14, 15, 0, time.UTC)

	tests := []struct {
		name string
		t    Time
		u    Time
		unit string
		abs  bool
		want float64
	}{
		// 月差
		{"月差：1月后", Date(2020, 9, 5, 13, 14, 15, 0, time.UTC), base, "M", false, 1},
		{"月差：1月前", Date(2020, 7, 5, 13, 14, 15, 0, time.UTC), base, "M", false, -1},

		// 日差
		{"日差：1天后", Date(2020, 8, 6, 13, 14, 15, 0, time.UTC), base, "d", false, 1},
		{"日差：1天前", Date(2020, 8, 4, 13, 14, 15, 0, time.UTC), base, "d", false, -1},

		// 时差
		{"时差：1小时后", Date(2020, 8, 5, 14, 14, 15, 0, time.UTC), base, "h", false, 1},
		{"时差：1小时前", Date(2020, 8, 5, 12, 14, 15, 0, time.UTC), base, "h", false, -1},

		// 分差
		{"分差：1分钟后", Date(2020, 8, 5, 13, 15, 15, 0, time.UTC), base, "m", false, 1},
		{"分差：1分钟前", Date(2020, 8, 5, 13, 13, 15, 0, time.UTC), base, "m", false, -1},

		// 秒差
		{"秒差：1秒后", Date(2020, 8, 5, 13, 14, 16, 0, time.UTC), base, "s", false, 1},
		{"秒差：1秒前", Date(2020, 8, 5, 13, 14, 14, 0, time.UTC), base, "s", false, -1},
		{"秒差：绝对值", Date(2020, 8, 5, 13, 14, 14, 0, time.UTC), base, "s", true, 1},
	}

	for _, tt := range tests {
		got := tt.t.Diff(tt.u, tt.unit, tt.abs)
		if got != tt.want {
			t.Errorf("%s: Diff(%q, abs=%v) = %v, want %v", tt.name, tt.unit, tt.abs, got, tt.want)
		}
	}

	// 年差单独测试（因闰年/非闰年天数差异，结果为浮点数）
	t.Run("年差精度测试", func(t *testing.T) {
		t1 := Date(2021, 8, 5, 13, 14, 15, 0, time.UTC)
		t2 := Date(2019, 8, 5, 13, 14, 15, 0, time.UTC)

		// 2020 是闰年(366天)，2021 是平年(365天)，存在微小偏差
		y1 := t1.Diff(base, "y")
		if y1 < 0.99 || y1 > 1.01 {
			t.Errorf("年差：1年后 = %v, 期望约为 1", y1)
		}

		y2 := t2.Diff(base, "y")
		if y2 > -0.99 || y2 < -1.01 {
			t.Errorf("年差：1年前 = %v, 期望约为 -1", y2)
		}

		// 绝对值测试
		y3 := t2.Diff(base, "y", true)
		if y3 < 0.99 || y3 > 1.01 {
			t.Errorf("年差：绝对值 = %v, 期望约为 1", y3)
		}
	})
}

func assert(t *testing.T, actual Time, expected string, msg string) {
	t.Helper()
	if actual.String() != expected {
		t.Errorf("%s, got [%s], want [%s]", msg, actual, expected)
	}
}
