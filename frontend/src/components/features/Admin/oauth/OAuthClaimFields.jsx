import { useTranslation } from 'react-i18next';
import { Input } from '../../../common';

export default function OAuthClaimFields({ form, change, mode }) {
  const { t } = useTranslation();

  return (
    <div className="space-y-3">
      <h3 className="text-neutral-200 font-medium">{t('admin.oauthProviders.form.mappingTitle')}</h3>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Input
          label={t('admin.oauthProviders.form.pictureClaimLabel')}
          value={form.picture_claim}
          onChange={(e) => change('picture_claim', e.target.value)}
          placeholder="{picture_url}"
          fullWidth
        />
        <Input
          label={t('admin.oauthProviders.form.nameClaimLabel')}
          value={form.name_claim}
          onChange={(e) => change('name_claim', e.target.value)}
          placeholder="{name}"
          fullWidth
          required={mode === 'create'}
        />
        <Input
          label={t('admin.oauthProviders.form.emailClaimLabel')}
          value={form.email_claim}
          onChange={(e) => change('email_claim', e.target.value)}
          placeholder="{email}"
          fullWidth
          required={mode === 'create'}
        />
        <Input
          label={t('admin.oauthProviders.form.descriptionClaimLabel')}
          value={form.description_claim}
          onChange={(e) => change('description_claim', e.target.value)}
          placeholder="{description}"
          fullWidth
        />
        <Input
          label={t('admin.oauthProviders.form.idClaimLabel')}
          value={form.id_claim}
          onChange={(e) => change('id_claim', e.target.value)}
          placeholder="{id}"
          fullWidth
          required={mode === 'create'}
        />
        <Input
          label={t('admin.oauthProviders.form.logoUrlLabel')}
          value={form.picture}
          onChange={(e) => change('picture', e.target.value)}
          placeholder={t('admin.oauthProviders.form.logoUrlPlaceholder')}
          fullWidth
        />
        <Input
          label={t('admin.oauthProviders.form.groupsClaimLabel')}
          value={form.groups_claim}
          onChange={(e) => change('groups_claim', e.target.value)}
          placeholder="{groups}"
          fullWidth
        />
        <Input
          label={t('admin.oauthProviders.form.adminGroupLabel')}
          value={form.admin_group}
          onChange={(e) => change('admin_group', e.target.value)}
          placeholder="admin"
          fullWidth
        />
      </div>
    </div>
  );
}
