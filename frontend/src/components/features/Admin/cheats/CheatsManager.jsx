import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { motion } from 'motion/react';
import { IconFilter, IconShieldCheck, IconTrash } from '@tabler/icons-react';
import { getContestCheats } from '../../../../api/admin/contest';
import { toast } from '../../../../utils/toast';
import { Button, StatCard } from '../../../common';
import { useUserDetailDialog } from '../details/useUserDetailDialog.jsx';
import { useTeamDetailDialog } from '../details/useTeamDetailDialog.jsx';
import CheatEvidenceDialog from './CheatEvidenceDialog';
import CheatReviewDialog from './CheatReviewDialog';
import CheatActionDialog from './CheatActionDialog';
import CheatsTable from './CheatsTable';
import { cheatListParams } from './payloads';

export default function CheatsManager({ contestId }) {
  const { t } = useTranslation();
  const [cheats, setCheats] = useState([]);
  const [total, setTotal] = useState(0);
  const [checked, setChecked] = useState(0);
  const [page, setPage] = useState(1);
  const [filterType, setFilterType] = useState('');
  const [filterReasonType, setFilterReasonType] = useState('');
  const [showFilter, setShowFilter] = useState(false);
  const [loading, setLoading] = useState(true);
  const [revision, setRevision] = useState(0);
  const [dialog, setDialog] = useState(null);
  const pageSize = 20;
  const { openUserDetail, renderUserDetailDialog } = useUserDetailDialog();
  const { openTeamDetail, renderTeamDetailDialog } = useTeamDetailDialog(contestId);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    getContestCheats(contestId, cheatListParams(page, pageSize, filterType, filterReasonType))
      .then((response) => {
        if (response.code !== 200) throw new Error(t('admin.contests.cheats.toast.fetchFailed'));
        if (cancelled) return;
        const count = response.data.count || 0;
        setCheats(response.data.cheats || []);
        setChecked(response.data.checked || 0);
        setTotal(count);
        setPage(Math.min(page, Math.max(1, Math.ceil(count / pageSize))));
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message || t('admin.contests.cheats.toast.fetchFailed') });
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [contestId, page, filterType, filterReasonType, revision, t]);

  const refresh = () => setRevision((value) => value + 1);
  const closeDialog = () => setDialog(null);
  const openAction = (action, cheat = null) => setDialog({ action, cheat });
  const filterRows = [
    {
      label: 'filterLabel',
      prefix: 'types',
      values: ['cheater', 'suspicious', 'pass'],
      value: filterType,
      setValue: setFilterType,
    },
    {
      label: 'reasonTypeFilterLabel',
      prefix: 'reasonTypes',
      values: ['same_web_ip', 'same_victim_ip', 'wrong_flag'],
      value: filterReasonType,
      setValue: setFilterReasonType,
    },
  ];

  return (
    <div className="w-full mx-auto space-y-6">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        <div className="flex flex-wrap items-center justify-end gap-2 mb-6">
          <Button size="sm" variant="ghost" onClick={() => openAction('check')} className="bg-black/30!">
            <IconShieldCheck size={16} />
            {t('admin.contests.cheats.actions.check')}
          </Button>
          <Button
            size="sm"
            variant="ghost"
            onClick={() => openAction('deleteAll')}
            className="bg-black/30! text-red-400! hover:text-red-300!"
          >
            <IconTrash size={16} />
            {t('admin.contests.cheats.actions.deleteAll')}
          </Button>
          <Button
            size="sm"
            variant="ghost"
            onClick={() => setShowFilter(!showFilter)}
            className={`bg-black/30! ${showFilter ? 'bg-geek-400/20!' : ''}`}
          >
            <IconFilter size={16} />
            {t('admin.contests.cheats.filterButton')}
          </Button>
        </div>
        {showFilter && (
          <motion.div
            className="border border-neutral-600 rounded-md bg-neutral-900 p-4 mb-6 space-y-3"
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
          >
            {filterRows.map((filter) => (
              <div key={filter.label} className="flex flex-col sm:flex-row sm:items-center gap-4">
                <span className="text-neutral-400 font-mono">{t(`admin.contests.cheats.${filter.label}`)}</span>
                <div className="flex flex-wrap gap-2">
                  {filter.values.map((value) => (
                    <Button
                      key={value}
                      size="sm"
                      variant={filter.value === value ? 'primary' : 'ghost'}
                      onClick={() => {
                        filter.setValue(value);
                        setPage(1);
                      }}
                      className={filter.value === value ? '' : 'bg-black/30!'}
                    >
                      {t(`admin.contests.cheats.${filter.prefix}.${value}`)}
                    </Button>
                  ))}
                </div>
              </div>
            ))}
          </motion.div>
        )}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <StatCard title={t('admin.contests.cheats.stats.total')} value={total} valueColor="text-neutral-50" />
          <StatCard
            title={t('admin.contests.cheats.stats.processed')}
            value={checked}
            valueColor="text-green-400"
            delay={0.1}
          />
          <StatCard
            title={t('admin.contests.cheats.stats.unprocessed')}
            value={total - checked}
            valueColor="text-yellow-400"
            delay={0.2}
          />
        </div>
        <CheatsTable
          cheats={cheats}
          total={total}
          page={page}
          pageSize={pageSize}
          loading={loading}
          onPageChange={setPage}
          onAction={openAction}
          openUserDetail={openUserDetail}
          openTeamDetail={openTeamDetail}
        />
      </motion.div>
      {dialog?.action === 'detail' && (
        <CheatEvidenceDialog
          cheat={dialog.cheat}
          onClose={closeDialog}
          openUserDetail={openUserDetail}
          openTeamDetail={openTeamDetail}
        />
      )}
      {dialog?.action === 'edit' && (
        <CheatReviewDialog
          key={dialog.cheat.id}
          contestId={contestId}
          cheat={dialog.cheat}
          onClose={closeDialog}
          onSaved={refresh}
        />
      )}
      {['delete', 'deleteAll', 'check'].includes(dialog?.action) && (
        <CheatActionDialog
          key={dialog.action}
          contestId={contestId}
          action={dialog.action}
          cheat={dialog.cheat}
          onClose={closeDialog}
          onCompleted={refresh}
        />
      )}
      {renderUserDetailDialog()}
      {renderTeamDetailDialog()}
    </div>
  );
}
