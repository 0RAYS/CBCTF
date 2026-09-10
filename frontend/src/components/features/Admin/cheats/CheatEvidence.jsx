import { useTranslation } from 'react-i18next';
import { Button, StatusTag } from '../../../common';
import IpLookup from '../network/IpLookup.jsx';

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
