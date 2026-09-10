import { DEFAULT_SLICE_MS } from './trafficPresentation.js';

export const createDemoTopology = (shift = 0, duration = DEFAULT_SLICE_MS, id = 'demo') => {
  const seed = shift % 12;
  const totalDuration = 90000;
  const nodes = [
    {
      id: '10.10.0.10',
      label: 'Victim 10.10.0.10',
      ip: '10.10.0.10',
      kind: 'victim',
      side: 'center',
      bytes: 396000,
      connections: 52,
      protocols: ['TCP', 'HTTP'],
    },
    {
      id: '10.10.0.11',
      label: 'Victim 10.10.0.11',
      ip: '10.10.0.11',
      kind: 'victim',
      side: 'center',
      bytes: 258000,
      connections: 33,
      protocols: ['TCP', 'HTTPS'],
    },
    {
      id: '172.20.1.34',
      label: 'Private 172.20.1.34',
      ip: '172.20.1.34',
      kind: 'peer',
      side: 'left',
      bytes: 202000 + seed * 11000,
      connections: 31,
      protocols: ['TCP', 'HTTP'],
    },
    {
      id: '192.168.31.23',
      label: 'Private 192.168.31.23',
      ip: '192.168.31.23',
      kind: 'peer',
      side: 'left',
      bytes: 131000,
      connections: 24,
      protocols: ['UDP', 'DNS'],
    },
    {
      id: '8.8.8.8',
      label: '8.8.8.8',
      ip: '8.8.8.8',
      kind: 'peer',
      side: 'right',
      bytes: 175000,
      connections: 29,
      protocols: ['UDP', 'DNS'],
    },
    {
      id: '104.26.6.171',
      label: '104.26.6.171',
      ip: '104.26.6.171',
      kind: 'peer',
      side: 'right',
      bytes: 102000 + seed * 6000,
      connections: 18,
      protocols: ['TCP', 'HTTPS'],
    },
  ];
  const edges = [
    {
      id: 'edge-a',
      source: '172.20.1.34',
      target: '10.10.0.10',
      direction: 'ingress',
      bytes: 202000 + seed * 11000,
      packets: 98,
      connections: 31,
      dominant_proto: 'TCP',
      dominant_app: 'HTTP',
      intensity: 1,
    },
    {
      id: 'edge-b',
      source: '192.168.31.23',
      target: '10.10.0.11',
      direction: 'ingress',
      bytes: 131000,
      packets: 77,
      connections: 24,
      dominant_proto: 'UDP',
      dominant_app: 'DNS',
      intensity: 0.64,
    },
    {
      id: 'edge-c',
      source: '10.10.0.10',
      target: '8.8.8.8',
      direction: 'egress',
      bytes: 175000,
      packets: 94,
      connections: 29,
      dominant_proto: 'UDP',
      dominant_app: 'DNS',
      intensity: 0.84,
    },
    {
      id: 'edge-d',
      source: '10.10.0.11',
      target: '104.26.6.171',
      direction: 'egress',
      bytes: 102000 + seed * 6000,
      packets: 56,
      connections: 18,
      dominant_proto: 'TCP',
      dominant_app: 'HTTPS',
      intensity: 0.48,
    },
    {
      id: 'edge-e',
      source: '10.10.0.10',
      target: '10.10.0.11',
      direction: 'internal',
      bytes: 82000,
      packets: 44,
      connections: 13,
      dominant_proto: 'TCP',
      dominant_app: 'PROXY',
      intensity: 0.36,
    },
  ];
  const timeline = Array.from({ length: 90 }).map((_, index) => {
    const timestampMs = index * 1000;
    const bytes = Math.round(16000 + Math.abs(Math.sin((index + seed) * 0.31)) * 54000 + ((index + seed) % 7) * 2100);
    return {
      timestamp_ms: timestampMs,
      bytes,
      packets: 7 + ((index + seed) % 11),
      ingress_bytes: Math.round(bytes * 0.54),
      egress_bytes: Math.round(bytes * 0.35),
    };
  });
  return {
    window: {
      start: shift,
      end: Math.min(totalDuration, shift + duration),
      duration,
      total: totalDuration,
      total_count: 115,
    },
    total_duration: totalDuration,
    available_slices: [1000, 5000, 15000, 30000, 60000, totalDuration],
    center: { label: `Victim #${id}`, ips: ['10.10.0.10', '10.10.0.11'], exposed: ['tcp://43.155.12.20:24001'] },
    summary: {
      total_bytes: 690000 + seed * 18000,
      ingress_bytes: 349000,
      egress_bytes: 256000,
      internal_bytes: 82000,
      visible_edges: 5,
      visible_nodes: 6,
    },
    nodes,
    edges,
    timeline,
    top_talkers: nodes
      .filter((node) => node.kind === 'peer')
      .sort((a, b) => b.bytes - a.bytes)
      .slice(0, 4),
    top_edges: edges
      .slice()
      .sort((a, b) => b.bytes - a.bytes)
      .slice(0, 4)
      .map((edge) => ({ ...edge, label: `${edge.source} -> ${edge.target}` })),
  };
};
