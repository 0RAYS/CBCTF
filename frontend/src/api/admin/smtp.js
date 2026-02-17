import request from '../request';

// 获取SMTP配置列表
export const getSmtpList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/smtp',
    method: 'GET',
    params,
  });
};

// 创建SMTP配置
export const createSmtp = (data) => {
  return request({
    url: '/admin/smtp',
    method: 'POST',
    data,
  });
};

// 更新SMTP配置
export const updateSmtp = (smtpId, data) => {
  return request({
    url: `/admin/smtp/${smtpId}`,
    method: 'PUT',
    data,
  });
};

// 删除SMTP配置
export const deleteSmtp = (smtpId) => {
  return request({
    url: `/admin/smtp/${smtpId}`,
    method: 'DELETE',
  });
};
