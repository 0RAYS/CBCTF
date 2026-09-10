import assert from 'node:assert/strict';
import test from 'node:test';
import { fileURLToPath } from 'node:url';
import {
  createResolver,
  findCycles,
  isSynchronousModule,
  readSources,
  scanImports,
} from './helpers/architecture.js';

// Resolve from this file, not process.cwd(): node --test works from either repo root or frontend.
const root = fileURLToPath(new URL('../', import.meta.url));
const sources = readSources(root);
const resolve = createResolver(root);
const dependencies = new Map();
const scanErrors = [];
for (const [file, source] of sources) {
  const { imports, errors } = scanImports(source);
  scanErrors.push(...errors.map((error) => `${file}: ${error}`));
  dependencies.set(
    file,
    imports.map((dependency) => ({ ...dependency, target: resolve(file, dependency.specifier) })),
  );
}

function assertNoViolations(violations, rule) {
  assert.equal(violations.length, 0, `${rule}\n${violations.join('\n')}`);
}

test('dependency scanner handles ESM declarations without treating examples as imports', () => {
  const source = [
    '// import Missing from "./comment.js";',
    '/* export * from "./comment-barrel.js"; */',
    'const text = "import(\'./example.js\')";',
    'const template = `export * from "./template.js"`;',
    'const pattern = /import("fake")/;',
    'import Default, {',
    '  named as renamed,',
    '} /* explanation */ from "./module";',
    'import "./style.css";',
    'import Worker from "./job.js?worker&inline";',
    'export { renamed };',
    'export { default as Widget } from "./Widget.jsx";',
    'export * from "./barrel";',
    'export * as namespace from "./namespace.js";',
    'const lazy = () => import(/* chunk */ "./lazy", { with: { type: "json" } });',
    'const literalTemplate = import(`./literal.js`);',
    'const url = import.meta.url;',
  ].join('\n');
  const { imports, errors } = scanImports(source);
  assert.deepEqual(errors, []);
  assert.deepEqual(
    imports.map(({ specifier, dynamic }) => [specifier, dynamic]),
    [
      ['./module', false],
      ['./style.css', false],
      ['./job.js?worker&inline', false],
      ['./Widget.jsx', false],
      ['./barrel', false],
      ['./namespace.js', false],
      ['./lazy', true],
      ['./literal.js', true],
    ],
  );
  assert.equal(imports[0].line, 6);
  assert.equal(scanImports('import("./" + name); import(`./${name}.js`);').errors.length, 2);
});

test('resolver preserves exact casing, extension/index resolution and asset queries', () => {
  // Existing test files make read-only fixtures; no temp files or installed packages needed.
  assert.equal(resolve('tests/entry.js', './helpers/architecture'), 'tests/helpers/architecture.js');
  assert.equal(resolve('tests/entry.js', './helpers/Architecture.js'), null);
  assert.equal(resolve('tests/entry.js', './Helpers/architecture.js'), null);
  assert.equal(resolve('tests/entry.js', './helpers/architecture.js?worker&inline'), 'tests/helpers/architecture.js');
  assert.equal(resolve('src/entry.js', './index.css?inline'), 'src/index.css');
  assert.equal(resolve('src/entry.js', './i18n/locales/en.json?raw'), 'src/i18n/locales/en.json');
  assert.equal(resolve('src/entry.js', './main'), 'src/main.jsx');
  assert.equal(resolve('src/entry.js', './store'), 'src/store/index.js');
  assert.equal(resolve('src/entry.js', './missing-file.js'), null);
  assert.equal(resolve('src/entry.js', 'monaco-editor/editor/editor.worker.js?worker'), null);
});

test('cycle detection returns real static paths, including barrel re-exports and self-imports', () => {
  const fixture = new Map([
    ['index.js', 'export * from "./leaf.js";'],
    ['leaf.js', 'import { value } from "./index.js";'],
    ['request.js', 'import "./store.js";'],
    ['store.js', 'import("./request.js");'],
    ['self.js', 'import "./self.js";'],
    ['asset.js', 'import code from "./asset.js?raw"; import Worker from "./asset.js?worker";'],
  ]);
  const graph = new Map(
    [...fixture].map(([file, source]) => [
      file,
      scanImports(source).imports
        .map((dependency) => ({ ...dependency, target: dependency.specifier.slice(2).split('?')[0] }))
        .filter(isSynchronousModule)
        .map(({ target }) => target),
    ]),
  );
  assert.deepEqual(findCycles(graph), [['index.js', 'leaf.js', 'index.js'], ['self.js', 'self.js']]);
});

