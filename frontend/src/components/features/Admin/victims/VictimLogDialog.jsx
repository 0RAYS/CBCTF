import { useTranslation } from 'react-i18next';
import { AnsiLog, Modal } from '../../../common';
import useVictimLogSession from './useVictimLogSession';

export default function VictimLogDialog({ victim, onClose, loadPods, loadLogs, translationKey }) {
  const { t } = useTranslation();
  // Keep the session outside Modal's retained exit subtree so closing invalidates requests immediately.
  const log = useVictimLogSession({ victim, loadPods, loadLogs, translationKey });
  return (
    <Modal
      isOpen={!!victim}
      onClose={onClose}
      title={t(`${translationKey}.logs.title`, { id: victim?.id ?? '' })}
      size="2xl"
      className="!max-w-[95vw]"
      bodyClassName="p-4 flex flex-col gap-3 max-h-[90vh] overflow-y-auto"
    >
      {log.podsLoading ? (
        <div className="flex justify-center py-12 text-neutral-400 text-sm">{t('common.loading')}</div>
      ) : log.pods.length === 0 ? (
        <p className="text-neutral-500 text-sm py-8 text-center">{t(`${translationKey}.logs.noPods`)}</p>
      ) : (
        <>
          <div className="flex flex-wrap items-end gap-3">
            <div className="flex flex-col gap-1 min-w-[180px]">
              <label className="text-xs text-neutral-400 font-mono">{t(`${translationKey}.logs.podName`)}</label>
              <select
                value={log.podName}
                onChange={(event) => log.selectPod(event.target.value)}
                className="bg-neutral-800 border border-neutral-700 rounded px-2 py-1.5 text-sm text-neutral-200 focus:outline-none focus:border-geek-400"
              >
                {log.pods.map((pod) => (
                  <option key={pod.name} value={pod.name} className="bg-neutral-900">
                    {pod.name}
                  </option>
                ))}
              </select>
            </div>
            <div className="flex flex-col gap-1 min-w-[140px]">
              <label className="text-xs text-neutral-400 font-mono">{t(`${translationKey}.logs.containerName`)}</label>
              <select
                value={log.containerName}
                onChange={(event) => log.selectContainer(event.target.value)}
                className="bg-neutral-800 border border-neutral-700 rounded px-2 py-1.5 text-sm text-neutral-200 focus:outline-none focus:border-geek-400"
              >
                {(log.pods.find((pod) => pod.name === log.podName)?.containers ?? []).map((name) => (
                  <option key={name} value={name} className="bg-neutral-900">
                    {name}
                  </option>
                ))}
              </select>
            </div>
            <div className="flex flex-col gap-1 w-24">
              <label className="text-xs text-neutral-400 font-mono">{t(`${translationKey}.logs.lines`)}</label>
              <input
                type="number"
                min={1}
                step={100}
                value={log.lines}
                onChange={(event) => log.setLines(Math.max(1, Number.parseInt(event.target.value, 10) || 1000))}
                className="bg-neutral-800 border border-neutral-700 rounded px-2 py-1.5 text-sm text-neutral-200 focus:outline-none focus:border-geek-400"
              />
            </div>
            {log.loading && <span className="text-xs text-neutral-500 font-mono pb-1">{t('common.loading')}</span>}
          </div>
          <AnsiLog
            content={log.content}
            loading={log.loading}
            empty={t(`${translationKey}.logs.empty`)}
            className="flex-1"
            scrollToBottom
          />
        </>
      )}
    </Modal>
  );
}
