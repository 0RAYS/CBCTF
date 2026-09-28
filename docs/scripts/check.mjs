import { existsSync, readFileSync, readdirSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = fileURLToPath(new URL('../docs/', import.meta.url));
const errors = [];
let pages = 0;
let links = 0;

function resolveTarget(file, target) {
  if (/^(?:[a-z]+:|#|\/\/)/i.test(target)) return;
  const pathname = decodeURIComponent(target.split(/[?#]/)[0]);
  if (!pathname) return;
  const absolute = pathname.startsWith('/');
  const base = path.resolve(absolute ? root : path.dirname(file), absolute ? pathname.slice(1) : pathname);
  const stem = base.replace(/\.(?:html|mdx?)$/, '');
  const candidates = [base, `${stem}.md`, `${stem}.mdx`, path.join(stem, 'index.md'), path.join(stem, 'index.mdx')];
  if (absolute) candidates.push(path.join(root, 'public', pathname));
  links++;
  if (!candidates.some(existsSync)) errors.push(`${path.relative(root, file)}: missing target ${target}`);
}

function walk(directory) {
  for (const entry of readdirSync(directory, { withFileTypes: true })) {
    const file = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      if (entry.name !== 'public') walk(file);
      continue;
    }
    if (entry.name === '_meta.json') {
      for (const item of JSON.parse(readFileSync(file, 'utf8'))) {
        if (typeof item === 'string') resolveTarget(file, item);
      }
    }
    if (!/\.mdx?$/.test(entry.name)) continue;
    pages++;
    const text = readFileSync(file, 'utf8');
    const frontmatter = text.match(/^---\r?\n([\s\S]*?)\r?\n---(?:\r?\n|$)/);
    for (const field of ['title', 'description']) {
      if (!frontmatter || !new RegExp(`^${field}:\\s*\\S`, 'm').test(frontmatter[1])) {
        errors.push(`${path.relative(root, file)}: missing initial frontmatter ${field}`);
      }
    }
    const prose = text.replace(/^```[^\n]*\n[\s\S]*?^```\s*$/gm, '');
    for (const match of prose.matchAll(/!?\[[^\]]*\]\(([^\s)]+)(?:\s+[^)]*)?\)/g)) resolveTarget(file, match[1]);
    for (const match of prose.matchAll(/\b(?:src|href)=["']([^"']+)["']/g)) resolveTarget(file, match[1]);
    for (const match of prose.matchAll(/^\[[^\]]+\]:\s+(\S+)/gm)) resolveTarget(file, match[1]);
  }
}

walk(root);
if (errors.length) {
  console.error(errors.join('\n'));
  process.exitCode = 1;
} else {
  console.log(`Checked ${pages} pages: frontmatter and ${links} local link/sidebar/asset targets passed.`);
}
