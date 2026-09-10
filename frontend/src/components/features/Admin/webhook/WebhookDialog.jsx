import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { createWebhook, updateWebhook, deleteWebhook, getEvents } from '../../../../api/admin/webhook';
import { toast } from '../../../../utils/toast';
import { Button, Modal } from '../../../common';
import { buildWebhookPayload, webhookForm } from './payload';
import WebhookFields from './WebhookFields';

export default function WebhookDialog({ mode, webhook, onClose, onSaved }) {
  const { t } = useTranslation();
  const [form, setForm] = useState(() => webhookForm(webhook));
  const [events, setEvents] = useState([]);

  useEffect(() => {
    if (mode === 'delete') return;
    let active = true;
    getEvents()
      .then((response) => {
        if (active && response.code === 200) setEvents(response.data || []);
      })
      .catch((error) => {
        if (active) toast.danger({ description: error.message || t('admin.webhook.toast.fetchEventsFailed') });
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
          ? await deleteWebhook(webhook.id)
          : mode === 'create'
            ? await createWebhook(buildWebhookPayload(form))
            : await updateWebhook(webhook.id, buildWebhookPayload(form, webhook));
      if (response.code === 200) {
        toast.success({ description: t(`admin.webhook.toast.${action}Success`) });
        onSaved();
      }
    } catch (error) {
      toast.danger({ description: error.message || t(`admin.webhook.toast.${action}Failed`) });
    }
  };

  return (
    <Modal
      isOpen
      onClose={onClose}
      title={t(`admin.webhook.modal.${mode}Title`)}
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
            {t('admin.webhook.modal.deletePrompt')} <span className="font-semibold text-red-400">{webhook.name}</span>?
          </p>
          <p className="text-neutral-400 text-sm">{t('admin.webhook.modal.deleteWarning')}</p>
        </div>
      ) : (
        <WebhookFields form={form} setForm={setForm} mode={mode} events={events} />
      )}
    </Modal>
  );
}
