import request from '../request';

export const getAdminBranding = () => {
  return request({
    url: '/admin/branding',
    method: 'GET',
  });
};

export const updateAdminBranding = (data) => {
  return request({
    url: '/admin/branding',
    method: 'PUT',
    data,
  });
};

export const uploadBrandingLogo = (file) => {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: '/admin/branding/logo',
    method: 'POST',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' },
  });
};
