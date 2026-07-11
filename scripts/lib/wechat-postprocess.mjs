/**
 * WeChat public account (公众号) HTML post-processing.
 *
 * The WeChat editor strips <style> blocks and class selectors, drops unknown
 * tags, forbids <script>, resets list markers, and caps same-tag nesting depth
 * at <15. juice already inlines CSS onto style="" attributes; this module makes
 * the remaining tree WeChat-safe. See references/wechat-constraints.md.
 *
 * Pure, synchronous, cheerio-based. No browser.
 */
import { load } from "cheerio";

// Tags the WeChat rich-text editor accepts. Anything else is unwrapped (its
// children kept) rather than dropped, so we never lose content.
const ALLOWED_TAGS = new Set([
  "p", "span", "div", "section", "img", "a",
  "ul", "ol", "li", "strong", "em", "b", "i", "br", "hr",
  "table", "thead", "tbody", "tr", "td", "th",
  "pre", "code", "blockquote",
  "h1", "h2", "h3", "h4", "h5", "h6",
  "sup", "sub", "mark",
]);

const NESTING_DEPTH_LIMIT = 15;

/**
 * Process parsed HTML for WeChat safety. Mutates and returns the cheerio $.
 * Exported for testing the individual transforms.
 */
export function postProcess($, options = {}) {
  stripDangerous($);
  unwrapUnknownTags($);
  enforceNestingDepth($);
  absolutizeImages($);
  hardenLists($);
  if (options.wrapSection !== false) wrapTopLevelInSection($);
  return $;
}

/** Remove <script>, <style>, <link>, <iframe>, <object>, <embed> entirely. */
function stripDangerous($) {
  $("script, style, link, iframe, object, embed, noscript").remove();
}

/** Replace disallowed tags with their children (keep content, lose the wrapper). */
function unwrapUnknownTags($) {
  // Iterate a bounded number of passes. Each pass unwraps disallowed tags by
  // replacing them with their own children. We do NOT use an unbounded
  // fixpoint loop because cheerio's unwrap() can leave transient state that
  // makes a "changed" flag stay true forever on some node shapes.
  for (let pass = 0; pass < 8; pass++) {
    let touched = false;
    // Materialize the list first; mutating during $("*").each is unsafe.
    const targets = $("*")
      .toArray()
      .filter((el) => el.type === "tag" && !ALLOWED_TAGS.has(el.tagName));
    if (targets.length === 0) break;
    for (const el of targets) {
      const $el = $(el);
      const kids = $el.contents();
      if (kids.length === 0) {
        $el.remove();
      } else {
        $el.replaceWith(kids);
      }
      touched = true;
    }
    if (!touched) break;
  }
}

/** Cap same-tag nesting depth at NESTING_DEPTH_LIMIT by flattening inner copies. */
function enforceNestingDepth($) {
  // WeChat rejects articles with >14 levels of the same tag nested. We collapse
  // deep <div>/<section> chains (the usual culprits) by unwrapping when the
  // ancestor chain of the same tag exceeds the limit.
  $("*").each((_, el) => {
    if (el.type !== "tag") return;
    const tag = el.tagName;
    let depth = 0;
    let ancestor = el.parent;
    while (ancestor && ancestor.type === "tag") {
      if (ancestor.tagName === tag) depth++;
      if (depth >= NESTING_DEPTH_LIMIT) {
        $(el).contents().unwrap();
        return;
      }
      ancestor = ancestor.parent;
    }
  });
}

/** Ensure <img> src is an absolute URL (WeChat re-hosts on paste). */
function absolutizeImages($) {
  $("img").each((_, img) => {
    let src = $(img).attr("src") || "";
    src = src.trim();
    if (!src) {
      $(img).remove();
      return;
    }
    // WeChat needs absolute URLs; relative paths can't resolve in its editor.
    if (src.startsWith("//")) src = "https:" + src;
    if (!/^(https?:|data:)/i.test(src)) {
      // Leave a visible placeholder rather than shipping a broken relative URL.
      $(img).attr("src", "");
      $(img).attr("data-bad-src", src);
    } else {
      $(img).attr("src", src);
    }
    // Block referer leakage on hosts that hotlink-protect.
    $(img).attr("referrerpolicy", "no-referrer");
  });
}

/**
 * Re-apply list marker styling inline because WeChat resets <ul>/<ol> default
 * markers. We add a stable visible marker so lists read correctly without CSS.
 */
function hardenLists($) {
  // Number <ol> manually and give <ul> a bullet via CSS on <li>.
  $("ol").each((_, ol) => {
    let n = 1;
    const start = parseInt($(ol).attr("start") || "1", 10);
    n = isNaN(start) ? 1 : start;
    $(ol)
      .find("> li")
      .each((__, li) => {
        const existing = $(li).attr("style") || "";
        $(li).attr("style", `${existing} list-style: decimal inside;`.trim());
        n++;
      });
  });
  $("ul").each((_, ul) => {
    $(ul)
      .find("> li")
      .each((__, li) => {
        const existing = $(li).attr("style") || "";
        $(li).attr("style", `${existing} list-style: disc inside;`.trim());
      });
  });
}

/**
 * Wrap the whole document body in a single <section>. WeChat's editor treats the
 * pasted blob as one block; a single root section gives theme CSS a stable hook
 * and prevents fragment drift on paste.
 */
function wrapTopLevelInSection($) {
  const section = $('<section id="easygzh-root"></section>');
  $("body").children().appendTo(section);
  $("body").append(section);
}

/**
 * Full convenience pipeline: take an already-inlined HTML string, return a
 * WeChat-safe inline-styled HTML string (just the inner fragment, no <html>).
 */
export function toWeChatHtml(inlinedHtml, options = {}) {
  const $ = load(`<div id="__root">${inlinedHtml}</div>`);
  const root = $("#__root");
  stripDangerous($);
  unwrapUnknownTags($);
  enforceNestingDepth($);
  absolutizeImages($);
  hardenLists($);
  // Return the fragment, wrapped in the root section.
  const inner = root.html() || "";
  if (options.wrapSection === false) return inner;
  return `<section id="easygzh-root">${inner}</section>`;
}
