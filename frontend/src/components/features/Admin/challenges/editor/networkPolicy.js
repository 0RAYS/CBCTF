export const defaultNetworkPolicy = (target = {}) => ({
  ...(target.service ? { service: target.service } : {}),
  from: [{ cidr: '', except: [''] }],
  to: [{ cidr: '0.0.0.0/0', except: ['10.0.0.0/8', '172.16.0.0/12', '192.168.0.0/16', '100.64.0.0/10'] }],
});

export const normalizePolicyRule = (rule = {}) => ({
  cidr: rule.cidr || rule.CIDR || '',
  except: Array.isArray(rule.except) ? rule.except : Array.isArray(rule.Except) ? rule.Except : [],
});

export const normalizePolicyRules = (rules = []) => (Array.isArray(rules) ? rules.map(normalizePolicyRule) : []);

export const normalizeNetworkPolicy = (policy = {}) => ({
  ...policy,
  from: normalizePolicyRules(policy.from),
  to: normalizePolicyRules(policy.to),
});

const ipToNumber = (ip) => {
  const parts = String(ip || '')
    .trim()
    .split('.')
    .map((part) => Number(part));
  if (parts.length !== 4 || parts.some((part) => !Number.isInteger(part) || part < 0 || part > 255)) return null;
  return parts.reduce((value, part) => value * 256 + part, 0);
};

const cidrContainsIp = (cidr, ip) => {
  const [baseIp, prefixText = '32'] = String(cidr || '')
    .trim()
    .split('/');
  const base = ipToNumber(baseIp);
  const target = ipToNumber(ip);
  const prefix = Number(prefixText);
  if (base === null || target === null || !Number.isInteger(prefix) || prefix < 0 || prefix > 32) return false;
  if (prefix === 0) return true;
  const mask = (0xffffffff << (32 - prefix)) >>> 0;
  return (base & mask) === (target & mask);
};

const getPolicyBlocks = (blocks) =>
  (blocks ?? [])
    .map((block) => ({
      cidr: String(block?.cidr || block?.CIDR || '').trim(),
      except: block?.except || block?.Except || [],
    }))
    .filter((block) => block.cidr);

const blockAllowsIp = (block, ip) => {
  if (!cidrContainsIp(block.cidr, ip)) return false;
  return !(block.except || []).some((exceptCidr) => cidrContainsIp(exceptCidr, ip));
};
const blocksAllowAnyIp = (blocks, ips) => ips.some((ip) => blocks.some((block) => blockAllowsIp(block, ip)));

export const buildNetworkTopology = (config, policies = []) => {
  const services = config.services || [];
  const nodes = services.map((service, index) => {
    const name = service.containerName.trim() || service.name.trim() || `service${index + 1}`;
    return {
      id: name,
      label: name,
      image: service.image,
      networks: (service.networks || [])
        .filter((network) => network.name || network.ipv4Address)
        .map((network) => ({ name: network.name || 'default', ip: network.ipv4Address || '-' })),
    };
  });
  const positionedNodes = nodes.map((node, index) => {
    if (nodes.length === 1) return { ...node, x: 50, y: 50 };
    const angle = (Math.PI * 2 * index) / nodes.length - Math.PI / 2;
    return { ...node, x: 50 + Math.cos(angle) * 34, y: 50 + Math.sin(angle) * 34 };
  });
  const policyByService = new Map((policies ?? []).map((policy) => [policy.service, policy]));
  const getNodeIps = (node) => node.networks.map((network) => network.ip).filter((ip) => ipToNumber(ip) !== null);
  const connections = [];
  positionedNodes.forEach((source) => {
    positionedNodes.forEach((target) => {
      if (source.id === target.id) return;
      const sourcePolicy = policyByService.get(source.id);
      const targetPolicy = policyByService.get(target.id);
      const targetIps = getNodeIps(target);
      const sourceIps = getNodeIps(source);
      const toBlocks = getPolicyBlocks(sourcePolicy?.to);
      const fromBlocks = getPolicyBlocks(targetPolicy?.from);
      const outboundAllowed = targetIps.length > 0 && blocksAllowAnyIp(toBlocks, targetIps);
      const inboundAllowed = fromBlocks.length === 0 || blocksAllowAnyIp(fromBlocks, sourceIps);
      const sharedNetworks = source.networks
        .map((network) => network.name)
        .filter((name) => target.networks.some((network) => network.name === name));
      connections.push({
        id: `${source.id}->${target.id}`,
        source,
        target,
        allowed: outboundAllowed && inboundAllowed,
        reasonKey: !outboundAllowed ? 'outboundBlocked' : !inboundAllowed ? 'inboundBlocked' : 'allowed',
        networks: sharedNetworks.length ? sharedNetworks : [],
      });
    });
  });
  return { nodes: positionedNodes, connections };
};
