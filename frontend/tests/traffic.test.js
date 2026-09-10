import assert from 'node:assert/strict';
import test from 'node:test';
import {
  buildLines,
  buildPositions,
  clampPanToViewport,
  getViewportMetrics,
  resolveEdgeLabels,
  VIEWBOX_HEIGHT,
  VIEWBOX_WIDTH,
} from '../src/components/features/Admin/traffic/trafficLayout.js';
import {
  buildCompactNodeLabel,
  buildCompactNodeMeta,
  edgeChipClass,
  edgeTone,
  ellipsis,
  filterTraffic,
  formatBytes,
  formatDurationMs,
  formatProcessLabel,
  resolveTrafficSelection,
  resolveVisibleInlineItems,
  sanitizeSlice,
} from '../src/components/features/Admin/traffic/trafficPresentation.js';
import { createDemoTopology } from '../src/components/features/Admin/traffic/trafficDemo.js';

test('byte formatting preserves units and precision', () => {
  for (const [input, expected] of [
    [undefined, '0 B'],
    [0, '0 B'],
    [1023, '1023 B'],
    ['1024', '1.0 KB'],
    [1536, '1.5 KB'],
    [1024 ** 2, '1.0 MB'],
    [1024 ** 3, '1.0 GB'],
  ])
    assert.equal(formatBytes(input), expected);
});

test('duration formatting and slice inputs stay in milliseconds', () => {
  for (const [input, expected] of [
    [undefined, '0 ms'],
    [1, '1 ms'],
    [999, '999 ms'],
    [1000, '1 s'],
    [1001, '1.00 s'],
    [1250, '1.25 s'],
    [60000, '1m'],
    [61000, '1m 1s'],
    [61250, '1m 1.25s'],
  ])
    assert.equal(formatDurationMs(input), expected);
  for (const [input, expected] of [
    [undefined, 1],
    ['', 1],
    ['invalid', 1],
    [-1, 1],
    [0, 1],
    ['1.6', 2],
    ['1500', 1500],
  ]) {
    assert.equal(sanitizeSlice(input), expected);
  }
});

test('compact labels keep service, process and PID precedence', () => {
  assert.equal(ellipsis(null), '');
  assert.equal(ellipsis('abcdef', 3), 'abc...');
  assert.equal(ellipsis('abc', 3), 'abc');
  assert.equal(buildCompactNodeLabel({ service: 'web', ip: '10.0.0.1' }), 'web');
  assert.equal(buildCompactNodeLabel({ kind: 'victim', label: 'Victim', ip: '10.0.0.1' }), 'Victim');
  assert.equal(buildCompactNodeLabel({ kind: 'peer', label: 'Peer', ip: '10.0.0.1' }), '10.0.0.1');
  assert.equal(formatProcessLabel(null), '');
  assert.equal(formatProcessLabel({ process_name: 'nginx', dominant_process: 'other', pid: 42 }), 'nginx (42)');
  assert.equal(formatProcessLabel({ dominant_process: 'nginx' }), 'nginx');
  assert.equal(formatProcessLabel({ pid: 42 }), 'PID 42');
  assert.equal(formatProcessLabel({}), '');
  assert.equal(buildCompactNodeMeta({ dominant_process: 'main', processes: [{ process_name: 'other' }] }), 'main');
  assert.equal(buildCompactNodeMeta({ processes: [{ process_name: 'nginx', pid: 42 }] }), 'nginx (42)');
  assert.equal(buildCompactNodeMeta({ service: 'web', ip: '10.0.0.1' }), '10.0.0.1');
  assert.equal(buildCompactNodeMeta({ bytes: 1024, connections: 2 }), '1.0 KB / 2 conn');
});

test('inline chips reserve overflow width without losing or mutating items', () => {
  const items = Object.freeze(['a', 'b', 'c']);
  assert.deepEqual(resolveVisibleInlineItems([], 200), { visibleItems: [], hiddenItems: [] });
  assert.deepEqual(resolveVisibleInlineItems(items, 0), { visibleItems: items, hiddenItems: [] });
  assert.deepEqual(resolveVisibleInlineItems(items, 119), { visibleItems: items, hiddenItems: [] });
  assert.deepEqual(resolveVisibleInlineItems(items, 120), { visibleItems: ['a'], hiddenItems: ['b', 'c'] });
  assert.deepEqual(resolveVisibleInlineItems(items, 184), { visibleItems: items, hiddenItems: [] });
  const longItems = ['a'.repeat(40), 'b'];
  assert.deepEqual(resolveVisibleInlineItems(longItems, 120), { visibleItems: [longItems[0]], hiddenItems: ['b'] });
});

