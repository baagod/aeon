# Aeon: A Zero-Allocation Go Time Library That Treats Time as "Containers" Rather Than "Offsets"

> I can't recall the exact moment it began. Perhaps it originated from my dissatisfaction with `time.Time` and existing time libraries — I developed an almost absurd obsession: **Why not write my own Go time library?**
>
> In ancient philosophy, Aeon represents "eternity" and "nested dimensions."
>
> I chose this name because I wanted to express a different logic of time — time is not a thin straight line. It is fluid, a universe that can be nested and penetrated.

## Why "Reinvent the Wheel"?

The dilemma of existing solutions: struggling between "linear arithmetic" and "heap allocation." Go's standard library `time.Time` is an engineering miracle — precise, stable, and thread-safe. But when we try to use it to handle business logic, the suffering begins.

### Cognitive Mismatch

In human intuition, time is **hierarchical**: we say "the third Friday of next month" or "the last day of this month." But in `time.Time`'s logic, time is **linear**: it's an accumulation of nanoseconds.

This creates a severe cognitive conflict. Imagine you want to find "the last n days of the next quarter." Using the standard library, you must perform a series of "mental calculations":

* First, which month is next quarter?
* How many days are in that month? Is it a leap year?
* Will `AddDate(0, 3, 0)` jump to the month after next because it starts from the 31st?

The code becomes a linear algebra problem filled with magic numbers like `(0, 1, -1)`, rather than an expression of business logic.

### Memory Allocation

To solve usability issues, the community has produced many excellent wrapper libraries (such as `now`, `carbon`). They provide fluent chain calls that read beautifully. But I cannot tolerate their underlying implementation: every call (or chain call) allocates memory on the heap!

```go
// The nightmare of most wrapper libraries
// New() -> Alloc
// AddMonth() -> Alloc
// StartOf() -> Alloc
Carbon.Now().AddMonth(1).StartOfWeek() // 3 heap allocations!
```

In a high-throughput concurrent system, these fragmented GC pressures are unforgivable!

### Feature Bloat

I glanced at the `carbon` library again — it's too "heavy." When I say "heavy," I don't mean it supports too many features. I mean it fails to systematically abstract and coalesce all those highly similar behaviors.

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

I don't want to memorize 300 method names. That's **"exhaustive enumeration"** — that's **"patching."** I need a divine sword that can precisely dissect time, cutting off all chaos at the root.

### Aeon is Born

Imagine — what if we defined the API like this?

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

This was my breaking point. I realized I didn't just want a better API — I wanted Zero-Alloc's ultimate performance. I wanted to leap through the timeline like a pointer, leaving no garbage behind.

And so, Aeon was born.

```bash
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
> The benchmark data above was measured **without** variadic parameters under single atomic operations. And **even with chain calls, it remains Zero-Alloc. The more complex the logic, the more dramatic Aeon's performance advantage becomes.**

> If you only want to quickly understand this library, you can stop here. Check out [Aeon](https://github.com/baagod/aeon) and its [complete documentation](https://zread.ai/baagod/aeon/1-overview) to learn more. Thank you sincerely for your support!

However, if you'd like to see how Aeon evolved step by step — how it was conceived and born — please continue.

## Just To Solve One "Overflow," I Accidentally Built a "Time Container" Model

Initially, I knew nothing of cascading indexes or time containers... I didn't even care about performance or zero allocation. I had only one simple wish: **let me handle months without day overflow..**

### The Long Exploration

So, I created Aeon's predecessor `thru` and simply implemented this functionality. At the time I realized that besides "adding," I might also need to directly "set" values. So the prototype of the `Go` method was born. For example, `GoMonth(1, 2)`: directly set the time to January 2nd while keeping the year, minutes, and seconds unchanged — and most importantly, suppress month overflow.

You can see how much trivial work I did just for this small "no overflow" feature. On Stack Overflow, this is an enduring complaint. Countless people ask: "Why did I only add one month, but the date landed in the month after next?"

**But the nightmare was just beginning.**

When I tried to extend this "patch-style" logic to weeks, quarters, years, and even more complex cross-century calculations, the code spiraled out of control.

I fell into an `if-else` hell. To make dates display correctly, I not only had to handle leap years, month lengths, start and end boundaries — but also cross-year weeks, quarter-end boundaries... When I finally plugged the "month" loophole, the "quarter" parameter crashed!

The entire method logic was fragmented. But at that time, I didn't know I was approaching a more fundamental truth...

### Atomic Operations

So, I stopped. I no longer tried to calculate all variadic parameters before returning. I did only one thing: **handle only the `Start` method, and restrict it to accepting only one parameter.**

> [!IMPORTANT]
>
> **I need to verify that under atomic operations with a single parameter, if the logic still breaks, it proves my arithmetic is fundamentally wrong.**
>
> *(This sentence is very important; it's the cornerstone of Aeon's entire navigation system.)*

I defined the `Start` method prototype as `t.Start(u Unit, ...n)`. For example, if I wanted to get the start time of a certain month, I would call `Start(aeon.Month, 5)` — with crystal-clear intent: locate to May, then flatten all its subordinate units (day, hour, minute, second, nanosecond).

Under this minimalist model, I finally detached from those trivial `if-else` blocks and focused on the logic of setting each passed unit. I defined a `switch-case` in the method:

```go
func applyAbs(u Unit, y, m, d int) Time {
    switch u {
    case Year:  // Only handle year positioning logic
    case Month: // Only handle month positioning logic
        if n > 0 {
            m = n
        } else if n < 0 {
            m = 13 + n // Negative number, reverse indexing
        }
    case ..
    }
}
```

If `0` is passed, I stay in the current month rather than setting a new value. But how do I know exactly how many years and months were added? What if passing `m=13` exceeds a year?

For this, I designed a month **automatic carry protocol**. Even if unconventional `13` or `-1` months are passed, it can automatically overflow to the year like flowing water, and finally return to the correct scale.

```go
// addMonth calculates year and month after adding/subtracting months (handles year carry/borrow)
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

