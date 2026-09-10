import assert from 'node:assert/strict';
import { registerHooks } from 'node:module';
import test from 'node:test';
import { mountHook, notices } from './helpers/pollingHookHarness.js';

const helperURL = new URL('./helpers/pollingHookHarness.js', import.meta.url).href;
const loader = registerHooks({
  resolve(specifier, context, nextResolve) {
    if (/\/(useGeneratorSession|useVictimList)\.js$/.test(context.parentURL ?? '')) {
      if (['react', 'react-i18next', '../../../../utils/toast'].includes(specifier)) {
        return { url: helperURL, shortCircuit: true };
      }
      if (specifier === './victimPayload') return nextResolve(`${specifier}.js`, context);
    }
    return nextResolve(specifier, context);
  },
});
const { default: useGeneratorSession } =
  await import('../src/components/features/Admin/generators/useGeneratorSession.js');
const { default: useVictimList } = await import('../src/components/features/Admin/victims/useVictimList.js');
loader.deregister();

function deferred() {
  let resolve;
  let reject;
  const promise = new Promise((yes, no) => {
    resolve = yes;
    reject = no;
  });
  return { promise, resolve, reject };
}

function fixture(t, kind) {
  t.mock.timers.enable({ apis: ['setInterval', 'setTimeout'] });
  notices.length = 0;
  const calls = [];
  const mutations = [];
  const load = (scope) => (params) => {
    const call = { ...deferred(), params, scope };
    calls.push(call);
    return call.promise;
  };
  const mutate = (operation) => (params) => {
    const call = { ...deferred(), operation, params };
    mutations.push(call);
    return call.promise;
  };
  const api = {
    list: load(1),
    challenges: async () => ({ code: 200, data: { challenges: [] } }),
    start: mutate('start'),
    stop: mutate('stop'),
  };
  let scope = { contestId: 1, translationKey: 'contest', loadVictims: load(1) };
  const render = () => (kind === 'generator' ? useGeneratorSession(api, (key) => key) : useVictimList(scope));
  let host = mountHook(render);
  t.after(() => host.unmount());
  return {
    calls,
    mutations,
    get host() {
      return host;
    },
    get value() {
      return host.value;
    },
    get rows() {
      return kind === 'generator' ? host.value.generators : host.value.containers;
    },
    async flush() {
      await host.flush();
    },
    async tick(ms) {
      t.mock.timers.tick(ms);
      await host.flush();
    },
    async resolve(index, id = `row-${index}`) {
      calls[index].resolve({
        code: 200,
        data: { generators: [{ id, status: 'running' }], victims: [{ id, status: 'running' }], count: 30, running: 1 },
      });
      await host.flush();
    },
    async page(page) {
      if (kind === 'generator') host.value.changePage(page);
      else host.value.onPageChange(page);
      await host.flush();
    },
    async filter() {
      if (kind === 'generator') host.value.toggleShowDeleted();
      else host.value.onFilterChange('challenge_id', '42');
      await host.flush();
    },
    async switchScope() {
      if (kind === 'generator') {
        // GeneratorManagement is keyed by contestId in the route.
        host.unmount();
        api.list = load(2);
        host = mountHook(render);
      } else {
        scope = { ...scope, contestId: 2, loadVictims: load(2) };
        host.render();
      }
      await host.flush();
    },
  };
}

