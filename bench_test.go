package aeon

import (
    "testing"
    "time"
)

func Benchmark_Func(b *testing.B) {
    b.Run("Aeon/Date", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = Parse("2024-01-02")
        }
    })

    b.Run("Std/Date", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _, _ = time.Parse(time.DateOnly, "2024-01-02")
        }
    })

    b.Run("Aeon/DateTime", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = Parse("2024-01-02 11:12:13")
        }
    })

    b.Run("Std/DateTime", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _, _ = time.Parse(time.DateTime, "2024-01-02 11:12:13")
        }
    })

    b.Run("Aeon/DateTimeMilli", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = Parse("2024-01-02 11:12:13.999")
        }
    })

    b.Run("Std/DateTimeMilli", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _, _ = time.Parse(time.DateTime+".000", "2024-01-02 11:12:13.999")
        }
    })
}
