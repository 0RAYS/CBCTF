import { useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import AdminNotice from '../../../components/features/Admin/Contests/AdminNotice';
import {
  createContestNotice,
  deleteContestNotice,
  getContestNotices,
  updateContestNotice,
} from '../../../api/admin/contest';
import { useTranslation } from 'react-i18next';

function AdminContestNotices() {
  const { id: contestId } = useParams();
  const [notices, setNotices] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [loadedPage, setLoadedPage] = useState(0);
  const requestRef = useRef(0);
  const pageSize = 20;
  const { t } = useTranslation();
  // 页码与已展示数据的页码不一致时才算加载中，刷新不会清空表格
  const loading = loadedPage !== currentPage;

  const fetchNotices = async () => {
    const version = ++requestRef.current;
    try {
      const response = await getContestNotices(parseInt(contestId), {
        limit: pageSize,
        offset: (currentPage - 1) * pageSize,
      });
      if (version !== requestRef.current || response.code !== 200) return;
      setNotices(response.data.notices);
      setTotalCount(response.data.count);
    } catch (error) {
      if (version === requestRef.current)
        toast.danger({ description: error.message || t('common.notices.fetchFailed') });
    } finally {
      if (version === requestRef.current) setLoadedPage(currentPage);
    }
  };

  useEffect(() => {
    fetchNotices();
    return () => {
      requestRef.current += 1;
    };
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
      loading={loading}
      onPageChange={setCurrentPage}
      onCreateNotice={handleCreateNotice}
      onUpdateNotice={handleUpdateNotice}
      onDeleteNotice={handleDeleteNotice}
    />
  );
}

export default AdminContestNotices;
