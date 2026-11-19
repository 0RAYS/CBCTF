import request from '../request';

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

export const deleteTeam = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me`,
    method: 'DELETE',
  });
};

export const kickTeamMember = (contestId, userId) => {
  return request({
    url: `/contests/${contestId}/teams/me/kick`,
    method: 'POST',
    data: {
      user_id: userId,
    },
  });
};

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

export const getTeamMembers = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me/users`,
    method: 'GET',
  });
};

export const updateTeamInfo = (contestId, data) => {
  return request({
    url: `/contests/${contestId}/teams/me`,
    method: 'PUT',
    data,
  });
};

export const createTeam = (contestId, data) => {
  return request({
    url: `/contests/${contestId}/teams/create`,
    method: 'POST',
    noToast: true,
    data,
  });
};

export const joinTeam = (contestId, data) => {
  return request({
    url: `/contests/${contestId}/teams/join`,
    method: 'POST',
    noToast: true,
    data,
  });
};
