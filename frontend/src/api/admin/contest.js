import request from '../request';

// 获取比赛列表
export function getContestList({ limit = 6, offset = 0 }) {
  return request({
    url: '/admin/contests',
    method: 'GET',
    params: {
      limit,
      offset,
    },
  });
}

// 创建比赛
export function createContest(data) {
  return request({
    url: '/admin/contests',
    method: 'POST',
    data,
  });
}

// 删除比赛
export function deleteContest(contestId) {
  return request({
    url: `/admin/contests/${contestId}`,
    method: 'DELETE',
  });
}

// 上传比赛封面
export function updateContestPicture(contestId, file) {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: `/admin/contests/${contestId}/picture`,
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}

// 获取比赛详情
export function getContestInfo(contestId) {
  return request({
    url: `/admin/contests/${contestId}`,
    method: 'GET',
  });
}

// 获取比赛题目列表
export function getContestChallenges(contestId, params = { limit: 20, offset: 0 }) {
  return request({
    url: `/admin/contests/${contestId}/challenges`,
    method: 'GET',
    params,
  });
}

// 添加题目到比赛
export function addContestChallenge(contestId, challenge_ids) {
  return request({
    url: `/admin/contests/${contestId}/challenges`,
    method: 'POST',
    data: {
      challenge_ids: challenge_ids,
    },
  });
}

// 更新比赛题目信息
export function updateContestChallenge(contestId, challengeId, data) {
  return request({
    url: `/admin/contests/${contestId}/challenges/${challengeId}`,
    method: 'PUT',
    data,
  });
}

// 移出比赛题目
export function removeContestChallenge(contestId, challengeId) {
  return request({
    url: `/admin/contests/${contestId}/challenges/${challengeId}`,
    method: 'DELETE',
  });
}

// 更新比赛信息
export function updateContestInfo(contestId, data) {
  return request({
    url: `/admin/contests/${contestId}`,
    method: 'PUT',
    data,
  });
}

// 获取比赛团队列表
export function getContestTeams(contestId, params = { limit: 20, offset: 0 }) {
  return request({
    url: `/admin/contests/${contestId}/teams`,
    method: 'GET',
    params,
  });
}

// 获取比赛单个团队信息
export function getContestTeam(contestId, teamId) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}`,
    method: 'GET',
  });
}

// 获取比赛排行榜
export function getContestRank(contestId, params = { limit: 20, offset: 0 }, noLoading = false) {
  return request({
    url: `/admin/contests/${contestId}/rank`,
    method: 'GET',
    params,
    noLoading,
  });
}

// 更新队伍头像
export function updateTeamPicture(contestId, teamId, file) {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/picture`,
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}

// 更新队伍信息
export function updateTeamInfo(contestId, teamId, data) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}`,
    method: 'PUT',
    data,
  });
}

// 删除队伍
export function deleteTeam(contestId, teamId) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}`,
    method: 'DELETE',
  });
}

// 移出队伍成员
export function kickTeamMember(contestId, teamId, userId) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/kick`,
    method: 'POST',
    data: {
      user_id: userId,
    },
  });
}

export function getTeamMembers(contestID, teamID) {
  return request({
    url: `/admin/contests/${contestID}/teams/${teamID}/users`,
    method: 'GET',
  });
}

// 获取团队容器列表
export const getTeamContainers = async (contestId, teamId, params) => {
  return request.get(`/admin/contests/${contestId}/teams/${teamId}/victims`, {
    params,
  });
};

// 获取容器流量数据
export function getContainerTraffic(contestId, teamId, containerId, params = { limit: 20, offset: 0 }) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/victims/${containerId}/traffic`,
    method: 'GET',
    params,
  });
}

