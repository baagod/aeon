> I can't recall the exact moment it began. Perhaps it stemmed from a dissatisfaction with `time.Time` and existing time libraries—I developed a near-absurd obsession: **Why not write a Go time library myself?**
>
> In ancient philosophy, Aeon represents "eternity" and "layered dimensions."
>
> I chose this name because I want to express a different logic of time—time is not a slender straight line; it is a flowing universe that can be nested and penetrated.

# Aeon

[![zread](https://img.shields.io/badge/Ask_Zread-_.svg?style=flat&color=00b0aa&labelColor=000000&logo=data%3Aimage%2Fsvg%2Bxml%3Bbase64%2CPHN2ZyB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIHZpZXdCb3g9IjAgMCAxNiAxNiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTQuOTYxNTYgMS42MDAxSDIuMjQxNTZDMS44ODgxIDEuNjAwMSAxLjYwMTU2IDEuODg2NjQgMS42MDE1NiAyLjI0MDFWNC45NjAxQzEuNjAxNTYgNS4zMTM1NiAxLjg4ODEgNS42MDAxIDIuMjQxNTYgNS42MDAxSDQuOTYxNTZDNS4zMTUwMiA1LjYwMDEgNS42MDE1NiA1LjMxMzU2IDUuNjAxNTYgNC45NjAxVjIuMjQwMUM1LjYwMTU2IDEuODg2NjQgNS4zMTUwMiAxLjYwMDEgNC45NjE1NiAxLjYwMDFaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik00Ljk2MTU2IDEwLjM5OTlIMi4yNDE1NkMxLjg4ODEgMTAuMzk5OSAxLjYwMTU2IDEwLjY4NjQgMS42MDE1NiAxMS4wMzk5VjEzLjc1OTlDMS42MDE1NiAxNC4xMTM0IDEuODg4MSAxNC4zOTk5IDIuMjQxNTYgMTQuMzk5OUg0Ljk2MTU2QzUuMzE1MDIgMTQuMzk5OSA1LjYwMTU2IDE0LjExMzQgNS42MDE1NiAxMy43NTk5VjExLjAzOTlDNS42MDE1NiAxMC42ODY0IDUuMzE1MDIgMTAuMzk5OSA0Ljk2MTU2IDEwLjM5OTlaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik0xMy43NTg0IDEuNjAwMUgxMS4wMzg0QzEwLjY4NSAxLjYwMDEgMTAuMzk4NCAxLjg4NjY0IDEwLjM5ODQgMi4yNDAxVjQuOTYwMUMxMC4zOTg0IDUuMzEzNTYgMTAuNjg1IDUuNjAwMSAxMS4wMzg0IDUuNjAwMUgxMy43NTg0QzE0LjExMTkgNS42MDAxIDE0LjM5ODQgNS4zMTM1NiAxNC4zOTg0IDQuOTYwMVYyLjI0MDFDMTQuMzk4NCAxLjg4NjY0IDE0LjExMTkgMS42MDAxIDEzLjc1ODQgMS42MDAxWiIgZmlsbD0iI2ZmZiIvPgo8cGF0aCBkPSJNNCAxMkwxMiA0TDQgMTJaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik00IDEyTDEyIDQiIHN0cm9rZT0iI2ZmZiIgc3Ryb2tlLXdpZHRoPSIxLjUiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIvPgo8L3N2Zz4K&logoColor=ffffff)](https://zread.ai/baagod/aeon)

🇨🇳 [中文](README_CN.md) | 🇺🇸 [English](README.md)

> Aeon is a **Zero-Allocation** time navigation library for Go based on **Time Container** theory. It replaces linear calculation with structured navigation, expressing complex time intentions in a way that aligns with human intuition.

## 🚀 Performance: A Dimensional Blow

Aeon achieves **True Zero Allocation** (Zero Alloc) and leverages a **Cascading Addressing** architecture. Whether you span multiple dimensions (from Millennium to Nanosecond), Aeon completes the operation in a **single atomic step**. The more complex the logic, the more staggering Aeon's lead becomes.

> [!NOTE]
> The following baseline data were obtained under single-atom operations without **using** cascade parameters.

```bash
Benchmark       | ns/op | allocs/op x B/op | up

New             |
Aeon            | 18.6 | 0 | x74
Carbon          | ████████████████████████████████████████ 1376 | 13x1600

Now             |
Aeon            | 7.8 | 0 | x177
Carbon          | ████████████████████████████████████████ 1384 | 13x1600

From Unix       |
Aeon            | 3.6 | 0 | x383
Carbon          | ████████████████████████████████████████ 1380 | 13x1600

From Std        |
Aeon            | 5.0 | 0 | x323
Carbon          | ████████████████████████████████████████ 1619 | 13x1600

Parse (Compact) |
Aeon            | 23.3 | 0 | x195
Carbon          | ████████████████████████████████████████ 4561 | 85x3922

Parse (ISO)     |
Aeon            | 19.6 | 0 | x91
Carbon          | ████████████████████████████████████████ 1794 | 15x1697

Start/End       |
Aeon            | █ 56.4 | 0 | x20
Carbon          | ████████████████████ 1141 | 7x1440

Add (Offset)    |
Aeon            | █ 56.5 | 0 | x2.5
Carbon          | ██ 142 | 2x128

Set (Position)  |
Aeon            | █ 58.7 | 0 | x2.6
Carbon          | ███ 156 | 2x128
```

## 📦 Installation

```bash
go get github.com/baagod/aeon
```

## 🧊 Core Concept: Containers

The core of Aeon is **Container Offset**. All **navigation** is essentially indexing within the **Parent Container** of the current unit (starting from `0`). For example:

- **`GoYear(5)`**: Not going to the year 5 AD, but indexing to the **5th year** within the **current Decade** (the parent container) ➜ `···5`.
- **`GoDecade(2)`**: Indexing to the **2nd Decade** of the **current Century** ➜ `··2·`.
- **`GoCentury(0)`**: Indexing to the **0th Century** of the **current Millennium** ➜ `·0··`.

```text
[Millennium]
  └─ [0...9 Century]
       └─ [0...9 Decade]
            └─ [0...9 Year]
                 └─ [1...12 Month]

Example: GoYear(5) Addressing Logic
         [-9]       [-8]            [-5]        [-4]             [-1]
2020 ─┬─ 2021 ──┬── 2022 ··· ──┬── [2025] ──┬── 2026 ─┬─ ··· ─┬─ 2029
[0]      [1]        [2]             [5]         [6]              [9]
```

## 🧭 Navigation Matrix

Aeon's API design is completely **Orthogonal**. You only need to remember **4 Actions**:

- `Go.. [·]` **Absolute Positioning:** `GoYear(5, 1)` ➜ 5th Year, 1st Month of current decade.
- `By.. [➜]` **Relative Offset:** `ByYear(1, 5)` ➜ Offset by 1 Year and 5 Months.
- `At.. [·, ➜]` **Position then Offset:** `AtYear(5, 1)` ➜ Locate 5th Year, then offset 1 Month.
- `In.. [➜, ·]` **Offset then Position:** `InYear(1, 5)` ➜ Next Year (Offset 1), then 5th Month.

> [!IMPORTANT]
>
> 1. `By` methods default to `1`. Others default to `0`.
> 2. Invalid `0` time (e.g., 0th Month) remains unchanged in **Positioning Mode** (but works in Offset mode).
> 3. `Go` positions only the target unit and preserves original time details as much as possible** (e.g., `GoWeek` automatically retains the weekday).
>   ```go 
>   t := Parse("2021-07-21 07:14:15") // Wed
>   t.GoMonth(1)  // 2021-01-21 07:14:15 (Set to Jan, time preserved)
>   t.GoWeek(1)   // 2021-06-30 07:14:15 (1st Week, Weekday preserved as Wed)
>   ```

---

Combined with `Start/End` prefixes to hit time boundaries:

- `StartYear()`: Start of this year (01-01 00:00:00...)
- `EndYear()`: End of this year (12-31 23:59:59...)

---

6 **Top-Level** methods allow the **first parameter** to enter **Absolute Year** mode:

1. `Go(2025, 2)` ➜ 2025-02
2. `At(2025, 2)` ➜ Position at 2025, then offset 2 months.
3. `Start(2025, 2)` ➜ 2025-02-01 00:00:00
4. `StartAt(2025, 1)` ➜ Position at 2025, offset 1 month, then Start of Month.
5. `End(2025, 2)` ➜ 2025-02-28 23:59:59...
6. `EndAt(2025, 1)` ➜ Position at 2025, offset 1 month, then End of Month.

---

### ♾️ Cascading Parameters

Method chaining? No, this is **Atomic Operation**! All methods support **Variadic Parameters** that cascade downwards. Parameters flow like water, completing complex positioning in **one line of code**.

Aeon automatically switches between 4 cascading sequences based on the **<u>Entry Unit</u>**:

1. **Year Sequence `Default`**: `Century ➜ Decade ➜ Year ➜ Month ➜ Day ➜ Hour.. ➜ Nanosecond`
2. **Quarter Flow `Quarter`**: `Quarter ➜ Month (in Quarter) ➜ Day ➜ Hour.. ➜ Nanosecond`
3. **Week Sequence `Week`🦬**: `Week (Smart Context) ➜ Weekday ➜ Hour.. ➜ Nanosecond`

   This is a **Transformer**! It automatically shifts shape based on the passed **Flags**:

   - `ISO`: ISO Week. Starts from the 1st ISO week of the year.
   - `Full`: Full Week. Starts from the 1st Monday of the month.
   - `Ord`: Ordinal Week. Starts from the 1st day of the month.
   - `Qtr`: Quarter Week. Starts from the 1st day of the quarter's first month.
   - `Default`: Calendar/Natural Week. Follows the calendar visual row.

4. **Weekday Flow `Weekday`**: `Weekday ➜ Hour.. ➜ Nanosecond`

```go
// Relative offset: 1 Year, 3 Months, 5 Days
ByYear(1, 3, 5)

// 2nd Tuesday of the current Quarter
GoWeek(aeon.Qtr|aeon.Ord, 2, 2)

// Last Friday of the current Quarter
GoWeek(aeon.Qtr|aeon.Ord, -1, 5)

// 2025, Feb, Last Day, 23:00
Go(2025).StartMonth(2, -1, 23)

// End of the 3rd Quarter, minus 1 month, minus 2 days
EndQuarter(3, -1, -2)

// 10th ISO Monday of 2025
Go(2025).StartWeek(aeon.ISO, 10, 1)

// 3rd Friday of this month (Ordinal week starting from 1st)
StartWeek(aeon.Ord, 3, 5)

// Last Friday of this month
GoWeek(aeon.Ord, -1, 5)

// End of previous Quarter
EndByQuarter(-1)

// 1st day of the last month of this Quarter
StartQuarter(0, -1, 1)

// This Friday at 18:00 (Happy Hour)
StartWeekday(5, 18)

// 3rd to last day of this month
StartDay(-3)

// Next Wednesday at 2 PM
StartInWeek(1, 3, 14)

// Yearly Archive: Start/End boundaries
StartYear() / EndYear()

// Last day of next month
EndInMonth(1, -1)
```

*Negative numbers are not just subtraction; they are **Reverse Indexing**, representing the **"N-th from last"** item in the container.*

### 🛡️ Overflow Protection

Aeon's core philosophy is **Intention First**. By default, navigation protects against day overflow for units **"Month and above"**.

```go
base := NewDate(2025, 1, 31)
base.GoMonth(2) // 2025-02-28 (Protected)
base.ByMonth(Overflow, 1) // 2025-03-03 (Overflow allowed)
base.ByMonth(1, 2) // 🛡️🦬 2025-03-02 (Protect to 2-28, then add 2 days)

// Leap Year Handling
leap := NewDate(2024, 2, 29)
leap.ByYear(1) // 2025-02-28 (Protected)
leap.ByYear(Overflow, 1) // 2025-03-01 (Overflow: Crosses month boundary)
leap.ByYear(4)           // 2028-02-29 (Next Leap Year)
```

## 🧰 Core API

### Create Time

```go
// Create time from optional time.Time(s)
Aeon(t ...time.Time) Time

// Create current time with optional location
Now(loc ...*time.Location) Time

// Precisely create time (y, m, d, h, mm, s), with optional nanosecond/location
New(y, m, d, h, mm, int, add ...any) Time

// Create from timestamp (auto-detects second/milli/micro/nano)
Unix(int64, ...*time.Location) Time

// Parse time string (ignores error)
Parse(string, ...*time.Location) Time

// Parse time string (returns error)
ParseE(string, ...*time.Location) (Time, error)

// Parse with layout (ignores error)
ParseBy(layout, value string, loc ...*time.Location) Time

// Parse with layout (returns error)
ParseByE(layout, value string, loc ...*time.Location) (Time, error)
```

### Get Time

```go
t.Year() int                       // Year
t.Month() int                      // Month (1-12)
t.Day() int                        // Day (1-31)
t.Hour() int                       // Hour (0-23)
t.Minute() int                     // Minute (0-59)
t.Second() int                     // Second (0-59)
t.Milli() int                      // Millisecond (0-999)
t.Micro() int                      // Microsecond (0-999999)
t.Nano() int                       // Nanosecond (0-999999999)
t.Clock() (h, mm, s int)           // Hour, Minute, Second
t.Weekday() time.Weekday           // Weekday
t.ISOWeek() (year, week int)       // ISO Week
t.YearDay() int                    // Day of Year (1-365/366)
t.YearDays() int                   // Total Days in Year
t.Days() int                       // Total Days in Month
t.Unix(n ...int) int64             // Timestamp, n specifies precision
t.Time() time.Time                 // Convert to Std Lib
t.Location() *time.Location        // Location
t.Zone() (name string, offset int) // Zone Name and Offset (Seconds)
```

### Compare & Judge

```go
t.Lt(u Time) bool          // t < u
t.Gt(u Time) bool          // t > u
t.Eq(u Time) bool          // t == u
t.Compare(u Time) int      // -1 / 0 / 1

t.IsSame(u Unit, target Time) bool        // Check if same unit (Year/Month/Day etc.)
t.Between(start, end Time, bound ...byte) bool  // Check if within interval
// bound: '=' inclusive, '!' exclusive, '[' left inclusive, ']' right inclusive

t.Diff(u Time, unit string, abs ...bool) float64
// unit: "y"Year, "M"Month, "d"Day, "h"Hour, "m"Minute, "s"Second
// abs: true returns absolute value

t.Sub(u Time) time.Duration  // t - u
t.Until() time.Duration      // t - Now()
```

### Extremum Selection

```go
Pick(op byte, times ...Time) Time
// op: '>' Max, '<' Min, '-' Near, '+' Far
```

### Zone Conversion

```go
t.UTC() Time                         // UTC Time
t.Local() Time                       // Local Time
t.To(loc *time.Location) Time        // Specific Location
t.WithWeekStarts(w time.Weekday) Time // Set Week Start Day
```

### Utility Functions

```go
IsLeapYear(y int) bool      // Is Leap Year
IsLongYear(y int) bool      // Is ISO Long Year (53 Weeks)
DaysIn(y int, m ...int) int // Days in y Year m Month, returns year days if m omitted

t.IsLeapYear() bool   // Is Leap Year
t.IsLongYear() bool   // Is ISO Long Year
t.IsWeekend() bool    // Is Weekend
t.IsAM() bool         // Is AM
t.IsDST() bool        // Is DST
t.IsZero() bool       // Is Zero
t.ZeroOr(u Time) Time // Return u if Zero
```

### Formatting

```go
t.String() string                 // Default Format
t.Format(layout string) string    // Custom Format
t.ToString(f ...string) string    // Custom Format (Short)
t.AppendFormat(b []byte, layout string) []byte
```

### Time Zones

```go
// Constant Zones (70+ IANA Zones)
const (
    UTC, Local, CET, EET, EST, GMT, MET, MST, WET
    Shanghai, Tokyo, NewYork, London, Paris, Berlin, Moscow
    // ...
)

// Return Zone by name and offset
Zone(string, ...int) *time.Location
```