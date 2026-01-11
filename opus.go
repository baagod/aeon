package aeon

import (
	"time"
)

var (
	stdSeq      = []Unit{Century, Decade, Year, Month, Day, Hour, Minute, Second, Millisecond, Microsecond, Nanosecond} // 标准流
	quarterSeq  = []Unit{Quarter, Month, Day, Hour, Minute, Second, Millisecond, Microsecond, Nanosecond}               // 季度流
	weekSeq     = []Unit{Week, Weekday, Hour, Minute, Second, Millisecond, Microsecond, Nanosecond}                     // 月周流
	yearWeekSeq = []Unit{YearWeek, Weekday, Hour, Minute, Second, Millisecond, Microsecond, Nanosecond}                 // 年周流
)

func (u Unit) seq() []Unit {
	switch u {
	case Quarter:
		return quarterSeq
	case Week:
		return weekSeq
	case YearWeek:
		return yearWeekSeq
	default:
		if u <= Nanosecond {
			return stdSeq[u:]
		}
	}
	return []Unit{u}
}

func (u Unit) factor() int {
	switch u {
	case Millisecond:
		return 1e6
	case Microsecond:
		return 1e3
	case Nanosecond:
		return 1
	default:
		return 0
	}
}

func apply(f from, c Flag, first bool, u, p Unit, n int, y, m, d, h, mm, s, ns int, w, sw time.Weekday) (int, int, int, int, int, int, int, time.Weekday) {
	switch f {
	case fromAdd:
		return applyOffset(c, u, p, n, y, m, d, h, mm, s, ns, w)
	case fromAbs, fromGoAbs:
		return applyAbs(c, u, p, n, y, m, d, h, mm, s, ns, w, sw)
	case fromRel, fromGoRel:
		return applyRel(c, u, p, n, y, m, d, h, mm, s, ns, w, sw)
	case fromAt, fromGoAt:
		if first {
			return applyAbs(c, u, p, n, y, m, d, h, mm, s, ns, w, sw)
		}
		return applyRel(c, u, p, n, y, m, d, h, mm, s, ns, w, sw)
	default: // fromIn, fromGoIn
		if first {
			return applyRel(c, u, p, n, y, m, d, h, mm, s, ns, w, sw)
		}
		return applyAbs(c, u, p, n, y, m, d, h, mm, s, ns, w, sw)
	}
}

