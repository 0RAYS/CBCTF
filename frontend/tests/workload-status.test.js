import assert from 'node:assert/strict';
import {registerHooks} from 'node:module';
import test from 'node:test';
import {mountHook} from './helpers/pollingHookHarness.js';

const helperURL = new URL('./helpers/pollingHookHarness.js', import.meta.url).href;
const loader = registerHooks({
  resolve(specifier, context, nextResolve) {
    if (context.parentURL?.endsWith('/useWorkloadStatus.js') && specifier === 'react') {
      return { url: helperURL, shortCircuit: true };
    }
    return nextResolve(specifier, context);
  },
});
const { default: useWorkloadStatus } = await import('../src/components/features/Admin/workloads/useWorkloadStatus.js');
loader.deregister();

test('workload status: serialized polling, stale responses, transport abort and error recovery', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout'] });
  const calls = [];
  let id = 1;
  const load = (resourceId, signal) => new Promise((resolve, reject) => calls.push({ resourceId, signal, resolve, reject }));
  const host = mountHook(() => useWorkloadStatus(id, load));
  t.after(() => host.unmount());
  await host.flush();
  host.value.refreshStatus();
  t.mock.timers.tick(10000);
  await host.flush();
  assert.equal(calls.length, 1);
  id = 2;
  host.render();
  await host.flush();
  assert.equal(calls[0].signal.aborted, true);
  calls[0].resolve({ code: 200, data: { pods: [{ name: 'stale' }] } });
  calls[1].resolve({ code: 200, data: { pods: [{ name: 'current' }] } });
  await host.flush();
  assert.equal(host.value.pods[0].name, 'current');
  assert.equal(host.value.refreshing, false);
  t.mock.timers.tick(5000);
  await host.flush();
  calls[2].reject(new Error('unavailable'));
  await host.flush();
  assert.equal(host.value.statusError, true);
  assert.equal(host.value.pods[0].name, 'current', 'stale data is preserved but marked');
  host.value.refreshStatus();
  await host.flush();
  calls[3].resolve({ code: 200, data: { pods: [] } });
  await host.flush();
  assert.equal(host.value.statusError, false);
  assert.deepEqual(host.value.pods, []);
  host.value.refreshStatus();
  await host.flush();
  host.unmount();
  assert.equal(calls[4].signal.aborted, true);
  calls[4].reject(new Error('aborted'));
  await host.flush();
  t.mock.timers.tick(10000);
  assert.equal(calls.length, 5);
  assert.equal(host.lateUpdates, 0);
});