I call it after the `switch`: `y, m = addMonth(y, m, n)`.

This way, if you call `addMonth(y, m, 12)`, it might return `y=y+1, m=1`, ensuring that the year and month I return to `time.Date()` are always correct.

---

But at this point, I still needed to handle month overflow. What should I do? The answer is, I wrote a "get month's maximum days" method.

```go
// DaysIn returns the maximum days in year y and month m, or returns total days in year y if m is ignored.
//
//   - Months 1, 3, 5, 7, 8, 10, 12 have 31 days; 4, 6, 9, 11 have 30 days.
//   - February has 28 days in common years, 29 in leap years.
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

This way, I can obtain the maximum days of month `m` in year `y`, and handle it at the end of `applyAbs`:

```go
// Unified overflow check: just judge whether the current operation's unit is at "month" level or above
if u <= Month {
    if dd := DaysIn(y, m); d > dd {
        d = dd
    }
}
```

Just like that, I completely ended the "month overflow" nightmare!

---

At this point, only one problem remained. How do I zero out all subordinate times?

For example, when I call `t.Start(Month, 5)`, I need to initialize from "day" to "nanosecond," producing: `y-05-01 00:00:00.000..`. But if I call `t.Start(Year, 5)`, I need to return `y-01-01 00:00:00.000..`.

I thought of a way to solve this problem: handle time boundaries uniformly at the **end** of the `switch` before returning:

```go
// align performs final time component alignment (zero or fill)
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

Linking all these methods together, I got a time atomic positioning method (code has been simplified):

