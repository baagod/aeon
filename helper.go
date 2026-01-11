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

// week1day 返回 y年m月d日 的星期几 (标准库的实现)
func weekday(y int, m int, d int) time.Weekday {
	days := dateToAbsDays(int64(y), time.Month(m), d)
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
