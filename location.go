package aeon

import (
    "sync"
    "time"
)

// 时区名称常量定义 (IANA 标准)
const (
    Local = "Local" // 本地时间
    UTC   = "UTC"   // 协调世界时间

    CET = "CET" // 中欧标准时间
    EET = "EET" // 东欧标准时间
    EST = "EST" // 东部标准时间
    GMT = "GMT" // 格林尼治标准时间
    MET = "MET" // 中欧时间
    MST = "MST" // 山地标准时间
    WET = "WET" // 西欧标准时间

    Cuba      = "Cuba"      // 古巴
    Egypt     = "Egypt"     // 埃及
    Eire      = "Eire"      // 爱尔兰
    Greenwich = "Greenwich" // 格林尼治
    Iceland   = "Iceland"   // 冰岛
    Iran      = "Iran"      // 伊朗
    Israel    = "Israel"    // 以色列
    Jamaica   = "Jamaica"   // 牙买加
    Japan     = "Japan"     // 日本
    Libya     = "Libya"     // 利比亚
    Poland    = "Poland"    // 波兰
    Portugal  = "Portugal"  // 葡萄牙
    PRC       = "PRC"       // 中国
    Singapore = "Singapore" // 新加坡
    Turkey    = "Turkey"    // 土耳其

    Shanghai   = "Asia/Shanghai"       // 上海
    Chongqing  = "Asia/Chongqing"      // 重庆
    Harbin     = "Asia/Harbin"         // 哈尔滨
    Urumqi     = "Asia/Urumqi"         // 乌鲁木齐
    HongKong   = "Asia/Hong_Kong"      // 香港
    Macao      = "Asia/Macao"          // 澳门
    Taipei     = "Asia/Taipei"         // 台北
    Tokyo      = "Asia/Tokyo"          // 东京
    HoChiMinh  = "Asia/Ho_Chi_Minh"    // 胡志明市
    Hanoi      = "Asia/Hanoi"          // 河内
    Saigon     = "Asia/Saigon"         // 西贡 (胡志明市)
    Seoul      = "Asia/Seoul"          // 首尔
    Pyongyang  = "Asia/Pyongyang"      // 平壤
    Bangkok    = "Asia/Bangkok"        // 曼谷
    Dubai      = "Asia/Dubai"          // 迪拜
    Qatar      = "Asia/Qatar"          // 卡塔尔
    Bangalore  = "Asia/Bangalore"      // 班加罗尔
    Kolkata    = "Asia/Kolkata"        // 加尔各答
    Mumbai     = "Asia/Mumbai"         // 孟买
    MexicoCity = "America/Mexico_City" // 墨西哥城
    NewYork    = "America/New_York"    // 纽约
    LosAngeles = "America/Los_Angeles" // 洛杉矶
    Chicago    = "America/Chicago"     // 芝加哥
    SaoPaulo   = "America/Sao_Paulo"   // 圣保罗
    Moscow     = "Europe/Moscow"       // 莫斯科
    London     = "Europe/London"       // 伦敦
    Berlin     = "Europe/Berlin"       // 柏林
    Paris      = "Europe/Paris"        // 巴黎
    Rome       = "Europe/Rome"         // 罗马
    Sydney     = "Australia/Sydney"    // 悉尼
    Melbourne  = "Australia/Melbourne" // 墨尔本
    Darwin     = "Australia/Darwin"    // 达尔文
)

var (
    offsetZone = &ZoneCache[int]{cache: make(map[int]*time.Location, 100)}
    fixedZone  = &ZoneCache[zoneKey]{cache: make(map[zoneKey]*time.Location, 100)}
)

type zoneKey struct {
    name   string
    offset int
}

type ZoneCache[K int | zoneKey] struct {
    sync.RWMutex
    cache map[K]*time.Location
}

func (c *ZoneCache[K]) Get(name string, k K) (loc *time.Location) {
    // 获取偏移量
    var off int
    switch v := any(k).(type) {
    case zoneKey:
        off = v.offset
    case int:
        off = v
    }

    if off == 0 {
        if name == "" || name == UTC {
            return time.UTC
        }
        if name == Local {
            return time.Local
        }
    }

    if off < -86400 || off > 86400 {
        // 这里必须分配内存，因为不能返回 nil。
        // 但因为没有写入 map，所以攻击者无法通过这个撑爆我们的内存。
        return &time.Location{}
    }

    c.RLock()
    if loc, _ = c.cache[k]; loc != nil { // OK
        c.RUnlock()
        return
    }
    c.RUnlock()

    // 加写锁
    c.Lock()
    defer c.Unlock()

    // 🔥 第二次检查 (必须)
    if loc, _ = c.cache[k]; loc != nil { // OK
        return
    }

    loc = time.FixedZone(name, off)
    c.cache[k] = loc
    return
}

// NewZone 返回指定时区
// name: 时区名称，offset: 固定偏移小时数
func NewZone(name string, offset ...int) *time.Location {
    var off int
    if len(offset) != 0 {
        off = offset[0]
    }
    return fixedZone.Get(name, zoneKey{name: name, offset: off})
}

// NewOffset 返回指定秒数偏移的固定时区
func NewOffset(offset int) *time.Location {
    return offsetZone.Get("", offset)
}
