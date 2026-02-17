import request from '../request';

// 获取题目分类
export const getChallengeCategories = (type) => {
  return request({
    url: '/admin/challenges/categories',
    method: 'GET',
    params: { type },
  });
};

export const getNotInContestChallengeList = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/contests/${contestId}/challenges/others`,
    method: 'GET',
    params,
  });
};

export const getContestChallengeCategories = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/challenges/categories`,
    method: 'GET',
  });
};

// 获取题目列表
export const getChallengeList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/challenges',
    method: 'GET',
    params,
  });
};

// 创建题目
export const createChallenge = (data) => {
  return request({
    url: '/admin/challenges',
    method: 'POST',
    data,
  });
};

// 更新题目
export const updateChallenge = (challengeId, data) => {
  return request({
    url: `/admin/challenges/${challengeId}`,
    method: 'PUT',
    data,
  });
};

// 删除题目
export const deleteChallenge = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}`,
    method: 'DELETE',
  });
};

// 上传题目附件
export const uploadChallengeFile = (challengeId, file) => {
  const formData = new FormData();
  formData.append('file', file);
  return request({
    url: `/admin/challenges/${challengeId}/upload`,
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 下载题目文件
export const downloadChallengeFile = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/download`,
    method: 'GET',
    responseType: 'blob',
  });
};

// 获取题目flag列表
export const getChallengeFlags = (contestID, challengeID) => {
  return request({
    url: `/admin/contests/${contestID}/challenges/${challengeID}/flags`,
    method: 'GET',
  });
};

// 更新题目flag
export const updateChallengeFlag = (contestID, challengeID, flagID, data) => {
  return request({
    url: `/admin/contests/${contestID}/challenges/${challengeID}/flags/${flagID}`,
    method: 'PUT',
    data,
  });
};

// 获取题目测试状态
export const getTestChallengeStatus = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/test`,
    method: 'GET',
  });
};

// 下载测试附件
export const downloadTestAttachment = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/test/attachment`,
    method: 'GET',
    responseType: 'blob',
  });
};

// 启动测试靶机
export const startTestVictim = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/test/start`,
    method: 'POST',
  });
};

// 停止测试靶机
export const stopTestVictim = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/test/stop`,
    method: 'POST',
  });
};
