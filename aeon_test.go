package aeon

import (
	"encoding/json"
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

func TestTime_SubSecond(t *testing.T) {
	tm := time.Date(2023, 1, 1, 12, 30, 45, 123456789, time.UTC)
	at := New(tm)

	if at.Second(3) != 123 {
		t.Errorf("Milli() = %d, want 123", at.Second(3))
	}
	if at.Second(6) != 123456 {
		t.Errorf("Micro() = %d, want 123456", at.Second(6))
	}
	if at.Second(9) != 123456789 {
		t.Errorf("Nano() = %d, want 123456789", at.Second(9))
	}
}

func TestTime_SubSecondCascading(t *testing.T) {
	base := Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	t.Run("Milli", func(t *testing.T) {
		// StartMilli(500) -> 12:00:00.500
		assert(t, base.StartMilli(500), "2024-01-01 12:00:00.5", "StartMilli(500)")
		// EndMilli(500) -> 12:00:00.500999999
		assert(t, base.EndMilli(500), "2024-01-01 12:00:00.500999999", "EndMilli(500)")
	})

	t.Run("Micro", func(t *testing.T) {
		// StartMicro(500000) -> 12:00:00.500
		assert(t, base.StartMicro(500000), "2024-01-01 12:00:00.5", "StartMicro(500000)")
		// EndMicro(500000) -> 12:00:00.500000999
		assert(t, base.EndMicro(500000), "2024-01-01 12:00:00.500000999", "EndMicro(500000)")
	})

	t.Run("Nano", func(t *testing.T) {
		// StartNano(500000000) -> 12:00:00.500
		assert(t, base.StartNano(500000000), "2024-01-01 12:00:00.5", "StartNano(500000000)")
		// EndNano(500000000) -> 12:00:00.5
		assert(t, base.EndNano(500000000), "2024-01-01 12:00:00.5", "EndNano(500000000)")
	})

	t.Run("Cross-unit Cascading", func(t *testing.T) {
		// StartMilli(100, 500) -> 100ms + 500us = 100.5ms
		assert(t, base.StartMilli(100, 500), "2024-01-01 12:00:00.1005", "StartMilli(100, 500)")
	})
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

// ---- F 泛型测试 ----

// 自定义格式类型
type testDateFormat struct{}

func (testDateFormat) Layout() string { return DateOnly }

type testDateTimeFormat struct{}

func (testDateTimeFormat) Layout() string { return DateTime }

func TestFormatted_MarshalJSON(t *testing.T) {
	tm := Date(2025, 1, 10, 14, 30, 45, 0, time.UTC)

	tests := []struct {
		name string
		time any
		want string
	}{
		{
			name: "DateOnly格式",
			time: F[testDateFormat]{Time: tm},
			want: `"2025-01-10"`,
		},
		{
			name: "DateTime格式",
			time: F[testDateTimeFormat]{Time: tm},
			want: `"2025-01-10 14:30:45"`,
		},
		{
			name: "零值返回null",
			time: F[testDateFormat]{},
			want: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.time)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestFormatted_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		format any
		want   string
	}{
		{
			name:   "DateOnly格式",
			input:  `"2025-01-10"`,
			format: &F[testDateFormat]{},
			want:   "2025-01-10 00:00:00",
		},
		{
			name:   "DateTime格式",
			input:  `"2025-01-10 14:30:45"`,
			format: &F[testDateTimeFormat]{},
			want:   "2025-01-10 14:30:45",
		},
		{
			name:   "null值",
			input:  `null`,
			format: &F[testDateFormat]{},
			want:   "0001-01-01 00:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tt.input), tt.format)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			var got string
			switch v := tt.format.(type) {
			case *F[testDateFormat]:
				got = v.String()
			case *F[testDateTimeFormat]:
				got = v.String()
			}

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestFormatted_StructField(t *testing.T) {
	type User struct {
		Name      string                `json:"name"`
		Birthday  F[testDateFormat]     `json:"birthday"`
		CreatedAt F[testDateTimeFormat] `json:"created_at"`
	}

	user := User{
		Name:      "张三",
		Birthday:  F[testDateFormat]{Time: Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)},
		CreatedAt: F[testDateTimeFormat]{Time: Date(2025, 1, 10, 14, 30, 45, 0, time.UTC)},
	}

	// 序列化
	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want := `{"name":"张三","birthday":"1990-05-15","created_at":"2025-01-10 14:30:45"}`
	if string(data) != want {
		t.Errorf("Marshal:\ngot  %s\nwant %s", data, want)
	}

	// 反序列化
	var user2 User
	err = json.Unmarshal(data, &user2)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if user2.Name != user.Name {
		t.Errorf("Name: got %s, want %s", user2.Name, user.Name)
	}
	if user2.Birthday.String() != user.Birthday.String() {
		t.Errorf("Birthday: got %s, want %s", user2.Birthday, user.Birthday)
	}
	if user2.CreatedAt.String() != user.CreatedAt.String() {
		t.Errorf("CreatedAt: got %s, want %s", user2.CreatedAt, user.CreatedAt)
	}
}
