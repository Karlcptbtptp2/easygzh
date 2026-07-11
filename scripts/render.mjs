#!/usr/bin/env node
/**
 * easyGZH render pipeline.
 *
 * Input  → inline-styled, WeChat-safe HTML.
 *
 * Pipeline (pure Node, no browser — mirrors mdnice's proven approach):
 *   Markdown ──markdown-it──▶ HTML
 *                             │
 *           theme CSS string ─┤
 *                             ▼
 *                       juice.inlineContent  → CSS inlined as style="..."
 *                             │
 *                             ▼
 *                       wechat-postprocess    → strip/unwrap/depth/images/lists
 *                             │
 *                             ▼
 *                       inline-styled HTML to stdout
 *
 * Usage:
 *   echo '<json>' | node scripts/render.mjs
 *   node scripts/render.mjs --markdown file.md --theme themes/default.css
 *   node scripts/render.mjs --test     # runs a built-in self-test
 *
 * JSON shape (stdin, one line):
 *   {
 *     "markdown": "# Hello\n\nWorld",
 *     "themeCss": "h1 { color: red } ...",   // OR
 *     "themePath": "themes/default.css",      // resolved relative to repo root
 *     "linkFootnotes": true                   // default true
 *   }
 */
import { readFileSync } from "node:fs";
import { resolve, dirname } from "node:path";
import { fileURLToPath } from "node:url";
import juice from "juice";
import { createRenderer } from "./lib/link-footnote.mjs";
import { toWeChatHtml } from "./lib/wechat-postprocess.mjs";

const __dirname = dirname(fileURLToPath(import.meta.url));
const REPO_ROOT = resolve(__dirname, "..");

// ---------- input helpers ----------

function readStdin() {
  return new Promise((res) => {
    let data = "";
    if (process.stdin.isTTY) return res("");
    process.stdin.setEncoding("utf8");
    process.stdin.on("data", (c) => (data += c));
    process.stdin.on("end", () => res(data));
  });
}

function parseArgs(argv) {
  const out = {};
  for (let i = 2; i < argv.length; i++) {
    const a = argv[i];
    if (a === "--test") out.test = true;
    else if (a === "--markdown") out.markdownPath = argv[++i];
    else if (a === "--theme") out.themePath = argv[++i];
    else if (a === "--no-footnotes") out.linkFootnotes = false;
    else if (a === "--help" || a === "-h") out.help = true;
  }
  return out;
}

function loadTheme(themePath, themeCss) {
  if (themeCss && themeCss.trim()) return themeCss;
  if (themePath) return readFileSync(resolve(REPO_ROOT, themePath), "utf8");
  // Default: ship a sane built-in so --test works without external theme.
  return DEFAULT_THEME;
}

// ---------- core pipeline ----------

export function render({ markdown, themeCss, themePath, linkFootnotes = true }) {
  const md = createRenderer({ linkFootnotes });
  const html = md.render(markdown);
  const css = loadTheme(themePath, themeCss);

  // Wrap in the #easygzh-root section BEFORE juicing, so the theme's
  // `#easygzh-root h1 { ... }` selectors actually match. juice inlines them
  // onto each element's style="".
  const scoped = `<section id="easygzh-root">${html}</section>`;
  const inlined = juice.inlineContent(scoped, css, {
    preserveImportant: true,
    inlinePseudoElements: true,
  });

  // Post-process for WeChat safety WITHOUT re-wrapping (already wrapped above).
  const wechatSafe = toWeChatHtml(inlined, { wrapSection: false });
  return wechatSafe;
}

// ---------- CLI ----------

async function main() {
  const args = parseArgs(process.argv);
  if (args.help) return void console.log(HELP);

  if (args.test) return runSelfTest();

  let markdown = "";
  let themeCss = "";
  let themePath = "";
  let linkFootnotes = true;

  const stdinData = (await readStdin()).trim();
  if (stdinData) {
    let payload;
    try {
      payload = JSON.parse(stdinData);
    } catch {
      // Treat raw stdin as plain markdown text (convenient for piping).
      payload = { markdown: stdinData };
    }
    markdown = payload.markdown ?? "";
    themeCss = payload.themeCss ?? "";
    themePath = payload.themePath ?? "";
    if (payload.linkFootnotes === false) linkFootnotes = false;
  }

  if (args.markdownPath) markdown = readFileSync(args.markdownPath, "utf8");
  if (args.themePath) themePath = args.themePath;
  if (args.linkFootnotes === false) linkFootnotes = false;

  if (!markdown.trim()) {
    process.stderr.write(
      "easyGZH: no markdown input. Pipe JSON or use --markdown FILE.\n"
    );
    process.exit(2);
  }

  const out = render({ markdown, themeCss, themePath, linkFootnotes });
  process.stdout.write(out);
  process.stdout.write("\n");
}

