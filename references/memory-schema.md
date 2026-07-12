# 记忆库 Schema 与 type 词汇表

easyGZH 的本地记忆库采用 **OKF 风格的 `type` 词汇表**，以 OpenKnowledge (OK) 结构组织。所有 Markdown（包括 `index.md` 和 `log.md`）都带 `title`、`description` 和 `type`，便于 AI 检索和管理。

## 1. OKF 的硬性要求（极简）

- 记忆库是一个目录（bundle），通常作为 git 仓库分发
- 每个 `.md` 文件有 YAML frontmatter，**唯一必填字段是 `type`**（非空字符串）
- 推荐字段：`title`、`description`、`tags`、`timestamp`
- `index.md` 使用 `type: index`，负责稳定导航；`log.md` 使用 `type: changelog`，按 `## YYYY-MM-DD` 倒序记录变更
- 用标准 Markdown 链接互联
- **禁止写入** API key、token、cookie、密码或高敏个人信息

## 2. OK 运行时的增量要求

- 项目根有 `.ok/config.yml`（内容就是 `content: { dir: "." }`）
- 每个 doc 的 frontmatter 要有 `title` + `description`
- 文件夹可用嵌套 `.ok/frontmatter.yml` 存文件夹元信息

## 3. easyGZH 的 type 词汇表

| type | 用途 | 存放位置 |
|---|---|---|
| `global-preferences` | 跨账号的输入/交付偏好 | `preferences.md` |
| `profile-identity` | 一个公众号的身份/受众/定位 | `profiles/<account>/identity.md` |
| `visual-tone` | 视觉调性（配色/字号/间距/分隔符/emoji）—— AI 生成 HTML 时直接据此写内联样式 | `profiles/<account>/visual-tone.md` |
| `structure-tone` | 结构习惯（开头/小标题/结尾模块） | `profiles/<account>/structure-tone.md` |
| `sample` | 满意的历史文章样本 | `profiles/<account>/samples/*.md` |
| `index` | 导航文件 | 各级 `index.md` |

## 4. 目录结构总览

```
~/.easygzh/memory/                    # 记忆库根（可用 EASYGZH_MEMORY_DIR 覆盖）
├── .ok/config.yml                    # OK 项目配置
├── index.md                          # 根导航
├── log.md                            # 变更日志
├── preferences.md                    # type: global-preferences
└── profiles/
    ├── index.md                      # type: index
    ├── .template/                    # 隐藏模板，不是激活 profile
    └── <account>/                    # 由 AI 创建的真实账号
        ├── .ok/frontmatter.yml       # 文件夹元信息
        ├── index.md                  # type: index
        ├── identity.md               # type: profile-identity
        ├── visual-tone.md            # type: visual-tone
        ├── structure-tone.md         # type: structure-tone
        └── samples/
            ├── index.md              # type: index
            └── YYYY-MM-DD-*.md       # type: sample
```

## 5. AI 如何读写记忆库

AI 用文件读写工具（Read / Write / Edit / Grep / Glob）直接操作记忆库文件。

- **读 profile**：用确定路径（如 `profiles/<account>/visual-tone.md`）稳定读取，保证每次调性一致
- **写样本**：满意的文章存为 `profiles/<account>/samples/YYYY-MM-DD-<slug>.md`
- **更新偏好**：用户说"我喜欢..."时，更新对应 `visual-tone.md` 或 `structure-tone.md`
- **记录变更**：每次重要变更在 `log.md` 加一条 `## YYYY-MM-DD` 倒序记录

## 6. 初始化

首次使用时，AI 把仓库的 `memory-scaffold/` 目录复制到 `~/.easygzh/memory/`：

```bash
cp -r memory-scaffold/ ~/.easygzh/memory/
```

或 AI 直接用文件工具创建等价的目录结构。

新建账号 profile 时，AI 把 `.template/` 复制为 `<account-name>/` 目录，然后引导用户填写 `identity.md` / `visual-tone.md` / `structure-tone.md`。

## 7. 校验检查清单

AI 修改记忆库后自检：
- [ ] 每个文件都有 `type` 字段（非空）
- [ ] 每个文件都有 `title` 和 `description`
- [ ] 文件内的 Markdown 链接指向的文件都存在
- [ ] 没有写入 API key / token / 密码（grep 检查 `sk-` / `password` / `secret` / `token`）
- [ ] `log.md` 有对应的变更记录
