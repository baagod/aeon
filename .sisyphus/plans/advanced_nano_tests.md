# 计划：补全纳秒级精度测试 (Start/End/By/At/In)

## 目标
在 `cascade_abs_test.go`, `cascade_rel_test.go`, `cascade_mix_test.go` 中添加针对 `Milli`, `Micro`, `Nano` 的测试用例。

## 1. cascade_abs_test.go
- **目标**: 验证绝对定位到毫秒/微秒/纳秒。
- **用例**:
    - `StartMilli(500)`: 预期 `.500000000` (如果单位是 Milli，后续 Micro/Nano 归零)
    - `EndMilli(500)`: 预期 `.500999999` (后续单位置满)
    - `StartMicro(500, 100)`: 毫秒=500, 微秒=100. -> `.500100000`
    - `StartNano(500, 100, 50)`: 毫秒=500, 微秒=100, 纳秒=50. -> `.500100050`

## 2. cascade_rel_test.go
- **目标**: 验证相对偏移并对齐到毫秒/微秒/纳秒。
- **用例**:
    - `StartByMilli(1)`: 当前时间 +1 毫秒，然后对齐到该毫秒的开始 (Micro/Nano = 0)。
    - `EndByMicro(1)`: 当前时间 +1 微秒，然后对齐到该微秒的结束 (Nano = 999)。

## 3. cascade_mix_test.go
- **目标**: 验证混合操作。
- **用例**:
    - `StartAtMilli(500, 1)`: 绝对定位到 500ms，然后相对偏移 1us (Micro) -> `.500001000`。
    - `EndInMilli(1, 500)`: 相对偏移 1ms，然后绝对定位到 500us (Micro) -> `ms+1, .000500999`。

## 实施步骤

### A. 更新 `cascade_abs_test.go`
```go
	t.Run("纳秒精度绝对定位 (Abs Nano)", func(t *testing.T) {
        // 基准: 秒级为0
		base := Parse("2024-01-01 00:00:00") 
        
        // 1. Milli 定位: 第 500 毫秒 (注意：参数 n 对应 1-based 还是 0-based？代码逻辑 n>0 -> (n-1)*1e6。所以 StartMilli(1) 是 .000)
        // 代码回顾: case Millisecond: if n>0 { ns = (n-1)*1e6 }
        // 所以 StartMilli(500) -> .499xxx
        // 验证一下通常习惯，如果我想定位于 .500，我应该传 501？
        // 标准时间库通常 Month 是 1-based, Hour/Min/Sec 是 0-based (但这里设计似乎是 1-based 用于定位第n个?)
        // 让我们复查 applyAbs 逻辑。
        // Hour: if n>0 h=n. (如果 n=0 保持不变). 所以 StartHour(1) -> 01:00. StartHour(0) -> 保持.
        // Millisecond: if n>0 { ns = (n-1)*1e6 }. 所以 StartMilli(1) -> 0ns. StartMilli(500) -> 499ms.
        // 这似乎有点反直觉。如果是 Month, 1月是第1个。
        // 如果是 Milli，第1个毫秒确实是 0-999999ns。所以 StartMilli(1) -> 0. 
        // 那么 StartMilli(500) -> 499ms.
        // 如果想要 .500，得传 501.
        
        assert(t, base.StartMilli(1), "2024-01-01 00:00:00.000000000", "StartMilli(1) -> 0ms")
        assert(t, base.StartMilli(501), "2024-01-01 00:00:00.500000000", "StartMilli(501) -> 500ms")
        
        // EndMilli(1): 第1个毫秒的结束 -> 0ms 999999ns -> .000999999
        assert(t, base.EndMilli(1), "2024-01-01 00:00:00.000999999", "EndMilli(1) -> 0ms end")
	})
```

### B. 更新 `cascade_rel_test.go`
```go
	t.Run("纳秒精度相对位移 (Rel Nano)", func(t *testing.T) {
        // 基准: .000
        ref := Parse("2024-01-01 00:00:00")
        
        // StartByMilli(1): +1ms -> .001. Start(归零) -> .001000000
        assert(t, ref.StartByMilli(1), "2024-01-01 00:00:00.001000000", "StartByMilli(1)")
        
        // EndByMicro(1): +1us -> .000001. End(置满) -> .000001999 (因为单位是 Micro，只置满 Nano)
        // 复查 align: case Microsecond: ns = (ns/1e3)*1e3 + 999
        assert(t, ref.EndByMicro(1), "2024-01-01 00:00:00.000001999", "EndByMicro(1)")
    })
```

### C. 更新 `cascade_mix_test.go`
```go
	t.Run("纳秒精度混合操作 (Mix Nano)", func(t *testing.T) {
        base := Parse("2024-01-01 00:00:00")
        
        // StartAtMilli(501, 1): Abs Milli(501) -> 500ms. Rel Micro(1) -> 500ms + 1us = .500001. Start -> .500001000
        assert(t, base.StartAtMilli(501, 1), "2024-01-01 00:00:00.500001000", "StartAtMilli(501, 1)")
    })
```
