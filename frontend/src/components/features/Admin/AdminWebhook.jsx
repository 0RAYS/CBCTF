import { motion } from 'motion/react';
import { List, StatusTag } from '../../common';
import { IconPlus, IconEdit, IconTrash, IconHistory } from '@tabler/icons-react';
import { Button, Pagination } from '../../../components/common';
import { useTranslation } from 'react-i18next';

/**
 * Webhook配置管理展示组件
 * @param {Object} props
 * @param {Array} props.webhooks - Webhook配置列表数据
 * @param {number} props.totalCount - 总数量
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {boolean} props.loading - 是否加载中
 * @param {function} props.onPageChange - 页码改变回调
 * @param {function} props.onCreateWebhook - 创建Webhook配置回调
 * @param {function} props.onEditWebhook - 编辑Webhook配置回调
 * @param {function} props.onDeleteWebhook - 删除Webhook配置回调
 * @param {function} props.onViewHistory - 查看历史记录回调
 * @param {function} props.onWebhookClick - Webhook配置点击回调
 */
function AdminWebhook({
  webhooks = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  loading = false,
  onPageChange,
  onCreateWebhook,
  onEditWebhook,
  onDeleteWebhook,
  onViewHistory,
  onWebhookClick,
}) {
  const { t, i18n } = useTranslation();
  // 列定义
  const columns = [
    { key: 'id', label: t('admin.webhook.list.columns.id'), width: '5%' },
    { key: 'name', label: t('admin.webhook.list.columns.name'), width: '15%' },
    { key: 'url', label: t('admin.webhook.list.columns.url'), width: '20%' },
    { key: 'method', label: t('admin.webhook.list.columns.method'), width: '8%' },
    { key: 'status', label: t('admin.webhook.list.columns.status'), width: '10%' },
    { key: 'events', label: t('admin.webhook.list.columns.events'), width: '15%' },
    { key: 'success', label: t('admin.webhook.list.columns.success'), width: '10%' },
    { key: 'failure', label: t('admin.webhook.list.columns.failure'), width: '10%' },
    { key: 'actions', label: t('admin.webhook.list.columns.actions'), width: '7%' },
  ];

  // 自定义单元格渲染
  const renderCell = (webhook, column) => {
    switch (column.key) {
      case 'id':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50 font-medium">#{webhook.id}</span>
          </div>
        );

      case 'name':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50 font-medium">{webhook.name}</span>
          </div>
        );

      case 'url':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-300 font-mono text-sm truncate max-w-48" title={webhook.url}>
              {webhook.url}
            </span>
          </div>
        );

      case 'method':
        return (
          <div className="flex flex-col">
            <span
              className={`text-sm font-mono px-2 py-1 rounded ${
                webhook.method === 'GET'
                  ? 'bg-blue-400/20 text-blue-400'
                  : webhook.method === 'POST'
                    ? 'bg-green-400/20 text-green-400'
                    : webhook.method === 'PUT'
                      ? 'bg-yellow-400/20 text-yellow-400'
                      : webhook.method === 'DELETE'
                        ? 'bg-red-400/20 text-red-400'
                        : 'bg-neutral-400/20 text-neutral-400'
              }`}
            >
              {webhook.method}
            </span>
          </div>
        );

      case 'status':
        return (
          <div className="flex flex-wrap gap-2">
            {webhook.on ? (
              <StatusTag type="success" text={t('admin.webhook.list.status.enabled')} />
            ) : (
              <StatusTag type="warning" text={t('admin.webhook.list.status.disabled')} />
            )}
          </div>
        );

      case 'events':
        return (
          <div className="flex flex-col">
            {webhook.events && webhook.events.length > 0 ? (
              <div className="flex flex-wrap gap-1">
                {webhook.events.slice(0, 2).map((event, index) => (
                  <span key={index} className="text-xs bg-neutral-700 text-neutral-300 px-1 py-0.5 rounded">
                    {event}
                  </span>
                ))}
                {webhook.events.length > 2 && (
                  <span className="text-xs text-neutral-400">+{webhook.events.length - 2}</span>
                )}
              </div>
            ) : (
              <span className="text-neutral-500 text-sm">{t('admin.webhook.list.events.all')}</span>
            )}
          </div>
        );

      case 'success':
        return (
          <div className="flex flex-col">
            <span className="text-green-400 font-medium">{webhook.success || 0}</span>
            {webhook.success_last && (
              <span className="text-xs text-neutral-400">
                {new Date(webhook.success_last).toLocaleString(i18n.language || 'en-US')}
              </span>
            )}
          </div>
        );

      case 'failure':
        return (
          <div className="flex flex-col">
            <span className="text-red-400 font-medium">{webhook.failure || 0}</span>
            {webhook.failure_last && (
              <span className="text-xs text-neutral-400">
                {new Date(webhook.failure_last).toLocaleString(i18n.language || 'en-US')}
              </span>
            )}
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
                onEditWebhook?.(webhook);
              }}
              className="p-1"
              title={t('admin.webhook.list.actions.edit')}
            >
              <IconEdit size={16} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                onViewHistory?.(webhook);
              }}
              className="p-1"
              title={t('admin.webhook.list.actions.viewHistory')}
            >
              <IconHistory size={16} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                onDeleteWebhook?.(webhook);
              }}
              className="p-1 text-red-400 hover:text-red-300"
              title={t('admin.webhook.list.actions.delete')}
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
          <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={onCreateWebhook}>
            {t('admin.webhook.list.actions.add')}
          </Button>
        </div>

        {/* 数据表格 */}
        <List
          data={webhooks}
          columns={columns}
          renderCell={renderCell}
          onRowClick={onWebhookClick}
          loading={loading}
          empty={webhooks.length === 0}
          emptyContent={t('admin.webhook.list.empty')}
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

export default AdminWebhook;
