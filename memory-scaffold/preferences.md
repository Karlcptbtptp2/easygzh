---
type: global-preferences
title: 全局偏好
description: 跨所有公众号 profile 的输入方式、交付方式、默认行为偏好
tags: [preferences, global]
timestamp: 2026-07-04T10:00:00Z
---
# 全局偏好

skill 每次启动会读取这里来决定默认行为。下面是初始模板，可在对话中让 easyGZH 更新。

## 输入方式偏好

（示例值，按需修改）

- 首选输入：粘贴文本
- 备选：本地 .md 文件路径

## 交付方式偏好

- 首选交付：剪贴板（写好后我手动粘贴到公众号后台）
- 若 agent 有浏览器能力：尝试自动填入后台，失败则降级回剪贴板

## 其他

- 渲染时是否默认把外链转成脚注引用：是（公众号正文外链不可点）
- 是否在结尾自动追加固定模块：否（由每个 profile 的 structure-tone 决定）
