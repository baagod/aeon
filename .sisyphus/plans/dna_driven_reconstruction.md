# DNA 驱动的极速解析器重构计划 (Final Revision)

## 1. 目标
- 回归 DNA 确定性架构，彻底消除 `ParseE` 中的盲猜路由。
- 利用 `getDNA` 产出的物理坐标加速解析，消除位置探测开销。
- 性能目标：~50ns (含 DNA 扫描)，100% 正确性，0 Panic 风险。

## 2. 核心改动点

### parse.go
- 修改 `fastParser` 类型：`func(string, []segment, *time.Location) (time.Time, error)`
- 清理 `ParseE`：
  - 移除 `switch ln` 全量代码块。
  - 修改 `master.fastParser` 的调用传参。
- 保持 `getDNA` 作为唯一路由依据。

### parse_fast.go
- 重写所有 `fastParse` 函数以接收 `segs`。
- 实现 `fastParseISO8601` 的 DNA 版本：根据 ID 序列识别 T/Z/. 符号。

## 3. 验收标准
- `go test -v .` 全量通过（无 Panic）。
- `go test -bench "Aeon/DateTime"` 性能稳定在 60ns 以内。
