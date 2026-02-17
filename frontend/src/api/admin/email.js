import request from '../request';

// 获取邮件发送历史记录列表
export const getEmailHistory = (smtpID, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/smtp/${smtpID}/email`,
    method: 'GET',
    params,
  });
};

export const getAllEmailHistory = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/email`,
    method: 'GET',
    params,
  });
};
