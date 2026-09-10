import { readdirSync, readFileSync } from 'node:fs';
import path from 'node:path';

// Deliberately dependency-free: these checks must work without Vite's transitive parsers.
// This small lexer covers the project's JS/JSX ESM syntax, not arbitrary JS evaluation.
// Comments, strings and regex literals are not dependency declarations. Computed import
// specifiers are reported rather than silently disappearing from the dependency graph.
export function scanImports(source) {
  const tokens = [];
  const lexeme = /\/\/[^\r\n]*|\/\*[\s\S]*?\*\/|'(?:\\[\s\S]|[^'\\])*'|"(?:\\[\s\S]|[^"\\])*"|`(?:\\[\s\S]|[^`\\])*`|[\w$]+|=>|\?\.|[^\s]/g;
  const regexLiteral = /\/(?:\\[^\r\n]|\[(?:\\[^\r\n]|[^\]\\\r\n])*\]|[^/\\[\r\n])+\/[dgimsuvy]*/y;
  for (let match; (match = lexeme.exec(source)); ) {
    const value = match[0];
    if (value.startsWith('//') || value.startsWith('/*')) continue;
    const previous = tokens.at(-1)?.value;
    if (value === '/' && (!previous || /^(?:=|\(|\[|,|:|;|!|\?|=>|return|case|throw)$/.test(previous))) {
      regexLiteral.lastIndex = match.index;
      const regex = regexLiteral.exec(source);
      if (regex) {
        lexeme.lastIndex = regexLiteral.lastIndex;
        tokens.push({ value: '<regex>', offset: match.index });
        continue;
      }
    }
    tokens.push({ value, offset: match.index, string: /^['"`]/.test(value) });
  }

  const imports = [];
  const errors = [];
  for (let index = 0; index < tokens.length; index += 1) {
    const token = tokens[index];
    if (!['import', 'export'].includes(token.value)) continue;
    if (['.', '?.'].includes(tokens[index - 1]?.value)) continue;
    const next = tokens[index + 1];
    if (next?.value === '.' || next?.value === ':') continue;
    const line = source.slice(0, token.offset).split(/\r\n|\n|\r/).length;
    const dynamic = token.value === 'import' && next?.value === '(';
    let specifier;
    if (dynamic) {
      specifier = tokens[index + 2];
      if (![',', ')'].includes(tokens[index + 3]?.value)) specifier = undefined;
    } else if (token.value === 'import' && next?.string) {
      specifier = next;
    } else {
      // Only export-from declarations create edges; local/default exports do not.
      if (token.value === 'export' && !['{', '*'].includes(next?.value)) continue;
      let cursor = index + 1;
      if (token.value === 'export' && next?.value === '{') {
        while (cursor < tokens.length && tokens[cursor].value !== '}') cursor += 1;
        if (tokens[cursor + 1]?.value !== 'from') continue;
      }
      for (; cursor < tokens.length; cursor += 1) {
        if (tokens[cursor].value === 'from' && tokens[cursor + 1]?.string) {
          specifier = tokens[cursor + 1];
          break;
        }
        if ([';', 'import', 'export'].includes(tokens[cursor].value)) break;
      }
    }
    if (!specifier?.string || /\\|\$\{/.test(specifier.value)) {
      errors.push(`line ${line}: unsupported ${dynamic ? 'computed/escaped import()' : 'module declaration'}`);
      continue;
    }
    imports.push({ specifier: specifier.value.slice(1, -1), dynamic, line });
  }
  return { imports, errors };
}

export function readSources(root, directory = 'src') {
  const sources = new Map();
  for (const entry of readdirSync(path.join(root, directory), { withFileTypes: true }).sort((a, b) =>
    a.name < b.name ? -1 : a.name > b.name ? 1 : 0,
  )) {
    const file = `${directory}/${entry.name}`;
    if (entry.isDirectory()) {
      for (const [name, source] of readSources(root, file)) sources.set(name, source);
    } else if (entry.isFile() && /\.(js|jsx)$/.test(file)) {
      sources.set(file, readFileSync(path.join(root, file), 'utf8'));
    }
  }
  return sources;
}

export function createResolver(root) {
  const directories = new Map();
  function exactFile(file) {
    const absolute = path.resolve(root, file);
    const { root: volume } = path.parse(absolute);
    let directory = volume;
    const segments = absolute.slice(volume.length).split(path.sep);
    for (let index = 0; index < segments.length; index += 1) {
      if (!directories.has(directory)) {
        directories.set(directory, readdirSync(directory, { withFileTypes: true }));
      }
      // existsSync alone incorrectly accepts casing mistakes on Windows/macOS.
      const entry = directories.get(directory).find((candidate) => candidate.name === segments[index]);
      if (!entry) return false;
      if (index === segments.length - 1) return entry.isFile();
      if (!entry.isDirectory()) return false;
      directory = path.join(directory, entry.name);
    }
    return false;
  }
  return (importer, specifier) => {
    if (!/^\.{1,2}\//.test(specifier)) return null;
    // CSS, JSON, images, fonts and Vite ?raw/?url/?worker share file resolution,
    // but only ordinary JS modules participate in synchronous execution cycles.
    const [pathname] = specifier.split(/[?#]/);
    const base = path.posix.normalize(path.posix.join(path.posix.dirname(importer), pathname));
    const candidates = [base, `${base}.js`, `${base}.jsx`, `${base}/index.js`, `${base}/index.jsx`];
    return candidates.find(exactFile) ?? null;
  };
}

export function isSynchronousModule(dependency) {
  return (
    !dependency.dynamic &&
    /\.(js|jsx)$/.test(dependency.target ?? '') &&
    !/[?&](?:raw|url|worker|sharedworker)(?:[=&#]|$)/.test(dependency.specifier)
  );
}

export function findCycles(graph) {
  const visited = new Set();
  const active = new Map();
  const stack = [];
  const cycles = [];
  function visit(file) {
    visited.add(file);
    active.set(file, stack.length);
    stack.push(file);
    for (const target of graph.get(file) ?? []) {
      if (active.has(target)) cycles.push([...stack.slice(active.get(target)), target]);
      else if (!visited.has(target)) visit(target);
    }
    stack.pop();
    active.delete(file);
  }
  // Back edges provide actual directed cycle witnesses, not a sorted SCC masquerading
  // as a path. Every cyclic component is detected; enumerating all simple cycles is exponential.
  for (const file of graph.keys()) if (!visited.has(file)) visit(file);
  return cycles;
}
