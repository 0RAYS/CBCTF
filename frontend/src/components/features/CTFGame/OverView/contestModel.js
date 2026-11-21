import { DEFAULT_CONTEST_IMAGE, getContestStatus, getContestTimeRange } from '../../../../config/contest.js';

export function transformContestData(apiData) {
  if (!apiData) return null;
  const { startTime, endTime } = getContestTimeRange(apiData.start, apiData.duration);

  return {
    title: `${apiData.prefix} ${apiData.name}`,
    description: apiData.description,
    image: apiData.picture || DEFAULT_CONTEST_IMAGE,
    status: apiData.paused ? 'paused' : getContestStatus(apiData.start, apiData.duration),
    startTime,
    // The server reports the effective end, which moves while the contest is paused.
    endTime: apiData.end ? new Date(apiData.end).toISOString() : endTime,
    participants: apiData.users ?? null,
    rules: apiData.rules ?? [],
    prizes: apiData.prizes ?? [],
    timeline: apiData.timelines ?? [],
    teamSize: apiData.size,
    teamsCount: apiData.teams,
    noticesCount: apiData.notices,
    isBlood: apiData.blood,
    isHidden: apiData.hidden,
  };
}
