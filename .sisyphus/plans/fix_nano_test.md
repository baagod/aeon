# 计划：修复 TestAbsSeriesDevilMatrix 中的纳秒定位错误

## 问题 1: StartMilli(1) 格式化问题
```
cascade_abs_test.go:89: StartMilli(1) -> 0ms, got [2024-01-01 00:00:00], want [2024-01-01 00:00:00.000000000]
```
- **实际值**: 时间确实是 `00:00:00` (纳秒为 0)。
- **期望值**: 字符串包含 `.000000000`。
- **原因**: `Time.String()` 方法（在 `format.go` 中）有一个优化：`if ns := t.time.Nanosecond(); ns == 0 { return t.time.Format(DateTime) }`。
- **分析**: 当纳秒为 0 时，默认不输出纳秒部分。测试断言强行要求输出 `.000000000`，导致失败。
- **修复**: 修改测试预期字符串，去除 `.000000000`。因为对于 `thru` 来说，`0ns` 和 `0s` 在字符串表示上是一样的（除非强制 Format）。

## 问题 2: StartMicro(501, 101) 逻辑/级联错误
```
cascade_abs_test.go:99: StartMicro(501, 101), got [2024-01-01 00:00:00.000500100], want [2024-01-01 00:00:00.500100000]
```
- **参数**: `StartMicro(501, 101)`。注意 `StartMicro` 是入口方法，意味着 `Unit` 序列从 `Microsecond` 开始。
- **期望 (Want)**: `.500100000` (500ms + 100us)。
- **实际 (Got)**: `.000500100` (500us + 100ns?)。
- **原因分析**:
    - 调用 `StartMicro` -> `cascade(..., Microsecond, ...)`。
    - `sequence(Microsecond)` 返回什么？查看 `opus.go`:
      ```go
      default:
          if u <= Nanosecond {
              return stdSeq[u:]
          }
      ```
      `stdSeq` 的定义是 `Century, ..., Second, Millisecond, Microsecond, Nanosecond`。
      所以 `stdSeq[Microsecond:]` 是 `[Microsecond, Nanosecond]`。
    - **关键错误**: `StartMicro` 的参数 `501` 被应用到了 `Microsecond` 上（即第 501 个微秒 = 500us），参数 `101` 被应用到了 `Nanosecond` 上（即第 101 个纳秒 = 100ns）。
    - 结果：`500us + 100ns` = `500100ns` = `.000500100`。
    - 我的原意（期望值）是想测试 `Milli=500` + `Micro=100`。但这应该调用 `StartMilli(501, 101)` (参数1->Milli, 参数2->Micro)。
    - **结论**: 测试代码写错了，应该调用 `StartMilli` 来实现期望的层级，或者调整 `StartMicro` 的预期值。
    - **修复**: 将测试调用改为 `StartMilli(501, 101)` 以匹配预期值 `.500100000`。

## 实施步骤
修改 `cascade_abs_test.go`:
1.  修正 `StartMilli(1)` 的预期值，去掉纳秒后缀。
2.  将 `StartMicro(501, 101)` 改为 `StartMilli(501, 101)`，因为目的是测试 `ms -> us` 的级联。

```go
		// 1. Milli 定位 (1-based index)
		// StartMilli(1) -> 第1个毫秒 -> 0ns (Format自动隐藏0ns)
		assert(t, base.StartMilli(1), "2024-01-01 00:00:00", "StartMilli(1) -> 0ms")
		
        // ...

		// 3. 级联定位 (Milli -> Micro)
		// StartMilli(501, 101) -> 500ms + 100us = .500100000
		assert(t, base.StartMilli(501, 101), "2024-01-01 00:00:00.500100000", "StartMilli(501, 101)")
```
