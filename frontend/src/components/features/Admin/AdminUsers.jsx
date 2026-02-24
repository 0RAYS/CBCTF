import { List, StatusTag, Avatar } from '../../common';
import { IconPlus, IconEdit, IconTrash, IconSearch } from '@tabler/icons-react';
import { Button, Input, Pagination, Spinner } from '../../../components/common';
import { useTranslation } from 'react-i18next';
import AdminUserDetailDialog from './AdminUserDetailDialog';

/**
 * 用户管理展示组件
 * @param {Object} props
 * @param {Array} props.users - 用户列表数据
 * @param {number} props.totalCount - 总用户数量
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {boolean} props.loading - 是否加载中
 * @param {function} props.onPageChange - 页码改变回调
 * @param {function} props.onCreateUser - 创建用户回调
 * @param {function} props.onEditUser - 编辑用户回调
 * @param {function} props.onDeleteUser - 删除用户回调
 * @param {function} props.onPictureUpload - 上传头像回调
 * @param {string} props.searchQuery - 搜索查询
 * @param {boolean} props.searchLoading - 搜索加载状态
 * @param {boolean} props.isSearchMode - 是否处于搜索模式
 * @param {function} props.onSearchChange - 搜索查询变化回调
 */
function AdminUsers({
  users = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 6,
  loading = false,
  onPageChange,
  onCreateUser,
  onEditUser,
  onDeleteUser,
  onPictureUpload,
  searchQuery = '',
  searchLoading = false,
  isSearchMode = false,
  onSearchChange,
  onRowClick,
  showDetailDialog = false,
  detailUser = null,
  onDetailClose,
}) {
  const { t } = useTranslation();

  // 列定义
  const columns = [
    { key: 'picture', label: t('admin.users.columns.picture'), width: '10%' },
    { key: 'name', label: t('admin.users.columns.name'), width: '15%' },
    { key: 'email', label: t('admin.users.columns.email'), width: '20%' },
    { key: 'status', label: t('admin.users.columns.status'), width: '15%' },
    { key: 'contests', label: t('admin.users.columns.contests'), width: '10%' },
    { key: 'teams', label: t('admin.users.columns.teams'), width: '10%' },
    { key: 'actions', label: t('admin.users.columns.actions'), width: '5%' },
  ];

  // 自定义单元格渲染
  const renderCell = (user, column) => {
    switch (column.key) {
      case 'picture':
        return (
          <div
            className="relative w-10 h-10 rounded-full overflow-hidden cursor-pointer group"
            onClick={(e) => {
              e.stopPropagation();
              onPictureUpload?.(user);
            }}
          >
            <Avatar src={user.picture} name={user.name} size="sm" shape="circle" />
            <div className="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
              <span className="text-neutral-300 text-xs">{t('admin.users.picture.replace')}</span>
            </div>
          </div>
        );

      case 'name':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50">{user.name}</span>
            <span className="text-xs text-neutral-400">{t('admin.users.id', { id: user.id })}</span>
          </div>
        );

      case 'status':
        return (
          <div className="flex flex-wrap gap-2">
            {user.verified && <StatusTag type="success" text={t('admin.users.status.verified')} />}
            {user.banned && <StatusTag type="error" text={t('admin.users.status.banned')} />}
            {user.hidden && <StatusTag type="warning" text={t('admin.users.status.hidden')} />}
          </div>
        );

      case 'actions':
        return (
          <div className="flex gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="!bg-yellow-400/20 !text-yellow-400 hover:!bg-yellow-400/30"
              onClick={(e) => {
                e.stopPropagation();
                onEditUser?.(user);
              }}
            >
              <IconEdit size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!bg-red-400/20 !text-red-400 hover:!bg-red-400/30"
              onClick={(e) => {
                e.stopPropagation();
                onDeleteUser?.(user);
              }}
            >
              <IconTrash size={18} />
            </Button>
          </div>
        );

      default:
        return user[column.key];
    }
  };

  // 分页组件
  const paginationComponent = totalCount > pageSize && !isSearchMode && (
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
        <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={onCreateUser}>
          {t('admin.users.actions.create')}
        </Button>
      </div>

      {/* 搜索框 */}
      <div className="mb-6">
        <div className="flex items-center gap-4">
          <div className="relative flex-1 max-w-md">
            <label className="block text-sm font-mono text-neutral-400 mb-2">{t('admin.users.search.label')}</label>
            <Input
              type="search"
              value={searchQuery}
              placeholder={t('admin.users.search.placeholder')}
              onChange={(e) => onSearchChange?.(e.target.value)}
              icon={<IconSearch size={16} />}
              iconRight={searchLoading && <Spinner size="sm" />}
            />
          </div>
        </div>
      </div>

      <List
        columns={columns}
        data={users}
        renderCell={renderCell}
        loading={loading}
        empty={users.length === 0}
        onRowClick={onRowClick}
        emptyContent={
          isSearchMode ? (
            <div className="flex flex-col items-center justify-center space-y-2">
              <svg className="w-12 h-12 text-neutral-300/30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={1.5}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
              <span className="font-mono">未找到匹配的用户</span>
            </div>
          ) : undefined
        }
        footer={paginationComponent}
        searchQuery={searchQuery}
        searchLoading={searchLoading}
        onSearchChange={onSearchChange}
      />

      <AdminUserDetailDialog isOpen={showDetailDialog} onClose={onDetailClose} user={detailUser} />
    </div>
  );
}

export default AdminUsers;
