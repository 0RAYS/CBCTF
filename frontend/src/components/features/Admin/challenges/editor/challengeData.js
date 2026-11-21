export function challengeDraft(challenge = {}) {
  return {
    id: challenge.id || '',
    name: challenge.name || '',
    description: challenge.description || '',
    category: challenge.category || '',
    type: challenge.type || 'static',
    contest_id: challenge.contest_id || 0,
    contest_name: challenge.contest_name || '',
    hidden: !!challenge.hidden,
    attempt: challenge.attempt || 0,
    hints: [...(challenge.hints || [])],
    tags: [...(challenge.tags || [])],
    generator_image: challenge.generator_image || '',
    flags: (challenge.flags || []).map((flag) =>
      typeof flag === 'string' ? { id: 0, value: flag } : { id: flag.id || 0, value: flag.value || '' }
    ),
    docker_compose: challenge.docker_compose || '',
    network_policies: challenge.network_policies || [],
  };
}

export function applyChallengeType(previous, next) {
  if (next.type === previous.type) return next;
  if (next.type === 'pods') return { ...next, flags: [] };
  if (previous.type === 'pods') return { ...next, flags: [{ id: 0, value: '' }] };
  return next;
}

export function flagUpdatePayload(flag) {
  return {
    score_type: flag.score_type,
    score: flag.score,
    decay: flag.decay,
    min_score: flag.min_score,
  };
}

export function isFlagDirty(flag, saved) {
  return Object.entries(flagUpdatePayload(flag)).some(([key, value]) => value !== saved?.[key]);
}