test('all src relative imports/re-exports resolve with exact path casing', () => {
  assert.ok(sources.size > 0, 'src must not be empty');
  const violations = [...scanErrors];
  for (const [file, imports] of dependencies) {
    for (const { specifier, target, line } of imports) {
      if (/^\.{1,2}\//.test(specifier) && !target) violations.push(`${file}:${line} -> ${specifier}`);
    }
  }
  assertNoViolations(violations, 'Unresolved imports or unsupported module syntax (never silently skipped):');
});

test('common components do not directly depend on admin APIs or business features', () => {
  const violations = [];
  for (const [file, imports] of dependencies) {
    if (!file.startsWith('src/components/common/')) continue;
    for (const { target, line } of imports) {
      // Shared hooks, i18n (LanguageSwitcher), config and utilities remain legitimate dependencies.
      if (/^src\/(?:api\/admin|components\/features)\//.test(target ?? '')) {
        violations.push(`${file}:${line} -> ${target}`);
      }
    }
  }
  assertNoViolations(violations, 'common must remain business-independent:');
});

test('global hooks do not import or re-export Admin business modules', () => {
  const violations = [];
  for (const [file, imports] of dependencies) {
    if (!file.startsWith('src/hooks/')) continue;
    for (const { target, line } of imports) {
      if (/^src\/(?:api\/admin|components\/features\/Admin|pages\/admin)\//.test(target ?? '')) {
        violations.push(`${file}:${line} -> ${target}`);
      }
    }
  }
  assertNoViolations(violations, 'Admin-specific hooks belong with their feature:');
});

test('components never depend on pages', () => {
  const violations = [];
  for (const [file, imports] of dependencies) {
    if (!file.startsWith('src/components/')) continue;
    for (const { target, line } of imports) {
      if (target?.startsWith('src/pages/')) violations.push(`${file}:${line} -> ${target}`);
    }
  }
  assertNoViolations(violations, 'Pages compose components, not the reverse:');
});

test('synchronous source dependencies are acyclic, including barrels', () => {
  // Dynamic request/store imports are checked for resolution/layering, but do not
  // imply synchronous initialization cycles. There is no legacy-cycle allowlist.
  const graph = new Map(
    [...dependencies].map(([file, imports]) => [
      file,
      [...new Set(imports.filter(isSynchronousModule).map(({ target }) => target))],
    ]),
  );
  const cycles = findCycles(graph).map((cycle) =>
    cycle.map((file, index) => {
      const edge = dependencies.get(file)?.find((dependency) =>
        isSynchronousModule(dependency) && dependency.target === cycle[index + 1],
      );
      return edge ? `${file}:${edge.line}` : file;
    }).join(' -> '),
  );
  assertNoViolations(cycles, 'Static dependency cycle witnesses (break the dependency, do not whitelist):');
});

test('each src JS/JSX file stays within 500 physical lines', (context) => {
  // Include blank lines/comments; exclude only the final newline's phantom line.
  // JSON translations, styles, assets and tests are outside this source-size budget.
  const violations = [];
  let largest = { file: '', lines: 0 };
  for (const [file, source] of sources) {
    const lines = source === '' ? 0 : source.replace(/(?:\r\n|\n|\r)$/, '').split(/\r\n|\n|\r/).length;
    if (lines > largest.lines) largest = { file, lines };
    if (lines > 500) violations.push(`${file}: ${lines} lines (limit 500)`);
  }
  context.diagnostic(`${sources.size} JS/JSX files checked; largest: ${largest.file} (${largest.lines} lines)`);
  assertNoViolations(violations, 'Split by responsibility; do not raise the budget to bless existing oversized files:');
});

test('extracted business models/parsers/serializers stay isolated from runtime libraries', () => {
  // Reviewed, closed dependency set: no React/ECharts/Monaco, API, store, hooks,
  // toast or other runtime adapter, even via a barrel or a lazy import. Add new pure
  // helpers here after review. Loaders and effectful session controllers are not models.
  const pureModules = new Set([
    'config/contest.js',
    'components/features/Admin/SystemConfig/configNormalizer.js',
    'components/features/Admin/SystemConfig/configPayloadBuilder.js',
    'components/features/Scoreboard/scoreboardModel.js',
    'components/features/Scoreboard/timelineModel.js',
    'components/features/Scoreboard/timelineOption.js',
    'components/features/CTFGame/Challenges/models/challengeViewModel.js',
    'components/features/CTFGame/Team/model.js',
    'components/features/Admin/Contests/challenges/challengeData.js',
    'components/features/Admin/Contests/editor/contestForm.js',
    'components/features/Admin/Contests/teams/teamDetailData.js',
    'components/features/Admin/challenges/editor/challengePayload.js',
    'components/features/Admin/challenges/editor/composeParser.js',
    'components/features/Admin/challenges/editor/composeSerializer.js',
    'components/features/Admin/challenges/editor/composeValidation.js',
    'components/features/Admin/challenges/editor/guideModel.js',
    'components/features/Admin/challenges/editor/networkPolicy.js',
    'components/features/Admin/cheats/payloads.js',
    'components/features/Admin/generators/generatorUtils.js',
    'components/features/Admin/images/imageModel.js',
    'components/features/Admin/oauth/payload.js',
    'components/features/Admin/rbac/payloads.js',
    'components/features/Admin/smtp/payload.js',
    'components/features/Admin/tasks/taskModel.js',
    'components/features/Admin/traffic/trafficDemo.js',
    'components/features/Admin/traffic/trafficLayout.js',
    'components/features/Admin/traffic/trafficPresentation.js',
    'components/features/Admin/victims/victimPayload.js',
    'components/features/Admin/webhook/payload.js',
  ].map((file) => `src/${file}`));
  const violations = [];
  for (const file of pureModules) {
    if (!sources.has(file)) violations.push(`${file}: missing pure module; review its replacement`);
    for (const { specifier, target, line } of dependencies.get(file) ?? []) {
      if (!pureModules.has(target) || /[?#]/.test(specifier)) {
        violations.push(`${file}:${line} -> ${specifier} (outside reviewed pure modules)`);
      }
    }
  }
  assertNoViolations(violations, 'Pure business modules must form a runtime-independent dependency closure:');
});
