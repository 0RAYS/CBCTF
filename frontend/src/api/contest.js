import request from '@/api/request';

export const getContestList = () => {
  return request({
    url: '/contests',
    method: 'GET',
  });
};

export const getContestInfo = (contestId) => {
  return request({
    url: `/contests/${contestId}`,
    method: 'GET',
  });
};

export const getContestLeaderboard = (contestId, limit = 3, offset = 0, showHidden = false) => {
  return request({
    url: `/contests/${contestId}/leaderboard`,
    method: 'GET',
    params: {
      limit,
      offset,
      show_hidden: showHidden ? 1 : 0,
    },
  });
};

export const getContestLeaderboardChallenges = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/contests/${contestId}/leaderboard/challenges`,
    method: 'GET',
    params,
  });
};

export const getContestLeaderboardTimeline = (contestId, showHidden = false) => {
  return request({
    url: `/contests/${contestId}/leaderboard/timeline`,
    method: 'GET',
    params: {
      show_hidden: showHidden ? 1 : 0,
    },
  });
};

export const getContestNotices = async (contestId, params = {}) => {
  return request.get(`/contests/${contestId}/notices`, { params });
};
