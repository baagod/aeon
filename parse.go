package aeon

import (
    "fmt"
    "sync"
    "time"
    "unsafe"
)

// bufPool 调整为 64 字节，匹配 CPU Cache Line 大小，且足够容纳所有标准时间布局。
var bufPool = sync.Pool{New: func() any {
    b := make([]byte, 64)
    return &b
}}

const (
    ANSIC      = "Mon Jan _2 15:04:05 2006"
    UnixD      = "Mon Jan _2 15:04:05 MST 2006"
    RubyD      = "Mon Jan 02 15:04:05 -0700 2006"
    RFC822     = "02 Jan 06 15:04 MST"
    RFC822Z    = "02 Jan 06 15:04 -0700"
    RFC850     = "Monday, 02-Jan-06 15:04:05 MST"
    RFC1123    = "Mon, 02 Jan 2006 15:04:05 MST"
    RFC1123Z   = "Mon, 02 Jan 2006 15:04:05 -0700"
    RFC3339    = "2006-01-02T15:04:05Z07:00"
    RFC3339Ns  = "2006-01-02T15:04:05.999999999Z07:00"
    Kitchen    = "3:04PM"
    Stamp      = "Jan _2 15:04:05"
    StampMilli = "Jan _2 15:04:05.000"
    StampMicro = "Jan _2 15:04:05.000000"
    StampNano  = "Jan _2 15:04:05.000000000"
    StampNs    = "Jan _2 15:04:05.999999999"

    // 核心布局 (DT/D 系列)
    DT       = "2006-01-02 15:04:05"
    DTMilli  = "2006-01-02 15:04:05.000"
    DTMicro  = "2006-01-02 15:04:05.000000"
    DTNano   = "2006-01-02 15:04:05.000000000"
    DTNs     = "2006-01-02 15:04:05.999999999"
    DateOnly = "2006-01-02"
    DMilli   = "2006-01-02.000"
    DMicro   = "2006-01-02.000000"
    DNano    = "2006-01-02.000000000"
    TimeOnly = "15:04:05"

    // 补强布局 (归一化复用版)
    DTFull         = "2006-01-02 15:04:05.999999999 -0700 MST"
    DTTZ           = "2006-01-02 15:04:05-07:00"
    DTTZShort      = "2006-01-02 15:04:05-07"
    DTISO          = "2006-01-02T15:04:05-07:00"
    DTPMMST        = "2006-01-02 15:04:05 PM MST"
    DTPMShortMST   = "2006-01-02 15:04:05PM MST"
    RFC3339Space   = "2006-01-02 15:04:05Z07:00"
    RFC3339NsSpace = "2006-01-02 15:04:05.999999999Z07:00"

    // 紧凑与特殊格式
    DTCompact     = "20060102150405"
    DCompact      = "20060102"
    TimeCompact   = "150405"
    DTVeryShort   = "2006-1-2 15:4:5"
    DTShort       = "2006-1-2 15:4"
    DOnlyShort    = "2006-1-2"
    TimeVeryShort = "15:4:5"
    TimeShort     = "15:4"
    MonthD        = "1-2"
    YearOnly      = "2006"

    DTCompactTZ     = "20060102150405-07:00"
    DTCompactZ      = "20060102150405Z07:00"
    DTCompactMilli  = "20060102150405.000"
    DHourShort      = "2006-1-2 15"
    HourOnly        = "15"
    DMonth          = "2006-1"
    FormattedD      = "Jan 2, 2006"
    FormattedDayD   = "Mon, Jan 2, 2006"
    DayDateTime     = "Mon, Jan 2, 2006 3:04 PM"
    Cookie          = "Monday, 02-Jan-2006 15:04:05 MST"
    Http            = "Mon, 02 Jan 2006 15:04:05 GMT"
    RFC1036         = "Mon, 02 Jan 06 15:04:05 -0700"
    RFC7231         = "Mon, 02 Jan 2006 15:04:05 MST"
    TimeTZShort     = "15:04:05-07"
    DTFullVeryShort = "2006-1-2 15:4:5 -0700 MST"
    DTNsVeryShort   = "2006-1-2 15:4:5.999999999"

    // ISO8601 家族
    ISO8601       = "2006-01-02T15:04:05-07:00"
    ISO8601Ns     = "2006-01-02T15:04:05.999999999-07:00"
    ISO8601Zulu   = "2006-01-02T15:04:05Z"
    ISO8601ZuluNs = "2006-01-02T15:04:05.999999999Z"
)

