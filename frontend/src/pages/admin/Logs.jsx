import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { toast } from '../../utils/toast';
import { getSystemLogs } from '../../api/admin/system';
import { getIpInfo } from '../../api/admin/contest';
import { Button, AnsiLog } from '../../components/common';
import { IconRefresh } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import Modal from '../../components/common/Modal';
import Card from '../../components/common/Card';

const LOG_LEVELS = ['DEBUG', 'INFO', 'WARNING', 'ERROR', 'FATAL', 'PANIC'];

// 匹配 IPv4 地址的正则（在 HTML 文本节点中）
const IPV4_RE = /\b(\d{1,3}\.){3}\d{1,3}\b/g;

// 判断是否为公网 IP（排除私有/保留地址）
function isPublicIp(ip) {
  const parts = ip.split('.').map(Number);
  if (parts.length !== 4 || parts.some((p) => isNaN(p) || p < 0 || p > 255)) return false;
  const [a, b] = parts;
  const isPrivate =
    a === 10 ||
    a === 127 ||
    (a === 172 && b >= 16 && b <= 31) ||
    (a === 192 && b === 168) ||
    (a === 169 && b === 254) ||
    (a === 100 && b >= 64 && b <= 127) ||
    a === 0 ||
    (a === 198 && (b === 18 || b === 19)) ||
    (a === 192 && b === 0 && parts[2] === 0) ||
    (a === 192 && b === 0 && parts[2] === 2) ||
    (a === 192 && b === 88 && parts[2] === 99) ||
    (a === 198 && b === 51 && parts[2] === 100) ||
    (a === 203 && b === 0 && parts[2] === 113) ||
    a >= 240;
  return !isPrivate;
}

// 在已转义的 HTML 字符串中，将公网 IP 替换为可点击的 span（data-ip 属性）
function injectClickableIps(html) {
  return html.replace(/(<[^>]+>)|([^<]+)/g, (match, tag, text) => {
    if (tag) return tag;
    if (!text) return match;
    return text.replace(IPV4_RE, (ip) => {
      if (!isPublicIp(ip)) return ip;
      return `<span data-ip="${ip}" class="ip-lookup-trigger" style="color:#597ef7;cursor:pointer;text-decoration:underline;text-decoration-style:dotted;">${ip}</span>`;
    });
  });
}

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

  // IP 反查弹窗状态
  const [ipModal, setIpModal] = useState({ open: false, ip: null, data: null, loading: false });
  const fetchLogs = useCallback(
    async (nextPage) => {
      if (loadingRef.current) return;
      loadingRef.current = true;
      try {
        const res = await getSystemLogs({ limit: pageSize, offset: (nextPage - 1) * pageSize, level });
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
    [level, t]
  );

  const handleRefresh = () => {
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

  useEffect(() => {
    handleRefresh();
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

    io.observe(sentinel);
    return () => io.disconnect();
  }, [fetchLogs]);

  // postProcess: 在每行 ANSI→HTML 后注入可点击 IP
  const postProcess = useMemo(() => injectClickableIps, []);

  // 事件委托：捕获日志区域内所有 data-ip 点击
  const handleLogClick = useCallback(
    async (e) => {
      const target = e.target.closest('[data-ip]');
      if (!target) return;
      const ip = target.getAttribute('data-ip');
      if (!ip) return;
      setIpModal({ open: true, ip, data: null, loading: true });
      try {
        const res = await getIpInfo(ip);
        if (res.code === 200) {
          setIpModal((prev) => ({ ...prev, data: { ip, ...res.data }, loading: false }));
        } else {
          setIpModal((prev) => ({ ...prev, data: { ip }, loading: false }));
        }
      } catch (error) {
        toast.danger({ description: error.message || t('admin.logs.toast.ipFetchFailed') });
        setIpModal((prev) => ({ ...prev, data: { ip }, loading: false }));
      }
    },
    [t]
  );

  const handleIpModalClose = () => {
    setIpModal({ open: false, ip: null, data: null, loading: false });
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

      <AnsiLog
        ref={containerRef}
        content={logs}
        postProcess={postProcess}
        allowedAttr={['data-ip']}
        onClick={handleLogClick}
        className="max-h-[70vh]"
        sentinel={
          <div ref={sentinelRef} className="h-8 flex items-center justify-center text-neutral-500 text-xs">
            {hasMore ? t('admin.logs.loadMore') : t('admin.logs.noMore')}
          </div>
        }
      />

      {/* IP 反查弹窗 */}
      <Modal isOpen={ipModal.open} onClose={handleIpModalClose} title={t('admin.logs.ipDetail.title')} size="sm">
        {ipModal.loading ? (
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-geek-400" />
          </div>
        ) : ipModal.data ? (
          <Card variant="default" padding="md">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">{t('admin.logs.ipDetail.ip')}</label>
                <p className="text-neutral-300 font-mono">{ipModal.data.ip || '-'}</p>
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">{t('admin.logs.ipDetail.iso')}</label>
                <p className="text-neutral-300">{ipModal.data.iso || '-'}</p>
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.logs.ipDetail.country')}
                </label>
                <p className="text-neutral-300">{ipModal.data.country || '-'}</p>
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.logs.ipDetail.subdivision')}
                </label>
                <p className="text-neutral-300">{ipModal.data.subdivision || '-'}</p>
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">{t('admin.logs.ipDetail.city')}</label>
                <p className="text-neutral-300">{ipModal.data.city || '-'}</p>
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.logs.ipDetail.timezone')}
                </label>
                <p className="text-neutral-300">{ipModal.data.timezone || '-'}</p>
              </div>
              <div className="col-span-2">
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.logs.ipDetail.coordinates')}
                </label>
                <p className="text-neutral-300 font-mono">
                  {ipModal.data.latitude != null && ipModal.data.longitude != null
                    ? `${ipModal.data.latitude}, ${ipModal.data.longitude}`
                    : '-'}
                </p>
              </div>
            </div>
          </Card>
        ) : null}
      </Modal>
    </div>
  );
}

export default AdminLogs;
