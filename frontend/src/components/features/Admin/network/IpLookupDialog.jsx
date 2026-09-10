import { useTranslation } from 'react-i18next';
import Modal from '../../../common/Modal';
import Card from '../../../common/Card';

const FIELDS = ['ip', 'iso', 'country', 'subdivision', 'city', 'timezone'];

export default function IpLookupDialog({ isOpen, onClose, data, loading }) {
  const { t } = useTranslation();
  const text = (key) => t(`admin.contests.cheats.ipDetail.${key}`);

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={text('title')} size="sm">
      {loading ? (
        <div className="flex items-center justify-center py-8" role="status" aria-label={t('common.loading')}>
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-geek-400" />
        </div>
      ) : data ? (
        <Card variant="default" padding="md">
          <dl className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {FIELDS.map((field) => (
              <div key={field} className="min-w-0">
                <dt className="block text-sm font-mono text-neutral-400 mb-1">{text(field)}</dt>
                <dd className={`text-neutral-300 break-words ${field === 'ip' ? 'font-mono' : ''}`}>
                  {data[field] || '-'}
                </dd>
              </div>
            ))}
            <div className="sm:col-span-2 min-w-0">
              <dt className="block text-sm font-mono text-neutral-400 mb-1">{text('coordinates')}</dt>
              <dd className="text-neutral-300 font-mono break-words">
                {data.latitude != null && data.longitude != null ? `${data.latitude}, ${data.longitude}` : '-'}
              </dd>
            </div>
          </dl>
        </Card>
      ) : null}
    </Modal>
  );
}
