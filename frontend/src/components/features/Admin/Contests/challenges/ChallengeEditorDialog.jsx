import { useEffect, useEffectEvent, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconPlus, IconTrash } from '@tabler/icons-react';
import { Button, Input, Modal, Spinner, Textarea } from '../../../../common';
import Checkbox from '../../../../common/Checkbox';
import { getChallengeFlags, updateChallengeFlag } from '../../../../../api/admin/challenge';
import { updateContestChallenge } from '../../../../../api/admin/contest';
import { toast } from '../../../../../utils/toast';
import FlagEditor from './FlagEditor';
import FlagSolversDialog from './FlagSolversDialog';
import { challengeDraft, flagUpdatePayload, isFlagDirty } from './challengeData.js';

export default function ChallengeEditorDialog({ contestId, challenge, onClose, onSaved }) {
  const { t } = useTranslation();
  const [draft, setDraft] = useState(() => challengeDraft(challenge));
  const [flags, setFlags] = useState([]);
  const [savedFlags, setSavedFlags] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState(null);
  const [saving, setSaving] = useState(false);
  const [solvers, setSolvers] = useState(null);
  const busy = useRef(false);
  const alive = useRef(true);
  const fetchFailure = useEffectEvent(() => t('admin.contests.challenges.toast.fetchListFailed'));
  useEffect(() => {
    alive.current = true;
    let cancelled = false;
    getChallengeFlags(contestId, challenge.id)
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
  }, [contestId, challenge.id]);
  const change = (field, value) => setDraft((previous) => ({ ...previous, [field]: value }));
  const changeFlag = (flag) => setFlags((previous) => previous.map((item) => (item.id === flag.id ? flag : item)));

  // Each flag is persisted independently. Refresh only its row so other flag drafts survive.
  const persistFlag = async (flag) => {
    const response = await updateChallengeFlag(contestId, challenge.id, flag.id, flagUpdatePayload(flag));
    if (response.code !== 200) throw new Error(response.msg || t('admin.contests.challenges.toast.updateFailed'));
    if (!alive.current) return;
    setSavedFlags((previous) => previous.map((item) => (item.id === flag.id ? flag : item)));
    toast.success({ description: t('admin.contests.challenges.toast.updateFlagSuccess') });
    const refresh = await getChallengeFlags(contestId, challenge.id);
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
        if (alive.current) onSaved();
      } else {
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
        if (!alive.current) return;
        const response = await updateContestChallenge(contestId, challenge.id, draft);
        if (response.code !== 200) throw new Error(response.msg || t('admin.contests.challenges.toast.updateFailed'));
        if (!alive.current) return;
        toast.success({ description: t('admin.contests.challenges.toast.updateSuccess') });
        onSaved();
        onClose();
      }
    } catch (error) {
      if (alive.current) toast.danger({ description: error.message });
    } finally {
      busy.current = false;
      if (alive.current) setSaving(false);
    }
  };
  const close = () => {
    if (!busy.current) onClose();
  };
  return (
    <>
      <Modal
        isOpen
        onClose={close}
        title={t('admin.contests.challengeModal.titleEdit')}
        size="2xl"
        footer={
          <>
            <Button size="sm" variant="ghost" onClick={close} disabled={saving}>
              {t('common.cancel')}
            </Button>
            <Button size="sm" variant="primary" onClick={() => save()} disabled={saving || loading || !!loadError}>
              {t('common.saveChanges')}
            </Button>
          </>
        }
      >
        <div className="grid grid-cols-1 xl:grid-cols-2 gap-4 [&_*]:min-w-0">
          <fieldset disabled={saving} className="space-y-4">
            <legend className="text-lg font-mono text-neutral-50">
              {t('admin.contests.challengeModal.sections.basic')}
            </legend>
            <Input
              label={t('admin.contests.challengeModal.labels.name')}
              value={draft.name}
              onChange={(e) => change('name', e.target.value)}
            />
            <Textarea
              label={t('admin.contests.challengeModal.labels.description')}
              value={draft.description}
              onChange={(e) => change('description', e.target.value)}
              rows={4}
            />
            <Input
              label={t('admin.contests.challengeModal.labels.attempts')}
              type="number"
              value={draft.attempt}
              onChange={(e) => change('attempt', parseInt(e.target.value) || 0)}
            />
            <Checkbox
              checked={draft.hidden}
              onChange={(e) => change('hidden', e.target.checked)}
              label={t('admin.contests.challengeModal.labels.hidden')}
            />
            {['tags', 'hints'].map((field) => {
              const singular = field === 'tags' ? 'Tag' : 'Hint';
              return (
                <section key={field} className="space-y-3 border-t border-neutral-700 pt-4">
                  <div className="flex flex-wrap justify-between gap-2">
                    <h3 className="text-lg font-mono text-neutral-50">
                      {t(`admin.contests.challengeModal.sections.${field}`)}
                    </h3>
                    <Button
                      size="sm"
                      variant="primary"
                      icon={<IconPlus size={16} />}
                      onClick={() => change(field, [...draft[field], ''])}
                    >
                      {t(`admin.contests.challengeModal.actions.add${singular}`)}
                    </Button>
                  </div>
                  {draft[field].map((value, index) => (
                    <div key={index} className="flex items-center gap-2">
                      <Input
                        aria-label={t(`admin.contests.challengeModal.placeholders.${singular.toLowerCase()}`, {
                          index: index + 1,
                        })}
                        value={value}
                        placeholder={t(`admin.contests.challengeModal.placeholders.${singular.toLowerCase()}`, {
                          index: index + 1,
                        })}
                        onChange={(e) =>
                          change(
                            field,
                            draft[field].map((item, i) => (i === index ? e.target.value : item))
                          )
                        }
                      />
                      <Button
                        size="icon"
                        variant="ghost"
                        onClick={() =>
                          change(
                            field,
                            draft[field].filter((_, i) => i !== index)
                          )
                        }
                      >
                        <IconTrash size={18} />
                      </Button>
                    </div>
                  ))}
                </section>
              );
            })}
          </fieldset>
          <section className="space-y-4">
            <h3 className="text-lg font-mono text-neutral-50">{t('admin.contests.challengeModal.sections.flags')}</h3>
            {loading ? (
              <Spinner />
            ) : loadError ? (
              <p role="alert" className="text-red-400">
                {loadError}
              </p>
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
          </section>
        </div>
      </Modal>
      {solvers && (
        <FlagSolversDialog
          key={solvers.flagId}
          isOpen
          onClose={() => setSolvers(null)}
          flagIndex={solvers.index}
          flagId={solvers.flagId}
          contestId={contestId}
          challengeId={challenge.id}
        />
      )}
    </>
  );
}
