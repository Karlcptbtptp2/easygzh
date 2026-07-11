# 主题编写指南

主题就是一个 **CSS 字符串**。easyGZH 把它连同渲染出的 HTML 一起喂给 `juice.inlineContent()`，juice 把每条规则内联到匹配元素的 `style=""` 上。

## 1. 主题文件位置

- **内置主题**（公开、共享）：仓库的 `themes/*.css`，如 `default.css`、`lively.css`。
- **私有主题**（个人化）：记忆库的 `themes/<name>.md`（`type: theme`），或在某个 profile 的 `current-theme.md` 里写覆盖 CSS。

## 2. 关键约定：所有规则用 `#easygzh-root` 作用域

`render.mjs` 会把整个文档片段包在 `<section id="easygzh-root">` 里。所以**所有主题规则都应以 `#easygzh-root` 开头**：

```css
#easygzh-root { font-size: 15px; color: #333; }
#easygzh-root h1 { color: #1a73e8; }
#easygzh-root blockquote { border-left: 3px solid #1a73e8; }
```

这样既能提高特异性（覆盖浏览器默认），也避免污染。juice 会带着这个作用域选择器正确匹配并内联。

## 3. 可选的元素

主题里可以对以下元素定义样式（markdown-it 会生成这些标签）：

```
h1 h2 h3 h4 h5 h6       标题
p                        段落
strong em b i            强调
a                        链接（外链会被转成脚注，见下）
blockquote               引用块
ul ol li                 列表
hr                       分隔线（---）
code pre                 行内代码 / 代码块
img                      图片
table thead tbody tr th td   表格
sup                      脚注上标 [n]
```

## 4. base + 覆盖 模式

profile 的 `current-theme.md` 可以这样写，先继承内置主题再覆盖：

```markdown
---
type: current-theme
base: default
---
```css
/* 只写要改的部分 */
#easygzh-root h1 { color: #00838f; }
```

渲染时 easyGZH 先加载 `themes/default.css`，再拼接这里的 CSS。后者优先级更高（CSS 后定义覆盖前定义），juice 合并内联。

## 5. 视觉调性 → CSS 的翻译对照

`visual-tone.md` 里写的自然语言规范，应能机械翻译成 CSS：

| visual-tone 字段 | 对应 CSS |
|---|---|
| 主色 | `#easygzh-root h1,h2,h3,strong { color: ... }` |
| 正文字号/行高 | `#easygzh-root { font-size; line-height }` |
| 段间距 | `#easygzh-root p { margin-bottom }` |
| H1 样式 | `#easygzh-root h1 { ... }` |
| 分隔符样式 | `#easygzh-root hr { ... }` 或 `hr::before` |
| emoji 风格 | 不进 CSS，进 structure-tone |

## 6. 编写建议

- **保持克制**。公众号长文审美 = 留白 + 层级清晰，不是花哨。
- **字号别太大**。正文 15px、H1 22px 左右是舒适区。
- **行高 1.7-1.8** 最适合中文长文阅读。
- **少用阴影/渐变**，除非定位是活泼/生活类号（参考 lively.css）。
- **测一下**：写完主题，用 `scripts/render.mjs --test` 或实际渲染一篇样本文章，复制到公众号后台看效果。

## 7. 三条主题来源路线（与 SKILL Stage 3 对应）

- **A. 自建**：对话问答收集偏好 → 直接生成 CSS。
- **B. 从满意文章提炼**：贴一篇满意文章 → AI 反推它的视觉特征 → 生成 CSS。
- **C. AI 优化现有**：读当前 CSS + 用户反馈 → 提议 diff → 确认 → 更新。

三条路线产出的都是上面这种 CSS 字符串，最终都进 `juice.inlineContent`。
