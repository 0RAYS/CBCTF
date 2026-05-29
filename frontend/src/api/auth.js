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

// 发送密码重置邮件
export const forgotPassword = (data) => {
  return request({
    url: '/password/forgot',
    method: 'POST',
    data: {
      email: data.email,
    },
  });
};

// 重置密码（使用邮件链接中的 token）
export const resetPassword = (data) => {
  return request({
    url: '/password/reset',
    method: 'POST',
    data: {
      token: data.token,
      id: data.id,
      password: data.password,
    },
  });
};
