import assert from 'node:assert/strict';
import { toast, registerToastHandler } from '../../src/utils/toast.js';

// Only the React lifecycle boundary is replaced; tests import the production hooks.
let rendering;

export const notices = [];
export { toast };
registerToastHandler((notice) => {
  assert.equal(typeof notice.description, 'string', 'notifications must supply their description as an options object');
  notices.push(notice);
});
export const useTranslation = () => ({ t: (key) => key });

export function useState(initial) {
  const host = rendering;
  const index = host.cursor++;
  if (!host.slots[index]) {
    const slot = { value: typeof initial === 'function' ? initial() : initial };
    slot.set = (next) => {
      if (!host.mounted) {
        host.lateUpdates += 1;
        return;
      }
      const value = typeof next === 'function' ? next(slot.value) : next;
      if (!Object.is(value, slot.value)) {
        slot.value = value;
        host.dirty = true;
      }
    };
    host.slots[index] = slot;
  }
  const slot = host.slots[index];
  return [slot.value, slot.set];
}

export function useRef(value) {
  const host = rendering;
  const index = host.cursor++;
  host.slots[index] ??= { current: value };
  return host.slots[index];
}

export function useEffectEvent(callback) {
  const host = rendering;
  const index = host.cursor++;
  const slot = (host.slots[index] ??= {});
  slot.callback = callback;
  slot.event ??= (...args) => slot.callback(...args);
  return slot.event;
}

export function useEffect(setup, deps) {
  const host = rendering;
  const index = host.cursor++;
  const previous = host.slots[index];
  if (!previous || deps.some((dep, i) => !Object.is(dep, previous.deps[i]))) {
    host.effects.push({ index, setup, deps });
  }
}

export function mountHook(render) {
  const host = {
    slots: [],
    effects: [],
    mounted: true,
    dirty: true,
    lateUpdates: 0,
    render() {
      host.dirty = false;
      host.cursor = 0;
      rendering = host;
      try {
        host.value = render();
      } finally {
        rendering = null;
      }
      const effects = host.effects.splice(0);
      for (const { index } of effects) host.slots[index]?.cleanup?.();
      for (const { index, setup, deps } of effects) host.slots[index] = { deps, cleanup: setup() };
    },
    async flush() {
      // Drain async requests and batched state updates without advancing the fake clock.
      for (let turn = 0; turn < 20; turn += 1) {
        await Promise.resolve();
        if (host.mounted && host.dirty) host.render();
      }
      assert.equal(host.dirty && host.mounted, false, 'hook updates must settle');
    },
    unmount() {
      if (!host.mounted) return;
      host.mounted = false;
      for (const slot of host.slots) slot?.cleanup?.();
    },
  };
  host.render();
  return host;
}
