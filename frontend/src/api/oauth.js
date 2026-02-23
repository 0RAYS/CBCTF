import request from './request';

// 获取可用的OAuth提供商列表
export const getOAuthProviders = () => {
  return request({
    url: '/oauth',
    method: 'GET',
  });
};

// 用一次性 code 换取真实 token
export const exchangeOauthCode = (code) => {
  return request({
    url: '/oauth/token',
    method: 'GET',
    params: { code },
  });
};
