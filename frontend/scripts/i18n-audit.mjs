import { readFileSync, readdirSync } from 'node:fs';
import { fileURLToPath, pathToFileURL } from 'node:url';
import path from 'node:path';
import { Linter } from 'eslint';

export function flatten(value, prefix = '', result = {}) {
  for (const [key, child] of Object.entries(value)) {
    const name = prefix ? `${prefix}.${key}` : key;
    if (typeof child === 'string') result[name] = child;
    else flatten(child, name, result);
  }
  return result;
}

const escape = (value) => value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
const pluralBase = (key) => key.replace(/_(?:zero|one|two|few|many|other)$/, '');
export const placeholders = (text = '') =>
  [...new Set([...text.matchAll(/{{-?\s*([\w.]+)(?:\s*,[^}]+)?\s*}}/g)].map((match) => match[1]))].sort();

// This is a conservative usage audit, not a JavaScript evaluator. Computed keys retain
// matching branches. Review those branches and bilingual semantics before merging keys.
export function auditMessages(resources, sources) {
  const languages = Object.keys(resources);
  const dictionaries = Object.fromEntries(languages.map((language) => [language, flatten(resources[language])]));
  const keys = [...new Set(Object.values(dictionaries).flatMap(Object.keys))];
  const literals = new Set();
  const references = new Map();
  const directKeys = new Set();
  const templates = [];
  const dynamic = [];
  const addReference = (key, location, direct = false) => {
    if (typeof key !== 'string') return;
    literals.add(key);
    if (!references.has(key)) references.set(key, new Set());
    references.get(key).add(location);
    if (direct) directKeys.add(key);
  };
  const linter = new Linter();
  for (const [file, source] of sources) {
    const location = (node) => `${file}:${node.loc.start.line}`;
    const collectDirect = (node) => {
      if (!node) return;
      if (node.type === 'Literal') addReference(node.value, location(node), true);
      if (node.type === 'TemplateLiteral' && !node.expressions.length) {
        addReference(node.quasis[0].value.cooked, location(node), true);
      }
      if (node.type === 'ConditionalExpression') {
        collectDirect(node.consequent);
        collectDirect(node.alternate);
      }
      if (node.type === 'ArrayExpression') node.elements.forEach(collectDirect);
    };
    const visitors = {
      Literal(node) {
        if (typeof node.value !== 'string') return;
        literals.add(node.value);
        const call = node.parent;
        if (
          call?.type === 'CallExpression' &&
          ['text', 'label'].includes(call.callee.name) &&
          call.arguments[0] === node
        ) {
          // The generator domain passes a scoped translator between its components/hooks.
          if (file.startsWith('src/components/features/Admin/generators/') && call.callee.name === 'text') {
            addReference(`admin.generators.${node.value}`, location(node), true);
          }
          return;
        }
        addReference(node.value, location(node));
      },
      TemplateLiteral(node) {
        if (!node.expressions.length) addReference(node.quasis[0].value.cooked, location(node));
        else
          templates.push({
            location: location(node),
            parts: node.quasis.map((part) => part.value.cooked),
            source: source.slice(...node.range),
          });
      },
      CallExpression(node) {
        if (node.callee.name !== 't' && !(node.callee.type === 'MemberExpression' && node.callee.property.name === 't'))
          return;
        collectDirect(node.arguments[0]);
        if (node.arguments[0] && node.arguments[0].type !== 'Literal') {
          dynamic.push({ location: location(node), source: source.slice(...node.arguments[0].range) });
        }
      },
      JSXAttribute(node) {
        if (node.name.name === 'i18nKey') collectDirect(node.value?.expression ?? node.value);
      },
    };
    const messages = linter.verify(
      source,
      [
        {
          languageOptions: {
            ecmaVersion: 'latest',
            sourceType: 'module',
            parserOptions: { ecmaFeatures: { jsx: true } },
          },
          plugins: { audit: { rules: { collect: { create: () => visitors } } } },
          rules: { 'audit/collect': 'error' },
        },
      ],
      { allowInlineConfig: false }
    );
    if (messages.length) throw new Error(`${file}: ${JSON.stringify(messages)}`);
  }

  for (const dictionary of Object.values(dictionaries)) {
    for (const [key, value] of Object.entries(dictionary)) {
      for (const match of value.matchAll(/\$t\(([^,)]+)/g)) addReference(match[1].trim(), key, true);
    }
  }
  const prefixes = [...literals].filter(
    (value) => value.includes('.') && keys.some((key) => key.startsWith(`${value}.`))
  );
  const patterns = templates
    .filter(({ parts }) => parts.join('').replaceAll('.', '').length > 1)
    .map((template) => ({ ...template, pattern: new RegExp(`^${template.parts.map(escape).join('.*')}$`) }));
  const unused = keys.filter(
    (key) =>
      !literals.has(key) &&
      !literals.has(pluralBase(key)) &&
      !prefixes.some((prefix) => key.startsWith(`${prefix}.`)) &&
      !patterns.some(({ pattern }) => pattern.test(key) || pattern.test(pluralBase(key)))
  );
  const groups = new Map();
  for (const key of keys) {
    const text = JSON.stringify(languages.map((language) => dictionaries[language][key]));
    if (!groups.has(text)) groups.set(text, []);
    groups.get(text).push(key);
  }
  const roots = new Set(keys.map((key) => key.split('.')[0]));
  const missingReferences = [...references]
    .filter(
      ([key]) =>
        (directKeys.has(key) || (/^[\w:-]+(?:\.[\w:-]+)+$/.test(key) && roots.has(key.split('.')[0]))) &&
        !keys.includes(key) &&
        !keys.some((existing) => existing.startsWith(`${key}.`) || pluralBase(existing) === key)
    )
    .map(([key, locations]) => ({ key, locations: [...locations] }));
  const mismatch = languages.flatMap((language) =>
    keys.filter((key) => !(key in dictionaries[language])).map((key) => `${language}:${key}`)
  );
  const placeholderMismatch = keys.filter(
    (key) =>
      languages.every((language) => key in dictionaries[language]) &&
      new Set(languages.map((language) => JSON.stringify(placeholders(dictionaries[language][key])))).size > 1
  );
  const aliases = keys.filter((key) =>
    languages.some((language) => /^\$t\([^)]+\)$/.test(dictionaries[language][key] ?? ''))
  );
  return {
    keys: Object.fromEntries(languages.map((language) => [language, Object.keys(dictionaries[language]).length])),
    unused,
    duplicates: [...groups.values()].filter((group) => group.length > 1),
    mismatch,
    missingReferences,
    placeholderMismatch,
    aliases,
    retainedPrefixes: prefixes,
    dynamic,
    patterns: patterns
      .filter(({ pattern }) => keys.some((key) => pattern.test(key)))
      .map(({ location, source }) => ({ location, source })),
  };
}

