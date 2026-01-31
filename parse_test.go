package aeon

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	// 设置默认时区为 UTC 以便结果幂等
	oldLoc := DefaultTimeZone
	DefaultTimeZone = time.UTC
	defer func() { DefaultTimeZone = oldLoc }()

	t.Run("Date", func(t *testing.T) {
		assert(t, Parse("2024"), "2024-01-01 00:00:00", "2024")
		assert(t, Parse("2024-5"), "2024-05-01 00:00:00", "2024-5")
		assert(t, Parse("2024-05"), "2024-05-01 00:00:00", "2024-05")
		assert(t, Parse("2024-5-2"), "2024-05-02 00:00:00", "2024-5-2")
		assert(t, Parse("2024-5-02"), "2024-05-02 00:00:00", "2024-5-02")
		assert(t, Parse("2020-02-29"), "2020-02-29 00:00:00", "2020-02-29")
	})

	t.Run("Time", func(t *testing.T) {
		assert(t, Parse("2:3"), "0000-01-01 02:03:00", "2:3")
		assert(t, Parse("2:13"), "0000-01-01 02:13:00", "2:13")
		assert(t, Parse("2:3:4"), "0000-01-01 02:03:04", "2:3:4")
		assert(t, Parse("2:3:14"), "0000-01-01 02:03:14", "2:3:14")
		assert(t, Parse("12:13:14.999"), "0000-01-01 12:13:14.999", "12:13:14.999")
		assert(t, Parse("12:03:04"), "0000-01-01 12:03:04", "12:03:04")
		assert(t, Parse("13:14"), "0000-01-01 13:14:00", "13:14")
	})

	t.Run("Number", func(t *testing.T) {
		assert(t, Parse("202405"), "2024-05-01 00:00:00", "202405")
		assert(t, Parse("20241220"), "2024-12-20 00:00:00", "20241220")
		assert(t, Parse("2024052015"), "2024-05-20 15:00:00", "2024052015")
		assert(t, Parse("202405201514"), "2024-05-20 15:14:00", "202405201514")
		assert(t, Parse("20240520150415"), "2024-05-20 15:04:15", "20240520150415")
		assert(t, Parse("20240520150415Z"), "2024-05-20 15:04:15", "20240520150415Z")
		assert(t, Parse("20240520T150415"), "2024-05-20 15:04:15", "20240520T150415")
		assert(t, Parse("20240520T150415Z"), "2024-05-20 15:04:15", "20240520T150415Z")

		// 紧凑带纳秒 (Waterline Probing)
		assert(t, Parse("2024.1"), "2024-01-01 00:00:00.1", "2024.1")
		assert(t, Parse("2024.123"), "2024-01-01 00:00:00.123", "2024.123")
		assert(t, Parse("202405.123"), "2024-05-01 00:00:00.123", "202405.123")
		assert(t, Parse("20240520.123"), "2024-05-20 00:00:00.123", "20240520.123")
		assert(t, Parse("2024052015.123"), "2024-05-20 15:00:00.123", "2024052015.123")
		assert(t, Parse("202405201504.123"), "2024-05-20 15:04:00.123", "202405201504.123")
		assert(t, Parse("20240520150405.123"), "2024-05-20 15:04:05.123", "20240520150405.123")
		assert(t, Parse("20240520T150405.123"), "2024-05-20 15:04:05.123", "20240520T150405.123")
	})

	t.Run("DateTime", func(t *testing.T) {
		assert(t, Parse("2024-05-20 13:14:15"), "2024-05-20 13:14:15", "2024-05-20 13:14:15")
		assert(t, Parse("2024-05-20T13:14:15"), "2024-05-20 13:14:15", "2024-05-20 13:14:15")
		assert(t, Parse("2024-05-20T13:14:15.999"), "2024-05-20 13:14:15.999", "2024-05-20 13:14:15.999")
		assert(t, Parse("2004-4-5 3:4:1"), "2004-04-05 03:04:01", "2004-4-5 3:4:1")
	})

	t.Run("Precision", func(t *testing.T) {
		assert(t, Parse("13:14:15.9"), "0000-01-01 13:14:15.9", "13:14:15.9")
		assert(t, Parse("13:14:15.999"), "0000-01-01 13:14:15.999", "13:14:15.999")
		assert(t, Parse("2020-08-05 13:14:15.123"), "2020-08-05 13:14:15.123", "2020-08-05 13:14:15.123")
		assert(t, Parse("2020-08-05 13:14:15.123456789"), "2020-08-05 13:14:15.123456789", "2020-08-05 13:14:15.123456789")
	})

	t.Run("Timezone", func(t *testing.T) {
		// UTC / Zulu
		resZulu := Parse("2024-05-20T15:04:05Z")
		assert(t, resZulu, "2024-05-20 15:04:05", "Zulu String")
		assertZone(t, resZulu, 0, "Zulu Offset")

		// +08:00
		resP8Colon := Parse("2024-05-20T15:04:05+08:00")
		assert(t, resP8Colon, "2024-05-20 15:04:05", "+08:00 String")
		assertZone(t, resP8Colon, 8*3600, "+08:00 Offset")

		// -0700
		resM7 := Parse("2024-05-20T15:04:05-0700")
		assert(t, resM7, "2024-05-20 15:04:05", "-0700 String")
		assertZone(t, resM7, -7*3600, "-0700 Offset")

		// +05:45 (Complex)
		resComplex := Parse("2024-05-20T15:04:05+05:45")
		assert(t, resComplex, "2024-05-20 15:04:05", "+05:45 String")
		assertZone(t, resComplex, 5*3600+45*60, "+05:45 Offset")

		// 纯时间 + 时区
		resTimeZ := Parse("13:14:15Z")
		assert(t, resTimeZ, "0000-01-01 13:14:15", "13:14:15Z String")
		assertZone(t, resTimeZ, 0, "13:14:15Z Offset")
	})

	t.Run("CompactTimezone", func(t *testing.T) {
		// Year-Offset
		resYOff := Parse("2024-0700")
		assert(t, resYOff, "2024-01-01 00:00:00", "2024-0700 String")
		assertZone(t, resYOff, -7*3600, "2024-0700 Offset")

		resYOffColon := Parse("2024+08:00")
		assert(t, resYOffColon, "2024-01-01 00:00:00", "2024+08:00 String")
		assertZone(t, resYOffColon, 8*3600, "2024+08:00 Offset")

		// YYYYMMDD + Offset
		resDOff := Parse("20240520+0800")
		assert(t, resDOff, "2024-05-20 00:00:00", "20240520+0800 String")
		assertZone(t, resDOff, 8*3600, "20240520+0800 Offset")

		// YYYYMMDDTHH + Offset
		resDHOff := Parse("20240520T15+0800")
		assert(t, resDHOff, "2024-05-20 15:00:00", "20240520T15+0800 String")
		assertZone(t, resDHOff, 8*3600, "20240520T15+0800 Offset")

		// YYYYMMDDTHHmm + Offset
		resDHMOff := Parse("20240520T1504+0800")
		assert(t, resDHMOff, "2024-05-20 15:04:00", "20240520T1504+0800 String")
		assertZone(t, resDHMOff, 8*3600, "20240520T1504+0800 Offset")

		// YYYYMMDDTHHmmss + Offset
		resDHMSOff := Parse("20240520T150415+08:00")
		assert(t, resDHMSOff, "2024-05-20 15:04:15", "20240520T150415+08:00 String")
		assertZone(t, resDHMSOff, 8*3600, "20240520T150415+08:00 Offset")

		// YYYYMMDDHHmmss + Offset (No T)
		resDHMSOffNoT := Parse("20240520150415-0700")
		assert(t, resDHMSOffNoT, "2024-05-20 15:04:15", "20240520150415-0700 String")
		assertZone(t, resDHMSOffNoT, -7*3600, "20240520150415-0700 Offset")
	})

	t.Run("Normalization", func(t *testing.T) {
		assert(t, Parse("2024-13-01"), "2025-01-01 00:00:00", "2024-13-01")
		assert(t, Parse("2024-01-32"), "2024-02-01 00:00:00", "2024-01-32")
	})

	t.Run("EdgeCase", func(t *testing.T) {
		assert(t, Parse(`"2024-05-20 15:04:05"`), "2024-05-20 15:04:05", `"2024-05-20 15:04:05"`)
		assert(t, Parse("null"), "0001-01-01 00:00:00", "null")
		assert(t, Parse(""), "0001-01-01 00:00:00", "empty")
	})
}
