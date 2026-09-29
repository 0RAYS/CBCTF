import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../common';
import { batchRows, batchTone } from './batchModel.js';

const PAGE_SIZE = 20;

export default function BatchResultPanel({ result, error, queued = false }) {
  const { t } = useTranslation();
  const [page, setPage] = useState(0);
  useEffect(() => setPage(0), [result]);
  if (!result)
    return error ? (
      <p role="alert" className="rounded border border-red-500/40 p-3 text-sm text-red-300">
        {error}
      </p>
    ) : null;
  const rows = batchRows(result);
  const tone = batchTone(result);
  const colors =
    tone === 'success'
      ? 'border-green-500/40 text-green-300'
      : tone === 'warning'
        ? 'border-amber-500/40 text-amber-200'
        : 'border-red-500/40 text-red-300';
  const lastPage = Math.max(0, Math.ceil(rows.length / PAGE_SIZE) - 1);
  const current = Math.min(page, lastPage);
  return (
    <section
      role="status"
      aria-live="polite"
      className={`min-w-0 rounded-lg border bg-neutral-900 p-3 text-sm ${colors}`}
    >
      <h3 className="font-semibold">{t(`admin.batch.states.${result.status}`)}</h3>
      <p className="mt-1 break-words">{t('admin.batch.summary', result)}</p>
      {queued ? <p className="mt-1 text-xs text-neutral-400">{t('admin.batch.queueHint')}</p> : null}
      {result.not_attempted > 0 ? <p className="mt-1 text-xs">{t('admin.batch.unattemptedHint')}</p> : null}
      <details className="mt-2">
        <summary className="cursor-pointer">{t('admin.batch.details', { count: rows.length })}</summary>
        <ul className="mt-2 grid max-h-72 gap-2 overflow-y-auto text-xs">
          {rows.slice(current * PAGE_SIZE, (current + 1) * PAGE_SIZE).map((item, index) => (
            <li key={`${current}-${index}`} className="rounded border border-neutral-700 p-2 text-neutral-200">
              <code className="block break-all">{item.path}</code>
              <span>
                {t(`admin.batch.itemStates.${item.status}`)} ·{' '}
                {t(`admin.batch.phases.${item.phase}`, { defaultValue: item.phase })}
              </span>
              {item.code ? (
                <div className="mt-1 break-all text-red-300">
                  {t(`admin.batch.errors.${item.code}`, { defaultValue: item.code })}
                </div>
              ) : null}
              {item.details ? <div>{t('admin.batch.summary', item.details)}</div> : null}
            </li>
          ))}
        </ul>
        {rows.length > PAGE_SIZE ? (
          <div className="mt-2 flex items-center justify-end gap-2">
            <Button size="sm" disabled={current === 0} onClick={() => setPage(current - 1)}>
              {t('common.previous')}
            </Button>
            <span>
              {current + 1} / {lastPage + 1}
            </span>
            <Button size="sm" disabled={current === lastPage} onClick={() => setPage(current + 1)}>
              {t('common.next')}
            </Button>
          </div>
        ) : null}
      </details>
    </section>
  );
}