test('demo topology is deterministic, independent, and uses millisecond windows', () => {
  const first = createDemoTopology(89999, 1000, 'unit');
  assert.deepEqual(first, createDemoTopology(89999, 1000, 'unit'));
  assert.deepEqual(first.window, { start: 89999, end: 90000, duration: 1000, total: 90000, total_count: 115 });
  assert.equal(first.center.label, 'Victim #unit');
  assert.equal(first.timeline.length, 90);
  assert.equal(first.timeline[1].timestamp_ms, 1000);
  assert.equal(first.timeline[89].timestamp_ms, 89000);
  assert.equal(first.top_talkers.length, 4);
  assert.equal(first.top_edges.length, 4);
  assert.ok(first.top_edges.every((edge) => edge.label === `${edge.source} -> ${edge.target}`));
  first.nodes[0].bytes = 0;
  assert.equal(createDemoTopology().nodes[0].bytes, 396000);
});

test('layout preserves original zones and does not mutate the universe', () => {
  const nodes = createDemoTopology().nodes.map((node) => Object.freeze(node));
  Object.freeze(nodes);
  const positions = buildPositions(nodes);
  assert.equal(positions.size, 6);
  assert.deepEqual(positions.get('172.20.1.34'), { x: 186, y: 115, w: 152, h: 46 });
  assert.deepEqual(positions.get('192.168.31.23'), { x: 186, y: 465, w: 152, h: 46 });
  assert.deepEqual(positions.get('10.10.0.10'), { x: 455, y: 290, w: 170, h: 54 });
  assert.deepEqual(positions.get('10.10.0.11'), { x: 665, y: 290, w: 170, h: 54 });
  assert.deepEqual(positions.get('8.8.8.8'), { x: 934, y: 115, w: 152, h: 46 });
  assert.deepEqual(buildPositions([...nodes].reverse()), positions);
  assert.equal(buildPositions([]).size, 0);
  assert.deepEqual(buildPositions([{ id: 'only', side: 'center' }]).get('only'), { x: 560, y: 290, w: 170, h: 54 });
});

test('dense and orbit layouts remain deterministic with finite dimensions', () => {
  const nodes = Array.from({ length: 72 }, (_, index) => ({
    id: String(index),
    side: ['left', 'right', 'center', 'orbit'][index % 4],
    bytes: index * 100,
    service: `service-${index % 3}`,
  }));
  const positions = buildPositions(nodes);
  assert.equal(positions.size, nodes.length);
  assert.deepEqual(buildPositions([...nodes].reverse()), positions);
  for (const position of positions.values()) {
    assert.ok(Object.values(position).every(Number.isFinite));
    assert.ok(position.w > 0 && position.h > 0);
  }
});

test('protocol filtering requires matching edges and endpoints, not dominant protocol', () => {
  const nodes = [
    { id: 'a', protocols: ['TCP'], kind: 'peer' },
    { id: 'b', protocols: ['TCP', 'HTTP'], kind: 'victim' },
    { id: 'c', protocols: ['UDP'], kind: 'peer' },
  ];
  const edges = [
    { id: 'ab', source: 'a', target: 'b', protocols: ['TCP'] },
    { id: 'ac', source: 'a', target: 'c', protocols: ['TCP'] },
    { id: 'dominant', source: 'a', target: 'b', dominant_proto: 'TCP' },
    { id: 'missing', source: 'a', target: 'missing', protocols: ['TCP'] },
  ];
  const topology = { nodes, edges, summary: { total_bytes: 12345 }, top_edges: edges };
  const before = structuredClone(topology);
  const universe = [...nodes, { id: 'ghost', protocols: ['DNS'] }];
  const unfiltered = filterTraffic(topology, null, new Set());
  assert.equal(unfiltered.nodes, nodes);
  assert.equal(unfiltered.edges, edges);
  const filtered = filterTraffic(topology, universe, new Set(['TCP']));
  assert.deepEqual(filtered.nodes, nodes.slice(0, 2));
  assert.deepEqual(filtered.edges, [edges[0]]);
  assert.deepEqual(filtered.universeProtocols, ['DNS', 'HTTP', 'TCP', 'UDP']);
  assert.deepEqual([...filtered.activeProtocolSet], ['HTTP', 'TCP', 'UDP']);
  assert.deepEqual(topology, before);
  assert.deepEqual(filterTraffic(null, null, new Set()).nodes, []);
  assert.deepEqual(filterTraffic(topology, [], new Set()).universeProtocols, []);
});

