# Aeon: 一个试图重塑时间操作体验的 Go 时间库，它零分配内存，将时间视为 “容器” 而非 “偏移量”。(寻求对这种思维模式的反馈)

**大家好，r/golang 的各位，我是新来的。**

我之前从未在 Reddit 上发过帖。最近，我写了一个 Go 时间库，我个人觉得非常实用。我问了一个 AI 程序，想把它分享给国际社区，它推荐我来这里，说这里聚集了最优秀的 Go 开发者。所以，我就来了。

我知道 AI 生成的内容或推荐可能不太受欢迎，如果我的帖子违反了任何不成文的规则，我深表歉意。我不是来刷屏的 —— 我真心觉得这个库有其独特之处。

**致版主**：我还在学习使用 Reddit。希望在决定删除帖子之前，你们能允许我介绍一下这个库。

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

**A. 认知的错位**

在人类的直觉里，时间是 **层级化** 的：我们会说 “下个月的第三个周五”。 但在 `time.Time` 的逻辑里，时间是 **线性** 的：它是纳秒的累加。

这就导致了极其严重的认知摩擦。试想一下，如果你想找到 "下个季度的最后 n 天"，使用标准库你必须进行一连串的 “心理算术题”：

* 先算出下个季度是几月？
* 那个月有多少天？是闰年吗？
* `AddDate(0, 3, 0)` 会不会因为从 31 号出发而跳到了下下个月？

代码变成了充满魔法数字 `(0, 1, -1)` 的线性代数题，而不是业务逻辑的表达。

**B. 包装库的 “原罪”：内存分配**

为了解决易用性问题，社区出现了很多优秀的包装库 (比如 `Now`, `Carbon`)。它们提供了链式调用，读起来很舒服。但是我无法忍受它们底层的实现方式：每一次链式调用，都在堆上分配新对象！

```go
// 大多数包装库的噩梦
// New() -> Alloc
// AddMonth() -> Alloc
// StartOf() -> Alloc
Carbon.Now().AddMonth(1).StartOfWeek() // 3 次堆分配！
```

这在一个高吞吐的并发系统中，这些细碎的 GC 压力是不可原谅的！

**C. 功能的臃肿**

我又看了一眼 `Carbon` 库，它太 “重” 了。我说 “重” 不是因为它支持的功能太多了，而是他没有将那些所有高度相似的行为系统化内聚和关联。

```go
IsSameYear(t)
IsSameMonth(t)
IsSameDay(t)
IsSameHour(t)
IsSameMinute(t)
IsSameSecond(t)

Between(start, end) // !
BetweenIncludedStart(start, end) // [
BetweenIncludedEnd(start, end) // ]
BetweenIncludedBoth(start, end) // []

Max(t1, t2)
Min(t1, t2)
Closest(t1, t2)
Farthest(t1, t2)

AddMonthsNoOverflow(1)
AddQuartersNoOverflow(1)
AddYearsNoOverflow(1)
```

我不需要为每个有着高度相似的行为单独定义 API，那是 **“穷举法”**。我只需要一把能精准解剖时间的手术刀，切开这一切混乱的根源！

如果是以下这样会如何？

```go
// u: aeon.Year, aeon.Month, aeon.Day..
t.IsSame(u Unit, target t) bool

// bound: "!", "[", "]", "[]"
t.Between(start, end Time, bound ...string) bool

// op: ">", "<", "+", "-"
Pick(op string, times ...Time) Time

ByMonth([aeon.Overflow], 1) // Default: NoOverflow
GoMonth(aeon.Ord, -1, 5) // Last Friday of the month
StartWeekday(5, 18) // This Friday at 18:00 (Happy Hour)
```

这就是我崩溃的临界点。我意识到，我不仅想要更好用的 API，我还想要 Zero-Alloc 的极致性能。我想要一个能像指针一样在时间轴上飞跃，却不留下任何垃圾回收负担的东西。

于是，Aeon 诞生了。

---

... 这里写我最初实现的一些思考到使用级联架构的范式转变

技术和使用细节..
