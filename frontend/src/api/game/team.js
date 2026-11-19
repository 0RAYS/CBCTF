import request from '../request';

// 上传队伍头像
// 请求参数: FormData, 包含picture字段
export const uploadTeamPicture = (contestId, file) => {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: `/contests/${contestId}/teams/me/picture`,
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};
// 解散队伍
export const deleteTeam = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me`,
    method: 'DELETE',
  });
};

// 踢出队员
export const kickTeamMember = (contestId, userId) => {
  return request({
    url: `/contests/${contestId}/teams/me/kick`,
    method: 'POST',
    data: {
      user_id: userId,
    },
  });
};

// 获取比赛队伍信息
export const getTeamInfo = (contestId, { noToast = false } = {}) => {
  return request({
    url: `/contests/${contestId}/teams/me`,
    method: 'GET',
    noToast,
  });
};

export const getTeamCaptcha = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me/captcha`,
    method: 'GET',
  });
};

export const updateTeamCaptcha = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me/captcha`,
    method: 'PUT',
  });
};

// 获取队伍成员列表
export const getTeamMembers = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me/users`,
    method: 'GET',
  });
};

// 更新队伍信息
export const updateTeamInfo = (contestId, data) => {
  return request({
    url: `/contests/${contestId}/teams/me`,
    method: 'PUT',
    data,
  });
};

// 创建队伍
export const createTeam = (contestId, data) => {
  return request({
    url: `/contests/${contestId}/teams/create`,
    method: 'POST',
    noToast: true,
    data,
  });
};

// 加入队伍
export const joinTeam = (contestId, data) => {
  return request({
    url: `/contests/${contestId}/teams/join`,
    method: 'POST',
    noToast: true,
    data,
  });
};
