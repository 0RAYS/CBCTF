import { useTranslation } from 'react-i18next';
import { List } from '../../../common';
import { formatDateTime } from './cheatFormat';

// actorLabel 展示事件发起方的身份，用户与队伍分开标注。
function actorLabel(item) {
  const prefix = item.actor_kind === 'Team' ? 'Team' : 'User';
  return `${prefix}-${item.actor_id}`;
}

// eventDetail 展示事件的具体内容，不同来源展示不同字段。
function eventDetail(item) {
  if (item.source === 'web_request') {
    return [item.method, item.path, item.status ? `HTTP ${item.status}` : null].filter(Boolean).join(' ');
  }
  if (item.source === 'flag_submission') return `${item.value || ''} ${item.detail || ''}`.trim();
  if (item.source === 'flag_solve') return `${item.value || ''} ${item.detail || ''}`.trim();
  return item.detail || '-';
}

// CheatEvidenceTimeline 按时间顺序展示判定背后的观测事件。
export default function CheatEvidenceTimeline({ timeline, language }) {
  const { t } = useTranslation();
  if (!timeline || timeline.length === 0) {
    return <div className="text-neutral-400 text-sm">{t('admin.contests.cheats.evidence.empty')}</div>;
  }
  const columns = [
    ['at', t('admin.contests.cheats.evidence.at'), '22%'],
    ['source', t('admin.contests.cheats.evidence.source'), '16%'],
    ['actor', t('admin.contests.cheats.evidence.actor'), '14%'],
    ['target', t('admin.contests.cheats.evidence.target'), '20%'],
    ['detail', t('admin.contests.cheats.evidence.detail'), '28%'],
  ].map(([key, label, width]) => ({ key, label, width }));

  function renderCell(item, column) {
    switch (column.key) {
      case 'at':
        return <span className="text-xs">{formatDateTime(item.at, language)}</span>;
      case 'source':
        return (
          <span className="text-xs font-mono text-neutral-300">
            {t(`admin.contests.cheats.evidence.sources.${item.source}`, item.source)}
          </span>
        );
      case 'actor':
        return <span className="text-xs font-mono text-geek-400">{actorLabel(item)}</span>;
      case 'target':
        return <span className="text-xs">{item.target || '-'}</span>;
      default:
        return <span className="text-xs break-all">{eventDetail(item)}</span>;
    }
  }

  return <List columns={columns} data={timeline} renderCell={renderCell} minWidth={720} animate={false} />;
}
