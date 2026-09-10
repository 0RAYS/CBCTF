import assert from 'node:assert/strict';
import { registerHooks } from 'node:module';
import test from 'node:test';
import { mountHook, notices } from './helpers/pollingHookHarness.js';
import { requests } from './helpers/logRequestHarness.js';

const helperURL = new URL('./helpers/pollingHookHarness.js', import.meta.url).href;
const requestURL = new URL('./helpers/logRequestHarness.js', import.meta.url).href;
const loader = registerHooks({
  resolve(specifier, context, nextResolve) {
    if (/\/(useGeneratorLogs|useVictimLogSession)\.js$/.test(context.parentURL ?? '')) {
      if (['react', 'react-i18next', '../../../../utils/toast'].includes(specifier)) {
        return { url: helperURL, shortCircuit: true };
      }
    }
    if (/\/api\/admin\/(generators|victims|contest)\.js$/.test(context.parentURL ?? '')) {
      if (specifier === '../request') return { url: requestURL, shortCircuit: true };
    }
    return nextResolve(specifier, context);
  },
});
const { default: useGeneratorLogs } = await import('../src/components/features/Admin/generators/useGeneratorLogs.js');
const { default: useVictimLogSession } = await import('../src/components/features/Admin/victims/useVictimLogSession.js');
const { getGeneratorLogs } = await import('../src/api/admin/generators.js');
const { getVictimPods, getVictimPodLogs } = await import('../src/api/admin/victims.js');
const { getContestVictimPods, getContestVictimPodLogs } = await import('../src/api/admin/contest.js');
loader.deregister();

const pods = [
  { name: 'pod-a', containers: ['web', 'sidecar'] },
  { name: 'pod-b', containers: ['worker'] },
];

function fixture(t, kind) {
  t.mock.timers.enable({ apis: ['setTimeout'] });
  notices.length = 0;
  requests.length = 0;
  let id = 11;
  const translationKey = kind === 'contest-victim' ? 'admin.contests.containers' : 'admin.victims';
  // Match the page adapters: neither victim page unwraps data or rejects business errors.
  const scope = {
    translationKey,
    loadPods: kind === 'contest-victim' ? (victimId) => getContestVictimPods(7, victimId) : getVictimPods,
    loadLogs:
      kind === 'contest-victim'
        ? (victimId, pod, container, lines) => getContestVictimPodLogs(7, victimId, pod, container, lines)
        : getVictimPodLogs,
  };
  const host = mountHook(() =>
    kind === 'generator'
      ? useGeneratorLogs({ logs: getGeneratorLogs }, id, (key) => `admin.generators.${key}`)
      : useVictimLogSession({ ...scope, victim: id ? { id } : null })
  );
  t.after(() => host.unmount());
  return {
    host,
    get value() {
      return host.value;
    },
    description(stage) {
      return kind === 'generator'
        ? 'admin.generators.toast.logFailed'
        : `${translationKey}.logs.${stage === 'pods' ? 'fetchPodsFailed' : 'fetchLogsFailed'}`;
    },
    async tick(ms) {
      t.mock.timers.tick(ms);
      await host.flush();
    },
    async resolve(index, data) {
      requests[index].resolve({ code: 200, data });
      await host.flush();
    },
    async ready() {
      if (kind !== 'generator') {
        await this.resolve(0, { pods });
        await this.tick(0);
      }
      return requests.length - 1;
    },
    async changeLines(lines) {
      if (kind === 'generator') host.value.changeLines(String(lines));
      else host.value.setLines(lines);
      await host.flush();
    },
    async switchScope() {
      id = 22;
      host.render();
      await host.flush();
    },
    async close() {
      if (kind === 'generator') host.unmount();
      else {
        id = null;
        host.render();
      }
      await host.flush();
    },
  };
}

function fail(index, result) {
  if (result === 'business') requests[index].resolve({ code: 500, msg: 'backend refused' });
  else requests[index].reject(new Error('HTTP request failed'));
}

