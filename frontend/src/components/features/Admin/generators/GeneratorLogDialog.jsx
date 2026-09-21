import { AnsiLog, Modal } from '../../../common';
import PodDiagnostics from '../workloads/PodDiagnostics';
import useWorkloadStatus from '../workloads/useWorkloadStatus';
import useGeneratorLogs from './useGeneratorLogs.js';

export default function GeneratorLogDialog({ api, generator, onClose, text, t }) {
  const status = useWorkloadStatus(generator.id, api.status);
  const logs = useGeneratorLogs(api, generator.id, text);

  return (
    <Modal
      isOpen
      onClose={onClose}
      title={text('logs.title', { name: generator.name ?? '' })}
      size="2xl"
      className="!max-w-[95vw]"
      bodyClassName="p-4 flex flex-col gap-3 max-h-[90vh] overflow-y-auto"
    >
      <PodDiagnostics {...status} />
      <div className="flex items-end gap-3">
        <label className="flex flex-col gap-1 w-24">
          <span className="text-xs text-neutral-400 font-mono shrink-0">{text('logs.lines')}</span>
          <input
            type="number"
            min={1}
            max={10000}
            step={100}
            value={logs.lines}
            onChange={(event) => logs.changeLines(event.target.value)}
            className="bg-neutral-800 border border-neutral-700 rounded px-2 py-1.5 text-sm text-neutral-200 text-center focus:outline-none focus:border-geek-400"
          />
        </label>
        <button
          type="button"
          onClick={logs.refreshLogs}
          disabled={logs.loading}
          className="text-xs text-geek-400 disabled:opacity-50 pb-1"
        >
          {t('admin.workloads.refreshLogs')}
        </button>
        {logs.loading && <span className="text-xs text-neutral-500 font-mono pb-1">{t('common.loading')}</span>}
      </div>
      <AnsiLog
        content={logs.content}
        loading={logs.loading}
        empty={text('logs.empty')}
        className="flex-1"
        scrollToBottom
      />
    </Modal>
  );
}
