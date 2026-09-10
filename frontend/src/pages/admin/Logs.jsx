import { useCallback, useEffect, useLayoutEffect, useRef, useState } from 'react';
import { toast } from '../../utils/toast';
import { getSystemLogs } from '../../api/admin/system';
import { Button, AnsiLog } from '../../components/common';
import { IconRefresh } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import IpLookupDialog from '../../components/features/Admin/network/IpLookupDialog';
import useIpLookup from '../../components/features/Admin/network/useIpLookup';
import { injectClickableIps, isPublicIp } from '../../components/features/Admin/network/logIps';

const LOG_LEVELS = ['DEBUG', 'INFO', 'WARNING', 'ERROR', 'FATAL', 'PANIC'];

const IP_ALLOWED_ATTR = ['data-ip', 'role', 'tabindex'];

function AdminLogs() {
  const [logs, setLogs] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [level, setLevel] = useState('INFO');
  const pageSize = 100;
  const containerRef = useRef(null);
  const sentinelRef = useRef(null);
  const { t } = useTranslation();
  const pageRef = useRef(1);
  const loadingRef = useRef(false);
  const hasMoreRef = useRef(true);
  const scopeRef = useRef(null);
  const { lookup, close: closeIpLookup, dialogProps } = useIpLookup(level);
  const fetchLogs = useCallback(
    async (nextPage) => {
      const scope = scopeRef.current;
      if (!scope || scope.level !== level || loadingRef.current) return;
      loadingRef.current = true;
      try {
        const res = await getSystemLogs({ limit: pageSize, offset: (nextPage - 1) * pageSize, level });
        if (scopeRef.current !== scope) return;
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
        if (scopeRef.current !== scope) return;
        toast.danger({ description: error.message || t('admin.logs.toast.fetchFailed') });
        hasMoreRef.current = false;
        setHasMore(false);
      } finally {
        if (scopeRef.current === scope) loadingRef.current = false;
      }
    },
    [level, t]
  );

  const handleRefresh = () => {
    scopeRef.current = { level };
    closeIpLookup();
    hasMoreRef.current = true;
    setHasMore(true);
    setLogs([]);
    pageRef.current = 1;
    loadingRef.current = false;
    if (containerRef.current) containerRef.current.scrollTop = 0;
    fetchLogs(1);
  };

  const handleLevelChange = (event) => {
    setLevel(event.target.value);
  };

  useLayoutEffect(() => {
    handleRefresh();
    return () => {
      scopeRef.current = null;
    };
  }, [level]);

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

    // Re-observe after each page so a still-visible sentinel can load the next one.
    io.observe(sentinel);
    return () => io.disconnect();
  }, [fetchLogs, logs]);

  const handleLogActivation = (event) => {
    if (event.type === 'keydown' && (event.repeat || !['Enter', ' '].includes(event.key))) return;
    const target = event.target.closest?.('.ip-lookup-trigger[data-ip]');
    if (!target || !containerRef.current?.contains(target)) return;
    const ip = target.getAttribute('data-ip');
    if (!isPublicIp(ip)) return;
    event.preventDefault();
    event.stopPropagation();
    lookup(ip);
  };

  return (
    <div className="w-full mx-auto">
      <div className="mb-8 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
        <label className="flex items-center gap-2 text-sm text-neutral-300">
          <span>{t('admin.logs.levelFilter')}</span>
          <select
            value={level}
            onChange={handleLevelChange}
            className="rounded-md border border-neutral-300/20 bg-black/30 px-3 py-2 text-sm text-neutral-100 outline-none focus:border-primary-400"
          >
            {LOG_LEVELS.map((item) => (
              <option key={item} value={item} className="bg-neutral-950 text-neutral-100">
                {item}
              </option>
            ))}
          </select>
        </label>
        <Button variant="primary" size="sm" align="icon-left" icon={<IconRefresh size={16} />} onClick={handleRefresh}>
          {t('common.refresh')}
        </Button>
      </div>

      <div onKeyDown={handleLogActivation}>
        <AnsiLog
          ref={containerRef}
          content={logs}
          postProcess={injectClickableIps}
          allowedAttr={IP_ALLOWED_ATTR}
          onClick={handleLogActivation}
          className="max-h-[70vh]"
          sentinel={
            <div ref={sentinelRef} className="h-8 flex items-center justify-center text-neutral-500 text-xs">
              {hasMore ? t('admin.logs.loadMore') : t('admin.logs.noMore')}
            </div>
          }
        />
      </div>

      <IpLookupDialog {...dialogProps} />
    </div>
  );
}

export default AdminLogs;
