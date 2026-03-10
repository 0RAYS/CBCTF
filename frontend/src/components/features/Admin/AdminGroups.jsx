import { List, StatusTag } from '../../common';
import { IconPlus, IconEdit, IconTrash, IconUsers } from '@tabler/icons-react';
import { Button, Pagination } from '../../../components/common';
import { useTranslation } from 'react-i18next';

function AdminGroups({
  groups = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  loading = false,
  onPageChange,
  onCreateGroup,
  onEditGroup,
  onDeleteGroup,
  onManageUsers,
  rolesMap = {},
}) {
  const { t } = useTranslation();

  const columns = [
    { key: 'id', label: t('admin.rbac.groups.columns.id'), width: '5%' },
    { key: 'name', label: t('admin.rbac.groups.columns.name'), width: '10%' },
    { key: 'description', label: t('admin.rbac.groups.columns.description'), width: '10%' },
    { key: 'role', label: t('admin.rbac.groups.columns.role'), width: '10%' },
    { key: 'users', label: t('admin.rbac.groups.columns.users'), width: '7%' },
    { key: 'default', label: t('admin.rbac.groups.columns.default'), width: '10%' },
    { key: 'actions', label: t('admin.rbac.groups.columns.actions'), width: '7%' },
  ];

  const renderCell = (group, column) => {
    switch (column.key) {
      case 'id':
        return <span className="text-neutral-400 font-mono">#{group.id}</span>;

      case 'name':
        return <span className="text-neutral-50">{group.name}</span>;

      case 'description':
        return <span className="text-neutral-300">{group.description || '-'}</span>;

      case 'role':
        return <span className="text-geek-400">{rolesMap[group.role_id] || '-'}</span>;

      case 'users':
        return <span className="text-neutral-300">{group.users}</span>;

      case 'default':
        return group.default ? (
          <StatusTag type="info" text={t('common.yes')} />
        ) : (
          <StatusTag type="error" text={t('common.no')} />
        );

      case 'actions':
        return (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-geek-400"
              onClick={(e) => {
                e.stopPropagation();
                onManageUsers?.(group);
              }}
            >
              <IconUsers size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-yellow-400"
              onClick={(e) => {
                e.stopPropagation();
                onEditGroup?.(group);
              }}
            >
              <IconEdit size={18} />
            </Button>
            {!group.default && (
              <Button
                variant="ghost"
                size="icon"
                className="!bg-transparent !text-red-400"
                onClick={(e) => {
                  e.stopPropagation();
                  onDeleteGroup?.(group);
                }}
              >
                <IconTrash size={18} />
              </Button>
            )}
          </div>
        );

      default:
        return group[column.key];
    }
  };

  const paginationComponent = totalCount > pageSize && (
    <Pagination
      total={Math.ceil(totalCount / pageSize)}
      current={currentPage}
      onChange={onPageChange}
      showTotal
      totalItems={totalCount}
      showJumpTo
    />
  );

  return (
    <div className="w-full mx-auto">
      <div className="flex items-center justify-end mb-8">
        <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={onCreateGroup}>
          {t('admin.rbac.groups.actions.create')}
        </Button>
      </div>

      <List
        columns={columns}
        data={groups}
        renderCell={renderCell}
        loading={loading}
        empty={groups.length === 0}
        footer={paginationComponent}
      />
    </div>
  );
}

export default AdminGroups;
