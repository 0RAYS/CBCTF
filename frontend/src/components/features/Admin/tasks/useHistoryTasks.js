import { useEffect, useEffectEvent, useReducer, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getTaskHistory } from '@/api/admin/task';
import { toast } from '@/utils/toast';
import { useDebounce } from '@/hooks/useDebounce';
import {
  createTaskQuery,
  normalizeTaskPayload,
  taskQueryParams,
  taskQueryReducer,
} from '@/components/features/Admin/tasks/taskModel';

export default function useHistoryTasks(active) {
  const { t } = useTranslation();
  const [query, dispatch] = useReducer(taskQueryReducer, undefined, createTaskQuery);
  const [data, setData] = useState(() => normalizeTaskPayload());
  const [loadedQuery, setLoadedQuery] = useState(null);
  const [revision, refresh] = useReducer((value) => value + 1, 0);
  // 输入过程中不打断已有结果
  const requestQuery = useDebounce(query, 300);
  // 刷新不会清空表格
  const loading = active && loadedQuery !== requestQuery;
  const reportError = useEffectEvent((error) => {
    toast.danger({ description: error.message || t('admin.tasks.toast.fetchHistoryFailed') });
  });

  useEffect(() => {
    if (!active) return;
    let current = true;
    // Defer one microtask so StrictMode's discarded effect never starts a request.
    Promise.resolve().then(async () => {
      if (!current) return;
      try {
        const response = await getTaskHistory(taskQueryParams(requestQuery));
        if (current && response.code === 200) setData(normalizeTaskPayload(response.data));
      } catch (error) {
        if (current) reportError(error);
      } finally {
        if (current) setLoadedQuery(requestQuery);
      }
    });
    return () => {
      current = false;
    };
  }, [active, requestQuery, revision]);

  return { ...data, query, dispatch, refresh, loading };
}