for (const kind of ['generator', 'victim', 'contest-victim']) {
  const stages = kind === 'generator' ? ['logs'] : ['pods', 'logs'];
  for (const stage of stages) {
    for (const result of ['business', 'HTTP reject']) {
      test(`${kind}: active ${stage} ${result} reports a danger description and settles loading`, async (t) => {
        const f = fixture(t, kind);
        const index = stage === 'pods' ? 0 : await f.ready();
        fail(index, result);
        await f.host.flush();
        assert.equal(stage === 'pods' ? f.value.podsLoading : f.value.loading, false);
        assert.equal(f.value.content, '');
        assert.deepEqual(
          notices.map(({ color, description }) => ({ color, description })),
          [{ color: 'danger', description: f.description(stage) }]
        );
        if (stage === 'pods') {
          await f.tick(1000);
          assert.deepEqual(f.value.pods, []);
          assert.equal(requests.length, 1, 'failed pod discovery must not start logs');
        }
      });

      for (const change of ['switchScope', 'close']) {
        test(`${kind}: ${stage} ${result} after ${change} is silent`, async (t) => {
          const f = fixture(t, kind);
          const index = stage === 'pods' ? 0 : await f.ready();
          await f[change]();
          fail(index, result);
          await f.host.flush();
          assert.deepEqual(notices, []);
          assert.equal(f.value.content, '');
          assert.equal(f.host.lateUpdates, 0);
          if (change === 'switchScope') {
            const current = requests.length - 1;
            if (kind === 'generator') assert.equal(f.value.loading, true);
            else {
              assert.equal(f.value.podsLoading, true);
              await f.resolve(current, { pods });
              await f.tick(0);
            }
            await f.resolve(requests.length - 1, { logs: 'current scope' });
            assert.equal(f.value.content, 'current scope');
          }
        });
      }
    }
  }

  test(`${kind}: real API parameters, successful content, and trailing 500ms debounce`, async (t) => {
    const f = fixture(t, kind);
    const index = await f.ready();
    const base =
      kind === 'contest-victim'
        ? '/admin/contests/7/victims/11'
        : `/admin/${kind === 'generator' ? 'generators' : 'victims'}/11`;
    assert.deepEqual(requests[index].config, {
      url: `${base}/${kind === 'generator' ? 'logs' : 'pods/logs'}`,
      method: 'GET',
      params: kind === 'generator' ? { lines: 1000 } : { pod_name: 'pod-a', container: 'web', lines: 1000 },
    });
    if (kind !== 'generator') {
      assert.deepEqual(requests[0].config, { url: `${base}/pods`, method: 'GET' });
      assert.deepEqual(f.value.pods, pods);
    }
    const content = '\u001b[32mready\u001b[0m\nsecond line\n';
    await f.resolve(index, { logs: content });
    assert.equal(f.value.content, content);
    assert.equal(f.value.loading, false);
    await f.changeLines(1000);
    await f.tick(500);
    assert.equal(requests.length, index + 1, 'unchanged line count does not refetch');
    await f.changeLines(20);
    await f.tick(499);
    assert.equal(requests.length, index + 1);
    await f.changeLines(30);
    await f.tick(499);
    assert.equal(requests.length, index + 1);
    await f.tick(1);
    assert.equal(requests.length, index + 2);
    assert.equal(requests[index + 1].config.params.lines, 30);
    await f.resolve(index + 1, { logs: 'latest lines' });
    assert.equal(f.value.content, 'latest lines');
    assert.equal(f.value.loading, false);
    assert.deepEqual(notices, []);
    await f.changeLines(40);
    await f.close();
    await f.tick(1000);
    assert.equal(requests.length, index + 2, 'closing cancels the debounce');
    assert.equal(f.host.lateUpdates, 0);
  });

  for (const result of ['business', 'HTTP reject', 'success']) {
    test(`${kind}: superseded logs ${result} is ignored during line-count debounce`, async (t) => {
      const f = fixture(t, kind);
      const index = await f.ready();
      await f.changeLines(25);
      if (result === 'success') await f.resolve(index, { logs: 'stale content' });
      else fail(index, result);
      await f.host.flush();
      assert.deepEqual(notices, []);
      assert.equal(f.value.content, '');
      await f.tick(499);
      assert.equal(requests.length, index + 1);
      await f.tick(1);
      await f.resolve(index + 1, { logs: 'new content' });
      assert.equal(f.value.content, 'new content');
      assert.equal(f.value.loading, false);
    });
  }
}

test('victim: pod/container selection loads immediately and invalidates the previous request', async (t) => {
  const f = fixture(t, 'victim');
  await f.ready();
  f.value.selectContainer('sidecar');
  await f.host.flush();
  await f.tick(0);
  assert.equal(requests[2].config.params.container, 'sidecar');
  requests[1].reject(new Error('old container'));
  await f.host.flush();
  assert.deepEqual(notices, []);
  await f.resolve(2, { logs: 'sidecar logs' });
  assert.equal(f.value.content, 'sidecar logs');
  f.value.selectPod('pod-b');
  await f.host.flush();
  assert.equal(f.value.content, '');
  await f.tick(0);
  assert.deepEqual(requests[3].config.params, { pod_name: 'pod-b', container: 'worker', lines: 1000 });
  await f.resolve(3, { logs: 'worker logs' });
  assert.equal(f.value.content, 'worker logs');
});
