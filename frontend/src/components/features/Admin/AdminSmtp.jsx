import { motion } from 'motion/react';
import { List, StatusTag } from '../../common';
import { IconPlus, IconEdit, IconTrash, IconHistory } from '@tabler/icons-react';
import { Button, Pagination } from '../../../components/common';
import { useTranslation } from 'react-i18next';

/**
 * SMTP配置管理展示组件
 * @param {Object} props
 * @param {Array} props.smtpConfigs - SMTP配置列表数据
 * @param {number} props.totalCount - 总数量
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {boolean} props.loading - 是否加载中
 * @param {function} props.onPageChange - 页码改变回调
 * @param {function} props.onCreateSmtp - 创建SMTP配置回调
 * @param {function} props.onEditSmtp - 编辑SMTP配置回调
 * @param {function} props.onDeleteSmtp - 删除SMTP配置回调
 * @param {function} props.onSmtpClick - SMTP配置点击回调
 */
function AdminSmtp({
  smtpConfigs = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  loading = false,
  onPageChange,
  onCreateSmtp,
  onEditSmtp,
  onDeleteSmtp,
  onSmtpClick,
  onViewHistory,
}) {
  const { t, i18n } = useTranslation();
  // 列定义
  const columns = [
    { key: 'id', label: t('admin.smtp.columns.id'), width: '5%' },
    { key: 'address', label: t('admin.smtp.columns.address'), width: '15%' },
    { key: 'host', label: t('admin.smtp.columns.host'), width: '15%' },
    { key: 'port', label: t('admin.smtp.columns.port'), width: '8%' },
    { key: 'status', label: t('admin.smtp.columns.status'), width: '10%' },
    { key: 'success', label: t('admin.smtp.columns.success'), width: '10%' },
    { key: 'failure', label: t('admin.smtp.columns.failure'), width: '10%' },
    { key: 'actions', label: t('admin.smtp.columns.actions'), width: '5%' },
  ];

  // 自定义单元格渲染
  const renderCell = (smtp, column) => {
    switch (column.key) {
      case 'id':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50 font-medium">#{smtp.id}</span>
          </div>
        );

      case 'address':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50 font-medium">{smtp.address}</span>
          </div>
        );

      case 'host':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-300 font-mono text-sm">{smtp.host}</span>
          </div>
        );

      case 'port':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-300 font-mono text-sm">{smtp.port}</span>
          </div>
        );

      case 'status':
        return (
          <div className="flex flex-wrap gap-2">
            {smtp.on ? (
              <StatusTag type="success" text={t('admin.smtp.status.enabled')} />
            ) : (
              <StatusTag type="warning" text={t('admin.smtp.status.disabled')} />
            )}
          </div>
        );

      case 'success':
        return (
          <div className="flex flex-col">
            <span className="text-green-400 font-medium">{smtp.success || 0}</span>
            {smtp.success_last && (
              <span className="text-xs text-neutral-400">
                {new Date(smtp.success_last).toLocaleString(i18n.language || 'en-US')}
              </span>
            )}
          </div>
        );

      case 'failure':
        return (
          <div className="flex flex-col">
            <span className="text-red-400 font-medium">{smtp.failure || 0}</span>
            {smtp.failure_last && (
              <span className="text-xs text-neutral-400">
                {new Date(smtp.failure_last).toLocaleString(i18n.language || 'en-US')}
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
                onEditSmtp?.(smtp);
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
                onViewHistory?.(smtp);
              }}
              className="p-1"
              title={t('admin.smtp.actions.viewHistory')}
            >
              <IconHistory size={16} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                onDeleteSmtp?.(smtp);
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
          <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={onCreateSmtp}>
            {t('admin.smtp.actions.add')}
          </Button>
        </div>

        {/* 数据表格 */}
        <List
          data={smtpConfigs}
          columns={columns}
          renderCell={renderCell}
          onRowClick={onSmtpClick}
          loading={loading}
          empty={smtpConfigs.length === 0}
          emptyContent={t('admin.smtp.empty')}
        />

        {/* 分页 */}
        {totalCount > pageSize && (
          <Pagination
            current={currentPage}
            total={Math.ceil(totalCount / pageSize)}
            onChange={onPageChange}
            showTotal
            totalItems={totalCount}
            showJumpTo
          />
        )}
      </motion.div>
    </div>
  );
}

export default AdminSmtp;
