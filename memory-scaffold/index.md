---
type: index
title: easyGZH 记忆库
description: easyGZH 的本地公众号调性记忆入口，导航全局偏好、账号 profile、主题和维护日志。
tags:
  - memory
  - easygzh
  - index
timestamp: 2026-07-10T00:00:00Z
---

# easyGZH 记忆库

本目录是 easyGZH 的本地调性记忆库，遵循 [OKF 协议](https://github.com/GoogleCloudPlatform/knowledge-catalog/blob/main/okf/SPEC.md) 与 OpenKnowledge (OK) 结构。

## 顶层概念

- [全局偏好](preferences.md) — 跨账号的输入/交付偏好
- [公众号 Profile](profiles/index.md) — 每个公众号一个目录，存放该号的视觉/结构调性
- [私有主题](themes/index.md) — 用户自定义或覆盖的主题 CSS

## 结构

```
.
├── preferences.md        # type: global-preferences
├── profiles/
│   ├── index.md
│   ├── .template/        # 隐藏模板，不会作为真实偏好读取
│   └── <account-name>/   # profile add 创建的真实账号
│       ├── identity.md         # type: profile-identity
│       ├── visual-tone.md      # type: visual-tone
│       ├── structure-tone.md   # type: structure-tone
│       ├── current-theme.md    # type: current-theme
│       └── samples/            # type: sample
└── themes/
    └── <theme-name>.md    # type: theme
```

每个 `.md` 文件都有 YAML frontmatter，必填字段是 `type`（值见 `references/memory-schema.md` 的词汇表）。OK 运行时还要求 `title` + `description`。

变更记录见 [log.md](log.md)。
