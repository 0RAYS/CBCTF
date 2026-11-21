import { existsSync, statSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

// Mirrors the Vite alias so `node --test` resolves the same "@/..." specifiers as the app.
const srcRoot = fileURLToPath(new URL('../src/', import.meta.url));

function resolveFile(target) {
  const candidates = [target, `${target}.js`, `${target}.jsx`, path.join(target, 'index.js'), path.join(target, 'index.jsx')];
  return candidates.find((candidate) => existsSync(candidate) && statSync(candidate).isFile()) ?? null;
}

export async function resolve(specifier, context, nextResolve) {
  if (!specifier.startsWith('@/')) return nextResolve(specifier, context);
  const queryIndex = specifier.search(/[?#]/);
  const pathname = queryIndex === -1 ? specifier : specifier.slice(0, queryIndex);
  const suffix = queryIndex === -1 ? '' : specifier.slice(queryIndex);
  const file = resolveFile(path.join(srcRoot, pathname.slice(2)));
  if (!file) return nextResolve(specifier, context);
  return nextResolve(`${pathToFileURL(file).href}${suffix}`, context);
}
