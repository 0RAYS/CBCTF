import request from '@/api/request';

export const getCaptcha = () => {
  return request({
    url: '/captcha',
    method: 'GET',
  });
};

export const login = (data) => {
  return request({
    url: '/login',
    method: 'POST',
    data: {
      name: data.name,
      password: data.password,
      captcha_id: data.captchaId,
      captcha: data.captcha,
    },
  });
};

export const register = (data) => {
  return request({
    url: '/register',
    method: 'POST',
    data: {
      name: data.name,
      email: data.email,
      password: data.password,
      captcha_id: data.captchaId,
      captcha: data.captcha,
    },
  });
};

export const verifyEmail = (data) => {
  return request({
    url: '/verify',
    method: 'POST',
    data: {
      token: data.token,
    },
  });
};

export const forgotPassword = (data) => {
  return request({
    url: '/password/forgot',
    method: 'POST',
    data: {
      email: data.email,
      captcha_id: data.captchaId,
      captcha: data.captcha,
    },
  });
};

export const resetPassword = (data) => {
  return request({
    url: '/password/reset',
    method: 'POST',
    data: {
      token: data.token,
      password: data.password,
    },
  });
};
