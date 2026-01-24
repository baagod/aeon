# 实施计划：Aeon DST 增强支持 (Project Twilight)

## 1. 目标 (Objectives)

在不破坏现有 API 和性能的前提下，赋予 Aeon 处理夏令时 **二义性时间 (Ambiguous Time)** 和 **缺失时间 (Missing Time)** 的能力。

*   **二义性时间**：秋季回拨（Fall Back），如 01:30 出现两次。需要支持选择 "Earlier" (第一次) 或 "Later" (第二次)。
*   **缺失时间**：春季拨快（Spring Forward），如 02:30 直接消失。需要支持 "Strict" (报错) 或 "Forward" (顺延)。

## 2. 数据结构变更 (aeon.go)

在 `aeon.Time` 结构体中增加 `dstPolicy` 字段。

```go
// aeon.go

type Time struct {
    time       time.Time
    weekStarts time.Weekday
    
    // dstPolicy 控制遇到 DST 边界时的行为
    // 0: Default (Go 标准库行为：通常取 First/Earlier，跳过 Missing)
    // 1: Strict  (严格模式：遇到 Missing 或 Ambiguous 需明确处理，否则报错)
    // 2: Later   (偏好后者：回拨时取第二次出现的时间)
    // 3: Earlier (偏好前者：回拨时取第一次出现的时间 - 显式)
    dstPolicy uint8 
}

// 增加常量定义
const (
    dstDefault = iota
    dstStrict
    dstLater
    dstEarlier
)
```

## 3. API 设计 (aeon.go)

增加一组流式配置方法，返回带有新策略的 `Time` 副本。

```go
// aeon.go

// DSTLater 策略：在 DST 回拨重叠期，优先选择“第二次”出现的时间（即较晚的时刻）。
func (t Time) DSTLater() Time {
    t.dstPolicy = dstLater
    return t
}

// DSTStrict 策略：在 DST 导致时间消失（Spring Forward）时返回错误（需配合 ParseE 或特定 API，待定）。
// 注：对于 Go/Sh 这种不返回 error 的链式 API，Strict 模式可能表现为 Panic 或返回零值，需慎重设计。
func (t Time) DSTStrict() Time {
    t.dstPolicy = dstStrict
    return t
}

// DSTEarlier 策略：显式要求“第一次”出现的时间（通常是默认行为）。
func (t Time) DSTEarlier() Time {
    t.dstPolicy = dstEarlier
    return t
}
```

## 4. 核心逻辑实现 (helper.go & cascade.go)

### 4.1 算法设计：`resolveDST` (helper.go)

这是核心的“纠偏”函数。它将在 `cascade` 级联运算生成基础时间**之后**被调用。建议放入 `helper.go` 或 `aeon.go` 中。

```go
// helper.go

func (t Time) resolveDST(candidate time.Time) time.Time {
    // 0. 如果是默认策略，直接返回（零性能损耗）
    if t.dstPolicy == dstDefault {
        return candidate
    }

    // 1. 处理 "Later" (想要回拨后的时间)
    if t.dstPolicy == dstLater {
        // 前提：candidate 是 IsDST=true (说明 Go 默认给了第一次)，且用户想要第二次
        if candidate.IsDST() {
            // 试探：往后推 1 小时 (物理时间)
            // 注意：这里需要获取当前时区的实际 dstOffset 差值，通常是 1h，严谨做法是计算 Zone 差
            // 简单实现：
            check := candidate.Add(1 * time.Hour)
            
            // 判据：如果物理时间变了，但“墙上时钟”没变，说明遇到了回拨分身
            h1, m1, _ := candidate.Clock()
            h2, m2, _ := check.Clock()
            if h1 == h2 && m1 == m2 && !check.IsDST() {
                return check // 偷梁换柱，返回第二次
            }
        }
    }

    // 2. 处理 "Earlier" (想要回拨前的时间)
    // Go 默认通常就是 Earlier，但在某些特殊 OS 或 Go 版本下，如果默认给了 Later，
    // 这里需要反向逻辑 (Add(-1 * time.Hour)) 把它纠回来。
    
    // 3. 处理 "Strict" (缺失时间检测)
    // 比如用户求 02:30，但 candidate 变成了 03:30 (被自动顺延了)
    if t.dstPolicy == dstStrict {
        // 需要传入用户原本期望的 h, m, s
        // 如果 candidate.Hour() != expectedHour，说明发生了跳跃
        // Panic or Log (具体实现需结合 cascade 中的预期值)
    }

    return candidate
}
```

### 4.2 集成点 (cascade.go)

修改 `cascade.go` 中的 `cascade` 函数。

```go
// cascade.go

func cascade(...) Time {
    // ... 原有的级联计算逻辑 ...
    // ... y, m, d, h, mm, s, ns, w = apply(...) ...
    // ... y, m, d, h, mm, s, ns = align(...) ...

    // 原逻辑：直接返回
    // return Time{time: time.Date(y, time.Month(m), d, h, mm, s, ns, t.Location()), ...}
    
    // 新逻辑：注入 DST 钩子
    rawTime := time.Date(y, time.Month(m), d, h, mm, s, ns, t.Location())
    finalTime := t.resolveDST(rawTime) // <--- 注入点
    
    return Time{
        time: finalTime,
        weekStarts: t.weekStarts,
        dstPolicy: t.dstPolicy, // 传递策略
    }
}
```

## 5. 测试策略 (Verification)

必须在测试中使用真实的 DST 时区（如 `America/New_York`），不能依赖本地时间。

*   **Test Case A (Fall Back / Later):**
    *   时间：2024-11-03 01:30:00 (NY)
    *   动作：`Go(2024, 11, 3, 1, 30)`
    *   断言：默认返回 `IsDST=true`, `Offset=-04`。
    *   动作：`DSTLater().Go(...)`
    *   断言：返回 `IsDST=false`, `Offset=-05`, UnixTimestamp 比前者大 3600s。

*   **Test Case B (Spring Forward / Gap):**
    *   时间：2024-03-10 02:30:00 (NY, 此时不存在)
    *   动作：`Go(..., 2, 30)`
    *   断言：Go 默认可能返回 03:30 或 01:30 (取决于实现)。
    *   动作：`DSTStrict().Go(...)`
    *   断言：检测到时分不匹配，触发 Panic/Error。
