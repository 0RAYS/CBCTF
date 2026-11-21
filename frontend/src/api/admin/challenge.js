import request from '@/api/request';

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

export const getChallengeUsages = (challengeId) => {
  return request({
    url: `/admin/challenges/${challengeId}/usages`,
    method: 'GET',
    noToast: true,
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

export const getChallengeFlags = (contestID, challengeID) => {
  return request({
    url: `/admin/contests/${contestID}/challenges/${challengeID}/flags`,
    method: 'GET',
  });
};

export const updateChallengeFlag = (contestID, challengeID, flagID, data) => {
  return request({
    url: `/admin/contests/${contestID}/challenges/${challengeID}/flags/${flagID}`,
    method: 'PUT',
    data,
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

export const getFlagSolvers = (contestId, challengeId, flagId) => {
  return request({
    url: `/admin/contests/${contestId}/challenges/${challengeId}/flags/${flagId}/solvers`,
    method: 'GET',
  });
};
