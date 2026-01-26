package aeon

import (
    "time"
    "unsafe"
)

// tryFastParse 进化版：大类分发 + 全能流式提取
func tryFastParse(s string, l *time.Location) (Time, bool, error) {
    ln := len(s)
    if ln < 4 {
        return Time{}, false, nil
    }

    c0, c1, c2, c3 := s[0], s[1], s[2], s[3]
    if isDigit(c0) && isDigit(c1) && isDigit(c2) && isDigit(c3) {
        y := int(c0-'0')*1000 + int(c1-'0')*100 + int(c2-'0')*10 + int(c3-'0')
        if ln == 4 {
            return Time{time: time.Date(y, 1, 1, 0, 0, 0, 0, l), weekStarts: DefaultWeekStarts}, true, nil
        }

        c4 := s[4]
        if isDigit(c4) {
            if ln == 8 {
                t, err := fastParseD8(s, l)
                return Time{time: t, weekStarts: DefaultWeekStarts}, true, err
            }
            if ln == 14 {
                t, err := fastParseDT14(s, l)
                return Time{time: t, weekStarts: DefaultWeekStarts}, true, err
            }
        } else {
            // 一级旁路：标准格式 SWAR
            switch ln {
            case 10:
                if !isDigit(s[7]) {
                    t, err := fastParseD10(s, l)
                    return Time{time: t, weekStarts: DefaultWeekStarts}, true, err
                }
            case 19:
                if !isDigit(s[7]) && (s[10] == ' ' || s[10] == 'T') {
                    t, err := fastParseDT19(s, l)
                    return Time{time: t, weekStarts: DefaultWeekStarts}, true, err
                }
            case 20:
                if s[19] == 'Z' {
                    t, err := fastParseZulu20(s, time.UTC)
                    return Time{time: t, weekStarts: DefaultWeekStarts}, true, err
                }
            case 25:
                if s[19] == '+' || s[19] == '-' {
                    t, err := fastParseRFC3339(s, l)
                    return Time{time: t, weekStarts: DefaultWeekStarts}, true, err
                }
            }

            // 二级旁路：全能 Seeker (支持变体分隔符)
            t, err := fastParseUniversalYearStarted(s, y, l)
            if err == nil {
                return Time{time: t, weekStarts: DefaultWeekStarts}, true, nil
            }
        }
    }

    return Time{}, false, nil
}

