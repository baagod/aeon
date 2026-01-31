package aeon

import (
	"time"
)

const (
	DT      = "2006-01-02 15:04:05"
	DTMilli = "2006-01-02 15:04:05.000"
	DTMicro = "2006-01-02 15:04:05.000000"
	DTNano  = "2006-01-02 15:04:05.000000000"
	DTNs    = "2006-01-02 15:04:05.999999999"
)

// ParseE 解析时间字符串，返回 Time 和 error
func ParseE(s string, loc ...*time.Location) (Time, error) {
	if s = trim(s); s == "" || s == "null" {
		return Time{}, nil
	}

	// 预处理时区，仅支持 Z, ±HH:mm, ±HHmm 三种标准格式。
	l := DefaultTimeZone
	if len(loc) > 0 && loc[0] != nil {
		l = loc[0]
	}

	n := len(s)
	if n > 1 && s[n-1] == 'Z' {
		s, l = s[:n-1], time.UTC
	} else if n >= 6 && s[n-3] == ':' { // ±HH:mm
		sign := int(s[n-6])
		if (sign == '+' || sign == '-') && isDigit2(s, n-5) && isDigit2(s, n-2) {
			offset := p2(s, n-5)*3600 + p2(s, n-2)*60
			s, l = s[:n-6], zOffset(offset*(44-sign))
		}
	} else if n >= 5 { // ±HHmm
		sign := int(s[n-5])
		if (sign == '+' || sign == '-') && isDigit4(s[n-4:]) {
			offset := p2(s, n-4)*3600 + p2(s, n-2)*60
			s, l = s[:n-5], zOffset(offset*(44-sign))
		}
	}

	// 解析时间并返回
	t, err := parseFast(trim(s), l)
	return Time{time: t, weekStarts: DefaultWeekStarts}, err
}

// Parse 解析时间字符串，忽略错误
func Parse(value string, loc ...*time.Location) Time {
	t, _ := ParseE(value, loc...)
	return t
}

// ParseByE 指定布局解析，返回 Time 和 error
func ParseByE(layout string, value string, loc ...*time.Location) (Time, error) {
	l := DefaultTimeZone
	if len(loc) > 0 && loc[0] != nil {
		l = loc[0]
	}
	pt, err := time.ParseInLocation(layout, value, l)
	return Time{time: pt, weekStarts: DefaultWeekStarts}, err
}

// ParseBy 指定布局解析，忽略错误
func ParseBy(layout string, value string, loc ...*time.Location) Time {
	t, _ := ParseByE(layout, value, loc...)
	return t
}

// parseFast 是 Aeon 的 L1 级分流决策树。
// 它通过探测 “特征位（isSep）” 实现对标准 ISO8601 家族的 O(1) 识别。
func parseFast(s string, loc *time.Location) (time.Time, error) {
	n := len(s)
	if n < 1 {
		return time.Time{}, nil
	}

	// 优先检查前 4 位是否为年份数字块
	if isDigit4(s) {
		y := p4(s)
		// 判定进入紧凑大类：长度为4，或者第5位是数字或小数点 (YYYYM... or YYYY.nnn)
		if n == 4 || isDigit(s[4]) || s[4] == '.' {
			return parseCompact(s, n, y, loc)
		}

		// --- 统一基因特征寻址 ---
		// 10位日期 (YYYY?MM?DD)
		if n == 10 && isSep2(s[4], s[7]) {
			return time.Date(y, time.Month(p2(s, 5)), p2(s, 8), 0, 0, 0, 0, loc), nil
		}

		// 16 位日期时间 (YYYY?MM?DD?HH?mm)
		if n == 16 && isSep4(s[4], s[7], s[10], s[13]) {
			return time.Date(y, time.Month(p2(s, 5)), p2(s, 8), p2(s, 11), p2(s, 14), 0, 0, loc), nil
		}

		// 19-23 位日期时间 (YYYY?MM?DD?HH?mm?ss[.SSS])
		if (n == 19 || n == 23) && isSep5(s[4], s[7], s[10], s[13], s[16]) {
			ns, _ := parseNanoseconds(s, n, 19)
			return time.Date(y, time.Month(p2(s, 5)), p2(s, 8), p2(s, 11), p2(s, 14), p2(s, 17), ns, loc), nil
		}

		// --- 通用寻址通道 (变长/异形) ---
		v, ns := parseGeneric(s, 4, 5)
		if v[0] > 0 {
			return time.Date(y, time.Month(max(1, v[0])), max(1, v[1]), v[2], v[3], v[4], ns, loc), nil
		}
	}

	if (n == 8 || n == 12) && s[2] == ':' && s[5] == ':' {
		ns, _ := parseNanoseconds(s, n, 8)
		return time.Date(0, 1, 1, p2(s, 0), p2(s, 3), p2(s, 6), ns, loc), nil
	}

	// 子类 A2：以时间开头 ("13:14..." 或 "2:3")
	if (n >= 3 && s[1] == ':') || (n >= 4 && isDigit(s[1]) && s[2] == ':') {
		v, ns := parseGeneric(s, 0, 3)
		return time.Date(0, 1, 1, v[0], v[1], v[2], ns, loc), nil
	}

	return time.Time{}, nil
}

