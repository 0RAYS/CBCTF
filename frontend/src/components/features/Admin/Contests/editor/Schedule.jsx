import { useTranslation } from 'react-i18next';
import { Card, DateTimeInput, Select } from '../../../../common';
import { formatDateForInput } from './contestForm.js';

export default function Schedule({ contest, onChange, errors }) {
  const { t } = useTranslation();
  return (
    <Card padding="lg" className="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Select
        label={t('admin.contests.editor.labels.status')}
        value={contest.status}
        disabled
        options={['upcoming', 'active', 'ended'].map((status) => ({
          value: status,
          label: t(`admin.contests.editor.status.${status}`),
        }))}
      />
      {['startTime', 'endTime'].map((field) => (
        <DateTimeInput
          key={field}
          label={t(`admin.contests.editor.labels.${field}`)}
          name={field}
          value={formatDateForInput(contest[field])}
          onChange={(e) => onChange(field, e.target.value ? new Date(e.target.value).toISOString() : '')}
          error={errors[field] ? t(`admin.contests.editor.validation.${errors[field]}`) : undefined}
        />
      ))}
    </Card>
  );
}
