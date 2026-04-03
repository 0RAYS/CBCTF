import { useEffect, useMemo, useState } from 'react';
import { motion } from 'motion/react';
import {
  IconActivity,
  IconArrowDownRight,
  IconArrowUpRight,
  IconDownload,
  IconRefresh,
  IconRoute,
  IconTopologyComplex,
} from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { getContestTeamTraffic, downloadContainerTraffic } from '../../../../api/admin/contest.js';
import { downloadVictimTraffic } from '../../../../api/admin/victims.js';
import { Button, Card, Chip, Modal } from '../../../../components/common';
import { downloadBlobResponse } from '../../../../utils/fileDownload';
import { toast } from '../../../../utils/toast';

const W = 920;
const H = 620;
const PW = 156;
const PH = 58;
const CW = 188;
const CH = 72;

const fmtBytes = (value) => {
  const bytes = Number(value || 0);
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(1)} GB`;
};

const tone = (direction, protocol = '') => {
  const base =
    direction === 'ingress'
      ? ['rgba(89,126,247,0.28)', '#597ef7']
      : direction === 'egress'
        ? ['rgba(148,163,184,0.24)', '#94a3b8']
        : ['rgba(115,115,115,0.24)', '#a3a3a3'];
  const proto = protocol.toUpperCase() === 'UDP' ? '#94a3b8' : protocol.toUpperCase() === 'TCP' ? '#597ef7' : base[1];
  return { soft: base[0], hard: proto };
};

const demoTopology = (shift = 0, duration = 15, id = 'demo') => {
  const s = shift % 10;
  const nodes = [
    {
      id: '10.10.0.10',
      label: 'Victim 10.10.0.10',
      ip: '10.10.0.10',
      kind: 'victim',
      side: 'center',
      bytes: 410000,
      connections: 58,
      protocols: ['TCP', 'HTTP'],
    },
    {
      id: '10.10.0.11',
      label: 'Victim 10.10.0.11',
      ip: '10.10.0.11',
      kind: 'victim',
      side: 'center',
      bytes: 252000,
      connections: 37,
      protocols: ['TCP', 'HTTPS'],
    },
    {
      id: '172.20.1.34',
      label: 'Private 172.20.1.34',
      ip: '172.20.1.34',
      kind: 'peer',
      side: 'left',
      bytes: 220000 + s * 12000,
      connections: 35,
      protocols: ['TCP', 'HTTP'],
    },
    {
      id: '192.168.31.23',
      label: 'Private 192.168.31.23',
      ip: '192.168.31.23',
      kind: 'peer',
      side: 'left',
      bytes: 124000,
      connections: 22,
      protocols: ['UDP', 'DNS'],
    },
    {
      id: '8.8.8.8',
      label: '8.8.8.8',
      ip: '8.8.8.8',
      kind: 'peer',
      side: 'right',
      bytes: 182000,
      connections: 30,
      protocols: ['UDP', 'DNS'],
    },
    {
      id: '104.26.6.171',
      label: '104.26.6.171',
      ip: '104.26.6.171',
      kind: 'peer',
      side: 'right',
      bytes: 98000 + s * 7000,
      connections: 19,
      protocols: ['TCP', 'HTTPS'],
    },
  ];
  const edges = [
    {
      id: 'a',
      source: '172.20.1.34',
      target: '10.10.0.10',
      direction: 'ingress',
      bytes: 220000 + s * 12000,
      packets: 118,
      connections: 35,
      dominant_proto: 'TCP',
      dominant_app: 'HTTP',
      intensity: 1,
    },
    {
      id: 'b',
      source: '192.168.31.23',
      target: '10.10.0.11',
      direction: 'ingress',
      bytes: 124000,
      packets: 76,
      connections: 22,
      dominant_proto: 'UDP',
      dominant_app: 'DNS',
      intensity: 0.55,
    },
    {
      id: 'c',
      source: '10.10.0.10',
      target: '8.8.8.8',
      direction: 'egress',
      bytes: 182000,
      packets: 101,
      connections: 30,
      dominant_proto: 'UDP',
      dominant_app: 'DNS',
      intensity: 0.8,
    },
    {
      id: 'd',
      source: '10.10.0.11',
      target: '104.26.6.171',
      direction: 'egress',
      bytes: 98000 + s * 7000,
      packets: 62,
      connections: 19,
      dominant_proto: 'TCP',
      dominant_app: 'HTTPS',
      intensity: 0.45,
    },
    {
      id: 'e',
      source: '10.10.0.10',
      target: '10.10.0.11',
      direction: 'internal',
      bytes: 86000,
      packets: 45,
      connections: 14,
      dominant_proto: 'TCP',
      dominant_app: 'PROXY',
      intensity: 0.32,
    },
  ];
  const timeline = Array.from({ length: 18 }).map((_, i) => {
    const bytes = Math.round(18000 + Math.abs(Math.sin((i + s) * 0.65)) * 64000);
    const ingress = Math.round(bytes * (0.45 + Math.abs(Math.cos(i * 0.4)) * 0.22));
    return {
      second: i,
      bytes,
      packets: 8 + ((i + s) % 12),
      ingress_bytes: ingress,
      egress_bytes: Math.max(0, bytes - ingress),
    };
  });
  return {
    window: { start: shift, end: Math.min(90, shift + duration), duration, total: 90, total_count: 120 },
    total_duration: 90,
    available_slices: [5, 15, 30, 60, 90],
    center: { label: `Victim #${id}`, ips: ['10.10.0.10', '10.10.0.11'], exposed: ['tcp://43.155.12.20:24001'] },
    summary: {
      total_bytes: 710000 + s * 20000,
      ingress_bytes: 362000,
      egress_bytes: 262000,
      internal_bytes: 86000,
      visible_edges: 5,
      visible_nodes: 6,
      peak_second: 9,
      peak_bytes: 82000,
    },
    nodes,
    edges,
    timeline,
    top_talkers: nodes
      .filter((n) => n.kind === 'peer')
      .sort((a, b) => b.bytes - a.bytes)
      .slice(0, 4)
      .map((n) => ({ ...n, direction: n.side === 'left' ? 'ingress' : 'egress' })),
    top_edges: edges
      .slice()
      .sort((a, b) => b.bytes - a.bytes)
      .slice(0, 4)
      .map((e) => ({ label: `${e.source} -> ${e.target}`, ...e })),
  };
};

