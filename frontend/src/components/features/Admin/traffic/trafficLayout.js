import { clamp, edgeTone, formatBytes } from './trafficPresentation.js';

export const VIEWBOX_WIDTH = 1120;
export const VIEWBOX_HEIGHT = 580;
const CENTER_X = VIEWBOX_WIDTH / 2;
export const CENTER_Y = VIEWBOX_HEIGHT / 2;
const TOP_BOUND = 92;
const BOTTOM_BOUND = VIEWBOX_HEIGHT - 92;
const PEER_WIDTH = 152;
const PEER_HEIGHT = 46;
const VICTIM_WIDTH = 170;
const VICTIM_HEIGHT = 54;
export const MIN_ZOOM = 0.78;
export const MAX_ZOOM = 1.9;
export const ZOOM_STEP = 0.12;
const LABEL_HEIGHT = 20;

const rectsOverlap = (a, b, padding = 0) =>
  a.x < b.x + b.w + padding && a.x + a.w + padding > b.x && a.y < b.y + b.h + padding && a.y + a.h + padding > b.y;

const sortByTraffic = (items) =>
  [...items].sort((a, b) => {
    const aService = a.service || (a.services || [])[0] || '';
    const bService = b.service || (b.services || [])[0] || '';
    if (aService !== bService) return aService.localeCompare(bService);
    if ((b.bytes || 0) !== (a.bytes || 0)) return (b.bytes || 0) - (a.bytes || 0);
    return String(a.ip || a.id || '').localeCompare(String(b.ip || b.id || ''));
  });

const distributeIntoLanes = (items, laneCount) => {
  const lanes = Array.from({ length: laneCount }, () => []);
  items.forEach((item, index) => lanes[index % laneCount].push(item));
  return lanes;
};

const buildColumnCandidates = (count, maxColumns, preferredColumns) => {
  const limit = Math.max(1, Math.min(maxColumns, count || 1));
  const preferred = clamp(preferredColumns || 1, 1, limit);
  const candidates = [];
  for (let offset = 0; offset < limit; offset += 1) {
    const higher = preferred + offset;
    const lower = preferred - offset;
    if (higher <= limit) candidates.push(higher);
    if (offset > 0 && lower >= 1) candidates.push(lower);
  }
  return [...new Set(candidates)];
};

const resolveGridLayout = (count, zone, config) => {
  const zoneWidth = zone.xMax - zone.xMin;
  const zoneHeight = zone.yMax - zone.yMin;
  const candidates = buildColumnCandidates(count, config.maxColumns, config.preferredColumns);
  for (const columns of candidates) {
    const rows = Math.ceil(count / columns);
    const width = Math.floor((zoneWidth - config.minGapX * Math.max(columns - 1, 0)) / columns);
    const height = Math.floor((zoneHeight - config.minGapY * Math.max(rows - 1, 0)) / rows);
    if (width < config.minWidth || height < config.minHeight) continue;
    const resolvedWidth = Math.min(config.baseWidth, width);
    const resolvedHeight = Math.min(config.baseHeight, height);
    const gapX = columns === 1 ? 0 : (zoneWidth - resolvedWidth * columns) / (columns - 1);
    const gapY = rows === 1 ? 0 : (zoneHeight - resolvedHeight * rows) / (rows - 1);
    return { columns, rows, width: resolvedWidth, height: resolvedHeight, gapX, gapY };
  }
  const fallbackColumns = candidates[0] || 1;
  const rows = Math.ceil(count / fallbackColumns);
  const width = Math.max(
    80,
    Math.floor((zoneWidth - config.minGapX * Math.max(fallbackColumns - 1, 0)) / fallbackColumns)
  );
  const height = Math.max(36, Math.floor((zoneHeight - config.minGapY * Math.max(rows - 1, 0)) / rows));
  const resolvedWidth = Math.min(config.baseWidth, width);
  const resolvedHeight = Math.min(config.baseHeight, height);
  return {
    columns: fallbackColumns,
    rows,
    width: resolvedWidth,
    height: resolvedHeight,
    gapX:
      fallbackColumns === 1 ? 0 : Math.max(8, (zoneWidth - resolvedWidth * fallbackColumns) / (fallbackColumns - 1)),
    gapY: rows === 1 ? 0 : Math.max(8, (zoneHeight - resolvedHeight * rows) / (rows - 1)),
  };
};

export const getViewportMetrics = (zoom) => ({ width: VIEWBOX_WIDTH / zoom, height: VIEWBOX_HEIGHT / zoom });

