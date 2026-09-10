import { useTranslation } from 'react-i18next';
import { IconDownload, IconGraph } from '@tabler/icons-react';
import { Button, Card, List, Pagination, Spinner, StatusTag } from '../../../../common';
import { containerStatus, TEAM_DETAIL_PAGE_SIZE } from './teamDetailData.js';

export default function TrafficPanel({
  containers = [],
  count = 0,
  page = 1,
  loading,
  onPageChange,
  onUserClick,
  onDownload,
  onViewGraph,
}) {
  const { t, i18n } = useTranslation();
  const columns = [
    ['id', 'containerId'],
    ['user_id', 'userId'],
    ['contest_challenge_name', 'challengeName'],
    ['start', 'startTime'],
    ['duration', 'duration'],
    ['status', 'status'],
    ['actions', 'actions'],
  ].map(([key, label]) => ({ key, label: t(`admin.contests.teamDetail.traffic.columns.${label}`) }));
  const durationLabel = (duration) => {
    const seconds = Number(duration) || 0;
    const unit = (name, count) => t(`utils.time.units.${name}`, { count });
    if (seconds < 60) return unit('second', Math.round(seconds));
    if (seconds < 3600) return unit('minute', Math.floor(seconds / 60)) + unit('second', Math.round(seconds % 60));
    const minutes = Math.floor((seconds % 3600) / 60);
    return unit('hour', Math.floor(seconds / 3600)) + (minutes ? unit('minute', minutes) : '');
  };
  const renderCell = (container, { key }) => {
    if (key === 'start') return new Date(container.start).toLocaleString(i18n.language);
    if (key === 'duration') return durationLabel(container.duration);
    if (key === 'status') {
      const status = containerStatus(container.start, container.duration);
      return (
        <StatusTag
          type={{ upcoming: 'info', ended: 'error', running: 'success' }[status]}
          text={t(`admin.contests.teamDetail.traffic.status.${status}`)}
        />
      );
    }
    if (key === 'user_id' && onUserClick)
      return (
        <Button variant="ghost" size="sm" onClick={() => onUserClick(container.user_id)}>
          {container.user_id}
        </Button>
      );
    if (key === 'actions')
      return (
        <div className="flex gap-2">
          {onViewGraph && (
            <Button
              variant="ghost"
              size="icon"
              onClick={() => onViewGraph(container)}
              title={t('admin.contests.teamDetail.traffic.actions.viewTraffic')}
            >
              <IconGraph size={18} />
            </Button>
          )}
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onDownload(container)}
            title={t('admin.contests.teamDetail.traffic.actions.downloadTraffic')}
          >
            <IconDownload size={18} />
          </Button>
        </div>
      );
    return container[key] ?? '-';
  };
  return (
    <section>
      <h2 className="text-xl font-mono text-neutral-50 mb-4">{t('admin.contests.teamDetail.sections.traffic')}</h2>
      {loading ? (
        <Card className="flex justify-center p-8">
          <Spinner />
        </Card>
      ) : (
        <List
          className="[&_table]:min-w-[900px]"
          columns={columns}
          data={containers}
          renderCell={renderCell}
          empty={!containers.length}
          emptyContent={t('admin.contests.teamDetail.empty.traffic')}
          footer={
            <Pagination
              total={Math.ceil(count / TEAM_DETAIL_PAGE_SIZE)}
              current={page}
              pageSize={TEAM_DETAIL_PAGE_SIZE}
              onChange={onPageChange}
              showTotal
              totalItems={count}
            />
          }
        />
      )}
    </section>
  );
}
