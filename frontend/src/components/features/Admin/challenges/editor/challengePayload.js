import { normalizeNetworkPolicy } from '@/components/features/Admin/challenges/editor/networkPolicy.js';

export const buildChallengePayload = (challenge, mode, initialCompose = []) => {
  const payload = {
    name: challenge.name,
    description: challenge.description,
    category: challenge.category,
    type: challenge.type,
  };
  let flags;
  if (challenge.type === 'pods') flags = [];
  else if (mode === 'add') flags = challenge.flags.map((flag) => (typeof flag === 'string' ? flag : flag.value || ''));
  else if (mode === 'edit')
    flags = challenge.flags.map((flag) =>
      typeof flag === 'string' ? { id: 0, value: flag } : { id: flag.id || 0, value: flag.value || '' }
    );

  if (challenge.type === 'static') payload.flags = flags || [];
  else if (challenge.type === 'dynamic') {
    payload.flags = flags || [];
    payload.generator_image = challenge.generator_image || '';
  } else if (challenge.type === 'pods') {
    if (mode === 'add') payload.docker_compose = challenge.docker_compose || '';
    else {
      // Compare raw text against the fetched snapshot, not guide-generated YAML.
      initialCompose.forEach((original) => {
        if (original.id === challenge.id && original.value !== challenge.docker_compose) {
          payload.docker_compose = challenge.docker_compose;
        }
      });
    }
    payload.network_policies =
      challenge.network_policies?.map((rawPolicy) => {
        const policy = normalizeNetworkPolicy(rawPolicy);
        return {
          service: policy.service,
          allow_internet: policy.allow_internet,
          egress_cidrs: policy.egress_cidrs,
          allow_external: policy.allow_external,
          peers: policy.peers,
          from: policy.from,
          to: policy.to,
        };
      }) || [];
  }
  return payload;
};
