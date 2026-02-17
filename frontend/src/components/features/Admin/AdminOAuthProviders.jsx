import { motion } from 'motion/react';
import { List, StatusTag } from '../../common';
import { IconPlus, IconEdit, IconTrash } from '@tabler/icons-react';
import { Button, Pagination } from '../../../components/common';
import { useTranslation } from 'react-i18next';

/**
 * OAuth Provider管理展示组件
 * @param {Object} props
 * @param {Array} props.providers - OAuth Provider列表数据
 * @param {number} props.totalCount - 总数量
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {boolean} props.loading - 是否加载中
 * @param {function} props.onPageChange - 页码改变回调
 * @param {function} props.onCreateProvider - 创建Provider回调
 * @param {function} props.onEditProvider - 编辑Provider回调
 * @param {function} props.onDeleteProvider - 删除Provider回调
 * @param {function} props.onProviderClick - Provider点击回调
 */
function AdminOAuthProviders({
  providers = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  loading = false,
  onPageChange,
  onCreateProvider,
  onEditProvider,
  onDeleteProvider,
  onProviderClick,
  onPictureUpload,
}) {
  const { t } = useTranslation();

  // 列定义
  const columns = [
    { key: 'picture', label: t('admin.oauthProviders.columns.logo'), width: '5%' },
    { key: 'provider', label: t('admin.oauthProviders.columns.provider'), width: '5%' },
    { key: 'uri', label: t('admin.oauthProviders.columns.uri'), width: '5%' },
    { key: 'status', label: t('admin.oauthProviders.columns.status'), width: '5%' },
    // { key: 'auth_url', label: 'AuthURL', width: '10%' },
    // { key: 'token_url', label: 'TokenURL', width: '10%' },
    // { key: 'user_info_url', label: 'UserInfoURL', width: '10%' },
    { key: 'callback_url', label: t('admin.oauthProviders.columns.callback'), width: '10%' },
    { key: 'actions', label: t('admin.oauthProviders.columns.actions'), width: '5%' },
  ];

  // 自定义单元格渲染
  const renderCell = (provider, column) => {
    switch (column.key) {
      case 'picture':
        return (
          <div
            className="relative w-10 h-10 rounded-full overflow-hidden cursor-pointer group"
            onClick={(e) => {
              e.stopPropagation();
              onPictureUpload?.(provider);
            }}
            title={t('admin.oauthProviders.actions.uploadLogo')}
          >
            <img src={provider.picture} alt={provider.provider} className="w-full h-full object-cover" />
            <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
              <IconEdit size={16} className="text-white" />
            </div>
          </div>
        );

      case 'provider':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50 font-medium">{provider.provider}</span>
            <span className="text-xs text-neutral-400">
              {t('admin.oauthProviders.columns.id', { id: provider.id })}
            </span>
          </div>
        );

      case 'uri':
        return <span className="text-neutral-300 font-mono text-sm">{provider.uri}</span>;

      case 'status':
        return (
          <div className="flex flex-wrap gap-2">
            {provider.on ? (
              <StatusTag type="success" text={t('admin.oauthProviders.status.enabled')} />
            ) : (
              <StatusTag type="warning" text={t('admin.oauthProviders.status.disabled')} />
            )}
          </div>
        );

      // case 'auth_url':
      //   return (
      //     <div className="max-w-12 truncate" title={provider.auth_url}>
      //       <span className="text-neutral-300 text-sm">
      //         {provider.auth_url.replace('http://', '').replace('https://', '')}
      //       </span>
      //     </div>
      //   );
      //
      // case 'token_url':
      //   return (
      //     <div className="max-w-12 truncate" title={provider.token_url}>
      //       <span className="text-neutral-300 text-sm">
      //         {provider.token_url.replace('http://', '').replace('https://', '')}
      //       </span>
      //     </div>
      //   );
      //
      // case 'user_info_url':
      //   return (
      //     <div className="max-w-12 truncate" title={provider.user_info_url}>
      //       <span className="text-neutral-300 text-sm">
      //         {provider.user_info_url.replace('http://', '').replace('https://', '')}
      //       </span>
      //     </div>
      //   );

      case 'callback_url':
        return (
          <div className="max-w-25 truncate" title={provider.callback_url}>
            <span className="text-neutral-300 text-sm">{provider.callback_url || ''}</span>
          </div>
        );

      case 'actions':
        return (
          <div className="flex gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                onEditProvider?.(provider);
              }}
              className="p-1"
            >
              <IconEdit size={16} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                onDeleteProvider?.(provider);
              }}
              className="p-1 text-red-400 hover:text-red-300"
            >
              <IconTrash size={16} />
            </Button>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="w-full mx-auto">
      <motion.div
        className="rounded-md bg-black/30 backdrop-blur-[2px] overflow-hidden"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <div className="flex items-center justify-end mb-6">
          <Button
            variant="primary"
            size="sm"
            align="icon-left"
            icon={<IconPlus size={16} />}
            onClick={onCreateProvider}
          >
            {t('admin.oauthProviders.actions.add')}
          </Button>
        </div>

        {/* 数据表格 */}
        <List
          data={providers}
          columns={columns}
          renderCell={renderCell}
          onRowClick={onProviderClick}
          loading={loading}
          emptyMessage={t('admin.oauthProviders.empty')}
        />

        {/* 分页 */}
        {totalCount > pageSize && (
          <Pagination
            current={currentPage}
            total={Math.ceil(totalCount / pageSize)}
            onChange={onPageChange}
            showTotal={true}
            totalItems={totalCount}
          />
        )}
      </motion.div>
    </div>
  );
}

export default AdminOAuthProviders;
