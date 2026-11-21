import request from '@/api/request';

export const getSystemStatus = (noLoading = false) => {
  return request({
    url: '/admin/system/status',
    method: 'GET',
    noLoading,
  });
};

export const getSystemConfig = () => {
  return request({
    url: '/admin/system/config',
    method: 'GET',
  });
};

export const updateSystemConfig = (data) => {
  return request({
    url: '/admin/system/config',
    method: 'PUT',
    data,
  });
};

export const uploadGeoCityDB = (file) => {
  const formData = new FormData();
  formData.append('file', file);
  return request({
    url: '/admin/system/geocity-db',
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

export const restartSystem = () => {
  return request({
    url: '/admin/system/restart',
    method: 'POST',
  });
};

// 返回带 ANSI 颜色的日志文本数组
export const getSystemLogs = (params = { limit: 100, offset: 0 }) => {
  return request({
    url: '/admin/logs',
    method: 'GET',
    params,
  });
};
