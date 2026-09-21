import { useTranslation } from 'react-i18next';

export default function PodDiagnostics({ pods = [], refreshing, statusError, updatedAt, refreshStatus }) {
  const { t } = useTranslation();
  const text = (key, options) => t(`admin.workloads.${key}`, options);
  return (
    <section className="rounded border border-neutral-700 bg-neutral-900 p-3 space-y-3" aria-label={text('title')}>
      <div className="flex flex-wrap items-center gap-3 text-xs">
        <h3 className="text-neutral-200 font-medium">{text('title')}</h3>
        <span className="text-neutral-500">{text('polling')}</span>
        {updatedAt && <span className="text-neutral-400">{text('updated', { time: updatedAt })}</span>}
        <button
          type="button"
          onClick={refreshStatus}
          disabled={refreshing}
          className="ml-auto text-geek-400 disabled:opacity-50 focus-visible:outline focus-visible:outline-geek-400"
        >
          {refreshing ? t('common.loading') : text('refresh')}
        </button>
      </div>
      {statusError && (
        <p role="status" className="text-xs text-amber-400">
          {text('failed')}
        </p>
      )}
      {!pods.length && !refreshing && <p className="text-xs text-neutral-400">{text('noPods')}</p>}
      {pods.map((pod) => (
        <div key={pod.uid} className="text-xs space-y-2 min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <span className="font-mono text-neutral-300 break-all">{pod.name}</span>
            <span className={pod.ready ? 'text-emerald-400' : 'text-amber-400'}>
              {text(pod.terminating ? 'terminating' : pod.ready ? 'ready' : 'notReady')}
            </span>
            <span className="text-neutral-400">
              {pod.status} {pod.reason}
            </span>
          </div>
          <div className="flex flex-wrap gap-x-4 gap-y-1 text-neutral-400">
            {pod.conditions
              .filter((condition) => condition.status !== 'True')
              .map((condition) => (
                <span key={condition.type} className="break-all">
                  {condition.type}: {condition.status} {condition.reason}
                </span>
              ))}
          </div>
          <ul className="space-y-1 text-neutral-400">
            {pod.container_statuses.map((container) => (
              <li key={`${container.init}:${container.name}`} className="flex flex-wrap gap-x-3 gap-y-1">
                <span className="font-mono break-all">
                  {container.name}
                  {container.init ? ' (init)' : ''}
                </span>
                <span>
                  {text(`state.${container.state}`)} {container.reason}
                </span>
                <span>{text('restarts', { count: container.restarts })}</span>
                {container.exit_code != null && <span>{text('exit', { code: container.exit_code })}</span>}
                {container.last_exit_code != null && (
                  <span>
                    {text('lastExit', { code: container.last_exit_code })} {container.last_reason}
                  </span>
                )}
              </li>
            ))}
          </ul>
        </div>
      ))}
      <p className="text-xs text-neutral-500">{text('logLimit')}</p>
    </section>
  );
}
