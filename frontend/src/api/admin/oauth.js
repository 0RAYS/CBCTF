import request from '../request';

// 获取OAuth Provider列表
export const getOAuthProviderList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/oauth',
    method: 'GET',
    params,
  });
};

// 创建OAuth Provider
export const createOAuthProvider = (data) => {
  return request({
    url: '/admin/oauth',
    method: 'POST',
    data,
  });
};

// 更新OAuth Provider
export const updateOAuthProvider = (providerId, data) => {
  return request({
    url: `/admin/oauth/${providerId}`,
    method: 'PUT',
    data,
  });
};

// 删除OAuth Provider
export const deleteOAuthProvider = (providerId) => {
  return request({
    url: `/admin/oauth/${providerId}`,
    method: 'DELETE',
  });
};

// 上传OAuth Provider Logo
export const uploadOAuthPicture = (providerId, file) => {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: `/admin/oauth/${providerId}/picture`,
    method: 'POST',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' },
  });
};
