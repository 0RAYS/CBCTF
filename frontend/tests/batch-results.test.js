import assert from 'node:assert/strict';
import test from 'node:test';
import { getBatchResult, hasBatchProgress, batchTone, remainingBatchIds, remainingGeneratorCounts, batchRows } from '../src/components/features/Admin/batch/batchModel.js';

test('partial responses are read even with a non-200 business code', () => {
  const batch = { status: 'partial', requested: 3, succeeded: 1, failed: 1, skipped: 0, not_attempted: 1, items: [] };
  assert.equal(getBatchResult({ code: 500, data: batch }), batch);
  assert.equal(getBatchResult({ code: 500, data: { batch } }, true), batch);
  assert.equal(getBatchResult({ code: 500, msg: 'preflight failed' }), null);
  assert.equal(batchTone(batch), 'warning');
});

test('retry selection retains failed and unattempted targets, not successes or skips', () => {
  const batch = { items: [{ id: '1', status: 'success' }, { id: '2', status: 'skipped' }, { id: '3', status: 'failed' }] };
  assert.deepEqual(remainingBatchIds([1, 2, 3, 4], batch), [3, 4]);
  assert.deepEqual(remainingGeneratorCounts(['a', 'a', 'b'], { items: [
    { id: 'challenge:a/instance:1/generator:9', status: 'success' },
    { id: 'challenge:a/instance:2', status: 'failed' },
  ] }), { a: 1, b: 1 });
});

test('nested failed scans still surface successful evidence and error paths', () => {
  const batch = { status: 'failed', succeeded: 0, skipped: 0, items: [{ id: 'web_ip', status: 'failed', details: {
    status: 'partial', succeeded: 1, skipped: 0, items: [{ id: '8.8.8.8', status: 'success' }, { id: '1.1.1.1', status: 'failed' }],
  } }] };
  assert.equal(hasBatchProgress(batch), true);
  assert.equal(batchTone(batch), 'warning');
  assert.equal(batchRows(batch)[2].path, 'web_ip / 1.1.1.1');
});
