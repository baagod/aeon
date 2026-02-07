# Aeon：一个零分配的 Go 语言时间库，它将时间视为 “容器” 而非 “偏移量”。

> 记不清是从哪一刻开始了，也许最初是源于对 `time.Time` 和现有时间库的不满 —— 我产生了一个近乎荒诞的执念：**为什么不自己写一个 Go 时间库呢？**
>
> 在古老的哲学中，Aeon 代表着 “永恒” 与 “层叠的维度”。
>
> 我选择这个名字，是因为我想表达时间不同的逻辑 —— **时间不是一条细长的直线，它是流动的、是可以被嵌套和穿透的宇宙。**


## 为什么重新 “造轮子”？

现有方案的困境：在 “线性算术” 与 “堆分配” 中挣扎。Go 的标准库 `time.Time` 是一个工程奇迹，它精准、稳定、且线程安全。但当我们试图用它来处理业务逻辑时，痛苦就开始了。

### 认知错位

在人类的直觉里，时间是 **层级化** 的：我们会说 “下个月的第三个周五”，“本月最后一天”。但在 `time.Time` 的逻辑里，时间是 **线性** 的：它是纳秒的累加。

这就导致了极其严重的认知冲突。试想一下，如果你想找到 “下个季度的最后 n 天”，使用标准库你必须进行一连串的 “心理预算”：

* 先算出下个季度是几月？
* 那个月有多少天？是闰年吗？
* `AddDate(0, 3, 0)` 会不会因为从 31 号出发而跳到了下下个月？

代码变成了充满魔法数字 `(0, 1, -1)` 的线性代数题，而不是业务逻辑的表达。

### 内存分配

为了解决易用性问题，社区出现了很多优秀的包装库 (例如 `now`, `carbon`)。它们提供了链式调用，读起来很舒服。但是我无法忍受它们的底层实现方式：每一次调用，都在堆上分配内存！

```go
// 大多数包装库的噩梦
// New() -> Alloc
// AddMonth() -> Alloc
// StartOf() -> Alloc
Carbon.Now().AddMonth(1).StartOfWeek() // 3 次堆分配！
```

在一个高吞吐的并发系统中，这些细碎的 GC 压力是不可原谅的！

### 功能臃肿

我又看了一眼 `carbon` 这个库，它太 “重” 了。我说 “重” 不是因为它支持太多功能，而是它没有将那些所有高度相似的行为，进行系统化抽象和内聚。

```go
IsSameYear(t)
IsSameMonth(t)
IsSameDay(t)
IsSameHour(t)
IsSameMinute(t)
IsSameSecond(t)

Between(start, end) // =
BetweenIncludedStart(start, end) // [
BetweenIncludedEnd(start, end) // ]
BetweenIncludedBoth(start, end) // !

Diff[Abs]InYears()
Diff[Abs]InMonths()
Diff[Abs]InWeeks()
Diff[Abs]InDays()
Diff[Abs]InHours()
Diff[Abs]InMinutes()
Diff[Abs]InSeconds()

Max(t1, t2)
Min(t1, t2)
Closest(t1, t2)
Farthest(t1, t2)

AddMonthsNoOverflow(1)
AddQuartersNoOverflow(1)
AddYearsNoOverflow(1)
```

我不想背诵 300 个方法名称，那是 **“穷举法”**，是 **“打补丁”**。我需要一把能够精准解剖时间的神剑，斩断一切混乱根源。

### Aeon 诞生

试想，如果是这样定义 API 会如何？

```go
// u: aeon.Year, aeon.Month, aeon.Day..
t.IsSame(u Unit, target t) bool

// bound: '=', '!', '[', ']'
t.Between(start, end Time, bound ...byte) bool

// unit: 'y', 'M', 'd', 'h', 'm', 's'
t.Diff(u Time, unit byte, abs ...bool) float64

// op: '>', '<', '+', '-'
Pick(op byte, times ...Time) Time

ByMonth([aeon.Overflow], 1) // Default: NoOverflow
GoMonth(aeon.Ord, -1, 5) // Last Friday of the month
StartWeekday(5, 18) // This Friday at 18:00 (Happy Hour)
```

