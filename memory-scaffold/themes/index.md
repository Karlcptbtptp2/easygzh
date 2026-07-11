---
type: index
title: 私有主题
description: 用户自定义或覆盖的主题 CSS。内置主题在仓库的 themes/ 目录，这里的主题是个人化的。
---
# 私有主题

存放只属于你的主题（CSS）。内置主题在 easyGZH 仓库的 `themes/` 目录（default、lively 等）；这里放的是你在内置主题基础上修改、或全新创作的私有主题。

新建主题：复制一份 `templates/visual-tone.md`，或让 easyGZH 在对话中帮你生成。

## 主题文件格式

```markdown
---
type: theme
title: 我的主题
description: 基于 default 微调，主色改为绿色
tags: [theme]
base: default          # 可选：基于哪个内置主题叠加
timestamp: 2026-07-04T10:00:00Z
---
#easygzh-root { ... }   # CSS 主体，会被 juice 内联
```

`base` 字段告诉 easyGZH：渲染时先加载仓库里的 `themes/<base>.css`，再叠加本文件里的 CSS（后者优先级更高）。
