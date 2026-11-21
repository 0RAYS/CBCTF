import { useEffect, useEffectEvent, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconRotateClockwise } from '@tabler/icons-react';
import { Button, Spinner } from '@/components/common';
import ConfirmModal from '@/components/common/ConfirmModal';
import { getChallengeFlags, resetChallengeSolves, updateChallengeFlag } from '@/api/admin/challenge';
import { toast } from '@/utils/toast';
import FlagEditor from '@/components/features/Admin/challenges/editor/FlagEditor.jsx';
import FlagSolversDialog from '@/components/features/Admin/challenges/editor/FlagSolversDialog.jsx';
import { flagUpdatePayload, isFlagDirty } from '@/components/features/Admin/challenges/editor/challengeData.js';

export default function ChallengeFlagPanel({ challengeId, contestId, version = 0, onChanged }) {
  const { t } = useTranslation();
  const [flags, setFlags] = useState([]);
  const [savedFlags, setSavedFlags] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState(null);
  const [saving, setSaving] = useState(false);
  const [solvers, setSolvers] = useState(null);
  const [resetOpen, setResetOpen] = useState(false);
  const [resetting, setResetting] = useState(false);
  const busy = useRef(false);
  const alive = useRef(true);
  const fetchFailure = useEffectEvent(() => t('admin.challengeModal.toast.fetchFlagsFailed'));

  useEffect(() => {
    alive.current = true;
    let cancelled = false;
    setLoading(true);
    setLoadError(null);
    getChallengeFlags(challengeId)
      .then((response) => {
        if (cancelled) return;
        if (response.code !== 200) throw new Error(response.msg || fetchFailure());
        setFlags(response.data || []);
        setSavedFlags(response.data || []);
      })
      .catch((error) => {
        if (!cancelled) setLoadError(error.message);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
      alive.current = false;
    };
  }, [challengeId, version]);

  const changeFlag = (flag) => setFlags((previous) => previous.map((item) => (item.id === flag.id ? flag : item)));

  // Each flag is persisted independently. Refresh only its row so other flag drafts survive.
  const persistFlag = async (flag) => {
    const response = await updateChallengeFlag(challengeId, flag.id, flagUpdatePayload(flag));
    if (response.code !== 200) throw new Error(response.msg || t('admin.challengeModal.toast.updateFlagFailed'));
    if (!alive.current) return;
    setSavedFlags((previous) => previous.map((item) => (item.id === flag.id ? flag : item)));
    toast.success({
      description: t('admin.challengeModal.toast.updateFlagSuccess'),
    });
    const refresh = await getChallengeFlags(challengeId);
    if (!alive.current || refresh.code !== 200) return;
    const updated = refresh.data?.find((item) => item.id === flag.id);
    if (updated) {
      changeFlag(updated);
      setSavedFlags((previous) => previous.map((item) => (item.id === flag.id ? updated : item)));
    }
  };

  const save = async (onlyFlag) => {
    if (busy.current || loading || loadError) return;
    busy.current = true;
    setSaving(true);
    try {
      if (onlyFlag) {
        await persistFlag(onlyFlag);
        if (alive.current) onChanged?.();
        return;
      }
      for (const flag of flags) {
        if (!alive.current) return;
        if (
          isFlagDirty(
            flag,
            savedFlags.find((item) => item.id === flag.id)
          )
        )
          await persistFlag(flag);
      }
      if (alive.current) onChanged?.();
    } catch (error) {
      if (alive.current) toast.danger({ description: error.message });
    } finally {
      busy.current = false;
      if (alive.current) setSaving(false);
    }
  };

  // 重置解题记录是破坏性操作，仅在此处显式触发
  const resetSolves = async () => {
    if (resetting) return;
    setResetting(true);
    try {
      const response = await resetChallengeSolves(challengeId);
      if (response.code !== 200) throw new Error(response.msg || t('admin.challengeModal.toast.resetSolvesFailed'));
      if (!alive.current) return;
      setResetOpen(false);
      toast.success({
        description: t('admin.challengeModal.toast.resetSolvesSuccess', {
          teams: response.data?.cleared_teams ?? 0,
          submissions: response.data?.cleared_submissions ?? 0,
        }),
      });
      const refresh = await getChallengeFlags(challengeId);
      if (!alive.current || refresh.code !== 200) return;
      setFlags(refresh.data || []);
      setSavedFlags(refresh.data || []);
      onChanged?.();
    } catch (error) {
      if (alive.current) toast.danger({ description: error.message });
    } finally {
      if (alive.current) setResetting(false);
    }
  };

  return (
    <>
      <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
        <h3 className="text-lg font-mono text-neutral-50">{t('admin.challengeModal.sections.scoring')}</h3>
        <Button
          size="sm"
          variant="danger"
          icon={<IconRotateClockwise size={15} />}
          disabled={saving || loading}
          onClick={() => setResetOpen(true)}
        >
          {t('admin.challengeModal.actions.resetSolves')}
        </Button>
      </div>
      {loading ? (
        <Spinner />
      ) : loadError ? (
        <p role="alert" className="text-red-400">
          {loadError}
        </p>
      ) : flags.length === 0 ? (
        <p className="text-sm font-mono text-neutral-500">{t('admin.challengeModal.scoring.empty')}</p>
      ) : (
        flags.map((flag, index) => (
          <FlagEditor
            key={flag.id}
            flag={flag}
            index={index}
            onChange={changeFlag}
            onSave={() => save(flag)}
            onViewSolvers={() => setSolvers({ flagId: flag.id, index })}
            disabled={saving}
            dirty={isFlagDirty(
              flag,
              savedFlags.find((item) => item.id === flag.id)
            )}
          />
        ))
      )}
      {solvers && (
        <FlagSolversDialog
          key={solvers.flagId}
          isOpen
          onClose={() => setSolvers(null)}
          flagIndex={solvers.index}
          flagId={solvers.flagId}
          contestId={contestId}
          challengeId={challengeId}
        />
      )}
      <ConfirmModal
        isOpen={resetOpen}
        onClose={() => setResetOpen(false)}
        onConfirm={resetSolves}
        type="danger"
        loading={resetting}
        title={t('admin.challengeModal.resetSolves.title')}
        confirmText={t('admin.challengeModal.resetSolves.confirm')}
      >
        <p className="text-sm text-neutral-300">{t('admin.challengeModal.resetSolves.description')}</p>
      </ConfirmModal>
    </>
  );
}
