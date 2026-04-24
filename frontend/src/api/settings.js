import request from './request';

export const getPublicConfig = () => {
  return request({
    url: '/config',
    method: 'GET',
    noLoading: true,
  });
};
