import { DEFAULT_CONTEST_IMAGE, getContestStatus, getContestTimeRange } from '../../../../config/contest.js';

export function transformContestData(apiData) {
  if (!apiData) return null;
  const { startTime, endTime } = getContestTimeRange(apiData.start, apiData.duration);

  return {
    title: `${apiData.prefix} ${apiData.name}`,
    description: apiData.description,
    image: apiData.picture || DEFAULT_CONTEST_IMAGE,
    status: getContestStatus(apiData.start, apiData.duration),
    startTime,
    endTime,
    participants: apiData.users || 0,
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
