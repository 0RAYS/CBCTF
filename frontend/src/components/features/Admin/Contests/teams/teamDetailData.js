export const TEAM_DETAIL_PAGE_SIZE = 20;

export function flagFilterOptions(challenges) {
  return {
    types: [...new Set(challenges.map((item) => item.type).filter(Boolean))].sort(),
    categories: [...new Set(challenges.map((item) => item.category).filter(Boolean))].sort(),
  };
}

export function filterTeamFlags(challenges, filters) {
  return challenges.filter((challenge) => {
    if (filters.name && !challenge.name?.toLowerCase().includes(filters.name.toLowerCase())) return false;
    if (filters.type && challenge.type !== filters.type) return false;
    if (filters.category && challenge.category !== filters.category) return false;
    const solved = (challenge.flags || []).some((flag) => flag.solved);
    if (filters.solved === 'true' && !solved) return false;
    if (filters.solved === 'false' && solved) return false;
    return true;
  });
}

export function formatFileSize(bytes) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}

export function containerStatus(startTime, duration, now = Date.now()) {
  const start = new Date(startTime).getTime();
  if (now < start) return 'upcoming';
  if (now > start + duration * 1000) return 'ended';
  return 'running';
}
