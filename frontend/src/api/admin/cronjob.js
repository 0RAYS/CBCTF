import request from '../request';

export const getCronJobList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/cronjobs',
    method: 'GET',
    params,
  });
};

export const updateCronJob = (cronJobId, data) => {
  return request({
    url: `/admin/cronjobs/${cronJobId}`,
    method: 'PUT',
    data,
  });
};