// parseCompact 负责解析不带分隔符的紧凑格式（如 YYYYMMDD），并支持在 4, 6, 8, 10, 12, 14 位后跟随小数点表示纳秒。
func parseCompact(s string, n int, y int, loc *time.Location) (time.Time, error) {
	if n == 4 {
		return time.Date(y, 1, 1, 0, 0, 0, 0, loc), nil
	}

	m, d, h, mm, sec, ns := 1, 1, 0, 0, 0, 0
	sep, dot := 0, -1 // 确定逻辑偏移与纳秒

	if n > 15 && s[15] == '.' {
		dot = 15
	} else if n > 14 && s[14] == '.' {
		dot = 14
	} else if n > 12 && s[12] == '.' {
		dot = 12
	} else if n > 10 && s[10] == '.' {
		dot = 10
	} else if n > 8 && s[8] == '.' {
		dot = 8
	} else if n > 6 && s[6] == '.' {
		dot = 6
	} else if n > 4 && s[4] == '.' {
		dot = 4
	}

	if n >= 9 && !isDigit(s[8]) && s[8] != '.' {
		sep = 1 // 只有长字符串才需要检查 s[8]
	}

	end := n
	if dot != -1 {
		end = dot
	}

	switch end - sep { // 统一寻址提取日期时间
	case 14:
		sec = p2(s, 12+sep)
		fallthrough
	case 12:
		mm = p2(s, 10+sep)
		fallthrough
	case 10:
		h = p2(s, 8+sep)
		fallthrough
	case 8:
		d = p2(s, 6)
		fallthrough
	case 6:
		m = p2(s, 4)
	}

	ns, _ = parseNanoseconds(s, n, dot)
	return time.Date(y, time.Month(m), d, h, mm, sec, ns, loc), nil
}

// parseGeneric 是通用时间解析器，采用 “贪吃蛇” 方式逐个提取时间组件。
//
// 参数说明：
//   - s: 待解析的时间字符串（如 "2025-01-02 15:04:05"）
//   - i: 起始解析位置（索引），通常传入已解析部分的结束位置
//     例如 parseFast 已提取年份 "2025"（4 字节），则传入 i=4
//   - limit: 最多提取的数字块数量（通常为 5 或 6）
//     提取顺序：月、日、时、分、秒、年（可选）
//
// 返回值说明：
//   - v [6]int: 提取的 6 个数字块数组
//     v[0]: 月份 (1-12)
//     v[1]: 日数 (1-31)
//     v[2]: 小时 (0-23)
//     v[3]: 分钟 (0-59)
//     v[4]: 秒数 (0-59)
//     v[5]: 年份（可选，部分格式可能提取）
//   - ns: 纳秒部分（如果有小数点）
//
// 算法特点：
//  1. 跳过所有非数字字符（分隔符：- : / 空格等）
//  2. 遇到数字开始贪婪提取，直到下一个非数字
//  3. 支持纳秒解析：仅当 id>=2（时间组件）时检测小数点
//  4. ASCII 优化：使用 s[i]-'0' > 9 替代 unicode.IsDigit，速度提升 10 倍
//
// 示例：
//
//	parseGeneric("01-02 15:04:05.123", 0, 6)
//	  → v = [1, 2, 15, 4, 5, 0], ns = 123
//
//	parseGeneric("15:04:05", 0, 3)
//	  → v = [15, 4, 5, 0, 0, 0], ns = 0
func parseGeneric(s string, i, limit int) (v [6]int, ns int) {
	n := len(s)
Loop:
	for id := 0; id < limit; id++ {
		// 1. 寻找组件起始位：跳过所有杂质字符，并在此探测纳秒。
		for ; i < n && s[i]-'0' > 9; i++ {
			if s[i] == '.' && id >= 2 {
				break Loop
			}
		}

		if i >= n {
			break
		}

		// 2. 极速提取数字块：利用指针 i 贪婪消费所有数字
		c := int(s[i] - '0')
		for i++; i < n; i++ {
			d := s[i] - '0'
			if d > 9 {
				break // 遇到非数字停止采集
			}
			c = c*10 + int(d)
		}
		v[id] = c
	}

	// 统一纳秒入口
	if i < n && s[i] == '.' {
		ns, i = parseNanoseconds(s, n, i)
	}

	return
}

func parseNanoseconds(s string, n, start int) (ns, i int) {
	if start == -1 || n == start {
		return
	}

	for i = start + 1; i < min(n, start+10); i++ {
		d := s[i] - '0'
		if d > 9 {
			break
		}
		ns = ns*10 + int(d)
	}

	return ns * pow10[10+start-i], i
}
