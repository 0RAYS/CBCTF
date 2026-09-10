import { useEffect, useEffectEvent, useReducer, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getLiveTasks } from '../../../../api/admin/task';
import { toast } from '../../../../utils/toast';
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
  const [refreshInterval, setRefreshInterval] = useState(10);
  const [revision, refresh] = useReducer((value) => value + 1, 0);
  const pending = useRef(false);
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
        const response = await getLiveTasks(taskQueryParams(query));
        if (current && response.code === 200) setData(normalizeTaskPayload(response.data));
      } catch (error) {
        if (current) reportError(error);
      } finally {
        if (current) pending.current = false;
      }
    });
    return () => {
      current = false;
      pending.current = false;
    };
  }, [active, query, revision]);

  useEffect(() => {
    const delay = livePollingDelay(active, refreshInterval);
    if (delay === null) return;
    const timer = setInterval(() => {
      if (!pending.current) refresh();
    }, delay);
    return () => clearInterval(timer);
  }, [active, refreshInterval]);

  return { ...data, query, dispatch, refresh, refreshInterval, setRefreshInterval };
}
