import { useId, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { updateContestCheat } from '../../../../api/admin/contest';
import { toast } from '../../../../utils/toast';
import { Button, FormField, Modal, Select, Textarea } from '../../../common';
import Checkbox from '../../../common/Checkbox';
import { cheatReviewForm, cheatUpdatePayload } from './payloads';

export default function CheatReviewDialog({ contestId, cheat, onClose, onSaved }) {
  const { t } = useTranslation();
  const id = useId();
  const [form, setForm] = useState(() => cheatReviewForm(cheat));
  const [pending, setPending] = useState(false);
  const busy = useRef(false);
  const close = () => {
    if (!busy.current) onClose();
  };

  async function save() {
    if (busy.current) return;
    const payload = cheatUpdatePayload(form, cheat);
    if (Object.keys(payload).length === 0) {
      onClose();
      return;
    }
    busy.current = true;
    setPending(true);
    try {
      const response = await updateContestCheat(contestId, cheat.id, payload);
      if (response.code !== 200) throw new Error(t('admin.contests.cheats.toast.updateFailed'));
      toast.success({ description: t('admin.contests.cheats.toast.updateSuccess') });
      onSaved();
      onClose();
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.cheats.toast.updateFailed') });
    } finally {
      busy.current = false;
      setPending(false);
    }
  }

  return (
    <Modal
      isOpen
      onClose={close}
      title={t('admin.contests.cheats.modals.editTitle')}
      size="md"
      footer={
        <>
          <Button size="sm" variant="ghost" onClick={close} disabled={pending}>
            {t('common.cancel')}
          </Button>
          <Button size="sm" variant="primary" onClick={save} disabled={pending}>
            {t('common.save')}
          </Button>
        </>
      }
    >
      <div className="space-y-4">
        <FormField label={t('admin.contests.cheats.form.verdict')} htmlFor={`${id}-verdict`}>
          <Select
            id={`${id}-verdict`}
            value={form.type}
            disabled={pending}
            onChange={(event) => setForm({ ...form, type: event.target.value })}
            options={['cheater', 'suspicious', 'pass'].map((type) => ({
              value: type,
              label: t(`admin.contests.cheats.types.${type}`),
            }))}
            fullWidth
          />
        </FormField>
        <FormField label={t('admin.contests.cheats.form.reason')} htmlFor={`${id}-reason`}>
          <Textarea
            id={`${id}-reason`}
            value={form.reason}
            disabled={pending}
            onChange={(event) => setForm({ ...form, reason: event.target.value })}
            placeholder={t('admin.contests.cheats.form.reasonPlaceholder')}
            rows={3}
            fullWidth
          />
        </FormField>
        <Checkbox
          label={t('admin.contests.cheats.status.processed')}
          checked={form.checked}
          disabled={pending}
          onChange={(event) => setForm({ ...form, checked: event.target.checked })}
        />
        <FormField label={t('admin.contests.cheats.form.comment')} htmlFor={`${id}-comment`}>
          <Textarea
            id={`${id}-comment`}
            value={form.comment}
            disabled={pending}
            onChange={(event) => setForm({ ...form, comment: event.target.value })}
            placeholder={t('admin.contests.cheats.form.commentPlaceholder')}
            rows={3}
            fullWidth
          />
        </FormField>
      </div>
    </Modal>
  );
}
