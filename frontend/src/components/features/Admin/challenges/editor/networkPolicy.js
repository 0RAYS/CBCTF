// 出公网范围：0.0.0.0/0 扣除私有网段，私网需要题目显式放行。
export const PUBLIC_EGRESS_EXCEPT = [
  '10.0.0.0/8',
  '172.16.0.0/12',
  '192.168.0.0/16',
  '100.64.0.0/10',
  '169.254.0.0/16',
];

const ipToNumber = (ip) => {
  const parts = String(ip || '')
    .trim()
    .split('.')
    .map((part) => Number(part));
  if (parts.length !== 4 || parts.some((part) => !Number.isInteger(part) || part < 0 || part > 255)) return null;
  return parts.reduce((value, part) => value * 256 + part, 0);
};

export const cidrContainsIp = (cidr, ip) => {
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

export const blockAllowsIp = (block, ip) => {
  if (!cidrContainsIp(block?.cidr, ip)) return false;
  return !(block?.except || []).some((exceptCidr) => cidrContainsIp(exceptCidr, ip));
};

export const isPublicIp = (ip) => !PUBLIC_EGRESS_EXCEPT.some((cidr) => cidrContainsIp(cidr, ip));

// isVpcIp 判断地址是否落在题目的 VPC 网段内。
export const isVpcIp = (ip, context = {}) => (context.vpcCIDRs || []).some((cidr) => cidrContainsIp(cidr, ip));

export const defaultNetworkPolicy = (service = '') => ({
  service,
  allow_internet: false,
  egress_cidrs: [],
  allow_external: true,
  restrict_ingress: true,
  peers: [],
  from: [],
  to: [],
});

export const normalizePolicyRule = (rule = {}) => ({
  cidr: rule.cidr || rule.CIDR || '',
  except: Array.isArray(rule.except) ? rule.except : Array.isArray(rule.Except) ? rule.Except : [],
});

export const normalizePolicyRules = (rules = []) => (Array.isArray(rules) ? rules.map(normalizePolicyRule) : []);

export const normalizeNetworkPolicy = (policy = {}) => {
  const from = normalizePolicyRules(policy.from);
  const peers = Array.isArray(policy.peers) ? policy.peers.filter(Boolean) : [];
  const declared = policy.allow_external;
  return {
    service: policy.service || '',
    allow_internet: policy.allow_internet === true,
    egress_cidrs: Array.isArray(policy.egress_cidrs) ? policy.egress_cidrs.filter(Boolean) : [],
    // null 表示旧条目没有声明入站限制，按旧语义处理。
    allow_external: declared === undefined || declared === null ? null : declared === true,
    peers,
    from,
    to: normalizePolicyRules(policy.to),
  };
};

// buildNetworkContext 计算容器网卡信息：所在网段、同 subnet 容器与地址。
export const buildNetworkContext = (config) => {
  const services = config.services || [];
  const subnetByName = new Map();
  (config.networks || []).forEach((network) => {
    const name = network.name?.trim();
    const subnet = network.subnet?.trim();
    if (name && subnet) subnetByName.set(name, subnet);
  });
  const entries = services.map((service, index) => {
    const name = service.containerName.trim() || service.name.trim() || `service${index + 1}`;
    const networks = (service.networks || []).filter((network) => network.name?.trim());
    return {
      name,
      networkNames: networks.map((network) => network.name.trim()),
      cidrs: [...new Set(networks.map((network) => subnetByName.get(network.name.trim())).filter(Boolean))],
      ips: networks.map((network) => network.ipv4Address?.trim()).filter(Boolean),
    };
  });
  const context = { ips: {}, subnetCIDRs: {}, subnetPeers: {}, vpcCIDRs: [] };
  entries.forEach((entry) => {
    context.ips[entry.name] = entry.ips;
    context.subnetCIDRs[entry.name] = entry.cidrs;
    context.subnetPeers[entry.name] = entries
      .filter(
        (other) =>
          other.name !== entry.name &&
          (other.networkNames.some((name) => entry.networkNames.includes(name)) ||
            other.cidrs.some((cidr) => entry.cidrs.includes(cidr)))
      )
      .map((other) => other.name);
  });
  context.vpcCIDRs = [...new Set(entries.flatMap((entry) => entry.cidrs))];
  return context;
};

// resolvePolicies 汇总每个容器的有效策略，语义与后端 ResolveNetworkPolicy 一致。
export const resolvePolicies = (policies, targets, context = {}) => {
  const list = (policies ?? []).map(normalizeNetworkPolicy);
  return targets.map((target) => resolveServicePolicy(list, target.service, context));
};

// resolveServicePolicy 汇总作用于单个容器的策略条目。
export const resolveServicePolicy = (list, service, context = {}) => {
  const matched = list.filter((policy) => !policy.service || policy.service === service);
  const merged =
    matched.length === 0
      ? defaultNetworkPolicy(service)
      : { ...defaultNetworkPolicy(service), ...mergePolicies(matched) };
  merged.peers = mergePeers(list, service, merged.peers);
  merged.egress_cidrs = [...new Set(merged.egress_cidrs)];
  merged.subnet_cidrs = context.subnetCIDRs?.[service] ?? [];
  merged.subnet_peers = context.subnetPeers?.[service] ?? [];
  merged.peer_ips = merged.peers.flatMap((peer) => context.ips?.[peer] ?? []);
  return merged;
};

// resolveSharedPolicy 汇总共享 Pod 网络下唯一 Pod 的策略，所有条目都作用于它。
export const resolveSharedPolicy = (policies) => {
  const list = (policies ?? []).map(normalizeNetworkPolicy);
  const merged = list.length === 0 ? defaultNetworkPolicy('') : { ...defaultNetworkPolicy(''), ...mergePolicies(list) };
  merged.peers = [...new Set(merged.peers)].sort();
  merged.egress_cidrs = [...new Set(merged.egress_cidrs)];
  merged.subnet_cidrs = [];
  merged.subnet_peers = [];
  merged.peer_ips = [];
  return merged;
};

// mergePolicies 合并命中的策略条目，规则之间取并集。
const mergePolicies = (policies) => {
  const merged = {
    allow_internet: false,
    allow_external: false,
    restrict_ingress: false,
    legacy_open: false,
    egress_cidrs: [],
    peers: [],
    from: [],
    to: [],
  };
  policies.forEach((policy) => {
    merged.allow_internet = merged.allow_internet || policy.allow_internet;
    if (policy.allow_external !== null) {
      merged.allow_external = merged.allow_external || policy.allow_external;
      merged.restrict_ingress = true;
    }
    if (policy.from.length > 0 || policy.peers.length > 0) merged.restrict_ingress = true;
    if (policy.allow_external === null && policy.from.length === 0 && policy.peers.length === 0)
      merged.legacy_open = true;
    merged.egress_cidrs = merged.egress_cidrs.concat(policy.egress_cidrs);
    merged.peers = merged.peers.concat(policy.peers);
    merged.from = merged.from.concat(policy.from);
    merged.to = merged.to.concat(policy.to);
  });
  if (merged.legacy_open) {
    // 旧条目没有声明入站限制，保持历史上的不限入站行为。
    merged.allow_external = true;
    merged.restrict_ingress = false;
  }
  return merged;
};

// mergePeers 取互通关系的对称闭包：任一侧声明的互通对两端同时生效。
const mergePeers = (policies, service, own) => {
  const peers = new Set(own.filter((peer) => peer && peer !== service));
  policies.forEach((policy) => {
    if (!policy.service || policy.service === service) return;
    if (policy.peers.includes(service)) peers.add(policy.service);
  });
  return [...peers].sort();
};

// policyAllowsEgressIp 判断容器能否访问目标 IP。
export const policyAllowsEgressIp = (policy, ip, context = {}) => {
  if (!ip) return false;
  if (policy.allow_internet && isPublicIp(ip) && !isVpcIp(ip, context)) return true;
  if (policy.subnet_cidrs.some((cidr) => cidrContainsIp(cidr, ip))) return true;
  if (policy.peer_ips.includes(ip)) return true;
  if (policy.egress_cidrs.some((cidr) => cidrContainsIp(cidr, ip))) return true;
  return policy.to.some((block) => blockAllowsIp(block, ip));
};

// policyAllowsIngressIp 判断容器是否接受来源 IP 的入站流量。
export const policyAllowsIngressIp = (policy, ip, context = {}) => {
  if (!policy.restrict_ingress) return true;
  if (!ip) return false;
  if (policy.subnet_cidrs.some((cidr) => cidrContainsIp(cidr, ip))) return true;
  if (policy.peer_ips.includes(ip)) return true;
  if (policy.from.some((block) => blockAllowsIp(block, ip))) return true;
  return policy.allow_external && !isVpcIp(ip, context);
};