// applyAbs 应用基于上级单位的绝对定位逻辑，它是 StartCentury() 系列方法的核心实现。
//
// 核心模式：
//   - n > 0：定位到第 n 个单元（绝对位置）
//   - n < 0：定位到倒数第 n 个单元
//   - n = 0：保持当前单元位置不变，仅对其下的子级单位进行对齐（StartCentury 系列归零，EndCentury 系列置满）。
//
// 定位规则：
//
//	定位单位（支持跨级定位，强制保护日期）：
//	- Century/Decade/Year/Quarter/Month：允许跨量级跳转。在跳转月份后会自动校正“天”分量
//	  （如 1月31日 跳转到 2月后校正为 28/29日），以保证内部状态 w (Weekday) 的合法性。
//
//	偏移单位（完全支持自然溢出）：
//	- Day/Hour/Minute/Second/Milli/Micro/Nano：作为细粒度偏移，不限制范围，允许产生跨月、跨天等自然时间溢出。
//	- Week：基于当前上下文日期，执行该日期所在周的周内定位。
//
// 参数级联：级联路径由入口方法决定。标准流为：世纪 → 年代 → 年 → 月 → 日 → 时 → 分 → 秒 → 毫秒 → 微秒 → 纳秒
// 季度流：季度 → 月 → 日 ...
// 周流：[Year] → [Year]Week → Weekday ...
func applyAbs(c Flag, u, p Unit, n, y, m, d, h, mm, sec, ns int, w, startsAt time.Weekday) (int, int, int, int, int, int, int, time.Weekday) {
	if c.iso {
		startsAt = time.Monday
	}

	switch u {
	case Century:
		// n >= 0: +n 个世纪
		// n < 0:  本千年倒数第 n 个世纪
		if c.abs {
			y = n * 100
		} else {
			if n >= 0 {
				y = y - (y % 100) + n*100
			} else {
				y = y - (y % 1000) + (10+n)*100
			}
		}
	case Decade:
		// n > 0: 本世纪第 n 个年代
		// n < 0: 本世纪倒数第 n 个年代
		// n = 0: 保持在上级（如果有）或当前年代
		if c.abs {
			y = n * 10
		} else if n > 0 {
			y = y - (y % 100) + n*10
		} else if n < 0 {
			y = y - (y % 100) + (10+n)*10
		} else {
			y = y / 10 * 10
		}
	case Year:
		// n > 0: 本年代第 n 年
		// n < 0: 定位到本年代倒数第 n 年
		// n = 0: 保持在上级（如果有）或当前年
		if c.abs {
			y = n
		} else if n != 0 {
			if y -= y % 10; n > 0 {
				y += n
			} else {
				y += 10 + n
			}
		}
	case Quarter:
		if n > 0 {
			m = (n-1)*3 + 1
		} else if n < 0 {
			m = (5+n-1)*3 + 1
		} else {
			m -= (m - 1) % 3
		}
	case Month:
		if p == Quarter {
			// 季度内月份对齐：先找回季度起始月 (1, 4, 7, 10)
			if q := ((m-1)/3)*3 + 1; n > 0 {
				m = q + n - 1
			} else if n < 0 {
				m = q + 3 + n
			}
		} else {
			if n > 0 {
				m = n
			} else if n < 0 {
				m = 13 + n
			}
		}
	case Week:
		if d -= int(w-startsAt+7) % 7; n > 0 {
			d += (n - 1) * 7
		} else if n < 0 {
			d += n * 7
		}
	case Day:
		if n > 0 {
			d = n
		} else if n < 0 {
			d = DaysIn(y, m) + n + 1
		}
	case Hour:
		if n > 0 {
			h = n
		} else if n < 0 {
			h = 24 + n
		}
	case Minute:
		if n > 0 {
			mm = n
		} else if n < 0 {
			mm = 60 + n
		}
	case Second:
		if n > 0 {
			sec = n
		} else if n < 0 {
			sec = 60 + n
		}
	case Millisecond, Microsecond, Nanosecond:
		f := u.factor()
		if pf := f * 1000; n > 0 {
			ns = (ns/pf)*pf + n*f
		} else if n < 0 {
			ns = (ns/pf)*pf + pf + n*f
		}
	case YearWeek:
		// YearWeek: 遵循 “主权原则”，W01 是本年首个星期一。
		// ISOWeek:  严格遵循 ISO 8601，锚点为 1月4日，强制周一起始。
		if n == 0 { // 表示当前周
			// 将当前日期 d 向前推，对齐到最近的一个 startsAt（周起始日）。
			d -= int(w-startsAt+7) % 7
			break
		}

		// 正向定位 (n > 0)：从年初开始找。
		// 逆向定位 (n < 0)：从年尾开始找。
		// ISO 特殊处理：根据 ISO 8601 标准，1月4日必在第一周，12月28日必在最后一周。
		//
		// 为了计算第 n 周，代码需要先找到一个起始参考点（锚点）
		if m, d = 1, 1; n < 0 { // 默认正向从 1月1日 开始
			m, d = 12, 31 // 负向（倒数）从 12月31日 开始
		}

		if c.iso {
			if d = 4; n < 0 { // ISO 正向锚点是 1月4日
				d = 28 // ISO 负向锚点是 12月28日（保证在最后一周内）
			}
		}

		// 这是最核心的计算逻辑，通过 wAnchor（锚点当日是星期几）来对齐周。
		wAnchor := weekday(y, m, d)

		// 它先找到 1月1日 之后的第一个周起始日作为 W01 的开头，然后再增加 (n-1) 周。
		if !c.iso && n > 0 {
			d += (int(startsAt) - int(wAnchor) + 7) % 7 // 找到本年第一个 startsAt 当天
			d += (n - 1) * 7                            // 累加周数
		} else {
			if d -= int(wAnchor-startsAt+7) % 7; n > 0 {
				d += (n - 1) * 7 // ISO 第一周已经包含了锚点，所以 + (n-1)。
			} else {
				d += (n + 1) * 7 // 逆向计算，n = -1 时即为最后一周，不需额外偏移。
			}
		}
	case Weekday:
		// 周内第几天：
		// n > 0: 周内第 n 天（1 = startsAt）。
		// n < 0: 周内倒数第 n 天。
		// n = 0: 保持当前星期几不变（也就是今天）
		if n != 0 {
			if d -= int(w-startsAt+7) % 7; n > 0 {
				d += n - 1
			} else {
				d += n + 7
			}
		}
	}

	if u == Quarter || u == Month {
		y, m = addMonth(y, m, 0)
	}

	y, m, d, w = final(c, u, y, m, d)
	return y, m, d, h, mm, sec, ns, w
}

