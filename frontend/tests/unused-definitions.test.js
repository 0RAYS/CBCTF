import assert from 'node:assert/strict';
import test from 'node:test';
import { fileURLToPath } from 'node:url';
import { ESLint } from 'eslint';

const root = fileURLToPath(new URL('../', import.meta.url));
const eslint = new ESLint({ cwd: root });

async function unusedDefinitions(source) {
  const [result] = await eslint.lintText(source, { filePath: 'src/unused-definition-fixture.js' });
  return result.messages.filter(({ ruleId }) => ruleId === 'no-unused-vars');
}

test('unused parameters before a used parameter are errors', async () => {
  const messages = await unusedDefinitions('export function login(providerName, loginUrl) { return loginUrl; }');
  assert.equal(messages.length, 1);
  assert.equal(messages[0].severity, 2);
  assert.match(messages[0].message, /providerName/);
});

test('positional callback placeholders are allowed, but unused local variables are not', async () => {
  const messages = await unusedDefinitions(
    'export function indices(items) { const unused = true; return items.map((_, index) => index); }',
  );
  assert.equal(messages.length, 1);
  assert.equal(messages[0].severity, 2);
  assert.match(messages[0].message, /unused/);
});

test('unused caught errors are checked and optional catch bindings are allowed', async () => {
  const withBinding = await unusedDefinitions(
    'export function attempt(run) { try { return run(); } catch (error) { return null; } }',
  );
  assert.equal(withBinding.length, 1);
  assert.equal(withBinding[0].severity, 2);
  assert.match(withBinding[0].message, /error/);
  assert.deepEqual(
    await unusedDefinitions('export function attempt(run) { try { return run(); } catch { return null; } }'),
    [],
  );
});
