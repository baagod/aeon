feat: 补全核心访问器并优化文本序列化与 `String` 表示

- 在 `aeon.go` 中补全 `ISOWeek` 和 `Zone` 方法，增强标准库兼容性。
- 在 `format.go` 中引入 `DateTimeNano` 常量，并优化 `String` 方法使用 `.999999999` 实现动态精度显示。
- 实现 `MarshalText` 和 `UnmarshalText` 接口，支持 Map Key、XML 及 YAML 序列化。
- 优化 `Formats` 列表顺序，并将 `DateTimeNano` 置于高优先级以提升解析性能。
- 更新测试用例中的断言，以适配新的动态精度字符串表示。
