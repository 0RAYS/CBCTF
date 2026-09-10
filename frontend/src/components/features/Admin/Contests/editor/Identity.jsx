import { useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconUpload } from '@tabler/icons-react';
import { Button, Card, Input, Textarea } from '../../../../common';
import Checkbox from '../../../../common/Checkbox';

export default function Identity({ contest, onChange, errors, onImageUpload }) {
  const { t } = useTranslation();
  const input = useRef(null);
  const [uploading, setUploading] = useState(false);
  const upload = async (event) => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file || uploading) return;
    setUploading(true);
    try {
      const image = onImageUpload
        ? await onImageUpload(file)
        : await new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = () => resolve(reader.result);
            reader.onerror = () => reject(reader.error);
            reader.readAsDataURL(file);
          });
      if (image) onChange('image', image);
    } finally {
      setUploading(false);
    }
  };
  return (
    <Card padding="lg" className="space-y-6">
      <Input
        label={t('admin.contests.editor.labels.title')}
        name="title"
        value={contest.title}
        onChange={(e) => onChange('title', e.target.value)}
        error={errors.title ? t(`admin.contests.editor.validation.${errors.title}`) : undefined}
        required
      />
      <Textarea
        label={t('admin.contests.editor.labels.description')}
        value={contest.description}
        onChange={(e) => onChange('description', e.target.value)}
        rows={4}
      />
      <fieldset className="min-w-0">
        <legend className="block text-sm font-medium text-neutral-400 mb-1">
          {t('admin.contests.editor.labels.backgroundImage')}
        </legend>
        {contest.image && (
          <img
            src={contest.image}
            alt={t('admin.contests.editor.labels.backgroundAlt')}
            className="w-full h-48 object-cover rounded-md mb-3"
          />
        )}
        <Button size="sm" disabled={uploading} icon={<IconUpload size={16} />} onClick={() => input.current?.click()}>
          {t(`admin.contests.editor.actions.${contest.image ? 'replaceImage' : 'uploadImage'}`)}
        </Button>
        <input
          ref={input}
          type="file"
          aria-label={t('admin.contests.editor.labels.backgroundImage')}
          accept="image/png,image/jpeg,image/jpg,image/gif"
          className="hidden"
          onChange={upload}
        />
      </fieldset>
      <h3 className="text-xl font-mono text-neutral-50">{t('admin.contests.editor.sections.advanced')}</h3>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[
          ['prefix', 'flagPrefix', 'text'],
          ['size', 'teamSize', 'number'],
          ['victims', 'teamVictims', 'number'],
          ['captcha', 'captcha', 'text'],
        ].map(([field, label, type]) => (
          <div key={field}>
            <Input
              label={t(`admin.contests.editor.labels.${label}`)}
              name={field}
              type={type}
              value={contest[field] ?? ''}
              min={type === 'number' ? 1 : undefined}
              max={type === 'number' ? 10 : undefined}
              onChange={(e) => onChange(field, type === 'number' ? parseInt(e.target.value) || 1 : e.target.value)}
              placeholder={t(`admin.contests.editor.placeholders.${label}`)}
            />
            <p className="mt-1 text-neutral-500 text-sm">{t(`admin.contests.editor.help.${label}`)}</p>
          </div>
        ))}
        {[
          ['hidden', 'hiddenContest'],
          ['blood', 'bloodBonus'],
        ].map(([field, label]) => (
          <div key={field} className="min-w-0">
            <Checkbox
              checked={!!contest[field]}
              onChange={(e) => onChange(field, e.target.checked)}
              label={t(`admin.contests.editor.labels.${label}`)}
            />
            <p className="mt-1 text-neutral-500 text-sm">{t(`admin.contests.editor.help.${label}`)}</p>
          </div>
        ))}
      </div>
    </Card>
  );
}
