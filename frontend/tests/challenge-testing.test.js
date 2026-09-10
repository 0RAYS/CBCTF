import assert from 'node:assert/strict';
import { test } from 'node:test';
import {
  createTestSession,
  hasReachedTarget,
  isInstanceTransitioning,
  normalizeInstanceStatus,
} from '../src/components/features/Admin/challenges/testing/testSession.js';

const response = (status, extra = {}) => ({
  code: 200,
  data: { file: 'test.zip', remote: { status, duration: 120, remaining: 90, target: ['127.0.0.1:80'] }, ...extra },
});
const flush = () => new Promise(setImmediate);

function setup(t, overrides = {}, challengeId = 1) {
  const changes = [];
  const notifications = [];
  const api = {
    getTestChallengeStatus: t.mock.fn(async () => response('')),
    startTestVictim: t.mock.fn(async () => ({ code: 200 })),
    stopTestVictim: t.mock.fn(async () => ({ code: 200 })),
    ...overrides,
  };
  const session = createTestSession({
    challengeId,
    api,
    onChange: (state) => changes.push(state),
    notify: (...args) => notifications.push(args),
  });
  t.after(() => session.dispose());
  return { session, api, changes, notifications, state: () => changes.at(-1) };
}

test('normalizes backend states and preserves start/stop target semantics', () => {
  for (const status of ['waiting', 'pending', 'terminating', 'running']) {
    assert.equal(normalizeInstanceStatus(status.toUpperCase()), status);
    assert.equal(hasReachedTarget(status, 'running'), status === 'running');
    assert.equal(hasReachedTarget(status, 'stopped'), false);
    assert.equal(isInstanceTransitioning(status), status !== 'running');
  }
  for (const status of [undefined, null, 42, '', 'stopped', 'unknown']) {
    assert.equal(normalizeInstanceStatus(status), '');
    assert.equal(hasReachedTarget(normalizeInstanceStatus(status), 'stopped'), true);
    assert.equal(hasReachedTarget(normalizeInstanceStatus(status), 'running'), false);
  }
});

test('loads the complete backend payload and requires numeric code 200', async (t) => {
  const payload = response('RUNNING');
  const h = setup(t, { getTestChallengeStatus: async () => payload });
  await h.session.load();
  assert.deepEqual(h.state().testStatus, payload.data);
  assert.equal(h.state().loading.status, false);

  for (const code of [400, 500, '200']) {
    const rejected = setup(t, { getTestChallengeStatus: async () => ({ ...payload, code }) });
    await rejected.session.load();
    assert.equal(rejected.state().testStatus, null);
    assert.equal(rejected.state().loading.status, false);
  }
});

test('launch preserves attachments and targets while waiting, then finishes at running', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  const h = setup(t);
  await h.session.load();
  const before = h.state().testStatus;
  await h.session.start();
  assert.equal(h.state().loading.starting, true);
  assert.equal(h.state().testStatus.remote.status, 'waiting');
  assert.equal(h.state().testStatus.remote.remaining, 0);
  assert.equal(h.state().testStatus.file, before.file);
  assert.deepEqual(h.state().testStatus.remote.target, before.remote.target);
  h.api.getTestChallengeStatus.mock.mockImplementation(async () => response('RUNNING'));
  t.mock.timers.tick(5000);
  await flush();
  assert.equal(h.state().testStatus.remote.remaining, 90);
  assert.equal(h.state().loading.starting, false);
  const count = h.api.getTestChallengeStatus.mock.callCount();
  t.mock.timers.tick(180000);
  assert.equal(h.api.getTestChallengeStatus.mock.callCount(), count);
});

