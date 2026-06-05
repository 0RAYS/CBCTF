import request from '../request';

export const getVictims = (params = {}) => request({ url: '/admin/victims', method: 'GET', params });

export const stopVictims = (victimIds) =>
  request({ url: '/admin/victims', method: 'DELETE', data: { victims: victimIds } });

export const getVictimTraffic = (victimId, params = {}) =>
  request({ url: `/admin/victims/${victimId}/traffic`, method: 'GET', params });

export const downloadVictimTraffic = (victimId) =>
  request({ url: `/admin/victims/${victimId}/traffic/download`, method: 'GET', responseType: 'blob' });

export const getVictimPods = (victimId) => request({ url: `/admin/victims/${victimId}/pods`, method: 'GET' });

export const getVictimPodLogs = (victimId, podName, container = '', lines = 1000) =>
  request({
    url: `/admin/victims/${victimId}/pods/logs`,
    method: 'GET',
    params: { pod_name: podName, container, lines },
  });
