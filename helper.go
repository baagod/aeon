package aeon

import (
	"cmp"
	"encoding/binary"
	"sync"
	"time"
	"unsafe"
)

const (
	absoluteYears = 292277022400
)

var (
	bt       [256]byte
	zeroArgs = []int{0}
	oneArgs  = []int{1}

	// maxDays 每个月的最大天数
	maxDays = [13]int{1, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// pow19 预定义的 10 的幂次方表，用于高精度计算
	pow19 = [...]int64{
		1, 10, 100, 1000, 10000, 100000, 1e6, 1e7, 1e8, 1e9,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
	}

	// pow10 纳秒精度缩放查找表
	pow10 = [10]int{1, 10, 100, 1000, 10000, 100000, 1e6, 1e7, 1e8, 1e9}

	// locCache 缓存常用的时区偏移，实现 0 Alloc 时区解析
	locCache sync.Map
)

const (
	kDigit  = 1 << iota // 0x01: 0-9
	kSep                // 0x02: 分隔符
	kAlpha              // 0x04: a-z, A-Z
	kTrim               // 0x08: 可 trim
	kSign               // 0x10: +, -
	kDot                // 0x20: .
	kNotSep = kDigit | kDot
)

func init() {
	for i := 0; i < 256; i++ {
		c := byte(i)
		if c-'0' <= 9 { // 数字
			bt[i] |= kDigit
		}

		if (c|0x20)-'a' <= 'z'-'a' { // 字母
			bt[i] |= kAlpha
		}

		switch c {
		case '-', ':', '/', 'T', ' ': // 分隔符
			bt[i] |= kSep
		}

		switch c {
		case ' ', '"', '\n', '\r', '\t': // trim
			bt[i] |= kTrim
		}

		if c == '-' || c == '+' { // 时区
			bt[i] |= kSign
		}

		if c == '.' { // 纳秒
			bt[i] |= kDot
		}
	}
}

// dateToAbsDays 从年月日中返回从绝对纪元到该天的天数 (参考标准库实现)
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

// weekday 返回指定年月日的周几 (参考标准库实现)
func weekday(y int, m int, d int) time.Weekday {
	days := dateToAbsDays(int64(y), time.Month(m), d)
	// 绝对纪元的基准日是星期三
	return time.Weekday((days + uint64(time.Wednesday)) % 7)
}

// addMonth 计算月份增加或减少后的年月（处理年进位/借位）
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

// zOffset 根据秒数偏移量获取 *time.Location。
// 它优先从缓存中读取，若无缓存则创建并存入，从而避免 time.FixedZone 带来的堆分配。
func zOffset(offset int) *time.Location {
	if offset == 0 {
		return time.UTC
	}
	if loc, ok := locCache.Load(offset); ok {
		return loc.(*time.Location)
	}
	newLoc := time.FixedZone("", offset)
	locCache.Store(offset, newLoc)
	return newLoc
}

func trim(s string) string {
	end := len(s)
	if end == 0 {
		return s
	}

	if (bt[s[0]]|bt[s[end-1]])&kTrim == 0 {
		return s
	}

	var start int
	for ; start < end && bt[s[start]]&kTrim != 0; start++ {
	}

	for ; end > start && bt[s[end-1]]&kTrim != 0; end-- {
	}

	return s[start:end]
}

// p2 从字符串位置 i 开始解析 2 位数字（十进制）。
// 使用位运算技巧：ASCII 数字 '0'-'9' 的低 4 位即为数字值（'0'=0, ..., '9'=9）。
// 示例：p2("2024", 0) → 20, p2("2024", 2) → 24
//
// ⚡ BCE (边界检查消除)：通过访问 s[i+1] 告诉编译器字符串至少有 2 个字符，
//
//	从而消除后续 s[i] 和 s[i+1] 访问的边界检查，提升性能。
func p2(s string, i int) int {
	_ = s[i+1]
	return int(s[i]&0xF)*10 + int(s[i+1]&0xF)
}

// p4 从字符串开头解析 4 位数字（十进制）。
// 使用 binary.LittleEndian.Uint32 一次性读取 4 个字节，然后通过位运算提取每个字节的低 4 位。
//
// 字节布局（小端序）：
//
//	s[0] (最低位) s[1] s[2] s[3] (最高位)
//	v&0xF         v>>8&0xF v>>16&0xF v>>24&0xF
//
// 示例：p4("2024") → 2024, p4("1530") → 1530
func p4(s string) int {
	v := binary.LittleEndian.Uint32(unsafe.Slice(unsafe.StringData(s), 4))
	return int(v&0xF)*1000 + int(v>>8&0xF)*100 + int(v>>16&0xF)*10 + int(v>>24&0xF)
}

// isDigit 检查单个字符 c 是否为 ASCII 数字 '0'-'9'。
//
// 原理：ASCII 数字 '0'-'9' 是连续的（48-57），使用减法技巧：
//
//	c-'0' <= 9 等价于 '0' <= c <= '9'
//
// 示例：isDigit('5') → true, isDigit('a') → false
func isDigit(c byte) bool {
	return c-'0' <= 9
}

// isDigit2 检查从位置 i 开始的 2 个字符是否都是 ASCII 数字 '0'-'9'。
//
// 位运算技巧：
// 1. 一次性读取 2 个字节为 uint16（LittleEndian）
// 2. 检查高 4 位：'0'-'9' 的 ASCII 高 4 位都是 0x3（0x30-0x39）
//   - v&0xF0F0 == 0x3030 → 确保高 4 位都是 0x3
//
// 3. 检查低 4 位：通过加 0x0606 检查是否溢出到下一个数字范围
//   - (v+0x0606)&0xF0F0 == 0x3030 → 确保低 4 位在 0-9 之间
//
// 示例：isDigit2("12", 0) → true, isDigit2("1a", 0) → false
func isDigit2(s string, i int) bool {
	if len(s) < i+2 {
		return false
	}
	v := binary.LittleEndian.Uint16(unsafe.Slice(unsafe.StringData(s[i:]), 2))
	return (v&0xF0F0 == 0x3030) && ((v+0x0606)&0xF0F0 == 0x3030)
}

// isDigit4 检查字符串开头的 4 个字符是否都是 ASCII 数字 '0'-'9'。
//
// 位运算技巧（扩展自 isDigit2）：
// 1. 一次性读取 4 个字节为 uint32（LittleEndian）
// 2. 检查高 4 位：确保每个字节的高 4 位都是 0x3
//   - v&0xF0F0F0F0 == 0x30303030
//
// 3. 检查低 4 位：通过加 0x06060606 检查是否溢出
//   - (v+0x06060606)&0xF0F0F0F0 == 0x30303030
//
// 示例：isDigit4("2024") → true, isDigit4("202a") → false
func isDigit4(s string) bool {
	if len(s) < 4 {
		return false
	}
	v := binary.LittleEndian.Uint32(unsafe.Slice(unsafe.StringData(s), 4))
	return (v&0xF0F0F0F0 == 0x30303030) && ((v+0x06060606)&0xF0F0F0F0 == 0x30303030)
}

// isSep2 判定两个位置是否都不是数字且不是小数点
func isSep2(c1, c2 byte) bool {
	return (bt[c1]|bt[c2])&kNotSep == 0
}

// isSep4 判定五个位置是否都不是数字且不是小数点
func isSep4(c1, c2, c3, c4 byte) bool {
	return (bt[c1]|bt[c2]|bt[c3]|bt[c4])&kNotSep == 0
}

// isSep5 判定五个位置是否都不是数字且不是小数点
func isSep5(c1, c2, c3, c4, c5 byte) bool {
	return (bt[c1]|bt[c2]|bt[c3]|bt[c4]|bt[c5])&kNotSep == 0
}

// max 返回两个整数中的较大值。
func max[T cmp.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func min[T cmp.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// clamp 将泛型数值限制在指定的最小值和最大值区间内
func clamp[T ~int | ~float32 | ~float64](value, minimum, maximum T) T {
	if value < minimum {
		return minimum
	} else if value > maximum {
		return maximum
	}
	return value
}
