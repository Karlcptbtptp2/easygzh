#!/usr/bin/env node
/**
 * easyGZH memory store initializer.
 *
 * Copies the memory-scaffold/ from the repo to a target location (default
 * ~/.easygzh/memory), creating a fresh OK-style + OK-conformant memory store ready
 * for the skill to read/write.
 *
 * Usage:
 *   node scripts/init-memory.mjs                    # → ~/.easyGZH/memory
 *   node scripts/init-memory.mjs --target ~/my-path # → custom location
 *   node scripts/init-memory.mjs --account my-account  # scaffold a new profile
 *   node scripts/init-memory.mjs --force            # overwrite if exists
 *
 * Corresponds to SKILL Stage 2b, Route 1 (highest automation).
 */
import {
  existsSync,
  cpSync,
  mkdirSync,
  readFileSync,
  writeFileSync,
  readdirSync,
  renameSync,
  statSync,
} from "node:fs";
import { resolve, dirname, join } from "node:path";
import { homedir } from "node:os";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const REPO_ROOT = resolve(__dirname, "..");
const SCAFFOLD = join(REPO_ROOT, "memory-scaffold");
const DEFAULT_TARGET = join(homedir(), ".easygzh", "memory");

function parseArgs(argv) {
  const out = { target: DEFAULT_TARGET };
  for (let i = 2; i < argv.length; i++) {
    const a = argv[i];
    if (a === "--target") out.target = resolve(argv[++i]);
    else if (a === "--account") out.account = argv[++i];
    else if (a === "--force") out.force = true;
    else if (a === "--help" || a === "-h") out.help = true;
  }
  return out;
}

function copyScaffold(target, force) {
  if (existsSync(target) && !force) {
    throw new Error(
      `Target already exists: ${target}\nUse --force to overwrite, or pick another --target.`
    );
  }
  if (existsSync(target) && force) {
    const stamp = new Date().toISOString().replace(/[-:]/g, "").replace(/\.\d{3}Z$/, "Z");
    let backup = `${target}.backup-${stamp}`;
    let suffix = 2;
    while (existsSync(backup)) backup = `${target}.backup-${stamp}-${suffix++}`;
    renameSync(target, backup);
    console.log(`✓ Preserved existing store at: ${backup}`);
  }
  // cpSync with recursive; create parent first.
  mkdirSync(dirname(target), { recursive: true });
  cpSync(SCAFFOLD, target, { recursive: true, force: false, errorOnExist: true });
}

/**
 * Scaffold a new profile directory under an existing memory store, by copying
 * the example main-account profile. Used when the store exists and the user
 * wants to add another account without re-running full init.
 */
function addProfile(target, accountName) {
  const store = resolve(target);
  if (!existsSync(store)) {
    throw new Error(`Memory store not found at ${store}. Run init without --account first.`);
  }
  if (!/^[a-z0-9-]+$/.test(accountName)) {
    throw new Error(`Account name must be lowercase-kebab (a-z, 0-9, -). Got: ${accountName}`);
  }
  const src = join(store, "profiles", "main-account");
  const dst = join(store, "profiles", accountName);
  if (!existsSync(src)) {
    throw new Error(`Template profile not found: ${src}. The scaffold may be incomplete.`);
  }
  if (existsSync(dst)) {
    throw new Error(`Profile already exists: ${dst}. Choose another name or remove it first.`);
  }
  cpSync(src, dst, { recursive: true });
  // Stamp the account field in each frontmatter of the new profile.
  stampAccount(dst, accountName);
  appendProfileIndex(store, accountName);
  prependLogEntry(store, accountName);
  return dst;
}

function stampAccount(profileDir, accountName) {
  for (const entry of readdirSync(profileDir)) {
    const p = join(profileDir, entry);
    if (statSync(p).isDirectory()) {
      stampAccount(p, accountName);
      continue;
    }
    if (!entry.endsWith(".md")) continue;
    let txt = readFileSync(p, "utf8");
    txt = txt.replaceAll("main-account", accountName).replaceAll("主号", accountName);
    writeFileSync(p, txt);
  }
}

function appendProfileIndex(store, accountName) {
  const path = join(store, "profiles", "index.md");
  let txt = readFileSync(path, "utf8");
  if (txt.includes(`](${accountName}/index.md)`)) return;
  txt = txt.trimEnd() + `\n- [${accountName}](${accountName}/index.md) — 由初始化工具创建\n`;
  writeFileSync(path, txt);
}

function prependLogEntry(store, accountName) {
  const path = join(store, "log.md");
  let txt = readFileSync(path, "utf8");
  const marker = "# 变更日志\n";
  const date = new Date().toISOString().slice(0, 10);
  if (!txt.includes(marker)) throw new Error("log.md is missing '# 变更日志' heading");
  txt = txt.replace(marker, `${marker}\n## ${date}\n\n**Profile** — 创建 \`${accountName}\` profile。\n`);
  writeFileSync(path, txt);
}

function main() {
  const args = parseArgs(process.argv);
  if (args.help) return void console.log(HELP);

  if (args.account) {
    const dst = addProfile(args.target, args.account);
    console.log(`✓ Created profile: ${dst}`);
    console.log(`  Edit the files there to define this account's tone.`);
    return;
  }

  copyScaffold(args.target, args.force);
  console.log(`✓ easyGZH memory store initialized at:`);
  console.log(`  ${args.target}`);
  console.log(`\nNext steps:`);
  console.log(`  - Edit profiles/main-account/*.md to match your account.`);
  console.log(`  - Or add another account: node scripts/init-memory.mjs --account <name> --target ${args.target}`);
  console.log(`  - Point the skill at this path (or keep the default ~/.easyGZH/memory).`);
}

const HELP = `easyGZH memory initializer.

Usage:
  node scripts/init-memory.mjs                              init ~/.easygzh/memory
  node scripts/init-memory.mjs --target ~/path              init at custom path
  node scripts/init-memory.mjs --account <name>             add a profile to existing store
  node scripts/init-memory.mjs --force                      overwrite existing target

The scaffold is OKF + OpenKnowledge conformant. See references/memory-schema.md.
`;

try {
  main();
} catch (err) {
  process.stderr.write("init-memory failed: " + (err?.message || err) + "\n");
  process.exit(1);
}
