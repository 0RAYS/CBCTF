import assert from 'node:assert/strict';
import test from 'node:test';
import { startVisiblePolling } from '../src/utils/visiblePolling.js';

function fixture(t, hidden = false) {
  t.mock.timers.enable({ apis: ['setInterval'] });
  const document = new EventTarget();
  const window = new EventTarget();
  const navigator = { onLine: true };
  document.visibilityState = hidden ? 'hidden' : 'visible';
  for (const [key, value] of Object.entries({ document, window, navigator })) {
    const descriptor = Object.getOwnPropertyDescriptor(globalThis, key);
    Object.defineProperty(globalThis, key, { configurable: true, value });
    t.after(() => {
      if (descriptor) Object.defineProperty(globalThis, key, descriptor);
      else delete globalThis[key];
    });
  }
  let calls = 0;
  const stop = startVisiblePolling(() => { calls += 1; }, 10000);
  t.after(stop);
  return {
    get calls() { return calls; },
    stop,
    tick: (ms) => t.mock.timers.tick(ms),
    visible(value) {
      document.visibilityState = value ? 'visible' : 'hidden';
      document.dispatchEvent(new Event('visibilitychange'));
    },
    online(value) {
      navigator.onLine = value;
      window.dispatchEvent(new Event(value ? 'online' : 'offline'));
    },
  };
}

test('background polling sends no requests, resumes once and restarts the full interval', (t) => {
  const f = fixture(t);
  assert.equal(f.calls, 0, 'caller owns the initial request');
  f.tick(10000);
  assert.equal(f.calls, 1);
  f.visible(false);
  f.tick(600000);
  assert.equal(f.calls, 1);
  f.visible(true);
  f.visible(true);
  assert.equal(f.calls, 2, 'duplicate events must not duplicate requests');
  f.tick(9999);
  assert.equal(f.calls, 2);
  f.tick(1);
  assert.equal(f.calls, 3);
});

test('hidden initial mount and reconnect only poll once both online and visible', (t) => {
  const f = fixture(t, true);
  f.tick(60000);
  assert.equal(f.calls, 0);
  f.online(false);
  f.visible(true);
  f.tick(60000);
  assert.equal(f.calls, 0);
  f.online(true);
  assert.equal(f.calls, 1);
  f.online(false);
  f.tick(60000);
  assert.equal(f.calls, 1);
  f.visible(false);
  f.online(true);
  assert.equal(f.calls, 1);
  f.visible(true);
  assert.equal(f.calls, 2);
});

test('cleanup removes timers and event listeners even when hidden', (t) => {
  const f = fixture(t);
  f.visible(false);
  f.stop();
  f.visible(true);
  f.online(false);
  f.online(true);
  f.tick(60000);
  assert.equal(f.calls, 0);
});

test('disabled polling never schedules work', (t) => {
  t.mock.timers.enable({ apis: ['setInterval'] });
  let calls = 0;
  const stop = startVisiblePolling(() => { calls += 1; }, 0);
  t.mock.timers.tick(60000);
  stop();
  assert.equal(calls, 0);
});
