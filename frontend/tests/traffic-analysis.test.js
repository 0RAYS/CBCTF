import assert from 'node:assert/strict';
import test from 'node:test';
import {
  evidenceOffset,
  filterAnalysisRows,
  victimOverlaps,
} from '../src/components/features/Admin/traffic/trafficAnalysisModel.js';

test('analysis filtering intersects half-open replay windows and matches reverse flows', () => {
  const startedAt = '2026-01-01T00:00:00Z';
  const row = (time, end_time, src_ip = 'a', dst_ip = 'b') => ({
    evidence: { time, end_time, src_ip, dst_ip },
  });
  const spanning = row(startedAt, '2026-01-01T00:00:02Z');
  const inside = row('2026-01-01T00:00:01Z', null, 'b', 'a');
  const outside = row('2026-01-01T00:00:02Z');
  const options = {
    startedAt,
    window: { start: 1000, end: 2000 },
    currentWindow: true,
    edge: { source: 'a', target: 'b' },
  };
  assert.deepEqual(filterAnalysisRows([spanning, inside, outside], options), [spanning, inside]);
  assert.deepEqual(filterAnalysisRows([inside], { ...options, nodeId: 'c' }), []);
  assert.equal(evidenceOffset(inside.evidence, startedAt), 1000);
  assert.equal(evidenceOffset(inside.evidence, undefined), null);
});

test('cross-team report excludes unrelated instances', () => {
  const related = {
    ip: '8.8.8.8',
    teams: [
      { team_id: 1, victim_ids: [42] },
      { team_id: 2, victim_ids: [99] },
    ],
  };
  assert.deepEqual(victimOverlaps([related, { teams: [{ victim_ids: [10] }] }], '42'), [related]);
});