export const clampPanToViewport = (pan, zoom) => {
  const { width, height } = getViewportMetrics(zoom);
  if (zoom <= 1) return { x: (VIEWBOX_WIDTH - width) / 2, y: (VIEWBOX_HEIGHT - height) / 2 };
  const minX = Math.min(0, VIEWBOX_WIDTH - width);
  const maxX = Math.max(0, VIEWBOX_WIDTH - width);
  const minY = Math.min(0, VIEWBOX_HEIGHT - height);
  const maxY = Math.max(0, VIEWBOX_HEIGHT - height);
  return {
    x: minX === maxX ? minX : clamp(pan.x, minX, maxX),
    y: minY === maxY ? minY : clamp(pan.y, minY, maxY),
  };
};

export function buildPositions(nodes) {
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
    if (leftNodes.length + index <= rightNodes.length + Math.floor(index / 2))
      leftNodes.push({ ...node, side: 'left' });
    else rightNodes.push({ ...node, side: 'right' });
  });
  const zones = {
    left: { xMin: 36, xMax: 336, yMin: TOP_BOUND, yMax: BOTTOM_BOUND },
    center: { xMin: 370, xMax: VIEWBOX_WIDTH - 370, yMin: TOP_BOUND, yMax: BOTTOM_BOUND },
    right: { xMin: VIEWBOX_WIDTH - 336, xMax: VIEWBOX_WIDTH - 36, yMin: TOP_BOUND, yMax: BOTTOM_BOUND },
  };
  const placeSideNodes = (items, zone) => {
    if (!items.length) return;
    const layout = resolveGridLayout(items.length, zone, {
      baseWidth: PEER_WIDTH,
      baseHeight: PEER_HEIGHT,
      minWidth: 104,
      minHeight: 40,
      minGapX: 16,
      minGapY: 12,
      maxColumns: 3,
      preferredColumns: Math.min(3, Math.max(1, Math.ceil(items.length / 6))),
    });
    const lanes = distributeIntoLanes(items, layout.columns);
    const zoneWidth = zone.xMax - zone.xMin;
    const zoneHeight = zone.yMax - zone.yMin;
    const totalWidth = layout.columns * layout.width + Math.max(0, layout.columns - 1) * layout.gapX;
    const startX = zone.xMin + (zoneWidth - totalWidth) / 2 + layout.width / 2;
    lanes.forEach((lane, laneIndex) => {
      if (!lane.length) return;
      const laneHeight = lane.length * layout.height + Math.max(0, lane.length - 1) * layout.gapY;
      const startY = zone.yMin + (zoneHeight - laneHeight) / 2 + layout.height / 2;
      lane.forEach((node, index) => {
        positions.set(node.id, {
          x: startX + laneIndex * (layout.width + layout.gapX),
          y: startY + index * (layout.height + layout.gapY),
          w: layout.width,
          h: layout.height,
        });
      });
    });
  };
  placeSideNodes(sortByTraffic(leftNodes), zones.left);
  placeSideNodes(sortByTraffic(rightNodes), zones.right);
  const orderedCenter = sortByTraffic(centerNodes);
  if (!orderedCenter.length) return positions;
  const centerLayout = resolveGridLayout(orderedCenter.length, zones.center, {
    baseWidth: VICTIM_WIDTH,
    baseHeight: VICTIM_HEIGHT,
    minWidth: 112,
    minHeight: 46,
    minGapX: 18,
    minGapY: 14,
    maxColumns: 3,
    preferredColumns: orderedCenter.length === 1 ? 1 : orderedCenter.length <= 4 ? 2 : 3,
  });
  const rowCount = Math.ceil(orderedCenter.length / centerLayout.columns);
  const zoneHeight = zones.center.yMax - zones.center.yMin;
  const totalHeight = rowCount * centerLayout.height + Math.max(0, rowCount - 1) * centerLayout.gapY;
  const startY = zones.center.yMin + (zoneHeight - totalHeight) / 2 + centerLayout.height / 2;
  for (let row = 0; row < rowCount; row += 1) {
    const rowItems = orderedCenter.slice(row * centerLayout.columns, (row + 1) * centerLayout.columns);
    const rowWidth = rowItems.length * centerLayout.width + Math.max(0, rowItems.length - 1) * centerLayout.gapX;
    const rowStartX = CENTER_X - rowWidth / 2 + centerLayout.width / 2;
    rowItems.forEach((node, column) => {
      positions.set(node.id, {
        x: rowStartX + column * (centerLayout.width + centerLayout.gapX),
        y: startY + row * (centerLayout.height + centerLayout.gapY),
        w: centerLayout.width,
        h: centerLayout.height,
      });
    });
  }
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
    items.forEach((item, index) => offsets.set(item.key, (index / (items.length - 1) - 0.5) * spread));
  });
  return offsets;
}

export function buildLines(edges, positions) {
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

export function resolveEdgeLabels(lines, positions, selectedEdgeId) {
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
      labelRects.push({ x: placement.x - 4, y: placement.y - 3, w: placement.w + 8, h: placement.h + 6 });
    }
  });
  return placements;
}
