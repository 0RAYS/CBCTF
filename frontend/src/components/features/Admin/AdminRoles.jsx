import { List, StatusTag } from '../../common';
import { IconPlus, IconEdit, IconTrash, IconShield } from '@tabler/icons-react';
import { Button, Pagination } from '../../../components/common';
import { useTranslation } from 'react-i18next';

function AdminRoles({
  roles = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  loading = false,
  onPageChange,
  onCreateRole,
  onEditRole,
  onDeleteRole,
  onManagePermissions,
}) {
  const { t } = useTranslation();

  const columns = [
    { key: 'id', label: t('admin.rbac.roles.columns.id'), width: '10%' },
    { key: 'name', label: t('admin.rbac.roles.columns.name'), width: '20%' },
    { key: 'description', label: t('admin.rbac.roles.columns.description'), width: '30%' },
    { key: 'default', label: t('admin.rbac.roles.columns.default'), width: '15%' },
    { key: 'actions', label: t('admin.rbac.roles.columns.actions'), width: '25%' },
  ];

  const renderCell = (role, column) => {
    switch (column.key) {
      case 'id':
        return <span className="text-neutral-400 font-mono">#{role.id}</span>;

      case 'name':
        return <span className="text-neutral-50">{role.name}</span>;

      case 'description':
        return <span className="text-neutral-300">{role.description || '-'}</span>;

      case 'default':
        return role.default ? (
          <StatusTag type="info" text={t('common.yes')} />
        ) : (
          <span className="text-neutral-500">{t('common.no')}</span>
        );

      case 'actions':
        return (
          <div className="flex gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="!bg-blue-400/20 !text-blue-400 hover:!bg-blue-400/30"
              onClick={(e) => {
                e.stopPropagation();
                onManagePermissions?.(role);
              }}
            >
              <IconShield size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!bg-yellow-400/20 !text-yellow-400 hover:!bg-yellow-400/30"
              onClick={(e) => {
                e.stopPropagation();
                onEditRole?.(role);
              }}
            >
              <IconEdit size={18} />
            </Button>
            {!role.default && (
              <Button
                variant="ghost"
                size="icon"
                className="!bg-red-400/20 !text-red-400 hover:!bg-red-400/30"
                onClick={(e) => {
                  e.stopPropagation();
                  onDeleteRole?.(role);
                }}
              >
                <IconTrash size={18} />
              </Button>
            )}
          </div>
        );

      default:
        return role[column.key];
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
        <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={onCreateRole}>
          {t('admin.rbac.roles.actions.create')}
        </Button>
      </div>

      <List
        columns={columns}
        data={roles}
        renderCell={renderCell}
        loading={loading}
        empty={roles.length === 0}
        footer={paginationComponent}
      />
    </div>
  );
}

export default AdminRoles;
