# 微信公众号 HTML 限制清单

公众号图文编辑器对粘贴的 HTML 有严格限制。easyGZH 的 `scripts/lib/wechat-postprocess.mjs` 已经处理了下面这些，但理解它们有助于解释行为、排查问题、编写主题。

## 1. 样式必须内联

- `<style>` 块和 `<link>` 样式表会被**整段剥掉**。
- class 选择器即使保留也无用（定义它的 `<style>` 没了）。
- → 所以 easyGZH 用 `juice.inlineContent()` 把所有 CSS 内联成 `style="..."` 写到每个元素上。

## 2. 标签白名单

安全标签（实测可保留）：`p, span, div, section, img, a, ul, ol, li, strong, em, b, i, br, hr, table, thead, tbody, tr, td, th, pre, code, blockquote, h1-h6, sup, sub, mark`。

未知标签会被微信**连同子内容一起丢弃**（危险！）。easyGZH 的做法是：未知标签**解包**（保留子内容，去掉外壳），绝不丢内容。

## 3. `<script>` 完全移除

任何 `<script>`、`<iframe>`、`<object>`、`<embed>` 都会被删。easyGZH 在后处理阶段主动移除。

## 4. 嵌套深度

同一标签嵌套超过 14 层会被拒。常见元凶是深层 `<div>`/`<section>`。easyGZH 检测并拍平超深的同类嵌套。

## 5. 行内元素属性限制

某些 CSS 属性只对行内元素生效。例如 `font-size` 在裸 `<p>` 上可能失效，要设在 `<span>` 上更稳。easyGZH 主题里关键尺寸尽量直接命中元素。

## 6. 图片

- 必须是**绝对 URL**（http/https）或 data URI。相对路径无法解析。
- 微信会在粘贴时**重新上传**图片到自己的 CDN。
- 部分图床有防盗链（referer 校验），会显示裂图。easyGZH 给 `<img>` 加 `referrerpolicy="no-referrer"`，并优先建议用宽松图床（如微博图床、阿里云 OSS、腾讯云 COS、mdnice 免费图床）。
- **MVP 不做 AI 配图**，用户自贴图。

## 7. 外部链接

- 标准文章（非原创/未开通白名单）正文里的**外部链接不可点击**。
- mdnice 的成熟解法：把外链转成上标 `[n]` + 文末编号引用列表。
- easyGZH 的 `scripts/lib/link-footnote.mjs` 实现了这个（markdown-it 核心规则）。
- 锚点（`#xxx` 同文档跳转）不转，保留。
- 在公众号后台开通了"原创"或"微信支付白名单"的号可以让外链可点，但 MVP 默认走安全路线（转脚注）。

## 8. 列表标记重置

微信会**重置 `<ul>`/`<ol>` 的默认标记**，导致列表项前面的圆点/编号消失。easyGZH 在 `<li>` 上硬编码 `list-style: disc/decimal inside`，并把 `<ol>` 的序号手动算出来，确保列表可读。

## 9. 字体

公众号不支持自定义 web font（`@font-face` 失效）。只能用系统字体。中文长文推荐：苹方 / PingFang SC（iOS/macOS）、Microsoft YaHei（Windows）、系统无衬线兜底。

## 10. 暗色模式

微信有暗色模式，会自动给某些元素重新着色。若想某节点不被重着色，加 `data-no-dark` 属性（MVP 暂不处理，留作后续）。

---

**来源依据**：mdnice `src/utils/converter.js`、lyricat/wechat-format、公众号编辑器实测、微信开发者文档。这些限制是 easyGZH 后处理模块存在的根本原因。
