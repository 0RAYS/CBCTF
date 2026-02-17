import request from '../request';

// 获取可搜索模型列表
export const getSearchModels = () => {
  return request({
    url: '/admin/models',
    method: 'GET',
  });
};

// 全局搜索
export const searchModels = (params) => {
  return request({
    url: '/admin/search',
    method: 'GET',
    params,
    noLoading: true,
  });
};
