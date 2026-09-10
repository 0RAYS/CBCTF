import { IconPlus, IconTrash } from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Input } from '../../../common';
import { renameHeader } from './payload';

export default function WebhookHeaders({ headers, setForm }) {
  const { t } = useTranslation();
  const add = () =>
    setForm((previous) => {
      let counter = 1;
      while (Object.hasOwn(previous.headers, `header_${counter}`)) counter++;
      return { ...previous, headers: { ...previous.headers, [`header_${counter}`]: '' } };
    });
  const remove = (key) =>
    setForm((previous) => {
      const next = { ...previous.headers };
      delete next[key];
      return { ...previous, headers: next };
    });

  return (
    <div>
      <div className="flex justify-between items-center gap-2 mb-3">
        <p className="text-neutral-300 text-sm font-medium">{t('admin.webhook.form.headers')}</p>
        <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={add}>
          {t('admin.webhook.form.addHeader')}
        </Button>
      </div>
      <div className="space-y-2">
        {Object.entries(headers).map(([key, value], index) => (
          <div key={index} className="flex gap-2 items-center">
            <HeaderKeyInput
              name={key}
              onRename={(nextKey) => {
                const renamed = renameHeader(headers, key, nextKey);
                // Commit on blur so an immediate Save includes the new key, using the latest values.
                setForm((previous) => ({ ...previous, headers: renameHeader(previous.headers, key, nextKey) }));
                return renamed === headers ? key : nextKey.trim();
              }}
            />
            <Input
              value={value}
              onChange={(e) => {
                const nextValue = e.target.value;
                setForm((previous) => ({ ...previous, headers: { ...previous.headers, [key]: nextValue } }));
              }}
              placeholder={t('admin.webhook.form.headerValuePlaceholder')}
              aria-label={t('admin.webhook.form.headerValuePlaceholder')}
              className="flex-1 min-w-0"
            />
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-red-400 hover:!text-red-300"
              onClick={() => remove(key)}
              aria-label={t('common.delete')}
            >
              <IconTrash size={18} />
            </Button>
          </div>
        ))}
      </div>
    </div>
  );
}

function HeaderKeyInput({ name, onRename }) {
  const { t } = useTranslation();
  const [draft, setDraft] = useState(name);
  useEffect(() => {
    setDraft(name);
  }, [name]);
  return (
    <Input
      value={draft}
      onChange={(e) => setDraft(e.target.value)}
      onBlur={() => setDraft(onRename(draft))}
      placeholder={t('admin.webhook.form.headerKeyPlaceholder')}
      aria-label={t('admin.webhook.form.headerKeyPlaceholder')}
      className="flex-1 min-w-0"
    />
  );
}
