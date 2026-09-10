import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconPlayerPlay, IconBan, IconCheck, IconX, IconTrash, IconRefresh } from '@tabler/icons-react';
import { Button, Pagination, StatCard } from '../../../common';
import AutoRefreshControl from '../../../common/AutoRefreshControl';
import GeneratorList from './GeneratorList';
import GeneratorStartDialog from './GeneratorStartDialog';
import GeneratorLogDialog from './GeneratorLogDialog';
import useGeneratorSession from './useGeneratorSession.js';
import { GENERATOR_PAGE_SIZE, getGeneratorPageStats } from './generatorUtils.js';

export default function GeneratorManagement({ api, textKey }) {
  const { t } = useTranslation();
  const text = (key, options) => t(`${textKey}.${key}`, options);
  const session = useGeneratorSession(api, text);
  const [logGenerator, setLogGenerator] = useState(null);
  const stats = getGeneratorPageStats(session.generators);

  return (
    <div className="flex flex-col min-w-0 gap-6 p-6">
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <StatCard title={text('stats.total')} value={session.totalCount} icon={<IconPlayerPlay size={20} />} />
        <StatCard title={text('stats.successes')} value={stats.successes} icon={<IconCheck size={20} />} />
        <StatCard title={text('stats.failures')} value={stats.failures} icon={<IconX size={20} />} />
      </div>

      <div className="flex flex-wrap gap-2 items-center">
        <AutoRefreshControl value={session.refreshInterval} onChange={session.setRefreshInterval} />
        <Button variant="ghost" size="sm" onClick={session.refresh} icon={<IconRefresh size={14} />}>
          {t('common.refresh')}
        </Button>
        <Button variant="primary" size="sm" onClick={session.openStart} icon={<IconPlayerPlay size={14} />}>
          {text('startButton')}
        </Button>
        <Button
          variant={session.showDeleted ? 'danger' : 'ghost'}
          size="sm"
          onClick={session.toggleShowDeleted}
          icon={<IconTrash size={14} />}
        >
          {text('showDeleted')}
        </Button>
        {session.selectedIds.length > 0 && (
          <Button
            variant="danger"
            size="sm"
            onClick={session.stop}
            disabled={Boolean(session.pendingOperation)}
            loading={session.pendingOperation === 'stop'}
            icon={<IconBan size={14} />}
          >
            {text('stopButton')} ({session.selectedIds.length})
          </Button>
        )}
      </div>

      <GeneratorList session={session} onViewLogs={setLogGenerator} text={text} t={t} />
      {session.totalCount > GENERATOR_PAGE_SIZE && (
        <Pagination
          current={session.currentPage}
          total={Math.ceil(session.totalCount / GENERATOR_PAGE_SIZE)}
          totalItems={session.totalCount}
          showTotal
          onChange={session.changePage}
        />
      )}

      {session.startModalOpen && (
        <GeneratorStartDialog
          challenges={session.dynamicChallenges}
          pendingOperation={session.pendingOperation}
          onStart={session.start}
          onClose={session.closeStart}
          text={text}
          t={t}
        />
      )}
      {logGenerator && (
        <GeneratorLogDialog
          key={logGenerator.id}
          api={api}
          generator={logGenerator}
          onClose={() => setLogGenerator(null)}
          text={text}
          t={t}
        />
      )}
    </div>
  );
}
