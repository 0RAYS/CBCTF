import request from '../request';

export const getTrafficAnalysis = ({ contestId, teamId, victimId }, signal) => {
  const scope = contestId && teamId ? `/contests/${contestId}/teams/${teamId}` : '';
  return request({
    url: `/admin${scope}/victims/${victimId}/traffic/analysis`,
    method: 'GET',
    signal,
    timeout: 70000,
    noLoading: true,
    noToast: true,
  });
};

export const getTrafficOverlaps = (contestId, signal) =>
  request({
    url: `/admin/contests/${contestId}/traffic/overlaps`,
    method: 'GET',
    signal,
    noLoading: true,
    noToast: true,
  });
