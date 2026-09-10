import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getWebhookList } from '../../api/admin/webhook';
import { Tabs } from '../../components/common';
import AdminWebhook from '../../components/features/Admin/webhook/AdminWebhook';
import WebhookDialog from '../../components/features/Admin/webhook/WebhookDialog';
import WebhookHistory from '../../components/features/Admin/webhook/WebhookHistory';
import { toast } from '../../utils/toast';

export default function WebhookManagement() {
  const { t } = useTranslation();
  const [webhooks, setWebhooks] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [activeTab, setActiveTab] = useState('webhook');
  const [historyTarget, setHistoryTarget] = useState(null);
  const [dialog, setDialog] = useState(null);
  const pageSize = 20;

  const fetchWebhooks = async () => {
    try {
      const response = await getWebhookList({ limit: pageSize, offset: (currentPage - 1) * pageSize });
      if (response.code === 200) {
        setWebhooks(response.data.webhooks || response.data);
        setTotalCount(response.data.count ?? response.data.length);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.webhook.toast.fetchListFailed') });
    }
  };

  useEffect(() => {
    if (activeTab === 'webhook') fetchWebhooks();
  }, [currentPage, activeTab]);

  const edit = (webhook) => setDialog({ mode: 'edit', webhook });
  const viewHistory = (webhook = null) => {
    setHistoryTarget(webhook);
    setActiveTab('history');
  };

  return (
    <>
      <Tabs
        value={activeTab}
        onChange={(next) => (next === 'webhook' ? setActiveTab(next) : viewHistory())}
        items={[
          { key: 'webhook', label: t('admin.webhook.tabs.webhook') },
          { key: 'history', label: t('admin.webhook.tabs.history') },
        ]}
      />
      {activeTab === 'webhook' ? (
        <AdminWebhook
          webhooks={webhooks}
          totalCount={totalCount}
          currentPage={currentPage}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onCreateWebhook={() => setDialog({ mode: 'create' })}
          onEditWebhook={edit}
          onDeleteWebhook={(webhook) => setDialog({ mode: 'delete', webhook })}
          onViewHistory={viewHistory}
          onWebhookClick={edit}
        />
      ) : (
        <WebhookHistory key={historyTarget?.id ?? 'all'} target={historyTarget} />
      )}
      {dialog && (
        <WebhookDialog
          {...dialog}
          onClose={() => setDialog(null)}
          onSaved={() => {
            setDialog(null);
            fetchWebhooks();
          }}
        />
      )}
    </>
  );
}