for (const kind of ['generator', 'victim']) {
  test(`${kind}: 12s requests survive 10s polling ticks and polling continues`, async (t) => {
    const f = fixture(t, kind);
    setTimeout(
      () =>
        f.calls[0].resolve({
          code: 200,
          data: { generators: [{ id: 'slow' }], victims: [{ id: 'slow' }] },
        }),
      12000
    );
    await f.tick(10000);
    assert.equal(f.calls.length, 1, 'polling must not supersede an in-flight request');
    await f.tick(2000);
    assert.deepEqual(f.rows, [{ id: 'slow' }]);
    if (kind === 'generator') assert.equal(f.value.loading, false);
    await f.tick(8000);
    assert.equal(f.calls.length, 2);
    await f.tick(10000);
    assert.equal(f.calls.length, 2);
    await f.tick(2000);
    await f.resolve(1, 'next');
    assert.equal(f.rows[0].id, 'next');
    if (kind === 'generator') assert.equal(f.value.loading, false);
    await f.tick(8000);
    assert.equal(f.calls.length, 3);
  });

  for (const change of ['page', 'filter', 'switchScope']) {
    test(`${kind}: ${change} immediately loads a new query and ignores the old response`, async (t) => {
      const f = fixture(t, kind);
      await f.tick(9000);
      if (change === 'page') await f.page(2);
      else await f[change]();
      assert.equal(f.calls.length, 2);
      if (change === 'page') assert.equal(f.calls[1].params.offset, 20);
      if (change === 'filter') {
        assert.equal(f.calls[1].params.offset, 0);
        assert.equal(
          f.calls[1].params[kind === 'generator' ? 'deleted' : 'challenge_id'],
          kind === 'generator' ? true : '42'
        );
      }
      if (change === 'switchScope') assert.equal(f.calls[1].scope, 2);
      await f.resolve(0, 'stale');
      assert.deepEqual(f.rows, []);
      if (kind === 'generator') assert.equal(f.value.loading, true);
      await f.tick(10000);
      assert.equal(f.calls.length, 2, 'old completion must not unlock the new request');
      await f.resolve(1, 'current');
      assert.equal(f.rows[0].id, 'current');
      assert.equal(f.value.totalCount, 30);
      if (kind === 'generator') assert.equal(f.value.loading, false);
      await f.tick(10000);
      assert.equal(f.calls.length, 3);
    });
  }

  test(`${kind}: manual refresh supersedes a slow request without waiting for a tick`, async (t) => {
    const f = fixture(t, kind);
    await f.tick(9000);
    f.value.refresh();
    await f.flush();
    assert.equal(f.calls.length, 2);
    await f.resolve(1, 'manual');
    await f.resolve(0, 'stale');
    assert.equal(f.rows[0].id, 'manual');
    await f.tick(10000);
    assert.equal(f.calls.length, 3);
  });

  test(`${kind}: stale errors are silent and cannot release the latest request`, async (t) => {
    const f = fixture(t, kind);
    f.value.refresh();
    await f.flush();
    f.calls[0].reject(new Error('stale'));
    await f.flush();
    assert.deepEqual(notices, []);
    await f.tick(10000);
    assert.equal(f.calls.length, 2);
    await f.resolve(1);
    await f.tick(10000);
    assert.equal(f.calls.length, 3);
  });

  test(`${kind}: a failed slow request releases polling for retry`, async (t) => {
    const f = fixture(t, kind);
    await f.tick(12000);
    assert.equal(f.calls.length, 1);
    f.calls[0].reject(new Error('offline'));
    await f.flush();
    assert.equal(notices.length, 1);
    if (kind === 'generator') assert.equal(f.value.loading, false);
    await f.tick(8000);
    assert.equal(f.calls.length, 2);
    await f.resolve(1);
    assert.equal(f.rows[0].id, 'row-1');
  });

  test(`${kind}: disabling polling clears ticks, manual refresh works, and re-enabling resumes`, async (t) => {
    const f = fixture(t, kind);
    f.value.setRefreshInterval(0);
    await f.flush();
    for (let i = 0; i < f.calls.length; i += 1) await f.resolve(i);
    const stoppedCount = f.calls.length;
    await f.tick(60000);
    assert.equal(f.calls.length, stoppedCount);
    f.value.refresh();
    await f.flush();
    assert.equal(f.calls.length, stoppedCount + 1);
    await f.resolve(stoppedCount, 'manual-while-off');
    await f.tick(60000);
    assert.equal(f.calls.length, stoppedCount + 1);
    f.value.setRefreshInterval(10);
    await f.flush();
    await f.tick(10000);
    assert.equal(f.calls.length, stoppedCount + 2);
    await f.resolve(stoppedCount + 1, 'resumed');
    assert.equal(f.rows[0].id, 'resumed');
  });

  for (const result of ['resolve', 'reject']) {
    test(`${kind}: unmount clears timers and ignores a late ${result}`, async (t) => {
      const f = fixture(t, kind);
      f.host.unmount();
      if (result === 'resolve') await f.resolve(0);
      else f.calls[0].reject(new Error('late'));
      await f.tick(60000);
      assert.equal(f.calls.length, 1);
      assert.equal(f.host.lateUpdates, 0);
      assert.deepEqual(notices, []);
    });
  }
}

