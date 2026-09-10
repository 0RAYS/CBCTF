export function challengeQuery({
  page = 1,
  pageSize = 10,
  type = 'all',
  category = 'all',
  name = '',
  description = '',
}) {
  return {
    limit: pageSize,
    offset: (page - 1) * pageSize,
    ...(type !== 'all' ? { type } : {}),
    ...(category !== 'all' ? { category } : {}),
    ...(name.trim() ? { name: name.trim() } : {}),
    ...(description.trim() ? { description: description.trim() } : {}),
  };
}

export function challengeDraft(challenge) {
  return {
    name: challenge.name || '',
    description: challenge.description || '',
    attempt: challenge.attempt || 0,
    hidden: !!challenge.hidden,
    hints: [...(challenge.hints || [])],
    tags: [...(challenge.tags || [])],
  };
}

export function toggleChallengeSelection(selected, challenge) {
  return selected.some((item) => item.id === challenge.id)
    ? selected.filter((item) => item.id !== challenge.id)
    : [...selected, challenge];
}

export function flagUpdatePayload(flag) {
  return {
    value: flag.value,
    score_type: flag.score_type,
    score: flag.score,
    decay: flag.decay,
    min_score: flag.min_score,
  };
}

export function isFlagDirty(flag, saved) {
  return Object.entries(flagUpdatePayload(flag)).some(([key, value]) => value !== saved?.[key]);
}
