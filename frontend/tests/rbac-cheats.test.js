import assert from 'node:assert/strict';
import test from 'node:test';
import {
  lastAvailablePage,
  successfulAssignmentIds,
  toggleVisibleCandidates,
  userUpdatePayload,
} from '../src/components/features/Admin/rbac/payloads.js';
import {
  cheatListParams,
  cheatReviewForm,
  cheatUpdatePayload,
} from '../src/components/features/Admin/cheats/payloads.js';

test('user update omits an empty password without mutating the form or dropping false flags', () => {
  const form = Object.freeze({ name: 'user', password: '', hidden: false, verified: false, banned: false });
  assert.deepEqual(userUpdatePayload(form), { name: 'user', hidden: false, verified: false, banned: false });
  assert.equal(form.password, '');
  assert.deepEqual(userUpdatePayload({ password: 'new-password' }), { password: 'new-password' });
  assert.deepEqual(userUpdatePayload({ password: ' ' }), { password: ' ' });
});

test('candidate select-all merges pages and deselects only visible IDs without duplicates', () => {
  const selected = Object.freeze([99, 1]);
  const firstPage = [{ id: 1 }, { id: 2 }];
  const merged = toggleVisibleCandidates(selected, firstPage);
  assert.deepEqual(merged, [99, 1, 2]);
  assert.deepEqual(toggleVisibleCandidates(merged, [{ id: 3 }]), [99, 1, 2, 3]);
  assert.deepEqual(toggleVisibleCandidates(merged, firstPage), [99]);
  assert.deepEqual(toggleVisibleCandidates(merged, []), merged);
  assert.deepEqual(selected, [99, 1]);
});

test('allSettled assignment policy counts only code 200, preserving rejected and API-failed IDs for retry', async () => {
  const ids = Object.freeze([41, 7, 103, 8]);
  const results = await Promise.allSettled([
    Promise.resolve({ code: 200 }),
    Promise.reject(new Error('network failure')),
    Promise.resolve({ code: 403 }),
    Promise.resolve({ code: 200 }),
  ]);
  const successful = successfulAssignmentIds(ids, results);
  assert.deepEqual(successful, [41, 8]);
  assert.deepEqual(
    ids.filter((id) => !successful.includes(id)),
    [7, 103]
  );
});

test('assignment outcomes stay aligned with input IDs rather than completion order', async () => {
  let finishFirst;
  const first = new Promise((resolve) => {
    finishFirst = resolve;
  });
  const results = Promise.allSettled([first, Promise.resolve({ code: 200 })]);
  finishFirst({ code: 409 });
  assert.deepEqual(successfulAssignmentIds([10, 20], await results), [20]);
});

test('assignment policy handles total success, total failure and empty selection', () => {
  const success = { status: 'fulfilled', value: { code: 200 } };
  const failure = { status: 'rejected', reason: new Error('failed') };
  assert.deepEqual(successfulAssignmentIds([1, 2], [success, success]), [1, 2]);
  assert.deepEqual(successfulAssignmentIds([1, 2], [failure, failure]), []);
  assert.deepEqual(successfulAssignmentIds([], []), []);
});

test('pagination falls back after removing the final member or assigning the final candidate', () => {
  assert.equal(lastAvailablePage(20, 10, 3), 2);
  assert.equal(lastAvailablePage(9, 10, 2), 1);
  assert.equal(lastAvailablePage(0, 10, 4), 1);
  assert.equal(lastAvailablePage(50, 10, 2), 2);
});

test('review form normalizes missing optional values without changing evidence', () => {
  const cheat = Object.freeze({ id: 4, reason: null, comment: null, checked: false, hash: 'evidence' });
  assert.deepEqual(cheatReviewForm(cheat), { reason: '', type: '', checked: false, comment: '' });
  assert.deepEqual(cheatUpdatePayload(cheatReviewForm(cheat), cheat), {});
  assert.equal(cheat.hash, 'evidence');
});

test('review updates retain explicit checked false and comment clearing as deltas', () => {
  const cheat = Object.freeze({ reason: 'reason', type: 'suspicious', checked: true, comment: 'old note' });
  const form = { ...cheatReviewForm(cheat), checked: false, comment: '' };
  assert.deepEqual(cheatUpdatePayload(form, cheat), { checked: false, comment: '' });
  assert.deepEqual(cheatUpdatePayload(cheatReviewForm(cheat), cheat), {});
});

test('review sends only changed review fields, never immutable evidence or IDs', () => {
  const cheat = {
    id: 8,
    reason: 'old',
    type: 'suspicious',
    checked: false,
    comment: '',
    hash: 'hash',
    ip: '127.0.0.1',
  };
  const form = Object.freeze({ ...cheatReviewForm(cheat), type: 'pass', reason: 'new', id: 99, hash: 'changed' });
  assert.deepEqual(cheatUpdatePayload(form, cheat), { reason: 'new', type: 'pass' });
  assert.deepEqual(cheatUpdatePayload({ ...cheatReviewForm(cheat), checked: true, comment: 'reviewed' }, cheat), {
    checked: true,
    comment: 'reviewed',
  });
});

test('cheat list payload preserves paging and independently optional filters', () => {
  assert.deepEqual(cheatListParams(1, 20, '', ''), { limit: 20, offset: 0 });
  assert.deepEqual(cheatListParams(3, 20, 'cheater', 'wrong_flag'), {
    limit: 20,
    offset: 40,
    type: 'cheater',
    reason_type: 'wrong_flag',
  });
  assert.deepEqual(cheatListParams(2, 20, '', 'same_web_ip'), { limit: 20, offset: 20, reason_type: 'same_web_ip' });
});
