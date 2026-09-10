import { useTranslation } from 'react-i18next';
import { Input } from '../../../common';

export default function SmtpFields({ form, setForm, mode }) {
  const { t } = useTranslation();
  const change = (key, value) => setForm((previous) => ({ ...previous, [key]: value }));
  const creating = mode === 'create';

  return (
    <div className="space-y-3">
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Input
          label={t('admin.smtp.form.addressLabel')}
          type="email"
          value={form.address}
          onChange={(e) => change('address', e.target.value)}
          placeholder={t('admin.smtp.form.addressPlaceholder')}
          fullWidth
          required={creating}
        />
        <Input
          label={t('admin.smtp.form.hostLabel')}
          value={form.host}
          onChange={(e) => change('host', e.target.value)}
          placeholder={t('admin.smtp.form.hostPlaceholder')}
          fullWidth
          required={creating}
        />
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Input
          label={t('admin.smtp.form.portLabel')}
          type="number"
          value={form.port}
          onChange={(e) => change('port', parseInt(e.target.value) || 587)}
          placeholder={t('admin.smtp.form.portPlaceholder')}
          fullWidth
          required={creating}
        />
        <Input
          label={creating ? t('admin.smtp.form.passwordLabelCreate') : t('admin.smtp.form.passwordLabelEdit')}
          type="password"
          value={form.pwd}
          onChange={(e) => change('pwd', e.target.value)}
          placeholder={
            creating ? t('admin.smtp.form.passwordPlaceholderCreate') : t('admin.smtp.form.passwordPlaceholderEdit')
          }
          fullWidth
          required={creating}
        />
      </div>
      <label className="flex items-center text-neutral-300">
        <input type="checkbox" checked={form.on} onChange={(e) => change('on', e.target.checked)} className="mr-2" />
        {t('admin.smtp.form.enable')}
      </label>
    </div>
  );
}
