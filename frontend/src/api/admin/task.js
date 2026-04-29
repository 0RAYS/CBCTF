import request from '../request';

export const getTaskHistory = (params = { limit: 20, offset: 0 }) =>
  request({
    url: '/admin/tasks',
    method: 'GET',
    params,
  });

export const getLiveTasks = (params = { limit: 20, offset: 0 }) =>
  request({
    url: '/admin/tasks/live',
    method: 'GET',
    params,
  });
