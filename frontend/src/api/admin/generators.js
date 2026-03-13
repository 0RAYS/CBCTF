import request from '../request';

export const getGenerators = (params = {}) => request({ url: '/admin/generators', method: 'GET', params });

export const startGenerators = (challenges) =>
  request({ url: '/admin/generators', method: 'POST', data: { challenges } });

export const stopGenerators = (generatorIds) =>
  request({ url: '/admin/generators', method: 'DELETE', data: { generators: generatorIds } });
