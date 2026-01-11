package aeon

import "time"

var (
	zeroArgs = []int{0}
	oneArgs  = []int{1}
	// maxDays 每个月的最大天数
	maxDays = [13]int{1, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
)

const (
	absoluteYears = 292277022400
)

// IsLeapYear 返回 year 是否闰年
func IsLeapYear[T ~int](year T) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DaysIn 返回 y 年天数或 y 年 m 月天数
//
// 各月份的最大天数：
//
//   - 1, 3, 5, 7, 8, 10, 12 月有 31 天；4, 6, 9, 11 月有 30 天；
//   - 平年 2 月有 28 天，闰年 29 天。
func DaysIn[T ~int](y T, m ...T) int {
	if len(m) > 0 {
		if m[0] == 2 && IsLeapYear(y) {
			return 29
		}
		return maxDays[m[0]]
	}

	if IsLeapYear(y) {
		return 366
	}

	return 365
}

// dateToAbsDays 从年月日中返回从绝对纪元到该天的天数 (标准库的实现)
func dateToAbsDays(year int64, month time.Month, day int) uint64 {
	amonth := uint32(month)
	janFeb := uint32(0)
	if amonth < 3 {
		janFeb = 1
	}

	amonth += 12 * janFeb
	y := uint64(year) - uint64(janFeb) + absoluteYears

	ayday := (979*amonth - 2919) >> 5
	century := y / 100
	cyear := uint32(y % 100)
	cday := 1461 * cyear / 4
	centurydays := 146097 * century / 4

	return centurydays + uint64(int64(cday+ayday)+int64(day)-1)
}

// weekday 返回由 days 指定的星期几 (标准库的实现)
func weekday(year int, month int, day int) time.Weekday {
	days := dateToAbsDays(int64(year), time.Month(month), day)
	// 绝对年份的 3 月 1 日 (如 2000 年 3 月 1 日) 是星期三
	return time.Weekday((days + uint64(time.Wednesday)) % 7)
}

// addMonth 计算月溢出
func addMonth(y, m, n int) (int, int) {
	months := m + n
	y += (months - 1) / 12
	m = (months-1)%12 + 1
	if m <= 0 {
		m += 12
		y--
	}
	return y, m
}

// clamp 将值限制在 [minimum, maximum] 范围内
func clamp[T ~int | ~float32 | ~float64](value, minimum, maximum T) T {
	if value < minimum {
		return minimum
	} else if value > maximum {
		return maximum
	}
	return value
}

func minimum[T ~int | ~float32 | ~float64](value, minimum T) T {
	if value < minimum {
		return minimum
	}
	return value
}
