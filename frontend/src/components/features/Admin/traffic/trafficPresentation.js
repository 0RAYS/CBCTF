export const DEFAULT_SLICE_MS = 1000;
export const MIN_SLICE_MS = 1;
export const INPUT_STEP_MS = 100;

export const clamp = (value, min, max) => Math.min(Math.max(value, min), max);

export const formatBytes = (value) => {
  const bytes = Number(value || 0);
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(1)} GB`;
};

export const formatDurationMs = (value) => {
  const durationMs = Number(value || 0);
  if (durationMs < 1000) return `${durationMs} ms`;
  if (durationMs < 60000) return `${(durationMs / 1000).toFixed(durationMs % 1000 === 0 ? 0 : 2)} s`;
  const minutes = Math.floor(durationMs / 60000);
  const seconds = (durationMs % 60000) / 1000;
  if (seconds === 0) return `${minutes}m`;
  return `${minutes}m ${seconds.toFixed(seconds % 1 === 0 ? 0 : 2)}s`;
};

export const ellipsis = (value, length = 18) => {
  const text = value || '';
  return text.length > length ? `${text.slice(0, length)}...` : text;
};

const estimateTagChipWidth = (label) => Math.max(56, Math.ceil(String(label || '').length * 7.4 + 24));

export const resolveVisibleInlineItems = (items, containerWidth) => {
  if (!items.length) return { visibleItems: [], hiddenItems: [] };
  if (!containerWidth || containerWidth < 120) return { visibleItems: items, hiddenItems: [] };
  const gap = 8;
  let usedWidth = 0;
  let visibleCount = 0;
  for (let index = 0; index < items.length; index += 1) {
    const chipWidth = estimateTagChipWidth(items[index]);
    const remaining = items.length - index - 1;
    const additionalGap = visibleCount > 0 ? gap : 0;
    const reserveWidth = remaining > 0 ? gap + estimateTagChipWidth(`+${remaining}`) : 0;
    if (usedWidth + additionalGap + chipWidth + reserveWidth > containerWidth) break;
    usedWidth += additionalGap + chipWidth;
    visibleCount = index + 1;
  }
  if (visibleCount === 0) visibleCount = Math.min(1, items.length);
  return { visibleItems: items.slice(0, visibleCount), hiddenItems: items.slice(visibleCount) };
};

export const edgeTone = (direction) => {
  if (direction === 'ingress') return { hard: '#597ef7', soft: 'rgba(89,126,247,0.24)' };
  if (direction === 'egress') return { hard: '#a3a3a3', soft: 'rgba(163,163,163,0.18)' };
  return { hard: '#737373', soft: 'rgba(115,115,115,0.18)' };
};

export const edgeChipClass = (direction) => {
  if (direction === 'ingress') return 'border-geek-400/30 bg-geek-400/10 text-geek-400';
  if (direction === 'egress') return 'border-neutral-400/30 bg-neutral-400/10 text-neutral-300';
  return 'border-neutral-500/30 bg-neutral-500/10 text-neutral-300';
};

export const buildCompactNodeLabel = (node) => {
  if (node.service) return ellipsis(node.service, 16);
  if (node.kind === 'victim') return ellipsis(node.label || node.ip || node.id, 16);
  return ellipsis(node.ip || node.label || node.id, 16);
};

export const formatProcessLabel = (process) => {
  if (!process) return '';
  const name = process.process_name || process.dominant_process || '';
  if (name && process.pid) return `${name} (${process.pid})`;
  if (name) return name;
  if (process.pid) return `PID ${process.pid}`;
  return '';
};

export const buildCompactNodeMeta = (node) => {
  const process = node.dominant_process || formatProcessLabel((node.processes || [])[0]);
  if (process) return ellipsis(process, 22);
  if (node.service) return ellipsis(node.ip || node.id, 22);
  return `${formatBytes(node.bytes)} / ${node.connections || 0} conn`;
};

export const sanitizeSlice = (value) => Math.max(MIN_SLICE_MS, Math.round(Number(value) || 0));

export function filterTraffic(topology, universeNodes, protocolFilter) {
  const rawNodes = topology?.nodes || [];
  const rawEdges = topology?.edges || [];
  const allProtocols = [...new Set([...rawNodes, ...rawEdges].flatMap((item) => item.protocols || []))].sort();
  const universeProtocols = universeNodes
    ? [...new Set(universeNodes.flatMap((node) => node.protocols || []))].sort()
    : allProtocols;
  const nodes = protocolFilter.size
    ? rawNodes.filter((node) => (node.protocols || []).some((p) => protocolFilter.has(p)))
    : rawNodes;
  const visibleNodeIds = new Set(nodes.map((node) => node.id));
  const edges = protocolFilter.size
    ? rawEdges.filter(
        (edge) =>
          (edge.protocols || []).some((p) => protocolFilter.has(p)) &&
          visibleNodeIds.has(edge.source) &&
          visibleNodeIds.has(edge.target)
      )
    : rawEdges;
  // Filtering affects the graph only; summary and rankings remain server aggregates.
  return { nodes, edges, universeProtocols, activeProtocolSet: new Set(allProtocols) };
}

export function resolveTrafficSelection(nodes, edges, selectedNodeId, selectedEdgeId) {
  return {
    selectedEdge: edges.find((edge) => edge.id === selectedEdgeId) || edges[0] || null,
    selectedNode:
      nodes.find((node) => node.id === selectedNodeId) ||
      nodes.find((node) => node.kind === 'victim') ||
      nodes[0] ||
      null,
  };
}
