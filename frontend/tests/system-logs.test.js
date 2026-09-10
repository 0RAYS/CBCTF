import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import { registerHooks } from 'node:module';
import test from 'node:test';
import { transformWithOxc } from 'vite';
import { mountHook, notices } from './helpers/pollingHookHarness.js';
import { requests } from './helpers/logRequestHarness.js';
import { findElement } from './helpers/systemLogsHarness.js';

const pageURL = new URL('../src/pages/admin/Logs.jsx', import.meta.url).href;
const helperURL = new URL('./helpers/systemLogsHarness.js', import.meta.url).href;
const requestURL = new URL('./helpers/logRequestHarness.js', import.meta.url).href;
const { code } = await transformWithOxc(await readFile(new URL(pageURL), 'utf8'), 'Logs.jsx', {
  jsx: { runtime: 'automatic' },
});
const loader = registerHooks({
  resolve(specifier, context, nextResolve) {
    if (context.parentURL === pageURL) {
      if (
        ['react', 'react-i18next', '@tabler/icons-react', '../../utils/toast', '../../components/common'].includes(
          specifier
        ) ||
        /\/(useIpLookup|IpLookupDialog)$/.test(specifier)
      ) {
        return { url: helperURL, shortCircuit: true };
      }
      if (specifier.startsWith('.')) return nextResolve(`${specifier}.js`, context);
    }
    if (context.parentURL?.endsWith('/api/admin/system.js') && specifier === '../request') {
      return { url: requestURL, shortCircuit: true };
    }
    return nextResolve(specifier, context);
  },
  load(url, context, nextLoad) {
    return url === pageURL ? { format: 'module', source: code, shortCircuit: true } : nextLoad(url, context);
  },
});
const { default: AdminLogs } = await import(pageURL);
loader.deregister();

function fixture(t) {
  requests.length = 0;
  notices.length = 0;
  const observers = [];
  t.mock.method(globalThis, 'IntersectionObserver', function (callback) {
    const observer = { callback, observe() {}, disconnect() {} };
    observers.push(observer);
    return observer;
  });
  const host = mountHook(() => {
    const element = AdminLogs();
    // Attach the two DOM refs before the existing harness runs the page's effects.
    findElement(element, (item) => {
      if (item.props?.ref) item.props.ref.current = { scrollTop: 0 };
      return false;
    });
    return element;
  });
  t.after(() => host.unmount());
  return {
    host,
    get logs() {
      return findElement(host.value, (item) => item.type === 'pre').props.content;
    },
    get sentinel() {
      return findElement(host.value, (item) => item.type === 'pre').props.sentinel;
    },
    async resolve(index, data) {
      requests[index].resolve({ code: 200, data });
      await host.flush();
    },
    async intersect() {
      observers.at(-1).callback([{ isIntersecting: true }]);
      await host.flush();
    },
    async refresh() {
      findElement(host.value, (item) => item.type === 'button').props.onClick();
      await host.flush();
    },
  };
}

// Node has no DOM; restore this placeholder after each fixture's mock.
globalThis.IntersectionObserver ??= class {};

for (const failure of ['business', 'transport']) {
  for (const page of [1, 2]) {
    test(`system logs: ${failure} failure on page ${page} pauses scrolling until Refresh`, async (t) => {
      const f = fixture(t);
      const firstPage = Array.from({ length: 100 }, (_, index) => `line ${index}`);
      if (page === 2) {
        await f.resolve(0, firstPage);
        await f.intersect();
      }
      const failedIndex = page - 1;
      assert.equal(requests[failedIndex].config.params.offset, failedIndex * 100);
      if (failure === 'business') requests[failedIndex].resolve({ code: 500, msg: 'Redis unavailable', data: null });
      else requests[failedIndex].reject(new Error('request failed'));
      await f.host.flush();

      assert.deepEqual(f.logs, page === 2 ? firstPage : []);
      assert.equal(f.sentinel.props.role, 'alert');
      assert.equal(f.sentinel.props.children, 'admin.logs.toast.fetchFailed');
      for (let repeat = 0; repeat < 3; repeat += 1) await f.intersect();
      assert.equal(requests.length, page, 'intersection callbacks must not retry or skip the failed page');

      await f.refresh();
      assert.equal(requests.length, page + 1);
      assert.equal(requests[page].config.params.offset, 0, 'Refresh restarts from the first page');
      assert.equal(f.sentinel.props.role, undefined);
      await f.resolve(page, firstPage);
      await f.intersect();
      assert.equal(requests[page + 1].config.params.offset, 100);
      await f.resolve(page + 1, ['last line']);
      assert.deepEqual(f.logs, [...firstPage, 'last line']);
      assert.equal(f.sentinel.props.children, 'admin.logs.noMore');
    });
  }
}

test('system logs: a successful empty page is exhausted without an error or automatic retry', async (t) => {
  const f = fixture(t);
  await f.resolve(0, []);
  assert.deepEqual(f.logs, []);
  assert.equal(f.sentinel.props.children, 'admin.logs.noMore');
  assert.deepEqual(notices, []);
  await f.intersect();
  assert.equal(requests.length, 1);
});

test('system logs: a stale failure cannot pause a refreshed scope', async (t) => {
  const f = fixture(t);
  await f.refresh();
  requests[0].resolve({ code: 500, data: null });
  await f.host.flush();
  assert.equal(f.sentinel.props.role, undefined);
  await f.resolve(1, Array(100).fill('current log'));
  await f.intersect();
  assert.equal(requests.length, 3);
  assert.equal(requests[2].config.params.offset, 100);
});

test('system logs: changing the level resumes a failed scope from page one', async (t) => {
  const f = fixture(t);
  requests[0].resolve({ code: 500, data: null });
  await f.host.flush();
  findElement(f.host.value, (item) => item.type === 'select').props.onChange({ target: { value: 'ERROR' } });
  await f.host.flush();
  assert.equal(f.sentinel.props.role, undefined);
  assert.deepEqual(requests[1].config.params, { limit: 100, offset: 0, level: 'ERROR' });
  await f.resolve(1, []);
  assert.equal(f.sentinel.props.children, 'admin.logs.noMore');
});
