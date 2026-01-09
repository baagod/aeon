test: 增加 Add 系列方法的高级边界测试用例

- 在 `cascade_add_test.go` 中新增 `TestAdd_Advanced` 测试集
- 验证纳秒级进位 (Rollover) 的正确性
- 验证夏令时 (DST) 切换时的“墙上时间”语义：`AddDay(1)` 保持小时数不变，而非物理 24 小时
- 验证大数值年份计算的稳定性
