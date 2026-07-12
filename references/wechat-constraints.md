# 微信公众号 HTML 约束清单（AI 作者必读）

你（AI）直接写公众号 HTML，这些是你**必须遵守的硬性约束**。违反任何一条，排版在公众号后台都会出问题。每次写完 HTML 后，对照这份清单自检。

## 1. 样式必须内联

`<style>` 块和 `<link>` 样式表会被微信**整段剥掉**。class 名也没用（定义它的 `<style>` 没了）。

**做法**：每个元素的样式直接写在 `style="..."` 上，不依赖任何外部 CSS 或 class。

```html
<!-- ✅ 正确 -->
<p style="font-size:15px;line-height:1.75;color:#333;margin:0 0 15px;">正文</p>

<!-- ❌ 错误：style 块会被剥掉 -->
<style>.lead { font-size:15px; }</style>
<p class="lead">正文</p>
```

## 2. 标签白名单

安全标签（实测可保留）：`p, span, div, section, img, a, ul, ol, li, strong, em, b, i, br, hr, table, thead, tbody, tr, td, th, pre, code, blockquote, h1-h6, sup, sub, mark`。

未知标签会被微信**连同子内容一起丢弃**。如果你用了白名单外的标签（如 `article`, `aside`, `figure`），内容会消失。

**做法**：只用白名单标签。需要分组时用 `section` 或 `div`。

## 3. 禁止 script / iframe

`<script>`、`<iframe>`、`<object>`、`<embed>` 会被删除。不要用。

## 4. 嵌套深度

同一标签嵌套超过 14 层会被拒。常见元凶是深层 `<div>`/`<section>` 嵌套。

**做法**：保持扁平结构。组件用 1-2 层 `section` 包裹就够，不要无限嵌套。

## 5. 行内元素属性限制

某些 CSS 属性只对行内元素生效。例如 `font-size` 在裸 `<p>` 上可能失效，设在 `<span>` 上更稳。

**做法**：关键字号、颜色，尽量用 `<span style="...">` 包裹文字，而不是只靠 `<p>` 的 style。

## 6. 图片

- 必须是**绝对 URL**（http/https 开头）或 data URI。相对路径（`images/xxx.jpg`）无法解析
- 微信会在粘贴时**重新上传**图片到自己的 CDN
- 部分图床有防盗链，会显示裂图
- 加 `referrerpolicy="no-referrer"` 可绕过部分防盗链

**做法**：图片用绝对 URL。优先用宽松图床（微博图床、阿里云 OSS、腾讯云 COS）。如果用户给的是本地图片，发布前需要先上传到某个可访问的图床。

```html
<img src="https://mmbiz.qpic.cn/..." referrerpolicy="no-referrer" style="max-width:100%;border-radius:8px;" />
```

## 7. 外部链接 → 脚注

标准文章正文里的**外部链接不可点击**（除非号开了"原创"白名单）。

**做法**：把外链转成正文中的上标 `[n]` + 文末编号引用列表：

```html
<!-- 正文中 -->
关于这个概念<sup style="font-size:0.75em;color:#888;">[1]</sup>的引用

<!-- 文末 -->
<hr style="border:none;border-top:1px dashed #ccc;margin:28px auto;width:50%;" />
<h3 style="font-size:16px;">引用</h3>
<p style="font-size:14px;color:#888;">[1] 链接文字 https://example.com</p>
```

锚点链接（`#xxx` 同文档跳转）不用转。

## 8. 列表标记重置

微信会**重置 `<ul>`/`<ol>` 的默认标记**，导致圆点/编号消失。

**做法**：在 `<li>` 上硬编码 `list-style`：

```html
<ul style="padding-left:22px;">
  <li style="list-style: disc inside;">列表项</li>
</ul>

<ol style="padding-left:22px;">
  <li style="list-style: decimal inside;">第一项</li>
  <li style="list-style: decimal inside;">第二项</li>
</ol>
```

## 9. 字体

公众号不支持 `@font-face` / web font。只能用系统字体。

**做法**：字体栈用系统无衬线：
```html
style="font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Microsoft YaHei', sans-serif;"
```

## 10. 暗色模式

微信有暗色模式，会自动给某些元素重新着色。如果不想某节点被重着色，加 `data-no-dark` 属性（高级用法，一般不用管）。

---

**来源依据**：mdnice `src/utils/converter.js`、lyricat/wechat-format、公众号编辑器实测。这些约束是 AI 直接写公众号 HTML 时必须满足的规范。
