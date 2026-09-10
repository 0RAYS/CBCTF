import {
  IconDownload,
  IconPlayerPause,
  IconPlayerPlay,
  IconPlayerTrackNext,
  IconPlayerTrackPrev,
  IconRefresh,
} from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button, Chip, Input } from '../../../common';
import { formatDurationMs, INPUT_STEP_MS, MIN_SLICE_MS } from './trafficPresentation.js';

export default function TrafficControls({ playback, isFetching, onRefresh, onDownload }) {
  const { t } = useTranslation();
  const {
    windowInfo,
    maxShift,
    shift,
    setShift,
    slice,
    sliceInput,
    setIsPlaying,
    handleSliceInputChange,
    commitSliceInput,
    isPlaying,
    playbackIndex,
    playbackFrames,
    stepPlayback,
    canStepBackward,
    canStepForward,
    togglePlayback,
    totalDuration,
  } = playback;
  return (
    <div className="rounded-2xl border border-neutral-600 bg-black/20 px-4 py-3">
      <div className="grid gap-2">
        <div>
          <div className="flex items-center justify-between text-[11px] text-neutral-400">
            <span>{t('admin.contests.trafficGraph.controls.timeShift')}</span>
            <span className="font-mono text-geek-400">
              {t('admin.contests.trafficGraph.hero.windowAt', {
                start: formatDurationMs(windowInfo.start),
                end: formatDurationMs(windowInfo.end),
              })}
            </span>
          </div>
          <input
            type="range"
            min="0"
            max={maxShift}
            value={shift}
            onChange={(event) => {
              setIsPlaying(false);
              setShift(Number(event.target.value));
            }}
            className="mt-1 w-full accent-geek-400"
          />
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <label className="mr-1 text-[11px] text-neutral-400" htmlFor="traffic-graph-slice-ms">
            {t('admin.contests.trafficGraph.controls.timeSlice')}
          </label>
          <div className="w-[132px]">
            <Input
              id="traffic-graph-slice-ms"
              type="number"
              min={MIN_SLICE_MS}
              step={INPUT_STEP_MS}
              value={sliceInput}
              onChange={handleSliceInputChange}
              onBlur={commitSliceInput}
              onKeyDown={(event) => {
                if (event.key === 'Enter') event.currentTarget.blur();
              }}
              className="!h-8 !border-neutral-600 !bg-black/20 !px-3 text-[11px] font-mono"
            />
          </div>
          <Chip
            label={t('admin.contests.trafficGraph.controls.milliseconds')}
            variant="tag"
            size="sm"
            colorClass="border-neutral-500/30 bg-black/20 text-neutral-300"
          />
          <div className="ml-auto flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon"
              className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
              onClick={onRefresh}
              title={t('common.refresh')}
            >
              <IconRefresh size={16} className={isFetching ? 'animate-spin' : ''} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
              onClick={onDownload}
              title={t('admin.contests.teamDetail.traffic.actions.downloadTraffic')}
            >
              <IconDownload size={16} />
            </Button>
          </div>
        </div>
        <div className="flex flex-wrap items-center gap-2 border-t border-neutral-700/80 pt-2">
          <Chip
            label={t('admin.contests.trafficGraph.footer.timeSlice', { count: formatDurationMs(slice) })}
            variant="tag"
            size="sm"
            colorClass="border-neutral-500/30 bg-black/20 text-neutral-300"
          />
          {isPlaying ? (
            <Chip
              label={t('admin.contests.trafficGraph.footer.playing', { current: playbackIndex, total: playbackFrames })}
              variant="tag"
              size="sm"
              colorClass="border-geek-400/30 bg-geek-400/10 text-geek-400"
            />
          ) : null}
          <div className="ml-auto flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon"
              className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
              onClick={() => stepPlayback(-1)}
              disabled={!canStepBackward || isFetching}
              title={t('common.previous')}
            >
              <IconPlayerTrackPrev size={16} />
            </Button>
            <Button
              variant={isPlaying ? 'outline' : 'primary'}
              size="sm"
              className="!h-8 !px-3"
              onClick={togglePlayback}
              disabled={totalDuration <= 0}
              icon={isPlaying ? <IconPlayerPause size={14} /> : <IconPlayerPlay size={14} />}
            >
              {isPlaying
                ? t('admin.contests.trafficGraph.controls.pauseReplay')
                : t('admin.contests.trafficGraph.controls.playReplay')}
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
              onClick={() => stepPlayback(1)}
              disabled={!canStepForward || isFetching}
              title={t('common.next')}
            >
              <IconPlayerTrackNext size={16} />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
