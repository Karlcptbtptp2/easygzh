---
type: index
title: 公众号 Profile 列表
description: 所有公众号 profile 的导航。每个子目录是一个公众号的完整调性记忆。
---
# 公众号 Profile 列表

每个公众号一个目录。新初始化的记忆库没有激活 profile，避免把示例值误当成真实用户偏好。

当前没有已创建的 profile。

## 如何新建一个号

告诉 AI「我要为我的公众号建一个 profile」，AI 会从隐藏的 `.template/` 复制完整结构到 `<account-name>/` 目录，然后在对话中引导你填写真实调性。

也可以手动创建：把 `.template/` 复制为你的账号名目录（英文小写连字符，如 `tech-notes`），修改里面的 `identity.md` / `visual-tone.md` / `structure-tone.md`。

每个 profile 目录包含：

- `identity.md` (`type: profile-identity`) — 号名、受众、定位
- `visual-tone.md` (`type: visual-tone`) — 配色、字号、间距、分隔符、emoji 风格（AI 生成 HTML 时直接据此写内联样式）
- `structure-tone.md` (`type: structure-tone`) — 开头破题、小标题命名、结尾固定模块
- `samples/` (`type: sample`) — 满意的历史文章样本
