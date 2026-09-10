import { useTranslation } from 'react-i18next';
import { Button, Modal } from '../../../common';

export default function WebhookHistoryDialog({ history, onClose }) {
  const { t } = useTranslation();
  const fields = [
    ['webhook', history.webhook],
    ['event', history.event],
    [
      'status',
      <span className={history.success ? 'text-green-400' : 'text-red-400'} key="status">
        {t(`admin.webhook.history.${history.success ? 'statusSuccess' : 'statusFailed'}`)}
      </span>,
    ],
    ['responseCode', history.resp || '-'],
    ['duration', history.duration ? `${history.duration}ms` : '-'],
    ['retry', history.retry],
  ];

  return (
    <Modal
      isOpen
      onClose={onClose}
      title={t('admin.webhook.history.detailTitle')}
      size="lg"
      footer={
        <Button size="sm" variant="ghost" onClick={onClose}>
          {t('admin.webhook.actions.close')}
        </Button>
      }
    >
      <div className="space-y-4">
        <dl className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {fields.map(([key, value]) => (
            <div key={key}>
              <dt className="text-neutral-300 text-sm font-medium mb-2">{t(`admin.webhook.history.${key}`)}</dt>
              <dd className="text-neutral-50 break-all">{value}</dd>
            </div>
          ))}
        </dl>
        {history.error && (
          <div>
            <p className="text-neutral-300 text-sm font-medium mb-2">{t('admin.webhook.history.error')}</p>
            <div className="bg-red-900/20 border border-red-700 rounded-lg p-3">
              <pre className="text-red-300 text-sm whitespace-pre-wrap break-words font-mono">{history.error}</pre>
            </div>
          </div>
        )}
      </div>
    </Modal>
  );
}
