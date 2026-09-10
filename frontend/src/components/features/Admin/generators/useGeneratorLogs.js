import { useEffect, useEffectEvent, useState } from 'react';
import { toast } from '../../../../utils/toast';
import { getGeneratorResponseData } from './generatorUtils.js';

export default function useGeneratorLogs(api, generatorId, text) {
  const [query, setQuery] = useState({ lines: 1000, delay: 0 });
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(true);
  const queryLogs = useEffectEvent((lines) => api.logs(generatorId, lines));
  const reportError = useEffectEvent(() => toast.danger({ description: text('toast.logFailed') }));

  useEffect(() => {
    let active = true;
    const fetchLogs = async () => {
      try {
        const response = await queryLogs(query.lines);
        if (active) setContent(getGeneratorResponseData(response).logs ?? '');
      } catch {
        if (active) reportError();
      } finally {
        if (active) setLoading(false);
      }
    };
    setContent('');
    setLoading(true);
    const timer = query.delay ? setTimeout(fetchLogs, query.delay) : null;
    if (!query.delay) fetchLogs();
    return () => {
      // Invalidate immediately, including while the replacement request is debouncing.
      active = false;
      clearTimeout(timer);
    };
  }, [generatorId, query]);

  const changeLines = (value) => {
    const lines = Math.max(1, Number.parseInt(value, 10) || 1000);
    setQuery((previous) => (previous.lines === lines ? previous : { lines, delay: 500 }));
  };

  return { lines: query.lines, changeLines, content, loading };
}
