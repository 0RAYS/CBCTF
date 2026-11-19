import request from './request';

// 获取用户信息
export const getUserInfo = () => {
  return request({
    url: '/me',
    method: 'GET',
  });
};

// 获取当前用户可访问的 API 路由列表
export const getAccessibleRoutes = () => {
  return request({
    url: '/me/permissions',
    method: 'GET',
  });
};

// 更新密码
export const updatePassword = (data) => {
  return request({
    url: '/me/password',
    method: 'PUT',
    data: data,
  });
};

// 更新用户信息
export const updateUserInfo = (data) => {
  return request({
    url: '/me',
    method: 'PUT',
    data: data,
  });
};

// 激活邮箱
export const activateEmail = (data) => {
  return request({
    url: '/me/activate',
    method: 'POST',
    data: data,
  });
};

// 删除账户
export const deleteAccount = (data) => {
  return request({
    url: '/me',
    method: 'DELETE',
    data: data,
  });
};

// 上传用户头像
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