// --- 基因 ID 定义 (4-bit, 0-15) ---
const (
    _         uint64 = iota
    idDigitS         // 1: 通用数字
    idYear           // 2: 4位数字
    idMonth          // 3: 变长月份
    idWeek           // 4: 变长星期
    idISO            // 5: T 或 Z 协议位
    idSep            // 6: 万能分隔符
    idDigit8         // 7: 8位纯数字
    idDigit14        // 8: 14位纯数字
    idMonthL         // 9: 全名月份
    idWeekL          // 10: 全名星期
)

const (
    dnaHeaderShift    = 56
    dnaSlotStartShift = 52
    dnaSlotBits       = 4
    dnaMaxSlots       = 14
)

type segment struct {
    id    uint64
    start int
    end   int
}

// layoutInfo 存储母版的完整信息
type layoutInfo struct {
    layout string
    segs   [32]segment
    count  int
}

var (
    dnaMap  = map[uint64]*layoutInfo{}
    formats = []string{
        // 核心母版 (YMD HMS Nano)
        "2006",
        "2006-",
        "2006-1",
        "2006-1-",
        "2006-1-2",
        "2006-1-2-",
        "2006-1-2 15:4",
        "2006-1-2 15:4:5",
        "2006-1-2 15:4:5.999999999",
        "2006-1-2T15:4:5",
        "2006-1-2T15:4:5Z",
        "2006-1-2T15:4:5-07:00",
        "2006-1-2T15:4:5.999999999Z07:00",
        "2006-1-2 15:4:5.999999999 -0700 MST",

        // 时间母版
        "3:4PM",
        "15:4",
        "15:4:",
        "15:4:5",
        "15:4:5:",
        "15:4:5.999999999",

        // 人文母版
        "2 Jan 2006",
        "2 Jan 6 15:4",
        "2 Jan 2006 15:4:5",
        "2 Jan 6 15:4 -0700",

        "Jan 2 15:4:5",
        "Jan 2 2006",
        "Jan 2, 2006",
        "Jan 2 15:4:5",
        "Jan 2 2006 15:4:5",
        "Jan 2 15:4:5.999999999",
        "Jan 2 2006 15:4:5.999999999",

        "January 2 2006",
        "January 2, 2006",
        "January 2 2006 15:4:5",
        "January 2 2006 15:4:5.999999999",

        "Mon, 2 Jan 2006 15:4:5 MST",
        "Mon, 2 Jan 2006 15:4:5 -0700",

        "Mon Jan 2 2006",
        "Mon Jan 2 15:4:5 2006",
        "Mon Jan 2 15:4:5 MST 2006",
        "Mon Jan 2 15:4:5 -0700 2006",

        "Monday January 2 2006",
        "Monday January 2 15:4:5 2006",
        "Monday, 2 January 2006 15:04:05 -0700",

        // 紧凑格式
        DTCompact, DCompact,

        // 纯纳秒
        ".999999999",
    }
)

// RegisterDNA 将布局字符串注册到 DNA 地图中
func RegisterDNA(layout string) {
    var segs [32]segment
    fp, count := getDNA(layout, &segs)
    if fp != 0 {
        dnaMap[fp] = &layoutInfo{layout: layout, segs: segs, count: count}
    }
}

func init() {
    for _, layout := range formats {
        RegisterDNA(layout)
    }
}

