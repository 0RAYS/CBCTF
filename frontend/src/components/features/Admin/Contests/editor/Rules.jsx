import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconEdit, IconTrash, IconPlus } from '@tabler/icons-react';
import { Button, Card, Input } from '../../../../common';

export default function Rules({ rules = [], onChange }) {
  const { t } = useTranslation();
  const [newRule, setNewRule] = useState('');
  const [editing, setEditing] = useState(null);
  return (
    <Card padding="lg" className="lg:col-span-2">
      <h2 className="text-2xl font-mono text-neutral-50 mb-6">{t('admin.contests.editor.sections.rules')}</h2>
      <div className="flex gap-2 mb-6">
        <Input
          value={newRule}
          onChange={(e) => setNewRule(e.target.value)}
          placeholder={t('admin.contests.editor.placeholders.newRule')}
        />
        <Button
          size="sm"
          icon={<IconPlus size={18} />}
          onClick={() => {
            if (!newRule.trim()) return;
            onChange([...rules, newRule]);
            setNewRule('');
          }}
        >
          {t('admin.contests.editor.actions.addRule')}
        </Button>
      </div>
      <div className="space-y-4">
        {rules.map((rule, index) => (
          <div key={index} className="flex items-center gap-3 text-neutral-300">
            <span className="text-geek-400 font-mono">{String(index + 1).padStart(2, '0')}</span>
            {editing === index ? (
              <Input
                autoFocus
                value={rule}
                onChange={(e) => onChange(rules.map((item, i) => (i === index ? e.target.value : item)))}
                onBlur={() => setEditing(null)}
              />
            ) : (
              <p className="flex-1 min-w-0 break-words">{rule}</p>
            )}
            <Button size="icon" variant="ghost" onClick={() => setEditing(index)}>
              <IconEdit size={18} />
            </Button>
            <Button
              size="icon"
              variant="ghost"
              className="text-red-400"
              onClick={() => {
                onChange(rules.filter((_, i) => i !== index));
                setEditing(null);
              }}
            >
              <IconTrash size={18} />
            </Button>
          </div>
        ))}
      </div>
    </Card>
  );
}