这就是我崩溃的临界点。我意识到，我不仅想要更好用的 API，我还要 **`Zero-Alloc`** 的极致性能，我想要像指针一样在时间轴上飞跃，不留下任何垃圾回收。

于是，Aeon 诞生了。

```sql
Benchmark       | ns/op | allocs x B/op | speedup

New             
Aeon            | 18.6  | 0             | x74
Carbon          | 1376  | 13x1600

Now             
Aeon            | 7.8   | 0             | x177
Carbon          | 1384  | 13x1600

From Unix       
Aeon            | 3.6   | 0             | x383
Carbon          | 1380  | 13x1600

From Std        
Aeon            | 5.0   | 0             | x323
Carbon          | 1619  | 13x1600

Parse (Compact) 
Aeon            | 23.3  | 0             | x195
Carbon          | 4561  | 85x3922

Parse (ISO)     
Aeon            | 19.6  | 0             | x91
Carbon          | 1794  | 15x1697

Start / End     
Aeon            | 56.4  | 0             | x20
Carbon          | 1141  | 7x1440

Add (Offset)    
Aeon            | 56.5  | 0             | x2.5
Carbon          | 142   | 2x128

Set (Position)  
Aeon            | 58.7  | 0             | x2.6
Carbon          | 156   | 2x128
```

> [!NOTE]
>
> 以上数据在 **未使用** 变长参数的单一原子操作下测得。并且 **即使链式调用，仍然 Zero-Alloc。逻辑越复杂，Aeon 的领先倍数就越惊人。**

