import { useTranslation } from 'react-i18next';
import { Input, Select } from '../../../common';
import WebhookHeaders from './WebhookHeaders';

export default function WebhookFields({ form, setForm, mode, events }) {
  const { t } = useTranslation();
  const change = (key, value) => setForm((previous) => ({ ...previous, [key]: value }));
  const toggleEvent = (event) =>
    setForm((previous) => ({
      ...previous,
      events: previous.events.includes(event)
        ? previous.events.filter((item) => item !== event)
        : [...previous.events, event],
    }));

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Input
          label={t('admin.webhook.form.nameLabel')}
          value={form.name}
          onChange={(e) => change('name', e.target.value)}
          placeholder={t('admin.webhook.form.namePlaceholder')}
          fullWidth
          required={mode === 'create'}
        />
        <Select
          label={t('admin.webhook.form.methodLabel')}
          value={form.method}
          onChange={(e) => change('method', e.target.value)}
          options={[
            { value: 'GET', label: 'GET' },
            { value: 'POST', label: 'POST' },
          ]}
          required={mode === 'create'}
        />
      </div>
      <Input
        label={t('admin.webhook.form.urlLabel')}
        value={form.url}
        onChange={(e) => change('url', e.target.value)}
        placeholder={t('admin.webhook.form.urlPlaceholder')}
        fullWidth
        required={mode === 'create'}
      />
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <Input
          label={t('admin.webhook.form.timeoutLabel')}
          type="number"
          value={form.timeout}
          onChange={(e) => change('timeout', parseInt(e.target.value))}
          fullWidth
        />
        <Input
          label={t('admin.webhook.form.retryLabel')}
          type="number"
          value={form.retry}
          onChange={(e) => change('retry', parseInt(e.target.value))}
          fullWidth
        />
      </div>
      <WebhookHeaders headers={form.headers} setForm={setForm} />
      <div>
        <p className="text-neutral-300 text-sm font-medium mb-2">{t('admin.webhook.form.events')}</p>
        <div className="bg-neutral-800 border border-neutral-700 rounded-lg p-3 max-h-40 overflow-y-auto">
          {events.length > 0 ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
              {events.map((event) => (
                <label key={event} className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={form.events.includes(event)}
                    onChange={() => toggleEvent(event)}
                    className="text-geek-400"
                  />
                  <span className="text-neutral-300 text-sm break-all">{event}</span>
                </label>
              ))}
            </div>
          ) : (
            <p className="text-neutral-500 text-sm">{t('admin.webhook.form.noEvents')}</p>
          )}
        </div>
        <p className="text-neutral-400 text-xs mt-1">{t('admin.webhook.form.eventsHint')}</p>
      </div>
      <label className="flex items-center text-neutral-300">
        <input type="checkbox" checked={form.on} onChange={(e) => change('on', e.target.checked)} className="mr-2" />
        {t('admin.webhook.form.enable')}
      </label>
    </div>
  );
}
