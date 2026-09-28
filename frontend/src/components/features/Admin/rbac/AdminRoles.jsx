import { Button, List, Pagination, StatusTag } from '../../../common';
import { IconEdit, IconPlus, IconShield, IconTrash } from '@tabler/icons-react';
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
    { key: 'id', label: t('admin.rbac.roles.columns.id'), width: '5%' },
    { key: 'name', label: t('admin.rbac.roles.columns.name'), width: '10%' },
    {
      key: 'description',
      label: t('admin.rbac.roles.columns.description'),
      width: '15%',
    },
    {
      key: 'default',
      label: t('admin.rbac.roles.columns.default'),
      width: '10%',
    },
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
          <StatusTag type="error" text={t('common.no')} />
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
        onRowClick={onEditRole}
        actionsColumn={{ label: t('admin.rbac.roles.columns.actions'), width: 88 }}
        getRowActions={(role) => [
          {
            key: 'permissions',
            inline: true,
            label: t('admin.rbac.roles.modal.permissionsTitle'),
            icon: <IconShield size={18} />,
            onClick: () => onManagePermissions?.(role),
          },
          {
            key: 'edit',
            label: t('common.edit'),
            icon: <IconEdit size={18} />,
            onClick: () => onEditRole?.(role),
          },
          {
            key: 'delete',
            label: t('common.delete'),
            icon: <IconTrash size={18} />,
            danger: true,
            hidden: role.default,
            onClick: () => onDeleteRole?.(role),
          },
        ]}
        loading={loading}
        empty={roles.length === 0}
        footer={paginationComponent}
      />
    </div>
  );
}

export default AdminRoles;
