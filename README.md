# easyGZH

> 一个编译型 Go CLI + Agent Skill，把写好的文字**稳定地**排成符合公众号调性的版式，并发布到草稿箱。本地渲染、本地发布、无付费绑定、有记忆。

## 它解决什么

公众号排版工具有两类老问题：

1. **md2wechat-skill 这类**：渲染靠付费闭源 API（"48 主题"在服务端，OSS 代码是空壳），连纯发布都绑付费 key。
2. **纯 AI 生成**：每次风格漂移，号没有统一调性；mdnice 这类靠人选主题，费时易撞衫。

easyGZH 的做法：**本地渲染**（goldmark + go-premailer CSS 内联，确定性，开源）+ **本地发布**（用户自己的 appid/secret 即可，零付费绑定）+ **OKF 记忆库**（记住每个号的调性，跨次复现）。

## 核心特点

- **🖥️ 本地渲染，确定性**：Markdown → goldmark → go-premailer CSS 内联 → 微信安全 HTML。同输入永远同输出，调性不漂移。**不依赖任何付费 API**。
- **🔐 本地发布，零绑定**：只需用户自己的 `WECHAT_APPID`/`WECHAT_SECRET` 即可推草稿。access_token 用文件持久化（修复上游内存缓存烧配额的弱点）。
- **🧠 本地记忆库**：每个公众号一个 profile（视觉/结构/主题/样本），采用 OKF 风格类型并由 OpenKnowledge 校验，越用越懂这个号。
- **⚙️ 编译型单二进制**：Go 编译，11MB，`brew install` 或下载即用，无需 Node/Python 环境。
- **🤖 Agent 原生**：是 Agent Skill（SKILL.md），跑在你已有的 agent 里，由 agent 指挥 CLI。
- **📖 MIT 开源**：对比 md2wechat 的 BUSL-1.1。

## 与 md2wechat-skill 的对比

| 维度 | easyGZH | md2wechat-skill |
|------|---------|----------------|
| **渲染** | ✅ 本地（goldmark + go-premailer），开源 | ❌ 付费闭源 API（OSS 是空壳） |
| **主题** | ✅ 真 CSS，可本地创作 | 服务端别名，无法本地定制 |
| **发布付费绑定** | ✅ 无，纯 appid/secret 即可 | ❌ 绑 md2wechat API key |
| **access_token 缓存** | ✅ 文件持久化 | 内存，每次重取（烧配额） |
| **微信错误码** | ✅ 全表 | 仅 4 个 |
| **记忆/调性** | ✅ OKF 记忆库（核心差异化） | ❌ 无（Brand Profile 是静态文件） |
| **AI 模式** | 委托 agent（不内置 LLM） | 只发 prompt 给 agent 不执行 |
| **图像生成** | 只做处理+上传，生成交 agent | 7 provider 但只是发 prompt |
| **技术栈** | Go 单二进制 | Go 二进制 + npm 包装 |
| **许可证** | MIT | BUSL-1.1（商业受限） |
| **newspic/小绿书** | 不做（无官方 API） | 手写 HTTP（脆弱） |

一句话：**md2wechat 的渲染是付费黑盒，easyGZH 是开源本地；md2wechat 无记忆，easyGZH 记住你的号。**

---

> 🤖 **通过 Agent（ZCode / Claude Code 等）使用？** 先看 [Agent 快速开始](docs/agent-quickstart.md) —— 你的 Agent 会一步步带你上手：说清能干嘛、问你要不要装、装好带你排第一篇，全程不用读技术手册。

---

## 安装

### 方式 1：从源码编译（当前）

```bash
git clone https://github.com/Karlcptbtptp2/easygzh.git
cd easygzh
make build          # 产出 ./easygzh
./easygzh version
```

需要 Go ≥ 1.22。`make build` 会处理依赖。

### 方式 2：下载预编译二进制（发布后）

见 GitHub Releases（macOS/Linux/Windows × amd64/arm64）。

### 配置微信发布（可选，仅发布需要）

