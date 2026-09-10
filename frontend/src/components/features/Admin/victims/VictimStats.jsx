import { IconBan, IconPlayerPlay, IconServer } from '@tabler/icons-react';
import { StatCard } from '../../../common';

export function VictimStats({ stats, t, translationKey }) {
  return (
    <div className="mb-8">
      <div className="mb-4">
        <p className="text-neutral-400 font-mono">{t(`${translationKey}.page.subtitle`)}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <StatCard
          title={t(`${translationKey}.stats.total`)}
          value={stats.totalContainers}
          valueColor="text-neutral-50"
          icon={<IconServer size={20} className="text-geek-400" />}
        />
        <StatCard
          title={t(`${translationKey}.stats.running`)}
          value={stats.runningContainers}
          valueColor="text-green-400"
          icon={<IconPlayerPlay size={20} className="text-green-400" />}
          iconBgClass="bg-green-400/20"
          delay={0.1}
        />
        <StatCard
          title={t(`${translationKey}.stats.stopped`)}
          value={stats.stoppedContainers}
          valueColor="text-red-400"
          icon={<IconBan size={20} className="text-red-400" />}
          iconBgClass="bg-red-400/20"
          delay={0.2}
        />
      </div>
    </div>
  );
}
