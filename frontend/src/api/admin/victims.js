import request from '../request';

export const getVictims = (params = {}) => request({ url: '/admin/victims', method: 'GET', params });

export const stopVictims = (victimIds) =>
  request({ url: '/admin/victims', method: 'DELETE', data: { victims: victimIds } });
