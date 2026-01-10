refactor: 全面重命名 `Jump` 系列为 `Go` 系列并同步语义

- 将 `Jump`、`JumpBy`、`JumpAt`、`JumpIn` 统一重构为 `Go`、`GoBy`、`GoAt`、`GoIn` 系列，消除语义歧义并提升 API 动词一致性。
- 同步更新级联单元方法（如 `GoYear`、`GoMonth` 等）及内部常量 (`fromGoAbs` 等)。
- 将 `Flag.jump` 字段重命名为 `Flag.goMode`，确保内部逻辑语义自洽。
- 修正级联导航系统中的注释与文档描述，对齐“Start, End, Add, Go”的核心动词矩阵。