// ParseE 解析时间字符串，返回 Time 和 error
func ParseE(s string, loc ...*time.Location) (Time, error) {
    start, end := 0, len(s)
    for start < end {
        if c := s[start]; c > ' ' && c != '"' {
            break
        }
        start++
    }

    for end > start {
        if c := s[end-1]; c > ' ' && c != '"' {
            break
        }
        end--
    }

    if s = s[start:end]; s == "" || s == "null" {
        return Time{}, nil
    }

    l := DefaultTimeZone
    if len(loc) > 0 && loc[0] != nil {
        l = loc[0]
    }

    // --- 1. 尝试 L1 决策树快轨 (parse_fast.go) ---
    if t, ok, err := tryFastParse(s, l); ok {
        return t, err
    }

    // --- 2. 原有 DNA 架构开始 (兜底) ---
    var inputSegs [32]segment
    fp, count := getDNA(s, &inputSegs)

    master, ok := dnaMap[fp]
    if !ok {
        return Time{}, fmt.Errorf("aeon 无法识别该时间指纹: [%x]", fp)
    }

    // 从池中获取缓冲区
    pBuf := bufPool.Get().(*[]byte)
    defer bufPool.Put(pBuf)
    buf := *pBuf

    // 动态合成
    dynamicLayout, identical, pos := morphLayout(s, inputSegs[:count], master, buf)

    // 如果不一致，构造临时 string。由于 time.Parse 不持有它，这是安全的。
    if !identical {
        dynamicLayout = unsafe.String(&buf[0], pos)
    }

    t, err := time.ParseInLocation(dynamicLayout, s, l)
    if err != nil {
        return Time{}, err
    }

    return Time{time: t, weekStarts: DefaultWeekStarts}, nil
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

// getDNA 核心测序仪：提取 DNA 并记忆物理位置
func getDNA(s string, segs *[32]segment) (uint64, int) {
    ln := len(s)
    if ln == 0 || ln > 64 {
        return 0, 0
    }

    var dna uint64
    count := 0

    for i := 0; i < ln && count < 32; {
        start, c := i, s[i]
        var id uint64
        consumed := 1

        switch {
        case isDigit(c):
            for i+consumed < ln && isDigit(s[i+consumed]) {
                consumed++
            }
            id = idDigitS
            if consumed == 4 {
                id = idYear
            } else if count == 0 {
                if consumed == 8 {
                    id = idDigit8
                } else if consumed == 14 {
                    id = idDigit14
                }
            }
        case isAlpha(c):
            if (c == 'T' || c == 'Z') && (i+1 == ln || !isAlpha(s[i+1])) {
                id = idISO
            } else {
                for i+consumed < ln && isAlpha(s[i+consumed]) {
                    consumed++
                }
                if consumed >= 3 {
                    sig := uint32(s[i]|0x20)<<16 | uint32(s[i+1]|0x20)<<8 | uint32(s[i+2]|0x20)
                    switch sig {
                    case 0x6a616e, 0x666562, 0x6d6172, 0x617072, 0x6d6179, 0x6a756e, 0x6a756c, 0x617567, 0x736570, 0x6f6374, 0x6e6f76, 0x646563:
                        id = idMonth
                        if consumed > 3 {
                            id = idMonthL
                        }
                    case 0x6d6f6e, 0x747565, 0x776564, 0x746875, 0x667269, 0x736174, 0x73756e:
                        id = idWeek
                        if consumed > 3 {
                            id = idWeekL
                        }
                    }
                }
                if id == 0 {
                    id = idSep
                }
            }
        default:
            for i+consumed < ln {
                nc := s[i+consumed]
                if isDigit(nc) || isAlpha(nc) {
                    break
                }
                consumed++
            }
            id = idSep
        }

        // 记录指纹 (DNA)
        if count < dnaMaxSlots {
            dna |= id << (dnaSlotStartShift - uint(count)*dnaSlotBits)
        } else {
            // 指纹槽位溢出：用当前 ID 覆盖最后一个槽位
            dna = (dna & ^uint64(0xF)) | id
        }

        segs[count] = segment{id: id, start: start, end: start + consumed}
        count++
        i += consumed
    }

    return dna | (uint64(count) << 56), count
}

// morphLayout 布局重组引擎：现场克隆母版
func morphLayout(input string, inputSegs []segment, master *layoutInfo, buf []byte) (string, bool, int) {
    pos := 0
    identical := true
    layout := master.layout

    // 此时 inputSegs 和 master.segs 的长度和 ID 序列是完全对齐的
    for i := 0; i < len(inputSegs); i++ {
        inSeg := inputSegs[i]
        msSeg := master.segs[i]

        // 检查是否为时区符号（如 -07:00 中的 -）。
        // 如果是 idSep 且紧接着 master 中的 07，则强制使用 master 的符号。
        var isTZSign bool
        if inSeg.id == idSep {
            if msSeg.end+2 <= len(layout) && layout[msSeg.end:msSeg.end+2] == "07" {
                isTZSign = true
            }
        }

        var n int
        if !isTZSign && (inSeg.id == idSep || inSeg.id == idISO) {
            // 物理/协议组件：从输入中拷贝用户写的真实字符
            n = copy(buf[pos:], input[inSeg.start:inSeg.end])
        } else {
            // 时间意图组件：从母版布局中拷贝标准 Token
            n = copy(buf[pos:], layout[msSeg.start:msSeg.end])
        }

        // 恒等检查：如果当前 Token 与母版对应位置不一致，则标记非恒等
        if identical {
            if pos+n > len(layout) || layout[pos:pos+n] != string(buf[pos:pos+n]) {
                identical = false
            }
        }

        pos += n
    }

    // 如果完全一致，直接返回常量母版字符串 (0 堆分配)
    if identical && pos == len(layout) {
        return layout, true, pos
    }

    return "", false, pos
}

// isAlpha 快速判定字母 (包含大小写，1次位运算 + 2次比较)
func isAlpha(c byte) bool {
    return (c|0x20) >= 'a' && (c|0x20) <= 'z'
}

// isDigit 快速判定数字 (减法溢出法，1次比较)
func isDigit(c byte) bool {
    return uint8(c-'0') <= 9
}
