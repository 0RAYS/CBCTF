import { existsSync, readFileSync, readdirSync } from 'node:fs';
import { gzipSync, gunzipSync } from 'node:zlib';

const dist = new URL('../dist/', import.meta.url);
const manifest = JSON.parse(readFileSync(new URL('.vite/manifest.json', dist), 'utf8'));
const entry = Object.keys(manifest).find((key) => manifest[key].isEntry);
// Validate every generated representation, including lazy vendors and worker assets.
const compressedFiles = readdirSync(dist, { recursive: true }).filter((file) => file.endsWith('.gz'));
for (const file of compressedFiles) {
  const compressed = readFileSync(new URL(file.replaceAll('\\', '/'), dist));
  const original = readFileSync(new URL(file.slice(0, -3).replaceAll('\\', '/'), dist));
  if (!gunzipSync(compressed).equals(original)) throw new Error(`Invalid precompressed asset: ${file}`);
}
const visited = new Set();
const files = new Set();

function visit(key) {
  if (visited.has(key)) return;
  visited.add(key);
  const chunk = manifest[key];
  files.add(chunk.file);
  for (const css of chunk.css ?? []) files.add(css);
  for (const dependency of chunk.imports ?? []) visit(dependency);
}

visit(entry);
const rows = [...files].map((file) => {
  const content = readFileSync(new URL(file, dist));
  const compressedFile = new URL(`${file}.gz`, dist);
  const compressed = existsSync(compressedFile) ? readFileSync(compressedFile) : null;
  if (compressed && !gunzipSync(compressed).equals(content)) throw new Error(`Invalid precompressed asset: ${file}`);
  return {
    file,
    bytes: content.length,
    gzip: gzipSync(content).length,
    servedWithGzip: compressed?.length ?? content.length,
  };
});
const totals = rows.reduce((sum, row) => ({
  bytes: sum.bytes + row.bytes,
  gzip: sum.gzip + row.gzip,
  servedWithGzip: sum.servedWithGzip + row.servedWithGzip,
}), {
  bytes: 0,
  gzip: 0,
  servedWithGzip: 0,
});
console.log(JSON.stringify({ scope: 'entry static JS + CSS (excludes lazy routes, fonts and API)',
  verifiedGzipFiles: compressedFiles.length, totals, files: rows }, null, 2));

// A route boundary is ineffective if its dependencies leak back into the entry graph.
const forbidden = /(?:vendor-(?:monaco|echarts)|Admin(?:Shell|Contests)?Layout|adminNavigation|admin-translations|Login-)/;
if (rows.some(({ file }) => forbidden.test(file))) {
  console.error('Heavy admin/chart/editor resources must stay outside the entry graph.');
  process.exitCode = 1;
}