> 如果你只想快速认识该库，那么到这里就可以结束了。你可以查看 [Aeon](https://github.com/baagod/aeon) 和它的 [完整文档](https://zread.ai/baagod/aeon/1-overview) 了解更多。衷心感激您的支持！

然而，如果你还想看看 Aeon 是如何通过一步步演进、构思及诞生的，请继续。


## 只是为了解决一个 “溢出”，我意外构建出一个 “时间容器” 模型。

最初，我根本不知道什么级联索引、时间容器.. 更没管什么性能与零分配。我只有一个非常单纯愿望：**让我操作月份时天数不要溢出..**

### 漫长的探索

于是，我创建了 Aeon 的前身 `thru` 并简单实现了该功能。当时我意识到除了 “增加”，我可能还需要直接 “设置”，于是就有了 `Go` 方法的原型。例如 `GoMonth(1, 2)`：直接把时间设为 1 月 2 日，并保持年、分、秒不变，最重要的是 —— 抑制月份溢出。

你可以看到，我为了这一个小小的 “不溢出”，做了多少琐碎工作。在 Stack Overflow 上，这是一个经久不息的抱怨。无数人在问：“为什么我只加了一个月，日期却来到了下下个月？”。

**但噩梦才刚刚开始。**

当我试图把这套 “补丁式” 的逻辑推广到周、季度、年甚至更复杂的跨世纪计算时，代码失控了。

我陷入了 `if-else` 的地狱：为了让日期显示正常，我不仅要判断闰年、大小月、开始和结束边界，还要处理跨年周、季度末的边界... 当我好不容易填补了 “月份” 漏洞，“季度” 传参又崩溃了！

整个方法块逻辑支离破碎。但那时我并不知道，我正在接近一个更本质的真相...

### 原子操作

于是，我停下来了。我不再尝试将所有的变长传参都算好时间再返回，我只做一件事，那就是：**只处理 `Start` 这个方法，并规定只能传一个参数。**

> [!IMPORTANT]
>
> **我需要验证在单个维度的原子操作下，如果逻辑依然崩溃，那就证明我的算数从根本上就是错误的。** *(这句话非常重要，它是构成 Aeon 整个导航系统的基石。)*

我定义 `Start` 方法的原型是 `t.Start(u Unit, ...n)`。例如，如果我想获取某月份的开始时间，就会调用 `Start(aeon.Month, 5)`，意图极其纯粹：定位到 5 月，然后将它所有下级单位 (日、时、分、秒、纳秒) 抹平。

在这个极简模型下，我终于从那些琐碎的 `if-else` 中抽离出来，专注于设置传入的每个单位逻辑。我在方法中定义了一个 `switch-case`：

```go
func applyAbs(u Unit, y, m, d int) Time {
    switch u {
    case Year:  // 只处理年定位逻辑
    case Month: // 只处理月定位逻辑
        if n > 0 {
            m = n
        } else if n < 0 {
            m = 13 + n // 负数，反向索引
        }
    case ..
    }
}
```

如果传 `0`，我会保持在当前月份，而不是设置新值。但是我怎么知道具体增加了几年几月呢？如果传入 `m=13` 超出了一年怎么办？

为此，我设计了一个月份 **自动进位协议**。即便传入了不合常理的 `13` 或 `-1` 月，它都能像水流一样，自动溢出到年份上，并最终归位到正确刻度。

```go
// addMonth 计算月份增加或减少后的年月 (处理年进位/借位)
func addMonth(u Unit, y, m, n int) (int, int) {
    months := m + n
    y += (months - 1) / 12
    if m = (months-1)%12 + 1; m <= 0 {
        m += 12
        y--
    }
    return y, m
}
```

我会在 `switch` 之后调用它：`y, m = addMonth(y, m, n)`。

这样，如果你调用 `addMonth(y, m, 12)`，它可能会返回 `y=y+1, m=1`，这确保了我返回给 `time.Date()` 的年月永远正确。

---

但此时我还需要处理月份溢出，我应该怎么做？答案是，我写了一个 “获取月份最大天数” 方法。

```go
// DaysIn 返回 y 年 m 月最大天数，如果忽略 m 则返回 y 年总天数。
//
//   - 1, 3, 5, 7, 8, 10, 12 月有 31 天；4, 6, 9, 11 月有 30 天。
//   - 平年 2 月有 28 天，闰年 29 天。
func DaysIn(y int, m ...int) int {
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
```

这样，我就能获得 `y` 年 `m` 月份的最大天数了，并且在 `applyAbs` 的最后处理：

```go
// 统一溢出检查：只需判断当前操作的单位是否在 “月” 级及以上
if u <= Month {
    if dd := DaysIn(y, m); d > dd {
        d = dd
    }
}
```

就这样，我彻底终结了 “月份溢出” 噩梦！

---

至此，仅剩下最后一个问题。那就是，我如何归零所有下级时间？

例如，当我调用 `t.Start(Month, 5)`，那么需要将 “天” 至  “纳秒” 初始化，即：`y-05-01 00:00:00.000..`。而如果我调用 `t.Start(Year, 5)`，需要返回 `y-01-01 00:00:00.000..`。

我想到了一个方法来解决这个问题，在 `switch` 的 **最后** 统一处理好时间边界再返回：

```go
// align 执行最终的时间分量对齐 (归零或置满)
func align(u Unit, y, m, d, h, mm, sec, ns int) (int, int, int, int, int, int, int) {
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
    return y, m, d, h, mm, sec, ns
}
```

将所有这一系列方法串联起来，我得到了一个时间原子化定位方法 (代码已经过简化)：

```go
func applyAbs(u Unit, y, m, d int) (int, int, int) {
    switch u {
    case Year:  // 只处理年定位逻辑
    case Month: // 只处理月定位逻辑
        if n > 0 {
            m = n
        } else if n < 0 {
            m = 13 + n // 负数，反向索引
        }
    case ..
    }
    
    y, m = addMonth(y, m, n)
    if u <= Month {
        if dd := DaysIn(y, m); d > dd {
            d = dd
        }
    }
    
    return align(u, y, m, d)
}
```

之后，我按照该逻辑，陆续增加了对更多 `Unit` 的 `case` 处理和各类边界的处理 (如结束边界)，确保其在单一参数下返回正确，并反复测试直到稳定为止。

终于，我拥有了一个绝对稳定的 **原子化时间操作引擎**，它就像一把从未出鞘的神剑，一但出鞘将震动整个时空！它能做什么呢？

### 维度的坍缩：级联架构

虽然我已经保证处理单一参数绝对正确，但我的终极目标是 **级联**，通过传入变长参数，在这个原子方法中只返回一个指定的 `time.Date`，而中间不创建任何 `Time` 对象辅助。

例如，我想要如下方法。无论级联多少参数，它只能创建一个 `Time` 对象并返回：

```go
GoMonth(1, 5, 3) // 1 月 5 日 3 点
```

思索许久，如果我继续通过传入的所有变长参来改变每个 `case` 的实现，必然又回到最初的令人恐惧那个 `if-else` 地狱。我应该怎么办？

现在我已经有了单一且稳定的时间操作方法，问题是，如何复用？

> 我凝视着屏幕上的 `GoMonth(1, 5, 3)`，时间仿佛凝固住了。一瞬间，仿如神灵附体，我似乎看见 `1, 5, 3` 这三个数字，分别从屏幕中飞了出来。**我知道了！**

> [!IMPORTANT]
>
> 如果我能调用 `GoMonth(1, 5, 3)`，那它和我链式调用 `GoMonth(1).GoDay(5).GoHour(3)` 有什么区别？
>
> 如果我能链式调用 `GoMonth(1).GoDay(5).GoHour(3)`，那它和我在一个循环体中，通过将每次结果返回，传递给下一个循环中的下一个时间单位 (如 Day)，又有什么区别！

一瞬间，我感觉整个时间维度坍缩了，它们变得不再琐碎和复杂，我可以将这一切都串联起来！也就是说，我可能需要在一个循环体中，陆续传入：

```go
i=0: y, m, d = applyAbs(Month, 1, y, m, d)
i=1: y, m, d = applyAbs(Day, 5, y, m, d)
i=2: y, m, d = applyAbs(Hour, 3, y, m, d)
```

这样，如果我将每个级联参数都向下传递，最终就能得到我想要的时间！(*方法已经过简化*)

```go
var years = []Unit{Century, Decade, Year, Month, Day, Hour, Minute, Second, Millisecond, Microsecond, Nanosecond}

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

// cascade 级联时间核心引擎
func cascade(t Time, u Unit, args ...int) Time {
    y, m, d := t.Date()
    h, mm, s := t.Clock()
    ns := t.time.Nanosecond()
    w := t.Weekday()
    sw := t.weekStarts

    seq := u.seq()
    if l := len(seq); len(args) > l {
        args = args[:l]
    }
    
    p, pN := u, args[0] // 父单元及其传值
    for i, n := range args {
        unit := seq[i]
        y, m, d, h, mm, s, ns, w = applyAbs(unit, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
        p, pN = unit, n
    }

    return Time{
        time:       time.Date(y, time.Month(m), d, h, mm, s, ns, t.Location()),
        weekStarts: t.weekStarts,
    }
}
```

这一瞬间，我意识到 Aeon 的灵魂诞生了。整个过程，没有一次内存分配！全都是数字计算！这在整个 CPU 运算中极其高效。

---

基于该原理，我构建了 **四个基础动作**：

1. **`Go`**：精准定位。
2. **`By`**：相对偏移。
3. **`Start`**：边界归零。
4. **`End`**：边界置满。

但是我发现如果只有这些，还不够！

例如，在当前动作中，我无法做到 **先偏移** 到某个时间点，再进行直接定位；反之，我也无法 **先定位** 到某个时间点，再执行相对偏移。

要做到这两点，我还需要为 Aeon 加入两个动作：

1. `At`: 先定位后偏移。例如：`At(5, 1, 1)` ➜ 定位到 5 月，再添加 1 天和 1 小时。
2. `In`: 先偏移后定位。例如：`In(1, 5, 1)` ➜ 明年 5 月 1 日。

**正是这 2 + 4 个完全正交的动作构成了 Aeon 整个导航系统的核心。** 我认为它几乎包含了所有可能的时间落脚点。该如何实现？

我依然使用 `cascade` 函数 (已简化代码)：

```go
for i, n := range args {
    unit := seq[i]
    if i == 0 { // 如果是首个参数，使用定位或偏移方法。
        y, m, d, h, mm, s, ns, w = applyRel(unit, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
    }else {
        y, m, d, h, mm, s, ns, w = applyAbs(unit, p, n, pN, y, m, d, h, mm, s, ns, w, sw)        
    }
    p, pN = unit, n
}
```

就这样，完了.. 解决方案之巧妙令我自己都感到惊讶，这就是 **正交** 的力量！


### 容器系统：时间索引

通过上述 4 个动作，我具备了几乎完整的时间导航系统。但是，我还想在此之上，更近一步。

我的目标是什么？我不想再将时间看作是线性的，例如 `AddDay().GoMonth().Start()`。这样每一次调用，都只是在上一个时间的基础上，做加减法而已，它们本质上仍然是线性的。

而我想要的是 **一个具有统一层级的时间结构**。

---

当我们调用 `Go/ByMonth()`，很自然会觉得这就是在设置本年的月份，包括 `XXDay/Hour..()` 亦如此。在直觉上我们会认为，`Month` 就应该属于 `Year`，`Day` 属于 `Month`，`Hour` 又跑在 `Day` 里。

是的，这才是符合 **人类直觉** 的时间观。它们一个层级套着另一个层级，是一个容器，大单位包含着小单位，小单位又包含子单位。

但是，`Year` 呢？`Year` 的上级又是什么？许多时间库到这里就结束了，**它们认为年是无限的**，所以它们没能给年定义一个归属，年是孤立的。

**但是，在 Aeon 中，我把这个概念彻底贯彻到了整个导航系统中！**

---

我为上述 4 个动作，都定义了一套时间单位方法，从世纪到纳秒。例如这是 `Go` 的：

```go
GoCentury(n ...int) Time
GoDecade(n ...int) Time
GoYear(n ...int) Time
GoMonth(n ...int) Time
GoDay(n ...int) Time
GoHour(n ...int) Time
GoMinute(n ...int) Time
GoSecond(n ...int) Time
GoMilli(n ...int) Time
GoMicro(n ...int) Time
GoNano(n ...int) Time
GoQuarter(n ...int) Time
GoWeek(n ...int) Time
GoWeekday(n ...int) Time
```

就这样，只需选择你想要的时间单位，它就会自动生成该单位到纳秒的 **级联序列**。每个参数会会像水流一样，穿透到下级。

```go
GoMonth() // Month, Day, Hour, .., Nanosecond
GoCentury() // Century, Decade, Year, Month, Day, Hour, .., Nanosecond
GoDecade(2, 5) // 本世纪第 2 个年代第 5 年 = 2025
GoYear(2) // 本年代第 2 年 = 2022
```

这就是 **时间容器**，你无需再记住每一个时间导航方法，你只需要记住 4 个动作加上你想要操作的时间单位。

并且在容器中，还能反向索引，从父容器的末尾开始定位。

```go
GoDay(-2, 23) // 本月最后 2 天 23 点
GoMonth(-1, -3) // 本年最后一个月倒数第三天
```

这极大降低了使用者的心智负担，他无需再为计算各类边界而感到烦恼。

**时间容器索引模型：**

```text
[Millennium]
  └─ [0...9 Century]
       └─ [0...9 Decade]
            └─ [0...9 Year]
                 └─ [1...12 Month]

Example：GoYear(5) Indexing logic
         [-9]       [-8]            [-5]        [-4]             [-1]
2020 ─┬─ 2021 ──┬── 2022 ··· ──┬── [2025] ──┬── 2026 ─┬─ ··· ─┬─ 2029
[0]      [1]        [2]             [5]         [6]              [9]
```

---

结束了吗？

是的。但是，为了不彻底将时间锁死在父容器中，我提供了 6 个顶级方法，允许传递的首个参数去往 **绝对年份**：

```go
Go(2025, 1).StartDay(-1, 23) // 2025-01-31 23:00:00
```

虽然 Aeon 定位的核心是 **“寻址”**，但它同时提供了 `Overflow` 标志，允许时间自然溢出。

```go
GoMonth(aeon.Overflow, 2) // 如果溢出，可能是 03-2/3。
```

## 尾声

> 这就是 Aeon 的故事。愿在你的代码中，时间永远流转，不再停顿。

如果你看到这里，衷心感谢您陪我走完这段心路历程，诚挚感谢！

最后，不仅仅是 **导航**，在 [Aeon](https://github.com/baagod/aeon) 的世界中，还有更多关于对时间的独特看法。
