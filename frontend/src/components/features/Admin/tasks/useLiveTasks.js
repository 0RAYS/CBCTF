import { useEffect, useEffectEvent, useReducer, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getLiveTasks } from '../../../../api/admin/task';
import { toast } from '../../../../utils/toast';
import { useDebounce } from '../../../../hooks/useDebounce';
import { startVisiblePolling } from '../../../../utils/visiblePolling.js';
import {
  createTaskQuery,
  livePollingDelay,
  normalizeTaskPayload,
  taskQueryParams,
  taskQueryReducer,
} from './taskModel';

export default function useLiveTasks(active) {
  const { t } = useTranslation();
  const [query, dispatch] = useReducer(taskQueryReducer, undefined, createTaskQuery);
  const [data, setData] = useState(() => normalizeTaskPayload());
  const [loadedQuery, setLoadedQuery] = useState(null);
  const [refreshInterval, setRefreshInterval] = useState(10);
  const [revision, refresh] = useReducer((value) => value + 1, 0);
  const pending = useRef(false);
  // 搜索词防抖后才发起请求，输入过程中不打断已有结果
  const requestQuery = useDebounce(query, 300);
  // 查询条件与已展示数据的条件不一致时才算加载中，轮询刷新不会清空表格
  const loading = active && loadedQuery !== requestQuery;
  const reportError = useEffectEvent((error) => {
    toast.danger({ description: error.message || t('admin.tasks.toast.fetchLiveFailed') });
  });

  useEffect(() => {
    if (!active) return;
    let current = true;
    // Query changes and refreshes have one request path, including the first mount.
    Promise.resolve().then(async () => {
      if (!current) return;
      pending.current = true;
      try {
        const response = await getLiveTasks(taskQueryParams(requestQuery));
        if (current && response.code === 200) setData(normalizeTaskPayload(response.data));
      } catch (error) {
        if (current) reportError(error);
      } finally {
        if (current) {
          pending.current = false;
          setLoadedQuery(requestQuery);
        }
      }
    });
    return () => {
      current = false;
      pending.current = false;
    };
  }, [active, requestQuery, revision]);

  useEffect(() => {
    const delay = livePollingDelay(active, refreshInterval);
    if (delay === null) return;
    return startVisiblePolling(() => {
      if (!pending.current) refresh();
    }, delay);
  }, [active, refreshInterval]);

  return { ...data, query, dispatch, refresh, refreshInterval, setRefreshInterval, loading };
}
