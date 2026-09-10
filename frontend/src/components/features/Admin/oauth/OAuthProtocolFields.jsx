import { useTranslation } from 'react-i18next';
import { Input, Select } from '../../../common';

export default function OAuthProtocolFields({ form, change, mode }) {
  const { t } = useTranslation();
  const creating = mode === 'create';

  return (
    <>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Input
          label={t('admin.oauthProviders.form.providerLabel')}
          value={form.provider}
          onChange={(e) => change('provider', e.target.value)}
          placeholder={t('admin.oauthProviders.form.providerPlaceholder')}
          fullWidth
          required={creating}
        />
        <Input
          label={t('admin.oauthProviders.form.uriLabel')}
          value={form.uri}
          onChange={(e) => change('uri', e.target.value)}
          placeholder={t('admin.oauthProviders.form.uriPlaceholder')}
          fullWidth
          required={creating}
        />
        <div className="sm:col-span-2">
          <Select
            label={t('admin.oauthProviders.form.protocolLabel')}
            value={form.protocol}
            onChange={(e) => change('protocol', e.target.value)}
            options={[
              { value: 'oauth2', label: t('admin.oauthProviders.form.protocolOAuth2') },
              { value: 'cas', label: t('admin.oauthProviders.form.protocolCAS') },
            ]}
          />
        </div>
        {form.protocol === 'oauth2' && (
          <div className="sm:col-span-2">
            <Input
              label={t('admin.oauthProviders.form.scopesLabel')}
              value={form.scopes}
              onChange={(e) => change('scopes', e.target.value)}
              placeholder={t('admin.oauthProviders.form.scopesPlaceholder')}
              fullWidth
            />
          </div>
        )}
      </div>
      <div className="space-y-3">
        <Input
          label={t('admin.oauthProviders.form.authUrlLabel')}
          value={form.auth_url}
          onChange={(e) => change('auth_url', e.target.value)}
          placeholder={t('admin.oauthProviders.form.authUrlPlaceholder')}
          fullWidth
          required={creating}
        />
        {form.protocol === 'oauth2' && (
          <Input
            label={t('admin.oauthProviders.form.tokenUrlLabel')}
            value={form.token_url}
            onChange={(e) => change('token_url', e.target.value)}
            placeholder={t('admin.oauthProviders.form.tokenUrlPlaceholder')}
            fullWidth
            required={creating}
          />
        )}
        <Input
          label={t('admin.oauthProviders.form.userInfoUrlLabel')}
          value={form.user_info_url}
          onChange={(e) => change('user_info_url', e.target.value)}
          placeholder={t('admin.oauthProviders.form.userInfoUrlPlaceholder')}
          fullWidth
          required={creating}
        />
        <Input
          label={t('admin.oauthProviders.form.callbackUrlLabel')}
          value={form.callback_url}
          onChange={(e) => change('callback_url', e.target.value)}
          placeholder={t('admin.oauthProviders.form.callbackUrlPlaceholder')}
          fullWidth
          required={creating}
        />
      </div>
      {form.protocol === 'oauth2' && (
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <Input
            label={t('admin.oauthProviders.form.clientIdLabel')}
            value={form.client_id}
            onChange={(e) => change('client_id', e.target.value)}
            placeholder={t('admin.oauthProviders.form.clientIdPlaceholder')}
            fullWidth
            required={creating}
          />
          <Input
            label={t('admin.oauthProviders.form.clientSecretLabel')}
            type="password"
            value={form.client_secret}
            onChange={(e) => change('client_secret', e.target.value)}
            placeholder={
              mode === 'edit' ? t('common.leaveBlankToKeep') : t('admin.oauthProviders.form.clientSecretPlaceholder')
            }
            fullWidth
            required={creating}
          />
        </div>
      )}
    </>
  );
}
