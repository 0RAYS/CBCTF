import { useEffect, useMemo, useRef, useState } from 'react';
import { motion } from 'motion/react';
import {
  IconActivity,
  IconArrowDownRight,
  IconArrowUpRight,
  IconDownload,
  IconPlayerPause,
  IconPlayerPlay,
  IconPlayerTrackNext,
  IconPlayerTrackPrev,
  IconRefresh,
  IconRoute,
  IconZoomIn,
  IconZoomOut,
} from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { getContestTeamTraffic, downloadContainerTraffic } from '../../../../api/admin/contest.js';
import { downloadVictimTraffic } from '../../../../api/admin/victims.js';
import { Button, Card, Chip, Input, Modal } from '../../../../components/common';
import { downloadBlobResponse } from '../../../../utils/fileDownload';
import { toast } from '../../../../utils/toast';

const VIEWBOX_WIDTH = 1120;
const VIEWBOX_HEIGHT = 580;
const CENTER_X = VIEWBOX_WIDTH / 2;
const CENTER_Y = VIEWBOX_HEIGHT / 2;
const TOP_BOUND = 92;
const BOTTOM_BOUND = VIEWBOX_HEIGHT - 92;
const PEER_WIDTH = 152;
const PEER_HEIGHT = 46;
const VICTIM_WIDTH = 170;
const VICTIM_HEIGHT = 54;
const MIN_ZOOM = 0.78;
const MAX_ZOOM = 1.9;
const LABEL_HEIGHT = 20;
const PLAYBACK_INTERVAL_MS = 900;
const ZOOM_STEP = 0.12;
const DEFAULT_SLICE_MS = 15000;
const MIN_SLICE_MS = 1;
const INPUT_STEP_MS = 100;

