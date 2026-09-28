import { useState } from 'react';
import UsersManager from '../../components/features/Admin/rbac/UsersManager';
import RolesManager from '../../components/features/Admin/rbac/RolesManager';
import GroupsManager from '../../components/features/Admin/rbac/GroupsManager';
import PermissionsManager from '../../components/features/Admin/rbac/PermissionsManager';
import { useTranslation } from 'react-i18next';
import { Tabs } from '../../components/common';

function RbacManagement() {
  const [activeTab, setActiveTab] = useState('users');
  const { t } = useTranslation();

  const tabs = [
    { key: 'users', label: t('admin.rbac.tabs.users') },
    { key: 'roles', label: t('admin.rbac.tabs.roles') },
    { key: 'groups', label: t('admin.rbac.tabs.groups') },
    { key: 'permissions', label: t('admin.rbac.tabs.permissions') },
  ];

  return (
    <>
      <Tabs items={tabs} value={activeTab} onChange={setActiveTab} />

      {activeTab === 'users' && <UsersManager />}
      {activeTab === 'roles' && <RolesManager />}
      {activeTab === 'groups' && <GroupsManager />}
      {activeTab === 'permissions' && <PermissionsManager />}
    </>
  );
}

export default RbacManagement;