// ---------- self-test ----------

function runSelfTest() {
  const sample = `# 欢迎使用 easyGZH

这是一段**正文**，用来验证渲染管线。正文里有[一个外链](https://example.com)，它会被转成脚注引用。

## 二级标题

- 列表项一
- 列表项二

> 引用块内容。\u00a0另一些文字。

正文末尾。\u00a0`;
  const out = render({ markdown: sample, themeCss: DEFAULT_THEME });
  const checks = [
    ["has root section", out.includes('<section id="easygzh-root"')],
    ["h1 inlined", /<h1[^>]*style="[^"]*"/.test(out)],
    ["no <script>", !out.includes("<script")],
    ["no <style>", !out.includes("<style")],
    ["external link → footnote ref", out.includes("[1]")],
    ["references section appended", out.includes("References")],
  ];
  let ok = true;
  for (const [name, pass] of checks) {
    console.log(`${pass ? "PASS" : "FAIL"}  ${name}`);
    if (!pass) ok = false;
  }
  if (!ok) {
    process.stderr.write("\nSelf-test FAILED. Dumping output:\n" + out + "\n");
    process.exit(1);
  }
  console.log("\nSelf-test passed. Inline HTML fragment written to stdout.");
  process.stdout.write(out + "\n");
}

const HELP = `easyGZH render — Markdown + theme CSS → WeChat-safe inline HTML.

Usage:
  echo '{"markdown":"# hi","themePath":"themes/default.css"}' | node scripts/render.mjs
  node scripts/render.mjs --markdown article.md --theme themes/default.css
  node scripts/render.mjs --test

Options:
  --markdown FILE   read markdown from a file
  --theme FILE      read theme CSS from a file (relative to repo root)
  --no-footnotes    keep external links as-is instead of converting to refs
  --test            run the built-in self-test
  --help            show this help
`;

// A minimal theme so the script works with zero external files.
// Real themes live in themes/*.css and are much richer.
const DEFAULT_THEME = `
#easygzh-root { font-family: -apple-system, "PingFang SC", "Helvetica Neue", sans-serif; font-size: 15px; line-height: 1.75; color: #333; word-break: break-word; }
#easygzh-root h1, #easygzh-root h2, #easygzh-root h3 { color: #1a73e8; line-height: 1.4; margin: 1.2em 0 0.6em; }
#easygzh-root h1 { font-size: 22px; border-bottom: 2px solid #1a73e8; padding-bottom: 0.3em; }
#easygzh-root h2 { font-size: 18px; }
#easygzh-root h3 { font-size: 16px; }
#easygzh-root p { margin: 0 0 15px; }
#easygzh-root strong { color: #1a73e8; }
#easygzh-root a { color: #1a73e8; text-decoration: none; border-bottom: 1px solid #ddd; }
#easygzh-root blockquote { border-left: 3px solid #1a73e8; padding: 4px 15px; margin: 15px 0; color: #666; background: #f6f8fa; }
#easygzh-root ul, #easygzh-root ol { margin: 0 0 15px; padding-left: 22px; }
#easygzh-root li { margin: 4px 0; }
#easygzh-root hr { border: none; border-top: 1px dashed #ccc; margin: 25px auto; width: 60%; }
#easygzh-root code { background: #f6f8fa; padding: 2px 5px; border-radius: 3px; font-size: 0.9em; }
#easygzh-root pre { background: #f6f8fa; padding: 12px; border-radius: 5px; overflow-x: auto; }
#easygzh-root img { max-width: 100%; height: auto; display: block; margin: 15px auto; }
`;

// Only run the CLI entrypoint when invoked directly, not when imported as a
// module (e.g. by tests). This guard prevents stdin from blocking on import.
const isMainEntry = import.meta.url === `file://${process.argv[1]}`;
if (isMainEntry) {
  main().catch((err) => {
    process.stderr.write("easyGZH render failed: " + (err?.stack || err) + "\n");
    process.exit(1);
  });
}
