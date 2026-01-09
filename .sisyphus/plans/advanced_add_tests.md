# 计划：Add 系列方法的高级边界测试

## 目标
在 `cascade_add_test.go` 中补充高级测试用例，重点验证 **极大值溢出**、**DST（夏令时）边界行为** 以及 **纳秒级精度进位**，确保 `Add` 系列方法在极端场景下的鲁棒性。

## 1. DST (夏令时) 边界测试
**背景**：
- 美国洛杉矶 (America/Los_Angeles) 在 2024-03-10 02:00:00 发生夏令时切换（跳过 02:00-03:00，直接变 03:00）。
- **预期行为验证**：
    - `AddDay(1)`：应保持“墙上时间”一致。即 03-09 10:00 -> 03-10 10:00 (物理间隔 23h)。
    - `AddHour(24)`：我们的实现 (`h += 24`) 也是基于 `time.Date` 的墙上时间逻辑，所以也应该是 03-10 10:00。
    - **对比**：标准库 `t.Add(24 * time.Hour)` 是物理时间，结果应为 03-10 11:00 (物理间隔 24h)。
- **测试价值**：明确区分 `thru` 的“日历时间”语义与标准库的“物理时间”语义。

## 2. 纳秒精度与进位测试
**背景**：
- `AddNano` 累加可能导致秒、分、时的连续进位。
- `time.Date` 的 `nsec` 参数范围通常是 [0, 999999999]。
- 我们的 `applyOffset` 中 `ns += n`。
- 关键点：`align` 函数或 `time.Date` 自身是否能处理 `ns` 超过 10^9 的情况？
- **验证**：`AddNano(1_000_000_001)` -> 是否正确进位 1秒 + 1纳秒？

## 3. 极值/溢出测试
**背景**：
- 测试大数值输入，验证不会导致 panic 或明显的计算错误（Go `time` 包有年份限制，通常 +/- 292年左右是安全的，但也取决于具体实现）。
- 验证 `AddYear(1000)` 等大跨度操作。

## 实施步骤

在 `cascade_add_test.go` 中追加 `TestAdd_Advanced` 函数：

```go
func TestAdd_Advanced(t *testing.T) {
	// 1. 纳秒进位测试
	// 2024-01-30 15:04:05.999999999 + 1ns -> 2024-01-30 15:04:06.000000000
	t.Run("Nanosecond Rollover", func(t *testing.T) {
		base := ParseByLayout("2006-01-02 15:04:05.000000000", "2024-01-30 15:04:05.999999999")
		// 加 2 ns -> 应该进位到下一秒的 .000000001
		next := base.AddNano(2)
		assert(t, next, "2024-01-30 15:04:06.000000001", "AddNano 进位测试")
	})

	// 2. DST (夏令时) 墙上时间语义测试
	// 依赖 time.LoadLocation，如果环境不支持可能需要跳过或 mock
	t.Run("DST Wall Clock Behavior", func(t *testing.T) {
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			t.Skip("Skipping DST test: America/Los_Angeles location not found")
		}
		
		// 2024-03-10 02:00:00 DST 开始，时间向前跳 1 小时 (02:00 -> 03:00)
		// 设定基准：2024-03-09 10:00:00 (PST)
		base := Date(2024, 3, 9, 10, 0, 0, 0, loc)
		
		// AddDay(1) -> 应该是 2024-03-10 10:00:00 (PDT)
		// 尽管实际上只过了 23 小时，但墙上时间保持 10:00
		dayAdded := base.AddDay(1)
		// 验证小时数仍为 10
		if h := dayAdded.Hour(); h != 10 {
			t.Errorf("AddDay(1) across DST: expected hour 10, got %d", h)
		}
		
		// 对比：标准库 Add(24h)
		// 24小时后应该是 2024-03-10 11:00:00 (PDT)
		stdAdd := base.Time().Add(24 * time.Hour)
		if h := stdAdd.Hour(); h != 11 {
			t.Errorf("StdLib Add(24h) across DST: expected hour 11, got %d", h)
		}
		
		// AddHour(24) -> 我们的 AddHour 是基于日历数字加法
		// 10 + 24 = 34. 34 % 24 = 10. +1天.
		// 所以应该是 2024-03-10 10:00:00 (PDT)
		hourAdded := base.AddHour(24)
		if h := hourAdded.Hour(); h != 10 {
			t.Errorf("AddHour(24) across DST: expected hour 10 (wall clock), got %d", h)
		}
	})

    // 3. 极值测试
    t.Run("Extreme Values", func(t *testing.T) {
        base := Now()
        // 加 1000 年
        future := base.AddYear(1000)
        if future.Year() != base.Year() + 1000 {
             t.Errorf("AddYear(1000): expected year %d, got %d", base.Year()+1000, future.Year())
        }
    })
}
```
