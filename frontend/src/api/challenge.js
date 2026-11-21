import request from '@/api/request';

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

export const getChallengeStatus = (contestId, challengeId, { noToast = false, noLoading = false } = {}) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}`,
    method: 'GET',
    noToast,
    noLoading,
  });
};

export const initChallenge = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/init`,
    method: 'POST',
  });
};

export const downloadChallengeAttachment = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/attachment`,
    method: 'GET',
    responseType: 'blob',
  }).then((response) => {
    return response;
  });
};

export const startRemoteTarget = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/start`,
    method: 'POST',
  });
};

export const extendContainerTime = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/extend`,
    method: 'POST',
  });
};

export const stopContainer = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/stop`,
    method: 'POST',
  });
};

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

export const getWriteups = (contestId) => {
  return request({
    url: `/contests/${contestId}/writeups`,
    method: 'GET',
  });
};
