import request from './request';

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

export const getContestRank = (contestId, limit = 3, offset = 0, showHidden = false) => {
  return request({
    url: `/contests/${contestId}/rank`,
    method: 'GET',
    params: {
      limit,
      offset,
      show_hidden: showHidden ? 1 : 0,
    },
  });
};

export const getContestScoreboard = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/contests/${contestId}/scoreboard`,
    method: 'GET',
    params,
  });
};

export const getContestTimeline = (contestId, showHidden = false) => {
  return request({
    url: `/contests/${contestId}/timeline`,
    method: 'GET',
    params: {
      show_hidden: showHidden ? 1 : 0,
    },
  });
};

export const getContestNotices = async (contestId, params = {}) => {
  return request.get(`/contests/${contestId}/notices`, { params });
};
