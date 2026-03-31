import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconClockPlay, IconRefresh, IconSearch, IconTimeline, IconTopologyStar3 } from '@tabler/icons-react';
import { Card, EmptyState, Input, Pagination, Select, Tabs, Button } from '../../components/common';
import { getLiveTasks, getTaskHistory } from '../../api/admin/task';
import { toast } from '../../utils/toast';

const PAGE_SIZE = 20;

const HISTORY_STATUS_OPTIONS = [
  { value: '', label: 'All' },
  { value: 'success', label: 'Success' },
  { value: 'failed', label: 'Failed' },
];

const LIVE_STATUS_OPTIONS = [
  { value: 'active', label: 'Active' },
  { value: 'pending', label: 'Pending' },
  { value: 'scheduled', label: 'Scheduled' },
  { value: 'retry', label: 'Retry' },
  { value: 'archived', label: 'Archived' },
  { value: 'completed', label: 'Completed' },
];

function formatDateTime(value, language) {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleString(language || 'en-US', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}

function formatPayload(value) {
  if (value === null || value === undefined || value === '') return '-';
  if (typeof value === 'string') return value;
  try {
    return JSON.stringify(value);
  } catch {
    return String(value);
  }
}

function StatusBadge({ value }) {
  const style =
    {
      success: 'bg-green-400/10 text-green-400 border-green-400/30',
      failed: 'bg-red-400/10 text-red-400 border-red-400/30',
      active: 'bg-geek-400/10 text-geek-400 border-geek-400/30',
      pending: 'bg-yellow-400/10 text-yellow-400 border-yellow-400/30',
      scheduled: 'bg-blue-400/10 text-blue-400 border-blue-400/30',
      retry: 'bg-orange-400/10 text-orange-300 border-orange-400/30',
      archived: 'bg-red-400/10 text-red-400 border-red-400/30',
      completed: 'bg-green-400/10 text-green-400 border-green-400/30',
    }[value] || 'bg-neutral-500/10 text-neutral-300 border-neutral-500/30';

  return <span className={`inline-block px-2 py-1 rounded border text-xs font-mono ${style}`}>{value || '-'}</span>;
}

function FilterBar({
  t,
  mode,
  filters,
  onFilterChange,
  onReset,
  queueOptions,
  typeOptions,
  refreshInterval,
  onRefreshIntervalChange,
  onRefresh,
}) {
  const statusOptions = mode === 'history' ? HISTORY_STATUS_OPTIONS : LIVE_STATUS_OPTIONS;
  const isHistory = mode === 'history';

  return (
    <Card className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2 text-neutral-300 font-mono">
          {mode === 'history' ? <IconTimeline size={18} /> : <IconTopologyStar3 size={18} />}
          <span>{mode === 'history' ? t('admin.tasks.history.title') : t('admin.tasks.live.title')}</span>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          {mode === 'live' && (
            <div className="flex items-center gap-1 px-2 h-8 rounded-md border border-neutral-700 bg-neutral-900">
              <IconClockPlay size={13} className="text-neutral-400 shrink-0" />
              <span className="text-xs text-neutral-400 shrink-0">{t('common.autoRefresh')}</span>
              <select
                value={refreshInterval}
                onChange={(e) => onRefreshIntervalChange(Number(e.target.value))}
                className="bg-transparent text-xs text-neutral-300 outline-none cursor-pointer"
              >
                {[5, 10, 30, 60].map((s) => (
                  <option key={s} value={s} className="bg-neutral-900">
                    {s}s
                  </option>
                ))}
                <option value={0} className="bg-neutral-900">
                  {t('common.autoRefreshOff')}
                </option>
              </select>
            </div>
          )}
          <Button variant="ghost" size="sm" icon={<IconRefresh size={14} />} onClick={onRefresh}>
            {t('common.refresh')}
          </Button>
          <Button variant="ghost" size="sm" onClick={onReset}>
            {t('admin.tasks.filters.reset')}
          </Button>
        </div>
      </div>

      <div className={`grid grid-cols-1 ${isHistory ? 'md:grid-cols-4' : 'md:grid-cols-2'} gap-3`}>
        {isHistory && (
          <Input
            type="search"
            value={filters.task_id}
            onChange={(e) => onFilterChange('task_id', e.target.value)}
            placeholder={t('admin.tasks.filters.taskId')}
            icon={<IconSearch size={14} />}
          />
        )}
        <Select
          value={filters.status}
          onChange={(e) => onFilterChange('status', e.target.value)}
          options={statusOptions.map((item) => ({
            value: item.value,
            label: item.value ? t(`admin.tasks.status.${item.value}`) : t('admin.tasks.filters.allStatus'),
          }))}
        />
        <Select
          value={filters.queue}
          onChange={(e) => onFilterChange('queue', e.target.value)}
          options={[
            { value: '', label: t('admin.tasks.filters.allQueues') },
            ...queueOptions.map((item) => ({ value: item, label: item })),
          ]}
        />
        {isHistory && (
          <Select
            value={filters.type}
            onChange={(e) => onFilterChange('type', e.target.value)}
            options={[
              { value: '', label: t('admin.tasks.filters.allTypes') },
              ...typeOptions.map((item) => ({ value: item, label: item })),
            ]}
          />
        )}
      </div>
    </Card>
  );
}

function HistoryTable({ rows, totalCount, currentPage, onPageChange, language, t }) {
  return (
    <Card padding="none" className="overflow-hidden">
      {rows.length === 0 ? (
        <EmptyState title={t('admin.tasks.history.empty')} className="py-20" />
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-sm text-neutral-300">
            <thead>
              <tr className="bg-black/30 border-b border-neutral-300/10">
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.taskId')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.type')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.queue')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.status')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.retry')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.processedAt')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.error')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.payload')}</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((item) => (
                <tr key={`${item.id}-${item.task_id}`} className="border-b border-neutral-300/10 align-top">
                  <td className="p-4 font-mono text-xs text-neutral-200 whitespace-nowrap">{item.task_id || '-'}</td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">{item.type}</td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">{item.queue}</td>
                  <td className="p-4 whitespace-nowrap">
                    <StatusBadge value={item.status} />
                  </td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">
                    {item.retry_count}/{item.max_retry}
                  </td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">
                    {formatDateTime(item.processed_at, language)}
                  </td>
                  <td className="p-4 max-w-80">
                    <div className="line-clamp-3 break-words text-xs text-red-300">{item.error || '-'}</div>
                  </td>
                  <td className="p-4 max-w-96">
                    <div className="line-clamp-3 break-all font-mono text-xs text-neutral-400">
                      {formatPayload(item.payload)}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {totalCount > PAGE_SIZE && (
        <div className="p-4 border-t border-neutral-300/10 bg-black/20 flex justify-center">
          <Pagination
            current={currentPage}
            total={Math.ceil(totalCount / PAGE_SIZE)}
            totalItems={totalCount}
            showTotal
            onChange={onPageChange}
          />
        </div>
      )}
    </Card>
  );
}

