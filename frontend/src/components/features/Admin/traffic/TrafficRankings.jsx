import { useTranslation } from 'react-i18next';
import { Card } from '../../../common';
import { formatBytes, formatDurationMs, formatProcessLabel } from './trafficPresentation.js';

export default function TrafficRankings({ topology, playback }) {
  const { t, i18n } = useTranslation();
  const topTalkers = topology?.top_talkers || [];
  const topEdges = topology?.top_edges || [];
  const summary = topology?.summary || {};
  const { windowInfo, slice, isPlaying, playbackIndex, playbackFrames } = playback;
  return (
    <Card padding="sm" className="rounded-2xl border-neutral-600 bg-neutral-900">
      <div className="pr-1">
        <div className="grid gap-3">
          <div>
            <div className="text-xs font-mono text-neutral-300">
              {t('admin.contests.trafficGraph.rankings.topTalkers')}
            </div>
            <div className="mt-2 grid max-h-[126px] gap-2 overflow-y-auto pr-1">
              {topTalkers.length > 0 &&
                topTalkers.map((item, index) => (
                  <div
                    key={`${item.ip}-${index}`}
                    className="flex items-center justify-between rounded-lg border border-neutral-600 bg-black/20 px-2.5 py-2"
                  >
                    <div className="min-w-0">
                      <div className="truncate font-mono text-[11px] text-white">{item.label || item.ip}</div>
                      <div className="mt-1 truncate text-[11px] text-neutral-500">
                        {item.dominant_process || formatProcessLabel((item.processes || [])[0]) || item.ip}
                      </div>
                    </div>
                    <div className="ml-3 text-right font-mono">
                      <div className="text-[11px] text-geek-400">{formatBytes(item.bytes)}</div>
                      <div className="mt-1 text-[11px] text-neutral-400">{item.connections || 0} conn</div>
                    </div>
                  </div>
                ))}
            </div>
          </div>
          <div className="border-t border-neutral-700 pt-3">
            <div className="text-xs font-mono text-neutral-300">
              {t('admin.contests.trafficGraph.rankings.topEdges')}
            </div>
            <div className="mt-2 grid max-h-[126px] gap-2 overflow-y-auto pr-1">
              {topEdges.length > 0 &&
                topEdges.map((item, index) => (
                  <div
                    key={`${item.id || item.label}-${index}`}
                    className="flex items-center justify-between rounded-lg border border-neutral-600 bg-black/20 px-2.5 py-2"
                  >
                    <div className="min-w-0">
                      <div className="truncate font-mono text-[11px] text-white">{item.label}</div>
                      <div className="mt-1 truncate text-[11px] text-neutral-500">
                        {item.dominant_process ||
                          formatProcessLabel((item.processes || [])[0]) ||
                          t(`admin.contests.trafficGraph.direction.${item.direction || 'internal'}`)}
                      </div>
                    </div>
                    <div className="ml-3 text-right font-mono">
                      <div className="text-[11px] text-geek-400">{formatBytes(item.bytes)}</div>
                      <div className="mt-1 text-[11px] text-neutral-400">{item.connections || 0} conn</div>
                    </div>
                  </div>
                ))}
            </div>
          </div>
        </div>
      </div>
      <div className="mt-3 border-t border-neutral-700 pt-3">
        <div className="flex flex-wrap gap-2 text-[11px] font-mono text-neutral-500">
          <span>
            {t('admin.contests.trafficGraph.footer.window', {
              start: formatDurationMs(windowInfo.start),
              end: formatDurationMs(windowInfo.end),
            })}
          </span>
          <span>{t('admin.contests.trafficGraph.footer.timeSlice', { count: formatDurationMs(slice) })}</span>
          <span>{t('admin.contests.trafficGraph.footer.connectionCount', { count: summary.visible_edges || 0 })}</span>
          <span>{t('admin.contests.trafficGraph.footer.ipCount', { count: summary.visible_nodes || 0 })}</span>
          <span>{t('admin.contests.trafficGraph.footer.maxDuration', { count: topology?.total_duration || 0 })}</span>
          {isPlaying ? (
            <span>
              {t('admin.contests.trafficGraph.footer.playing', { current: playbackIndex, total: playbackFrames })}
            </span>
          ) : null}
        </div>
        <div className="mt-2 text-[11px] font-mono text-neutral-500">
          {t('admin.contests.trafficGraph.footer.updatedAt', {
            time: new Date().toLocaleString(i18n.language || 'en-US'),
          })}
        </div>
      </div>
    </Card>
  );
}
