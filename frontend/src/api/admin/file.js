import request from '../request';

// 获取文件列表
export function getFileList({ limit = 12, offset = 0, type }) {
  return request({
    url: '/admin/files',
    method: 'GET',
    params: {
      limit,
      offset,
      type,
    },
  });
}

// 批量删除文件
export function batchDeleteFiles({ file_ids = [] }) {
  return request({
    url: '/admin/files',
    method: 'DELETE',
    data: {
      file_ids,
    },
  });
}

// 获取文件URL
export function getFileUrl(fileId, type = 'file') {
  const baseUrl = request.defaults.baseURL;
  switch (type) {
    case 'picture':
      return `${baseUrl}/pictures/${fileId}`;
    default:
      return `${baseUrl}/admin/files/${fileId}`;
  }
}

export function downloadFile(fileId) {
  return request({
    url: `/admin/files/${fileId}`,
    method: 'GET',
    responseType: 'blob',
  });
}