```go
func applyAbs(u Unit, y, m, d int) (int, int, int) {
    switch u {
    case Year:  // Only handle year positioning logic
    case Month: // Only handle month positioning logic
        if n > 0 {
            m = n
        } else if n < 0 {
            m = 13 + n // Negative number, reverse indexing
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

Afterwards, following this logic, I successively added `case` handling for more `Unit`s and various boundary conditions (such as end boundaries), ensuring correct returns with single parameters, and tested repeatedly until stable.

Finally, I possessed an absolutely stable **atomic time operation engine**. It's like a divine sword never unsheathed — once unsheathed, it will shake all of spacetime! What can it do?

### Dimensional Collapse: Cascade Architecture

Although I had guaranteed that handling single parameters was absolutely correct, my ultimate goal was **cascading**: by passing variadic parameters, return only one specified `time.Date` from this atomic method, without creating any intermediate `Time` objects.

For example, I wanted a method like this. No matter how many parameters are cascaded, it should create only one `Time` object and return:

```go
GoMonth(1, 5, 3) // January 5th, 3 AM
```

I thought for a long time. If I continued changing each `case` implementation by passing all variadic parameters, I would inevitably return to that terrifying `if-else` hell. What should I do?

Now I had a single, stable time operation method. The question was: how to reuse it?

I stared blankly at the `GoMonth(1, 5, 3)` method on the screen for a long time. Suddenly, as if possessed by a deity, I seemed to see these three parameters `1, 5, 3` fly out of the screen. **I got it!**

> [!IMPORTANT]
>
> If I can call `GoMonth(1, 5, 3)`, what's the difference from chain calling `GoMonth(1).GoDay(5).GoHour(3)`?
>
> And if I can chain call `GoMonth(1).GoDay(5).GoHour(3)`, what's the difference from, in a loop, passing each call's result back through `applyAbs` to the next loop's next time unit (such as Day)?

In an instant, I felt the entire time dimension collapse. It was no longer complex and trivial. I could link everything together in each atomic operation! That is to say, I might need to successively pass in a loop body:

```go
i=0: y, m, d = applyAbs(Month, 1, y, m, d)
i=1: y, m, d = applyAbs(Day, 5, y, m, d)
i=2: y, m, d = applyAbs(Hour, 3, y, m, d)
```

This way, if I pass each cascade parameter downward, I can finally obtain the desired time!

*(Methods have been simplified)*

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

// cascade: core engine for time cascading
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

    p, pN := u, args[0] // Parent unit and its passed value
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

In this instant, I realized Aeon's soul was born. The entire process involved not a single memory allocation! All numerical calculations! This is extremely efficient in CPU operations.

---

Based on this principle, I implemented 4 methods: `Go`, `By`, `Start`, `End`. But I found that if I only had these, it still wasn't enough!

For example, in the current navigation system, I cannot **first offset** to a time point, then position on top of that; conversely, I also cannot **first position** to a time point, then offset on top of that.

To achieve these two capabilities, I needed to add two more actions to Aeon:

1. `At`: First position, then offset. For example: `At(5, 1, 1)` ➜ position to May, then add 1 day and 1 hour.
2. `In`: First offset, then position. For example: `In(1, 5, 1)` ➜ May 1st of next year.

It is these 4 completely orthogonal actions that form the core of Aeon's entire navigation system. I believe they encompass almost all possible time landing points. How to implement them?

I still used the `cascade` function (code has been simplified):

```go
for i, n := range args {
    unit := seq[i]
    if i == 0 { // If it's the first parameter, use positioning or offset method.
        y, m, d, h, mm, s, ns, w = applyRel(unit, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
    } else {
        y, m, d, h, mm, s, ns, w = applyAbs(unit, p, n, pN, y, m, d, h, mm, s, ns, w, sw)
    }
    p, pN = unit, n
}
```

And just like that, it was done. I thought it would be very complex... but the solution's elegance surprised even myself. This is the power of **orthogonality**!

*Note: `applyRel` principle is similar to `applyAbs`, but it uses **full relative offset** rather than direct positioning (setting).*

### Container System: Time Indexing

Through the above 4 actions, I possessed an almost complete time navigation system. But I wanted to go one step further.

What was my goal? I no longer wanted to see time as linear — like `AddDay().GoMonth().Start()`. Each call just does addition or subtraction on top of the previous time. It's essentially still linear calculation. I wanted **a unified hierarchical time structure**.

---

When we call `Go/ByMonth()`, we naturally think of it as setting the month of the current year. The same applies to `XXDay/Hour..()`. In our intuition, we think `Month` should belong to `Year`, `Day` belongs to `Month`, and `Hour` runs within `Day`.

Yes, this is the time view that conforms to **human intuition**. They are hierarchies nested within hierarchies — a container where large units contain small units, and small units contain sub-units. But what about `Year`? What's `Year`'s superior? Many time libraries stop here. They don't define a belonging for years — years are isolated.

**But in Aeon, I thoroughly implemented this concept throughout the entire navigation system!**

---

I defined a set of time unit methods for each of the above 4 actions, from century to nanosecond. For example, here's `Go`'s:

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

Just like that, simply choose the time unit you want, and it will automatically generate the **cascade sequence** from that unit to nanosecond. Each parameter will penetrate to subordinate levels like flowing water.

```go
GoMonth() // Month, Day, Hour, .., Nanosecond
GoCentury() // Century, Decade, Year, Month, Day, Hour, .., Nanosecond
GoDecade(2, 5) // 2nd decade, 5th year of this century = 2025
GoYear(2) // 2nd year of this decade = 2022
```

This is the **time container**. You no longer need to memorize every time navigation method. You only need to remember 4 actions plus the time unit you want to operate.

And in the container, reverse indexing is also possible — positioning from the end of the parent container.

```go
GoDay(-2, 23) // Last 2 days of this month, 23rd hour
GoMonth(-1, -3) // Last month of this year, 3rd day from end
```

This greatly reduces the user's mental burden. They no longer need to worry about calculating various boundaries.

**Time container indexing model:**

```text
[Millennium]
  └─ [0...9 Century]
       └─ [0...9 Decade]
            └─ [0...9 Year]
                 └─ [1...12 Month]

Example: GoYear(5) Indexing logic
         [-9]       [-8]            [-5]        [-4]             [-1]
2020 ─┬─ 2021 ──┬── 2022 ··· ──┬── [2025] ──┬── 2026 ─┬─ ··· ─┬─ 2029
[0]      [1]        [2]             [5]         [6]              [9]
```

---

Is it finished?

Yes. But to not completely lock time in the parent container, I provided 6 top-level methods that allow the first passed parameter to go to an **absolute year**:

```go
Go(2025, 1).StartDay(-1, 23) // 2025-01-31 23:00:00
```

Although Aeon's positioning core is "addressing," it also provides an `Overflow` flag that allows time to naturally overflow.

```go
GoMonth(aeon.Overflow, 2) // If overflow, could be 03-2/3.
```

---

> This is Aeon's story. It's not just a library — it's my reimagining of time logic.
>
> If you've read this far, I am deeply grateful for accompanying me on this journey. Thank you from the bottom of my heart!

Finally, beyond just **navigation**, in [Aeon](https://github.com/baagod/aeon)'s world, there are more unique perspectives on time.
