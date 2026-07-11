// Regenerates golden fixtures from Node juice. Run after editing themes.
import { readFileSync, writeFileSync, mkdirSync } from "node:fs";
import { resolve, dirname } from "node:path";
import { fileURLToPath } from "node:url";
import juice from "juice";
import MarkdownIt from "markdown-it";

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT = resolve(__dirname, "..");
const FIXTURES = resolve(ROOT, "testdata/fixtures");
const THEMES = resolve(ROOT, "themes");
const GOLDEN = resolve(ROOT, "testdata/golden");
mkdirSync(GOLDEN, { recursive: true });

const md = new MarkdownIt({ html: false, breaks: true, linkify: true, typographer: true });
const samples = ["sample1", "sample2"];
const themes = ["default", "lively"];

for (const sample of samples) {
  const markdown = readFileSync(resolve(FIXTURES, `${sample}.md`), "utf8");
  const html = `<section id="easygzh-root">${md.render(markdown)}</section>`;
  for (const theme of themes) {
    const css = readFileSync(resolve(THEMES, `${theme}.css`), "utf8");
    const inlined = juice.inlineContent(html, css, { preserveImportant: true, inlinePseudoElements: true });
    writeFileSync(resolve(GOLDEN, `${sample}@${theme}.html`), inlined);
    console.log(`wrote ${sample}@${theme} (${inlined.length} bytes)`);
  }
}
