/**
 * markdown-it plugin: convert external links to footnote-style references.
 *
 * Why: standard WeChat public account articles (非原创/未开通微信支付的白名单)
 * cannot have clickable external links inside the body. mdnice solves this by
 * turning each <a href="https://...">text</a> into text followed by a superscript
 * reference [n], with a numbered link list appended at the end.
 *
 * This plugin does the same: rewrites the markdown-it "link_open" / "link_text"
 * / "link_close" rules so that external links become inline superscript refs,
 * and appends a "## References" section at the end of the document.
 *
 * Internal links (same-document anchors) are left untouched.
 */
import MarkdownIt from "markdown-it";

/**
 * Build a configured markdown-it instance with this plugin applied.
 * Pass { linkFootnotes: false } in options to disable.
 */
export function createRenderer(options = {}) {
  const md = new MarkdownIt({
    html: false, // WeChat sanitizes anyway; keep source clean
    breaks: true, // single line break → <br>, matches WeChat writing habits
    linkify: true,
    typographer: true,
  });

  if (options.linkFootnotes !== false) {
    applyLinkFootnotePlugin(md);
  }

  return md;
}

/**
 * Collect external links encountered during render and append a reference list.
 * Each external link renders as: <text><sup>[n]</sup>. The list is appended once
 * at the end of the document via a core rule.
 */
export function applyLinkFootnotePlugin(md) {
  md.core.ruler.push("easygzh_collect_links", (state) => {
    const Token = state.Token;
    let counter = 1;
    const links = []; // { n, text, href }
    const seen = new Map(); // href -> n (dedupe)

    // Link tokens live INSIDE inline tokens' .children, not at the top level.
    // Walk every inline token's children, find link_open ... link_close pairs,
    // stamp a data-ref on the link_open, and insert a <sup>[n]</sup> right after
    // the matching link_close.
    for (const tok of state.tokens) {
      if (tok.type !== "inline" || !tok.children) continue;
      const children = tok.children;
      for (let i = 0; i < children.length; i++) {
        const child = children[i];
        if (child.type !== "link_open") continue;
        const href = child.attrGet("href") || "";
        if (!isExternal(href)) continue;

        // Capture link text between link_open and link_close.
        let text = "";
        let j = i + 1;
        while (j < children.length && children[j].type !== "link_close") {
          if (children[j].type === "text") text += children[j].content;
          j++;
        }

        // Dedupe by href.
        let n;
        if (seen.has(href)) {
          n = seen.get(href);
        } else {
          n = counter++;
          seen.set(href, n);
          links.push({ n, text, href });
        }

        child.attrSet("data-ref", String(n));

        // Inject superscript marker right after link_close, if found.
        if (j < children.length && children[j].type === "link_close") {
          const sup = new Token("html_inline", "", 0);
          sup.content = `<sup style="font-size:0.75em;color:#888;">[${n}]</sup>`;
          sup.level = children[j].level;
          children.splice(j + 1, 0, sup);
          i = j + 1; // skip past the close + the inserted sup
        }
      }
    }

    // Append the reference list if any links were collected.
    if (links.length > 0) {
      const hr = new Token("hr", "hr", 0);
      hr.markup = "---";
      state.tokens.push(hr);

      const hOpen = new Token("heading_open", "h3", 1);
      hOpen.markup = "###";
      state.tokens.push(hOpen);
      const hInline = new Token("inline", "", 0);
      hInline.content = "References / 引用";
      hInline.children = [makeTextToken(Token, "References / 引用")];
      state.tokens.push(hInline);
      const hClose = new Token("heading_close", "h3", -1);
      hClose.markup = "###";
      state.tokens.push(hClose);

      for (const link of links) {
        state.tokens.push(new Token("paragraph_open", "p", 1));
        const inline = new Token("inline", "", 0);
        const body = `[${link.n}] ${link.text} ${link.href}`;
        inline.content = body;
        inline.children = [makeTextToken(Token, body)];
        state.tokens.push(inline);
        state.tokens.push(new Token("paragraph_close", "p", -1));
      }
    }
  });
}

function makeTextToken(Token, content) {
  const t = new Token("text", "", 0);
  t.content = content;
  return t;
}

function isExternal(href) {
  if (!href) return false;
  // Anchors inside the same document are not external.
  if (href.startsWith("#")) return false;
  // mailto: / tel: are not web links for our purposes.
  if (/^(mailto:|tel:)/i.test(href)) return false;
  // Protocol-relative and http/https are external.
  return /^(https?:)?\/\//i.test(href);
}
