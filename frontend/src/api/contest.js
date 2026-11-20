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

export const getContestRank = (contestId, limit = 3, offset = 0) => {
  return request({
    url: `/contests/${contestId}/rank`,
    method: 'GET',
    params: {
      limit,
      offset,
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

export const exportContestScoreboard = (contestId) => {
  return request({
    url: `/contests/${contestId}/scoreboard/export`,
    method: 'GET',
    responseType: 'blob',
  });
};

export const getContestTimeline = (contestId) => {
  return request({
    url: `/contests/${contestId}/timeline`,
    method: 'GET',
  });
};

export const getContestNotices = async (contestId, params = {}) => {
  return request.get(`/contests/${contestId}/notices`, { params });
};
