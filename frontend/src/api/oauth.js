import request from './request';

// 获取可用的OAuth提供商列表
export const getOAuthProviders = () => {
  return request({
    url: '/oauth',
    method: 'GET',
  });
};