// 下载容器流量文件
export function downloadContainerTraffic(contestId, teamId, containerId) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/victims/${containerId}/traffic/download`,
    method: 'GET',
    responseType: 'blob',
  });
}

// 获取公告列表
export const getContestNotices = (contestId, params = {}) => {
  return request({
    url: `/admin/contests/${contestId}/notices`,
    method: 'GET',
    params,
  });
};

// 创建公告
export const createContestNotice = (contestId, data) => {
  return request({
    url: `/admin/contests/${contestId}/notices`,
    method: 'POST',
    data,
  });
};

// 更新公告
export const updateContestNotice = (contestId, noticeId, data) => {
  return request({
    url: `/admin/contests/${contestId}/notices/${noticeId}`,
    method: 'PUT',
    data,
  });
};

// 删除公告
export const deleteContestNotice = (contestId, noticeId) => {
  return request({
    url: `/admin/contests/${contestId}/notices/${noticeId}`,
    method: 'DELETE',
  });
};

// 获取比赛团队提交列表
export const getContestTeamSubmissions = (contestId, teamId, params) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/submissions`,
    method: 'GET',
    params,
  });
};

// 获取比赛团队题解列表
export const getContestTeamWriteups = (contestId, teamId, params) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/writeups`,
    method: 'GET',
    params,
  });
};

// 下载比赛团队题解
export const downloadContestTeamWriteup = (contestId, teamId, writeupId) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/writeups/${writeupId}`,
    method: 'GET',
    responseType: 'blob',
  }).then((response) => {
    return response;
  });
};

// 获取比赛团队流量列表
export const getContestTeamTraffic = (contestId, teamId, victimID, params) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/victims/${victimID}/traffic`,
    method: 'GET',
    params,
  });
};

// 获取预热镜像状态
export const getContestWarmupImages = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/images`,
    method: 'GET',
  });
};

// 执行预热镜像
export const warmupContestImages = (contestId, data) => {
  return request({
    url: `/admin/contests/${contestId}/images`,
    method: 'POST',
    data,
  });
};

// 获取积分榜表格视图
export const getContestScoreboard = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/contests/${contestId}/scoreboard`,
    method: 'GET',
    params,
  });
};

// 获取比赛分数时间线
export const getContestTimeline = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/timeline`,
    method: 'GET',
  });
};

// 获取容器列表
export const getContestVictims = (contestId, params = {}) => {
  return request({
    url: `/admin/contests/${contestId}/victims`,
    method: 'GET',
    params,
  });
};

// 停止容器
export const stopContestVictims = (contestId, victimIds) => {
  return request({
    url: `/admin/contests/${contestId}/victims`,
    method: 'DELETE',
    data: {
      victims: victimIds,
    },
  });
};

// 开启容器
export const startContestVictims = (contestId, challenges, teams) => {
  return request({
    url: `/admin/contests/${contestId}/victims`,
    method: 'POST',
    data: {
      challenges,
      teams,
    },
  });
};

// 获取比赛作弊事件列表
export const getContestCheats = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/contests/${contestId}/cheats`,
    method: 'GET',
    params,
  });
};

// 获取比赛团队Flag列表
export const getContestTeamFlags = (contestId, teamId) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/flags`,
    method: 'GET',
  });
};

// 更新作弊事件信息
export const updateContestCheat = (contestId, cheatId, data) => {
  return request({
    url: `/admin/contests/${contestId}/cheats/${cheatId}`,
    method: 'PUT',
    data,
  });
};

// 删除单个作弊事件
export const deleteContestCheat = (contestId, cheatId) => {
  return request({
    url: `/admin/contests/${contestId}/cheats/${cheatId}`,
    method: 'DELETE',
  });
};

// 删除比赛所有作弊事件
export const deleteAllContestCheats = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/cheats`,
    method: 'DELETE',
  });
};

// 手动触发作弊检测
export const checkContestCheats = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/cheats`,
    method: 'POST',
  });
};

// IP 地理位置反查
export const getIpInfo = (ip) => {
  return request({
    url: '/admin/ip',
    method: 'GET',
    params: { ip },
  });
};
