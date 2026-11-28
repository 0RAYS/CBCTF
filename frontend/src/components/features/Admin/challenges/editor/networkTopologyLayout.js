// 张力布局：按连通关系做应力优化（Kamada-Kawai 式 Guttman 变换），让节点散开成有张力的图形，而不是排成一条直线。

const HOP_LENGTH = 150;
const MIN_HOP_LENGTH = 54;
const NODE_GAP = 16;
const MARGIN = 18;
const INTERNET_GAP = 46;
const INTERNET_WIDTH = 108;
const INTERNET_HEIGHT = 28;
const COMPONENT_GAP = 44;
const SHELF_GAP = 40;
const ITERATIONS = 260;
const DAMPING = 0.6;
const MAX_SCALE = 1.25;
const LAN_WEIGHT = 2.2;
const LAN_LENGTH = 0.65;
const PEER_WEIGHT = 1.6;
const PEER_LENGTH = 1;

const NODE_CARD_WIDTH = 152;

const estimateNodeHeight = (node) => {
  const chips = node.networks.length ? node.networks.length * 40 + (node.networks.length - 1) * 4 : 16;
  return 16 + 20 + (node.image ? 18 : 0) + 8 + chips + 6;
};

// buildGraph 把同网段容器与互通容器连成无向图，权重表达关系亲疏。
const buildGraph = (topology) => {
  const adjacency = new Map();
  const ensure = (id) => {
    if (!adjacency.has(id)) adjacency.set(id, new Map());
    return adjacency.get(id);
  };
  const link = (a, b, weight, length) => {
    ensure(a);
    ensure(b);
    const existing = adjacency.get(a).get(b);
    if (existing) {
      existing.weight = Math.max(existing.weight, weight);
      existing.length = Math.min(existing.length, length);
      return;
    }
    adjacency.get(a).set(b, { weight, length });
    adjacency.get(b).set(a, { weight, length });
  };
  topology.nodes.forEach((node) => ensure(node.id));
  topology.lanes.forEach((lane) => {
    lane.members.forEach((source, index) => {
      lane.members.slice(index + 1).forEach((target) => link(source, target, LAN_WEIGHT, LAN_LENGTH));
    });
  });
  topology.connections.forEach((connection) => {
    if (connection.allowed) link(connection.source.id, connection.target.id, PEER_WEIGHT, PEER_LENGTH);
  });
  return adjacency;
};

const findComponents = (adjacency) => {
  const visited = new Set();
  const components = [];
  adjacency.forEach((_, id) => {
    if (visited.has(id)) return;
    const stack = [id];
    const members = [];
    visited.add(id);
    while (stack.length) {
      const current = stack.pop();
      members.push(current);
      adjacency.get(current).forEach((_, neighbor) => {
        if (visited.has(neighbor)) return;
        visited.add(neighbor);
        stack.push(neighbor);
      });
    }
    components.push(members);
  });
  return components;
};

const shortestPaths = (adjacency, source) => {
  const distances = new Map([[source, 0]]);
  const settled = new Set();
  for (;;) {
    let current = null;
    let best = Infinity;
    distances.forEach((value, id) => {
      if (settled.has(id) || value >= best) return;
      best = value;
      current = id;
    });
    if (current === null) break;
    settled.add(current);
    adjacency.get(current).forEach((edge, neighbor) => {
      const candidate = best + edge.length;
      if (candidate < (distances.get(neighbor) ?? Infinity)) distances.set(neighbor, candidate);
    });
  }
  return distances;
};

const measureStress = (positions, relations) => {
  let stress = 0;
  relations.forEach((row, id) => {
    const origin = positions.get(id);
    row.forEach(({ distance, weight }, other) => {
      const target = positions.get(other);
      stress += weight * (Math.hypot(origin.x - target.x, origin.y - target.y) - distance) ** 2;
    });
  });
  return stress;
};

