import { getGuideServiceTargets } from '@/components/features/Admin/challenges/editor/guideModel.js';
import {
  buildNetworkContext,
  policyAllowsEgressIp,
  policyAllowsIngressIp,
  resolvePolicies,
  resolveSharedPolicy,
} from '@/components/features/Admin/challenges/editor/networkPolicy.js';

export const INTERNET_ID = 'internet';

const LANE_COLORS = ['#38bdf8', '#a78bfa', '#fbbf24', '#34d399', '#fb7185', '#22d3ee', '#c084fc', '#a3e635'];

const nodeIps = (node) => node.networks.map((network) => network.ip).filter((ip) => ip && ip !== '-');

// 靶场节点角色：网关由 x-range 决定，其余取显式 x-role，未声明视为节点。
const rangeRole = (config, service) => {
  const range = config.range;
  if (!range) return '';
  if (service.name.trim() === (range.gateway || '').trim()) return 'gateway';
  return (service.role || '').trim() || 'node';
};

// buildNetworkTopology 根据容器网卡生成拓扑，并按网络策略标注可达性。
export const buildNetworkTopology = (config, policies = []) => {
  const services = config.services || [];
  const targets = getGuideServiceTargets(config);
  const context = buildNetworkContext(config);
  const sharedNetwork = (config.networks || []).length === 0;
  const resolved = sharedNetwork
    ? targets.map((target) => ({ ...resolveSharedPolicy(policies), service: target.service }))
    : resolvePolicies(policies, targets, context);
  const policyByService = new Map(resolved.map((policy) => [policy.service, policy]));

  const nodes = services.map((service, index) => ({
    id: targets[index].service,
    label: targets[index].label,
    image: service.image,
    role: rangeRole(config, service),
    networks: (service.networks || [])
      .filter((network) => network.name || network.ipv4Address)
      .map((network) => ({ name: network.name || 'default', ip: network.ipv4Address || '-' })),
  }));

  const laneNames = [];
  (config.networks || []).forEach((network) => {
    const name = network.name?.trim();
    if (name && !laneNames.includes(name)) laneNames.push(name);
  });
  nodes.forEach((node) => {
    node.networks.forEach((network) => {
      if (!laneNames.includes(network.name)) laneNames.push(network.name);
    });
  });

  const lanes = laneNames.map((name, index) => ({
    id: name,
    name,
    color: LANE_COLORS[index % LANE_COLORS.length],
    members: nodes.filter((node) => node.networks.some((network) => network.name === name)).map((node) => node.id),
  }));
  const colorByLane = new Map(lanes.map((lane) => [lane.id, lane.color]));
  nodes.forEach((node) => {
    node.networks.forEach((network) => {
      network.color = colorByLane.get(network.name);
    });
  });

  const connections = [];
  nodes.forEach((source) => {
    nodes.forEach((target) => {
      if (source.id === target.id) return;
      const sourcePolicy = policyByService.get(source.id);
      const targetPolicy = policyByService.get(target.id);
      const outbound =
        sourcePolicy.peers.includes(target.id) ||
        sourcePolicy.subnet_peers.includes(target.id) ||
        nodeIps(target).some((ip) => policyAllowsEgressIp(sourcePolicy, ip, context));
      const inbound =
        targetPolicy.peers.includes(source.id) ||
        targetPolicy.subnet_peers.includes(source.id) ||
        nodeIps(source).some((ip) => policyAllowsIngressIp(targetPolicy, ip, context));
      const allowed = sharedNetwork || (outbound && inbound);
      const reasonKey = sharedNetwork
        ? 'sharedNetwork'
        : !outbound
          ? 'outboundBlocked'
          : !inbound
            ? 'inboundBlocked'
            : 'allowed';
      connections.push({
        id: `${source.id}->${target.id}`,
        source,
        target,
        allowed,
        reasonKey,
        networks: source.networks
          .map((network) => network.name)
          .filter((name) => target.networks.some((network) => network.name === name)),
      });
    });
  });

  const egress = nodes.map((node) => {
    const policy = policyByService.get(node.id);
    return {
      id: `${node.id}->${INTERNET_ID}`,
      source: node,
      target: { id: INTERNET_ID, label: 'internet' },
      allowed: policy.allow_internet,
      reasonKey: policy.allow_internet ? 'egressAllowed' : 'egressBlocked',
      ranges: policy.egress_cidrs,
    };
  });

  // 靶场声明的网段互通是显式的，按网段各取一个代表节点连成虚线。
  const routes = ((config.range || {}).routes || [])
    .map((route) => {
      const from = (route.from || '').trim();
      const to = (route.to || '').trim();
      const source = lanes.find((lane) => lane.id === from)?.members[0];
      const target = lanes.find((lane) => lane.id === to)?.members[0];
      if (!source || !target || source === target) return null;
      return { id: `${from}->${to}`, from, to, source, target };
    })
    .filter(Boolean);

  return {
    nodes,
    lanes,
    internet: { id: INTERNET_ID, label: 'internet' },
    connections,
    egress,
    routes,
    sharedNetwork,
  };
};
