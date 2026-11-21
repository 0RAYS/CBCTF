import { useTranslation } from 'react-i18next';
import { Button, StatusTag } from '@/components/common';
import IpLookup from '@/components/features/Admin/network/IpLookup.jsx';
import { formatDateTime, formatGap } from '@/components/features/Admin/cheats/cheatFormat';

export function CheatModels({ model, openUserDetail, openTeamDetail }) {
  if (!model || Object.keys(model).length === 0) return '-';
  return (
    <div className="flex flex-wrap gap-1">
      {Object.entries(model).map(([key, values]) => {
        const ids = Array.isArray(values) ? values : [values];
        const open = key === 'User' ? openUserDetail : key === 'Team' ? openTeamDetail : null;
        if (!open)
          return (
            <span key={key} className="text-neutral-300">
              {key}: {ids.join(', ')}
            </span>
          );
        return ids.map((id) => (
          <Button
            key={`${key}-${id}`}
            size="sm"
            variant="ghost"
            className="text-geek-400! p-0! h-auto!"
            onClick={(event) => {
              event.stopPropagation();
              open(id);
            }}
          >
            {key}-{id};
          </Button>
        ));
      })}
    </div>
  );
}

export function CheatIp({ ip }) {
  const isIp = ip && (/^(\d{1,3}\.){3}\d{1,3}(\/\d+)?$/.test(ip) || ip.includes(':'));
  return isIp ? (
    <IpLookup ip={ip} className="text-xs" />
  ) : (
    <span className="font-mono text-xs text-neutral-300">{ip || '-'}</span>
  );
}

export function CheatVerdict({ type }) {
  const { t } = useTranslation();
  return (
    <StatusTag
      type={{ cheater: 'error', suspicious: 'warning', pass: 'default' }[type] || 'default'}
      text={t(`admin.contests.cheats.types.${type}`, type)}
    />
  );
}

export function CheatStatus({ checked }) {
  const { t } = useTranslation();
  return (
    <StatusTag
      type={checked ? 'success' : 'warning'}
      text={t(`admin.contests.cheats.status.${checked ? 'processed' : 'unprocessed'}`)}
    />
  );
}

// CheatEvidenceSummary 展示线索的证据摘要：判定对象、参与角色、事件规模与时间跨度。
export function CheatEvidenceSummary({ summary, language }) {
  const { t } = useTranslation();
  if (!summary || !summary.subject) return null;
  const items = [
    [
      t('admin.contests.cheats.evidence.actors'),
      t('admin.contests.cheats.evidence.actorCount', { count: summary.actors }),
    ],
    [t('admin.contests.cheats.evidence.events'), summary.events],
    [t('admin.contests.cheats.evidence.span'), formatGap(summary.span_seconds)],
    [t('admin.contests.cheats.evidence.minGap'), formatGap(summary.min_gap_seconds)],
  ];
  return (
    <div className="space-y-2">
      <div className="flex flex-wrap gap-2 text-xs text-neutral-300">
        {items.map(([label, value]) => (
          <span key={label} className="bg-neutral-800 rounded px-2 py-1">
            <span className="text-neutral-400">{label}</span>
            <span className="ml-1 font-mono">{value}</span>
          </span>
        ))}
      </div>
      <div className="text-xs text-neutral-400 font-mono">
        {t('admin.contests.cheats.evidence.window')}: {formatDateTime(summary.first_seen, language)} →{' '}
        {formatDateTime(summary.last_seen, language)}
      </div>
      {summary.truncated && (
        <div className="text-xs text-yellow-400">{t('admin.contests.cheats.evidence.truncated')}</div>
      )}
    </div>
  );
}

// CheatEvidenceActors 列出判定涉及的用户或队伍及其活跃区间。
export function CheatEvidenceActors({ actors, openUserDetail, openTeamDetail }) {
  const { t } = useTranslation();
  if (!actors || actors.length === 0) return null;
  return (
    <div className="space-y-2">
      <div className="text-sm font-mono text-neutral-400">{t('admin.contests.cheats.evidence.actorList')}</div>
      <ul className="space-y-1">
        {actors.map((actor) => {
          const open = actor.kind === 'User' ? openUserDetail : openTeamDetail;
          return (
            <li key={`${actor.kind}-${actor.id}`} className="flex flex-wrap items-center gap-2 text-xs">
              <Button
                size="sm"
                variant="ghost"
                className="text-geek-400! p-0! h-auto!"
                onClick={(event) => {
                  event.stopPropagation();
                  open?.(actor.id);
                }}
              >
                {actor.label || `${actor.kind}-${actor.id}`}
              </Button>
              {actor.team_ids?.length > 0 && (
                <span className="text-neutral-400 font-mono">
                  {t('admin.contests.cheats.evidence.team')}: {actor.team_ids.join(', ')}
                </span>
              )}
              <span className="text-neutral-400 font-mono">
                {t('admin.contests.cheats.evidence.count')}: {actor.count}
              </span>
              {actor.first_seen && (
                <span className="text-neutral-500 font-mono">
                  {formatDateTime(actor.first_seen)} → {formatDateTime(actor.last_seen)}
                </span>
              )}
            </li>
          );
        })}
      </ul>
    </div>
  );
}
