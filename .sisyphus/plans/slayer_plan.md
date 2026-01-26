# “屠神”计划：超越 iso8601 性能方案

## 目标
1. 实现 `Aeon` 在处理 ISO8601/DT 格式时，性能进入 **30ns** 以内。
2. 彻底超越 `github.com/relvacode/iso8601` (约 30-40ns)。
3. 保持 **0 Allocations**。

## 核心技术点
- **Fast Path Hook**: 利用现有的 `getDNA` 预判，绕过 `time.Parse` 解释器。
- **SWAR / Loop Unrolling**: 手写数字解析，消除所有 `for` 循环和 `switch` 分支。
- **Direct time.Date**: 解析出整数后直接构造 `time.Time`，跳过所有中间层。

## 详细步骤

### 1. 基础设施 (parse.go)
- [ ] 定义 `type fastParser func(string, *time.Location) (time.Time, error)`
- [ ] 在 `layoutInfo` 结构体中添加 `fastParser fastParser` 字段。
- [ ] 修改 `ParseE`：在获取 `master` 后，优先检查并执行 `master.fastParser`。

### 2. 高性能解析引擎 (parser_fast.go)
- [ ] 实现 `parse2(s string) int`: 解析 2 位数字（月、日、时、分、秒）。
- [ ] 实现 `parse4(s string) int`: 解析 4 位数字（年份）。
- [ ] 实现 `fastParseISO8601(s string, loc *time.Location) (time.Time, error)`。
- [ ] 实现 `fastParseDT(s string, loc *time.Location) (time.Time, error)`。

### 3. 指纹注册 (parse.go)
- [ ] 在 `init()` 或 `RegisterDNA` 中，通过字符串匹配为常用布局（DT, ISO8601, DateOnly）手动分配对应的 `fastParser`。

### 4. 验证与对决 (bench_test.go)
- [ ] 增加与 `iso8601.ParseString` 的直接 Benchmark 对比。
- [ ] 验证时区解析逻辑的正确性（特别是 Z 和 +08:00）。
