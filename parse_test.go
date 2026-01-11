package aeon

import (
	"strings"
	"testing"
	"time"
)

func TestParseE(t *testing.T) {
	// 设置默认时区为 UTC 以便测试结果幂等
	DefaultTimeZone = time.UTC
	defer func() { DefaultTimeZone = time.Local }()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// 1. 核心格式 (DT, D)
		{"Standard DT", "2020-08-05 13:14:15", "2020-08-05 13:14:15", false},
		{"Standard D", "2020-08-05", "2020-08-05 00:00:00", false},

		// 2. 归一化测试 (Slash, Dot, T)
		{"Slash Normalization", "2020/08/05 13:14:15", "2020-08-05 13:14:15", false},
		{"Dot Normalization (Date)", "2020.08.05", "2020-08-05 00:00:00", false},
		{"T Normalization", "2020-08-05T13:14:15", "2020-08-05 13:14:15", false},
		{"Complex Normalization", "2020.08.05T13:14:15", "2020-08-05 13:14:15", false},

		// 3. 子秒精度 (.999 机制)
		{"Milli Precision", "2020-08-05 13:14:15.123", "2020-08-05 13:14:15.123", false},
		{"Micro Precision", "2020-08-05 13:14:15.123456", "2020-08-05 13:14:15.123456", false},
		{"Nano Precision", "2020-08-05 13:14:15.123456789", "2020-08-05 13:14:15.123456789", false},
		{"Dot Date with Nano", "2020.08.05 13:14:15.123456789", "2020-08-05 13:14:15.123456789", false},

		// 4. 动态格式 (Weekdays, Months, MST)
		{"ANSIC", "Wed Aug  5 13:14:15 2020", "2020-08-05 13:14:15", false},
		{"UnixD", "Wed Aug  5 13:14:15 UTC 2020", "2020-08-05 13:14:15", false},
		{"RFC1123Z", "Wed, 05 Aug 2020 13:14:15 +0000", "2020-08-05 13:14:15", false},
		{"Full with MST", "2020-08-05 13:14:15.999 +0000 UTC", "2020-08-05 13:14:15.999", false},

		// 5. 宽松/紧凑格式
		{"DCompact", "20200805", "2020-08-05 00:00:00", false},
		{"DTShort", "2020-8-5 13:14", "2020-08-05 13:14:00", false},

		// 6. 空值情况 (按设计应返回零值 Time 且无错误)
		{"Empty", "", "", false},
		{"Null", "null", "", false},

		// 7. 错误情况
		{"Invalid", "invalid-time", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseE(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.input == "" || tt.input == "null" {
					if !got.IsZero() {
						t.Errorf("ParseE(%q) should return zero time", tt.input)
					}
					return
				}
				gotStr := got.Format(DTNs)
				if !strings.Contains(gotStr, tt.want) {
					t.Errorf("ParseE(%q) got = %v, want %v", tt.input, gotStr, tt.want)
				}
			}
		})
	}
}

func TestParse_BucketHits(t *testing.T) {
	inputs := []string{
		"2020-08-05",                // len 10
		"2020-08-05 13:14:15",       // len 19
		"2020-08-05 13:14:15.123",   // len 23
		"Wed Aug  5 13:14:15 2020", // len 24 (Dynamic)
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			res := Parse(input)
			if res.IsZero() {
				t.Errorf("Parse(%q) failed to match any bucket or dynamic format", input)
			}
		})
	}
}