// majorize 一步 Guttman 变换：节点朝邻居的加权质心与理想距离方向靠拢。
const majorize = (positions, relations) => {
  const next = new Map();
  positions.forEach((point, id) => {
    let numeratorX = 0;
    let numeratorY = 0;
    let total = 0;
    relations.get(id).forEach(({ distance, weight }, other) => {
      const target = positions.get(other);
      let dx = point.x - target.x;
      let dy = point.y - target.y;
      let gap = Math.hypot(dx, dy);
      if (gap < 0.01) {
        dx = 0.01;
        dy = 0;
        gap = 0.01;
      }
      numeratorX += weight * (target.x + (distance * dx) / gap);
      numeratorY += weight * (target.y + (distance * dy) / gap);
      total += weight;
    });
    if (total <= 0) {
      next.set(id, point);
      return;
    }
    next.set(id, {
      x: point.x + DAMPING * (numeratorX / total - point.x),
      y: point.y + DAMPING * (numeratorY / total - point.y),
    });
  });
  return next;
};

const relaxOverlaps = (positions, members, sizes) => {
  for (let iteration = 0; iteration < 80; iteration += 1) {
    let moved = false;
    for (let i = 0; i < members.length; i += 1) {
      for (let j = i + 1; j < members.length; j += 1) {
        const origin = positions.get(members[i]);
        const target = positions.get(members[j]);
        const sizeA = sizes.get(members[i]);
        const sizeB = sizes.get(members[j]);
        const overlapX = (sizeA.width + sizeB.width) / 2 + NODE_GAP - Math.abs(target.x - origin.x);
        const overlapY = (sizeA.height + sizeB.height) / 2 + NODE_GAP - Math.abs(target.y - origin.y);
        if (overlapX <= 0 || overlapY <= 0) continue;
        moved = true;
        if (overlapX < overlapY) {
          const shift = (Math.sign(target.x - origin.x) || 1) * overlapX * 0.35;
          origin.x -= shift;
          target.x += shift;
        } else {
          const shift = (Math.sign(target.y - origin.y) || 1) * overlapY * 0.35;
          origin.y -= shift;
          target.y += shift;
        }
      }
    }
    if (!moved) break;
  }
};

// layoutComponent 在一个连通分量内部求应力最小的排布。
const layoutComponent = (members, adjacency, sizes, hopLength) => {
  const positions = new Map();
  const radius = (hopLength * members.length) / (2 * Math.PI);
  members.forEach((id, index) => {
    if (members.length === 1) {
      positions.set(id, { x: 0, y: 0 });
      return;
    }
    const angle = (2 * Math.PI * index) / members.length - (Math.PI * 3) / 4;
    positions.set(id, { x: Math.cos(angle) * radius, y: Math.sin(angle) * radius });
  });

  const relations = new Map();
  members.forEach((source) => {
    const distances = shortestPaths(adjacency, source);
    const row = new Map();
    members.forEach((target) => {
      if (source === target) return;
      const hops = distances.get(target);
      if (hops === undefined) return;
      const edge = adjacency.get(source).get(target);
      const size = sizes.get(source);
      const other = sizes.get(target);
      const minimum = Math.max((size.width + other.width) / 2, (size.height + other.height) / 2) + NODE_GAP;
      row.set(target, {
        distance: Math.max(hops * hopLength, minimum),
        weight: edge ? edge.weight : 0.6,
      });
    });
    relations.set(source, row);
  });

  let current = positions;
  let best = new Map(positions);
  let bestStress = Infinity;
  for (let iteration = 0; iteration < ITERATIONS; iteration += 1) {
    current = majorize(current, relations);
    const stress = measureStress(current, relations);
    if (stress < bestStress) {
      bestStress = stress;
      best = new Map(current);
    }
  }

  relaxOverlaps(best, members, sizes);
  const bounds = { minX: Infinity, minY: Infinity, maxX: -Infinity, maxY: -Infinity };
  members.forEach((id) => {
    const point = best.get(id);
    const size = sizes.get(id);
    bounds.minX = Math.min(bounds.minX, point.x - size.width / 2);
    bounds.maxX = Math.max(bounds.maxX, point.x + size.width / 2);
    bounds.minY = Math.min(bounds.minY, point.y - size.height / 2);
    bounds.maxY = Math.max(bounds.maxY, point.y + size.height / 2);
  });
  members.forEach((id) => {
    const point = best.get(id);
    point.x -= bounds.minX;
    point.y -= bounds.minY;
  });
  return { positions: best, width: bounds.maxX - bounds.minX, height: bounds.maxY - bounds.minY };
};

