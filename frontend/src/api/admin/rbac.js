import request from '../request';

// ========== Permissions ==========

export const getPermissionList = (params = { limit: 50, offset: 0 }) => {
  return request({
    url: '/admin/permissions',
    method: 'GET',
    params,
  });
};

export const updatePermission = (permissionId, data) => {
  return request({
    url: `/admin/permissions/${permissionId}`,
    method: 'PUT',
    data,
  });
};

// ========== Roles ==========

export const getRoleList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/roles',
    method: 'GET',
    params,
  });
};

export const createRole = (data) => {
  return request({
    url: '/admin/roles',
    method: 'POST',
    data,
  });
};

export const updateRole = (roleId, data) => {
  return request({
    url: `/admin/roles/${roleId}`,
    method: 'PUT',
    data,
  });
};

export const deleteRole = (roleId) => {
  return request({
    url: `/admin/roles/${roleId}`,
    method: 'DELETE',
  });
};

export const assignPermissionToRole = (roleId, data) => {
  return request({
    url: `/admin/roles/${roleId}/permissions`,
    method: 'POST',
    data,
  });
};

export const revokePermissionFromRole = (roleId, data) => {
  return request({
    url: `/admin/roles/${roleId}/permissions`,
    method: 'DELETE',
    data,
  });
};

export const getRolePermissions = (roleId) => {
  return request({
    url: `/admin/roles/${roleId}/permissions`,
    method: 'GET',
  });
};

// ========== Groups ==========

export const getGroupList = (params = { limit: 20, offset: 0 }) => {
  return request({
    url: '/admin/groups',
    method: 'GET',
    params,
  });
};

export const createGroup = (data) => {
  return request({
    url: '/admin/groups',
    method: 'POST',
    data,
  });
};

export const updateGroup = (groupId, data) => {
  return request({
    url: `/admin/groups/${groupId}`,
    method: 'PUT',
    data,
  });
};

export const deleteGroup = (groupId) => {
  return request({
    url: `/admin/groups/${groupId}`,
    method: 'DELETE',
  });
};

export const assignUserToGroup = (groupId, data) => {
  return request({
    url: `/admin/groups/${groupId}/users`,
    method: 'POST',
    data,
  });
};

export const removeUserFromGroup = (groupId, data) => {
  return request({
    url: `/admin/groups/${groupId}/users`,
    method: 'DELETE',
    data,
  });
};

export const getGroupUsers = (groupId, params = { limit: 10, offset: 0 }) => {
  return request({
    url: `/admin/groups/${groupId}/users`,
    method: 'GET',
    params,
  });
};
