import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { registerHooks } from 'node:module';
import test from 'node:test';
import { api, mountHook } from './helpers/challengeSessionHarness.js';
import { registerToastHandler } from '../src/utils/toast.js';
import { POLL_TIMEOUT } from '../src/config/workload.js';

const helperURL = new URL('./helpers/challengeSessionHarness.js', import.meta.url).href;
const loader = registerHooks({
  resolve(specifier, context, nextResolve) {
    if (context.parentURL?.endsWith('/useChallengeSession.js')) {
      if (
        ['react', 'react-i18next'].includes(specifier) ||
        /\/(api\/challenge|utils\/toast|utils\/fileDownload)$/.test(specifier)
      ) {
        return { url: helperURL, shortCircuit: true };
      }
      if (specifier === '../models/challengeViewModel') return nextResolve(`${specifier}.js`, context);
    }
    return nextResolve(specifier, context);
  },
});
const { default: useChallengeSession } =
  await import('../src/components/features/CTFGame/Challenges/hooks/useChallengeSession.js');
loader.deregister();
registerToastHandler(() => {});

const status = (file = '', remote = '') => ({ code: 200, data: { file, init: true, remote: { status: remote } } });

test('frontend polling budgets cover the backend startup and attachment deadlines', () => {
  const victim = readFileSync(new URL('../../internal/task/victim.go', import.meta.url), 'utf8');
  const attachment = readFileSync(new URL('../../internal/task/attachment.go', import.meta.url), 'utf8');
  const startup = victim.match(/ReadyDeadline\s*=\s*time\.Now\(\)\.Add\((\d+)\s*\*\s*time\.Minute\)/);
  const queued = attachment.match(/asynq\.Deadline\(time\.Now\(\)\.Add\((\d+)\s*\*\s*time\.Minute\)\)/);
  assert.ok(startup && queued, 'review the polling contract when backend deadline definitions change');
  assert.ok(POLL_TIMEOUT.running > Number(startup[1]) * 60000);
  assert.ok(POLL_TIMEOUT.attachment > Number(queued[1]) * 60000);
});

function fixture(t, challenge = { id: 'one', type: 'dynamic' }) {
  t.mock.timers.enable({ apis: ['setTimeout'] });
  api.status = t.mock.fn(async () => status());
  for (const name of ['init', 'reset', 'start', 'stop']) api[name] = t.mock.fn(async () => ({ code: 200 }));
  const host = mountHook(() => useChallengeSession('contest', { updateChallenge() {}, onSolved() {} }));
  t.after(() => host.unmount());
  return {
    host,
    async open() {
      await host.value.openChallenge(challenge);
      await host.flush();
    },
    async tick(ms) {
      t.mock.timers.tick(ms);
      await host.flush();
    },
  };
}

test('dynamic initialization refreshes until the generated attachment is ready', async (t) => {
  const f = fixture(t);
  await f.open();
  await f.host.value.modalProps.onInitialize('one');
  await f.host.flush();
  assert.equal(f.host.value.selectedChallenge.attachment, '');
  api.status.mock.mockImplementation(async () => status('attachment.zip'));
  await f.tick(5000);
  assert.equal(f.host.value.selectedChallenge.attachment, 'attachment.zip');
  const calls = api.status.mock.callCount();
  await f.tick(60000);
  assert.equal(api.status.mock.callCount(), calls);
});

test('opening an initialized dynamic challenge resumes queued attachment polling beyond three minutes', async (t) => {
  const f = fixture(t);
  await f.open();
  const pending = Promise.withResolvers();
  api.status.mock.mockImplementation(() => pending.promise);
  await f.tick(210000);
  pending.resolve(status('attachment.zip'));
  await f.host.flush();
  assert.equal(f.host.value.selectedChallenge.attachment, 'attachment.zip');
});

test('reset clears the previous attachment and rejects its in-flight status response', async (t) => {
  const f = fixture(t);
  await f.open();
  const pending = Promise.withResolvers();
  api.status.mock.mockImplementation(() => pending.promise);
  await f.tick(5000);
  api.status.mock.mockImplementation(async () => status());
  await f.host.value.modalProps.onReset('one');
  await f.host.flush();
  pending.resolve(status('old.zip'));
  await f.host.flush();
  assert.equal(f.host.value.selectedChallenge.attachment, '');
  api.status.mock.mockImplementation(async () => status('new.zip'));
  await f.tick(5000);
  assert.equal(f.host.value.selectedChallenge.attachment, 'new.zip');
});

test('closing an attachment dialog prevents late results and further polling', async (t) => {
  const f = fixture(t);
  await f.open();
  const pending = Promise.withResolvers();
  api.status.mock.mockImplementation(() => pending.promise);
  await f.tick(5000);
  f.host.value.closeChallenge();
  pending.resolve(status('late.zip'));
  await f.host.flush();
  const calls = api.status.mock.callCount();
  await f.tick(60000);
  assert.equal(f.host.value.selectedChallenge, null);
  assert.equal(api.status.mock.callCount(), calls);
});

test('resetting a Pod challenge follows asynchronous shutdown to completion', async (t) => {
  const f = fixture(t, { id: 'one', type: 'pods' });
  api.status.mock.mockImplementation(async () => status('', 'running'));
  await f.open();
  await f.host.value.modalProps.onReset('one');
  await f.host.flush();
  api.status.mock.mockImplementation(async () => status('', 'terminating'));
  await f.tick(5000);
  assert.equal(f.host.value.selectedChallenge.instanceTerminating, true);
  api.status.mock.mockImplementation(async () => status());
  await f.tick(5000);
  assert.equal(f.host.value.selectedChallenge.instanceStatus, '');
});
