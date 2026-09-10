import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { createSmtp, updateSmtp, deleteSmtp } from '../../../../api/admin/smtp';
import { toast } from '../../../../utils/toast';
import { Button, Modal } from '../../../common';
import { buildSmtpPayload, smtpForm } from './payload';
import SmtpFields from './SmtpFields';

export default function SmtpDialog({ mode, smtp, onClose, onSaved }) {
  const { t } = useTranslation();
  const [form, setForm] = useState(() => smtpForm(smtp));

  const submit = async () => {
    const action = mode === 'edit' ? 'update' : mode;
    try {
      const response =
        mode === 'delete'
          ? await deleteSmtp(smtp.id)
          : mode === 'create'
            ? await createSmtp(buildSmtpPayload(form))
            : await updateSmtp(smtp.id, buildSmtpPayload(form, smtp));
      if (response.code === 200) {
        toast.success({ description: t(`admin.smtp.toast.${action}Success`) });
        onSaved();
      }
    } catch (error) {
      toast.danger({ description: error.message || t(`admin.smtp.toast.${action}Failed`) });
    }
  };

  return (
    <Modal
      isOpen
      onClose={onClose}
      title={t(`admin.smtp.modal.${mode}Title`)}
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
            {t('admin.smtp.modal.deletePrompt')} <span className="font-semibold text-red-400">{smtp.address}</span>?
          </p>
          <p className="text-neutral-400 text-sm">{t('admin.smtp.modal.deleteWarning')}</p>
        </div>
      ) : (
        <SmtpFields form={form} setForm={setForm} mode={mode} />
      )}
    </Modal>
  );
}
