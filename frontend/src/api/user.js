import request from '@/api/request';

export const getUserInfo = () => {
  return request({
    url: '/me',
    method: 'GET',
  });
};

export const getAccessibleRoutes = () => {
  return request({
    url: '/me/permissions',
    method: 'GET',
  });
};

export const getUserLeaderboard = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/users/leaderboard',
    method: 'GET',
    params,
  });
};

export const updatePassword = (data) => {
  return request({
    url: '/me/password',
    method: 'PUT',
    data: data,
  });
};

export const updateUserInfo = (data) => {
  return request({
    url: '/me',
    method: 'PUT',
    data: data,
  });
};

export const activateEmail = (data) => {
  return request({
    url: '/me/activate',
    method: 'POST',
    data: data,
  });
};

export const deleteAccount = (data) => {
  return request({
    url: '/me',
    method: 'DELETE',
    data: data,
  });
};

export const uploadPicture = (file) => {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: '/me/picture',
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};
