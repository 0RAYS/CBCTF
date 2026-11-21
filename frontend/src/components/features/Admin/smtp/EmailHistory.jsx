import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getAllEmailHistory, getEmailHistory } from '@/api/admin/email';
import { toast } from '@/utils/toast';
import AdminEmailHistory from '@/components/features/Admin/smtp/AdminEmailHistory';
import EmailHistoryDialog from '@/components/features/Admin/smtp/EmailHistoryDialog';

export default function EmailHistory({ target }) {
  const { t } = useTranslation();
  const [emails, setEmails] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [selected, setSelected] = useState(null);
  const [loading, setLoading] = useState(true);
  const pageSize = 20;

  useEffect(() => {
    let active = true;
    setLoading(true);
    const params = { limit: pageSize, offset: (currentPage - 1) * pageSize };
    const request = target ? getEmailHistory(target.id, params) : getAllEmailHistory(params);
    request
      .then((response) => {
        if (active && response.code === 200) {
          setEmails(response.data.emails || response.data);
          setTotalCount(response.data.count ?? response.data.length);
        }
      })
      .catch((error) => {
        if (active) toast.danger({ description: error.message || t('admin.smtp.toast.fetchHistoryFailed') });
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [currentPage, target?.id]);

  return (
    <>
      <AdminEmailHistory
        emailHistory={emails}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        loading={loading}
        onPageChange={setCurrentPage}
        onViewEmail={setSelected}
        onEmailClick={setSelected}
        smtpAddress={target?.address}
      />
      {selected && <EmailHistoryDialog email={selected} onClose={() => setSelected(null)} />}
    </>
  );
}