test('generator: successful stop and start immediately refresh even while a poll is pending', async (t) => {
  const f = fixture(t, 'generator');
  await f.resolve(0, 'running');
  f.value.toggleSelect('running');
  await f.flush();
  await f.tick(10000);
  assert.equal(f.calls.length, 2);
  const stopping = f.value.stop();
  assert.deepEqual(f.mutations[0].params, ['running']);
  f.mutations[0].resolve({ code: 200 });
  await stopping;
  await f.flush();
  assert.equal(f.calls.length, 3);
  assert.deepEqual(f.value.selectedIds, []);
  await f.resolve(2, 'after-stop');
  await f.resolve(1, 'before-stop');
  assert.equal(f.rows[0].id, 'after-stop');

  await f.page(2);
  const starting = f.value.start({ challenge: 2 });
  assert.deepEqual(f.mutations[1].params, ['challenge', 'challenge']);
  f.mutations[1].resolve({ code: 200 });
  await starting;
  await f.flush();
  assert.equal(f.calls.length, 5);
  assert.equal(f.value.currentPage, 1);
  assert.equal(f.calls[4].params.offset, 0);
  await f.resolve(4, 'after-start');
  await f.resolve(3, 'before-start');
  assert.equal(f.rows[0].id, 'after-start');
  assert.equal(f.value.loading, false);
  assert.deepEqual(
    notices.map(({ color, description }) => ({ color, description })),
    [{ color: 'success', description: 'toast.stopSuccess' }]
  );
});

for (const operation of ['start', 'stop']) {
  for (const result of ['business', 'HTTP reject']) {
    test(`generator: ${operation} ${result} preserves state, reports a danger description, and allows retry`, async (t) => {
      const f = fixture(t, 'generator');
      await f.resolve(0, 'running');
      f.value.toggleSelect('running');
      f.value.openStart();
      await f.flush();
      const pending = f.value[operation]({ challenge: 1 });
      await f.flush();
      assert.equal(f.value.pendingOperation, operation);
      await f.value[operation]({ challenge: 1 });
      assert.equal(f.mutations.length, 1, 'pending operations cannot be submitted twice');
      if (result === 'business') f.mutations[0].resolve({ code: 500, msg: 'refused' });
      else f.mutations[0].reject(new Error('offline'));
      await pending;
      await f.flush();
      assert.equal(f.value.pendingOperation, null);
      assert.equal(f.value.startModalOpen, true);
      assert.deepEqual(f.value.selectedIds, ['running']);
      assert.equal(f.calls.length, 1, 'failed mutations must not refresh');
      assert.deepEqual(
        notices.map(({ color, description }) => ({ color, description })),
        [{ color: 'danger', description: `toast.${operation}Failed` }]
      );
      const retry = f.value[operation]({ challenge: 1 });
      assert.equal(f.mutations.length, 2);
      f.mutations[1].resolve({ code: 200 });
      await retry;
      await f.flush();
      assert.equal(f.value.pendingOperation, null);
      assert.equal(f.calls.length, 2);
    });

    for (const change of ['switchScope', 'unmount']) {
      test(`generator: ${operation} ${result} after ${change} is silent`, async (t) => {
        const f = fixture(t, 'generator');
        await f.resolve(0, 'running');
        f.value.toggleSelect('running');
        await f.flush();
        const pending = f.value[operation]({ challenge: 1 });
        const oldHost = f.host;
        if (change === 'switchScope') await f.switchScope();
        else f.host.unmount();
        if (result === 'business') f.mutations[0].resolve({ code: 500, msg: 'late refusal' });
        else f.mutations[0].reject(new Error('late failure'));
        await pending;
        await f.flush();
        assert.deepEqual(notices, []);
        assert.equal(oldHost.lateUpdates, 0);
        assert.equal(f.calls.length, change === 'switchScope' ? 2 : 1);
        if (change === 'switchScope') assert.equal(f.value.pendingOperation, null);
      });
    }
  }
}

test('generator: empty start selection emits a warning description without sending a mutation', async (t) => {
  const f = fixture(t, 'generator');
  await f.resolve(0);
  await f.value.start({ challenge: 0 });
  assert.equal(f.mutations.length, 0);
  assert.equal(f.value.pendingOperation, null);
  assert.deepEqual(
    notices.map(({ color, description }) => ({ color, description })),
    [{ color: 'warning', description: 'toast.selectRequired' }]
  );
});
