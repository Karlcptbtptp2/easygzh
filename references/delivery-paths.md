# 交付路径实现细节

按 SKILL 核心契约，从自动化最高开始试，不通逐级降级。剪贴板永远是兜底。

## 路径 1：浏览器自动填入公众号后台（最高自动化）

**前提**：agent 具备浏览器控制能力（MCP browser 工具 / playwright / puppeteer 类）。

**步骤**：
1. 打开 `https://mp.weixin.qq.com`，确认已登录（未登录则提示用户登录，不强求）。
2. 进入"图文素材 → 新建图文"。
3. 定位编辑器内容区（`rich-text-editor` 类的 contenteditable div）。
4. 把渲染好的 inline HTML 注入（设置 innerHTML，或模拟 paste 事件）。
5. 交给用户点"保存/发布"。

**风险与降级**：
- 公众号后台是强校验 SPA，有反自动化检测，很可能失败。
- 失败信号：登录态丢失、编辑器没响应、内容被清洗。
- **失败即降级到路径 2，不要卡住、不要重试到死。**

**MVP 状态**：保留为加分路径，主流程不依赖它。

## 路径 2：剪贴板（兜底主路径）

**前提**：能调用系统剪贴板命令。

**按平台**：
- **macOS**：`echo "$HTML" | pbcopy`（pbcopy 支持 HTML? —— 注意：pbcopy 默认是纯文本。要写 HTML 到剪贴板让公众号识别为富文本，需要用 macOS 的富文本剪贴板，见下方"富文本剪贴板"。）
- **Windows**：`echo "$HTML" | clip`（同样是纯文本）。
- **Linux**：`xclip -selection clipboard -t text/html`（支持 HTML MIME type）。

**富文本剪贴板（关键）**：

公众号编辑器粘贴时，识别的是剪贴板的 `text/html` MIME 类型，不是纯文本。纯文本粘贴进去就是一串 HTML 源码，不会渲染。

- **macOS** 写 HTML 到富文本剪贴板：
  ```bash
  printf '%s' "$HTML" | pbcopy  # 这是纯文本，不对
  # 正确做法：用 osascript 写 public.html
  osascript -e "set the clipboard to (read (POSIX file \"$TMPFILE\") as «class HTML»)"
  ```
- **跨平台最稳的做法**：把 HTML 写进一个临时 `.html` 文件，让用户**在浏览器打开 → Ctrl+A 全选 → Ctrl+C 复制 → 粘贴到公众号**。这绕过剪贴板 MIME 的复杂性，兼容性最好。

**easyGZH 的实现策略（核心契约：自动化优先 + 降级）**：
1. 先尝试富文本剪贴板（平台特定命令）。
2. 不行 → 写临时 HTML 文件 + 自动用浏览器打开 + 提示用户"全选复制粘贴"（路径 3 的变体）。

## 路径 3：本地文件

**前提**：什么都不可用时。

**步骤**：
1. 把 inline HTML 写到 `~/.easyGZH/output/<timestamp>.html`。
2. 输出文件路径。
3. 给预览方式：`open <file>`（macOS）/ `start <file>`（Windows）/ `xdg-open <file>`（Linux）。
4. 指引用户：浏览器打开 → 全选 → 复制 → 粘贴到公众号后台。

## 决策流程图

```
agent 有 browser 能力?
├─ 是 → 尝试自动填后台
│       ├─ 成功 → 用户点发布 (最丝滑)
│       └─ 失败 → 降级
└─ 否 ↓
能写富文本剪贴板?
├─ 是 → 写剪贴板 → 提示粘贴到后台
└─ 否 ↓
写本地 HTML 文件 + 打开浏览器 → 用户全选复制粘贴 (最终兜底)
```

## 注意

- **永远给用户一个能用的产物**。哪怕所有自动化都失败，本地 HTML 文件 + 手动复制粘贴也一定可用。
- **不要把失败藏起来**。降级时告诉用户"自动填后台没成功，已改用剪贴板，原因是 X"。
- 渲染产物本身（inline HTML）在任何路径下都是同一个字符串，只是交付方式不同。
