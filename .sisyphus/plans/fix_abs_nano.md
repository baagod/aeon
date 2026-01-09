# 计划：修复 TestAbsSeriesDevilMatrix 中的纳秒格式化错误

## 问题分析
`TestAbsSeriesDevilMatrix` 依然有失败：
```
cascade_abs_test.go:89: StartMilli(1) -> 0ms, got [...0.001000000], want [...0.000000000]
cascade_abs_test.go:91: StartMilli(501) -> 500ms, got [...0.501000000], want [...0.500000000]
cascade_abs_test.go:95: EndMilli(1) -> 0ms end, got [...0.001999999], want [...0.000999999]
cascade_abs_test.go:102: StartMilli(501, 101), got [...0.501100000], want [...0.501101000]
```

## 原因
1.  **Format 问题**: `StartMilli(1)` 的期望值我之前为了适配旧逻辑改成了 `.000000000` (假设n=1是0ms)，但现在逻辑已修正为 n=1ms (`.001`)。
2.  **数值逻辑不匹配**: 我之前在 `edit` 操作中修正代码为 `ns = n * 1e6`，但测试预期还没完全跟上这个变化。
3.  **StartMilli(501, 101) 的微秒计算**:
    - Got: `.501100000` (501ms + 100us?)
    - Want: `.501101000` (501ms + 101us?)
    - 让我们看 `applyAbs` 里的 `Microsecond` 逻辑：
      ```go
      case Microsecond:
          if n > 0 {
              ns = (ns/1e6)*1e6 + (n-1)*1e3  // <--- 这里还在用 n-1 !!!
          }
      ```
    - 我只修正了 `Millisecond` 的 `n-1`，忘了修正 `Microsecond` 和 `Nanosecond` 的！
    - 所以 `Micro(101)` 实际上变成了 `100us`。

## 修复方案
1.  **全面修正 `opus.go`**: 将 `Microsecond` and `Nanosecond` 的 `(n-1)` 也去掉。
2.  **全面修正 `cascade_abs_test.go`**:
    - `StartMilli(1)` -> Expect `.001000000`
    - `StartMilli(501)` -> Expect `.501000000`
    - `EndMilli(1)` -> Expect `.001999999`
    - `StartMilli(501, 101)` -> Expect `.501101000` (501ms + 101us)

## 实施步骤
1.  **修改 `opus.go`**: 去除 `applyAbs` 中 `Microsecond` 和 `Nanosecond` 分支的 `(n-1)`。
2.  **修改 `cascade_abs_test.go`**: 更新预期值以匹配 0-based 数值语义。
