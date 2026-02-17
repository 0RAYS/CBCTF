import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { toast } from '../../utils/toast';
import { getSystemLogs } from '../../api/admin/system';
import { ansiToHtml } from '../../utils/ansi';
import { Button } from '../../components/common';
import { IconRefresh } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';

function AdminLogs() {
  const [logs, setLogs] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const pageSize = 100;
  const containerRef = useRef(null);
  const sentinelRef = useRef(null);
  const { t } = useTranslation();
  const pageRef = useRef(1);
  const loadingRef = useRef(false);
  const hasMoreRef = useRef(true);

  const fetchLogs = useCallback(
    async (nextPage) => {
      if (loadingRef.current) return;
      loadingRef.current = true;
      try {
        const res = await getSystemLogs({ limit: pageSize, offset: (nextPage - 1) * pageSize });
        if (res.code === 200) {
          const list = Array.isArray(res.data) ? res.data : [];
          setLogs((prev) => (nextPage === 1 ? list : [...prev, ...list]));
          pageRef.current = nextPage;
          if (list.length < pageSize) {
            hasMoreRef.current = false;
            setHasMore(false);
          }
        } else {
          hasMoreRef.current = false;
          setHasMore(false);
        }
      } catch (error) {
        toast.danger({ description: error.message || t('admin.logs.toast.fetchFailed') });
        hasMoreRef.current = false;
        setHasMore(false);
      } finally {
        loadingRef.current = false;
      }
    },
    [t]
  );

  useEffect(() => {
    fetchLogs(1);
  }, [fetchLogs]);

  const handleRefresh = () => {
    hasMoreRef.current = true;
    setHasMore(true);
    setLogs([]);
    pageRef.current = 1;
    loadingRef.current = false;
    if (containerRef.current) containerRef.current.scrollTop = 0;
    fetchLogs(1);
  };

  useEffect(() => {
    const container = containerRef.current;
    const sentinel = sentinelRef.current;
    if (!container || !sentinel) return;

    const io = new IntersectionObserver(
      (entries) => {
        const entry = entries[0];
        if (entry.isIntersecting && hasMoreRef.current && !loadingRef.current) {
          const next = pageRef.current + 1;
          fetchLogs(next);
        }
      },
      { root: container, threshold: 1.0 }
    );

    io.observe(sentinel);
    return () => io.disconnect();
  }, [fetchLogs]);

  const rendered = useMemo(() => {
    return logs
      .map((line) => {
        const html = ansiToHtml(line).replace(/\n$/, '');
        return `<div class="whitespace-pre-wrap break-words leading-6 font-mono text-sm">${html}</div>`;
      })
      .join('');
  }, [logs]);

  return (
    <div className="w-full mx-auto">
      <div className="mb-8 flex items-center justify-end">
        <Button variant="primary" size="sm" align="icon-left" icon={<IconRefresh size={16} />} onClick={handleRefresh}>
          {t('common.refresh')}
        </Button>
      </div>

      <div
        ref={containerRef}
        className="border border-neutral-300/20 rounded-md bg-black/30 p-3 overflow-auto max-h-[70vh]"
      >
        <div className="max-w-full">
          <div dangerouslySetInnerHTML={{ __html: rendered }} />
        </div>
        <div ref={sentinelRef} className="h-8 flex items-center justify-center text-neutral-500 text-xs">
          {hasMore ? t('admin.logs.loadMore') : t('admin.logs.noMore')}
        </div>
      </div>
    </div>
  );
}

export default AdminLogs;
