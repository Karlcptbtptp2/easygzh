# 记忆库 Schema 与 type 词汇表

easyGZH 的本地记忆库采用 **OKF 风格的 `type` 词汇表**，并以 **OpenKnowledge (OK) 运行时**作为实际读写和校验契约。为保证当前 OK 版本能完整索引，所有 Markdown（包括 `index.md` 和 `log.md`）都带 `title`、`description` 和 `type`。

## 1. OKF 的硬性要求（极简）

- 记忆库是一个目录（bundle），通常作为 git 仓库分发。
- 每个 `.md` 文件有 YAML frontmatter，**唯一必填字段是 `type`**（非空字符串，值自定义，无需注册）。
- 推荐字段：`title`、`description`、`tags`、`timestamp`。
- `index.md` 使用 `type: index`，负责稳定导航；`log.md` 使用 `type: changelog`，按 `## YYYY-MM-DD` 倒序记录变更。
- 用标准 Markdown 链接互联；初始化和维护后必须运行 `easygzh memory validate`，不能把坏链留给消费者。
- 所有文档都保留开放 frontmatter，消费者必须保留未知字段。

## 2. OK 运行时的增量要求

- 项目根有 `.ok/config.yml`（内容就是 `content: { dir: "." }`）。
- 每个 doc 的 frontmatter 要有 `title` + `description`（OK 比 OKF 严），`tags` 推荐，`type` 为 OKF 可移植性保留。
- 文件夹可用嵌套 `.ok/frontmatter.yml` 存文件夹元信息（**self-only，不级联**），`.ok/templates/` 存新文档模板。

## 3. easyGZH 的 type 词汇表

这些 `type` 值是 easyGZH 自定义的（OKF 允许任意值）：

| type | 用途 | 存放位置 |
|---|---|---|
| `global-preferences` | 跨账号的输入/交付偏好 | `preferences.md` |
| `profile-identity` | 一个公众号的身份/受众/定位 | `profiles/<account>/identity.md` |
| `visual-tone` | 视觉调性（配色/字号/间距/分隔符/emoji） | `profiles/<account>/visual-tone.md` |
| `structure-tone` | 结构习惯（开头/小标题/结尾模块） | `profiles/<account>/structure-tone.md` |
| `current-theme` | 当前主题引用 + 覆盖 CSS | `profiles/<account>/current-theme.md` |
| `theme` | 独立的私有主题 CSS | `themes/<name>.md` |
| `sample` | 满意的历史文章样本 | `profiles/<account>/samples/*.md` |
| `index` | 导航文件（OK 习惯，OKF 里 index.md 无 frontmatter，但 OK 允许子目录 index 带 frontmatter） | 各级 `index.md` |

## 4. 目录结构总览

```
~/.easygzh/memory/                    # 记忆库根（可用 EASYGZH_MEMORY_DIR 覆盖）
├── .ok/config.yml                    # OK 项目配置
├── index.md                          # 根导航
├── log.md                            # 变更日志
├── preferences.md                    # type: global-preferences
├── themes/
│   ├── index.md                      # type: index
│   └── <name>.md                     # type: theme
└── profiles/
    ├── index.md                      # type: index
    ├── .template/                    # 隐藏模板，不是激活 profile
    └── <account>/                    # 由 profile add 创建的真实账号
        ├── .ok/frontmatter.yml       # 文件夹元信息（self-only）
        ├── index.md                  # type: index
        ├── identity.md               # type: profile-identity
        ├── visual-tone.md            # type: visual-tone
        ├── structure-tone.md         # type: structure-tone
        ├── current-theme.md          # type: current-theme
        └── samples/
            ├── index.md              # type: index
            └── YYYY-MM-DD-*.md       # type: sample
```

## 5. 如何被 AI 快速检索

- **有 OK MCP 时**：用 `mcp__open-knowledge__exec({ command: "cat <path>" })` 精确读，或 `search({ query })` 排序检索（title boost + BM25 + recency）。
- **无 OK MCP 时**：降级为 `Read`/`Grep` 直接读文件。目录结构本身可被 grep 定位（`Grep "type: visual-tone"` 能找到所有视觉调性文件）。
- **稳定路径访问**：SKILL Stage 2 用确定路径（如 `profiles/<account>/visual-tone.md`）读，保证每次调性一致。

## 6. frontmatter 示例（visual-tone）

```yaml
---
type: visual-tone
title: 主号视觉调性
description: 主公众号配色、字号、间距规范
tags: [visual, main-account]
timestamp: 2026-07-04T10:00:00Z
account: main-account
---
```

`account` 是 easyGZH 额外加的便捷字段（OKF 允许任意额外 key，消费者必须保留）。它让检索/过滤更直接。

## 7. 初始化路径（与 SKILL Stage 2b 对应）

- **路线 1**：`easygzh memory init` 从二进制内嵌脚手架初始化，并立即运行 `easygzh memory validate`（最高自动化）。
- **路线 2**：`easygzh memory profile add <account>` 创建新账号 profile，再由 Agent 引导填写。
- **路线 3**：从满意文章反推（贴文章 → AI 提炼 → 生成 profile）。

诊断命令：`easygzh memory status --json` 会报告目录是否存在、profile 列表和校验问题。
