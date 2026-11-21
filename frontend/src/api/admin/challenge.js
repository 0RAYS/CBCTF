import request from '@/api/request';

export const getChallengeCategories = (params = {}) => {
  return request({
    url: '/admin/challenges/categories',
    method: 'GET',
    params,
  });
};

export const getChallengeList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/challenges',
    method: 'GET',
    params,
  });
};

export const getChallenge = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}`,
    method: 'GET',
  });
};

export const createChallenge = (data) => {
  return request({
    url: '/admin/challenges',
    method: 'POST',
    data,
  });
};

export const updateChallenge = (challengeId, data) => {
  return request({
    url: `/admin/challenges/${challengeId}`,
    method: 'PUT',
    data,
  });
};

export const deleteChallenge = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}`,
    method: 'DELETE',
  });
};

export const copyChallenge = (challengeId, contestId) => {
  return request({
    url: `/admin/challenges/${challengeId}/copy`,
    method: 'POST',
    data: { contest_id: contestId },
  });
};

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

export const downloadChallengeFile = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/download`,
    method: 'GET',
    responseType: 'blob',
  });
};

export const getChallengeFlags = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/flags`,
    method: 'GET',
  });
};

export const updateChallengeFlag = (challengeId, flagId, data) => {
  return request({
    url: `/admin/challenges/${challengeId}/flags/${flagId}`,
    method: 'PUT',
    data,
  });
};

export const getFlagSolvers = (challengeId, flagId) => {
  return request({
    url: `/admin/challenges/${challengeId}/flags/${flagId}/solvers`,
    method: 'GET',
  });
};

export const resetChallengeSolves = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/reset-solves`,
    method: 'POST',
  });
};

export const getTestChallengeStatus = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/test`,
    method: 'GET',
  });
};

export const downloadTestAttachment = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/test/attachment`,
    method: 'GET',
    responseType: 'blob',
  });
};

export const startTestVictim = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/test/start`,
    method: 'POST',
  });
};

export const stopTestVictim = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/test/stop`,
    method: 'POST',
  });
};
