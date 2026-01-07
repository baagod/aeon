# thru

Golang 时间包

## 安装

```shell
go get -u github.com/baagod/thru@latest
```

## 使用

### 创建时间

使用 `thru` 包中提供的函数创建时间对象。

```go
now := thru.Now()                              // 现在。通过快捷函数创建
thru.New(time.Now())                           // 通过 time.Time 创建
thru.Unix(now.Unix())                          // 通过秒时间戳创建
thru.Unix(now.Unix(3))                         // 通过毫秒时间戳创建
thru.Parse("2023-03-19")                       // 通过解析自动字符串创建
thru.ParseByLayout("2006-01-02", "2023-03-19") // 通过布局格式创建
```

### 偏移时间

主要使用 `AddYear(ymd ...int)` 方法进行操作。该函数跟 `time.AddDate()` 相似，但是在不传入参数 d (要偏移的天数) 时，年、月不会溢出。
`y, m, d` 可以为负数。

现在假设当前时间是：`2023-01-31 13:14:15 +0800 CST` ( 中国标准时间 )

```go
now.AddYear()       // 添加一年
now.AddMonth(-2, 3) // 在当前时间上偏移 -2月3天
now.AddDay(3)       // 三天后
now.Add(time.Hour)  // 1小时后
```

指定 `d` 参数时可能会导致日期溢出，例如：

```go
now.AddYear(1, 1, 0)
```

该函数在指定 `d` 参数时的行为和 `time.AddDate()` 完全相同，因此返回 `2024-03-02 13:14:15`。这导致我们以为结果和预期不符，其实就是日期溢出了。

让我来尝试解释一下，这是因为 2024 年 2 月一共有 29 天。当我们预期从 2023-01-31 移动到 2024-02-31 时，实际上天数已经溢出 2
天，结果就是它跑到 2024-02-29 时发现后面还有 2 天要继续增加，最终返回 2024-03-02。

### 选择时间

主要使用 `Go(y int, md ...int)` 方法进行操作。该函数偏移 `±y` 年并选择 `m` 月 `d` 日，如果 `m, d`
为负数，则从最后的月、日开始选择，并且不会溢出。例如：

- `Go(1, 13, 1)` 会返回 `2024-12-01 13:14:15`。
- `Go(1, 2, 31)` 返回二月最后一天 `2024-02-29 13:14:15`。

```go
now.Go(-1, 11, -3) // 去年11月最后三天
now.GoMonth(-1, 2) // 年末第二天
now.GoDay(10)      // 本月10号
```

为了扩展，我还添加了一个 `now.GoYear(y int, md ...int)` 方法。跟 `Go()` 一样，但 `y` 指定为确切年份而非偏移。例如：

```go
now.GoYear(2024, 03, 15) // 返回 2024-03-15 13:14:15
```

但是这就跟直接创建差不多了？有用吗？也许吧。

### 开始时间、结束时间

我们经常需要定位要时间边界，因此定义开始时间和结束时间的相关函数。

开始时间主要使用 `Start(ymd ...int)` 进行操作，它和 `AddYear()` 类似，但是返回开始时间。

#### 周起始日配置

`thru` 支持自定义周起始日，默认为周一。提供两种配置方式：

**1. 全局配置：**
```go
// 设置全局默认周起始日为周日
thru.DefaultWeekStartsAt = time.Sunday

// 所有新创建的 Time 实例都会使用全局周起始日
t := thru.Now()  // 使用全局周起始日
```

**2. 实例覆盖：**
```go
// 创建临时使用周一起始日的实例（不影响全局配置）
t := thru.Now().WithWeekStartsAt(time.Sunday)
```

**支持的周起始日：**
- `time.Sunday`：周日
- `time.Monday`：周一（默认）
- `time.Tuesday`：周二
- `time.Wednesday`：周三
- `time.Thursday`：周四
- `time.Friday`：周五
- `time.Saturday`：周六

**示例：**
```go
// 默认周一
now := thru.Now()
now.StartWeek()  // 返回本周周一 00:00:00

// 全局设置为周日
thru.SetWeekStartsAt(time.Sunday)
now.StartWeek()  // 返回本周周日 00:00:00

// 实例覆盖（不影响全局）
now.WithWeekStartsAt(time.Saturday).StartWeek()  // 临时使用周六作为周起始日
now.StartWeek()  // 恢复使用全局配置
```


