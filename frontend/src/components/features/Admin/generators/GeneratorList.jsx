import { motion } from 'motion/react';
import { IconFileText } from '@tabler/icons-react';
import { Card, EmptyState } from '../../../common';
import useContextMenu from '../../../common/useContextMenu';
import RowActions from '../../../common/RowActions';
import { isGeneratorStoppable } from './generatorUtils.js';
import TruncatedText from '../../../common/TruncatedText';

const COLUMN_WIDTHS = [48, 80, 200, 256, 112, 176, 96, 184, 96, 184, 128, 88];

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
  'image',
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
  const getRowActions = (generator) => [
    {
      key: 'logs',
      inline: true,
      label: text('logs.viewLogs'),
      icon: <IconFileText size={18} />,
      hidden: !['pending', 'running', 'terminating'].includes(generator.status),
      onClick: () => onViewLogs(generator),
    },
  ];
  const contextMenu = useContextMenu(generators, getRowActions, loading);

  return (
    <Card>
      {contextMenu.hint}
      {contextMenu.menu}
      {loading ? (
        <div className="flex justify-center py-12 text-neutral-400 text-sm">{t('common.loading')}</div>
      ) : generators.length === 0 ? (
        <EmptyState title={text('noGenerators')} description={text('noGeneratorsDesc')} />
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
                    <TruncatedText>{text(`columns.${column}`)}</TruncatedText>
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {generators.map((generator) => (
                <motion.tr
                  key={generator.id}
                  {...contextMenu.getRowProps(generator)}
                  className="border-b border-neutral-800 hover:bg-neutral-800/40 transition-colors focus-visible:outline focus-visible:outline-geek-400 focus-visible:-outline-offset-2"
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
                  <td className="py-3 px-4 font-mono text-xs text-neutral-500">
                    <TruncatedText>{generator.id}</TruncatedText>
                  </td>
                  <td className="py-3 px-4 font-mono text-xs text-neutral-200">
                    <TruncatedText>{generator.name}</TruncatedText>
                  </td>
                  <td className="py-3 px-4 font-mono text-xs text-neutral-400">
                    <span className="block max-w-64 truncate" title={generator.image}>
                      {generator.image || '\u2014'}
                    </span>
                  </td>
                  <td className="py-3 px-4 text-neutral-400">
                    <TruncatedText>{generator.contest_id}</TruncatedText>
                  </td>
                  <td className="py-3 px-4 text-neutral-400">
                    <TruncatedText>{generator.challenge_name || generator.challenge_id}</TruncatedText>
                  </td>
                  <td className="py-3 px-4 text-green-400">
                    <TruncatedText>{generator.success ?? 0}</TruncatedText>
                  </td>
                  <td className="py-3 px-4 text-neutral-400 text-xs">
                    <TruncatedText>{formatTime(generator.success_last)}</TruncatedText>
                  </td>
                  <td className="py-3 px-4 text-red-400">
                    <TruncatedText>{generator.failure ?? 0}</TruncatedText>
                  </td>
                  <td className="py-3 px-4 text-neutral-400 text-xs">
                    <TruncatedText>{formatTime(generator.failure_last)}</TruncatedText>
                  </td>
                  <td className="py-3 px-4">
                    <span
                      className={`inline-block max-w-full truncate px-2 py-0.5 rounded border text-xs font-mono ${STATUS_STYLES[generator.status] ?? STATUS_STYLES.stopped}`}
                      title={t(`admin.contests.generators.status.${generator.status}`, generator.status)}
                    >
                      {t(`admin.contests.generators.status.${generator.status}`, generator.status)}
                    </span>
                  </td>
                  <td className="py-3 px-4">
                    <RowActions actions={getRowActions(generator)} />
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
