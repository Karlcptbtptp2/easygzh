---
name: easyGZH
description: "Use easyGZH to format Markdown into WeChat-public-account (公众号) HTML and publish it. Trigger when the user says '公众号排版''排版我的文章''给这篇文章套样式''把这段排到公众号''easygzh', or wants a stable per-account visual/structural tone. Operate the `easygzh` CLI (a compiled Go binary): convert/preview/inspect/publish/image/memory/theme/doctor/skills. Keep a per-account tone stable via the local OKF memory store. Do NOT use for Xiaohongshu/Zhihu formatting. Unlike md2wechat-skill, easyGZH renders LOCALLY (no paid API), needs only the user's own appid/secret to publish, and remembers tone across runs."
---

# easyGZH — 公众号调性排版

easyGZH 把写好的文字稳定地排成符合公众号调性的版式并发布到草稿箱。它是一个**编译型 Go CLI**（`easygzh` 命令），由你（agent）指挥。核心差异化：

- **本地渲染**（goldmark + go-premailer CSS 内联），确定性，调性不漂移 —— 不依赖任何付费 API
- **本地发布**，只需用户自己的 appid/secret —— 不绑任何第三方付费 key
- **OKF 记忆库**记住每个号的视觉/结构调性 —— 这是它独有的、md2wechat 等竞品完全空白的能力

---

## 核心行为契约（统管全局）

**每个决策点都遵循：多路径 + 最高自动化优先 + 优雅降级 + 用户终决。**

1. **列路径** — 先想清楚当前环境有哪几条可行路径。
2. **最高自动化优先** — 默认尝试自动化程度最高的那条。
3. **逐级降级** — 走不通就降级，再不通再降级，直到跑通。
4. **用户终决** — 不替用户默认。给选项让他选；同时从最强的那条开始试。

原话精神：*"对于有多种实现可能的事，不要替用户默认，给出选项让他选；同时默认从最强的那个开始试，不通就降级，再不通再降级，直到跑通。"*

easygzh CLI 每个命令都输出结构化 JSON（`--json`，含 `code`/`data`/`next_actions`）——用这个做决策。

---

## CLI 与你的分工

```
你（agent，读这份 SKILL）           easygzh CLI（确定性执行）
   意图路由 ────────────────────────►  inspect（就绪检查）
   读 OK 记忆库定调性 ──────────────►
   决定模板/主题 ───────────────────►  convert / preview（渲染）
   AI 能力（改标题/润色/去AI味）        image（处理+上传，生成交给你）
   ← easygzh 只给 prompt 模板
   确认发布 ─────────────────────────►  publish（推草稿）
   反馈写回 OK 记忆库                  doctor（诊断）
```

**原则**：CLI 做确定性、有副作用、需凭证的事（渲染/上传/发布/诊断）；你做决策、记忆维护、AI 能力。

---

## Stage 0 — 探测（每次先做）

跑 `easygzh doctor --json` 探测环境：版本、平台、可用主题、微信凭证状态，以及记忆库是否存在、是否通过校验。据此决定后续路径。

确认 `easygzh` 在 PATH（用户已 `brew install` 或下载二进制）。若不在，引导安装。

---

## Stage 1 — 获取输入

**询问用户怎么输入，不替默认**（对装了 agent 的用户，输入是小事，把选择权给他）：
- 粘贴文本到对话
- 给本地 `.md` 文件路径
- 给一个 URL（若有 web 能力抓取）

拿到文字后，若非 Markdown，做最小结构化（识别标题/列表/段落），**保留作者原意，不擅自改写**。

---

## Stage 2 — 读取调性（稳定访问记忆库）

调性稳定的关键：**每次用确定路径读同一组文件**。

先跑 `easygzh memory profiles --json` 看有哪些号。读 `~/.easygzh/memory/profiles/<account>/` 下的：
- `visual-tone.md`（配色/字号/间距/分隔符/emoji）
- `structure-tone.md`（开头/小标题/结尾习惯）
- `current-theme.md`（引用主题 + 个性化 CSS 覆盖）

把 `current-theme.md` 里的 `base` + 覆盖 CSS 传给 `convert --css`，或在 `--theme` 指定。

