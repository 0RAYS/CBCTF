import request from '../request';

export const getAdminPullImages = () => {
  return request({
    url: '/admin/images',
    method: 'GET',
  });
};

export const pullAdminImages = (data) => {
  return request({
    url: '/admin/images',
    method: 'POST',
    noToast: true,
    data,
  });
};

// 手动触发镜像检测
export const detectAdminImages = () => {
  return request({
    url: '/admin/images/detect',
    method: 'POST',
    noToast: true,
  });
};
