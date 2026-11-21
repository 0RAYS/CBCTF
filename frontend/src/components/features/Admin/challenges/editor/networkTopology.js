import { getGuideServiceTargets } from '@/components/features/Admin/challenges/editor/guideModel.js';
import {
  buildNetworkContext,
  policyAllowsEgressIp,
  policyAllowsIngressIp,
  resolvePolicies,
  resolveSharedPolicy,
} from '@/components/features/Admin/challenges/editor/networkPolicy.js';

export const INTERNET_ID = 'internet';

const TOP = 26;
const BOTTOM = 4;

const nodeIps = (node) => node.networks.map((network) => network.ip).filter((ip) => ip && ip !== '-');

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
    const name = node.networks[0]?.name || 'default';
    if (!laneNames.includes(name)) laneNames.push(name);
  });

  const laneWidth = 100 / Math.max(laneNames.length, 1);
  const lanes = laneNames.map((name, index) => ({
    id: name,
    name,
    x: laneWidth * (index + 0.5),
    width: laneWidth,
    nodes: [],
  }));
  nodes.forEach((node) => {
    const lane = lanes.find((item) => item.id === (node.networks[0]?.name || 'default'));
    lane.nodes.push(node);
  });
  lanes.forEach((lane) => {
    lane.nodes.forEach((node, index) => {
      node.x = lane.x;
      node.y = TOP + ((index + 0.5) * (100 - TOP - BOTTOM)) / lane.nodes.length;
    });
  });

  const internet = { id: INTERNET_ID, label: 'internet', x: 50, y: 5, networks: [] };

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
      target: internet,
      allowed: policy.allow_internet,
      reasonKey: policy.allow_internet ? 'egressAllowed' : 'egressBlocked',
      ranges: policy.egress_cidrs,
    };
  });

  return { nodes, lanes, internet, connections, egress, sharedNetwork };
};
