import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconDeviceFloppy } from '@tabler/icons-react';
import { Button } from '../../../../common';
import Identity from './Identity';
import Schedule from './Schedule';
import Rules from './Rules';
import Prizes from './Prizes';
import Timeline from './Timeline';
import { createContestDraft, validateContestDraft } from './contestForm.js';

export default function ContestEditor({ contest: initialContest, onSave, onCancel, onImageUpload }) {
  const { t } = useTranslation();
  // The draft belongs to this editor session, not to picture-upload refreshes.
  const [contest, setContest] = useState(() => initialContest || createContestDraft());
  const [errors, setErrors] = useState({});
  const [saving, setSaving] = useState(false);
  const change = (field, value) => {
    setContest((previous) => ({ ...previous, [field]: value }));
    setErrors((previous) => ({ ...previous, [field]: undefined }));
  };
  const submit = async (event) => {
    event.preventDefault();
    if (saving) return;
    const validation = validateContestDraft(contest);
    setErrors(validation);
    const firstField = Object.keys(validation)[0];
    if (firstField) {
      event.currentTarget.elements.namedItem(firstField)?.focus();
      return;
    }
    setSaving(true);
    try {
      await onSave(contest);
    } finally {
      setSaving(false);
    }
  };
  return (
    <form onSubmit={submit} className="w-full mx-auto space-y-6">
      <Identity contest={contest} onChange={change} errors={errors} onImageUpload={onImageUpload} />
      <Schedule contest={contest} onChange={change} errors={errors} />
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Rules rules={contest.rules} onChange={(rules) => change('rules', rules)} />
        <Prizes prizes={contest.prizes} onChange={(prizes) => change('prizes', prizes)} />
      </div>
      <Timeline timeline={contest.timeline} onChange={(timeline) => change('timeline', timeline)} />
      <div className="flex flex-wrap justify-end gap-4">
        <Button variant="ghost" size="sm" onClick={onCancel}>
          {t('common.cancel')}
        </Button>
        <Button type="submit" variant="primary" size="sm" disabled={saving} icon={<IconDeviceFloppy size={18} />}>
          {t('common.saveChanges')}
        </Button>
      </div>
    </form>
  );
}
