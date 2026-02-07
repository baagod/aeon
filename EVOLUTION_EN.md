# Aeon: A Zero-Alloc Go Time Library That Reimagines Time as a "Container" Rather Than an "Offset"

> I can’t recall exactly when it began. Perhaps it originated from a simmering frustration with `time.Time` and the landscape of existing libraries — but eventually, it crystallized into a nearly absurd obsession: **Why not write my own Go time library?**
>
> In ancient philosophy, Aeon represents "eternity" and "layered dimensions."
>
> I chose this name to express a different fundamental logic: time is not a slender, linear thread. It is a flowing universe, capable of being nested, layered, and penetrated.

## Why Re-invent the Wheel?

The dilemma of existing solutions: a constant struggle between "linear arithmetic" and "heap allocation." Go’s standard library, `time.Time`, is an engineering marvel — precise, stable, and thread-safe. But the moment we attempt to wield it for complex business logic, the friction becomes unbearable.

### Cognitive Mismatch

To the human mind, time is inherently **hierarchical**. We think in terms of "the third Friday of next month" or "the last day of this month." Yet, in the logic of `time.Time`, time is **linear**: a mere accumulation of nanoseconds.

This creates a severe cognitive mismatch. Imagine trying to locate "the last n days of the next quarter." Using the standard library requires a series of "mental gymnastics":

*   Which months constitute the next quarter?
*   How many days are in those months? Is it a leap year?
*   Will `AddDate(0, 3, 0)` overshoot because I started on the 31st?

Code quickly degenerates into a linear algebra problem riddled with magic numbers `(0, 1, -1)`, obscuring the business logic it was meant to express.

### Memory Allocation

To solve usability hurdles, the community has produced excellent wrappers like `now` and `carbon`. They provide fluent, chainable APIs that are a joy to read. However, I could not tolerate their underlying cost: every single call (or link in a chain) triggers a heap allocation.

```go
// The nightmare of most wrapper libraries:
// New()      -> Alloc
// AddMonth() -> Alloc
// StartOf()  -> Alloc
Carbon.Now().AddMonth(1).StartOfWeek() // 3 heap allocations!
```

In a high-throughput, concurrent system, these fragmented GC pressures are unforgivable!

### Feature Bloat

I took another look at the `carbon` library; it is too "heavy." By "heavy," I don't mean it supports too many features, but that it fails to systematically abstract and unify those highly similar behaviors.

```go
IsSameYear(t)
IsSameMonth(t)
IsSameDay(t)
// ...and hundreds more...

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

I don't want to memorize 300 method names. That is **"brute-force enumeration"**; it is **"patchwork programming."** I need a divine sword — a single, precise instrument capable of dissecting time and severing the roots of chaos.

### The Genesis of Aeon

Imagine if the API were defined by intent rather than enumeration:

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
GoMonth(aeon.Ord, -1, 5)    // Last Friday of the month
StartWeekday(5, 18)         // This Friday at 18:00 (Happy Hour)
```

This was my breaking point. I realized I didn't just want a cleaner API; I wanted the extreme performance of **Zero-Alloc**. I wanted to leap across the timeline with the precision of a pointer, leaving no trace for the garbage collector to find.

Thus, Aeon was born.

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
> The metrics above were recorded during single atomic operations **without** variadic parameters. **Crucially, even with complex chained calls, Aeon remains Zero-Alloc. The more complex the logic, the more astonishing the performance gap becomes.**

> If you're looking for a quick introduction, you can stop here. Visit the [Aeon repository](https://github.com/baagod/aeon) or explore the [full documentation](https://zread.ai/baagod/aeon/1-overview) to dive deeper. Thank you for your support!

However, if you wish to see how Aeon evolved — from conception to birth — please read on.

## From a Simple "Overflow" Fix to the "Time Container" Model

In the beginning, I had no grand theories about cascade indexing or time containers. I certainly wasn't obsessing over zero-allocation performance. I had only one modest wish: **Let me manipulate the months without the days overflowing.**

### The Long Journey of Exploration

