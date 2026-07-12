---
type: changelog
title: 变更日志
description: 按时间倒序记录 easyGZH 记忆库的初始化、profile 变更和校验活动。
tags:
  - memory
  - changelog
timestamp: 2026-07-10T00:00:00Z
---

# 变更日志

## 2026-07-12

**纯 Skill 重构** — 删除 Go 引擎和 Node 脚本，easyGZH 变为纯 Agent Skill。AI 直接写公众号 inline-styled HTML，不再走 Markdown→HTML 转换。移除 CSS 主题系统（`current-theme.md` / `themes/`），视觉规范由 `visual-tone.md` 承载，AI 生成 HTML 时直接参照写内联样式。

## 2026-07-10

**Health repair** — 补齐 OpenKnowledge frontmatter；将示例 profile 移入隐藏 `.template/`，避免被当作真实偏好。

## 2026-07-04

**Creation** — 初始化 easyGZH 记忆库脚手架。