func fastParseUniversalYearStarted(s string, y int, l *time.Location) (time.Time, error) {
    ln := len(s)
    var v [5]int
    vIdx := 0
    cur := 0
    foundDigit := false
    i := 4

    var nsec int
    var zone *time.Location = l

    for ; i < ln && vIdx < 5; i++ {
        c := s[i]
        if isDigit(c) {
            cur = cur*10 + int(c-'0')
            foundDigit = true
        } else {
            if foundDigit {
                v[vIdx] = cur
                vIdx++
                cur = 0
                foundDigit = false
                // 区分纳秒与分隔符
                if c == '.' && (vIdx == 2 || vIdx == 5) {
                    goto ExtractNano
                }
                // 区分时区与日期分隔符：只有在解析完 D 之后，且符号后紧跟数字时，才判定为时区
                if (c == '+' || c == '-') && vIdx >= 2 && i+1 < ln && isDigit(s[i+1]) {
                    goto ProcessZone
                }
                if c == 'Z' {
                    goto ProcessZone
                }
            }
        }
    }
    if foundDigit && vIdx < 5 {
        v[vIdx] = cur
        vIdx++
    }

ExtractNano:
    if i < ln && s[i] == '.' {
        i++
        start := i
        val := 0
        for i < ln && isDigit(s[i]) {
            val = val*10 + int(s[i]-'0')
            i++
        }
        dLen := i - start
        if dLen > 0 {
            nsec = val
            for k := dLen; k < 9; k++ {
                nsec *= 10
            }
            for k := dLen; k > 9; k-- {
                nsec /= 10
            }
        }
    }

ProcessZone:
    for ; i < ln; i++ {
        c := s[i]
        if c == 'Z' {
            zone = time.UTC
            break
        }
        if c == '+' || c == '-' {
            sign := 1
            if c == '-' {
                sign = -1
            }
            i++
            oh, om := 0, 0
            cnt := 0
            for i < ln && isDigit(s[i]) && cnt < 2 {
                oh = oh*10 + int(s[i]-'0')
                i++
                cnt++
            }
            if i < ln && s[i] == ':' {
                i++
            }
            cnt = 0
            for i < ln && isDigit(s[i]) && cnt < 2 {
                om = om*10 + int(s[i]-'0')
                i++
                cnt++
            }
            zone = time.FixedZone("", sign*(oh*3600+om*60))
            break
        }
    }

    return time.Date(y, time.Month(max(1, v[0])), max(1, v[1]), v[2], v[3], v[4], nsec, zone), nil
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func parse4(s string) int {
    v := *(*uint32)(unsafe.Pointer(unsafe.StringData(s)))
    return int((v&0xF)*1000 + ((v>>8)&0xF)*100 + ((v>>16)&0xF)*10 + (v>>24)&0xF)
}

func fastParseD8(s string, l *time.Location) (time.Time, error) {
    return time.Date(parse4(s[0:4]), time.Month(int((s[4]-'0')*10+s[5]-'0')), int((s[6]-'0')*10+s[7]-'0'), 0, 0, 0, 0, l), nil
}

func fastParseDT14(s string, l *time.Location) (time.Time, error) {
    return time.Date(parse4(s[0:4]), time.Month(int((s[4]-'0')*10+s[5]-'0')), int((s[6]-'0')*10+s[7]-'0'),
        int((s[8]-'0')*10+s[9]-'0'), int((s[10]-'0')*10+s[11]-'0'), int((s[12]-'0')*10+s[13]-'0'), 0, l), nil
}

func fastParseDT19(s string, l *time.Location) (time.Time, error) {
    return time.Date(parse4(s[0:4]), time.Month(int((s[5]-'0')*10+s[6]-'0')), int((s[8]-'0')*10+s[9]-'0'),
        int((s[11]-'0')*10+s[12]-'0'), int((s[14]-'0')*10+s[15]-'0'), int((s[17]-'0')*10+s[18]-'0'), 0, l), nil
}

func fastParseD10(s string, l *time.Location) (time.Time, error) {
    return time.Date(parse4(s[0:4]), time.Month(int((s[5]-'0')*10+s[6]-'0')), int((s[8]-'0')*10+s[9]-'0'), 0, 0, 0, 0, l), nil
}

func fastParseZulu20(s string, l *time.Location) (time.Time, error) {
    return time.Date(parse4(s[0:4]), time.Month(int((s[5]-'0')*10+s[6]-'0')), int((s[8]-'0')*10+s[9]-'0'),
        int((s[11]-'0')*10+s[12]-'0'), int((s[14]-'0')*10+s[15]-'0'), int((s[17]-'0')*10+s[18]-'0'), 0, l), nil
}

func fastParseRFC3339(s string, l *time.Location) (time.Time, error) {
    y := parse4(s[0:4])
    m := int((s[5]-'0')*10 + s[6] - '0')
    d := int((s[8]-'0')*10 + s[9] - '0')
    h := int((s[11]-'0')*10 + s[12] - '0')
    min := int((s[14]-'0')*10 + s[15] - '0')
    sec := int((s[17]-'0')*10 + s[18] - '0')
    sign := 1
    if s[19] == '-' {
        sign = -1
    }
    oh := int((s[20]-'0')*10 + s[21] - '0')
    om := int((s[23]-'0')*10 + s[24] - '0')
    return time.Date(y, time.Month(m), d, h, min, sec, 0, time.FixedZone("", sign*(oh*3600+om*60))), nil
}