test('polling waits five seconds after completion and never overlaps slow requests', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  const pending = Promise.withResolvers();
  const getStatus = t.mock.fn(() => pending.promise);
  const h = setup(t, { getTestChallengeStatus: getStatus });
  await h.session.start();
  t.mock.timers.tick(4999);
  assert.equal(getStatus.mock.callCount(), 0);
  t.mock.timers.tick(1);
  assert.equal(getStatus.mock.callCount(), 1);
  t.mock.timers.tick(20000);
  assert.equal(getStatus.mock.callCount(), 1);
  pending.resolve(response('pending'));
  await flush();
  t.mock.timers.tick(4999);
  assert.equal(getStatus.mock.callCount(), 1);
  t.mock.timers.tick(1);
  assert.equal(getStatus.mock.callCount(), 2);
});

test('the three-minute deadline clears loading and rejects late poll responses', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  const pending = Promise.withResolvers();
  const h = setup(t, { getTestChallengeStatus: () => pending.promise });
  await h.session.start();
  t.mock.timers.tick(5000);
  t.mock.timers.tick(175000);
  assert.equal(h.state().loading.starting, false);
  const count = h.changes.length;
  pending.resolve(response('running'));
  await flush();
  assert.equal(h.changes.length, count);
  assert.equal(h.state().testStatus.remote.status, 'waiting');
});

test('stop refreshes immediately and does not poll once the instance is stopped', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  const h = setup(t);
  await h.session.stop();
  assert.equal(h.api.getTestChallengeStatus.mock.callCount(), 1);
  assert.deepEqual(h.state().loading, { status: false, starting: false, stopping: false });
  t.mock.timers.tick(180000);
  assert.equal(h.api.getTestChallengeStatus.mock.callCount(), 1);
});

test('stop waits through running and terminating until no active state remains', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  const statuses = ['terminating', 'running', ''];
  const h = setup(t, { getTestChallengeStatus: async () => response(statuses.shift()) });
  await h.session.stop();
  assert.equal(h.state().loading.stopping, true);
  t.mock.timers.tick(5000);
  await flush();
  assert.equal(h.state().loading.stopping, true);
  t.mock.timers.tick(5000);
  await flush();
  assert.equal(h.state().loading.stopping, false);
});

test('a slow immediate stop refresh does not extend the deadline', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  const pending = Promise.withResolvers();
  const getStatus = t.mock.fn(() => pending.promise);
  const h = setup(t, { getTestChallengeStatus: getStatus });
  const stopping = h.session.stop();
  await flush();
  t.mock.timers.tick(179000);
  assert.equal(getStatus.mock.callCount(), 1);
  pending.resolve(response('terminating'));
  await stopping;
  t.mock.timers.tick(1000);
  assert.deepEqual(h.state().loading, { status: false, starting: false, stopping: false });
  t.mock.timers.tick(5000);
  assert.equal(getStatus.mock.callCount(), 1);
});

test('a hung stop refresh releases all loading at the deadline and ignores its late result', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  const pending = Promise.withResolvers();
  const h = setup(t, { getTestChallengeStatus: () => pending.promise });
  const stopping = h.session.stop();
  await flush();
  assert.equal(h.state().loading.status, true);
  t.mock.timers.tick(180000);
  assert.deepEqual(h.state().loading, { status: false, starting: false, stopping: false });
  const count = h.changes.length;
  pending.resolve(response(''));
  await stopping;
  assert.equal(h.changes.length, count);
});

test('restores polling for waiting, pending and terminating when opening', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  for (const status of ['waiting', 'pending', 'terminating']) {
    const getStatus = t.mock.fn(async () => response(status));
    const h = setup(t, { getTestChallengeStatus: getStatus });
    await h.session.load();
    getStatus.mock.mockImplementation(async () => response(status === 'terminating' ? '' : 'running'));
    t.mock.timers.tick(5000);
    await flush();
    assert.equal(getStatus.mock.callCount(), 2);
    t.mock.timers.tick(180000);
    assert.equal(getStatus.mock.callCount(), 2);
    h.session.dispose();
  }
});

