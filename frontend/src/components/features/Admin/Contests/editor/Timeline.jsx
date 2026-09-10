import { useTranslation } from 'react-i18next';
import { IconPlus, IconX } from '@tabler/icons-react';
import { Button, Card, DateTimeInput, Input, Textarea } from '../../../../common';
import { formatDateForInput } from './contestForm.js';

export default function Timeline({ timeline = [], onChange }) {
  const { t } = useTranslation();
  const update = (index, field, value) =>
    onChange(timeline.map((item, i) => (i === index ? { ...item, [field]: value } : item)));
  return (
    <Card padding="lg">
      <div className="flex flex-wrap justify-between gap-3 mb-6">
        <h2 className="text-2xl font-mono text-neutral-50">{t('admin.contests.editor.sections.timeline')}</h2>
        <Button
          size="sm"
          icon={<IconPlus size={16} />}
          onClick={() => onChange([...timeline, { date: '', title: '', description: '' }])}
        >
          {t('admin.contests.editor.actions.addTimeline')}
        </Button>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {timeline.map((item, index) => (
          <div key={index} className="p-4 border border-neutral-700 rounded-md space-y-4">
            <div className="flex justify-end">
              <Button
                size="icon"
                variant="ghost"
                className="text-red-400"
                onClick={() => onChange(timeline.filter((_, i) => i !== index))}
              >
                <IconX size={18} />
              </Button>
            </div>
            <DateTimeInput
              label={t('admin.contests.editor.labels.timelineDate')}
              value={formatDateForInput(item.date)}
              onChange={(e) => update(index, 'date', e.target.value)}
            />
            <Input
              label={t('admin.contests.editor.labels.timelineTitle')}
              value={item.title}
              onChange={(e) => update(index, 'title', e.target.value)}
            />
            <Textarea
              label={t('admin.contests.editor.labels.timelineDescription')}
              value={item.description}
              onChange={(e) => update(index, 'description', e.target.value)}
            />
          </div>
        ))}
      </div>
    </Card>
  );
}