function layoutGraph(nodes, edges) {
  const pos = new Map();
  const left = nodes.filter((n) => n.side === 'left');
  const right = nodes.filter((n) => n.side === 'right');
  const center = nodes.filter((n) => n.side === 'center');
  left.forEach((node, index) => pos.set(node.id, { x: 126, y: (H / (left.length + 1)) * (index + 1), w: PW, h: PH }));
  right.forEach((node, index) =>
    pos.set(node.id, { x: W - 126, y: (H / (right.length + 1)) * (index + 1), w: PW, h: PH })
  );
  if (center.length === 1) pos.set(center[0].id, { x: W / 2, y: H / 2, w: CW, h: CH });
  if (center.length > 1) {
    const startX = W / 2 - ((center.length - 1) * 96) / 2;
    center.forEach((node, index) =>
      pos.set(node.id, { x: startX + index * 96, y: H / 2 + (index % 2 === 0 ? -40 : 40), w: 164, h: 62 })
    );
  }
  return edges
    .map((edge, index) => {
      const s = pos.get(edge.source);
      const t = pos.get(edge.target);
      if (!s || !t) return null;
      const ax = edge.direction === 'ingress' ? s.x + s.w / 2 : edge.direction === 'egress' ? s.x : s.x;
      const bx = edge.direction === 'ingress' ? t.x : edge.direction === 'egress' ? t.x - t.w / 2 : t.x;
      const ay = s.y;
      const by = t.y;
      const dx = bx - ax;
      const bend = edge.direction === 'internal' ? -88 : edge.direction === 'ingress' ? -92 : 92;
      const c1x = ax + dx * 0.32;
      const c2x = bx - dx * 0.32;
      const path = `M ${ax} ${ay} C ${c1x} ${ay + bend}, ${c2x} ${by + bend}, ${bx} ${by}`;
      return {
        ...edge,
        path,
        lx: (ax + bx) / 2,
        ly: (ay + by) / 2 + bend * 0.4,
        tone: tone(edge.direction, edge.dominant_proto),
        delay: index * 0.16,
      };
    })
    .filter(Boolean);
}

