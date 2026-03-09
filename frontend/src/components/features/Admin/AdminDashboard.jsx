/**
 * 管理员仪表盘组件
 * @param {Object} props
 * @param {Object} props.status - 系统状态数据
 * @param {string} props.status.ip - 服务器IP
 * @param {number} props.status.users - 用户数量
 * @param {number} props.status.contests - 赛事数量
 * @param {number} props.status.challenges - 题目数量
 * @param {number} props.status.requests - 请求数量
 * @param {number} props.status.duration - 平均响应时间(ms)
 * @param {number} props.status.sent - 下行流量(bytes)
 * @param {number} props.status.recv - 上行流量(bytes)
 * @param {number} props.status.io - 总流量(bytes)
 * @param {number} props.status.cache - 缓存总量
 * @param {number} props.status.victims - 靶机数量
 * @param {number} props.status.submissions - 提交次数
 * @param {React.ReactNode} props.chartContent - 图表内容
 * @param {React.ReactNode} props.extraContent - 额外内容，如表格等
 */

import { motion } from 'motion/react';
import { Card, StatCard } from '../../common';
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

        {/* 图表区域 - 由外部传入 */}
        {chartContent && <div className="mb-5 border border-neutral-600 rounded-md bg-neutral-900">{chartContent}</div>}

        {/* 状态卡片网格 */}
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <StatCard title={t('admin.dashboard.cards.serverIp')} value={status?.ip || t('common.notAvailable')} />
          <StatCard title={t('admin.dashboard.cards.users')} value={status?.users || 0} />
          <StatCard title={t('admin.dashboard.cards.contests')} value={status?.contests || 0} />
          <StatCard title={t('admin.dashboard.cards.challenges')} value={status?.challenges || 0} />
          <StatCard title={t('admin.dashboard.cards.victims')} value={status?.victims || 0} />
          <StatCard title={t('admin.dashboard.cards.submissions')} value={status?.submissions || 0} />
          <StatCard title={t('admin.dashboard.cards.requests')} value={status?.requests || 0} />
          <StatCard
            title={t('admin.dashboard.cards.responseTime')}
            value={status?.duration ? `${status.duration} ms` : t('common.notAvailable')}
          />
          <StatCard title={t('admin.dashboard.cards.downlink')} value={formatBytes(status?.sent)} />
          <StatCard title={t('admin.dashboard.cards.uplink')} value={formatBytes(status?.recv)} />
          <StatCard title={t('admin.dashboard.cards.totalTraffic')} value={formatBytes(status?.io)} />
          <StatCard title={t('admin.dashboard.cards.cacheSize')} value={status?.cache || 0} />
        </div>
      </motion.div>

      {/* 额外内容区域 - 由外部传入 */}
      {extraContent && (
        <Card variant="default" padding="md" animate>
          {extraContent}
        </Card>
      )}
    </div>
  );
}

export default AdminDashboard;
