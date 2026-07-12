---
type: index
title: easyGZH 记忆库
description: easyGZH 的本地公众号调性记忆入口，导航全局偏好、账号 profile 和维护日志。
tags:
  - memory
  - easygzh
  - index
timestamp: 2026-07-10T00:00:00Z
---

# easyGZH 记忆库

本目录是 easyGZH 的本地调性记忆库，遵循 [OKF 协议](https://github.com/GoogleCloudPlatform/knowledge-catalog/blob/main/okf/SPEC.md) 与 OpenKnowledge (OK) 结构。AI 每次排版前读这里理解用户，排版后把满意结果和偏好更新写回这里。

## 顶层概念

- [全局偏好](preferences.md) — 跨账号的输入/交付偏好
- [公众号 Profile](profiles/index.md) — 每个公众号一个目录，存放该号的视觉/结构调性

## 结构

```
.
├── preferences.md        # type: global-preferences
├── profiles/
│   ├── index.md
│   ├── .template/        # 隐藏模板，不会作为真实偏好读取
│   └── <account-name>/   # 创建的真实账号
│       ├── identity.md         # type: profile-identity
│       ├── visual-tone.md      # type: visual-tone（AI 写 HTML 时参照的内联样式规范）
│       ├── structure-tone.md   # type: structure-tone
│       └── samples/            # type: sample（满意的历史文章）
```

每个 `.md` 文件都有 YAML frontmatter，必填字段是 `type`（值见 `references/memory-schema.md` 的词汇表）。OK 运行时还要求 `title` + `description`。

变更记录见 [log.md](log.md)。
