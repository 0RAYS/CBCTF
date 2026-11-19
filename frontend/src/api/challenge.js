import request from './request';

// 获取题目列表
export const getChallengeList = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/contests/${contestId}/challenges`,
    method: 'GET',
    params,
  });
};

export const getChallengeCategories = (contestId) => {
  return request({
    url: `/contests/${contestId}/challenges/categories`,
    method: 'GET',
  });
};

// 获取题目状态
export const getChallengeStatus = (contestId, challengeId, { noToast = false, noLoading = false } = {}) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}`,
    method: 'GET',
    noToast,
    noLoading,
  });
};

// 初始化题目
export const initChallenge = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/init`,
    method: 'POST',
  });
};

// 下载题目附件
// 响应类型: blob
export const downloadChallengeAttachment = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/attachment`,
    method: 'GET',
    responseType: 'blob',
  }).then((response) => {
    return response;
  });
};

// 启动靶机
export const startRemoteTarget = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/start`,
    method: 'POST',
  });
};

// 延长靶机时间
export const extendContainerTime = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/extend`,
    method: 'POST',
  });
};

// 关闭靶机
export const stopContainer = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/stop`,
    method: 'POST',
  });
};

// 提交 flag
export const submitFlag = (contestId, challengeId, data) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/submit`,
    method: 'POST',
    noToast: true,
    data,
  });
};

export const resetChallenge = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/reset`,
    method: 'POST',
  });
};

// 上传题解
// 请求参数为FormData, 包含writeup文件
export const uploadWriteup = (contestId, file) => {
  const formData = new FormData();
  formData.append('writeup', file);

  return request({
    url: `/contests/${contestId}/writeups`,
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 获取已上传题解列表
export const getWriteups = (contestId) => {
  return request({
    url: `/contests/${contestId}/writeups`,
    method: 'GET',
  });
};
