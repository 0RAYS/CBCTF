import assert from 'node:assert/strict';
import { registerHooks } from 'node:module';
import test from 'node:test';
import { mountHook, notices } from './helpers/pollingHookHarness.js';

const helperURL = new URL('./helpers/pollingHookHarness.js', import.meta.url).href;
const loader = registerHooks({
  resolve(specifier, context, nextResolve) {
    if (context.parentURL?.endsWith('/useBatchAction.js') &&
        ['react', 'react-i18next', '../../../../utils/toast.js'].includes(specifier)) {
      return { url: helperURL, shortCircuit: true };
    }
    return nextResolve(specifier, context);
  },
});
const { default: useBatchAction } = await import('../src/components/features/Admin/batch/useBatchAction.js');
loader.deregister();

const batch = {
  status: 'partial', requested: 3, succeeded: 1, skipped: 0, failed: 1, not_attempted: 1,
  items: [{ id: '1', status: 'success', phase: 'queued' }, { id: '2', status: 'failed', phase: 'enqueue', code: 'task.enqueueError' }],
};

test('partial business failure preserves details, calls refresh and prevents duplicate submissions', async () => {
  notices.length = 0;
  const host = mountHook(() => useBatchAction('contest:1'));
  let resolve;
  let requests = 0;
  const results = [];
  const request = () => { requests++; return new Promise((done) => { resolve = done; }); };
  const options = { onResult: (result) => results.push(result) };
  const pending = host.value.run(request, options);
  await host.value.run(request, options);
  await host.flush();
  assert.equal(requests, 1);
  assert.equal(host.value.pending, true);
  resolve({ code: 500, data: batch });
  await pending;
  await host.flush();
  assert.equal(host.value.pending, false);
  assert.deepEqual(host.value.result, batch);
  assert.deepEqual(results, [batch]);
  assert.equal(notices.length, 1);
  assert.equal(notices[0].color, 'warning');
  host.unmount();
});

test('nested import result is preserved and a request failure has an explicit unknown outcome', async () => {
  const host = mountHook(() => useBatchAction(1));
  await host.value.run(async () => ({ code: 500, data: { batch } }), { nested: true });
  await host.flush();
  assert.deepEqual(host.value.result, batch);
  let refreshed = false;
  await host.value.run(async () => { throw new Error('network'); }, { onResult: () => { refreshed = true; } });
  await host.flush();
  assert.equal(host.value.result, null);
  assert.equal(host.value.error, 'admin.batch.unknownOutcome');
  assert.equal(refreshed, false);
  assert.equal(host.value.pending, false);
  host.unmount();
});

for (const close of ['scope change', 'unmount']) {
  test(`late batch response is ignored after ${close}`, async () => {
    notices.length = 0;
    let scope = 1;
    const host = mountHook(() => useBatchAction(scope));
    let resolve;
    let refreshed = false;
    const pending = host.value.run(() => new Promise((done) => { resolve = done; }), {
      onResult: () => { refreshed = true; },
    });
    if (close === 'unmount') host.unmount();
    else { scope = 2; host.render(); await host.flush(); }
    resolve({ code: 500, data: batch });
    await pending;
    await host.flush();
    assert.equal(refreshed, false);
    assert.equal(notices.length, 0);
    assert.equal(host.lateUpdates, 0);
    if (close === 'scope change') assert.equal(host.value.result, null);
    host.unmount();
  });
}
