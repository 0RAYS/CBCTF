export const isVictimStoppable = (victim) => victim?.status === 'running';

export function buildVictimListParams(filters, page, deleted, pageSize = 20) {
  const params = { ...filters, limit: pageSize, offset: (page - 1) * pageSize };
  delete params.deleted;
  if (deleted) params.deleted = true;
  return Object.fromEntries(Object.entries(params).filter(([, value]) => value !== ''));
}

export function toggleVictimSelection(selected, victims, id) {
  if (!isVictimStoppable(victims.find((victim) => victim.id === id))) return selected;
  return selected.includes(id) ? selected.filter((value) => value !== id) : [...selected, id];
}

export function selectPageVictims(selected, victims) {
  const ids = victims.filter(isVictimStoppable).map((victim) => victim.id);
  return selected.length === ids.length ? [] : ids;
}

export function updateChallengeSelection(selected, id, checked) {
  return checked ? (selected.includes(id) ? selected : [...selected, id]) : selected.filter((value) => value !== id);
}

export function buildVictimStartPayload(challenges, percentage, durationInput) {
  return {
    challenges: [...challenges],
    team_ratio: percentage / 100,
    duration: Number.parseInt(durationInput, 10) || 0,
  };
}

export function estimateVictimTeams(total, percentage) {
  return total > 0 && percentage > 0 && percentage < 100 ? Math.max(1, Math.floor((total * percentage) / 100)) : 0;
}
