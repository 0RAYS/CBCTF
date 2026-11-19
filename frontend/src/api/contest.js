import request from './request';

// 获取竞赛列表
export const getContestList = () => {
  return request({
    url: '/contests',
    method: 'GET',
  });
};

// 获取比赛详情
export const getContestInfo = (contestId) => {
  return request({
    url: `/contests/${contestId}`,
    method: 'GET',
  });
};

// 获取比赛排行榜
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

// 获取比赛积分榜表格视图
export const getContestScoreboard = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/contests/${contestId}/scoreboard`,
    method: 'GET',
    params,
  });
};

// 获取比赛分数时间线
export const getContestTimeline = (contestId) => {
  return request({
    url: `/contests/${contestId}/timeline`,
    method: 'GET',
  });
};

export const getContestNotices = async (contestId, params = {}) => {
  return request.get(`/contests/${contestId}/notices`, { params });
};
