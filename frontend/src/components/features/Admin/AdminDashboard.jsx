/** 管理员仪表盘组件 */

import { motion } from 'motion/react';
import { Card, StatCard } from '@/components/common';
import { useTranslation } from 'react-i18next';

function AdminDashboard({ status, chartContent, extraContent }) {
  const { t } = useTranslation();

  const formatBytes = (bytes) => {
    if (!bytes && bytes !== 0) return t('common.notAvailable');
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="w-full mx-auto space-y-6">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        <div className="flex items-center justify-end mb-6">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-green-400 animate-pulse"></div>
            <span className="text-neutral-400 text-sm">{t('admin.dashboard.realtime')}</span>
          </div>
        </div>
        {status?.unavailable?.length || status?.metrics_skipped ? (
          <p role="status" className="mb-4 text-sm text-amber-300">
            {t('admin.dashboard.partialStats')}
          </p>
        ) : null}

        {chartContent && <div className="mb-5 border border-neutral-600 rounded-md bg-neutral-900">{chartContent}</div>}

        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <StatCard title={t('admin.dashboard.cards.requestIp')} value={status?.ip ?? t('common.notAvailable')} />
          <StatCard title={t('admin.dashboard.cards.users')} value={status?.users ?? t('common.notAvailable')} />
          <StatCard title={t('admin.dashboard.cards.contests')} value={status?.contests ?? t('common.notAvailable')} />
          <StatCard
            title={t('admin.dashboard.cards.challenges')}
            value={status?.challenges ?? t('common.notAvailable')}
          />
          <StatCard title={t('admin.dashboard.cards.victims')} value={status?.victims ?? t('common.notAvailable')} />
          <StatCard
            title={t('admin.dashboard.cards.submissions')}
            value={status?.submissions ?? t('common.notAvailable')}
          />
          <StatCard title={t('admin.dashboard.cards.requests')} value={status?.requests ?? t('common.notAvailable')} />
          <StatCard
            title={t('admin.dashboard.cards.responseTime')}
            value={status?.duration ? `${status.duration} ms` : t('common.notAvailable')}
          />
          <StatCard title={t('admin.dashboard.cards.downlink')} value={formatBytes(status?.sent)} />
          <StatCard title={t('admin.dashboard.cards.uplink')} value={formatBytes(status?.recv)} />
          <StatCard title={t('admin.dashboard.cards.totalTraffic')} value={formatBytes(status?.io)} />
          <StatCard title={t('admin.dashboard.cards.cacheSize')} value={status?.cache ?? t('common.notAvailable')} />
        </div>
      </motion.div>

      {extraContent && (
        <Card variant="default" padding="md" animate>
          {extraContent}
        </Card>
      )}
    </div>
  );
}

export default AdminDashboard;