function LiveTable({ rows, totalCount, currentPage, onPageChange, language, t }) {
  return (
    <Card padding="none" className="overflow-hidden">
      {rows.length === 0 ? (
        <EmptyState title={t('admin.tasks.live.empty')} className="py-20" />
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-sm text-neutral-300">
            <thead>
              <tr className="bg-black/30 border-b border-neutral-300/10">
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.taskId')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.type')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.queue')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.status')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.retry')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.nextProcessAt')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.error')}</th>
                <th className="p-4 text-left font-mono text-neutral-400">{t('admin.tasks.columns.payload')}</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((item) => (
                <tr key={`${item.queue}-${item.task_id}`} className="border-b border-neutral-300/10 align-top">
                  <td className="p-4 font-mono text-xs text-neutral-200 whitespace-nowrap">{item.task_id || '-'}</td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">{item.type}</td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">{item.queue}</td>
                  <td className="p-4 whitespace-nowrap">
                    <StatusBadge value={item.status} />
                  </td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">
                    {item.retry_count}/{item.max_retry}
                  </td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">
                    {formatDateTime(item.next_process_at || item.completed_at || item.last_failed_at, language)}
                  </td>
                  <td className="p-4 max-w-80">
                    <div className="line-clamp-3 break-words text-xs text-red-300">{item.error || '-'}</div>
                  </td>
                  <td className="p-4 max-w-96">
                    <div className="line-clamp-3 break-all font-mono text-xs text-neutral-400">
                      {formatPayload(item.payload)}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {totalCount > PAGE_SIZE && (
        <div className="p-4 border-t border-neutral-300/10 bg-black/20 flex justify-center">
          <Pagination
            current={currentPage}
            total={Math.ceil(totalCount / PAGE_SIZE)}
            totalItems={totalCount}
            showTotal
            onChange={onPageChange}
          />
        </div>
      )}
    </Card>
  );
}

function TasksPage() {
  const { t, i18n } = useTranslation();
  const [tab, setTab] = useState('history');

  const [historyRows, setHistoryRows] = useState([]);
  const [historyTotal, setHistoryTotal] = useState(0);
  const [historyPage, setHistoryPage] = useState(1);
  const [historyFilters, setHistoryFilters] = useState({ task_id: '', type: '', queue: '', status: '' });
  const [historyQueues, setHistoryQueues] = useState([]);
  const [historyTypes, setHistoryTypes] = useState([]);

  const [liveRows, setLiveRows] = useState([]);
  const [liveTotal, setLiveTotal] = useState(0);
  const [livePage, setLivePage] = useState(1);
  const [liveFilters, setLiveFilters] = useState({ queue: '', status: 'active' });
  const [liveQueues, setLiveQueues] = useState([]);
  const [refreshInterval, setRefreshInterval] = useState(10);

  const livePageRef = useRef(livePage);
  const liveFiltersRef = useRef(liveFilters);

  useEffect(() => {
    livePageRef.current = livePage;
  }, [livePage]);

  useEffect(() => {
    liveFiltersRef.current = liveFilters;
  }, [liveFilters]);

  const fetchHistory = async (page = historyPage, filters = historyFilters) => {
    try {
      const res = await getTaskHistory({
        ...filters,
        limit: PAGE_SIZE,
        offset: (page - 1) * PAGE_SIZE,
      });
      if (res.code === 200) {
        const payload = res.data || {};
        const tasks = Array.isArray(payload.tasks) ? payload.tasks : [];
        setHistoryRows(tasks);
        setHistoryTotal(payload.count || 0);
        setHistoryQueues(Array.isArray(payload.queues) ? payload.queues : []);
        setHistoryTypes(Array.isArray(payload.types) ? payload.types : []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.tasks.toast.fetchHistoryFailed') });
    }
  };

  const fetchLive = async (page = livePage, filters = liveFilters) => {
    try {
      const res = await getLiveTasks({
        ...filters,
        limit: PAGE_SIZE,
        offset: (page - 1) * PAGE_SIZE,
      });
      if (res.code === 200) {
        const payload = res.data || {};
        setLiveRows(Array.isArray(payload.tasks) ? payload.tasks : []);
        setLiveTotal(payload.count || 0);
        setLiveQueues(Array.isArray(payload.queues) ? payload.queues : []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.tasks.toast.fetchLiveFailed') });
    }
  };

  useEffect(() => {
    fetchHistory(1, historyFilters);
  }, []);

  useEffect(() => {
    fetchLive(1, liveFilters);
  }, []);

  useEffect(() => {
    fetchHistory(historyPage, historyFilters);
  }, [historyPage]);

  useEffect(() => {
    fetchLive(livePage, liveFilters);
  }, [livePage]);

  useEffect(() => {
    if (refreshInterval <= 0) return;
    const timer = setInterval(() => {
      fetchLive(livePageRef.current, liveFiltersRef.current);
    }, refreshInterval * 1000);
    return () => clearInterval(timer);
  }, [refreshInterval]);

  const handleHistoryFilterChange = (key, value) => {
    const next = { ...historyFilters, [key]: value };
    setHistoryFilters(next);
    setHistoryPage(1);
    fetchHistory(1, next);
  };

  const handleLiveFilterChange = (key, value) => {
    const next = { ...liveFilters, [key]: value };
    setLiveFilters(next);
    setLivePage(1);
    fetchLive(1, next);
  };

  const resetHistory = () => {
    const next = { task_id: '', type: '', queue: '', status: '' };
    setHistoryFilters(next);
    setHistoryPage(1);
    fetchHistory(1, next);
  };

  const resetLive = () => {
    const next = { queue: '', status: 'active' };
    setLiveFilters(next);
    setLivePage(1);
    fetchLive(1, next);
  };

  const tabItems = useMemo(
    () => [
      { key: 'history', label: t('admin.tasks.tabs.history') },
      { key: 'live', label: t('admin.tasks.tabs.live') },
    ],
    [t]
  );

  return (
    <div className="w-full mx-auto space-y-6">
      <div>
        <p className="text-neutral-400 font-mono">{t('admin.tasks.subtitle')}</p>
      </div>

      <Tabs items={tabItems} value={tab} onChange={setTab} variant="compact" wrapperClassName="w-full" />

      {tab === 'history' ? (
        <>
          <FilterBar
            t={t}
            mode="history"
            filters={historyFilters}
            onFilterChange={handleHistoryFilterChange}
            onReset={resetHistory}
            queueOptions={historyQueues}
            typeOptions={historyTypes}
            refreshInterval={0}
            onRefreshIntervalChange={() => {}}
            onRefresh={() => fetchHistory(historyPage, historyFilters)}
          />
          <HistoryTable
            rows={historyRows}
            totalCount={historyTotal}
            currentPage={historyPage}
            onPageChange={setHistoryPage}
            language={i18n.language}
            t={t}
          />
        </>
      ) : (
        <>
          <FilterBar
            t={t}
            mode="live"
            filters={liveFilters}
            onFilterChange={handleLiveFilterChange}
            onReset={resetLive}
            queueOptions={liveQueues}
            typeOptions={[]}
            refreshInterval={refreshInterval}
            onRefreshIntervalChange={setRefreshInterval}
            onRefresh={() => fetchLive(livePage, liveFilters)}
          />
          <LiveTable
            rows={liveRows}
            totalCount={liveTotal}
            currentPage={livePage}
            onPageChange={setLivePage}
            language={i18n.language}
            t={t}
          />
        </>
      )}
    </div>
  );
}

export default TasksPage;
