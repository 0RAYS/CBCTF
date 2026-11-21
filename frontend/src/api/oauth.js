import request from '@/api/request';

export const getOAuthProviders = () => {
  return request({
    url: '/oauth',
    method: 'GET',
  });
};

export const exchangeOauthCode = (code) => {
  return request({
    url: '/oauth/token',
    method: 'GET',
    params: { code },
  });
};
