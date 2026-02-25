import { useState } from 'react';
import UsersTab from './rbac/UsersTab';
import RolesTab from './rbac/RolesTab';
import GroupsTab from './rbac/GroupsTab';
import PermissionsTab from './rbac/PermissionsTab';
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

      {activeTab === 'users' && <UsersTab />}
      {activeTab === 'roles' && <RolesTab />}
      {activeTab === 'groups' && <GroupsTab />}
      {activeTab === 'permissions' && <PermissionsTab />}
    </>
  );
}

export default RbacManagement;
