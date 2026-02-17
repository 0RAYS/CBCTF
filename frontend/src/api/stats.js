import request from './request';

export const getStats = () => {
  return request({
    url: '/stats',
    method: 'GET',
  });
};
