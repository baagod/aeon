perf: 实现 O(1) 级分桶解析引擎并统一 DT/D 缩写体系

- 引入“双轨解析架构”：在 `init()` 中建立长度索引桶 (`buckets`) 与变长布局轨道 (`Formats`)，实现解析路径的极速分发。
- 优化 `ParseE` 逻辑：优先执行长度精准匹配， fallback 至动态特征轨道，彻底消灭全量循环遍历。
- 全面实施布局常量重命名：将 `DateTime` 简化为 `DT`，`DateOnly` 简化为 `D`，提升代码视觉密度。
- 补齐 `Time` 结构体的 `Format`、`AppendFormat` 及 `ToString` 方法，对齐高性能输出能力。
- 新增 `parse_test.go` 覆盖全量归一化与分桶逻辑；更新 Benchmark 验证解析性能优势（最高领先 Carbon 38倍）。
- 修正 `init()` 中的变长因子识别，支持对不带前导零的宽松格式进行动态分流。
