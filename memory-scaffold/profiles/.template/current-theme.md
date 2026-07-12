---
type: current-theme
title: 主号当前主题
description: 该号当前使用的主题和结构模板——引用仓库内置 + 个性化 CSS 覆盖
tags: [theme, template, main-account]
timestamp: 2026-07-04T10:00:00Z
account: main-account
base: default
template: ""
---
# 主号当前主题

**基础主题**：`default`（来自仓库 `themes/default.css`）

**结构模板**：`template` 字段为空表示不使用结构模板（线性 Markdown 渲染）。如需场景化结构（品牌标识区/钩子/CTA/收束），设为 `mindful-journal`、`book-club` 或 `product-launch`。用 `easygzh template list` 查看全部可用模板。

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
