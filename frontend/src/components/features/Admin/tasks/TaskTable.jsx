import { useTranslation } from 'react-i18next';
import { Card, EmptyState, Pagination } from '../../../common';
import { formatDateTime, formatPayload, PAGE_SIZE, taskRowKey, taskTimestamp } from './taskModel';
import TruncatedText from '../../../common/TruncatedText';

const COLUMN_WIDTHS = [220, 160, 140, 128, 112, 200, 240, 280];

function StatusBadge({ value }) {
  const style =
    {
      success: 'bg-green-400/10 text-green-400 border-green-400/30',
      failed: 'bg-red-400/10 text-red-400 border-red-400/30',
      active: 'bg-geek-400/10 text-geek-400 border-geek-400/30',
      pending: 'bg-yellow-400/10 text-yellow-400 border-yellow-400/30',
      scheduled: 'bg-blue-400/10 text-blue-400 border-blue-400/30',
      retry: 'bg-orange-400/10 text-orange-300 border-orange-400/30',
      archived: 'bg-red-400/10 text-red-400 border-red-400/30',
      completed: 'bg-green-400/10 text-green-400 border-green-400/30',
    }[value] || 'bg-neutral-500/10 text-neutral-300 border-neutral-500/30';

  return (
    <span
      className={`inline-block max-w-full truncate px-2 py-1 rounded border text-xs font-mono ${style}`}
      title={value}
    >
      {value || '-'}
    </span>
  );
}

export default function TaskTable({ mode, rows, totalCount, currentPage, onPageChange }) {
  const { t, i18n } = useTranslation();
  const columns = [
    'taskId',
    'type',
    'queue',
    'status',
    'retry',
    mode === 'history' ? 'processedAt' : 'nextProcessAt',
    'error',
    'payload',
  ];

  return (
    <Card padding="none" className="overflow-hidden">
      {rows.length === 0 ? (
        <EmptyState title={t(`admin.tasks.${mode}.empty`)} className="py-20" />
      ) : (
        <div className="overflow-x-auto">
          <table
            className="w-full table-fixed text-sm text-neutral-300"
            style={{ minWidth: COLUMN_WIDTHS.reduce((a, b) => a + b, 0) }}
          >
            <colgroup>
              {COLUMN_WIDTHS.map((width, index) => (
                <col key={index} style={{ width }} />
              ))}
            </colgroup>
            <thead>
              <tr className="bg-black/30 border-b border-neutral-300/10">
                {columns.map((column) => (
                  <th key={column} className="p-4 text-left font-mono text-neutral-400">
                    <TruncatedText>{t(`admin.tasks.columns.${column}`)}</TruncatedText>
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((item) => {
                const payload = formatPayload(item.payload);
                return (
                  <tr key={taskRowKey(item, mode)} className="border-b border-neutral-300/10 align-top">
                    <td className="p-4 font-mono text-xs text-neutral-200">
                      <TruncatedText>{item.task_id || '-'}</TruncatedText>
                    </td>
                    <td className="p-4 font-mono text-xs">
                      <TruncatedText>{item.type}</TruncatedText>
                    </td>
                    <td className="p-4 font-mono text-xs">
                      <TruncatedText>{item.queue}</TruncatedText>
                    </td>
                    <td className="p-4 whitespace-nowrap">
                      <StatusBadge value={item.status} />
                    </td>
                    <td className="p-4 font-mono text-xs whitespace-nowrap">
                      <TruncatedText>{`${item.retry_count}/${item.max_retry}`}</TruncatedText>
                    </td>
                    <td className="p-4 font-mono text-xs whitespace-nowrap">
                      <TruncatedText>{formatDateTime(taskTimestamp(item, mode), i18n.language)}</TruncatedText>
                    </td>
                    <td className="p-4 max-w-80">
                      <div className="line-clamp-3 [overflow-wrap:anywhere] text-xs text-red-300" title={item.error}>
                        {item.error || '-'}
                      </div>
                    </td>
                    <td className="p-4 max-w-96">
                      <div className="line-clamp-3 break-all font-mono text-xs text-neutral-400" title={payload}>
                        {payload}
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
      {totalCount > PAGE_SIZE && (
        <div className="p-4 border-t border-neutral-300/10 bg-black/20 flex justify-center">
          <Pagination
            current={currentPage}
            total={Math.ceil(totalCount / PAGE_SIZE)}
            totalItems={totalCount}
            showTotal
            onChange={onPageChange}
          />
        </div>
      )}
    </Card>
  );
}
