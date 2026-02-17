import request from '../request';

// 获取Webhook配置列表
export const getWebhookList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/webhook',
    method: 'GET',
    params,
  });
};

// 创建Webhook配置
export const createWebhook = (data) => {
  return request({
    url: '/admin/webhook',
    method: 'POST',
    data,
  });
};

// 更新Webhook配置
export const updateWebhook = (webhookId, data) => {
  return request({
    url: `/admin/webhook/${webhookId}`,
    method: 'PUT',
    data,
  });
};

// 删除Webhook配置
export const deleteWebhook = (webhookId) => {
  return request({
    url: `/admin/webhook/${webhookId}`,
    method: 'DELETE',
  });
};

// 获取Webhook历史记录
export const getAllWebhookHistory = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/webhook/history`,
    method: 'GET',
    params,
  });
};

// 获取Webhook历史记录
export const getWebhookHistory = (webhookID, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/webhook/${webhookID}/history`,
    method: 'GET',
    params,
  });
};

// 获取可用的事件列表
export const getEvents = () => {
  return request({
    url: '/admin/webhook/events',
    method: 'GET',
  });
};