test('poll errors and non-200 responses retry silently without accepting data', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  let calls = 0;
  const h = setup(t, {
    getTestChallengeStatus: async () => {
      if (++calls === 1) throw new Error('offline');
      return calls === 2 ? { ...response('running'), code: 500 } : response('running');
    },
  });
  await h.session.start();
  for (let i = 0; i < 2; i++) {
    t.mock.timers.tick(5000);
    await flush();
    assert.equal(h.state().testStatus.remote.status, 'waiting');
  }
  t.mock.timers.tick(5000);
  await flush();
  assert.equal(h.state().testStatus.remote.status, 'running');
  assert.deepEqual(h.notifications, [['success', 'actionSuccess']]);
});

test('failed actions release loading and do not start polling', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  for (const action of ['start', 'stop']) {
    for (const failure of ['code', 'network']) {
      const h = setup(t, {
        [action === 'start' ? 'startTestVictim' : 'stopTestVictim']: async () => {
          if (failure === 'network') throw new Error('offline');
          return { code: 500 };
        },
      });
      await h.session[action]();
      assert.deepEqual(h.state().loading, { status: false, starting: false, stopping: false });
      assert.equal(h.notifications.length, failure === 'network' ? 1 : 0);
      t.mock.timers.tick(180000);
      assert.equal(h.api.getTestChallengeStatus.mock.callCount(), 0);
    }
  }
});

test('closing and reopening or switching challenges isolates late load/start/stop responses', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  for (const nextId of [1, 2]) {
    for (const action of ['load', 'start', 'stop']) {
      const pending = Promise.withResolvers();
      const apiName = { load: 'getTestChallengeStatus', start: 'startTestVictim', stop: 'stopTestVictim' }[action];
      const old = setup(t, { [apiName]: () => pending.promise });
      const request = old.session[action]();
      old.session.dispose();
      const count = old.changes.length;
      const reopened = setup(t, {}, nextId);
      await reopened.session.load();
      pending.resolve(response('running'));
      await request;
      t.mock.timers.tick(180000);
      await flush();
      assert.equal(old.changes.length, count);
      assert.equal(old.notifications.length, 0);
      assert.equal(reopened.state().testStatus.remote.status, '');
      assert.equal(reopened.api.getTestChallengeStatus.mock.calls[0].arguments[0], nextId);
    }
  }
});

test('disposing clears scheduled polls and suppresses late errors', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  const h = setup(t);
  await h.session.start();
  h.session.dispose();
  const count = h.changes.length;
  t.mock.timers.tick(180000);
  assert.equal(h.api.getTestChallengeStatus.mock.callCount(), 0);
  assert.equal(h.changes.length, count);

  const pending = Promise.withResolvers();
  const closing = setup(t, { getTestChallengeStatus: () => pending.promise });
  const request = closing.session.load();
  closing.session.dispose();
  pending.reject(new Error('late failure'));
  await request;
  assert.equal(closing.notifications.length, 0);
});

test('an in-flight poll cannot update or restart polling after closing and reopening', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'Date'] });
  const pending = Promise.withResolvers();
  const getStatus = t.mock.fn(() => pending.promise);
  const old = setup(t, { getTestChallengeStatus: getStatus });
  await old.session.start();
  t.mock.timers.tick(5000);
  old.session.dispose();
  const count = old.changes.length;
  const reopened = setup(t);
  await reopened.session.load();
  pending.resolve(response('pending'));
  await flush();
  t.mock.timers.tick(180000);
  assert.equal(old.changes.length, count);
  assert.equal(getStatus.mock.callCount(), 1);
  assert.equal(reopened.state().testStatus.remote.status, '');
});

test('repeated clicks cannot issue duplicate start/stop requests', async (t) => {
  const pending = Promise.withResolvers();
  const start = t.mock.fn(() => pending.promise);
  const h = setup(t, { startTestVictim: start });
  const first = h.session.start();
  await h.session.start();
  await h.session.stop();
  assert.equal(start.mock.callCount(), 1);
  assert.equal(h.api.stopTestVictim.mock.callCount(), 0);
  pending.resolve({ code: 500 });
  await first;
});