// applyRel 相对坐标对齐逻辑
func applyRel(c Flag, u, p Unit, n, y, m, d, h, mm, sec, ns int, w, startsAt time.Weekday) (int, int, int, int, int, int, int, time.Weekday) {
	if c.iso {
		startsAt = time.Monday
	}

	switch u {
	case Century:
		y -= y % 100
		y += n * 100
	case Decade:
		y -= y % 10
		y += n * 10
	case Year:
		y += n
	case Quarter:
		m -= (m - 1) % 3
		y, m = addMonth(y, m, n*3)
	case Month:
		y, m = addMonth(y, m, n)
	case Week:
		d -= int(w-startsAt+7) % 7
		d += n * 7
	case Day:
		d += n
	case Hour:
		h += n
	case Minute:
		mm += n
	case Second:
		sec += n
	case Millisecond, Microsecond, Nanosecond:
		ns += n * u.factor()
	case YearWeek:
		d -= int(w-startsAt+7) % 7
		d += n * 7
	case Weekday:
		if n != 0 {
			d -= int(w-startsAt+7) % 7 // 回到本周周初
			d += n                     // 偏移正负 n 天
		}
	}

	y, m, d, w = final(c, u, y, m, d)
	return y, m, d, h, mm, sec, ns, w
}

// applyOffset 相对坐标偏移逻辑
func applyOffset(c Flag, u, p Unit, n, y, m, d, h, mm, sec, ns int, w time.Weekday) (int, int, int, int, int, int, int, time.Weekday) {
	switch u {
	case Century:
		y += n * 100
	case Decade:
		y += n * 10
	case Year:
		y += n
	case Quarter:
		y, m = addMonth(y, m, n*3)
	case Month:
		y, m = addMonth(y, m, n)
	case Week, YearWeek:
		// 纯偏移模式下，所有周相关的单位都等同于增加 n * 7 天。
		d += n * 7
	case Day, Weekday:
		d += n
	case Hour:
		h += n
	case Minute:
		mm += n
	case Second:
		sec += n
	case Millisecond, Microsecond, Nanosecond:
		ns += n * u.factor()
	}

	y, m, d, w = final(c, u, y, m, d)
	return y, m, d, h, mm, sec, ns, w
}

func final(c Flag, u Unit, y, m, d int) (int, int, int, time.Weekday) {
	if !c.overflow && u <= Month {
		// 仅针对这些时间单元做天数溢出处理
		if dd := DaysIn(y, m); d > dd {
			d = dd
		}
	}

	if c.fill {
		switch u {
		case Century:
			y += 99
		case Decade:
			y += 9
		case Quarter:
			m += 2
		case Week, YearWeek:
			d += 6
		default:
		}
	}

	return y, m, d, weekday(y, m, d)
}

// align 执行最终的时间分量对齐（归零或置满）
func align(c Flag, last Unit, y, m, d, h, mm, sec, ns int) (int, int, int, int, int, int, int) {
	if !c.fill {
		switch last {
		case Century, Decade, Year:
			m, d, h, mm, sec, ns = 1, 1, 0, 0, 0, 0
		case Quarter, Month:
			d, h, mm, sec, ns = 1, 0, 0, 0, 0
		case YearWeek, Week, Weekday, Day:
			h, mm, sec, ns = 0, 0, 0, 0
		case Hour:
			mm, sec, ns = 0, 0, 0
		case Minute:
			sec, ns = 0, 0
		case Second:
			ns = 0
		case Millisecond, Microsecond, Nanosecond:
			f := last.factor()
			ns = (ns / f) * f
		}
	} else {
		switch last {
		case Century, Decade, Year:
			m, d, h, mm, sec, ns = 12, 31, 23, 59, 59, 999999999
		case Quarter, Month:
			d, h, mm, sec, ns = DaysIn(y, m), 23, 59, 59, 999999999
		case YearWeek, Week, Weekday, Day:
			h, mm, sec, ns = 23, 59, 59, 999999999
		case Hour:
			mm, sec, ns = 59, 59, 999999999
		case Minute:
			sec, ns = 59, 999999999
		case Second:
			ns = 999999999
		case Millisecond, Microsecond, Nanosecond:
			f := last.factor()
			ns = (ns/f)*f + (f - 1)
		}
	}

	return y, m, d, h, mm, sec, ns
}
