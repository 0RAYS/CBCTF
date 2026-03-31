import request from '../request';

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
    data,
  });
};
