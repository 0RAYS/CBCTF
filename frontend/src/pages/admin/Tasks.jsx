import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Tabs } from '../../components/common';
import TaskFilters from '../../components/features/Admin/tasks/TaskFilters';
import TaskTable from '../../components/features/Admin/tasks/TaskTable';
import useHistoryTasks from '../../components/features/Admin/tasks/useHistoryTasks';
import useLiveTasks from '../../components/features/Admin/tasks/useLiveTasks';

export default function TasksPage() {
  const { t } = useTranslation();
  const [tab, setTab] = useState('history');
  // Keep both queries mounted so switching tabs preserves each filter and page.
  const history = useHistoryTasks(tab === 'history');
  const live = useLiveTasks(tab === 'live');
  const state = tab === 'history' ? history : live;

  return (
    <div className="w-full mx-auto space-y-6">
      <div>
        <p className="text-neutral-400 font-mono">{t('admin.tasks.subtitle')}</p>
      </div>
      <Tabs
        items={[
          { key: 'history', label: t('admin.tasks.tabs.history') },
          { key: 'live', label: t('admin.tasks.tabs.live') },
        ]}
        value={tab}
        onChange={setTab}
        variant="compact"
        wrapperClassName="w-full"
      />
      <TaskFilters mode={tab} state={state} />
      <TaskTable
        mode={tab}
        rows={state.rows}
        totalCount={state.totalCount}
        currentPage={state.query.page}
        onPageChange={(page) => state.dispatch({ type: 'page', page })}
      />
    </div>
  );
}
