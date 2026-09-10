import { useTranslation } from 'react-i18next';
import { Card, EmptyState, Pagination } from '../../../common';
import { formatDateTime, formatPayload, PAGE_SIZE, taskRowKey, taskTimestamp } from './taskModel';

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

  return <span className={`inline-block px-2 py-1 rounded border text-xs font-mono ${style}`}>{value || '-'}</span>;
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
          <table className="w-full text-sm text-neutral-300">
            <thead>
              <tr className="bg-black/30 border-b border-neutral-300/10">
                {columns.map((column) => (
                  <th key={column} className="p-4 text-left font-mono text-neutral-400">
                    {t(`admin.tasks.columns.${column}`)}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((item) => (
                <tr key={taskRowKey(item, mode)} className="border-b border-neutral-300/10 align-top">
                  <td className="p-4 font-mono text-xs text-neutral-200 whitespace-nowrap">{item.task_id || '-'}</td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">{item.type}</td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">{item.queue}</td>
                  <td className="p-4 whitespace-nowrap">
                    <StatusBadge value={item.status} />
                  </td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">
                    {item.retry_count}/{item.max_retry}
                  </td>
                  <td className="p-4 font-mono text-xs whitespace-nowrap">
                    {formatDateTime(taskTimestamp(item, mode), i18n.language)}
                  </td>
                  <td className="p-4 max-w-80">
                    <div className="line-clamp-3 break-words text-xs text-red-300">{item.error || '-'}</div>
                  </td>
                  <td className="p-4 max-w-96">
                    <div className="line-clamp-3 break-all font-mono text-xs text-neutral-400">
                      {formatPayload(item.payload)}
                    </div>
                  </td>
                </tr>
              ))}
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
