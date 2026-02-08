package aeon

import (
    "time"
)

type Unit int

const (
    Century Unit = iota
    Decade
    Year
    Month
    Day
    Hour
    Minute
    Second
    Millisecond
    Microsecond
    Nanosecond
    Quarter // 季度流
    Week    // 月周流
    Weekday // 星期流
)

var (
    years    = []Unit{Century, Decade, Year, Month, Day, Hour, Minute, Second, Millisecond, Microsecond, Nanosecond} // 年序列
    quarters = []Unit{Quarter, Month, Day, Hour, Minute, Second, Millisecond, Microsecond, Nanosecond}               // 季度流
    weeks    = []Unit{Week, Weekday, Hour, Minute, Second, Millisecond, Microsecond, Nanosecond}                     // 月周流
    weekdays = []Unit{Weekday, Hour, Minute, Second, Millisecond, Microsecond, Nanosecond}                           // 星期流
)

func (u Unit) seq() []Unit {
    switch u {
    case Quarter:
        return quarters
    case Week:
        return weeks
    case Weekday:
        return weekdays
    default:
        return years[u:]
    }
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

func apply(f Action, c Flag, first bool, u, p Unit, n, pN int, y, m, d, h, mm, s, ns int, w, sw time.Weekday) (int, int, int, int, int, int, int, time.Weekday) {
    switch f {
    case seAbs, goAbs: // 全绝对
        return applyAbs(c, u, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
    case seRel, goRel: // 全相对
        return applyRel(c, u, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
    case seAt, goAt: // 首绝后相
        if first {
            return applyAbs(c, u, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
        }
        return applyRel(c, u, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
    default: // seIn, goIn
        if first {
            return applyRel(c, u, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
        }
        return applyAbs(c, u, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
    }
}

// applyAbs 是级联引擎的 “空间定位器”，负责执行基于父容器的绝对定位逻辑。
//
// 核心逻辑：
//   - 索引机制: 根据 n 值实现容器内的绝对定位（正向、倒数或保持当前）。
//   - 模式切换: 根据 c.goMode 在 “对齐边界（归零或置满）” 与 “高保真跳转（保留子级精度）” 间切换。
//   - 溢出保护: 在跳转模式下，自动校正由于单位变更导致的日期溢出（如月末对齐）。
func applyAbs(c Flag, u, p Unit, n, pN, y, m, d, h, mm, sec, ns int, w, sw time.Weekday) (int, int, int, int, int, int, int, time.Weekday) {
    if c.isoWeek {
        sw = time.Monday
    }

    switch u {
    case Century:
        // n >= 0: +n 个世纪
        // n < 0:  本千年倒数第 n 个世纪
        if c.abs {
            y = n * 100
            break
        }

        if n < 0 { // 预平移：为千年之内的倒数寻址留出空间
            y += 1000
        }

        if c.goMode { // Go 模式：去指定 “世纪”，保留 “年代” 和 “年位”。
            y = (y - y%1000) + (y % 100) + n*100
        } else { // se 模式：抹平世纪以下的位
            if n == 0 {
                y = y - y%100
            } else {
                y = (y - y%1000) + n*100
            }
        }
    case Decade:
        if n < 0 {
            y += 100
        }

        // 假设日期：2021-02-02
        if c.goMode { // Go 模式：去指定 “年代”，保留 “年位”。
            // 例如 in(0)=2001, 1=2011, 2=2021, -1=2091
            y = (y - y%100) + (y % 10) + n*10
        } else { // Start/End 模式
            // n > 0: 本世纪第 n 个年代
            // n < 0: 本世纪倒数第 n 个年代
            // n = 0: 保持当前年代
            // 示例：0=2020/2029, 1=2010/2019, 2=2020/2029, -1=2029/2099
            if n == 0 {
                y = y - y%10
            } else {
                y = (y - y%100) + n*10
            }
        }
    case Year:
        // n > 0: 本年代第 n 年
        // n < 0: 定位到本年代倒数第 n 年
        // n = 0: 保持在上级（如果有）或当前年
        if c.abs {
            y = n
        } else if c.goMode || n != 0 {
            // 偏移补偿算法：如果 n < 0，补偿 10 年。
            y = (y - y%10) + (10 & (n >> 63)) + n
        }
    case Quarter:
        if n > 0 {
            m = (n-1)*3 + 1
        } else if n < 0 {
            m = (5+n-1)*3 + 1
        } else if !c.goMode {
            m -= (m - 1) % 3
        }
    case Month:
        if p == Quarter { // 季内月
            if n != 0 {
                if m -= (m - 1) % 3; n > 0 {
                    m += n - 1
                } else {
                    m += 3 + n
                }
            }
        } else {
            if n > 0 {
                m = n
            } else if n < 0 {
                m = 13 + n
            }
        }
    case Week:
        if c.goMode && n == 0 {
            break
        }

        if c.qtrWeek { // 季内周
            if n == 0 { // se 模式：回到周首
                d -= int(w-sw+7) % 7
                break
            }

            if m -= (m - 1) % 3; n > 0 { // 回到季首
                d = 1 + (n-1)*7 // 正数从季初开始数周
            } else {
                m += 2                         // 定位到月末
                d = DaysIn(y, m) + (n+1)*7 - 6 // 从季末开始数周
                if !c.ordWeek {                // 如果是日历周 (自然周)
                    // 从序数锚点修正到自然周首 (sw)
                    d += (int(sw-weekday(y, m, d+6)-1) + 7) % 7
                }
            }

            if !c.ordWeek && c.goMode {
                // 统一对齐公式：不论正反，d 此时均在寻址周的起始日。
                d += (int(w-weekday(y, m, d)) + 7) % 7
            }
        } else if c.fullWeek {
            // 完整周：从本月第 1 个周一开始
            if n > 0 {
                // 正向：找到第一个 sw（周一）再加偏移
                d = 1 + (int(sw)-int(weekday(y, m, 1))+7)%7 + (n-1)*7
            } else if n < 0 {
                // 逆向：从月末找到最后一个 sw 再减偏移
                last := DaysIn(y, m)
                d = last - (int(weekday(y, m, last))-int(sw)+7)%7 + (n+1)*7
            } else {
                d -= int(w-sw+7) % 7 // n=0 归位到当前周首
            }

            if c.goMode {
                d += int(w-sw+7) % 7
            }
        } else if c.ordWeek {
            // 序数周：从本月 1 日开始数周
            if n == 0 {
                d -= int(w-sw+7) % 7 // se 模式：回到周首
            } else if d = 1 + (n-1)*7; n < 0 {
                d = DaysIn(y, m) + (n+1)*7 - 6
            }
        } else if c.isoWeek {
            // ISO 年周：遵循 ISO 8601
            if n == 0 { // 回到周首
                d -= int(w-sw+7) % 7
                break
            }

            // 选择锚点：默认正向从1月1日开始，负向从12月31日开始
            // ISO特殊处理：正向用1月4日，逆向用12月28日（保证都在目标周内）
            if m, d = 1, 1; n < 0 {
                m, d = 12, 31
            }

            if d = 4; n < 0 { // ISO 正向锚点是1月4日
                d = 28 // ISO 负向锚点是12月28日（保证在最后一周内）
            }

            // 从锚点对齐到周首，再加上周偏移。
            wAnchor := weekday(y, m, d)
            if d -= int(wAnchor-sw+7) % 7; n > 0 {
                d += (n - 1) * 7
            } else {
                d += (n + 1) * 7
            }

            if c.goMode { // Go 模式：恢复周内偏移，保持当前星期几。
                d += int(w-sw+7) % 7
            }
        } else { // 日历周 (默认)
            if n > 0 {
                d = 1 + (n-1)*7
            } else if n < 0 {
                d = DaysIn(y, m) + (n+1)*7
            }

            tw := weekday(y, m, d) // 归位到周首
            if d -= int(tw-sw+7) % 7; c.goMode {
                d += int(w-sw+7) % 7
            }
        }
    case Weekday:
        // 周内第几天：
        // n > 0: 周内第 n 天
        // n < 0: 周内倒数第 n 天
        // n = 0: 保持当前星期几不变
        if n == 0 {
            break
        }
        // 识别是否处于 Week 的序数周级联中
        if c.ordWeek && p == Week && pN != 0 {
            // 🦬 预平移坐标参考系：
            // 如果 n < 0，将参考点向左拨 1 天，让 0-6 索引与序数意图完美重合。
            if cur := int(weekday(y, m, d+n>>63)); pN >= 0 {
                d += (n - cur + 7) % 7
            } else {
                d -= (cur - n + 7) % 7
            }
        } else {
            if d -= int(w-sw+7) % 7; n > 0 {
                d += n - 1
            } else {
                d += n + 7
            }
        }
    case Day:
        if n > 0 {
            d = n
        } else if n < 0 {
            d = DaysIn(y, m) + n + 1
        }
    case Hour:
        if n >= 0 {
            h = n
        } else {
            h = 24 + n
        }
    case Minute:
        if n >= 0 {
            mm = n
        } else {
            mm = 60 + n
        }
    case Second:
        if n >= 0 {
            sec = n
        } else {
            sec = 60 + n
        }
    case Millisecond, Microsecond, Nanosecond:
        f := u.factor()
        if pf := f * 1000; n >= 0 {
            ns = (ns/pf)*pf + n*f
        } else {
            ns = (ns/pf)*pf + pf + n*f
        }
    }

    if u == Quarter || u == Month {
        y, m = addMonth(y, m, 0)
    }

    y, m, d, w = final(c, u, n, y, m, d)
    return y, m, d, h, mm, sec, ns, w
}

// applyRel 相对坐标对齐逻辑
func applyRel(c Flag, u, p Unit, n, pN, y, m, d, h, mm, sec, ns int, w, sw time.Weekday) (int, int, int, int, int, int, int, time.Weekday) {
    if c.isoWeek {
        sw = time.Monday
    }

    switch u {
    case Century:
        if c.abs {
            y = n
            break
        }
        if c.goMode {
            y += n * 100
            break
        }
        y += (y - y%100) + n*100
    case Decade:
        if c.goMode {
            y += n * 10
            break
        }
        y += (y - y%10) + n*10
    case Year:
        if c.abs {
            y = n
            break
        }
        y += n
    case Quarter:
        if !c.goMode { // se 模式：回到季首
            m -= (m - 1) % 3
        }
        y, m = addMonth(y, m, n*3)
    case Month:
        y, m = addMonth(y, m, n)
    case Week:
        if !c.goMode { // se 模式：回到周初
            d -= int(w-sw+7) % 7
        }
        d += n * 7
    case Day, Weekday: // 在偏移场景下，都是偏移 ±n 天
        if u == Weekday && !c.goMode && n != 0 {
            d -= int(w-sw+7) % 7 // se 模式：回到周初
        }
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

    y, m, d, w = final(c, u, n, y, m, d)
    return y, m, d, h, mm, sec, ns, w
}

func final(c Flag, u Unit, n, y, m, d int) (int, int, int, time.Weekday) {
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
        case Week:
            d += 6
        default:
        }
    } else {
        // Go 模式下，序数周会定位到从月末倒数周的周初，
        // 所以需要加上 6 天，来到本月最后一天。
        if u == Week && c.ordWeek && n < 0 && c.goMode {
            d += 6
        }
    }

    return y, m, d, weekday(y, m, d)
}

// align 执行最终的时间分量对齐（归零或置满）
func align(c Flag, u Unit, y, m, d, h, mm, sec, ns int) (int, int, int, int, int, int, int) {
    if !c.fill {
        switch u {
        case Century, Decade, Year:
            m, d, h, mm, sec, ns = 1, 1, 0, 0, 0, 0
        case Quarter, Month:
            d, h, mm, sec, ns = 1, 0, 0, 0, 0
        case Week, Weekday, Day:
            h, mm, sec, ns = 0, 0, 0, 0
        case Hour:
            mm, sec, ns = 0, 0, 0
        case Minute:
            sec, ns = 0, 0
        case Second:
            ns = 0
        case Millisecond, Microsecond, Nanosecond:
            f := u.factor()
            ns = (ns / f) * f
        }
    } else {
        switch u {
        case Century, Decade, Year:
            m, d, h, mm, sec, ns = 12, 31, 23, 59, 59, 999999999
        case Quarter, Month:
            d, h, mm, sec, ns = DaysIn(y, m), 23, 59, 59, 999999999
        case Week, Weekday, Day:
            h, mm, sec, ns = 23, 59, 59, 999999999
        case Hour:
            mm, sec, ns = 59, 59, 999999999
        case Minute:
            sec, ns = 59, 999999999
        case Second:
            ns = 999999999
        case Millisecond, Microsecond, Nanosecond:
            f := u.factor()
            ns = (ns/f)*f + (f - 1)
        }
    }

    return y, m, d, h, mm, sec, ns
}
