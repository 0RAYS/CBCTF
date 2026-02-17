import request from './request';

// 获取竞赛列表
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": {
//        "contests": [
//            {
//                "picture": "",
//                "description": "test123456",
//                "duration": 7200000000000,
//                "hidden": false,
//                "id": 12,
//                "name": "test123456",
//                "prefix": "CBCTF",
//                "size": 1,
//                "start": "2025-02-01T05:45:26.735Z",
//                "teams": 1,
//                "users": 1
//            }
//        ],
//        "count": 1
//    },
//    "trace": "bb5fa197-5aac-47eb-866f-111d38ca2e80"
// }
export const getContestList = () => {
  return request({
    url: '/contests',
    method: 'GET',
  });
};

// 获取比赛详情
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": {
//        "picture": "https://test.jbnrz.com.cn/picture/df20aa8b-b15f-4bc7-811b-16efa185ee44",
//        "description": "contest11",
//        "duration": 25920000000000000,
//        "hidden": false,
//        "id": 11,
//        "name": "contest11",
//        "prefix": "CBCTF",
//        "size": 4,
//        "start": "2025-01-29T14:35:24.572+08:00",
//        "teams": 1,
//        "users": 1
//    },
//    "trace": "3136b1ed-fc3c-4be7-a150-3f4ed3f49db4"
// }
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
