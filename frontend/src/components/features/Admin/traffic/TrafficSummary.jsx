import { IconActivity, IconArrowDownRight, IconArrowUpRight, IconRoute } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { formatBytes } from './trafficPresentation.js';

export default function TrafficSummary({ summary }) {
  const { t } = useTranslation();
  return (
    <div className="grid grid-cols-2 gap-2 xl:grid-cols-4">
      {[
        {
          icon: <IconActivity size={15} className="text-geek-400" />,
          label: t('admin.contests.trafficGraph.stats.totalTraffic'),
          value: formatBytes(summary.total_bytes),
          tone: 'text-geek-400',
        },
        {
          icon: <IconArrowDownRight size={15} className="text-neutral-300" />,
          label: t('admin.contests.trafficGraph.stats.ingress'),
          value: formatBytes(summary.ingress_bytes),
          tone: 'text-neutral-100',
        },
        {
          icon: <IconArrowUpRight size={15} className="text-neutral-300" />,
          label: t('admin.contests.trafficGraph.stats.egress'),
          value: formatBytes(summary.egress_bytes),
          tone: 'text-neutral-100',
        },
        {
          icon: <IconRoute size={15} className="text-neutral-300" />,
          label: t('admin.contests.trafficGraph.stats.processes'),
          value: summary.process_count || 0,
          tone: 'text-neutral-100',
        },
      ].map((item) => (
        <div key={item.label} className="rounded-xl border border-neutral-600 bg-black/20 px-3 py-2">
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-neutral-800">{item.icon}</div>
            <div className="min-w-0">
              <div className="text-[11px] font-mono text-neutral-500">{item.label}</div>
              <div className={`truncate font-mono text-base ${item.tone}`}>{item.value}</div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
