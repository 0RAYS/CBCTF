import { motion } from 'motion/react';
import { Button, List, Pagination, StatusTag } from '@/components/common';
import { IconEdit, IconPlus, IconTrash } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';

/** OAuth Provider管理展示组件 */
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

  const columns = [
    {
      key: 'picture',
      label: t('admin.oauthProviders.columns.logo'),
      width: '5%',
    },
    {
      key: 'provider',
      label: t('admin.oauthProviders.columns.provider'),
      width: '10%',
    },
    { key: 'uri', label: t('admin.oauthProviders.columns.uri'), width: '10%' },
    {
      key: 'protocol',
      label: t('admin.oauthProviders.columns.protocol'),
      width: '10%',
    },
    {
      key: 'status',
      label: t('admin.oauthProviders.columns.status'),
      width: '8%',
    },
    {
      key: 'callback_url',
      label: t('admin.oauthProviders.columns.callback'),
      width: '25%',
    },
  ];

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
            <img src={provider.picture} alt={provider.provider} loading="lazy" className="w-full h-full object-cover" />
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

      case 'protocol':
        return (
          <span className="text-neutral-300 font-mono text-sm">
            {provider.protocol === 'cas'
              ? t('admin.oauthProviders.form.protocolCAS')
              : t('admin.oauthProviders.form.protocolOAuth2')}
          </span>
        );

      case 'status':
        return (
          <div className="flex flex-wrap gap-2">
            {provider.on ? (
              <StatusTag type="success" text={t('common.toggle.enabled')} />
            ) : (
              <StatusTag type="warning" text={t('common.toggle.disabled')} />
            )}
          </div>
        );

      case 'callback_url': {
        const url = provider.callback_url || '';
        return (
          <div className="flex items-center gap-2 min-w-0 w-full">
            <span className="text-neutral-300 font-mono text-sm truncate block min-w-0 flex-1" title={url}>
              {url}
            </span>
          </div>
        );
      }

      default:
        return null;
    }
  };

  return (
    <div className="w-full mx-auto">
      <motion.div
        className="rounded-md bg-neutral-900 overflow-hidden"
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

        <List
          data={providers}
          columns={columns}
          renderCell={renderCell}
          getRowActions={(provider) => [
            {
              key: 'edit',
              label: t('common.edit'),
              icon: <IconEdit size={18} />,
              onClick: () => onEditProvider?.(provider),
            },
            {
              key: 'delete',
              label: t('common.delete'),
              icon: <IconTrash size={18} />,
              danger: true,
              onClick: () => onDeleteProvider?.(provider),
            },
          ]}
          onRowClick={onEditProvider || onProviderClick}
          loading={loading}
          emptyContent={t('admin.oauthProviders.empty')}
        />

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
