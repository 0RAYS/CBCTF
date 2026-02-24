export const DEFAULT_CONTEST_IMAGE =
  'https://images.unsplash.com/photo-1562813733-b31f71025d54?q=80&w=2069&auto=format&fit=crop';

export function getContestStatus(startTime, durationSeconds, nowMs = Date.now()) {
  const start = new Date(startTime).getTime();
  const end = start + durationSeconds * 1000;

  if (nowMs < start) return 'upcoming';
  if (nowMs > end) return 'ended';
  return 'active';
}

export function getContestTimeRange(startTime, durationSeconds) {
  const startDate = new Date(startTime);
  const endDate = new Date(startDate.getTime() + durationSeconds * 1000);
  return { startTime: startDate.toISOString(), endTime: endDate.toISOString() };
}

export function isContestEnded(startTime, durationSeconds, nowMs = Date.now()) {
  const start = new Date(startTime).getTime();
  const end = start + durationSeconds * 1000;
  return nowMs > end;
}

