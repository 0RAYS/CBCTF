import request from '../request';

export const getUserList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/users',
    method: 'GET',
    params,
  });
};

export const getUserInfo = (userId) => {
  return request({
    url: `/admin/users/${userId}`,
    method: 'GET',
  });
};

export const createUser = (data) => {
  return request({
    url: '/admin/users',
    method: 'POST',
    data,
  });
};

export const updateUser = (userId, data) => {
  return request({
    url: `/admin/users/${userId}`,
    method: 'PUT',
    data,
  });
};

export const deleteUser = (userId) => {
  return request({
    url: `/admin/users/${userId}`,
    method: 'DELETE',
  });
};

export const updateUserPicture = (userId, file) => {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: `/admin/users/${userId}/picture`,
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};
