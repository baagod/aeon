package aeon

import (
    "testing"
    "time"

    "github.com/dromara/carbon/v2"
    "github.com/relvacode/iso8601"
)

// --- 创建 ---

func Benchmark_New(b *testing.B) {
    b.Run("Aeon/Full", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = New(2020, 8, 5, 13, 14, 15, 999999999, UTC)
        }
    })

    b.Run("Carbon/Full", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = carbon.CreateFromDateTimeNano(2020, 8, 5, 13, 14, 15, 999999999, UTC)
        }
    })

    b.Run("Aeon/Now", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = Now(time.UTC)
        }
    })

    b.Run("Carbon/Now", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = carbon.Now(UTC)
        }
    })

    b.Run("Aeon/Unix", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = Unix(1596633255, time.UTC)
        }
    })

    b.Run("Carbon/Unix", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = carbon.CreateFromTimestamp(1596633255, UTC)
        }
    })

    b.Run("Aeon/Std", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = Aeon(time.Now().In(time.UTC))
        }
    })

    b.Run("Carbon/Std", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = carbon.CreateFromStdTime(time.Now(), UTC)
        }
    })

    b.Run("Aeon/ToStd", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = Now(time.UTC).Time()
        }
    })

    b.Run("Carbon/ToStd", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = carbon.Now(UTC).StdTime()
        }
    })
}

// --- 偏移 (触发溢出) ---

func Benchmark_Add_Overflow(b *testing.B) {
    t := New(2025, 1, 31, 0, 0, 0)
    c := carbon.CreateFromDate(2025, 1, 31)

    b.Run("Aeon/Century", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByCentury(Overflow, 1)
        }
    })

    b.Run("Carbon/Century", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddCenturies(1)
        }
    })

    b.Run("Aeon/Decade", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByDecade(Overflow, 1)
        }
    })

    b.Run("Carbon/Decade", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddDecades(1)
        }
    })

    b.Run("Aeon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByYear(Overflow, 1)
        }
    })

    b.Run("Carbon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddYears(1)
        }
    })

    b.Run("Aeon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByMonth(Overflow, 1)
        }
    })

    b.Run("Carbon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddMonths(1)
        }
    })

    b.Run("Aeon/Week", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByWeek(1)
        }
    })

    b.Run("Carbon/Week", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddWeeks(1)
        }
    })

    b.Run("Aeon/Day", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByDay(1)
        }
    })

    b.Run("Carbon/Day", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddDays(1)
        }
    })

    b.Run("Aeon/Hour", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByHour(1)
        }
    })

    b.Run("Carbon/Hour", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddHours(1)
        }
    })

    b.Run("Aeon/Minute", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByMinute(1)
        }
    })

    b.Run("Carbon/Minute", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddMinutes(1)
        }
    })

    b.Run("Aeon/Second", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.BySecond(1)
        }
    })

    b.Run("Carbon/Second", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddSeconds(1)
        }
    })

    b.Run("Aeon/Milli", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByMilli(1)
        }
    })

    b.Run("Carbon/Milli", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddMilliseconds(1)
        }
    })

    b.Run("Aeon/Micro", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByMicro(1)
        }
    })

    b.Run("Carbon/Micro", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddMicroseconds(1)
        }
    })

    b.Run("Aeon/Nano", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByNano(1)
        }
    })

    b.Run("Carbon/Nano", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddNanoseconds(1)
        }
    })
}

// --- 偏移 (保护溢出) ---

func Benchmark_Add(b *testing.B) {
    t := New(2025, 1, 31, 0, 0, 0)
    c := carbon.CreateFromDate(2025, 1, 31)

    b.Run("Aeon/Century", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByCentury(1)
        }
    })

    b.Run("Carbon/Century", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddCenturyNoOverflow()
        }
    })

    b.Run("Aeon/Decade", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByDecade(1)
        }
    })

    b.Run("Carbon/Decade", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddDecadesNoOverflow(1)
        }
    })

    b.Run("Aeon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByYear(1)
        }
    })

    b.Run("Carbon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddYearNoOverflow()
        }
    })

    b.Run("Aeon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.ByMonth(1)
        }
    })

    b.Run("Carbon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.AddMonthNoOverflow()
        }
    })
}

// --- 设置 (触发溢出) ---

func Benchmark_Set_Overflow(b *testing.B) {
    t := New(2025, 1, 31, 0, 0, 0)
    c := carbon.CreateFromDate(2025, 1, 31)

    b.Run("Aeon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoYear(Overflow, 2024)
        }
    })

    b.Run("Carbon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetYear(2024)
        }
    })

    b.Run("Aeon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoMonth(Overflow, 2)
        }
    })

    b.Run("Carbon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetMonth(2)
        }
    })

    b.Run("Aeon/Day", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoDay(15)
        }
    })

    b.Run("Carbon/Day", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetDay(15)
        }
    })

    b.Run("Aeon/Hour", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoHour(14)
        }
    })

    b.Run("Carbon/Hour", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetHour(14)
        }
    })

    b.Run("Aeon/Minute", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoMinute(30)
        }
    })

    b.Run("Carbon/Minute", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetMinute(30)
        }
    })

    b.Run("Aeon/Second", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoSecond(45)
        }
    })

    b.Run("Carbon/Second", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetSecond(45)
        }
    })

    b.Run("Aeon/Milli", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoMilli(123)
        }
    })

    b.Run("Carbon/Milli", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetMillisecond(123)
        }
    })

    b.Run("Aeon/Micro", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoMicro(123456)
        }
    })

    b.Run("Carbon/Micro", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetMicrosecond(123456)
        }
    })

    b.Run("Aeon/Nano", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoNano(123456789)
        }
    })

    b.Run("Carbon/Nano", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetNanosecond(123456789)
        }
    })
}

