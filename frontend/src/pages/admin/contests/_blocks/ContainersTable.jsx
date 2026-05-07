import {
  IconBan,
  IconClockPlay,
  IconDownload,
  IconGraph,
  IconRefresh,
  IconTable,
  IconTrash,
} from '@tabler/icons-react';
import { motion } from 'motion/react';
import { Button, Card, EmptyState, Pagination } from '../../../../components/common';

export function ContainersTable({
  t,
  containers,
  runningCount,
  selectedContainers,
  refreshInterval,
  showDeleted,
  currentPage,
  pageSize,
  onRefreshIntervalChange,
  onRefresh,
  onToggleShowDeleted,
  onOpenStopModal,
  onSelectAll,
  onContainerSelect,
  onPageChange,
  onViewTrafficGraph,
  onDownloadTraffic,
  isVictimStoppable,
  formatTime,
  formatRemaining,
  getContainerStatusStyle,
  VictimStatusBadge,
}) {
  const stoppableCount = containers.filter(isVictimStoppable).length;

  return (
    <Card variant="default" padding="none" className="overflow-hidden">
      <div className="p-4 bg-black/20 border-b border-neutral-300/30 space-y-4">
        <div className="flex flex-wrap items-center gap-2">
          <IconTable size={20} className="text-neutral-400" />
          <h3 className="text-lg font-mono text-neutral-50">{t('admin.contests.containers.table.title')}</h3>
          <span className="text-sm font-mono text-neutral-400">
            {t('admin.contests.containers.table.total', { count: runningCount })}
          </span>
        </div>

        <div className="flex flex-wrap gap-2 items-center">
          <div className="flex items-center gap-1 px-2 h-8 rounded-md border border-neutral-700 bg-neutral-900">
            <IconClockPlay size={13} className="text-neutral-400 shrink-0" />
            <span className="text-xs text-neutral-400 shrink-0">{t('common.autoRefresh')}</span>
            <select
              value={refreshInterval}
              onChange={(e) => onRefreshIntervalChange(Number(e.target.value))}
              className="bg-transparent text-xs text-neutral-300 outline-none cursor-pointer"
            >
              {[5, 10, 30, 60].map((seconds) => (
                <option key={seconds} value={seconds} className="bg-neutral-900">
                  {seconds}s
                </option>
              ))}
              <option value={0} className="bg-neutral-900">
                {t('common.autoRefreshOff')}
              </option>
            </select>
          </div>
          <Button variant="ghost" size="sm" icon={<IconRefresh size={14} />} onClick={onRefresh}>
            {t('common.refresh')}
          </Button>
          <Button
            variant={showDeleted ? 'danger' : 'ghost'}
            size="sm"
            icon={<IconTrash size={14} />}
            onClick={onToggleShowDeleted}
          >
            {t('admin.contests.containers.showDeleted')}
          </Button>
          {selectedContainers.length > 0 && (
            <Button variant="danger" size="sm" icon={<IconBan size={14} />} onClick={onOpenStopModal}>
              {t('admin.contests.containers.table.stopButton')} ({selectedContainers.length})
            </Button>
          )}
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="bg-black/40">
              <th className="p-4 text-left text-neutral-400 font-mono">
                <input
                  type="checkbox"
                  checked={stoppableCount > 0 && selectedContainers.length === stoppableCount}
                  onChange={onSelectAll}
                  className="w-4 h-4 rounded border-neutral-300/30 text-geek-400 focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
                  aria-label={t('common.selectAll')}
                />
              </th>
              <TableHeader>{t('admin.contests.containers.table.columns.id')}</TableHeader>
              <TableHeader>{t('admin.contests.containers.table.columns.challenge')}</TableHeader>
              <TableHeader>{t('admin.contests.containers.table.columns.team')}</TableHeader>
              <TableHeader>{t('admin.contests.containers.table.columns.user')}</TableHeader>
              <TableHeader>{t('admin.contests.containers.table.columns.remote')}</TableHeader>
              <TableHeader>{t('admin.contests.containers.table.columns.startTime')}</TableHeader>
              <TableHeader>{t('admin.contests.containers.table.columns.status')}</TableHeader>
              <TableHeader>{t('admin.contests.containers.table.columns.remaining')}</TableHeader>
              <TableHeader>{t('admin.contests.teamDetail.traffic.columns.actions')}</TableHeader>
            </tr>
          </thead>
          <tbody>
            {containers.length === 0 ? (
              <tr>
                <td colSpan="11">
                  <EmptyState title={t('admin.contests.containers.table.empty')} />
                </td>
              </tr>
            ) : (
              containers.map((container, index) => (
                <motion.tr
                  key={container.id}
                  className="border-t border-neutral-300/10 hover:bg-black/40 transition-colors"
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.02 }}
                  whileHover={{ backgroundColor: 'rgba(0, 0, 0, 0.5)' }}
                >
                  <td className="p-4 text-neutral-300 font-mono">
                    <input
                      type="checkbox"
                      checked={selectedContainers.includes(container.id)}
                      disabled={!isVictimStoppable(container)}
                      onChange={() => onContainerSelect(container.id)}
                      className="w-4 h-4 rounded border-neutral-300/30 text-geek-400 focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
                      aria-label={t('admin.contests.containers.table.selectContainer', { id: container.id })}
                    />
                  </td>
                  <TableCell>{container.id}</TableCell>
                  <TableCell>{container.challenge}</TableCell>
                  <TableCell>{container.team}</TableCell>
                  <TableCell>{container.user}</TableCell>
                  <td className="p-4 text-neutral-300 font-mono">
                    {container.remote && container.remote.length > 0 ? (
                      <div className="space-y-1">
                        {container.remote.map((addr, remoteIndex) => (
                          <div
                            key={remoteIndex}
                            className="text-xs bg-black/30 text-geek-400 px-2 py-1 rounded border border-geek-400/30"
                          >
                            {addr}
                          </div>
                        ))}
                      </div>
                    ) : (
                      <span className="text-neutral-400">-</span>
                    )}
                  </td>
                  <TableCell>{formatTime(container.start)}</TableCell>
                  <td className="p-4 whitespace-nowrap">
                    <VictimStatusBadge status={container.status} t={t} />
                  </td>
                  <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">
                    <span
                      className={`px-2 py-1 rounded-md text-xs font-mono border ${getContainerStatusStyle(container.remaining)}`}
                    >
                      {formatRemaining(container.remaining)}
                    </span>
                  </td>
                  <td className="p-4 whitespace-nowrap">
                    <div className="flex items-center gap-2">
                      <Button
                        variant="ghost"
                        size="icon"
                        className="!text-geek-400 hover:!text-geek-300"
                        onClick={() => onViewTrafficGraph(container)}
                        aria-label={t('admin.contests.teamDetail.traffic.actions.viewTraffic')}
                        title={t('admin.contests.teamDetail.traffic.actions.viewTraffic')}
                      >
                        <IconGraph size={18} />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="!text-geek-400 hover:!text-geek-300"
                        onClick={() => onDownloadTraffic(container)}
                        aria-label={t('admin.contests.teamDetail.traffic.actions.downloadTraffic')}
                        title={t('admin.contests.teamDetail.traffic.actions.downloadTraffic')}
                      >
                        <IconDownload size={18} />
                      </Button>
                    </div>
                  </td>
                </motion.tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {runningCount > 0 && (
        <div className="p-4 border-t border-neutral-300/30 bg-black/20 flex justify-center">
          <Pagination
            total={Math.ceil(runningCount / pageSize)}
            current={currentPage}
            pageSize={pageSize}
            onChange={onPageChange}
            showTotal
            totalItems={runningCount}
          />
        </div>
      )}
    </Card>
  );
}

function TableHeader({ children }) {
  return <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap">{children}</th>;
}

function TableCell({ children }) {
  return <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">{children}</td>;
}
