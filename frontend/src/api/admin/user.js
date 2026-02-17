import request from '../request';

// 获取用户列表
export const getUserList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/users',
    method: 'GET',
    params,
  });
};

// 获取用户详情
export const getUserInfo = (userId) => {
  return request({
    url: `/admin/users/${userId}`,
    method: 'GET',
  });
};

// 创建用户
export const createUser = (data) => {
  return request({
    url: '/admin/users',
    method: 'POST',
    data,
  });
};

// 更新用户信息
export const updateUser = (userId, data) => {
  return request({
    url: `/admin/users/${userId}`,
    method: 'PUT',
    data,
  });
};

// 删除用户
export const deleteUser = (userId) => {
  return request({
    url: `/admin/users/${userId}`,
    method: 'DELETE',
  });
};

// 更新用户头像
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
