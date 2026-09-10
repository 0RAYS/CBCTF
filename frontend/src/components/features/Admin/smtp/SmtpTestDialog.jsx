import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { testSmtp } from '../../../../api/admin/smtp';
import { toast } from '../../../../utils/toast';
import { Button, Input, Modal } from '../../../common';

export default function SmtpTestDialog({ smtp, onClose, onSent }) {
  const { t } = useTranslation();
  const [to, setTo] = useState('');
  const [sending, setSending] = useState(false);

  const send = async () => {
    if (!to || sending) return;
    setSending(true);
    try {
      const response = await testSmtp(smtp.id, { to });
      if (response.code === 200) {
        toast.success({ description: t('admin.smtp.toast.testSuccess') });
        onSent();
      } else {
        toast.danger({ description: response.msg || t('admin.smtp.toast.testFailed') });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.smtp.toast.testFailed') });
    } finally {
      setSending(false);
    }
  };

  return (
    <Modal
      isOpen
      onClose={onClose}
      title={t('admin.smtp.actions.testModalTitle')}
      size="sm"
      footer={
        <>
          <Button size="sm" variant="ghost" onClick={onClose} disabled={sending}>
            {t('common.cancel')}
          </Button>
          <Button size="sm" variant="primary" onClick={send} disabled={sending || !to}>
            {sending ? t('admin.smtp.actions.sending') : t('admin.smtp.actions.sendTest')}
          </Button>
        </>
      }
    >
      <div className="space-y-4">
        <div className="bg-neutral-800 border border-neutral-700 rounded-lg px-4 py-3 break-all">
          <p className="text-xs text-neutral-400 mb-1">{t('admin.smtp.columns.address')}</p>
          <p className="text-sm text-neutral-100 font-mono">{smtp.address}</p>
          <p className="text-xs text-neutral-400 mt-2 mb-1">{t('admin.smtp.columns.host')}</p>
          <p className="text-sm text-neutral-100 font-mono">
            {smtp.host}:{smtp.port}
          </p>
        </div>
        <Input
          label={t('admin.smtp.actions.testEmailLabel')}
          type="email"
          value={to}
          onChange={(e) => setTo(e.target.value)}
          placeholder={t('admin.smtp.actions.testEmailPlaceholder')}
          fullWidth
          disabled={sending}
          onKeyDown={(e) => {
            if (e.key === 'Enter') send();
          }}
        />
      </div>
    </Modal>
  );
}
