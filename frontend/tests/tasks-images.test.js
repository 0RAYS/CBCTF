import test from 'node:test';
import assert from 'node:assert/strict';
import {
  createTaskQuery,
  formatDateTime,
  formatPayload,
  livePollingDelay,
  normalizeTaskPayload,
  PAGE_SIZE,
  TASK_STATUSES,
  taskQueryParams,
  taskQueryReducer,
  taskRowKey,
  taskTimestamp,
} from '../src/components/features/Admin/tasks/taskModel.js';
import {
  buildTargetKey,
  buildTargets,
  hasImageTag,
  missingTargetKeys,
  normalizeNodes,
  normalizePayload,
  normalizeTargetImages,
  parseManualImages,
  parseTargetKey,
} from '../src/components/features/Admin/images/imageModel.js';

test('task queries retain independent pages and filters', () => {
  let history = createTaskQuery();
  let live = createTaskQuery();
  history = taskQueryReducer(history, { type: 'filter', key: 'status', value: 'failed' });
  history = taskQueryReducer(history, { type: 'page', page: 4 });
  live = taskQueryReducer(live, { type: 'filter', key: 'queue', value: 'critical' });
  live = taskQueryReducer(live, { type: 'page', page: 2 });
  assert.deepEqual(history, { page: 4, filters: { task_id: '', queue: '', status: 'failed' } });
  assert.deepEqual(live, { page: 2, filters: { task_id: '', queue: 'critical', status: '' } });
  assert.notStrictEqual(history.filters, live.filters);
});

test('filter changes atomically reset the page without mutating the previous query', () => {
  const previous = { page: 5, filters: { task_id: 'task', queue: 'default', status: 'active' } };
  const next = taskQueryReducer(previous, { type: 'filter', key: 'queue', value: 'critical' });
  assert.equal(next.page, 1);
  assert.deepEqual(next.filters, { task_id: 'task', queue: 'critical', status: 'active' });
  assert.equal(previous.page, 5);
  assert.equal(previous.filters.queue, 'default');
  assert.deepEqual(taskQueryReducer(next, { type: 'reset' }), createTaskQuery());
});

test('task request parameters preserve filters and convert pages to offsets', () => {
  const filters = { task_id: 'abc', queue: 'default', status: 'retry' };
  assert.deepEqual(taskQueryParams({ page: 3, filters }), { ...filters, limit: PAGE_SIZE, offset: 40 });
  assert.equal(taskQueryParams(createTaskQuery()).offset, 0);
});

test('task response normalization tolerates missing and malformed collections', () => {
  const empty = { rows: [], totalCount: 0, queues: [] };
  assert.deepEqual(normalizeTaskPayload(), empty);
  assert.deepEqual(normalizeTaskPayload({ tasks: {}, queues: null }), empty);
  const tasks = [{ task_id: 'abc' }];
  assert.deepEqual(normalizeTaskPayload({ tasks, count: 42, queues: ['default'] }), {
    rows: tasks,
    totalCount: 42,
    queues: ['default'],
  });
});

test('history and live expose their actual status sets', () => {
  assert.deepEqual(TASK_STATUSES.history, ['', 'success', 'failed']);
  assert.deepEqual(TASK_STATUSES.live, ['', 'active', 'pending', 'scheduled', 'retry', 'archived', 'completed']);
});

test('shared task table retains history timestamps and live fallback order', () => {
  const task = { processed_at: 'history', next_process_at: 'next', completed_at: 'done', last_failed_at: 'failed' };
  assert.equal(taskTimestamp(task, 'history'), 'history');
  assert.equal(taskTimestamp(task, 'live'), 'next');
  assert.equal(taskTimestamp({ ...task, next_process_at: '' }, 'live'), 'done');
  assert.equal(taskTimestamp({ last_failed_at: 'failed' }, 'live'), 'failed');
  assert.equal(taskTimestamp({ completed_at: 'done' }, 'history'), undefined);
});

test('task row keys distinguish attempts, queues, and delimiter-like IDs', () => {
  const task = { id: 1, task_id: 'abc', queue: 'default' };
  assert.notEqual(taskRowKey(task, 'history'), taskRowKey({ ...task, id: 2 }, 'history'));
  assert.equal(taskRowKey(task, 'live'), taskRowKey({ ...task, id: 2 }, 'live'));
  assert.notEqual(taskRowKey(task, 'live'), taskRowKey({ ...task, queue: 'critical' }, 'live'));
  assert.notEqual(
    taskRowKey({ queue: 'a-b', task_id: 'c' }, 'live'),
    taskRowKey({ queue: 'a', task_id: 'b-c' }, 'live')
  );
});