I started by creating `thru`, the predecessor to Aeon. I wanted more than just "adding"; I needed to "set" time directly, which led to the first `Go` method prototypes. For instance, `GoMonth(1, 2)` would set the date to January 2nd while preserving the year and clock time — and most importantly, it would suppress the month overflow.

You can see how much tedious labor went into this one "no overflow" fix. On Stack Overflow, this remains a perennial grievance: "Why did adding one month land me in the month after next?"

**But the nightmare was only just beginning.**

As I tried to scale this "patch-style" logic to handle weeks, quarters, years, and even centuries, the code spiraled into chaos.

I fell into an `if-else` hell. To keep the dates correct, I had to juggle leap years, varying month lengths, and boundary conditions for quarter-ends and cross-year weeks. Every time I patched the "month" hole, the "quarter" logic would buckle.

The architecture was in tatters. But at that moment, I didn't realize I was standing on the threshold of a deeper truth.

### The Power of Atomic Operations

So, I stopped. I abandoned the attempt to pre-calculate everything for every possible variadic parameter. I decided to do exactly one thing: **focus on the `Start` method and mandate that it accept only one parameter.**

> [!IMPORTANT]
>
> **I needed to prove that if the logic failed even at a single-dimension atomic level, then my fundamental arithmetic was flawed.** *(This realization became the cornerstone of Aeon's navigation system.)*

I defined `Start` as `t.Start(u Unit, ...n)`. If I wanted the start of a month, I’d call `Start(aeon.Month, 5)`. The intent was pure: pinpoint May and then "flatten" all subordinate units — day, hour, minute, second, and nanosecond.

In this minimalist model, the `if-else` clutter vanished. I could focus entirely on the logic of each individual unit using a clean `switch-case`:

```go
func applyAbs(u Unit, y, m, d int) Time {
    switch u {
    case Year:  // Handle year positioning
    case Month: // Handle month positioning
        if n > 0 {
            m = n
        } else if n < 0 {
            m = 13 + n // Negative indexing (reverse)
        }
    case ..
    }
}
```

If `0` is passed, we stay in the current unit. But how to handle carries? What if someone passes `m=13`?

To solve this, I designed an **automatic carry protocol**. Even if a nonsensical month like `13` or `-1` is provided, the value flows like water, automatically overflowing into the year until it settles into the correct tick mark.

```go
// addMonth calculates the year and month after an offset (handling carries/borrows)
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

This logic is applied immediately after the `switch`, ensuring the `time.Date()` parameters are always valid.

---

The final hurdle was the "month overflow" itself. I addressed this with a dedicated "Max Days" method:

```go
// DaysIn returns the maximum days in month m of year y.
func DaysIn(y int, m ...int) int {
    if len(m) > 0 {
        if m[0] == 2 && IsLeapYear(y) {
            return 29
        }
        return maxDays[m[0]]
    }
    return IsLeapYear(y) ? 366 : 365
}
```

By applying this at the end of `applyAbs`, I ensured that the day count would always be clamped to the month's limit:

```go
// Unified overflow check
if u <= Month {
    if dd := DaysIn(y, m); d > dd {
        d = dd
    }
}
```

With that, the "month overflow" nightmare was finally over.

---

But one question remained: how do I zero out the subordinate components?

When I call `Start(Month, 5)`, I need to reset everything from the "day" down to the "nanosecond" (e.g., `y-05-01 00:00:00.000...`). I solved this by unifying the boundary handling at the very end of the `switch`:

```go
// align handles the final alignment (zeroing or filling)
func align(u Unit, y, m, d, h, mm, sec, ns int) (int, int, int, int, int, int, int) {
    switch u {
    case Century, Decade, Year:
        m, d, h, mm, sec, ns = 1, 1, 0, 0, 0, 0
    case Quarter, Month:
        d, h, mm, sec, ns = 1, 0, 0, 0, 0
    case Week, Weekday, Day:
        h, mm, sec, ns = 0, 0, 0, 0
    // ...etc...
    }
    return y, m, d, h, mm, sec, ns
}
```

By chaining these atomized methods together, I created an **atomic time operation engine**. It was a stable, precise instrument — a sword ready to be drawn. But what could it truly do?

### Dimensional Collapse: The Cascade Architecture

Atomic operations were correct, but my ultimate goal was **Cascading**: passing variadic parameters to return a specific `time.Date` in one go, without creating intermediate `Time` objects.

I wanted this:

```go
GoMonth(1, 5, 3) // January 5th, 03:00
```

I stared at that method signature on the screen for a long time. Suddenly, as if struck by divine inspiration, I saw the numbers `1`, `5`, and `3` **fly out of the screen**, detaching themselves from the line. **I got it!**

> [!IMPORTANT]
>
> If I can call `GoMonth(1, 5, 3)`, how is that fundamentally different from chaining `GoMonth(1).GoDay(5).GoHour(3)`?
>
> And if those are equivalent, how is a chain different from a simple loop that takes the result of one atomic operation and feeds it into the next unit?

In that instant, **the dimensions of time collapsed**. The complexity vanished. I realized I could string everything together within each atomic operation! I just needed to pass the baton in a sequence:

```go
i=0: y, m, d = applyAbs(Month, 1, y, m, d)
i=1: y, m, d = applyAbs(Day, 5, y, m, d)
i=2: y, m, d = applyAbs(Hour, 3, y, m, d)
```

By passing each cascading parameter downward, I could calculate the final time in a single pass.

```go
// cascade: The core engine of Aeon
func cascade(t Time, u Unit, args ...int) Time {
    // ... setup ...
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

Aeon's soul was born. No memory allocations, just raw numerical calculation. Peak efficiency.

---

Based on this principle, I constructed the **four fundamental actions**:

1.  **`Go`**: Precise positioning.
2.  **`By`**: Relative offset.
3.  **`Start`**: Boundary zeroing.
4.  **`End`**: Boundary filling.

But reality is often mixed. Sometimes you need to offset *then* position, or vice versa. To handle these cases, I added two compound actions:

*   **`At`** (Position then Offset): e.g., `At(5, 1)` ➜ May, then +1 day.
*   **`In`** (Offset then Position): e.g., `In(1, 5)` ➜ Next year, then May.

These **4 + 2** actions form an orthogonal universe, covering nearly every conceivable point in time. The implementation was shockingly elegant—the power of **orthogonality** at work.

### The Container System: Intuitive Time Indexing

With these actions, I had a complete navigation system. But I wanted more. I didn't want time to feel linear anymore. I wanted a **unified hierarchy**.

---

When we call `GoMonth()`, it feels natural that we are setting the month *within* the current year. Similarly, a day belongs to a month, and an hour belongs to a day.

This is the time view that aligns with **human intuition**. Units are nested like containers. Yet, many libraries treat the `Year` as a top-level, isolated unit. They view the year as infinite and unattached.

**In Aeon, I implemented this hierarchy to its logical conclusion.**

---

I defined a full set of unit methods for all actions, from Centuries to Nanoseconds.

```go
GoCentury()
GoDecade()
GoYear()
GoMonth()
// ...
```

Simply pick your unit, and Aeon generates a **cascade sequence** down to the nanosecond. Parameters flow like water through the layers.

```go
GoMonth()      // Month -> Day -> Hour...
GoDecade(2, 5) // 2nd decade, 5th year = 2025
GoYear(2) // 2nd year of this decade = 2022
```

This is the **Time Container**. You don't need to memorize hundreds of methods; you only need to know 4 actions and the unit you want to target. You can even index backward from the end of a container:

```go
GoDay(-2, 23)   // 2nd to last day of the month, 23:00
GoMonth(-1, -3) // Last month of the year, 3rd to last day
```

This erases the mental overhead of boundary calculations.

**The Time Container Indexing Model:**

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

Is that all?

Almost. To ensure flexibility, I provided 6 top-level methods that allow the first parameter to point to an **absolute year**, and added an `Overflow` flag for those who prefer natural time spillover.

## Epilogue

> This is the story of Aeon. May time flow eternally within your code, uninterrupted.

If you have made it this far, thank you for witnessing the evolution of Aeon. I truly appreciate your journey with me.

Finally, beyond **navigation**, the world of [Aeon](https://github.com/baagod/aeon) holds many more unique perspectives on the nature of time.
