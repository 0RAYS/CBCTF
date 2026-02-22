import { useState } from 'react';
import UsersTab from './rbac/UsersTab';
import RolesTab from './rbac/RolesTab';
import GroupsTab from './rbac/GroupsTab';
import PermissionsTab from './rbac/PermissionsTab';
import { useTranslation } from 'react-i18next';

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
      <div className="w-full mx-auto mb-6">
        <div className="flex border-b border-neutral-700">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              className={`px-6 py-3 text-sm font-medium transition-colors ${
                activeTab === tab.key
                  ? 'text-blue-400 border-b-2 border-blue-400'
                  : 'text-neutral-400 hover:text-neutral-300'
              }`}
              onClick={() => setActiveTab(tab.key)}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {activeTab === 'users' && <UsersTab />}
      {activeTab === 'roles' && <RolesTab />}
      {activeTab === 'groups' && <GroupsTab />}
      {activeTab === 'permissions' && <PermissionsTab />}
    </>
  );
}

export default RbacManagement;
