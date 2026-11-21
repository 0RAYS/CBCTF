import request from '@/api/request';

export const getPublicConfig = () => {
  return request({
    url: '/config',
    method: 'GET',
    noLoading: true,
  });
};
