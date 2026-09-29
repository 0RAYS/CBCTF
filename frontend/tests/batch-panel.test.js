import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import { registerHooks } from 'node:module';
import test from 'node:test';
import { createElement } from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import { createInstance } from 'i18next';
import { I18nextProvider } from 'react-i18next';
import { transformWithOxc } from 'vite';
import en from '../src/i18n/locales/en.json' with { type: 'json' };
import zh from '../src/i18n/locales/zh-CN.json' with { type: 'json' };

const panelURL = new URL('../src/components/features/Admin/batch/BatchResultPanel.jsx', import.meta.url).href;
const { code } = await transformWithOxc(await readFile(new URL(panelURL), 'utf8'), 'BatchResultPanel.jsx', {
  jsx: { runtime: 'automatic' },
});
const loader = registerHooks({
  resolve(specifier, context, nextResolve) {
    if (context.parentURL === panelURL && specifier === '../../../common') {
      return { url: 'data:text/javascript,export const Button = "button";', shortCircuit: true };
    }
    return nextResolve(specifier, context);
  },
  load(url, context, nextLoad) {
    return url === panelURL ? { format: 'module', source: code, shortCircuit: true } : nextLoad(url, context);
  },
});
const { default: BatchResultPanel } = await import(panelURL);
loader.deregister();

const result = {
  status: 'partial', requested: 3, succeeded: 1, failed: 1, skipped: 0, not_attempted: 1,
  items: [
    { id: 'node-a', phase: 'queued', status: 'success' },
    { id: 'node-b', phase: 'enqueue', status: 'failed', code: 'task.enqueueError' },
  ],
};

async function render(lng, props) {
  const i18n = createInstance();
  await i18n.init({ lng, resources: { en: { translation: en }, 'zh-CN': { translation: zh } }, interpolation: { escapeValue: false } });
  return renderToStaticMarkup(createElement(I18nextProvider, { i18n }, createElement(BatchResultPanel, props)));
}

test('partial feedback displays counts, item failure and queued hint in both languages', async () => {
  const english = await render('en', { result, queued: true });
  assert.match(english, /Partially completed/);
  assert.match(english, /Succeeded 1.*Failed 1.*Not attempted 1/);
  assert.match(english, /node-b/);
  assert.match(english, /Task could not be queued/);
  assert.match(english, /does not mean the workload is ready/);
  const chinese = await render('zh-CN', { result, queued: true });
  assert.match(chinese, /部分/);
  assert.match(chinese, /成功 1/);
  assert.match(chinese, /失败 1/);
  assert.match(chinese, /任务未能入队/);
  assert.doesNotMatch(chinese, /admin\.batch\./);
});

test('nested scan details retain evidence paths and escape untrusted target text', async () => {
  const html = await render('en', { result: { ...result, items: [
    { id: '<script>alert(1)</script>', phase: 'scan', status: 'failed', details: result },
  ] } });
  assert.match(html, /&lt;script&gt;alert\(1\)&lt;\/script&gt; \/ node-b/);
  assert.match(html, /Task could not be queued/);
  assert.doesNotMatch(html, /<script>/);
});

test('unknown request outcome stays visible as an accessible error', async () => {
  const html = await render('en', { error: 'Refresh before retrying' });
  assert.match(html, /role="alert"/);
  assert.match(html, /Refresh before retrying/);
});
