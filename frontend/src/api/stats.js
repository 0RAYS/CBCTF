import request from '@/api/request';

export const getStats = () => {
  return request({
    url: '/stats',
    method: 'GET',
    noLoading: true,
  });
};
