import { motion } from 'motion/react';
import { List, StatusTag } from '../../common';
import { IconEye, IconWebhook } from '@tabler/icons-react';
import { Button, Pagination, EmptyState } from '../../../components/common';
import { useTranslation } from 'react-i18next';

/**
 * Webhook历史记录展示组件
 * @param {Object} props
 * @param {Array} props.webhookHistory - Webhook历史记录列表数据
 * @param {number} props.totalCount - 总数量
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {boolean} props.loading - 是否加载中
 * @param {function} props.onPageChange - 页码改变回调
 * @param {function} props.onViewDetail - 查看详情回调
 * @param {function} props.onHistoryClick - 历史记录点击回调
 * @param {string} props.webhookName - 当前查看的Webhook名称（可选）
 */
function AdminWebhookHistory({
  webhookHistory = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  loading = false,
  onPageChange,
  onViewDetail,
  onHistoryClick,
  webhookName,
}) {
  const { t, i18n } = useTranslation();
  // 列定义
  const columns = [
    { key: 'id', label: t('admin.webhook.historyList.columns.id'), width: '5%' },
    { key: 'webhook', label: t('admin.webhook.historyList.columns.webhook'), width: '15%' },
    { key: 'event', label: t('admin.webhook.historyList.columns.event'), width: '15%' },
    { key: 'status', label: t('admin.webhook.historyList.columns.status'), width: '10%' },
    { key: 'resp', label: t('admin.webhook.historyList.columns.responseCode'), width: '10%' },
    { key: 'duration', label: t('admin.webhook.historyList.columns.duration'), width: '10%' },
    { key: 'time', label: t('admin.webhook.historyList.columns.time'), width: '12%' },
    { key: 'actions', label: t('admin.webhook.historyList.columns.actions'), width: '5%' },
  ];

  // 自定义单元格渲染
  const renderCell = (history, column) => {
    switch (column.key) {
      case 'id':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50 font-medium">#{history.id}</span>
          </div>
        );

      case 'webhook':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50 font-medium">{history.webhook}</span>
          </div>
        );

      case 'event':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50 font-medium">{history.event}</span>
          </div>
        );

      case 'status':
        return (
          <div className="flex flex-wrap gap-2">
            {history.success ? (
              <StatusTag type="success" text={t('admin.webhook.history.statusSuccess')} />
            ) : (
              <StatusTag type="error" text={t('admin.webhook.history.statusFailed')} />
            )}
          </div>
        );

      case 'resp':
        return (
          <div className="flex flex-col">
            <span
              className={`text-sm font-mono ${
                history.resp >= 200 && history.resp < 300
                  ? 'text-green-400'
                  : history.resp >= 400 && history.resp < 500
                    ? 'text-yellow-400'
                    : history.resp >= 500
                      ? 'text-red-400'
                      : 'text-neutral-400'
              }`}
            >
              {history.resp || '-'}
            </span>
          </div>
        );

      case 'duration':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-300 text-sm">{history.duration ? `${history.duration}ms` : '-'}</span>
          </div>
        );

      case 'time':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-300 text-sm">
              {new Date(history.time).toLocaleString(i18n.language || 'en-US')}
            </span>
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
                onViewDetail?.(history);
              }}
              className="p-1"
              title={t('admin.webhook.historyList.actions.viewDetail')}
            >
              <IconEye size={16} />
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
        className="rounded-md bg-neutral-900 overflow-hidden"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        {webhookName && (
          <div className="flex items-center gap-3 mb-6">
            <IconWebhook size={20} className="text-neutral-400" />
            <p className="text-sm text-neutral-400">
              {t('admin.webhook.historyList.currentViewing', { name: webhookName })}
            </p>
          </div>
        )}

        {/* 数据表格 */}
        <List
          data={webhookHistory}
          columns={columns}
          renderCell={renderCell}
          onRowClick={onHistoryClick}
          loading={loading}
          empty={webhookHistory.length === 0}
          emptyContent={<EmptyState title={t('admin.webhook.historyList.empty')} />}
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

export default AdminWebhookHistory;
