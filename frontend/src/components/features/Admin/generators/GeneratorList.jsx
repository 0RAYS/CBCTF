import { motion } from 'motion/react';
import { IconFileText } from '@tabler/icons-react';
import { Button, Card, EmptyState } from '../../../common';
import { isGeneratorStoppable } from './generatorUtils.js';

const STATUS_STYLES = {
  waiting: 'bg-yellow-400/10 text-yellow-400 border-yellow-400/30',
  pending: 'bg-geek-400/10 text-geek-400 border-geek-400/30',
  terminating: 'bg-orange-400/10 text-orange-300 border-orange-400/30',
  running: 'bg-green-400/10 text-green-400 border-green-400/30',
  stopped: 'bg-neutral-500/10 text-neutral-400 border-neutral-500/30',
};

const COLUMNS = [
  'id',
  'name',
  'contestId',
  'challengeId',
  'success',
  'successLast',
  'failure',
  'failureLast',
  'status',
  'actions',
];

const formatTime = (timestamp) => (timestamp ? new Date(timestamp).toLocaleString() : '\u2014');

export default function GeneratorList({ session, onViewLogs, text, t }) {
  const { generators, selectedIds, loading, toggleSelect, toggleSelectAll } = session;
  const stoppable = generators.filter(isGeneratorStoppable);

  return (
    <Card>
      {loading ? (
        <div className="flex justify-center py-12 text-neutral-400 text-sm">{t('common.loading')}</div>
      ) : generators.length === 0 ? (
        <EmptyState title={text('noGenerators')} description={text('noGeneratorsDesc')} />
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-sm text-neutral-300">
            <thead>
              <tr className="border-b border-neutral-700 text-neutral-400 text-xs uppercase tracking-wider">
                <th className="py-3 px-4 text-left w-10" scope="col">
                  <input
                    type="checkbox"
                    className="accent-geek-400"
                    aria-label={text('stopButton')}
                    checked={stoppable.length > 0 && stoppable.every((generator) => selectedIds.includes(generator.id))}
                    onChange={toggleSelectAll}
                  />
                </th>
                {COLUMNS.map((column) => (
                  <th key={column} className="py-3 px-4 text-left" scope="col">
                    {text(`columns.${column}`)}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {generators.map((generator) => (
                <motion.tr
                  key={generator.id}
                  className="border-b border-neutral-800 hover:bg-neutral-800/40 transition-colors"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                >
                  <td className="py-3 px-4">
                    <input
                      type="checkbox"
                      className="accent-geek-400"
                      aria-label={`${text('stopButton')} ${generator.name}`}
                      checked={selectedIds.includes(generator.id)}
                      disabled={!isGeneratorStoppable(generator)}
                      onChange={() => toggleSelect(generator.id)}
                    />
                  </td>
                  <td className="py-3 px-4 font-mono text-xs text-neutral-500">{generator.id}</td>
                  <td className="py-3 px-4 font-mono text-xs text-neutral-200">{generator.name}</td>
                  <td className="py-3 px-4 text-neutral-400">{generator.contest_id}</td>
                  <td className="py-3 px-4 text-neutral-400">{generator.challenge_name || generator.challenge_id}</td>
                  <td className="py-3 px-4 text-green-400">{generator.success ?? 0}</td>
                  <td className="py-3 px-4 text-neutral-400 text-xs">{formatTime(generator.success_last)}</td>
                  <td className="py-3 px-4 text-red-400">{generator.failure ?? 0}</td>
                  <td className="py-3 px-4 text-neutral-400 text-xs">{formatTime(generator.failure_last)}</td>
                  <td className="py-3 px-4">
                    <span
                      className={`inline-block px-2 py-0.5 rounded border text-xs font-mono ${STATUS_STYLES[generator.status] ?? STATUS_STYLES.stopped}`}
                    >
                      {t(`admin.contests.generators.status.${generator.status}`, generator.status)}
                    </span>
                  </td>
                  <td className="py-3 px-4">
                    {['pending', 'running', 'terminating'].includes(generator.status) && (
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => onViewLogs(generator)}
                        title={text('logs.viewLogs')}
                        aria-label={text('logs.viewLogs')}
                      >
                        <IconFileText size={16} />
                      </Button>
                    )}
                  </td>
                </motion.tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Card>
  );
}
