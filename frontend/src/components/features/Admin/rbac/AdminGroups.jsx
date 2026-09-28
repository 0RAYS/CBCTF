import { Button, List, Pagination, StatusTag } from '../../../common';
import { IconEdit, IconPlus, IconTrash, IconUsers } from '@tabler/icons-react';
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
    {
      key: 'description',
      label: t('admin.rbac.groups.columns.description'),
      width: '10%',
    },
    { key: 'role', label: t('admin.rbac.groups.columns.role'), width: '10%' },
    { key: 'users', label: t('admin.rbac.groups.columns.users'), width: '7%' },
    {
      key: 'default',
      label: t('admin.rbac.groups.columns.default'),
      width: '10%',
    },
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
        onRowClick={onEditGroup}
        actionsColumn={{ label: t('admin.rbac.groups.columns.actions'), width: 88 }}
        getRowActions={(group) => [
          {
            key: 'members',
            inline: true,
            label: t('admin.rbac.groups.modal.usersTitle'),
            icon: <IconUsers size={18} />,
            onClick: () => onManageUsers?.(group),
          },
          {
            key: 'edit',
            label: t('common.edit'),
            icon: <IconEdit size={18} />,
            onClick: () => onEditGroup?.(group),
          },
          {
            key: 'delete',
            label: t('common.delete'),
            icon: <IconTrash size={18} />,
            danger: true,
            hidden: group.default,
            onClick: () => onDeleteGroup?.(group),
          },
        ]}
        loading={loading}
        empty={groups.length === 0}
        footer={paginationComponent}
      />
    </div>
  );
}

export default AdminGroups;