test('full-duration nodes retain positions across frames and protocol filters', () => {
  const universe = createDemoTopology().nodes;
  const first = filterTraffic({ nodes: universe.slice(0, 3) }, universe, new Set());
  const second = filterTraffic({ nodes: universe.slice(2) }, universe, new Set(['UDP']));
  const layout = buildPositions(universe);
  assert.notDeepEqual(first.nodes, second.nodes);
  assert.deepEqual(buildPositions(universe), layout);
  assert.ok(universe.some((node) => !second.nodes.some((active) => active.id === node.id)));
});

test('selection falls back to the first edge and victim, then first node, then null', () => {
  const nodes = [{ id: 'peer' }, { id: 'victim', kind: 'victim' }];
  const edges = [{ id: 'one' }, { id: 'two' }];
  assert.deepEqual(resolveTrafficSelection(nodes, edges, 'peer', 'two'), {
    selectedNode: nodes[0],
    selectedEdge: edges[1],
  });
  assert.deepEqual(resolveTrafficSelection(nodes, edges, 'gone', 'gone'), {
    selectedNode: nodes[1],
    selectedEdge: edges[0],
  });
  assert.deepEqual(resolveTrafficSelection(nodes.slice(0, 1), [], '', ''), {
    selectedNode: nodes[0],
    selectedEdge: null,
  });
  assert.deepEqual(resolveTrafficSelection([], [], '', ''), { selectedNode: null, selectedEdge: null });
});

test('curves preserve direction, traffic order and missing-endpoint handling', () => {
  const positions = new Map([
    ['a', { x: 100, y: 100, w: 100, h: 50 }],
    ['b', { x: 400, y: 100, w: 100, h: 50 }],
  ]);
  const edge = Object.freeze({ id: 'ab', source: 'a', target: 'b', direction: 'ingress', bytes: 1024 });
  const lines = buildLines(Object.freeze([edge]), positions);
  assert.equal(lines[0].path, 'M 150 100 C 220 61, 280 61, 350 100');
  assert.equal(lines[0].labelBaseX, 250);
  assert.equal(lines[0].tone.hard, '#597ef7');
  assert.deepEqual(buildLines([{ ...edge, target: 'missing' }], positions), []);
  const demo = createDemoTopology();
  const sorted = buildLines(demo.edges, buildPositions(demo.nodes));
  assert.deepEqual(
    sorted.map((line) => line.id),
    ['edge-e', 'edge-d', 'edge-b', 'edge-c', 'edge-a']
  );
  assert.ok(sorted.every((line) => !line.path.includes('NaN')));
  assert.equal(edgeTone('egress').hard, '#a3a3a3');
  assert.equal(edgeTone('internal').hard, '#737373');
  assert.match(edgeChipClass('ingress'), /text-geek-400/);
});

test('edge labels prioritize selection and bound fallback when collisions fill the canvas', () => {
  const lines = Array.from({ length: 8 }, (_, index) => ({
    id: String(index),
    bytes: index * 1024,
    labelBaseX: 100 + index * 120,
    labelBaseY: 100,
  }));
  const labels = resolveEdgeLabels(lines, new Map(), '0');
  assert.equal(labels.size, 6);
  assert.equal([...labels.keys()][0], '0');
  assert.ok(labels.has('0'));
  const coveringNode = new Map([['cover', { x: 560, y: 290, w: 2240, h: 1160 }]]);
  assert.equal(resolveEdgeLabels(lines, coveringNode, '').size, 0);
  const forced = resolveEdgeLabels(
    [{ id: 'selected', bytes: 1024, labelBaseX: -100, labelBaseY: 9999 }],
    coveringNode,
    'selected'
  );
  const label = forced.get('selected');
  assert.equal(label.text, '1.0 KB');
  assert.ok(label.x >= 12 && label.x + label.w <= VIEWBOX_WIDTH - 12);
  assert.ok(label.y >= 12 && label.y + label.h <= VIEWBOX_HEIGHT - 12);
});

test('viewport centers zoom-out and clamps zoom-in pan', () => {
  assert.deepEqual(getViewportMetrics(1), { width: 1120, height: 580 });
  assert.deepEqual(getViewportMetrics(2), { width: 560, height: 290 });
  assert.deepEqual(clampPanToViewport({ x: 999, y: -999 }, 1), { x: 0, y: 0 });
  assert.deepEqual(clampPanToViewport({ x: 999, y: -999 }, 0.8), { x: -140, y: -72.5 });
  assert.deepEqual(clampPanToViewport({ x: 999, y: -999 }, 2), { x: 560, y: 0 });
  assert.deepEqual(clampPanToViewport({ x: 123, y: 45 }, 2), { x: 123, y: 45 });
});
