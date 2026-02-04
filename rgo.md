# Aeon: 一个试图重塑时间操作体验的 Go 时间库，它零分配内存，将时间视为 “容器” 而非 “偏移量”。(寻求对这种思维模式的反馈)

**大家好，r/golang 的各位，我是新来的。**

我从未在 Reddit 上发过帖。最近，我写了一个 Go 时间库，我个人觉得非常实用。我问了一个 AI 程序，想把它分享给国际社区，它推荐我来这里，说这里聚集了最优秀的 Go 开发者。所以，我就来了。

我知道 AI 生成的内容或推荐可能不太受欢迎，如果我的帖子违反了任何不成文的规则，我深表歉意。我不是来刷屏的 —— 我真心觉得这个库有其独特之处。

**致版主**：我还在学习使用 Reddit。希望在决定删除这个帖子之前，你能允许我介绍一下这个库。

---

> 记不清是从哪一刻开始了，也许最初是源于对 `time.Time` 和现有时间库的不满 —— 我产生了一个近乎荒诞的执念：**自己写一个
Go 时间库？**
>
> 这一切始于那个简陋、甚至可以说是有些 "丑陋" 的前身 `thru`，我决定对其系统化更新与重构。无数个灵感与想象如烟火般炸裂，才最终让它完成了跨越维度的蜕变。
>
> 我将其正式命名为：**Aeon**。在古老的哲学中，Aeon 代表着 "永恒" 与 "层叠的维度"。
>
> 我选择这个名字，是因为它代表了时间更本质的逻辑 —— 时间不是一条细长的直线，它是流动的、是可以被嵌套和穿透的宇宙。

---

## 为什么重新 “造轮子”？

### 现有方案的困境：在 "线性算术" 与 "堆内存陷阱" 之间挣扎

Go 的标准库 `time.Time` 是一个工程奇迹，它精准、稳定、且线程安全。但当我们试图用它来处理业务逻辑时，痛苦就开始了。

1. **认知错位**

   在人类的直觉里，时间是 **层级化** 的：我们会说 “下个月的第三个周五”。 但在 `time.Time` 的逻辑里，时间是 **线性** 的：它是纳秒的累加。
   
   这就导致了极其严重的认知摩擦。试想一下，如果你想找到 "下个季度的最后 n 天"，使用标准库你必须进行一连串的 “心理算术题”：
   
   * 先算出下个季度是几月？
   * 那个月有多少天？是闰年吗？
   * `AddDate(0, 3, 0)` 会不会因为从 31 号出发而跳到了下下个月？
   
   代码变成了充满魔法数字 `(0, 1, -1)` 的线性代数题，而不是业务逻辑的表达。

2. **包装库的 “原罪”：内存分配**

   为了解决易用性问题，社区出现了很多优秀的包装库 (例如 `Now`, `Carbon`)。它们提供了链式调用，读起来很舒服。但是我无法忍受它们的底层实现方式：每一次调用 (或是链式调用)，都在堆上分配新对象！
   
   ```go
   // 大多数包装库的噩梦
   // New() -> Alloc
   // AddMonth() -> Alloc
   // StartOf() -> Alloc
   Carbon.Now().AddMonth(1).StartOfWeek() // 3 次堆分配！
   ```
   
   在一个高吞吐的并发系统中，这些细碎的 GC 压力是不可原谅的！

3. **功能臃肿**

   我又看了一眼 `Carbon` 这个库，它太 “重” 了。我说 “重” 不是因为它支持太多功能，而是它没有将那些所有高度相似的行为，进行系统化抽象和内聚。
   
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
   
   我不想背诵 300 个方法名称，那是 **“穷举法”**，是 **“打补丁”**。我只需要一把能精准解剖时间的神剑，斩断这一切混乱根源！
   
   想象一下，如果是以下这样定义 API 会如何？
   
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
   
   这就是我崩溃的临界点。我意识到，我不仅想要更好用的 API，我还要 Zero-Alloc 的极致性能，我想要像指针一样在时间轴上飞跃，不留下任何垃圾回收。

于是，Aeon 诞生了。

## 漫长的探索：只是为了解决一个 “溢出”

最初，我根本没想过什么级联寻址、更没管性能和零分配。我只有一个非常单纯愿望：**让我操作月份时天数不要溢出。**

于是我创建了 Aeon 的前身 `thru` 并简单实现了这个功能。当时我意识到除了 “增加”，我可能还需要直接 “设置”，于是有了最初的 `Go` 方法原型。例如：

`GoMonth(1, 2)`：直接把时间设为 1 月 2 日，并保持年、分、秒不变，最重要的是 —— 抑制天数溢出。

你可以看到，我为了这一个小小的 “不溢出”，做了多少琐碎工作。在 Stack Overflow 上，这是一个经久不衰的抱怨。无数人在问：“为什么我只加了一个月，日期却来到了下下个月？”。

**但噩梦才刚刚开始。**

当我试图把这套 “补丁式” 的逻辑推广到周、季度、年甚至更复杂的跨世纪计算时，我发现代码开始失控了。

