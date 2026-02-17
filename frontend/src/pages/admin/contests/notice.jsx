import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import AdminNotice from '../../../components/features/Admin/Contests/AdminNotice';
import {
  getContestNotices,
  createContestNotice,
  updateContestNotice,
  deleteContestNotice,
} from '../../../api/admin/contest';
import { useTranslation } from 'react-i18next';

function AdminContestNotices() {
  const { id: contestId } = useParams();
  const [notices, setNotices] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;
  const { t } = useTranslation();

  const fetchNotices = async () => {
    try {
      const response = await getContestNotices(parseInt(contestId), {
        limit: pageSize,
        offset: (currentPage - 1) * pageSize,
      });
      if (response.code === 200) {
        setNotices(response.data.notices);
        setTotalCount(response.data.count);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.notices.toast.fetchFailed') });
    }
  };

  useEffect(() => {
    fetchNotices();
  }, [contestId, currentPage]);

  const handleCreateNotice = async (form) => {
    try {
      const response = await createContestNotice(parseInt(contestId), form);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.notices.toast.createSuccess') });
        fetchNotices();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.notices.toast.createFailed') });
    }
  };

  const handleUpdateNotice = async (noticeId, form) => {
    try {
      const response = await updateContestNotice(parseInt(contestId), noticeId, form);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.notices.toast.updateSuccess') });
        fetchNotices();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.notices.toast.updateFailed') });
    }
  };

  const handleDeleteNotice = async (noticeId) => {
    try {
      const response = await deleteContestNotice(parseInt(contestId), noticeId);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.notices.toast.deleteSuccess') });
        fetchNotices();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.notices.toast.deleteFailed') });
    }
  };

  return (
    <AdminNotice
      notices={notices}
      totalCount={totalCount}
      currentPage={currentPage}
      pageSize={pageSize}
      onPageChange={setCurrentPage}
      onCreateNotice={handleCreateNotice}
      onUpdateNotice={handleUpdateNotice}
      onDeleteNotice={handleDeleteNotice}
    />
  );
}

export default AdminContestNotices;
