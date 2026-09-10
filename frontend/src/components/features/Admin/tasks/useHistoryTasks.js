import { useEffect, useEffectEvent, useReducer, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getTaskHistory } from '../../../../api/admin/task';
import { toast } from '../../../../utils/toast';
import { createTaskQuery, normalizeTaskPayload, taskQueryParams, taskQueryReducer } from './taskModel';

export default function useHistoryTasks(active) {
  const { t } = useTranslation();
  const [query, dispatch] = useReducer(taskQueryReducer, undefined, createTaskQuery);
  const [data, setData] = useState(() => normalizeTaskPayload());
  const [revision, refresh] = useReducer((value) => value + 1, 0);
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
        const response = await getTaskHistory(taskQueryParams(query));
        if (current && response.code === 200) setData(normalizeTaskPayload(response.data));
      } catch (error) {
        if (current) reportError(error);
      }
    });
    return () => {
      current = false;
    };
  }, [active, query, revision]);

  return { ...data, query, dispatch, refresh };
}