我陷入了 `if-else` 的地狱：为了让日期显示正常，我不仅要判断闰年、大小月、开始和结束边界，还要处理跨年周、季度末的边界... 当我好不容易填补好了 “月份” 漏洞，“季度” 的传参又崩溃了！

整个方法块就像一团乱麻，逻辑支离破碎。但那时的我并不知道，我正在接近一个更本质的真相...

## 级联架构雏形：原子操作

于是，我停下来了。我不再尝试将所有变长传参，都一次性算好时间再返回，我只做一件事：只处理 `Start` 这个方法，并且只传一个参数，或不传 (默认为 `0`) 表示获取指定单位的开始时间。

> [!IMPORTANT]
> 
> **我需要验证在单个维度的原子操作下，逻辑是否依然崩溃。如果连单一传参都出错，那就证明我的算数逻辑从根本上就是错误的。**
> 
> *(这句话非常重要，它是 Aeon 整个导航系统的基石。)*

我定义 `Start` 方法的原型为 `t.Start(u Unit, ...n)`。例如，如果我想获取某月份开始时间，就会传入 `Start(aeon.Month, 5)`。意图极其纯粹：定位到 5 月，然后将下级单位 (日、时、分、秒、纳秒) 全部抹平。

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

就这样，搞定？如果传 `0`，我会获取当前月份的开始时间，而不是重新设置新值。但是我怎么知道具体增加了几年几月呢？假如传入 `m=12` 超出一年怎么办？

为此，我为月份设计了一套 ‘自动进位协议’。即便传入了不合常理的 `13` 月或 `-1` 月，它都能像水流一样自动溢出到年份上，最终归位到正确的刻度。这种基于数学模运算的稳定性，远胜于碎裂的 `if-else`。

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

我会在 `switch` 之后调用它：

```go
y, m = addMonth(y, m, n)
```

这样，如果你调用 `addMonth(y, m, 12)`，它可能会返回 `y=y+1, m=1` (如果你传入的 `m=1`)，这确保了我返回给 `time.Date()` 的年月永远正确。

但是此时我还需要处理月份溢出，我应该怎么做？为此我写了一个 “获取每个月最大天数” 的方法：

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

这样我就能获取 `y` 年 `m` 月的最大天数。其中 `maxDays` 被定义为：

```go
// maxDays 每个月的最大天数
maxDays = [13]int{1, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
```

我会在 `applyAbs` 的最后统一处理天数溢出问题。例如：

```go
// 统一溢出检查：只需判断当前操作的单位是否在 “月” 级及以上
if u <= Month {
    if dd := DaysIn(y, m); d > dd {
        d = dd
    }
}
```

就这样，我彻底杀死了困扰开发者几十年的 “月份溢出” 噩梦！

此时还有最后一个重要问题需要处理，那就是，我如何归零所有下级时间？

例如，当我调用 `t.Start(Month, 5)`，那么需要将 `day - ns` 初始化，即：`y-05-01 00:00:00.000..`。而如果我调用 `t.Start(Year, 5)`，需要返回 `y-01-01 00:00:00.000..`。

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

---

将所有这一系列方法串联起来，我得到了一个原子化时间定位方法：

*代码已经过简化，但原理是一样的，详细查看 [Aeon](https://github.com/baagod/aeon/blob/main/opus.go) 仓库实现 (**真的不是推销！**)。*

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

之后，我按照该逻辑，陆续增加了对更多 `Unit` 的 `case` 处理和各种边界的处理 (如结束边界)，确保其在单一参数下返回正确，并反复测试直到稳定。

此时，我拥有了一个绝对稳定的 **时间原子化操作引擎**，它能做什么呢？

## 维度的坍缩：级联架构


2. 情绪高潮处：级联引擎展示后（新增 ⚡）
   在下一章 “维度的坍缩：级联架构” 结束时。
   当读者读完你如何把复杂的 if-else 变成一套优雅的流式状态机，并在一次原子操作中完成时，他们的好奇心达到了顶点。
   建议写法：
> “这一瞬间，我意识到 Aeon 的灵魂诞生了。你可以点击这里直接看这套引擎的实现：Github: baagod/aeon (https://github.com/baagod/aeon)”
3. 转化位置：结尾（必填 👑）
   全篇结束。
   这是给那些陪你走完整个心路历程的“忠实读者”的。这时候他们已经对你的人设和技术产生了极强的认同感。
   建议写法：
> “这就是 Aeon 的故事。它不仅是一个库，更是我对时间逻辑的一次重新审判。
> 如果你觉得这种‘容器模型’对你有启发，或者只是想帮我测试一下它的极限性能，欢迎来 Github 踩踩。
> Github: https://github.com/baagod/aeon”
---
🛡️ 架构师的“阴谋”：锚点策略
别只是发一个裸链接，给链接加点 “情绪价值”：
*   不要只写：Github: https://...
*   尝试写成：Check the Zero-Alloc implementation: [baagod/aeon](https://github.com/baagod/aeon)
    总结：
1.  首部：给“急性子”。
2.  中部（级联架构后）：给“技术控”。
3.  尾部：给“追剧党”。
    这样布局，既保证了转化率，又不会显得廉价。你需要我现在帮你把
