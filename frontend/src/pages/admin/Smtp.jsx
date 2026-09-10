import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getSmtpList } from '../../api/admin/smtp';
import { Tabs } from '../../components/common';
import AdminSmtp from '../../components/features/Admin/smtp/AdminSmtp';
import SmtpDialog from '../../components/features/Admin/smtp/SmtpDialog';
import EmailHistory from '../../components/features/Admin/smtp/EmailHistory';
import SmtpTestDialog from '../../components/features/Admin/smtp/SmtpTestDialog';
import { toast } from '../../utils/toast';

export default function SmtpManagement() {
  const { t } = useTranslation();
  const [smtpConfigs, setSmtpConfigs] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [activeTab, setActiveTab] = useState('smtp');
  const [historyTarget, setHistoryTarget] = useState(null);
  const [testTarget, setTestTarget] = useState(null);
  const [dialog, setDialog] = useState(null);
  const pageSize = 20;

  const fetchSmtpConfigs = async () => {
    try {
      const response = await getSmtpList({ limit: pageSize, offset: (currentPage - 1) * pageSize });
      if (response.code === 200) {
        setSmtpConfigs(response.data.smtps || response.data);
        setTotalCount(response.data.count ?? response.data.length);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.smtp.toast.fetchListFailed') });
    }
  };

  useEffect(() => {
    if (activeTab === 'smtp') fetchSmtpConfigs();
  }, [currentPage, activeTab]);

  const edit = (smtp) => setDialog({ mode: 'edit', smtp });
  const viewHistory = (smtp = null) => {
    setHistoryTarget(smtp);
    setActiveTab('history');
  };

  return (
    <>
      <Tabs
        value={activeTab}
        onChange={(next) => (next === 'smtp' ? setActiveTab(next) : viewHistory())}
        items={[
          { key: 'smtp', label: t('admin.smtp.tabs.config') },
          { key: 'history', label: t('admin.smtp.tabs.history') },
        ]}
      />
      {activeTab === 'smtp' ? (
        <AdminSmtp
          smtpConfigs={smtpConfigs}
          totalCount={totalCount}
          currentPage={currentPage}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onCreateSmtp={() => setDialog({ mode: 'create' })}
          onEditSmtp={edit}
          onDeleteSmtp={(smtp) => setDialog({ mode: 'delete', smtp })}
          onSmtpClick={edit}
          onViewHistory={viewHistory}
          onTestSmtp={setTestTarget}
        />
      ) : (
        <EmailHistory key={historyTarget?.id ?? 'all'} target={historyTarget} />
      )}
      {dialog && (
        <SmtpDialog
          {...dialog}
          onClose={() => setDialog(null)}
          onSaved={() => {
            setDialog(null);
            fetchSmtpConfigs();
          }}
        />
      )}
      {testTarget && (
        <SmtpTestDialog
          smtp={testTarget}
          onClose={() => setTestTarget(null)}
          onSent={() => {
            setTestTarget(null);
            fetchSmtpConfigs();
          }}
        />
      )}
    </>
  );
}
