import { useTranslation } from 'react-i18next';
import { IconRefresh, IconSearch, IconTimeline, IconTopologyStar3 } from '@tabler/icons-react';
import { Button, Card, Input, Select } from '../../../common';
import AutoRefreshControl from '../../../common/AutoRefreshControl';
import { TASK_STATUSES } from './taskModel';

export default function TaskFilters({ mode, state }) {
  const { t } = useTranslation();
  const { query, dispatch, queues, refresh, refreshInterval, setRefreshInterval } = state;
  const onFilterChange = (key, value) => dispatch({ type: 'filter', key, value });

  return (
    <Card className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2 text-neutral-300 font-mono">
          {mode === 'history' ? <IconTimeline size={18} /> : <IconTopologyStar3 size={18} />}
          <span>{t(`admin.tasks.${mode}.title`)}</span>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          {mode === 'live' && <AutoRefreshControl value={refreshInterval} onChange={setRefreshInterval} />}
          <Button variant="ghost" size="sm" icon={<IconRefresh size={14} />} onClick={refresh}>
            {t('common.refresh')}
          </Button>
          <Button variant="ghost" size="sm" onClick={() => dispatch({ type: 'reset' })}>
            {t('admin.tasks.filters.reset')}
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
        <Input
          type="search"
          value={query.filters.task_id}
          onChange={(e) => onFilterChange('task_id', e.target.value)}
          placeholder={t('admin.tasks.filters.taskId')}
          icon={<IconSearch size={14} />}
        />
        <Select
          value={query.filters.status}
          onChange={(e) => onFilterChange('status', e.target.value)}
          options={TASK_STATUSES[mode].map((value) => ({
            value,
            label: value ? t(`admin.tasks.status.${value}`) : t('admin.tasks.filters.allStatus'),
          }))}
        />
        <Select
          value={query.filters.queue}
          onChange={(e) => onFilterChange('queue', e.target.value)}
          options={[
            { value: '', label: t('admin.tasks.filters.allQueues') },
            ...queues.map((value) => ({ value, label: value })),
          ]}
        />
      </div>
    </Card>
  );
}