// packComponents 把互不连通的分组按行摆放，组与组之间留出空隙。
const packComponents = (components, available) => {
  const placed = [];
  let cursorX = 0;
  let cursorY = 0;
  let shelfHeight = 0;
  components.forEach((component) => {
    if (cursorX > 0 && cursorX + component.width > available) {
      cursorY += shelfHeight + SHELF_GAP;
      cursorX = 0;
      shelfHeight = 0;
    }
    placed.push({ component, x: cursorX, y: cursorY });
    cursorX += component.width + COMPONENT_GAP;
    shelfHeight = Math.max(shelfHeight, component.height);
  });
  return {
    placed,
    width: placed.reduce((value, item) => Math.max(value, item.x + item.component.width), 0),
    height: placed.reduce((value, item) => Math.max(value, item.y + item.component.height), 0),
  };
};

const boundaryPoint = (node, ux, uy) => {
  const rx = node.width / 2;
  const ry = node.height / 2;
  const factor = 1 / Math.sqrt((ux * ux) / (rx * rx) + (uy * uy) / (ry * ry));
  return { x: node.x + ux * factor, y: node.y + uy * factor };
};

const quadratic = (start, control, end, t) => ({
  x: (1 - t) * (1 - t) * start.x + 2 * t * (1 - t) * control.x + t * t * end.x,
  y: (1 - t) * (1 - t) * start.y + 2 * t * (1 - t) * control.y + t * t * end.y,
});

// edgePath 在两张卡片之间画一条曲线，offset 把多条连线分到两侧。
const edgePath = (from, to, offset) => {
  const dx = to.x - from.x;
  const dy = to.y - from.y;
  const distance = Math.hypot(dx, dy) || 1;
  const nx = -dy / distance;
  const ny = dx / distance;
  const startCenter = { x: from.x + nx * offset, y: from.y + ny * offset };
  const endCenter = { x: to.x + nx * offset, y: to.y + ny * offset };
  const ux = (endCenter.x - startCenter.x) / distance;
  const uy = (endCenter.y - startCenter.y) / distance;
  const start = boundaryPoint(from, ux, uy);
  const end = boundaryPoint(to, -ux, -uy);
  const bow = Math.min(14, distance * 0.04);
  const control = { x: (start.x + end.x) / 2 - uy * bow, y: (start.y + end.y) / 2 + ux * bow };
  const first = quadratic(start, control, end, 1 / 3);
  const second = quadratic(start, control, end, 2 / 3);
  const mid = quadratic(start, control, end, 0.5);
  return `M ${start.x} ${start.y} Q ${first.x} ${first.y} ${mid.x} ${mid.y} Q ${second.x} ${second.y} ${end.x} ${end.y}`;
};

// buildEdges 每对容器只画一条连线：双向可达画双向箭头，单向可达按放行方向画箭头，其余画阻断标记。
const buildEdges = (topology, positions, internet) => {
  const edges = [];
  const pairs = new Map();
  topology.connections.forEach((connection) => {
    const key = [connection.source.id, connection.target.id].sort().join('|');
    const pair = pairs.get(key) || { key, members: [] };
    pair.members.push(connection);
    pairs.set(key, pair);
  });
  pairs.forEach((pair) => {
    const reachable = pair.members.filter((member) => member.allowed);
    const forward = reachable[0] || pair.members[0];
    const source = positions[forward.source.id];
    const target = positions[forward.target.id];
    if (!source || !target) return;
    const path = edgePath(source, target, 0);
    edges.push({
      id: pair.key,
      source: forward.source.id,
      target: forward.target.id,
      allowed: reachable.length > 0,
      bidirectional: reachable.length === pair.members.length,
      d: path,
    });
  });
  (topology.routes || []).forEach((route) => {
    const source = positions[route.source];
    const target = positions[route.target];
    if (!source || !target) return;
    edges.push({
      id: route.id,
      source: route.source,
      target: route.target,
      allowed: true,
      routed: true,
      bidirectional: true,
      d: edgePath(source, target, 0),
    });
  });
  const egress = topology.egress
    .map((connection) => {
      const source = positions[connection.source.id];
      if (!source) return null;
      return { connection, source, angle: Math.atan2(source.y - internet.y, source.x - internet.x) };
    })
    .filter(Boolean)
    .sort((a, b) => a.angle - b.angle);
  egress.forEach((item, index) => {
    const spread = egress.length > 1 ? (index / (egress.length - 1) - 0.5) * 44 : 0;
    const path = edgePath(item.source, internet, spread);
    edges.push({
      id: item.connection.id,
      source: item.connection.source.id,
      target: 'internet',
      allowed: item.connection.allowed,
      bidirectional: false,
      d: path,
    });
  });
  return edges;
};

