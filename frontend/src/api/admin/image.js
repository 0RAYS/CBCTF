import request from '@/api/request';

export const getAdminPullImages = () => {
  return request({
    url: '/admin/images',
    method: 'GET',
  });
};

export const pullAdminImages = (data) => {
  return request({
    url: '/admin/images',
    method: 'POST',
    noToast: true,
    data,
  });
};

export const detectAdminImages = () => {
  return request({
    url: '/admin/images/detect',
    method: 'POST',
    noToast: true,
  });
};
