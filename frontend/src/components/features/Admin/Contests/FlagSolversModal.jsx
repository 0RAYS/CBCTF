import { useEffect, useState } from 'react';
import { motion } from 'motion/react';
import { IconX, IconTrophy } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../../components/common';
import { getFlagSolvers } from '../../../../api/admin/challenge';
import { toast } from '../../../../utils/toast';
import { useUserDetailDialog } from '../../../../hooks/useUserDetailDialog.jsx';
import { useTeamDetailDialog } from '../../../../hooks/useTeamDetailDialog.jsx';

function FlagSolversModal({ isOpen, onClose, flagIndex, contestId, challengeId, flagId }) {
  const { t } = useTranslation();
  const [solvers, setSolvers] = useState([]);
  const [loading, setLoading] = useState(false);

  const { openUserDetail, renderUserDetailDialog } = useUserDetailDialog();
  const { openTeamDetail, renderTeamDetailDialog } = useTeamDetailDialog(contestId);

  useEffect(() => {
    if (!isOpen || !flagId) return;
    setLoading(true);
    setSolvers([]);
    getFlagSolvers(contestId, challengeId, flagId)
      .then((res) => {
        if (res.code === 200) {
          setSolvers(res.data?.solvers ?? []);
        }
      })
      .catch((err) => {
        toast.danger({ description: err.message });
      })
      .finally(() => setLoading(false));
  }, [isOpen, flagId, contestId, challengeId]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center bg-black/70 backdrop-blur-sm p-4">
      <motion.div
        className="w-full max-w-2xl bg-neutral-900 border border-neutral-700 rounded-md overflow-hidden"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.2 }}
      >
        {/* Header */}
        <div className="flex justify-between items-center p-4 border-b border-neutral-700">
          <h2 className="text-lg font-mono text-neutral-50">
            {t('admin.contests.challengeModal.solversModal.title', { index: flagIndex + 1 })}
          </h2>
          <Button
            variant="ghost"
            size="icon"
            className="!text-neutral-400 hover:!text-neutral-300"
            aria-label={t('common.close')}
            onClick={onClose}
          >
            <IconX size={18} />
          </Button>
        </div>

        {/* Body */}
        <div className="p-4 max-h-[60vh] overflow-y-auto">
          {loading ? (
            <p className="text-center text-sm font-mono text-neutral-400 py-8">
              {t('admin.contests.challengeModal.solversModal.loading')}
            </p>
          ) : solvers.length === 0 ? (
            <p className="text-center text-sm font-mono text-neutral-500 py-8">
              {t('admin.contests.challengeModal.solversModal.empty')}
            </p>
          ) : (
            <table className="w-full text-sm font-mono">
              <thead>
                <tr className="border-b border-neutral-700 text-neutral-400">
                  <th className="text-left py-2 pr-4 w-12" scope="col">
                    {t('admin.contests.challengeModal.solversModal.columns.rank')}
                  </th>
                  <th className="text-left py-2 pr-4" scope="col">
                    {t('admin.contests.challengeModal.solversModal.columns.user')}
                  </th>
                  <th className="text-left py-2 pr-4" scope="col">
                    {t('admin.contests.challengeModal.solversModal.columns.team')}
                  </th>
                  <th className="text-right py-2 pr-4 w-20" scope="col">
                    {t('admin.contests.challengeModal.solversModal.columns.score')}
                  </th>
                  <th className="text-right py-2 w-36" scope="col">
                    {t('admin.contests.challengeModal.solversModal.columns.solvedAt')}
                  </th>
                </tr>
              </thead>
              <tbody>
                {solvers.map((solver, i) => (
                  <tr key={i} className="border-b border-neutral-800 hover:bg-white/5 transition-colors">
                    <td className="py-2 pr-4">
                      {i === 0 ? (
                        <IconTrophy size={16} className="text-yellow-400" />
                      ) : i === 1 ? (
                        <IconTrophy size={16} className="text-neutral-300" />
                      ) : i === 2 ? (
                        <IconTrophy size={16} className="text-amber-600" />
                      ) : (
                        <span className="text-neutral-400">{i + 1}</span>
                      )}
                    </td>
                    <td className="py-2 pr-4">
                      {solver.user_id ? (
                        <button
                          className="text-neutral-200 hover:text-geek-400 transition-colors cursor-pointer text-left"
                          onClick={() => openUserDetail(solver.user_id)}
                        >
                          {solver.user_name || '—'}
                        </button>
                      ) : (
                        <span className="text-neutral-200">{solver.user_name || '—'}</span>
                      )}
                    </td>
                    <td className="py-2 pr-4">
                      {solver.team_id ? (
                        <button
                          className="text-neutral-300 hover:text-geek-400 transition-colors cursor-pointer text-left"
                          onClick={() => openTeamDetail(solver.team_id)}
                        >
                          {solver.team_name || '—'}
                        </button>
                      ) : (
                        <span className="text-neutral-300">{solver.team_name || '—'}</span>
                      )}
                    </td>
                    <td className="py-2 pr-4 text-right text-geek-400">{solver.score}</td>
                    <td className="py-2 text-right text-neutral-400 text-xs">
                      {solver.solved_at ? new Date(solver.solved_at).toLocaleString() : '—'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </motion.div>

      {renderUserDetailDialog()}
      {renderTeamDetailDialog()}
    </div>
  );
}

export default FlagSolversModal;