```go
now.Start()         // 今年开始时间 2023-01-01 00:00:00
now.Start(-1)       // 去年开始时间 2022-01-01 00:00:00
now.Start(1)        // 明年开始时间 2024-01-01 00:00:00
now.Start(1, 1, 1)  // 偏移后的开始时间 2024-03-03 00:00:00

now.StartWeek(-1) // 上周开始时间
now.StartWeek()   // 本周开始时间
now.StartWeek(1)  // 下周开始时间
now.StartWeek(2)  // 两周后的星期一开始时间
```

另外还提供以下方便方法，和 `Start(ymd ...int)` 一样：

- `StartMonth(md ...int)`: 从 `m` 开始。
- `StartDay(d int)`: 从 `d` 开始。

还有与开始时间 `Start()` 系列方法对应的所有结束时间函数: `End`, `EndMoth`, `EndDay`, `EndWeek`。

### 获取时间

获得各种时间的函数，例如要获得年、月、日、时、分、秒等。

```go
Year() int   // 今年
Month() int  // 本月
Day() int    // 返回本月的第几天
Hour() int   // 返回小时 [0, 23]
Minute() int // 返回分钟 [0, 59]
Second() int // 返回秒 [0, 59]
Second(3)int // 毫秒 (3位)
Second(6)int // 微秒 (6位)
Second(9)int // 纳秒 (9位)

Clock() (hour, min, sec int) // 一天中的小时、分钟和秒
YearDay() int   // 一年中的第几天，非闰年范围 [1, 365]，闰年范围 [1, 366]。
Days() int      // 本年总天数
MonthDays() int // 本月总天数

Unix()  int64  // 秒时间戳 (10位)
Unix(3) int64  // 毫秒时间戳 (13位)
Unix(6) int64  // 微秒时间戳 (16位)
Unix(9) int64  // 纳秒时间戳 (19位)

Round(d time.Duration) Time    // 将时间四舍五入到最近接 d 的时间 (详情查阅注释)
Truncate(d time.Duration) Time // 永远向下取整到最近接 d 的时间 (详情查阅注释)

UTC() Time                  // 返回 UTC 时间
Local() Time                // 返回本地时间
In(loc *time.Location) Time // 返回指定的 loc 时间
Location() *time.Location   // 返回时区信息
Time() Time                 // 返回 time.Time
```

### 比较时间

比较时间差主要由以下两个函数返回，可指定要比较的时间单位。

- `DiffIn(u Time, uint string) float64`
- `DiffAbsIn(u Time, uint string) float64`

```go
DiffIn('y') // 返回年差
DiffIn('M') // 返回月差
DiffIn('d') // 返回天差
DiffIn('h') // 返回时差
DiffIn('m') // 返回分差
DiffIn('s') // 返回秒差
```

另外补充具有和 `time` 包中相同的行为的函数：

- `Before(u Time) bool`: 返回 `t < u`
- `After(u Time) bool`: 返回 `t > u`
- `Equal(u Time) bool`: 返回 `t == u`
- `Compare(u Time) int`: 如果 t 小于 u，返回 -1；大于返回 1；等于返回 0
- `Sub(u Time) time.Duration`: 返回 `t - u` 的时间差
- `thru.Since(Time)`: 返回自 `t` 以来经过的时间。它是 `Now().Sub(t)` 的简写。
- `thru.Until(Time)`: 返回直到 `t` 的持续时间。它是 `t.Sub(Now())` 的简写。

### 时间字符串

```go
now.String()   // 返回日期时间字符串: 2023-03-19 15:05:27
now.Format("2006年01年02日") // 指定时间布局格式返回
```

### 其他

```go
IsZero() bool // 返回时间是否零时，即 0001-01-01 00:00:00 UTC。
ZeroOr(u Time) // 使用 u 代替零时 (isZero() is true) 时间
thru.IsLeapYear(y int) // 返回 y 是否闰年
thru.DaysIn(y int, ...m int) // 返回 y 年天数或 y 年 m 月天数
```
