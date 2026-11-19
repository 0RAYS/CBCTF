import request from '../request';

export const getWebhookList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/webhook',
    method: 'GET',
    params,
  });
};

export const createWebhook = (data) => {
  return request({
    url: '/admin/webhook',
    method: 'POST',
    data,
  });
};

export const updateWebhook = (webhookId, data) => {
  return request({
    url: `/admin/webhook/${webhookId}`,
    method: 'PUT',
    data,
  });
};

export const deleteWebhook = (webhookId) => {
  return request({
    url: `/admin/webhook/${webhookId}`,
    method: 'DELETE',
  });
};

export const getAllWebhookHistory = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/webhook/history`,
    method: 'GET',
    params,
  });
};

export const getWebhookHistory = (webhookID, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/webhook/${webhookID}/history`,
    method: 'GET',
    params,
  });
};

export const getEvents = () => {
  return request({
    url: '/admin/webhook/events',
    method: 'GET',
  });
};
