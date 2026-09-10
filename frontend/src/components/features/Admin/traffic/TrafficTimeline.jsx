import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import { Card } from '../../../common';
import { clamp, formatBytes, formatDurationMs } from './trafficPresentation.js';

export default function TrafficTimeline({ timeline, windowInfo, setShift }) {
  const { t } = useTranslation();
  const peakTimeline = Math.max(1, ...timeline.map((bucket) => bucket.bytes || 0));
  return (
    <Card padding="sm" className="rounded-2xl border-neutral-600 bg-neutral-900">
      <div className="flex items-center justify-between gap-2">
        <div className="text-sm font-mono text-neutral-300">{t('admin.contests.trafficGraph.timeline.title')}</div>
        <div className="text-[11px] font-mono text-neutral-500">
          {t('admin.contests.trafficGraph.hero.windowAt', {
            start: formatDurationMs(windowInfo.start),
            end: formatDurationMs(windowInfo.end),
          })}
        </div>
      </div>
      {timeline.length > 0 && (
        <>
          <div className="mt-3 flex h-16 items-end gap-px">
            {timeline.map((bucket) => {
              const active =
                bucket.timestamp_ms >= windowInfo.start &&
                bucket.timestamp_ms < Math.max(windowInfo.end, windowInfo.start + 1);
              const ratio = clamp((bucket.bytes || 0) / peakTimeline, 0.08, 1);
              return (
                <button
                  key={bucket.timestamp_ms}
                  type="button"
                  title={`${t('admin.contests.trafficGraph.timeline.at', { value: formatDurationMs(bucket.timestamp_ms) })} / ${formatBytes(bucket.bytes)}`}
                  onClick={() =>
                    setShift(Math.min(bucket.timestamp_ms, Math.max(0, windowInfo.total - windowInfo.duration)))
                  }
                  className="group relative flex min-w-0 flex-1 items-end"
                >
                  <motion.span
                    className={`block w-full rounded-t-[3px] ${active ? 'bg-geek-400' : 'bg-neutral-600 group-hover:bg-neutral-500'}`}
                    initial={{ height: '0%' }}
                    animate={{ height: `${ratio * 100}%` }}
                    transition={{ duration: 0.28, ease: 'easeOut' }}
                  />
                </button>
              );
            })}
          </div>
          <div className="mt-2 flex items-center justify-between text-[11px] font-mono text-neutral-500">
            <span>
              {timeline[0]
                ? t('admin.contests.trafficGraph.timeline.at', { value: formatDurationMs(timeline[0].timestamp_ms) })
                : '--'}
            </span>
            <span>{formatBytes(peakTimeline)}</span>
            <span>
              {timeline[timeline.length - 1]
                ? t('admin.contests.trafficGraph.timeline.at', {
                    value: formatDurationMs(timeline[timeline.length - 1].timestamp_ms),
                  })
                : '--'}
            </span>
          </div>
        </>
      )}
    </Card>
  );
}
