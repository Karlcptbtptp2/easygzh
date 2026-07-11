---
type: current-theme
title: 主号当前主题
description: 该号当前使用的主题——引用仓库内置主题 + 个性化 CSS 覆盖
tags: [theme, main-account]
timestamp: 2026-07-04T10:00:00Z
account: main-account
base: default
---
# 主号当前主题

**基础主题**：`default`（来自仓库 `themes/default.css`）

**个性化覆盖**：在 default 基础上做下面的微调。渲染时先加载 default.css，再叠加下面的 CSS（后者优先级更高，由 juice 合并内联）。

```css
/* 个性化覆盖示例：把主色从蓝改成深青 */
#easygzh-root h1,
#easygzh-root h2,
#easygzh-root h3,
#easygzh-root strong {
  color: #00838f;
}
#easygzh-root h1 {
  border-bottom-color: #00838f;
}
#easygzh-root h2 {
  border-left-color: #00838f;
}
#easygzh-root blockquote {
  border-left-color: #00838f;
}
```

如果某天想换调性，改这里的 `base` 和覆盖 CSS 即可，历史样本不受影响。
