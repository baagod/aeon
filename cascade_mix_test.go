package aeon

import (
	"testing"
)

func TestMixSeries(t *testing.T) {
	base := Parse("2024-04-15 12:00:00")

	t.Run("StartAt/EndAt 系列", func(t *testing.T) {
		// StartAtYear(5, 1): 定位到本年代第 5 年 (2025)，然后加 1 个月，并归零
		assert(t, base.StartAtYear(5, 1), "2025-05-01 00:00:00", "StartAtYear(5, 1)")

		// EndAtMonth(6, 5): 定位到 6 月，然后加 5 天，并置满
		assert(t, base.EndAtMonth(6, 5), "2024-06-20 23:59:59.999999999", "EndAtMonth(6, 5)")
	})

	t.Run("High Fidelity At/In 系列", func(t *testing.T) {
		// AtYear(2025, 1): 绝对定位到 2025 年，加 1 个月，保留 12:00:00
		assert(t, base.At(2025, 1), "2025-05-15 12:00:00", "At(2025, 1) 保真")
	})
}
