import request from '@/api/request';

export const getSmtpList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/smtp',
    method: 'GET',
    params,
  });
};

export const createSmtp = (data) => {
  return request({
    url: '/admin/smtp',
    method: 'POST',
    data,
  });
};

export const updateSmtp = (smtpId, data) => {
  return request({
    url: `/admin/smtp/${smtpId}`,
    method: 'PUT',
    data,
  });
};

export const deleteSmtp = (smtpId) => {
  return request({
    url: `/admin/smtp/${smtpId}`,
    method: 'DELETE',
  });
};

export const testSmtp = (smtpId, data) => {
  return request({
    url: `/admin/smtp/${smtpId}/test`,
    method: 'POST',
    data,
  });
};