function main() {
  const root = fileURLToPath(new URL('../', import.meta.url));
  const sources = new Map();
  function walk(directory) {
    for (const entry of readdirSync(directory, { withFileTypes: true })) {
      const file = path.join(directory, entry.name);
      if (entry.isDirectory()) walk(file);
      else if (/\.jsx?$/.test(entry.name))
        sources.set(path.relative(root, file).replaceAll('\\', '/'), readFileSync(file, 'utf8'));
    }
  }
  walk(path.join(root, 'src'));
  const resources = Object.fromEntries(
    ['en', 'zh-CN'].map((language) => [
      language,
      JSON.parse(readFileSync(path.join(root, 'src/i18n/locales', `${language}.json`), 'utf8')),
    ])
  );
  const report = auditMessages(resources, sources);
  const { keys, mismatch, missingReferences, placeholderMismatch, aliases } = report;
  const summary = {
    keys,
    unusedCandidates: report.unused.length,
    duplicateGroups: report.duplicates.length,
    mismatch,
    missingReferences,
    placeholderMismatch,
    aliases,
  };
  console.log(
    JSON.stringify(
      process.argv.includes('--details')
        ? report
        : process.argv.includes('--unused')
          ? { ...summary, unused: report.unused }
          : summary,
      null,
      2
    )
  );
  if (
    process.argv.includes('--check') &&
    [mismatch, missingReferences, placeholderMismatch, aliases, report.unused].some((items) => items.length)
  ) {
    process.exitCode = 1;
  }
}

if (process.argv[1] && import.meta.url === pathToFileURL(path.resolve(process.argv[1])).href) main();