```bash
export WECHAT_APPID=你的appid
export WECHAT_SECRET=你的secret
# access_token 会持久化到 ~/.easygzh/token.json
```

## 使用

### 命令一览

```bash
easygzh convert <article.md> --theme default        # 渲染为公众号 HTML
easygzh preview <article.md> --theme lively          # 浏览器预览
easygzh inspect <article.md> --json                  # 就绪检查
easygzh publish <article.md> --title "标题" --json    # 推草稿（需凭证）
easygzh publish <article.md> --save-draft out.json   # 只生成本地草稿JSON（不推微信）
easygzh image <img.jpg> --upload                     # 处理并上传图片
easygzh memory init                                  # 初始化本地记忆库
easygzh memory status --json                         # 查看存在性、profiles 和健康状态
easygzh memory validate --json                       # 校验元数据、链接和秘密扫描
easygzh memory profile add <account>                 # 新增公众号 profile
easygzh theme list                                   # 列内置主题
easygzh doctor --json                                # 环境诊断
easygzh skills read easygzh                          # 输出 SKILL.md（供 agent 读）
```

所有命令支持 `--json` 输出结构化响应（含 `code`/`data`/`next_actions`），便于 agent 消费。

### 作为 Agent Skill 使用

把 `SKILL.md` 加到你的 agent（ZCode / Claude Code 等），然后对话：

```
你：用 easyGZH 帮我排版这篇（贴文字/给文件）
Agent：[doctor 探测] [读记忆库调性] [convert 渲染]
       已生成 inline HTML。满意吗？满意我存到样本库里。
       要发布到草稿箱吗？（你已配 WECHAT_APPID/SECRET）
```

第一次为某号排版会引导建立 profile（三条路线：自建 / 从满意文章提炼 / AI 优化）。

## 渲染管线（纯本地）

```
Markdown ──goldmark(+GFM)──▶ HTML ──go-premailer 内联CSS──▶ 微信安全后处理 ──▶ inline-styled HTML
```

- **goldmark**：CommonMark + 表格/删除线/linkify/任务列表
- **go-premailer**：CSS 内联到 `style=""`（goquery + Cascadia，支持 `#easygzh-root h1` 后代选择器）
- **微信后处理**：标签白名单、嵌套深度<15、图片绝对化、列表标记硬化、外链转脚注

确定性是调性稳定的技术根基。CSS 内联层经 golden fixture 测试与 Node `juice` 对齐验证。

## 目录结构

```
easyGZH/
├── cmd/easygzh/           # CLI（cobra）：convert/preview/publish/image/memory/theme/doctor/skills
├── internal/
│   ├── render/            # 渲染管线（markdown/inline/wechat/render）
│   ├── wechat/            # 发布（FileCache/client/errors/publish）— silenceper SDK
│   ├── image/             # 图像处理与上传（imaging）
│   ├── memory/            # OKF 记忆库读写
│   ├── theme/             # 主题加载
│   └── cli/               # JSON 响应契约
├── themes/                # 内置主题（default.css / lively.css，纯 CSS）
├── memory-scaffold/       # OKF 记忆库骨架
├── SKILL.md               # Agent 指挥手册
├── references/            # 渐进式披露文档
├── testdata/              # golden fixtures（juice vs go-premailer 对齐）
└── Makefile
```

## 技术栈

| 模块 | 库 |
|------|-----|
| Markdown | [yuin/goldmark](https://github.com/yuin/goldmark) |
| CSS 内联 | [vanng822/go-premailer](https://github.com/vanng822/go-premailer) |
| 微信 SDK | [silenceper/wechat/v2](https://github.com/silenceper/wechat) |
| 图像 | [disintegration/imaging](https://github.com/disintegration/imaging) |
| CLI | [spf13/cobra](https://github.com/spf13/cobra) |
| 记忆库 | [OKF 协议](https://github.com/GoogleCloudPlatform/knowledge-catalog/blob/main/okf/SPEC.md) + [OpenKnowledge](https://github.com/inkeep/open-knowledge) 结构 |

## License

MIT
