import assert from 'node:assert/strict';
import test from 'node:test';
import {
  buildVictimListParams,
  buildVictimStartPayload,
  estimateVictimTeams,
  isVictimStoppable,
  selectPageVictims,
  toggleVictimSelection,
  updateChallengeSelection,
} from '../src/components/features/Admin/victims/victimPayload.js';

test('list parameters retain scope filters, paginate and only include deleted when requested', () => {
  const filters = { user_id: '', team_id: '12', challenge_id: '8', limit: 1, offset: 99 };
  assert.deepEqual(buildVictimListParams(filters, 3, false), {
    team_id: '12',
    challenge_id: '8',
    limit: 20,
    offset: 40,
  });
  assert.deepEqual(buildVictimListParams({ user_id: '4', challenge_id: '' }, 1, true), {
    user_id: '4',
    limit: 20,
    offset: 0,
    deleted: true,
  });
  assert.equal(filters.offset, 99);
  assert.equal('deleted' in buildVictimListParams({ deleted: true }, 1, false), false);
});

test('only running victims can be selected regardless of remaining time', () => {
  const victims = ['waiting', 'pending', 'terminating', 'running', 'stopped'].map((status, id) => ({
    id,
    status,
    remaining: 10,
  }));
  for (const victim of victims) assert.equal(isVictimStoppable(victim), victim.status === 'running');
  assert.equal(isVictimStoppable({ status: 'running', remaining: 0 }), true);
  assert.equal(isVictimStoppable(null), false);
  assert.deepEqual(toggleVictimSelection([99], victims, 0), [99]);
  assert.deepEqual(toggleVictimSelection([99], victims, 3), [99, 3]);
  assert.deepEqual(toggleVictimSelection([99, 3], victims, 3), [99]);
  assert.deepEqual(toggleVictimSelection([99], victims, 100), [99]);
});

test('stop select-all preserves the existing page replacement and count-based toggle policy', () => {
  const victims = [
    { id: 1, status: 'running' },
    { id: 2, status: 'running' },
    { id: 3, status: 'stopped' },
  ];
  assert.deepEqual(selectPageVictims([99], victims), [1, 2]);
  assert.deepEqual(selectPageVictims([1, 2], victims), []);
  assert.deepEqual(selectPageVictims([98, 99], victims), []);
  assert.deepEqual(selectPageVictims([99], []), []);
});

test('challenge selection retains IDs across pages and searches without duplicates', () => {
  let selected = updateChallengeSelection([], 10, true);
  selected = updateChallengeSelection(selected, 40, true);
  selected = updateChallengeSelection(selected, 10, true);
  assert.deepEqual(selected, [10, 40]);
  assert.deepEqual(updateChallengeSelection(selected, 40, false), [10]);
  assert.deepEqual(selected, [10, 40]);
});

test('start payload converts percentage to ratio and keeps duration in seconds', () => {
  const selected = [10, 40];
  const payload = buildVictimStartPayload(selected, 25, '7200');
  assert.deepEqual(payload, { challenges: [10, 40], team_ratio: 0.25, duration: 7200 });
  assert.notEqual(payload.challenges, selected);
  assert.equal(buildVictimStartPayload([], 1, '1').duration, 1);
  assert.equal(buildVictimStartPayload([], 99, '').duration, 0);
  assert.equal(buildVictimStartPayload([], 50, 'invalid').duration, 0);
  assert.equal(buildVictimStartPayload([], 50, '-1').duration, -1);
  assert.equal(buildVictimStartPayload([], 50, '1.9').duration, 1);
});

test('team estimate rejects endpoints and floors with a minimum of one team', () => {
  assert.equal(estimateVictimTeams(10, 25), 2);
  assert.equal(estimateVictimTeams(1, 1), 1);
  for (const percentage of [0, 100, -1, 101]) assert.equal(estimateVictimTeams(10, percentage), 0);
  assert.equal(estimateVictimTeams(0, 50), 0);
});