// --- 设置 (保护溢出) ---

func Benchmark_Set(b *testing.B) {
    t := New(2025, 1, 31, 0, 0, 0)
    c := carbon.CreateFromDate(2025, 1, 31)

    b.Run("Aeon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoYear(2024)
        }
    })

    b.Run("Carbon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetYearNoOverflow(2024)
        }
    })

    b.Run("Aeon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.GoMonth(2)
        }
    })

    b.Run("Carbon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.SetMonthNoOverflow(2)
        }
    })
}

// --- Start/End ---

func Benchmark_Start(b *testing.B) {
    t := New(2025, 6, 15, 12, 30, 45, 0, UTC)
    c := carbon.CreateFromDateTimeNano(2025, 6, 15, 12, 30, 45, 0, UTC)

    b.Run("Aeon/Century", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.StartCentury()
        }
    })

    b.Run("Carbon/Century", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.StartOfCentury()
        }
    })

    b.Run("Aeon/Decade", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.StartDecade()
        }
    })

    b.Run("Carbon/Decade", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.StartOfDecade()
        }
    })

    b.Run("Aeon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.StartYear()
        }
    })

    b.Run("Carbon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.StartOfYear()
        }
    })

    b.Run("Aeon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.StartMonth()
        }
    })

    b.Run("Carbon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.StartOfMonth()
        }
    })

    b.Run("Aeon/Day", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.StartDay()
        }
    })

    b.Run("Carbon/Day", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.StartOfDay()
        }
    })

    b.Run("Aeon/Hour", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.StartHour()
        }
    })

    b.Run("Carbon/Hour", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.StartOfHour()
        }
    })

    b.Run("Aeon/Minute", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.StartMinute()
        }
    })

    b.Run("Carbon/Minute", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.StartOfMinute()
        }
    })

    b.Run("Aeon/Second", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.StartSecond()
        }
    })

    b.Run("Carbon/Second", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.StartOfSecond()
        }
    })
}

func Benchmark_End(b *testing.B) {
    t := New(2025, 6, 15, 12, 30, 45, 0, UTC)
    c := carbon.CreateFromDateTimeNano(2025, 6, 15, 12, 30, 45, 0, UTC)

    b.Run("Aeon/Century", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.EndCentury()
        }
    })

    b.Run("Carbon/Century", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.EndOfCentury()
        }
    })

    b.Run("Aeon/Decade", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.EndDecade()
        }
    })

    b.Run("Carbon/Decade", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.EndOfDecade()
        }
    })

    b.Run("Aeon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.EndYear()
        }
    })

    b.Run("Carbon/Year", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.EndOfYear()
        }
    })

    b.Run("Aeon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.EndMonth()
        }
    })

    b.Run("Carbon/Month", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.EndOfMonth()
        }
    })

    b.Run("Aeon/Day", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.EndDay()
        }
    })

    b.Run("Carbon/Day", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.EndOfDay()
        }
    })

    b.Run("Aeon/Hour", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.EndHour()
        }
    })

    b.Run("Carbon/Hour", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.EndOfHour()
        }
    })

    b.Run("Aeon/Minute", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.EndMinute()
        }
    })

    b.Run("Carbon/Minute", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.EndOfMinute()
        }
    })

    b.Run("Aeon/Second", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = t.EndSecond()
        }
    })

    b.Run("Carbon/Second", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = c.EndOfSecond()
        }
    })
}

// --- Parse ---

func Benchmark_Parse(b *testing.B) {
    // Date
    b.Run("Aeon/Date", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = Parse("2020-08-05")
        }
    })

    b.Run("Iso8601/Date", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _, _ = iso8601.ParseString("2020-08-05")
        }
    })

    b.Run("Aeon/DateShort", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = Parse("2020-8-5")
        }
    })

    b.Run("iso8601/DateShort", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _, _ = iso8601.ParseString("2020-8-5")
        }
    })

    // DateTime
    b.Run("Aeon/DateTime", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = Parse("2020-08-05 03:14:05")
        }
    })

    b.Run("Iso8601/DateTime", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _, _ = iso8601.ParseString("2020-08-05 03:14:05")
        }
    })

    b.Run("Aeon/DateTimeShort", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = Parse("2020-8-5 3:1:5")
        }
    })

    b.Run("iso8601/DateTimeShort", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _, _ = iso8601.ParseString("2020-8-5 3:1:5")
        }
    })

    b.Run("Aeon/DateTimeNanosecond", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = Parse("2020-08-05 03:14:05.123456789")
        }
    })

    b.Run("Iso8601/DateTimeNanosecond", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _, _ = iso8601.ParseString("2020-08-05 03:14:05.123456789")
        }
    })

    // Short
    b.Run("Aeon/Short", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = Parse("20200805131415")
        }
    })

    b.Run("Iso8601/Short", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _, _ = iso8601.ParseString("20200805131415")
        }
    })

    // ---

    b.Run("Aeon/ISO8601Zulu", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = Parse("2020-08-05T13:14:15Z")
        }
    })

    b.Run("Iso8601/ISO8601Zulu", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _, _ = iso8601.ParseString("2020-08-05T13:14:15Z")
        }
    })

    b.Run("Aeon/RFC3339", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = Parse("2020-08-05T13:14:15+08:00")
        }
    })

    b.Run("Iso8601/RFC3339", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _, _ = iso8601.ParseString("2020-08-05T13:14:15+08:00")
        }
    })
}
