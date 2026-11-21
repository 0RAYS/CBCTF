import request from '@/api/request';

export const getBranding = () => {
  return request({
    url: '/branding',
    method: 'GET',
    noLoading: true,
  });
};