### Stage 2b — 无 profile 时初始化

**遵循核心契约，列选项让用户选：**
- **路线 1（最高自动化）**：`easygzh memory init` 从二进制内嵌脚手架初始化并校验，再引导填关键问题。
- **路线 2（新增账号）**：`easygzh memory profile add <account>` 创建 profile，再逐项填写。
- **路线 3（从满意文章反推）**：用户贴 1-3 篇满意文章，你提炼调性生成 profile。

初始化或修改后运行 `easygzh memory validate --json`；校验不通过时不能把该 profile 当作稳定记忆使用。

---

## Stage 3 — 确定模板/主题

跑 `easygzh theme list --json` 看内置主题。**列出让用户选，默认用已有 profile 的 current-theme（最高自动化）：**

| 路线 | 做法 | 适用 |
|------|------|------|
| **A. 用已有 profile** | 读 `current-theme.md`（base + 覆盖 CSS） | 已有 profile（默认） |
| **B. 自建** | 对话问答（主色？字号？严肃/活泼？emoji？）→ 生成主题 CSS | 全新号 |
| **C. 从满意文章提炼** | 用户给 1-3 篇满意文章 → 你提取视觉特征 → 生成主题 | 有积累的号 |
| **D. AI 优化现有** | 读当前主题 + 用户反馈 → 你提议 diff → 确认 → 更新 | 迭代调性 |

额外选项：**无模板**（`convert --css ""`，最快但调性不保证）；**纯模板**（convert 直接渲染）；**模板 + AI 协作**（视觉层锁模板，AI 补过渡句/摘要）。

---

## Stage 4 — 渲染

```
easygzh convert <article.md> --theme <name> --json
# 或带个性化覆盖：
easygzh convert <article.md> --css "<覆盖CSS>" --json
```

渲染管线（纯本地，无网络）：Markdown → goldmark → go-premailer 内联 CSS → 微信安全后处理 → inline-styled HTML。确定性：同输入永远同输出（调性稳定的根基）。

`--no-footnotes` 保留外链不转脚注（默认转，因公众号正文外链不可点）。

---

## Stage 5 — 交付 / 发布（最高自动化优先降级）

```
1. preview（无副作用，先看效果）
   easygzh preview <article.md> --theme <name>
2. inspect（就绪检查：凭证/内容/封面齐备？）
   easygzh inspect <article.md> --json
3. publish（推草稿，需 WECHAT_APPID/SECRET）
   easygzh publish <article.md> --title "..." --cover cover.jpg --json
   # 只测到本地？用 --save-draft out.json（不推微信，生成草稿JSON）
4. 兜底：剪贴板/文件
   convert -o out.html，提示用户粘贴到公众号后台
```

**剪贴板/文件永远是兜底。** 浏览器自动填后台因反爬风险，仅作加分路径，不通就降级，不赌核心体验。

---

## Stage 6 — 反馈与记忆更新（长期伙伴）

交付后问满意度。满意则：
1. 把文章存为 `type: sample`（`~/.easygzh/memory/profiles/<account>/samples/YYYY-MM-DD-<slug>.md`）。
2. 用户表达的新偏好提炼成特征，更新 visual-tone/structure-tone。
3. 更新 log.md（`## YYYY-MM-DD` 倒序）。
4. 运行 `easygzh memory validate --json`，确认 frontmatter、链接和秘密扫描全部通过。

下次为同号排版，Stage 2 读到这些积累 → 越用越懂。

---

## AI 能力（你来做，easygzh 不内置 LLM）

这些 easygzh 不做，由你（agent）执行，必要时可让 easygzh 提供 prompt 模板：
- 改标题（`title suggest` 风格的头脑风暴）
- 去 AI 味（改写得更自然）
- 润色 / 补过渡句 / 生成摘要
- 配图建议（实际生成交给你的图像工具，easygzh 只做处理+上传）

---

## 提醒

- 渲染确定性 ≠ 调性死板。视觉层锁模板，你在内容层增量。
- 记忆库私有，不要把 profile 内容发外部服务。
- MVP 不做：小红书/知乎、代码块高亮、公式、newspic（无官方API）、内置 LLM。
- 不确定就问，别猜。
