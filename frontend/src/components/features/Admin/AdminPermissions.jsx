import { List } from '../../common';
import { IconEdit } from '@tabler/icons-react';
import { Button, Pagination } from '../../../components/common';
import { useTranslation } from 'react-i18next';

function AdminPermissions({
  permissions = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 50,
  loading = false,
  onPageChange,
  onEditPermission,
}) {
  const { t } = useTranslation();

  const columns = [
    { key: 'id', label: t('admin.rbac.permissions.columns.id'), width: '8%' },
    { key: 'name', label: t('admin.rbac.permissions.columns.name'), width: '22%' },
    { key: 'resource', label: t('admin.rbac.permissions.columns.resource'), width: '15%' },
    { key: 'operation', label: t('admin.rbac.permissions.columns.operation'), width: '15%' },
    { key: 'description', label: t('admin.rbac.permissions.columns.description'), width: '30%' },
    { key: 'actions', label: t('admin.rbac.permissions.columns.actions'), width: '10%' },
  ];

  const renderCell = (permission, column) => {
    switch (column.key) {
      case 'id':
        return <span className="text-neutral-400 font-mono">#{permission.id}</span>;

      case 'name':
        return <span className="text-neutral-50 font-mono">{permission.name}</span>;

      case 'resource':
        return <span className="text-blue-400">{permission.resource}</span>;

      case 'operation':
        return <span className="text-neutral-300">{permission.operation}</span>;

      case 'description':
        return <span className="text-neutral-300">{permission.description || '-'}</span>;

      case 'actions':
        return (
          <div className="flex gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="!bg-yellow-400/20 !text-yellow-400 hover:!bg-yellow-400/30"
              onClick={(e) => {
                e.stopPropagation();
                onEditPermission?.(permission);
              }}
            >
              <IconEdit size={18} />
            </Button>
          </div>
        );

      default:
        return permission[column.key];
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
      <List
        columns={columns}
        data={permissions}
        renderCell={renderCell}
        loading={loading}
        empty={permissions.length === 0}
        footer={paginationComponent}
      />
    </div>
  );
}

export default AdminPermissions;
