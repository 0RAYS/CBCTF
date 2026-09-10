import { useTranslation } from 'react-i18next';
import { Button, Modal } from '../../../common';

export default function EmailHistoryDialog({ email, onClose }) {
  const { t, i18n } = useTranslation();

  return (
    <Modal
      isOpen
      onClose={onClose}
      title={t('admin.smtp.email.detailTitle')}
      size="lg"
      footer={
        <Button size="sm" variant="ghost" onClick={onClose}>
          {t('admin.smtp.actions.close')}
        </Button>
      }
    >
      <dl className="space-y-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <dt className="text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.id')}</dt>
            <dd className="text-neutral-50 font-mono">#{email.id}</dd>
          </div>
          <div>
            <dt className="text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.status')}</dt>
            <dd className={email.success ? 'text-green-400' : 'text-red-400'}>
              {email.success ? t('admin.smtp.email.statusSuccess') : t('admin.smtp.email.statusFailed')}
            </dd>
          </div>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <dt className="text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.from')}</dt>
            <dd className="text-neutral-50 break-all">{email.from}</dd>
          </div>
          <div>
            <dt className="text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.to')}</dt>
            <dd className="text-neutral-50 break-all">{email.to}</dd>
          </div>
        </div>
        <div>
          <dt className="text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.subject')}</dt>
          <dd className="text-neutral-50 break-words">{email.subject || t('admin.smtp.email.noSubject')}</dd>
        </div>
        <div>
          <dt className="text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.time')}</dt>
          <dd className="text-neutral-50">{new Date(email.time).toLocaleString(i18n.language || 'en-US')}</dd>
        </div>
        <div>
          <dt className="text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.content')}</dt>
          <dd className="bg-neutral-800 border border-neutral-700 rounded-lg p-3 max-h-60 overflow-y-auto">
            <pre className="text-neutral-300 text-sm whitespace-pre-wrap break-words font-mono">
              {email.content || t('admin.smtp.email.noContent')}
            </pre>
          </dd>
        </div>
      </dl>
    </Modal>
  );
}
