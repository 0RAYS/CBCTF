import request from '@/api/request';

export const getOAuthProviderList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/oauth',
    method: 'GET',
    params,
  });
};

export const createOAuthProvider = (data) => {
  const { protocol, ...body } = data;
  return request({
    url: '/admin/oauth',
    method: 'POST',
    params: protocol ? { protocol } : undefined,
    data: body,
  });
};

export const updateOAuthProvider = (providerId, data) => {
  const { protocol, ...body } = data;
  return request({
    url: `/admin/oauth/${providerId}`,
    method: 'PUT',
    params: protocol ? { protocol } : undefined,
    data: body,
  });
};

export const deleteOAuthProvider = (providerId) => {
  return request({
    url: `/admin/oauth/${providerId}`,
    method: 'DELETE',
  });
};

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
