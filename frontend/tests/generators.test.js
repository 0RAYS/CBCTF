import assert from 'node:assert/strict';
import test from 'node:test';
import {
  expandStartCounts,
  getGeneratorPageStats,
  getGeneratorResponseData,
  isGeneratorStoppable,
  normalizeStartCount,
  retainStoppableSelection,
} from '../src/components/features/Admin/generators/generatorUtils.js';

test('start counts normalize empty, negative, fractional and invalid input', () => {
  for (const value of ['', undefined, null, 'invalid', -3, Infinity]) {
    assert.equal(normalizeStartCount(value), 0);
  }
  assert.equal(normalizeStartCount('3.9'), 3);
  assert.equal(normalizeStartCount('12'), 12);
  assert.equal(normalizeStartCount(0), 0);
});

test('starting N instances preserves N copies of the challenge ID', () => {
  const counts = { 'challenge-a': 3, 'challenge-b': 2, 'challenge-c': 0 };
  assert.deepEqual(expandStartCounts(counts), [
    'challenge-a',
    'challenge-a',
    'challenge-a',
    'challenge-b',
    'challenge-b',
  ]);
  assert.deepEqual(counts, { 'challenge-a': 3, 'challenge-b': 2, 'challenge-c': 0 });
});

test('count expansion keeps numeric IDs as strings and ignores nonpositive counts', () => {
  assert.deepEqual(expandStartCounts({ 42: '2', other: -1, empty: '', invalid: 'oops' }), ['42', '42']);
  assert.deepEqual(expandStartCounts({}), []);
});

test('success and failure statistics aggregate only the supplied page', () => {
  const page = [{ success: 4, failure: 2 }, { success: 3 }, { failure: 1 }, { success: null, failure: null }];
  assert.deepEqual(getGeneratorPageStats(page), { successes: 7, failures: 3 });
  assert.deepEqual(getGeneratorPageStats([{ success: 100, failure: 20 }]), { successes: 100, failures: 20 });
  assert.deepEqual(getGeneratorPageStats([]), { successes: 0, failures: 0 });
});

test('only running generators are stoppable', () => {
  assert.equal(isGeneratorStoppable({ status: 'running' }), true);
  for (const status of ['waiting', 'pending', 'terminating', 'stopped', 'unknown', undefined]) {
    assert.equal(isGeneratorStoppable({ status }), false);
  }
  assert.equal(isGeneratorStoppable(null), false);
});

test('selection retains only running IDs on the current page after refresh', () => {
  const selected = ['running', 'stopped', 'off-page'];
  const page = [
    { id: 'running', status: 'running' },
    { id: 'stopped', status: 'terminating' },
    { id: 'unselected', status: 'running' },
  ];
  assert.deepEqual(retainStoppableSelection(selected, page), ['running']);
  assert.deepEqual(retainStoppableSelection(selected, []), []);
  assert.deepEqual(selected, ['running', 'stopped', 'off-page']);
});

test('successful responses allow missing or null data, including mutation responses', () => {
  assert.deepEqual(getGeneratorResponseData({ code: 200 }), {});
  assert.deepEqual(getGeneratorResponseData({ code: 200, data: null }), {});
  const data = { generators: [], count: 0 };
  assert.equal(getGeneratorResponseData({ code: 200, data }), data);
  assert.deepEqual(getGeneratorResponseData({ code: 200, data: { logs: '' } }), { logs: '' });
});

test('resolved business errors and absent responses cannot take success paths', () => {
  for (const response of [undefined, null, {}, { code: 500 }, { code: 403, data: { generators: [] } }]) {
    assert.throws(() => getGeneratorResponseData(response), /Generator request failed/);
  }
  assert.throws(() => getGeneratorResponseData({ code: 400, msg: 'denied' }), /denied/);
});
