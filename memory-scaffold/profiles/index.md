---
type: index
title: 公众号 Profile 列表
description: 所有公众号 profile 的导航。每个子目录是一个公众号的完整调性记忆。
---
# 公众号 Profile 列表

每个公众号一个目录。新初始化的记忆库没有激活 profile，避免把示例值误当成真实用户偏好。

当前没有已创建的 profile。`easygzh memory profile add <account>` 会在这里追加真实账号。

## 如何新建一个号

运行 `easygzh memory profile add <account>`（英文小写连字符，如 `tech-notes`），工具会从隐藏的 `.template/` 创建完整结构并立即校验。然后在对话中让 easyGZH 引导你填写真实调性。

每个 profile 目录应包含：

- `identity.md` (`type: profile-identity`) — 号名、受众、定位
- `visual-tone.md` (`type: visual-tone`) — 配色、字号、间距、分隔符、emoji 风格
- `structure-tone.md` (`type: structure-tone`) — 开头破题、小标题命名、结尾固定模块
- `current-theme.md` (`type: current-theme`) — 引用哪个主题 + 个性化 CSS 覆盖
- `samples/` (`type: sample`) — 满意的历史文章样本