function TrafficGraphModal({ isOpen, onClose, container, contestId, teamId, fetchTraffic: customFetchTraffic }) {
  const { t, i18n } = useTranslation();
  const [topology, setTopology] = useState(null);
  const [shift, setShift] = useState(0);
  const [slice, setSlice] = useState(15);
  const [loading, setLoading] = useState(false);
  const [demoMode, setDemoMode] = useState(false);
  const [edgeId, setEdgeId] = useState('');
  const [nodeId, setNodeId] = useState('');

  const fetchData = async (nextShift = shift, nextSlice = slice, forceLive = false) => {
    if (!container?.id) return;
    setLoading(true);
    try {
      const response = customFetchTraffic
        ? await customFetchTraffic(container, { time_shift: nextShift, duration: nextSlice })
        : await getContestTeamTraffic(contestId, teamId, container.id, { time_shift: nextShift, duration: nextSlice });
      if (response.code !== 200) throw new Error(t('admin.contests.trafficGraph.toast.fetchFailed'));
      setTopology(response.data);
      setDemoMode(false);
    } catch {
      if (!forceLive) toast.warning({ description: t('admin.contests.trafficGraph.toast.demoFallback') });
      setTopology(demoTopology(nextShift, nextSlice, container.id));
      setDemoMode(true);
    } finally {
      setLoading(false);
      setEdgeId('');
      setNodeId('');
    }
  };

  useEffect(() => {
    if (!isOpen || !container?.id) return;
    fetchData(shift, slice);
  }, [isOpen, container?.id]);

  useEffect(() => {
    if (!isOpen || !container?.id) return;
    fetchData(shift, slice, demoMode);
  }, [shift, slice]);

  const nodes = topology?.nodes || [];
  const edges = topology?.edges || [];
  const lines = useMemo(() => layoutGraph(nodes, edges), [nodes, edges]);
  const positions = useMemo(() => {
    const map = new Map();
    nodes
      .filter((n) => n.side === 'left')
      .forEach((n, i, arr) => map.set(n.id, { x: 126, y: (H / (arr.length + 1)) * (i + 1), w: PW, h: PH }));
    nodes
      .filter((n) => n.side === 'right')
      .forEach((n, i, arr) => map.set(n.id, { x: W - 126, y: (H / (arr.length + 1)) * (i + 1), w: PW, h: PH }));
    const center = nodes.filter((n) => n.side === 'center');
    if (center.length === 1) map.set(center[0].id, { x: W / 2, y: H / 2, w: CW, h: CH });
    if (center.length > 1) {
      const startX = W / 2 - ((center.length - 1) * 96) / 2;
      center.forEach((n, i) =>
        map.set(n.id, { x: startX + i * 96, y: H / 2 + (i % 2 === 0 ? -40 : 40), w: 164, h: 62 })
      );
    }
    return map;
  }, [nodes]);

  const selectedEdge = useMemo(() => edges.find((e) => e.id === edgeId) || edges[0] || null, [edges, edgeId]);
  const selectedNode = useMemo(
    () => nodes.find((n) => n.id === nodeId) || nodes.find((n) => n.kind === 'victim') || nodes[0] || null,
    [nodes, nodeId]
  );
  const summary = topology?.summary || {};
  const windowInfo = topology?.window || { start: 0, end: 0, duration: slice, total: 0 };
  const timeline = topology?.timeline || [];
  const peak = Math.max(...timeline.map((i) => i.bytes || 0), 1);

  const downloadTraffic = async () => {
    if (!container?.id) return;
    try {
      const response =
        contestId && teamId
          ? await downloadContainerTraffic(contestId, teamId, container.id)
          : await downloadVictimTraffic(container.id);
      if (response.headers?.['file'] === 'true') downloadBlobResponse(response, `traffic_${container.id}.zip`);
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teamContainers.toast.downloadTrafficFailed') });
    }
  };

  if (!isOpen) return null;

  return (
    <>
      <style>
        {`
          @keyframes traffic-flow { from { stroke-dashoffset: 0; } to { stroke-dashoffset: -120; } }
          @keyframes traffic-pulse { 0%,100% { opacity: .35; transform: scale(1); } 50% { opacity: .92; transform: scale(1.06); } }
          .traffic-line { stroke-dasharray: 12 10; animation: traffic-flow 10s linear infinite; }
          .traffic-pulse { animation: traffic-pulse 4s ease-in-out infinite; transform-origin: center; }
        `}
      </style>
      <Modal
        isOpen={isOpen}
        onClose={onClose}
        title={t('admin.contests.trafficGraph.title')}
        size="2xl"
        className="!bg-neutral-900/95 !border-neutral-600"
      >
        <div className="space-y-5 text-neutral-100">
          <div className="overflow-hidden rounded-2xl border border-neutral-600 bg-[radial-gradient(circle_at_top_left,_rgba(89,126,247,0.12),_transparent_30%),linear-gradient(160deg,_rgba(10,10,10,0.98),_rgba(23,23,23,0.96)_48%,_rgba(15,23,42,0.88))] p-5">
            <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-xs uppercase tracking-[0.3em] text-geek-400/80">
                  <IconTopologyComplex size={16} />
                  <span>{t('admin.contests.trafficGraph.hero.kicker')}</span>
                </div>
                <div className="text-2xl font-['Maple_UI'] text-white">
                  {topology?.center?.label || `Victim #${container?.id || '-'}`}
                </div>
                <div className="text-sm text-neutral-300">
                  {t('admin.contests.trafficGraph.hero.subtitle', {
                    challenge: container?.challenge || container?.contest_challenge_name || `#${container?.id}`,
                  })}
                </div>
                <div className="flex flex-wrap gap-2">
                  {(topology?.center?.ips || []).map((ip) => (
                    <Chip
                      key={ip}
                      label={ip}
                      variant="tag"
                      size="sm"
                      colorClass="border-geek-400/30 bg-geek-400/10 text-geek-400"
                    />
                  ))}
                  {demoMode ? (
                    <Chip
                      label={t('admin.contests.trafficGraph.hero.demoMode')}
                      variant="tag"
                      size="sm"
                      colorClass="border-neutral-400/30 bg-neutral-400/10 text-neutral-300"
                    />
                  ) : null}
                </div>
              </div>
              <div className="flex flex-wrap items-center gap-3">
                <div className="rounded-xl border border-neutral-600 bg-black/20 px-3 py-2 text-sm">
                  <div className="mb-1 text-xs text-neutral-400">
                    {t('admin.contests.trafficGraph.controls.timeShift')}
                  </div>
                  <input
                    type="range"
                    min="0"
                    max={Math.max(0, windowInfo.total - windowInfo.duration)}
                    value={shift}
                    onChange={(e) => setShift(Number(e.target.value))}
                    className="w-44 accent-geek-400"
                  />
                  <div className="mt-1 font-mono text-geek-400">
                    {t('admin.contests.trafficGraph.hero.windowAt', { start: windowInfo.start, end: windowInfo.end })}
                  </div>
                </div>
                <div className="rounded-xl border border-neutral-600 bg-black/20 px-3 py-2 text-sm">
                  <div className="mb-1 text-xs text-neutral-400">
                    {t('admin.contests.trafficGraph.controls.timeSlice')}
                  </div>
                  <div className="flex gap-2">
                    {(topology?.available_slices || [5, 15, 30, 60]).map((item) => (
                      <button
                        key={item}
                        type="button"
                        onClick={() => setSlice(item)}
                        className={`rounded-md border px-2 py-1 text-xs font-mono ${slice === item ? 'border-geek-400/60 bg-geek-400/15 text-geek-400' : 'border-neutral-600 bg-black/20 text-neutral-300'}`}
                      >
                        {item}s
                      </button>
                    ))}
                  </div>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  className="!text-neutral-300 hover:!text-neutral-100"
                  onClick={() => fetchData(shift, slice, true)}
                >
                  <IconRefresh size={18} />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="!text-neutral-300 hover:!text-neutral-100"
                  onClick={downloadTraffic}
                >
                  <IconDownload size={18} />
                </Button>
              </div>
            </div>
            <div className="mt-4 grid grid-cols-1 gap-3 md:grid-cols-4">
              {[
                {
                  icon: <IconActivity size={18} className="text-geek-400" />,
                  label: t('admin.contests.trafficGraph.stats.totalTraffic'),
                  value: fmtBytes(summary.total_bytes),
                  valueClass: 'text-geek-400',
                  box: 'bg-geek-400/15',
                },
                {
                  icon: <IconArrowDownRight size={18} className="text-neutral-300" />,
                  label: t('admin.contests.trafficGraph.stats.ingress'),
                  value: fmtBytes(summary.ingress_bytes),
                  valueClass: 'text-neutral-100',
                  box: 'bg-neutral-400/15',
                },
                {
                  icon: <IconArrowUpRight size={18} className="text-neutral-300" />,
                  label: t('admin.contests.trafficGraph.stats.egress'),
                  value: fmtBytes(summary.egress_bytes),
                  valueClass: 'text-neutral-100',
                  box: 'bg-neutral-400/15',
                },
                {
                  icon: <IconRoute size={18} className="text-neutral-300" />,
                  label: t('admin.contests.trafficGraph.stats.peers'),
                  value: summary.visible_nodes || 0,
                  valueClass: 'text-neutral-100',
                  box: 'bg-neutral-400/15',
                },
              ].map((item) => (
                <div key={item.label} className="rounded-xl border border-neutral-600 bg-black/20 p-4">
                  <div className="flex items-center gap-3">
                    <div className={`flex h-10 w-10 items-center justify-center rounded-xl ${item.box}`}>
                      {item.icon}
                    </div>
                    <div>
                      <div className="text-xs font-mono text-neutral-400">{item.label}</div>
                      <div className={`font-mono text-xl ${item.valueClass}`}>{item.value}</div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {demoMode ? (
            <div className="rounded-xl border border-neutral-500/30 bg-neutral-500/10 px-4 py-3 text-sm text-neutral-300">
              {t('admin.contests.trafficGraph.demoNotice')}
            </div>
          ) : null}

          <div className="grid grid-cols-1 gap-5 xl:grid-cols-[minmax(0,1.5fr)_380px]">
            <Card padding="none" className="overflow-hidden rounded-2xl border-neutral-600 bg-neutral-900">
              <div className="flex items-center justify-between border-b border-neutral-600 px-5 py-4">
                <div>
                  <div className="text-sm font-mono text-neutral-300">
                    {t('admin.contests.trafficGraph.canvas.title')}
                  </div>
                  <div className="mt-1 text-xs text-neutral-500">
                    {t('admin.contests.trafficGraph.canvas.subtitle')}
                  </div>
                </div>
                <div className="flex gap-2 text-xs">
                  <Chip
                    label={t('admin.contests.trafficGraph.legend.ingress')}
                    variant="tag"
                    size="sm"
                    colorClass="border-geek-400/30 bg-geek-400/10 text-geek-400"
                  />
                  <Chip
                    label={t('admin.contests.trafficGraph.legend.egress')}
                    variant="tag"
                    size="sm"
                    colorClass="border-neutral-400/30 bg-neutral-400/10 text-neutral-300"
                  />
                  <Chip
                    label={t('admin.contests.trafficGraph.legend.internal')}
                    variant="tag"
                    size="sm"
                    colorClass="border-neutral-500/30 bg-neutral-500/10 text-neutral-300"
                  />
                </div>
              </div>
              <div className="overflow-x-auto p-4">
                <div className="min-w-[920px]">
                  <svg viewBox={`0 0 ${W} ${H}`} className="h-[620px] w-full">
                    <g opacity="0.18">
                      <path d="M 50 125 H 870" stroke="#404040" strokeWidth="1" />
                      <path d="M 50 310 H 870" stroke="#404040" strokeWidth="1" />
                      <path d="M 50 495 H 870" stroke="#404040" strokeWidth="1" />
                    </g>

                    {lines.map((edge) => (
                      <g key={edge.id}>
                        <path
                          d={edge.path}
                          fill="none"
                          stroke={edge.tone.hard}
                          strokeWidth={9}
                          opacity={selectedEdge?.id === edge.id ? 0.26 : 0.12 + edge.intensity * 0.14}
                        />
                        <path
                          d={edge.path}
                          fill="none"
                          stroke={edge.tone.soft}
                          strokeWidth={2 + edge.intensity * 5}
                          className="traffic-line"
                          opacity={0.55 + edge.intensity * 0.25}
                          style={{ animationDelay: `${edge.delay}s` }}
                          onClick={() => setEdgeId(edge.id)}
                        />
                        <g
                          transform={`translate(${edge.lx}, ${edge.ly})`}
                          className="cursor-pointer"
                          onClick={() => setEdgeId(edge.id)}
                        >
                          <rect
                            x={-42}
                            y={-12}
                            width={84}
                            height={24}
                            rx={12}
                            fill="rgba(10,10,10,.92)"
                            stroke={selectedEdge?.id === edge.id ? edge.tone.hard : 'rgba(115,115,115,.35)'}
                          />
                          <text x="0" y="5" textAnchor="middle" fill="#e2e8f0" fontSize="11" fontFamily="Maple Mono">
                            {fmtBytes(edge.bytes)}
                          </text>
                        </g>
                      </g>
                    ))}

                    {nodes.map((node) => {
                      const pos = positions.get(node.id);
                      if (!pos) return null;
                      const active = selectedNode?.id === node.id;
                      const x = pos.x - pos.w / 2;
                      const y = pos.y - pos.h / 2;
                      const fill = node.kind === 'victim' ? 'rgba(89,126,247,.14)' : 'rgba(38,38,38,.9)';
                      const stroke = node.kind === 'victim' ? '#597ef7' : '#737373';
                      return (
                        <g
                          key={node.id}
                          transform={`translate(${x}, ${y})`}
                          className="cursor-pointer"
                          onClick={() => setNodeId(node.id)}
                        >
                          <rect
                            className={node.kind === 'victim' ? 'traffic-pulse' : ''}
                            x={active ? -4 : 0}
                            y={active ? -4 : 0}
                            width={active ? pos.w + 8 : pos.w}
                            height={active ? pos.h + 8 : pos.h}
                            rx="18"
                            fill={fill}
                            stroke={active ? '#d4d4d4' : stroke}
                            strokeWidth={active ? 2 : 1.2}
                          />
                          <circle
                            cx="22"
                            cy={pos.h / 2}
                            r={node.kind === 'victim' ? 13 : 11}
                            fill={fill}
                            stroke={stroke}
                          />
                          <text x="48" y="24" fill="#f8fafc" fontSize="12" fontFamily="Maple Mono">
                            {node.label.length > 18 ? `${node.label.slice(0, 18)}...` : node.label}
                          </text>
                          <text x="48" y="42" fill="#94a3b8" fontSize="11" fontFamily="Maple Mono">
                            {node.ip}
                          </text>
                          <text
                            x={pos.w - 10}
                            y="24"
                            textAnchor="end"
                            fill="#cbd5e1"
                            fontSize="11"
                            fontFamily="Maple Mono"
                          >
                            {fmtBytes(node.bytes)}
                          </text>
                          <text
                            x={pos.w - 10}
                            y="42"
                            textAnchor="end"
                            fill="#64748b"
                            fontSize="10"
                            fontFamily="Maple Mono"
                          >
                            {node.connections} conn
                          </text>
                        </g>
                      );
                    })}
                  </svg>
                </div>
              </div>
            </Card>

            <div className="space-y-4">
              <Card className="rounded-2xl border-neutral-600 bg-neutral-900">
                <div className="flex items-center justify-between">
                  <div>
                    <div className="text-sm font-mono text-neutral-300">
                      {t('admin.contests.trafficGraph.panel.selectedFlow')}
                    </div>
                    <div className="mt-1 text-xs text-neutral-500">
                      {t('admin.contests.trafficGraph.panel.selectedFlowHint')}
                    </div>
                  </div>
                  {selectedEdge ? (
                    <Chip
                      label={t(`admin.contests.trafficGraph.direction.${selectedEdge.direction || 'internal'}`)}
                      variant="tag"
                      size="sm"
                      colorClass={
                        selectedEdge.direction === 'ingress'
                          ? 'border-geek-400/30 bg-geek-400/10 text-geek-400'
                          : selectedEdge.direction === 'egress'
                            ? 'border-neutral-500/30 bg-neutral-500/10 text-neutral-300'
                            : 'border-neutral-600 bg-black/20 text-neutral-300'
                      }
                    />
                  ) : null}
                </div>
                {selectedEdge ? (
                  <div className="mt-4 space-y-3">
                    <div className="rounded-xl border border-neutral-600 bg-black/20 p-3">
                      <div className="text-xs text-neutral-500">{t('admin.contests.trafficGraph.panel.edgePath')}</div>
                      <div className="mt-2 break-all font-mono text-sm text-white">
                        {selectedEdge.source} {'->'} {selectedEdge.target}
                      </div>
                    </div>
                    <div className="grid grid-cols-2 gap-3">
                      <div className="rounded-xl border border-neutral-600 bg-black/20 p-3">
                        <div className="text-xs text-neutral-500">
                          {t('admin.contests.trafficGraph.panel.edgeBytes')}
                        </div>
                        <div className="mt-1 font-mono text-lg text-geek-400">{fmtBytes(selectedEdge.bytes)}</div>
                      </div>
                      <div className="rounded-xl border border-neutral-600 bg-black/20 p-3">
                        <div className="text-xs text-neutral-500">
                          {t('admin.contests.trafficGraph.panel.edgePackets')}
                        </div>
                        <div className="mt-1 font-mono text-lg text-slate-100">{selectedEdge.packets}</div>
                      </div>
                    </div>
                    <div className="flex flex-wrap gap-2">
                      {(selectedEdge.protocols && selectedEdge.protocols.length > 0
                        ? selectedEdge.protocols
                        : [selectedEdge.dominant_proto, selectedEdge.dominant_app]
                      )
                        .filter(Boolean)
                        .map((item) => (
                          <Chip key={item} label={item} size="sm" colorClass="bg-neutral-400/10 text-neutral-300" />
                        ))}
                    </div>
                  </div>
                ) : (
                  <div className="mt-4 rounded-xl border border-dashed border-neutral-600 px-4 py-5 text-sm text-neutral-400">
                    {t('admin.contests.trafficGraph.empty')}
                  </div>
                )}
              </Card>

              <Card className="rounded-2xl border-neutral-600 bg-neutral-900">
                <div className="text-sm font-mono text-neutral-300">
                  {t('admin.contests.trafficGraph.panel.selectedNode')}
                </div>
                {selectedNode ? (
                  <div className="mt-4 space-y-3">
                    <div className="rounded-xl border border-neutral-600 bg-black/20 p-3">
                      <div className="text-xs text-neutral-500">{t('admin.contests.trafficGraph.panel.nodeLabel')}</div>
                      <div className="mt-1 font-mono text-sm text-white">{selectedNode.label}</div>
                      <div className="mt-2 text-xs text-neutral-400">{selectedNode.ip}</div>
                    </div>
                    <div className="grid grid-cols-2 gap-3">
                      <div className="rounded-xl border border-neutral-600 bg-black/20 p-3">
                        <div className="text-xs text-neutral-500">
                          {t('admin.contests.trafficGraph.panel.nodeTraffic')}
                        </div>
                        <div className="mt-1 font-mono text-lg text-geek-400">{fmtBytes(selectedNode.bytes)}</div>
                      </div>
                      <div className="rounded-xl border border-neutral-600 bg-black/20 p-3">
                        <div className="text-xs text-neutral-500">
                          {t('admin.contests.trafficGraph.panel.nodeConnections')}
                        </div>
                        <div className="mt-1 font-mono text-lg text-slate-100">{selectedNode.connections}</div>
                      </div>
                    </div>
                    <div className="flex flex-wrap gap-2">
                      {(selectedNode.protocols || []).map((item) => (
                        <Chip key={item} label={item} size="sm" colorClass="bg-neutral-400/10 text-neutral-300" />
                      ))}
                    </div>
                  </div>
                ) : null}
              </Card>

              <Card className="rounded-2xl border-neutral-600 bg-neutral-900">
                <div className="text-sm font-mono text-neutral-300">
                  {t('admin.contests.trafficGraph.timeline.title')}
                </div>
                <div className="mt-4 space-y-2">
                  {timeline.length === 0 ? (
                    <div className="rounded-xl border border-dashed border-neutral-600 px-4 py-5 text-sm text-neutral-400">
                      {t('admin.contests.trafficGraph.empty')}
                    </div>
                  ) : (
                    timeline.map((bucket) => {
                      const ratio = Math.max((bucket.bytes || 0) / peak, 0.02);
                      const active =
                        bucket.second >= windowInfo.start &&
                        bucket.second <= Math.max(windowInfo.start, windowInfo.end - 1);
                      return (
                        <button
                          key={bucket.second}
                          type="button"
                          onClick={() => setShift(bucket.second)}
                          className={`w-full rounded-xl border px-3 py-2 text-left ${active ? 'border-geek-400/30 bg-geek-400/10' : 'border-neutral-600 bg-black/20'}`}
                        >
                          <div className="flex items-center justify-between text-xs text-neutral-400">
                            <span className="font-mono">
                              {t('admin.contests.trafficGraph.timeline.second', { value: bucket.second })}
                            </span>
                            <span className="font-mono">{fmtBytes(bucket.bytes)}</span>
                          </div>
                          <div className="mt-2 h-2 overflow-hidden rounded-full bg-neutral-700/50">
                            <motion.div
                              className="h-full rounded-full bg-geek-400"
                              initial={{ width: 0 }}
                              animate={{ width: `${ratio * 100}%` }}
                              transition={{ duration: 0.45, ease: 'easeOut' }}
                            />
                          </div>
                        </button>
                      );
                    })
                  )}
                </div>
              </Card>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-5 xl:grid-cols-2">
            <Card className="rounded-2xl border-neutral-600 bg-neutral-900">
              <div className="text-sm font-mono text-neutral-300">
                {t('admin.contests.trafficGraph.rankings.topTalkers')}
              </div>
              <div className="mt-4 space-y-2">
                {(topology?.top_talkers || []).map((item, index) => (
                  <div
                    key={`${item.ip}-${index}`}
                    className="flex items-center justify-between rounded-xl border border-neutral-600 bg-black/20 px-3 py-3"
                  >
                    <div>
                      <div className="font-mono text-sm text-white">{item.label}</div>
                      <div className="mt-1 text-xs text-neutral-500">{item.ip}</div>
                    </div>
                    <div className="text-right">
                      <div className="font-mono text-sm text-geek-400">{fmtBytes(item.bytes)}</div>
                      <div className="mt-1 text-xs text-neutral-400">{item.connections} conn</div>
                    </div>
                  </div>
                ))}
              </div>
            </Card>

            <Card className="rounded-2xl border-neutral-600 bg-neutral-900">
              <div className="text-sm font-mono text-neutral-300">
                {t('admin.contests.trafficGraph.rankings.topEdges')}
              </div>
              <div className="mt-4 space-y-2">
                {(topology?.top_edges || []).map((item, index) => (
                  <div
                    key={`${item.id || item.ip}-${index}`}
                    className="flex items-center justify-between rounded-xl border border-neutral-600 bg-black/20 px-3 py-3"
                  >
                    <div>
                      <div className="font-mono text-sm text-white">{item.label}</div>
                      <div className="mt-1 text-xs text-neutral-500">
                        {t(`admin.contests.trafficGraph.direction.${item.direction || 'internal'}`)}
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="font-mono text-sm text-geek-400">{fmtBytes(item.bytes)}</div>
                      <div className="mt-1 text-xs text-neutral-400">{item.connections} conn</div>
                    </div>
                  </div>
                ))}
              </div>
            </Card>
          </div>

          <div className="flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-neutral-600 bg-black/20 px-4 py-3 text-xs text-neutral-400">
            <div className="flex flex-wrap gap-4 font-mono">
              <span>
                {t('admin.contests.trafficGraph.footer.window', { start: windowInfo.start, end: windowInfo.end })}
              </span>
              <span>
                {t('admin.contests.trafficGraph.footer.connectionCount', { count: summary.visible_edges || 0 })}
              </span>
              <span>{t('admin.contests.trafficGraph.footer.ipCount', { count: summary.visible_nodes || 0 })}</span>
              <span>
                {t('admin.contests.trafficGraph.footer.maxDuration', { count: topology?.total_duration || 0 })}
              </span>
            </div>
            <div className="font-mono">
              {t('admin.contests.trafficGraph.footer.updatedAt', {
                time: new Date().toLocaleString(i18n.language || 'en-US'),
              })}
            </div>
          </div>
        </div>
      </Modal>
    </>
  );
}

export default TrafficGraphModal;
