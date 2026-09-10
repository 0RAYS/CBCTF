import { useTranslation } from 'react-i18next';
import { IconPlus, IconTrash } from '@tabler/icons-react';
import { Button, Card, Input, Textarea } from '../../../../common';

export default function Prizes({ prizes = [], onChange }) {
  const { t } = useTranslation();
  const update = (index, field, value) =>
    onChange(prizes.map((prize, i) => (i === index ? { ...prize, [field]: value } : prize)));
  return (
    <Card padding="lg">
      <div className="flex flex-wrap justify-between gap-3 mb-6">
        <h2 className="text-2xl font-mono text-neutral-50">{t('game.detail.labels.prizes')}</h2>
        <Button
          size="sm"
          icon={<IconPlus size={16} />}
          onClick={() => onChange([...prizes, { amount: '$0', description: '' }])}
        >
          {t('admin.contests.editor.actions.addPrize')}
        </Button>
      </div>
      <div className="space-y-6">
        {prizes.map((prize, index) => (
          <div key={index} className="p-4 border border-neutral-700 rounded-md space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-geek-400 font-mono">
                {t(`admin.contests.editor.prizeRank.${['first', 'second', 'third'][index] || 'other'}`, {
                  rank: index + 1,
                })}
              </span>
              <Button
                size="icon"
                variant="ghost"
                className="text-red-400"
                onClick={() => onChange(prizes.filter((_, i) => i !== index))}
              >
                <IconTrash size={18} />
              </Button>
            </div>
            <Input
              label={t('admin.contests.editor.labels.prizeAmount')}
              value={prize.amount}
              onChange={(e) => update(index, 'amount', e.target.value)}
            />
            <Textarea
              label={t('admin.contests.editor.labels.prizeDescription')}
              value={prize.description || ''}
              onChange={(e) => update(index, 'description', e.target.value)}
            />
          </div>
        ))}
      </div>
    </Card>
  );
}