const formatBytes = (value) => {
  const bytes = Number(value || 0);
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(1)} GB`;
};

const formatDurationMs = (value) => {
  const durationMs = Number(value || 0);
  if (durationMs < 1000) return `${durationMs} ms`;
  if (durationMs < 60000) return `${(durationMs / 1000).toFixed(durationMs % 1000 === 0 ? 0 : 2)} s`;
  const minutes = Math.floor(durationMs / 60000);
  const seconds = (durationMs % 60000) / 1000;
  if (seconds === 0) return `${minutes}m`;
  return `${minutes}m ${seconds.toFixed(seconds % 1 === 0 ? 0 : 2)}s`;
};

const clamp = (value, min, max) => Math.min(Math.max(value, min), max);

const ellipsis = (value, length = 18) => {
  const text = value || '';
  return text.length > length ? `${text.slice(0, length)}...` : text;
};

const edgeTone = (direction) => {
  if (direction === 'ingress') return { hard: '#597ef7', soft: 'rgba(89,126,247,0.24)' };
  if (direction === 'egress') return { hard: '#a3a3a3', soft: 'rgba(163,163,163,0.18)' };
  return { hard: '#737373', soft: 'rgba(115,115,115,0.18)' };
};

const edgeChipClass = (direction) => {
  if (direction === 'ingress') return 'border-geek-400/30 bg-geek-400/10 text-geek-400';
  if (direction === 'egress') return 'border-neutral-400/30 bg-neutral-400/10 text-neutral-300';
  return 'border-neutral-500/30 bg-neutral-500/10 text-neutral-300';
};

const rectsOverlap = (a, b, padding = 0) =>
  a.x < b.x + b.w + padding && a.x + a.w + padding > b.x && a.y < b.y + b.h + padding && a.y + a.h + padding > b.y;

const sortByTraffic = (items) =>
  [...items].sort((a, b) => {
    if ((b.bytes || 0) !== (a.bytes || 0)) return (b.bytes || 0) - (a.bytes || 0);
    return String(a.ip || a.id || '').localeCompare(String(b.ip || b.id || ''));
  });

const distributeIntoLanes = (items, laneCount) => {
  const lanes = Array.from({ length: laneCount }, () => []);
  items.forEach((item, index) => {
    lanes[index % laneCount].push(item);
  });
  return lanes;
};

const getViewportMetrics = (zoom) => ({
  width: VIEWBOX_WIDTH / zoom,
  height: VIEWBOX_HEIGHT / zoom,
});

const clampPanToViewport = (pan, zoom) => {
  const { width, height } = getViewportMetrics(zoom);
  if (zoom <= 1) {
    return {
      x: (VIEWBOX_WIDTH - width) / 2,
      y: (VIEWBOX_HEIGHT - height) / 2,
    };
  }
  const minX = Math.min(0, VIEWBOX_WIDTH - width);
  const maxX = Math.max(0, VIEWBOX_WIDTH - width);
  const minY = Math.min(0, VIEWBOX_HEIGHT - height);
  const maxY = Math.max(0, VIEWBOX_HEIGHT - height);
  return {
    x: minX === maxX ? minX : clamp(pan.x, minX, maxX),
    y: minY === maxY ? minY : clamp(pan.y, minY, maxY),
  };
};

const buildCompactNodeLabel = (node) => {
  if (node.kind === 'victim') {
    return ellipsis(node.label || node.ip || node.id, 16);
  }
  return ellipsis(node.ip || node.label || node.id, 16);
};

const buildCompactNodeMeta = (node) => `${formatBytes(node.bytes)} / ${node.connections || 0} conn`;

const sanitizeSlice = (value) => {
  const nextValue = Math.round(Number(value) || 0);
  return Math.max(MIN_SLICE_MS, nextValue);
};

const createDemoTopology = (shift = 0, duration = DEFAULT_SLICE_MS, id = 'demo') => {
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
    window: { start: shift, end: Math.min(totalDuration, shift + duration), duration, total: totalDuration, total_count: 115 },
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

function buildPositions(nodes) {
  const positions = new Map();
  const leftNodes = [];
  const rightNodes = [];
  const centerNodes = [];
  const orbitNodes = [];

  sortByTraffic(nodes).forEach((node) => {
    if (node.side === 'left') leftNodes.push(node);
    else if (node.side === 'right') rightNodes.push(node);
    else if (node.side === 'center') centerNodes.push(node);
    else orbitNodes.push(node);
  });

  orbitNodes.forEach((node, index) => {
    if (leftNodes.length + index <= rightNodes.length + Math.floor(index / 2)) {
      leftNodes.push({ ...node, side: 'left' });
    } else {
      rightNodes.push({ ...node, side: 'right' });
    }
  });

  const placeSideNodes = (items, singleX, multiX) => {
    if (!items.length) return;
    const laneCount = items.length > 4 ? 2 : 1;
    const lanes = distributeIntoLanes(items, laneCount);
    lanes.forEach((lane, laneIndex) => {
      if (!lane.length) return;
      const x = laneCount === 1 ? singleX : multiX[laneIndex];
      const gap = lane.length === 1 ? 0 : clamp((BOTTOM_BOUND - TOP_BOUND) / Math.max(lane.length - 1, 1), 62, 108);
      const totalHeight = gap * Math.max(0, lane.length - 1);
      const startY = CENTER_Y - totalHeight / 2;
      lane.forEach((node, index) => {
        positions.set(node.id, {
          x,
          y: startY + index * gap,
          w: laneCount === 1 ? PEER_WIDTH : PEER_WIDTH - 8,
          h: PEER_HEIGHT,
        });
      });
    });
  };

  placeSideNodes(sortByTraffic(leftNodes), 176, [110, 240]);
  placeSideNodes(sortByTraffic(rightNodes), VIEWBOX_WIDTH - 176, [VIEWBOX_WIDTH - 240, VIEWBOX_WIDTH - 110]);

  const orderedCenter = sortByTraffic(centerNodes);
  if (orderedCenter.length === 1) {
    positions.set(orderedCenter[0].id, { x: CENTER_X, y: CENTER_Y, w: VICTIM_WIDTH, h: VICTIM_HEIGHT });
    return positions;
  }

  if (orderedCenter.length === 2) {
    orderedCenter.forEach((node, index) => {
      positions.set(node.id, {
        x: CENTER_X + (index === 0 ? -88 : 88),
        y: CENTER_Y,
        w: VICTIM_WIDTH - 2,
        h: VICTIM_HEIGHT,
      });
    });
    return positions;
  }

  const columnCount = orderedCenter.length <= 4 ? 2 : 3;
  const rowCount = Math.ceil(orderedCenter.length / columnCount);
  const xGap = columnCount === 2 ? 124 : 106;
  const yGap = rowCount === 1 ? 0 : clamp(138 / Math.max(rowCount - 1, 1), 72, 92);
  const startX = CENTER_X - ((columnCount - 1) * xGap) / 2;
  const startY = CENTER_Y - ((rowCount - 1) * yGap) / 2;
  const width = orderedCenter.length > 4 ? VICTIM_WIDTH - 18 : VICTIM_WIDTH - 8;
  const height = orderedCenter.length > 4 ? VICTIM_HEIGHT - 2 : VICTIM_HEIGHT;

  orderedCenter.forEach((node, index) => {
    const row = Math.floor(index / columnCount);
    const column = index % columnCount;
    positions.set(node.id, {
      x: startX + column * xGap + (columnCount === 3 && row % 2 === 1 ? 8 : 0),
      y: startY + row * yGap,
      w: width,
      h: height,
    });
  });

  return positions;
}

function buildPortOffsets(edges, positions) {
  const groups = new Map();

  edges.forEach((edge) => {
    const source = positions.get(edge.source);
    const target = positions.get(edge.target);
    if (!source || !target) return;

    const sourceList = groups.get(edge.source) || [];
    sourceList.push({ key: `${edge.id}:source`, otherX: target.x, otherY: target.y });
    groups.set(edge.source, sourceList);

    const targetList = groups.get(edge.target) || [];
    targetList.push({ key: `${edge.id}:target`, otherX: source.x, otherY: source.y });
    groups.set(edge.target, targetList);
  });

  const offsets = new Map();
  groups.forEach((items) => {
    items.sort((a, b) => a.otherY - b.otherY || a.otherX - b.otherX);
    if (items.length === 1) {
      offsets.set(items[0].key, 0);
      return;
    }
    const spread = items.length > 4 ? 24 : 32;
    items.forEach((item, index) => {
      const ratio = index / (items.length - 1);
      offsets.set(item.key, (ratio - 0.5) * spread);
    });
  });

  return offsets;
}

function buildLines(edges, positions) {
  const offsets = buildPortOffsets(edges, positions);
  const bandCounts = { ingress: 0, egress: 0, internal: 0, external: 0 };

  return [...edges]
    .sort((a, b) => {
      if ((a.bytes || 0) !== (b.bytes || 0)) return (a.bytes || 0) - (b.bytes || 0);
      return String(a.id || '').localeCompare(String(b.id || ''));
    })
    .map((edge) => {
      const source = positions.get(edge.source);
      const target = positions.get(edge.target);
      if (!source || !target) return null;

      const directionX = Math.sign(target.x - source.x) || (edge.direction === 'egress' ? 1 : -1);
      const fromX = source.x + (directionX * source.w) / 2;
      const toX = target.x - (directionX * target.w) / 2;
      const fromY = source.y + (offsets.get(`${edge.id}:source`) || 0);
      const toY = target.y + (offsets.get(`${edge.id}:target`) || 0);
      const dx = toX - fromX;
      const dy = toY - fromY;
      const band = bandCounts[edge.direction] || 0;
      bandCounts[edge.direction] = band + 1;

      let c1y;
      let c2y;
      let labelBaseY;

      if (edge.direction === 'internal') {
        const apexY = Math.min(fromY, toY) - 34 - (band % 3) * 18;
        c1y = apexY;
        c2y = apexY;
        labelBaseY = apexY - 14;
      } else {
        const sign = edge.direction === 'ingress' ? -1 : 1;
        const arc = 30 + Math.abs(dx) * 0.045 + (band % 4) * 12;
        c1y = fromY + sign * arc - dy * 0.08;
        c2y = toY + sign * arc + dy * 0.08;
        labelBaseY = (fromY + toY) / 2 + sign * (arc * 0.82);
      }

      return {
        ...edge,
        path: `M ${fromX} ${fromY} C ${fromX + dx * 0.35} ${c1y}, ${toX - dx * 0.35} ${c2y}, ${toX} ${toY}`,
        labelBaseX: fromX + dx * 0.5,
        labelBaseY,
        tone: edgeTone(edge.direction),
      };
    })
    .filter(Boolean);
}

function resolveEdgeLabels(lines, positions, selectedEdgeId) {
  const nodeBoxes = Array.from(positions.values()).map((pos) => ({
    x: pos.x - pos.w / 2 - 6,
    y: pos.y - pos.h / 2 - 6,
    w: pos.w + 12,
    h: pos.h + 12,
  }));
  const labelRects = [];
  const placements = new Map();
  const labelCandidates = [...lines]
    .sort((a, b) => {
      const aPriority = a.id === selectedEdgeId ? 1 : 0;
      const bPriority = b.id === selectedEdgeId ? 1 : 0;
      if (aPriority !== bPriority) return bPriority - aPriority;
      if ((b.bytes || 0) !== (a.bytes || 0)) return (b.bytes || 0) - (a.bytes || 0);
      return String(a.id || '').localeCompare(String(b.id || ''));
    })
    .filter((line, index) => line.id === selectedEdgeId || index < 6);

  labelCandidates.forEach((line) => {
    const text = formatBytes(line.bytes);
    const width = clamp(text.length * 6.7 + 20, 56, 86);
    const candidates = [
      { x: line.labelBaseX - width / 2, y: line.labelBaseY - LABEL_HEIGHT / 2 },
      { x: line.labelBaseX - width / 2, y: line.labelBaseY - 24 },
      { x: line.labelBaseX - width / 2, y: line.labelBaseY + 5 },
      { x: line.labelBaseX + 10, y: line.labelBaseY - LABEL_HEIGHT / 2 },
      { x: line.labelBaseX - width - 10, y: line.labelBaseY - LABEL_HEIGHT / 2 },
    ];

    let placement = null;
    for (const candidate of candidates) {
      const rect = {
        x: clamp(candidate.x, 12, VIEWBOX_WIDTH - width - 12),
        y: clamp(candidate.y, 12, VIEWBOX_HEIGHT - LABEL_HEIGHT - 12),
        w: width,
        h: LABEL_HEIGHT,
      };
      const nodeCollision = nodeBoxes.some((box) => rectsOverlap(rect, box, line.id === selectedEdgeId ? 2 : 4));
      const labelCollision = labelRects.some((box) => rectsOverlap(rect, box, 4));
      if (!nodeCollision && !labelCollision) {
        placement = rect;
        break;
      }
    }

    if (!placement && line.id === selectedEdgeId) {
      placement = {
        x: clamp(line.labelBaseX - width / 2, 12, VIEWBOX_WIDTH - width - 12),
        y: clamp(line.labelBaseY - LABEL_HEIGHT / 2, 12, VIEWBOX_HEIGHT - LABEL_HEIGHT - 12),
        w: width,
        h: LABEL_HEIGHT,
      };
    }

    if (placement) {
      placements.set(line.id, { ...placement, text });
      labelRects.push({
        x: placement.x - 4,
        y: placement.y - 3,
        w: placement.w + 8,
        h: placement.h + 6,
      });
    }
  });

  return placements;
}

function TrafficGraphModal({ isOpen, onClose, container, contestId, teamId, fetchTraffic: customFetchTraffic }) {
  const { t, i18n } = useTranslation();
  const canvasRef = useRef(null);
  const dragRef = useRef({ active: false, moved: false, x: 0, y: 0, panX: 0, panY: 0 });
  const requestSequenceRef = useRef(0);
  const [topology, setTopology] = useState(null);
  const [shift, setShift] = useState(0);
  const [slice, setSlice] = useState(DEFAULT_SLICE_MS);
  const [sliceInput, setSliceInput] = useState(String(DEFAULT_SLICE_MS));
  const [demoMode, setDemoMode] = useState(false);
  const [isFetching, setIsFetching] = useState(false);
  const [isPlaying, setIsPlaying] = useState(false);
  const [selectedEdgeId, setSelectedEdgeId] = useState('');
  const [selectedNodeId, setSelectedNodeId] = useState('');
  const [zoom, setZoom] = useState(1);
  const [pan, setPan] = useState({ x: 0, y: 0 });

  const fetchData = async ({ nextShift = shift, nextSlice = slice, forceLive = false } = {}) => {
    if (!container?.id) return;
    const resolvedSlice = sanitizeSlice(nextSlice);
    const requestId = requestSequenceRef.current + 1;
    requestSequenceRef.current = requestId;
    setIsFetching(true);
    try {
      const response = customFetchTraffic
        ? await customFetchTraffic(container, { time_shift: nextShift, duration: resolvedSlice })
        : await getContestTeamTraffic(contestId, teamId, container.id, { time_shift: nextShift, duration: resolvedSlice });
      if (requestSequenceRef.current !== requestId) return;
      if (response.code !== 200) throw new Error(t('admin.contests.trafficGraph.toast.fetchFailed'));
      setTopology(response.data);
      setDemoMode(false);
    } catch {
      if (requestSequenceRef.current !== requestId) return;
      if (!forceLive) {
        toast.warning({ description: t('admin.contests.trafficGraph.toast.demoFallback') });
      }
      setTopology(createDemoTopology(nextShift, resolvedSlice, container.id));
      setDemoMode(true);
    } finally {
      if (requestSequenceRef.current === requestId) {
        setIsFetching(false);
        setSelectedEdgeId('');
        setSelectedNodeId('');
      }
    }
  };

  useEffect(() => {
    if (isOpen) return;
    requestSequenceRef.current += 1;
    dragRef.current.active = false;
    setDemoMode(false);
    setIsFetching(false);
    setIsPlaying(false);
  }, [isOpen]);

  useEffect(() => {
    if (!isOpen) return;
    setIsPlaying(false);
    setSelectedEdgeId('');
    setSelectedNodeId('');
    setZoom(1);
    setPan({ x: 0, y: 0 });
    setSlice(DEFAULT_SLICE_MS);
    setSliceInput(String(DEFAULT_SLICE_MS));
  }, [isOpen, container?.id]);

  useEffect(() => {
    if (!isOpen || !container?.id) return;
    fetchData({ nextShift: shift, nextSlice: slice, forceLive: demoMode });
  }, [isOpen, container?.id, shift, slice]);

  const nodes = topology?.nodes || [];
  const edges = topology?.edges || [];
  const summary = topology?.summary || {};
  const windowInfo = topology?.window || { start: 0, end: 0, duration: slice, total: 0 };
  const timeline = topology?.timeline || [];
  const peakTimeline = Math.max(1, ...timeline.map((bucket) => bucket.bytes || 0));
  const totalDuration = Math.max(windowInfo.total || 0, topology?.total_duration || 0);
  const maxShift = Math.max(0, totalDuration - slice);
  const playbackFrames = Math.max(1, Math.floor(maxShift / Math.max(slice, 1)) + 1);
  const playbackIndex = Math.min(playbackFrames, Math.floor(Math.max(shift, 0) / Math.max(slice, 1)) + 1);
  const canStepBackward = shift > 0;
  const canStepForward = shift < maxShift;
  const viewportMetrics = useMemo(() => getViewportMetrics(zoom), [zoom]);
  const viewBox = `${pan.x} ${pan.y} ${viewportMetrics.width} ${viewportMetrics.height}`;
  const positions = useMemo(() => buildPositions(nodes), [nodes]);
  const lines = useMemo(() => buildLines(edges, positions), [edges, positions]);
  const labelPlacements = useMemo(
    () => resolveEdgeLabels(lines, positions, selectedEdgeId),
    [lines, positions, selectedEdgeId]
  );
  const orderedLines = useMemo(
    () =>
      [...lines].sort((a, b) => {
        const aSelected = a.id === selectedEdgeId ? 1 : 0;
        const bSelected = b.id === selectedEdgeId ? 1 : 0;
        if (aSelected !== bSelected) return aSelected - bSelected;
        if ((a.bytes || 0) !== (b.bytes || 0)) return (a.bytes || 0) - (b.bytes || 0);
        return String(a.id || '').localeCompare(String(b.id || ''));
      }),
    [lines, selectedEdgeId]
  );
  const selectedEdge = useMemo(
    () => edges.find((edge) => edge.id === selectedEdgeId) || edges[0] || null,
    [edges, selectedEdgeId]
  );
  const selectedNode = useMemo(
    () =>
      nodes.find((node) => node.id === selectedNodeId) ||
      nodes.find((node) => node.kind === 'victim') ||
      nodes[0] ||
      null,
    [nodes, selectedNodeId]
  );

  const visibleIPs = (topology?.center?.ips || []).slice(0, 3);
  const hiddenIpCount = Math.max(0, (topology?.center?.ips || []).length - visibleIPs.length);
  const visibleExposed = (topology?.center?.exposed || []).slice(0, 1);
  const topTalkers = (topology?.top_talkers || []).slice(0, 2);
  const topEdges = (topology?.top_edges || []).slice(0, 2);

  useEffect(() => {
    setShift((current) => Math.min(current, maxShift));
  }, [maxShift]);

  useEffect(() => {
    setSliceInput(String(slice));
  }, [slice]);

  useEffect(() => {
    if (!isPlaying || !isOpen || !container?.id) return;
    if (maxShift <= 0 || shift >= maxShift) {
      setIsPlaying(false);
      return;
    }
    if (isFetching) return;
    const timer = window.setTimeout(() => {
      setShift((current) => Math.min(current + slice, maxShift));
    }, PLAYBACK_INTERVAL_MS);
    return () => window.clearTimeout(timer);
  }, [container?.id, isFetching, isOpen, isPlaying, maxShift, shift, slice]);

  const applyZoom = (targetZoom, anchor) => {
    const rect = canvasRef.current?.getBoundingClientRect();
    const nextZoom = clamp(targetZoom, MIN_ZOOM, MAX_ZOOM);
    if (!rect?.width || !rect?.height) {
      setZoom(nextZoom);
      setPan((current) => clampPanToViewport(current, nextZoom));
      return;
    }
    const anchorX = anchor?.x ?? rect.width / 2;
    const anchorY = anchor?.y ?? rect.height / 2;
    const currentViewport = getViewportMetrics(zoom);
    const nextViewport = getViewportMetrics(nextZoom);
    const worldX = pan.x + (anchorX / rect.width) * currentViewport.width;
    const worldY = pan.y + (anchorY / rect.height) * currentViewport.height;
    const nextPan = clampPanToViewport(
      {
        x: worldX - (anchorX / rect.width) * nextViewport.width,
        y: worldY - (anchorY / rect.height) * nextViewport.height,
      },
      nextZoom
    );
    setZoom(nextZoom);
    setPan(nextPan);
  };

  const handleWheel = (event) => {
    event.preventDefault();
    const rect = canvasRef.current?.getBoundingClientRect();
    if (!rect) return;
    const nextZoom = clamp(zoom + (event.deltaY < 0 ? ZOOM_STEP : -ZOOM_STEP), MIN_ZOOM, MAX_ZOOM);
    applyZoom(nextZoom, { x: event.clientX - rect.left, y: event.clientY - rect.top });
  };

  const startDrag = (event) => {
    if (event.button !== 0) return;
    dragRef.current = {
      active: true,
      moved: false,
      x: event.clientX,
      y: event.clientY,
      panX: pan.x,
      panY: pan.y,
    };
  };

  const onDrag = (event) => {
    if (!dragRef.current.active) return;
    const rect = canvasRef.current?.getBoundingClientRect();
    if (!rect?.width || !rect?.height) return;
    const dx = event.clientX - dragRef.current.x;
    const dy = event.clientY - dragRef.current.y;
    if (Math.abs(dx) > 3 || Math.abs(dy) > 3) {
      dragRef.current.moved = true;
    }
    const currentViewport = getViewportMetrics(zoom);
    setPan(
      clampPanToViewport(
        {
          x: dragRef.current.panX - dx * (currentViewport.width / rect.width),
          y: dragRef.current.panY - dy * (currentViewport.height / rect.height),
        },
        zoom
      )
    );
  };

  const stopDrag = () => {
    dragRef.current.active = false;
  };

  const resetView = () => {
    setZoom(1);
    setPan({ x: 0, y: 0 });
  };

  const handleSelectEdge = (edgeId) => {
    if (dragRef.current.moved) return;
    setSelectedEdgeId(edgeId);
  };

  const handleSelectNode = (nodeId) => {
    if (dragRef.current.moved) return;
    setSelectedNodeId(nodeId);
  };

  const stepPlayback = (direction) => {
    setIsPlaying(false);
    setShift((current) => clamp(current + direction * slice, 0, maxShift));
  };

  const handleSliceInputChange = (event) => {
    setSliceInput(event.target.value);
  };

  const commitSliceInput = () => {
    const nextSlice = sanitizeSlice(sliceInput);
    setSliceInput(String(nextSlice));
    setSlice(nextSlice);
    setIsPlaying(false);
  };

  const togglePlayback = () => {
    if (isPlaying) {
      setIsPlaying(false);
      return;
    }
    if (shift >= maxShift && maxShift > 0) {
      setShift(0);
    }
    setIsPlaying(true);
  };

  const downloadTraffic = async () => {
    if (!container?.id) return;
    try {
      const response =
        contestId && teamId
          ? await downloadContainerTraffic(contestId, teamId, container.id)
          : await downloadVictimTraffic(container.id);
      if (response.headers?.['file'] === 'true') {
        downloadBlobResponse(response, `traffic_${container.id}.zip`);
      }
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
          .traffic-line { stroke-dasharray: 11 10; animation: traffic-flow 10s linear infinite; }
        `}
      </style>
      <Modal
        isOpen={isOpen}
        onClose={onClose}
        title={t('admin.contests.trafficGraph.title')}
        size="2xl"
        className="!bg-neutral-900/95 !border-neutral-600"
        bodyClassName="p-4 h-[74vh] max-h-[820px] overflow-hidden"
      >
        <div className="flex h-full min-h-0 flex-col gap-3 text-neutral-100">
          <div className="grid grid-cols-1 gap-3 xl:grid-cols-[minmax(0,1fr)_360px]">
            <div className="rounded-2xl border border-neutral-600 bg-black/20 px-4 py-3">
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div className="min-w-0">
                  <div className="text-[11px] uppercase tracking-[0.24em] text-geek-400/80">
                    {t('admin.contests.trafficGraph.hero.kicker')}
                  </div>
                  <div className="mt-1 text-lg font-['Maple_UI'] text-white">
                    {topology?.center?.label || `Victim #${container?.id || '-'}`}
                  </div>
                  <div className="mt-1 max-w-[64ch] text-xs leading-5 text-neutral-400">
                    {t('admin.contests.trafficGraph.hero.subtitle', {
                      challenge: container?.challenge || container?.contest_challenge_name || `#${container?.id}`,
                    })}
                  </div>
                </div>
                <div className="flex flex-wrap gap-2">
                  {visibleIPs.map((ip) => (
                    <Chip
                      key={ip}
                      label={ip}
                      variant="tag"
                      size="sm"
                      colorClass="border-geek-400/30 bg-geek-400/10 text-geek-400"
                    />
                  ))}
                  {hiddenIpCount > 0 ? (
                    <Chip
                      label={`+${hiddenIpCount}`}
                      variant="tag"
                      size="sm"
                      colorClass="border-neutral-500/30 bg-neutral-500/10 text-neutral-300"
                    />
                  ) : null}
                  {visibleExposed.map((item) => (
                    <Chip
                      key={item}
                      label={ellipsis(item, 18)}
                      variant="tag"
                      size="sm"
                      colorClass="border-neutral-500/30 bg-black/20 text-neutral-300"
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
            </div>

            <div className="rounded-2xl border border-neutral-600 bg-black/20 px-4 py-3">
              <div className="grid gap-2">
                <div>
                  <div className="flex items-center justify-between text-[11px] text-neutral-400">
                    <span>{t('admin.contests.trafficGraph.controls.timeShift')}</span>
                    <span className="font-mono text-geek-400">
                      {t('admin.contests.trafficGraph.hero.windowAt', {
                        start: formatDurationMs(windowInfo.start),
                        end: formatDurationMs(windowInfo.end),
                      })}
                    </span>
                  </div>
                  <input
                    type="range"
                    min="0"
                    max={maxShift}
                    value={shift}
                    onChange={(event) => {
                      setIsPlaying(false);
                      setShift(Number(event.target.value));
                    }}
                    className="mt-1 w-full accent-geek-400"
                  />
                </div>
                <div className="flex flex-wrap items-center gap-2">
                  <label className="mr-1 text-[11px] text-neutral-400" htmlFor="traffic-graph-slice-ms">
                    {t('admin.contests.trafficGraph.controls.timeSlice')}
                  </label>
                  <div className="w-[132px]">
                    <Input
                      id="traffic-graph-slice-ms"
                      type="number"
                      min={MIN_SLICE_MS}
                      step={INPUT_STEP_MS}
                      value={sliceInput}
                      onChange={handleSliceInputChange}
                      onBlur={commitSliceInput}
                      onKeyDown={(event) => {
                        if (event.key === 'Enter') {
                          event.currentTarget.blur();
                        }
                      }}
                      className="!h-8 !border-neutral-600 !bg-black/20 !px-3 text-[11px] font-mono"
                    />
                  </div>
                  <Chip
                    label={t('admin.contests.trafficGraph.controls.milliseconds')}
                    variant="tag"
                    size="sm"
                    colorClass="border-neutral-500/30 bg-black/20 text-neutral-300"
                  />
                  <div className="ml-auto flex items-center gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
                      onClick={() => fetchData({ nextShift: shift, nextSlice: slice, forceLive: true })}
                      title={t('common.refresh')}
                    >
                      <IconRefresh size={16} className={isFetching ? 'animate-spin' : ''} />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
                      onClick={downloadTraffic}
                      title={t('admin.contests.teamDetail.traffic.actions.downloadTraffic')}
                    >
                      <IconDownload size={16} />
                    </Button>
                  </div>
                </div>
                <div className="flex flex-wrap items-center gap-2 border-t border-neutral-700/80 pt-2">
                  <Chip
                    label={t('admin.contests.trafficGraph.footer.timeSlice', { count: formatDurationMs(slice) })}
                    variant="tag"
                    size="sm"
                    colorClass="border-neutral-500/30 bg-black/20 text-neutral-300"
                  />
                  {isPlaying ? (
                    <Chip
                      label={t('admin.contests.trafficGraph.footer.playing', {
                        current: playbackIndex,
                        total: playbackFrames,
                      })}
                      variant="tag"
                      size="sm"
                      colorClass="border-geek-400/30 bg-geek-400/10 text-geek-400"
                    />
                  ) : null}
                  <div className="ml-auto flex items-center gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
                      onClick={() => stepPlayback(-1)}
                      disabled={!canStepBackward || isFetching}
                      title={t('common.previous')}
                    >
                      <IconPlayerTrackPrev size={16} />
                    </Button>
                    <Button
                      variant={isPlaying ? 'outline' : 'primary'}
                      size="sm"
                      className="!h-8 !px-3"
                      onClick={togglePlayback}
                      disabled={totalDuration <= 0}
                      icon={isPlaying ? <IconPlayerPause size={14} /> : <IconPlayerPlay size={14} />}
                    >
                      {isPlaying
                        ? t('admin.contests.trafficGraph.controls.pauseReplay')
                        : t('admin.contests.trafficGraph.controls.playReplay')}
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
                      onClick={() => stepPlayback(1)}
                      disabled={!canStepForward || isFetching}
                      title={t('common.next')}
                    >
                      <IconPlayerTrackNext size={16} />
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-2 xl:grid-cols-4">
            {[
              {
                icon: <IconActivity size={15} className="text-geek-400" />,
                label: t('admin.contests.trafficGraph.stats.totalTraffic'),
                value: formatBytes(summary.total_bytes),
                tone: 'text-geek-400',
              },
              {
                icon: <IconArrowDownRight size={15} className="text-neutral-300" />,
                label: t('admin.contests.trafficGraph.stats.ingress'),
                value: formatBytes(summary.ingress_bytes),
                tone: 'text-neutral-100',
              },
              {
                icon: <IconArrowUpRight size={15} className="text-neutral-300" />,
                label: t('admin.contests.trafficGraph.stats.egress'),
                value: formatBytes(summary.egress_bytes),
                tone: 'text-neutral-100',
              },
              {
                icon: <IconRoute size={15} className="text-neutral-300" />,
                label: t('admin.contests.trafficGraph.stats.peers'),
                value: summary.visible_nodes || 0,
                tone: 'text-neutral-100',
              },
            ].map((item) => (
              <div key={item.label} className="rounded-xl border border-neutral-600 bg-black/20 px-3 py-2">
                <div className="flex items-center gap-3">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-neutral-800">{item.icon}</div>
                  <div className="min-w-0">
                    <div className="text-[11px] font-mono text-neutral-500">{item.label}</div>
                    <div className={`truncate font-mono text-base ${item.tone}`}>{item.value}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>

          <div className="grid min-h-0 flex-1 grid-cols-1 gap-3 xl:grid-cols-[minmax(0,1.55fr)_320px]">
            <Card padding="none" className="flex min-h-0 overflow-hidden rounded-2xl border-neutral-600 bg-neutral-900">
              <div className="flex min-h-0 flex-1 flex-col">
                <div className="flex flex-wrap items-start justify-between gap-2 border-b border-neutral-600 px-4 py-3">
                  <div>
                    <div className="text-sm font-mono text-neutral-300">
                      {t('admin.contests.trafficGraph.canvas.title')}
                    </div>
                    <div className="mt-1 text-[11px] text-neutral-500">
                      {t('admin.contests.trafficGraph.canvas.subtitle')}
                    </div>
                  </div>
                  <div className="flex flex-wrap items-center justify-end gap-2">
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
                    <Chip
                      label={`x${zoom.toFixed(2)}`}
                      variant="tag"
                      size="sm"
                      colorClass="border-neutral-500/30 bg-black/20 text-neutral-300"
                    />
                    <Button
                      variant="ghost"
                      size="icon"
                      className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
                      onClick={() => applyZoom(zoom - ZOOM_STEP)}
                    >
                      <IconZoomOut size={16} />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
                      onClick={() => applyZoom(zoom + ZOOM_STEP)}
                    >
                      <IconZoomIn size={16} />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="!h-8 !px-3 !text-neutral-300 hover:!text-neutral-100"
                      onClick={resetView}
                    >
                      {t('common.reset')}
                    </Button>
                  </div>
                </div>

                <div
                  ref={canvasRef}
                  className="relative min-h-0 flex-1 overflow-hidden bg-[linear-gradient(180deg,rgba(0,0,0,.14),rgba(0,0,0,.05))] cursor-grab active:cursor-grabbing"
                  onWheel={handleWheel}
                  onPointerMove={onDrag}
                  onPointerUp={stopDrag}
                  onPointerLeave={stopDrag}
                >
                  <div className="pointer-events-none absolute inset-0 opacity-20 [background-image:linear-gradient(rgba(82,82,82,.45)_1px,transparent_1px),linear-gradient(90deg,rgba(82,82,82,.45)_1px,transparent_1px)] [background-size:34px_34px]" />
                  <div className="absolute left-3 top-3 z-10 text-[11px] text-neutral-400">
                    {t('admin.contests.trafficGraph.canvas.zoomHint')}
                  </div>

                  {nodes.length === 0 ? (
                    <div className="absolute inset-0 flex items-center justify-center px-6 text-center text-sm text-neutral-500">
                      {t('admin.contests.trafficGraph.empty')}
                    </div>
                  ) : null}

                  <svg
                    viewBox={viewBox}
                    className="absolute inset-0 h-full w-full select-none touch-none"
                    shapeRendering="geometricPrecision"
                    textRendering="geometricPrecision"
                    onPointerDown={startDrag}
                  >
                    <g opacity="0.18">
                      <path
                        d={`M 34 ${CENTER_Y - 132} H ${VIEWBOX_WIDTH - 34}`}
                        stroke="#404040"
                        strokeWidth="1"
                        vectorEffect="non-scaling-stroke"
                      />
                      <path
                        d={`M 34 ${CENTER_Y} H ${VIEWBOX_WIDTH - 34}`}
                        stroke="#404040"
                        strokeWidth="1"
                        vectorEffect="non-scaling-stroke"
                      />
                      <path
                        d={`M 34 ${CENTER_Y + 132} H ${VIEWBOX_WIDTH - 34}`}
                        stroke="#404040"
                        strokeWidth="1"
                        vectorEffect="non-scaling-stroke"
                      />
                    </g>

                    {orderedLines.map((edge) => {
                      const label = labelPlacements.get(edge.id);
                      const selected = selectedEdge?.id === edge.id;
                      return (
                        <g key={edge.id}>
                          <path
                            d={edge.path}
                            fill="none"
                            stroke={edge.tone.hard}
                            strokeWidth={selected ? 8.5 : 7}
                            opacity={selected ? 0.22 : 0.08 + edge.intensity * 0.1}
                            vectorEffect="non-scaling-stroke"
                          />
                          <path
                            d={edge.path}
                            fill="none"
                            stroke={edge.tone.soft}
                            strokeWidth={2.2 + edge.intensity * 3.2}
                            className="traffic-line"
                            opacity={selected ? 0.95 : 0.78}
                            onClick={() => handleSelectEdge(edge.id)}
                            vectorEffect="non-scaling-stroke"
                          />
                          {label ? (
                            <g
                              transform={`translate(${label.x + label.w / 2}, ${label.y + label.h / 2})`}
                              onClick={() => handleSelectEdge(edge.id)}
                            >
                              <rect
                                x={-label.w / 2}
                                y={-label.h / 2}
                                width={label.w}
                                height={label.h}
                                rx={10}
                                fill="rgba(10,10,10,.95)"
                                stroke={selected ? edge.tone.hard : 'rgba(115,115,115,.35)'}
                                vectorEffect="non-scaling-stroke"
                              />
                              <text
                                x="0"
                                y="4"
                                textAnchor="middle"
                                fill="#d4d4d4"
                                fontSize="10.2"
                                fontFamily="Maple Mono"
                              >
                                {label.text}
                              </text>
                            </g>
                          ) : null}
                        </g>
                      );
                    })}

                    {nodes.map((node) => {
                      const pos = positions.get(node.id);
                      if (!pos) return null;
                      const active = selectedNode?.id === node.id;
                      const isVictim = node.kind === 'victim';
                      const x = pos.x - pos.w / 2;
                      const y = pos.y - pos.h / 2;
                      const title = buildCompactNodeLabel(node);
                      const meta = buildCompactNodeMeta(node);
                      return (
                        <g key={node.id} transform={`translate(${x}, ${y})`} onClick={() => handleSelectNode(node.id)}>
                          <rect
                            x={active ? -2.5 : 0}
                            y={active ? -2.5 : 0}
                            width={active ? pos.w + 5 : pos.w}
                            height={active ? pos.h + 5 : pos.h}
                            rx="15"
                            fill={isVictim ? 'rgba(89,126,247,.12)' : 'rgba(24,24,27,.94)'}
                            stroke={active ? '#f5f5f5' : isVictim ? '#597ef7' : '#666666'}
                            strokeWidth={active ? 1.7 : 1.1}
                            vectorEffect="non-scaling-stroke"
                          />
                          <circle
                            cx="18"
                            cy={pos.h / 2}
                            r={isVictim ? 9.5 : 8}
                            fill={isVictim ? 'rgba(89,126,247,.14)' : 'rgba(82,82,91,.62)'}
                            stroke={isVictim ? '#597ef7' : '#8a8a8a'}
                            vectorEffect="non-scaling-stroke"
                          />
                          <text x="34" y={pos.h / 2 - 4} fill="#f5f5f5" fontSize="10.7" fontFamily="Maple Mono">
                            {title}
                          </text>
                          <text x="34" y={pos.h / 2 + 11} fill="#a3a3a3" fontSize="9.4" fontFamily="Maple Mono">
                            {ellipsis(meta, 20)}
                          </text>
                        </g>
                      );
                    })}
                  </svg>
                </div>
              </div>
            </Card>

            <div className="grid min-h-0 grid-rows-[auto_auto_minmax(0,1fr)] gap-3">
              <Card padding="sm" className="rounded-2xl border-neutral-600 bg-neutral-900">
                <div className="grid gap-2">
                  <div className="rounded-xl border border-neutral-600 bg-black/20 p-3">
                    <div className="flex items-start justify-between gap-2">
                      <div className="min-w-0">
                        <div className="text-xs font-mono text-neutral-300">
                          {t('admin.contests.trafficGraph.panel.selectedFlow')}
                        </div>
                        <div className="mt-1 truncate text-[11px] text-neutral-500">
                          {selectedEdge
                            ? `${selectedEdge.source} -> ${selectedEdge.target}`
                            : t('admin.contests.trafficGraph.panel.selectedFlowHint')}
                        </div>
                      </div>
                      {selectedEdge ? (
                        <Chip
                          label={t(`admin.contests.trafficGraph.direction.${selectedEdge.direction || 'internal'}`)}
                          variant="tag"
                          size="sm"
                          colorClass={edgeChipClass(selectedEdge.direction)}
                        />
                      ) : null}
                    </div>
                    {selectedEdge ? (
                      <div className="mt-3 grid grid-cols-3 gap-2 text-[11px] font-mono">
                        <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                          <div className="text-neutral-500">{t('admin.contests.trafficGraph.panel.edgeBytes')}</div>
                          <div className="mt-1 text-geek-400">{formatBytes(selectedEdge.bytes)}</div>
                        </div>
                        <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                          <div className="text-neutral-500">{t('admin.contests.trafficGraph.panel.edgePackets')}</div>
                          <div className="mt-1 text-neutral-100">{selectedEdge.packets || 0}</div>
                        </div>
                        <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                          <div className="text-neutral-500">Proto</div>
                          <div className="mt-1 truncate text-neutral-100">
                            {selectedEdge.dominant_proto || selectedEdge.dominant_app || '--'}
                          </div>
                        </div>
                      </div>
                    ) : null}
                  </div>

                  <div className="rounded-xl border border-neutral-600 bg-black/20 p-3">
                    <div className="text-xs font-mono text-neutral-300">
                      {t('admin.contests.trafficGraph.panel.selectedNode')}
                    </div>
                    {selectedNode ? (
                      <>
                        <div className="mt-1 truncate text-[11px] text-white">
                          {selectedNode.label || selectedNode.ip}
                        </div>
                        <div className="mt-1 truncate text-[11px] text-neutral-500">{selectedNode.ip}</div>
                        <div className="mt-3 grid grid-cols-3 gap-2 text-[11px] font-mono">
                          <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                            <div className="text-neutral-500">{t('admin.contests.trafficGraph.panel.nodeTraffic')}</div>
                            <div className="mt-1 text-geek-400">{formatBytes(selectedNode.bytes)}</div>
                          </div>
                          <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                            <div className="text-neutral-500">
                              {t('admin.contests.trafficGraph.panel.nodeConnections')}
                            </div>
                            <div className="mt-1 text-neutral-100">{selectedNode.connections || 0}</div>
                          </div>
                          <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                            <div className="text-neutral-500">Proto</div>
                            <div className="mt-1 truncate text-neutral-100">
                              {(selectedNode.protocols || []).slice(0, 2).join(' / ') || '--'}
                            </div>
                          </div>
                        </div>
                      </>
                    ) : null}
                  </div>
                </div>
              </Card>

              <Card padding="sm" className="rounded-2xl border-neutral-600 bg-neutral-900">
                <div className="flex items-center justify-between gap-2">
                  <div className="text-sm font-mono text-neutral-300">
                    {t('admin.contests.trafficGraph.timeline.title')}
                  </div>
                  <div className="text-[11px] font-mono text-neutral-500">
                    {t('admin.contests.trafficGraph.hero.windowAt', {
                      start: formatDurationMs(windowInfo.start),
                      end: formatDurationMs(windowInfo.end),
                    })}
                  </div>
                </div>
                {timeline.length > 0 ? (
                  <>
                    <div className="mt-3 flex h-16 items-end gap-px">
                      {timeline.map((bucket) => {
                        const active =
                          bucket.timestamp_ms >= windowInfo.start &&
                          bucket.timestamp_ms < Math.max(windowInfo.end, windowInfo.start + 1);
                        const ratio = clamp((bucket.bytes || 0) / peakTimeline, 0.08, 1);
                        return (
                          <button
                            key={bucket.timestamp_ms}
                            type="button"
                            title={`${t('admin.contests.trafficGraph.timeline.at', {
                              value: formatDurationMs(bucket.timestamp_ms),
                            })} / ${formatBytes(bucket.bytes)}`}
                            onClick={() =>
                              setShift(Math.min(bucket.timestamp_ms, Math.max(0, windowInfo.total - windowInfo.duration)))
                            }
                            className="group relative flex min-w-0 flex-1 items-end"
                          >
                            <motion.span
                              className={`block w-full rounded-t-[3px] ${active ? 'bg-geek-400' : 'bg-neutral-600 group-hover:bg-neutral-500'}`}
                              initial={{ height: '0%' }}
                              animate={{ height: `${ratio * 100}%` }}
                              transition={{ duration: 0.28, ease: 'easeOut' }}
                            />
                          </button>
                        );
                      })}
                    </div>
                    <div className="mt-2 flex items-center justify-between text-[11px] font-mono text-neutral-500">
                      <span>
                        {timeline[0]
                          ? t('admin.contests.trafficGraph.timeline.at', {
                              value: formatDurationMs(timeline[0].timestamp_ms),
                            })
                          : '--'}
                      </span>
                      <span>{formatBytes(peakTimeline)}</span>
                      <span>
                        {timeline[timeline.length - 1]
                          ? t('admin.contests.trafficGraph.timeline.at', {
                              value: formatDurationMs(timeline[timeline.length - 1].timestamp_ms),
                            })
                          : '--'}
                      </span>
                    </div>
                  </>
                ) : (
                  <div className="mt-3 text-sm text-neutral-500">{t('admin.contests.trafficGraph.empty')}</div>
                )}
              </Card>

              <Card padding="sm" className="rounded-2xl border-neutral-600 bg-neutral-900">
                <div className="grid gap-3">
                  <div>
                    <div className="text-xs font-mono text-neutral-300">
                      {t('admin.contests.trafficGraph.rankings.topTalkers')}
                    </div>
                    <div className="mt-2 grid gap-2">
                      {topTalkers.length > 0 ? (
                        topTalkers.map((item, index) => (
                          <div
                            key={`${item.ip}-${index}`}
                            className="flex items-center justify-between rounded-lg border border-neutral-600 bg-black/20 px-2.5 py-2"
                          >
                            <div className="min-w-0">
                              <div className="truncate font-mono text-[11px] text-white">{item.label || item.ip}</div>
                              <div className="mt-1 text-[11px] text-neutral-500">{item.ip}</div>
                            </div>
                            <div className="ml-3 text-right font-mono">
                              <div className="text-[11px] text-geek-400">{formatBytes(item.bytes)}</div>
                              <div className="mt-1 text-[11px] text-neutral-400">{item.connections || 0} conn</div>
                            </div>
                          </div>
                        ))
                      ) : (
                        <div className="rounded-lg border border-neutral-600 bg-black/20 px-2.5 py-2 text-[11px] text-neutral-500">
                          {t('admin.contests.trafficGraph.empty')}
                        </div>
                      )}
                    </div>
                  </div>

                  <div className="border-t border-neutral-700 pt-3">
                    <div className="text-xs font-mono text-neutral-300">
                      {t('admin.contests.trafficGraph.rankings.topEdges')}
                    </div>
                    <div className="mt-2 grid gap-2">
                      {topEdges.length > 0 ? (
                        topEdges.map((item, index) => (
                          <div
                            key={`${item.id || item.label}-${index}`}
                            className="flex items-center justify-between rounded-lg border border-neutral-600 bg-black/20 px-2.5 py-2"
                          >
                            <div className="min-w-0">
                              <div className="truncate font-mono text-[11px] text-white">{item.label}</div>
                              <div className="mt-1 text-[11px] text-neutral-500">
                                {t(`admin.contests.trafficGraph.direction.${item.direction || 'internal'}`)}
                              </div>
                            </div>
                            <div className="ml-3 text-right font-mono">
                              <div className="text-[11px] text-geek-400">{formatBytes(item.bytes)}</div>
                              <div className="mt-1 text-[11px] text-neutral-400">{item.connections || 0} conn</div>
                            </div>
                          </div>
                        ))
                      ) : (
                        <div className="rounded-lg border border-neutral-600 bg-black/20 px-2.5 py-2 text-[11px] text-neutral-500">
                          {t('admin.contests.trafficGraph.empty')}
                        </div>
                      )}
                    </div>
                  </div>
                </div>

                <div className="mt-3 border-t border-neutral-700 pt-3">
                  <div className="flex flex-wrap gap-2 text-[11px] font-mono text-neutral-500">
                    <span>
                      {t('admin.contests.trafficGraph.footer.window', {
                        start: formatDurationMs(windowInfo.start),
                        end: formatDurationMs(windowInfo.end),
                      })}
                    </span>
                    <span>{t('admin.contests.trafficGraph.footer.timeSlice', { count: formatDurationMs(slice) })}</span>
                    <span>
                      {t('admin.contests.trafficGraph.footer.connectionCount', { count: summary.visible_edges || 0 })}
                    </span>
                    <span>
                      {t('admin.contests.trafficGraph.footer.ipCount', { count: summary.visible_nodes || 0 })}
                    </span>
                    <span>
                      {t('admin.contests.trafficGraph.footer.maxDuration', { count: topology?.total_duration || 0 })}
                    </span>
                    {isPlaying ? (
                      <span>
                        {t('admin.contests.trafficGraph.footer.playing', {
                          current: playbackIndex,
                          total: playbackFrames,
                        })}
                      </span>
                    ) : null}
                  </div>
                  <div className="mt-2 text-[11px] font-mono text-neutral-500">
                    {t('admin.contests.trafficGraph.footer.updatedAt', {
                      time: new Date().toLocaleString(i18n.language || 'en-US'),
                    })}
                  </div>
                </div>
              </Card>
            </div>
          </div>
        </div>
      </Modal>
    </>
  );
}

export default TrafficGraphModal;