// layoutTopology 计算拓扑图的像素坐标：分组排布、应力优化、公网节点置顶。
export const layoutTopology = (topology, options = {}) => {
  const services = topology.nodes || [];
  const width = Math.max(options.width || 0, 240);
  if (!services.length) return { width, height: 0, positions: {}, internet: null, edges: [] };

  const sizes = new Map(
    services.map((node) => [node.id, { width: NODE_CARD_WIDTH, height: estimateNodeHeight(node) }])
  );
  const adjacency = buildGraph(topology);
  const groups = findComponents(adjacency);
  const available = Math.max(width - MARGIN * 2, 240);

  let hopLength = HOP_LENGTH;
  let packed = null;
  for (let attempt = 0; attempt < 4; attempt += 1) {
    const components = groups.map((members) => layoutComponent(members, adjacency, sizes, hopLength));
    packed = packComponents(components, available);
    if (packed.width <= available) break;
    hopLength = Math.max(MIN_HOP_LENGTH, hopLength * (available / packed.width) * 0.95);
  }

  const scale = Math.max(1, Math.min(MAX_SCALE, available / Math.max(packed.width, 1)));
  const placed = {};
  packed.placed.forEach(({ component, x, y }) => {
    component.positions.forEach((point, id) => {
      placed[id] = { x: x + point.x * scale, y: y + point.y * scale };
    });
  });

  const bounds = services.reduce(
    (accumulator, node) => {
      const point = placed[node.id];
      const size = sizes.get(node.id);
      accumulator.minX = Math.min(accumulator.minX, point.x - size.width / 2);
      accumulator.maxX = Math.max(accumulator.maxX, point.x + size.width / 2);
      accumulator.minY = Math.min(accumulator.minY, point.y - size.height / 2);
      accumulator.maxY = Math.max(accumulator.maxY, point.y + size.height / 2);
      return accumulator;
    },
    { minX: Infinity, minY: Infinity, maxX: -Infinity, maxY: -Infinity }
  );

  const canvasWidth = Math.max(width, bounds.maxX - bounds.minX + MARGIN * 2);
  const contentTop = MARGIN + INTERNET_HEIGHT + INTERNET_GAP;
  const shiftX = MARGIN + (canvasWidth - MARGIN * 2 - (bounds.maxX - bounds.minX)) / 2 - bounds.minX;
  const shiftY = contentTop - bounds.minY;
  const positions = {};
  services.forEach((node) => {
    const point = placed[node.id];
    const size = sizes.get(node.id);
    positions[node.id] = {
      x: point.x + shiftX,
      y: point.y + shiftY,
      width: size.width,
      height: size.height,
    };
  });

  const egressTargets = topology.egress
    .filter((connection) => connection.allowed)
    .map((connection) => connection.source.id);
  const anchor = egressTargets.length ? egressTargets : services.map((node) => node.id);
  const anchorX = anchor.reduce((value, id) => value + positions[id].x, 0) / anchor.length;
  const internet = {
    id: 'internet',
    x: anchorX,
    y: MARGIN + INTERNET_HEIGHT / 2,
    width: INTERNET_WIDTH,
    height: INTERNET_HEIGHT,
  };
  const height = contentTop + (bounds.maxY - bounds.minY) + MARGIN;

  return { width: canvasWidth, height, positions, internet, edges: buildEdges(topology, positions, internet) };
};
