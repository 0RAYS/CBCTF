import { useTranslation } from 'react-i18next';
import { IconRefresh, IconServerBolt } from '@tabler/icons-react';
import { Button, Card, EmptyState } from '@/components/common';
import TruncatedText from '@/components/common/TruncatedText';
import VictimStatusBadge from '@/components/features/Admin/victims/VictimStatusBadge';
import useRangePools from '@/components/features/Admin/victims/useRangePools';

// RangePoolPanel 展示靶场实例池的常驻数量与每个实例的运行状态。
export default function RangePoolPanel({ contestId }) {
  const { t, i18n } = useTranslation();
  const { pools, loading, prewarming, prewarm, refresh } = useRangePools(contestId);

  const formatTime = (value) =>
    value && !String(value).startsWith('0001-')
      ? new Date(value).toLocaleString(i18n.language || 'en-US', {
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
        })
      : '-';

  return (
    <Card variant="default" padding="none" className="overflow-hidden">
      <div className="flex flex-wrap items-center justify-between gap-3 p-4 bg-black/20 border-b border-neutral-300/30">
        <div className="flex items-center gap-2">
          <IconServerBolt size={20} className="text-amber-400" />
          <h3 className="text-lg font-mono text-neutral-50">{t('admin.victims.rangePool.title')}</h3>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="sm" icon={<IconRefresh size={14} />} onClick={refresh}>
            {t('common.refresh')}
          </Button>
          <Button variant="primary" size="sm" loading={prewarming} onClick={prewarm}>
            {t('admin.victims.rangePool.prewarm')}
          </Button>
        </div>
      </div>
      {loading ? (
        <EmptyState title={t('admin.victims.rangePool.loading')} />
      ) : pools.length === 0 ? (
        <EmptyState title={t('admin.victims.rangePool.empty')} description={t('admin.victims.rangePool.emptyHint')} />
      ) : (
        <div className="divide-y divide-neutral-300/10">
          {pools.map((pool) => (
            <div key={pool.challenge_id} className="p-4 space-y-3">
              <div className="flex flex-wrap items-center gap-3 text-sm font-mono">
                <span className="text-neutral-100">{pool.challenge_name}</span>
                <span className="text-neutral-400">
                  {t('admin.victims.rangePool.summary', {
                    size: pool.pool_size,
                    leased: pool.leased,
                    max: pool.max_leased,
                  })}
                </span>
              </div>
              <div className="table-scroll">
                <table className="table-adaptive" style={{ '--table-min-width': '640px' }}>
                  <thead>
                    <tr className="bg-black/40">
                      <th className="p-3 text-left text-neutral-400 font-mono">
                        {t('admin.victims.rangePool.columns.id')}
                      </th>
                      <th className="p-3 text-left text-neutral-400 font-mono">
                        {t('admin.victims.rangePool.columns.status')}
                      </th>
                      <th className="p-3 text-left text-neutral-400 font-mono">
                        {t('admin.victims.rangePool.columns.team')}
                      </th>
                      <th className="p-3 text-left text-neutral-400 font-mono">
                        {t('admin.victims.rangePool.columns.nodes')}
                      </th>
                      <th className="p-3 text-left text-neutral-400 font-mono">
                        {t('admin.victims.rangePool.columns.leaseExpireAt')}
                      </th>
                      <th className="p-3 text-left text-neutral-400 font-mono">
                        {t('admin.victims.rangePool.columns.templateKey')}
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {pool.instances.length === 0 ? (
                      <tr>
                        <td colSpan={6} className="p-4 text-neutral-400 font-mono">
                          {t('admin.victims.rangePool.noInstances')}
                        </td>
                      </tr>
                    ) : (
                      pool.instances.map((instance) => (
                        <tr key={instance.id} className="border-t border-neutral-300/10">
                          <td className="p-3 text-neutral-300 font-mono">{instance.id}</td>
                          <td className="p-3 font-mono">
                            <VictimStatusBadge status={instance.status} t={t} translationKey="admin.victims" />
                          </td>
                          <td className="p-3 text-neutral-300 font-mono">
                            <TruncatedText>{instance.team_name || '-'}</TruncatedText>
                          </td>
                          <td className="p-3 text-neutral-300 font-mono">{instance.nodes}</td>
                          <td className="p-3 text-neutral-300 font-mono whitespace-nowrap">
                            {instance.team_id ? formatTime(instance.lease_expire_at) : '-'}
                          </td>
                          <td className="p-3 text-neutral-400 font-mono">
                            <TruncatedText>{instance.template_key}</TruncatedText>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          ))}
        </div>
      )}
    </Card>
  );
}
