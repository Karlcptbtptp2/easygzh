# 公众号内联样式指南（AI 作者参考）

AI 直接写公众号 HTML 时，所有样式都要内联到 `style="..."` 上。这份指南给出经过验证的经验值和常用组件的写法。

## 1. 基础参数（中文长文舒适区）

| 参数 | 推荐值 | 说明 |
|------|--------|------|
| 正文字号 | 15px | 公众号长文最舒适 |
| 行高 | 1.75 | 中文阅读舒适区 |
| 字间距 | 0.05em | 微微透气 |
| 正文色 | #333 | 不用纯黑 #000 |
| 辅助色 | #666 / #888 | 导语、引用、注释 |
| 字体栈 | -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Microsoft YaHei', sans-serif | 系统无衬线 |

**正文段落示例**：
```html
<p style="font-family:-apple-system,BlinkMacSystemFont,'PingFang SC','Microsoft YaHei',sans-serif;font-size:15px;line-height:1.75;color:#333;letter-spacing:0.05em;margin:0 0 15px;padding:0 18px;">
  正文内容
</p>
```

## 2. 标题层级

| 级别 | 字号 | 样式建议 |
|------|------|---------|
| H1（文章标题） | 22px | 居中，底部加 2px 色线 |
| H2（小标题） | 18px | 左侧 4px 色条，节奏控制器 |
| H3 | 16px | 纯色字 |

**H1 示例**：
```html
<h1 style="font-size:22px;color:#2D2D2D;text-align:center;border-bottom:2px solid #7A8B6E;padding-bottom:0.4em;margin:0 0 30px;">
  文章标题
</h1>
```

**H2 示例（小标题 = 节奏控制器）**：
```html
<h2 style="font-size:18px;color:#2D2D2D;border-left:4px solid #7A8B6E;padding-left:12px;margin:40px 0 20px;">
  像买彩票一样
</h2>
```

## 3. 视觉锚点组件

### 导语区

标题下方、正文上方。字号 14px，颜色 #888。建立"这是我们共同的事"的联结感。

```html
<p style="font-size:14px;color:#888;line-height:1.8;margin:0 0 32px;padding:0 18px;">
  这个月，我从一次偶然的散步开始，重新学会了慢下来。
</p>
```

### 金句引用块

核心观点、情绪高潮、转折处。左侧 3px 色线 + 浅灰底 + 斜体。

```html
<section style="margin:28px 0;padding:18px 18px 18px 22px;background:#f7f7f7;border-left:3px solid #7A8B6E;margin-left:18px;margin-right:18px;">
  <p style="font-size:15px;color:#555;line-height:1.8;font-style:italic;margin:0;">
    "真正的探索，不在于抵达远方，而在于重新看见眼前。"
  </p>
</section>
```

### 分割线

克制使用，一篇 1-2 条。居中虚线或实线。

```html
<!-- 细实线 -->
<p style="margin:24px 18px;padding:0;border-top:1px solid #D5CFC6;font-size:0;line-height:0;">&nbsp;</p>

<!-- 虚线 -->
<hr style="border:none;border-top:1px dashed #ccc;margin:28px auto;width:50%;" />
```

### 留白（最重要的视觉元素）

连续文字超过 4 行必须打断。情绪转折处、观点切换处，用空行制造"顿号"。

```html
<!-- 情绪留白 -->
<section style="margin:24px 0;">&nbsp;</section>
```

### 提示卡片

```html
<section style="margin:28px 18px;padding:20px;background:#f7f7f7;border-radius:8px;">
  <p style="font-size:18px;margin:0 0 8px;">💡</p>
  <p style="font-size:15px;color:#555;line-height:1.8;margin:0;">
    尝试在今天找到三个你从未注意过的细节。
  </p>
</section>
```

### 胶囊按钮（CTA）

```html
<section style="text-align:center;margin:30px 18px;">
  <span style="display:inline-block;background:#E8923C;color:#fff;font-size:14px;letter-spacing:1px;padding:10px 36px;border-radius:24px;">
    立即加入
  </span>
</section>
```

### 气泡标签

```html
<section style="margin:12px 18px;">
  <span style="display:inline-block;background:#E8F0E8;color:#5A7D5A;font-size:14px;padding:8px 18px;border-radius:16px;">
    📖 总是抽不出时间阅读？
  </span>
</section>
```

## 4. 配色经验

### 莫兰迪绿系（适合个人成长/复盘/情绪随笔）
- 主色：#7A8B6E（橄榄绿）
- 强调：#B8956A（金棕）
- 背景锚点：#F5F3EF（暖米）
- 正文：#3A3A3A

### 暖橙系（适合生活/社群/轻松内容）
- 主色：#E8923C（暖橙）
- 气泡绿：#E8F0E8 / #5A7D5A
- 气泡粉：#F5E0D6 / #9C6B4F
- 背景：#FFF5EF

### 蓝灰系（适合技术/正式/公告）
- 主色：#1a73e8
- 背景：#f6f8fa
- 正文：#333

## 5. 配图样式

```html
<!-- 标题下方封面图 -->
<img src="https://..." referrerpolicy="no-referrer"
  style="max-width:100%;border-radius:8px;margin:20px auto;display:block;" />

<!-- 文中情绪配图 -->
<img src="https://..." referrerpolicy="no-referrer"
  style="max-width:100%;border-radius:8px;margin:32px auto;display:block;" />
```

## 6. 写作原则

- **保持克制**：公众号长文审美 = 留白 + 层级清晰，不是花哨
- **字号别太大**：正文 15px、H1 22px 左右是舒适区
- **行高 1.75**：最适合中文长文阅读
- **少用阴影/渐变**：除非内容定位是活泼/生活类
- **不是纯白底**：用浅色背景锚点（#f7f7f7 / #F5F3EF）制造视觉变化，避免纯文字白底
- **每 3-4 段插锚点**：金句/留白/小标题，控制读者滑动节奏

## 7. visual-tone.md → 内联样式映射

用户的 `visual-tone.md` 里写的自然语言规范，AI 读完后直接翻译成内联样式：

| visual-tone 字段 | 对应内联样式 |
|---|---|
| 主色 #7A8B6E | `style="color:#7A8B6E"` 加到标题/强调 |
| 正文字号 15px | `style="font-size:15px"` 加到 `<p>` |
| 行高 1.75 | `style="line-height:1.75"` 加到正文容器 |
| 段间距 15px | `style="margin:0 0 15px"` 加到 `<p>` |
| H1 居中带下划线 | `style="text-align:center;border-bottom:2px solid ..."` |
| 引用块浅灰底 | `style="background:#f7f7f7;border-left:3px solid ..."` |
