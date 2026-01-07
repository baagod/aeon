# Aeon 逻辑法典 (Project Aeon Logic Definitions)

## 1. By 系列核心语义 (Relative Offset + Alignment)
- **定义**: `By` 系列 = `AddTime` (相对偏移) + `Start/End` (边界对齐)。
- **对齐优先**: 当操作单位大于 Day 时，`Start` 系列会先执行偏移保护，随后立即将子级单位归零。
    - 示例: `StartByMonth(1)` 从 1月31日开始 -> 2月29日(保护) -> 2月1日(对齐)。

## 2. At 系列核心语义 (Absolute Positioning + Relative Offset)
- **定义**: `At` 系列 = `Start/End` (绝对定位) + `AddTime` (相对偏移)。
- **锚点优先**: 先执行绝对定位锚点，再进行相对偏移。
- **参数顺序**: 第一个参数为绝对定位值，后续参数为相对偏移值。
    - 示例: `StartAtYear(2024, 1, 2)` -> 锚定到 2024 年年初 -> 偏移 1 月 2 天 -> `2024-02-03 00:00:00`。

## 3. In 系列核心语义 (Relative Offset + Absolute Positioning)
- **定义**: `In` 系列 = `AddTime` (相对偏移) + `Start/End` (绝对定位)。
- **偏移优先**: 先执行相对偏移，再进行绝对定位锚点。
- **参数顺序**: 第一个参数为相对偏移值，后续参数为绝对定位值。
    - 示例: `StartInYear(1, 6, 15)` -> 偏移 1 年 -> 锚定到 6 月 15 日 -> `2025-06-15 00:00:00`。

## 4. 溢出与保护 (Overflow Flag)
- **容器单位 (Month/Year等)**: 始终执行日期保护（不溢出）。
- **位移单位 (Day/Hour等)**: 默认允许自然溢出（Add 灵魂）。
- **强制保护标志 (`Overflow`)**: 若传入此标志，非容器单位也将执行强制日期保护（截断到当前月末）。
    - 示例: 4月15日 `StartByDay(Overflow, 45)` -> 4月60日 -> 截断至 4月30日。

## 5. 周逻辑
- **基准线对齐**: 位移前先对齐到当前周起始点。
- **YearWeek(ISO)**: 自动映射为 `ISOYearWeek`。

## 6. 四位一体总结
- **Start/End**: 全绝对定位
- **StartBy/EndBy**: 全相对偏移 + 边界对齐
- **StartAt/EndAt**: 先绝对定位，后相对偏移
- **StartIn/EndIn**: 先相对偏移，后绝对定位
