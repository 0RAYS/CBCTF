import request from '../request';

// 获取系统状态
export const getSystemStatus = (noLoading = false) => {
  return request({
    url: '/admin/system/status',
    method: 'GET',
    noLoading,
  });
};

// 获取系统配置
export const getSystemConfig = () => {
  return request({
    url: '/admin/system/config',
    method: 'GET',
  });
};

// 更新系统配置
export const updateSystemConfig = (data) => {
  return request({
    url: '/admin/system/config',
    method: 'PUT',
    data,
  });
};

export const restartSystem = () => {
  return request({
    url: '/admin/system/restart',
    method: 'POST',
  });
};

// 更新管理员密码
export const updateAdminPassword = (data) => {
  return request({
    url: '/admin/me/password',
    method: 'PUT',
    data,
  });
};

// 更新管理员信息
export const updateAdminInfo = (data) => {
  return request({
    url: '/admin/me',
    method: 'PUT',
    data,
  });
};

// 更新管理员头像
export const updateAdminPicture = (file) => {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: '/admin/me/picture',
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 获取系统运行日志（带颜色的ANSI文本数组）
// 支持分页参数：limit，offset
export const getSystemLogs = (params = { limit: 100, offset: 0 }) => {
  return request({
    url: '/admin/logs',
    method: 'GET',
    params,
  });
};
