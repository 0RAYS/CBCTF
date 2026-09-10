import { useLayoutEffect, useRef, useState } from 'react';
import { motion } from 'motion/react';
import { IconCheck, IconClipboard, IconLoader2 } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import Button from '../../../../common/Button';
import { isInstanceTransitioning, normalizeInstanceStatus } from './testSession';

const formatTimeLeft = (seconds) => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);
  return [hours, minutes, secs].map((value) => String(value).padStart(2, '0')).join(':');
};

export default function TestInstancePanel({ testStatus, loading, timeLeft, start, stop }) {
  const { t } = useTranslation();
  const [copied, setCopied] = useState({});
  const lifetime = useRef(null);
  useLayoutEffect(() => {
    const session = { active: true, timers: new Map() };
    lifetime.current = session;
    return () => {
      session.active = false;
      session.timers.forEach(clearTimeout);
    };
  }, []);

  async function copy(target) {
    const session = lifetime.current;
    if (!session?.active) return;
    try {
      await navigator.clipboard.writeText(target);
      if (!session.active) return;
      clearTimeout(session.timers.get(target));
      setCopied((prev) => ({ ...prev, [target]: true }));
      session.timers.set(
        target,
        setTimeout(() => {
          session.timers.delete(target);
          setCopied((prev) => ({ ...prev, [target]: false }));
        }, 2000)
      );
    } catch {
      // Clipboard access can be denied; do not show a false success indicator.
    }
  }

  const status = normalizeInstanceStatus(testStatus?.remote?.status);
  const running = status === 'running';
  const transitioning = isInstanceTransitioning(status);
  const duration = Number(testStatus?.remote?.duration) || 0;
  const progress = duration > 0 ? Math.max(0, Math.min(100, (timeLeft / duration) * 100)) : 0;
  const launchLabel =
    status === 'waiting' || status === 'terminating'
      ? `instance.${status}`
      : status === 'pending'
        ? 'actions.launching'
        : 'actions.launch';
  const indicatorClass = running
    ? 'bg-green-400'
    : status === 'terminating'
      ? 'bg-orange-400 animate-pulse'
      : status === 'pending'
        ? 'bg-yellow-400 animate-pulse'
        : status === 'waiting'
          ? 'bg-yellow-400'
          : 'bg-neutral-400';

  return (
    <section className="space-y-1.5">
      <h3 className="text-neutral-400 font-mono text-sm">{t('admin.challenge.testModal.sections.instance')}</h3>
      {(testStatus || loading.starting) && (
        <div className="space-y-3">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="flex flex-wrap items-center gap-4">
              <div className="flex items-center gap-2" role="status">
                <span className={`w-2 h-2 rounded-full transition-colors duration-300 ${indicatorClass}`} />
                <span className="text-neutral-50 font-mono text-sm">
                  {t(`admin.challenge.testModal.instance.${status || 'notRunning'}`)}
                </span>
              </div>
              {running && (
                <div className="flex items-center gap-2">
                  <span className="text-neutral-400 text-sm">{t('admin.challenge.testModal.instance.time')}</span>
                  <span className="text-yellow-400 font-mono text-sm">{formatTimeLeft(timeLeft)}</span>
                </div>
              )}
            </div>
            {running ? (
              <Button variant="danger" size="sm" onClick={stop} disabled={loading.stopping || loading.status}>
                <span className="flex items-center gap-2">
                  {loading.stopping && <IconLoader2 size={16} className="animate-spin" aria-hidden="true" />}
                  {t(`admin.challenge.testModal.actions.${loading.stopping ? 'stopping' : 'stop'}`)}
                </span>
              </Button>
            ) : (
              <Button
                variant="primary"
                size="sm"
                onClick={start}
                disabled={loading.starting || loading.stopping || loading.status || transitioning}
                className={transitioning ? 'border-yellow-400 text-yellow-400' : ''}
              >
                {t(`admin.challenge.testModal.${launchLabel}`)}
              </Button>
            )}
          </div>
          {(running || transitioning) && (
            <div className="h-1.5 bg-neutral-700 rounded-full overflow-hidden">
              {running ? (
                <motion.div
                  className="h-full bg-yellow-400"
                  initial={{ width: 0 }}
                  animate={{ width: `${progress}%` }}
                  transition={{ duration: 0.5 }}
                />
              ) : status === 'waiting' ? (
                <div className="h-full w-full bg-yellow-400/35" />
              ) : (
                <motion.div
                  className={`h-full w-2/5 rounded-full ${status === 'terminating' ? 'bg-orange-400/70' : 'bg-yellow-400/60'}`}
                  animate={{ x: ['-100%', '350%'] }}
                  transition={{ duration: status === 'terminating' ? 1.1 : 1.4, repeat: Infinity, ease: 'easeInOut' }}
                />
              )}
            </div>
          )}
          {running && testStatus.remote?.target?.length > 0 && (
            <div>
              <div className="p-1 text-neutral-400 text-xs">{t('admin.challenge.testModal.instance.address')}</div>
              {testStatus.remote.target.map((target, index) => (
                <div key={`${target}-${index}`} className="flex items-center justify-between gap-2 bg-neutral-900">
                  <button
                    type="button"
                    className="min-w-0 text-left font-mono text-neutral-50 text-sm cursor-pointer [overflow-wrap:anywhere]"
                    onClick={() => copy(target)}
                    aria-label={`${t('common.copy')} ${target}`}
                  >
                    {target}
                  </button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="shrink-0 !text-neutral-400 hover:!text-geek-400"
                    onClick={() => copy(target)}
                    aria-label={`${t('common.copy')} ${target}`}
                  >
                    {copied[target] ? <IconCheck size={16} /> : <IconClipboard size={16} />}
                  </Button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </section>
  );
}
