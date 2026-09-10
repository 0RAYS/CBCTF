import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { createOAuthProvider, updateOAuthProvider, deleteOAuthProvider } from '../../../../api/admin/oauth';
import { getGroupList } from '../../../../api/admin/rbac';
import { toast } from '../../../../utils/toast';
import { Button, Modal, Select } from '../../../common';
import { buildOAuthPayload, oauthForm } from './payload';
import OAuthProtocolFields from './OAuthProtocolFields';
import OAuthClaimFields from './OAuthClaimFields';

export default function OAuthProviderDialog({ mode, provider, onClose, onSaved }) {
  const { t } = useTranslation();
  const [form, setForm] = useState(() => oauthForm(provider));
  const [groups, setGroups] = useState([]);
  const change = (key, value) => setForm((previous) => ({ ...previous, [key]: value }));

  useEffect(() => {
    if (mode === 'delete') return;
    let active = true;
    getGroupList({ limit: 100, offset: 0 })
      .then((response) => {
        if (active && response.code === 200) setGroups(response.data.groups);
      })
      .catch((error) => {
        if (active) toast.danger({ description: error.message || t('admin.oauthProviders.toast.fetchFailed') });
      });
    return () => {
      active = false;
    };
  }, [mode]);

  const submit = async () => {
    const action = mode === 'edit' ? 'update' : mode;
    try {
      const response =
        mode === 'delete'
          ? await deleteOAuthProvider(provider.id)
          : mode === 'create'
            ? await createOAuthProvider(buildOAuthPayload(form))
            : await updateOAuthProvider(provider.id, buildOAuthPayload(form));
      if (response.code === 200) {
        toast.success({ description: t(`admin.oauthProviders.toast.${action}Success`) });
        onSaved();
      }
    } catch (error) {
      toast.danger({ description: error.message || t(`admin.oauthProviders.toast.${action}Failed`) });
    }
  };

  return (
    <Modal
      isOpen
      onClose={onClose}
      title={t(`admin.oauthProviders.modal.${mode}Title`)}
      size={mode === 'delete' ? 'sm' : 'lg'}
      footer={
        <>
          <Button size="sm" variant="ghost" onClick={onClose}>
            {t('common.cancel')}
          </Button>
          <Button size="sm" variant={mode === 'delete' ? 'danger' : 'primary'} onClick={submit}>
            {t(`common.${mode === 'edit' ? 'save' : mode}`)}
          </Button>
        </>
      }
    >
      {mode === 'delete' ? (
        <div className="text-center">
          <p className="text-neutral-300 mb-4">
            {t('admin.oauthProviders.modal.deletePrompt')}{' '}
            <span className="font-semibold text-red-400">{provider.provider}</span>?
          </p>
          <p className="text-neutral-400 text-sm">{t('admin.oauthProviders.modal.deleteWarning')}</p>
        </div>
      ) : (
        <div className="space-y-3">
          <OAuthProtocolFields form={form} change={change} mode={mode} />
          <OAuthClaimFields form={form} change={change} mode={mode} />
          <Select
            label={t('admin.oauthProviders.form.defaultGroupLabel')}
            value={form.default_group}
            onChange={(e) => change('default_group', Number(e.target.value))}
            options={[
              { value: 0, label: t('admin.oauthProviders.form.defaultGroupNone') },
              ...groups.map((group) => ({ value: group.id, label: group.name })),
            ]}
          />
          <label className="flex items-center text-neutral-300">
            <input
              type="checkbox"
              checked={form.on}
              onChange={(e) => change('on', e.target.checked)}
              className="mr-2"
            />
            {t('admin.oauthProviders.form.enable')}
          </label>
        </div>
      )}
    </Modal>
  );
}
