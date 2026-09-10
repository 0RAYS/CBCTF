import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getAllWebhookHistory, getWebhookHistory } from '../../../../api/admin/webhook';
import { toast } from '../../../../utils/toast';
import AdminWebhookHistory from './AdminWebhookHistory';
import WebhookHistoryDialog from './WebhookHistoryDialog';

export default function WebhookHistory({ target }) {
  const { t } = useTranslation();
  const [histories, setHistories] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [selected, setSelected] = useState(null);
  const pageSize = 20;

  useEffect(() => {
    let active = true;
    const params = { limit: pageSize, offset: (currentPage - 1) * pageSize };
    const request = target ? getWebhookHistory(target.id, params) : getAllWebhookHistory(params);
    request
      .then((response) => {
        if (active && response.code === 200) {
          setHistories(response.data.histories || response.data);
          setTotalCount(response.data.count ?? response.data.length);
        }
      })
      .catch((error) => {
        if (active) toast.danger({ description: error.message || t('admin.webhook.toast.fetchHistoryFailed') });
      });
    return () => {
      active = false;
    };
  }, [currentPage, target?.id]);

  return (
    <>
      <AdminWebhookHistory
        webhookHistory={histories}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onViewDetail={setSelected}
        onHistoryClick={setSelected}
        webhookName={target?.name}
      />
      {selected && <WebhookHistoryDialog history={selected} onClose={() => setSelected(null)} />}
    </>
  );
}
