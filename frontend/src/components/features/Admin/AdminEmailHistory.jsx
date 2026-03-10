import { motion } from 'motion/react';
import { List, StatusTag } from '../../common';
import { IconEye, IconMail } from '@tabler/icons-react';
import { Button, Pagination, EmptyState } from '../../../components/common';
import { useTranslation } from 'react-i18next';

/**
 * 邮件发送历史记录展示组件
 * @param {Object} props
 * @param {Array} props.emailHistory - 邮件历史记录列表数据
 * @param {number} props.totalCount - 总数量
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {boolean} props.loading - 是否加载中
 * @param {function} props.onPageChange - 页码改变回调
 * @param {function} props.onViewEmail - 查看邮件详情回调
 * @param {function} props.onEmailClick - 邮件点击回调
 * @param {string} props.smtpAddress - 当前查看的SMTP地址（可选）
 */
function AdminEmailHistory({
  emailHistory = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  loading = false,
  onPageChange,
  onViewEmail,
  onEmailClick,
  smtpAddress,
}) {
  const { t, i18n } = useTranslation();
  // 列定义
  const columns = [
    { key: 'id', label: t('admin.smtp.history.columns.id'), width: '5%' },
    { key: 'from', label: t('admin.smtp.history.columns.from'), width: '15%' },
    { key: 'to', label: t('admin.smtp.history.columns.to'), width: '15%' },
    { key: 'subject', label: t('admin.smtp.history.columns.subject'), width: '10%' },
    { key: 'status', label: t('admin.smtp.history.columns.status'), width: '10%' },
    { key: 'time', label: t('admin.smtp.history.columns.time'), width: '15%' },
    { key: 'actions', label: t('admin.smtp.history.columns.actions'), width: '5%' },
  ];

  // 自定义单元格渲染
  const renderCell = (email, column) => {
    switch (column.key) {
      case 'id':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50 font-medium">#{email.id}</span>
          </div>
        );

      case 'from':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-50 font-medium">{email.from}</span>
          </div>
        );

      case 'to':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-300 font-mono text-sm">{email.to}</span>
          </div>
        );

      case 'subject':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-300 text-sm truncate max-w-48" title={email.subject}>
              {email.subject || t('admin.smtp.email.noSubject')}
            </span>
          </div>
        );

      case 'status':
        return (
          <div className="flex flex-wrap gap-2">
            {email.success ? (
              <StatusTag type="success" text={t('admin.smtp.email.statusSuccess')} />
            ) : (
              <StatusTag type="error" text={t('admin.smtp.email.statusFailed')} />
            )}
          </div>
        );

      case 'time':
        return (
          <div className="flex flex-col">
            <span className="text-neutral-300 text-sm">
              {new Date(email.time).toLocaleString(i18n.language || 'en-US')}
            </span>
          </div>
        );

      case 'actions':
        return (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              onClick={(e) => {
                e.stopPropagation();
                onViewEmail?.(email);
              }}
              className="p-1"
              title={t('admin.smtp.history.actions.viewDetail')}
            >
              <IconEye size={18} />
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
        {smtpAddress && (
          <div className="flex items-center gap-3 mb-6">
            <IconMail size={20} className="text-neutral-400" />
            <p className="text-sm text-neutral-400">
              {t('admin.smtp.history.currentViewing', { address: smtpAddress })}
            </p>
          </div>
        )}

        {/* 数据表格 */}
        <List
          data={emailHistory}
          columns={columns}
          renderCell={renderCell}
          onRowClick={onEmailClick}
          loading={loading}
          empty={emailHistory.length === 0}
          emptyContent={<EmptyState title={t('admin.smtp.history.empty')} />}
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

export default AdminEmailHistory;