test('only the active live tab polls, using the existing seconds selector', () => {
  for (const seconds of [5, 10, 30, 60]) {
    assert.equal(livePollingDelay(true, seconds), seconds * 1000);
    assert.equal(livePollingDelay(false, seconds), null);
  }
  assert.equal(livePollingDelay(true, 0), null);
  assert.equal(livePollingDelay(true, -1), null);
});

test('task formatting retains empty, invalid, structured, and circular payload fallbacks', () => {
  assert.equal(formatDateTime(null, 'en'), '-');
  assert.equal(formatDateTime('not-a-date', 'zh-CN'), 'not-a-date');
  const value = '2026-01-02T03:04:05Z';
  for (const language of ['en', 'zh-CN']) {
    assert.equal(
      formatDateTime(value, language),
      new Date(value).toLocaleString(language, {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
      })
    );
  }
  for (const empty of [null, undefined, '']) assert.equal(formatPayload(empty), '-');
  assert.equal(formatPayload('raw'), 'raw');
  assert.equal(formatPayload({ team: 3 }), '{"team":3}');
  assert.equal(formatPayload(0), '0');
  const circular = {};
  circular.self = circular;
  assert.equal(formatPayload(circular), '[object Object]');
});

test('image reference detection preserves tags, registry ports, and digests', () => {
  assert.equal(hasImageTag(' alpine:3 '), true);
  assert.equal(hasImageTag('registry:5000/team/app:latest'), true);
  assert.equal(hasImageTag('registry:5000/team/app'), false);
  assert.equal(hasImageTag('team/app@sha256:abcd'), true);
  for (const empty of ['', null, undefined, 'alpine']) assert.equal(hasImageTag(empty), false);
});

test('node normalization trims, filters, deduplicates, and sorts without changing input', () => {
  const input = [
    { node: 'z', images: [' b:2 ', 'a:1', 'b:2', 'untagged', null] },
    null,
    { images: ['ignored:1'] },
    { node: 'a', images: null },
  ];
  assert.deepEqual(normalizeNodes(input), [
    { node: 'a', images: [] },
    { node: 'z', images: ['a:1', 'b:2'] },
  ]);
  assert.equal(input[0].images[0], ' b:2 ');
  assert.deepEqual(normalizeNodes({}), []);
});

test('image payloads support node arrays and explicit contest targets, including an empty target list', () => {
  const nodes = [
    { node: 'a', images: ['a:1'] },
    { node: 'b', images: ['b:2', 'a:1'] },
  ];
  assert.deepEqual(normalizePayload(nodes), { nodes: normalizeNodes(nodes), targetImages: ['a:1', 'b:2'] });
  assert.deepEqual(normalizePayload({ nodes, target_images: [' c:3 ', 'c:3', 'untagged'] }).targetImages, ['c:3']);
  assert.deepEqual(normalizePayload({ nodes, target_images: [] }).targetImages, []);
  assert.deepEqual(normalizeTargetImages(undefined, nodes), ['a:1', 'b:2']);
  assert.deepEqual(normalizePayload(null), { nodes: [], targetImages: [] });
});

test('manual images use one tagged and deduplicated list for counts and submission', () => {
  assert.deepEqual(parseManualImages(' a:1\r\n\nb:2\na:1\nuntagged\nregistry:5000/app\n'), ['a:1', 'b:2']);
  assert.deepEqual(parseManualImages(' \n'), []);
});

test('target keys round-trip and distinguish ambiguous hyphenated combinations', () => {
  const target = { node: 'worker-a', image: 'registry:5000/team/app:1' };
  assert.deepEqual(parseTargetKey(buildTargetKey(target.node, target.image)), target);
  assert.notEqual(buildTargetKey('worker-a', 'b:1'), buildTargetKey('worker', 'a-b:1'));
});

test('missing targets exclude present images while manual targets retain every chosen pair', () => {
  const nodes = [
    { node: 'a', images: ['app:1'] },
    { node: 'b', images: [] },
  ];
  assert.deepEqual(missingTargetKeys(nodes, ['app:1', 'db:2']).map(parseTargetKey), [
    { node: 'b', image: 'app:1' },
    { node: 'a', image: 'db:2' },
    { node: 'b', image: 'db:2' },
  ]);
  assert.deepEqual(buildTargets(['a', 'b'], ['app:1', 'db:2']), [
    { node: 'a', image: 'app:1' },
    { node: 'a', image: 'db:2' },
    { node: 'b', image: 'app:1' },
    { node: 'b', image: 'db:2' },
  ]);
  assert.deepEqual(missingTargetKeys([], ['app:1']), []);
  assert.deepEqual(buildTargets(['a'], []), []);
});
