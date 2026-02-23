import request from './request';

// 用户登录
export const login = (data) => {
  return request({
    url: '/login',
    method: 'POST',
    data: {
      name: data.name,
      password: data.password,
    },
  });
};

// 用户注册
export const register = (data) => {
  return request({
    url: '/register',
    method: 'POST',
    data: {
      name: data.name,
      email: data.email,
      password: data.password,
    },
  });
};
