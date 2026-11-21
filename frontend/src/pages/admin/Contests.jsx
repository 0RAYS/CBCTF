import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from '@/utils/toast';
import AdminContests from '@/components/features/Admin/AdminContests';
import { DateTimeInput, Input, Modal, Textarea } from '@/components/common';
import CRUDModalFooter from '@/components/common/CRUDModalFooter';
import DeleteConfirmation from '@/components/common/DeleteConfirmation';
import { createContest, deleteContest, getContestList, updateContestPicture } from '@/api/admin/contest';
import { formatDateForInput } from '@/components/features/Admin/Contests/editor/contestForm.js';
import { useTranslation } from 'react-i18next';

const CONTEST_DURATION_HOURS = 48;

// The picker edits local calendar fields, so the default must be rendered in local time too.
const defaultStartTime = () => formatDateForInput(new Date(Date.now() + 24 * 60 * 60 * 1000));

const defaultCreateForm = () => ({
  name: '',
  description: '',
  captcha: '',
  prefix: 'CBCTF',
  size: 4,
  start: defaultStartTime(),
  duration: CONTEST_DURATION_HOURS,
});

function ContestsManagement() {
  const [contests, setContests] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [mode, setMode] = useState('create'); // 'create' | 'delete'
  const [selectedContest, setSelectedContest] = useState(null);
  const [createForm, setCreateForm] = useState(defaultCreateForm);

  const fileInputRef = useRef(null);
  const pageSize = 10;
  const navigate = useNavigate();
  const { t } = useTranslation();

  const fetchContests = async () => {
    setLoading(true);
    try {
      const response = await getContestList({
        limit: pageSize,
        offset: (currentPage - 1) * pageSize,
      });
      if (response.code === 200) {
        setContests(response.data.contests);
        setTotalCount(response.data.count);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.toast.fetchFailed') });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchContests();
  }, [currentPage]);

  const handleCreateClick = () => {
    setMode('create');
    setCreateForm(defaultCreateForm());
    setIsModalOpen(true);
  };

  const handleDeleteClick = (contest) => {
    setMode('delete');
    setSelectedContest(contest);
    setIsModalOpen(true);
  };

  const handleContestClick = (contest) => {
    navigate(`/admin/contests/${contest.id}`);
  };

  const handlePictureUpload = (contest) => {
    setSelectedContest(contest);
    fileInputRef.current?.click();
  };

  const handlePictureChange = async (event) => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file) return;

    try {
      const response = await updateContestPicture(selectedContest.id, file);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.toast.coverUpdated') });
        fetchContests();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.toast.coverUpdateFailed') });
    }
  };

  const handleSubmit = async () => {
    try {
      if (mode === 'create') {
        const formData = {
          ...createForm,
          start: new Date(createForm.start).toISOString(),
          duration: parseInt(createForm.duration) * 3600, // 转换为秒
        };
        const response = await createContest(formData);
        if (response.code === 200) {
          toast.success({ description: t('admin.contests.toast.createSuccess') });
          setIsModalOpen(false);
          fetchContests();
        }
      } else {
        const response = await deleteContest(selectedContest.id);
        if (response.code === 200) {
          toast.success({ description: t('admin.contests.toast.deleteSuccess') });
          setIsModalOpen(false);
          fetchContests();
        }
      }
    } catch (error) {
      toast.danger({
        description:
          error.message ||
          (mode === 'create' ? t('admin.contests.toast.createFailed') : t('admin.contests.toast.deleteFailed')),
      });
    }
  };

  const renderModalContent = () => {
    if (mode === 'create') {
      return (
        <div className="space-y-4">
          <Input
            label={t('admin.contests.form.name')}
            type="text"
            value={createForm.name}
            onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
            fullWidth
          />

          <Textarea
            label={t('admin.contests.form.description')}
            value={createForm.description}
            onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
            rows={4}
            fullWidth
          />

          <div className="grid grid-cols-2 gap-4">
            <Input
              label={t('admin.contests.form.flagPrefix')}
              type="text"
              value={createForm.prefix}
              onChange={(e) => setCreateForm({ ...createForm, prefix: e.target.value })}
              fullWidth
            />
            <Input
              label={t('admin.contests.form.teamSize')}
              type="number"
              value={createForm.size}
              onChange={(e) => setCreateForm({ ...createForm, size: parseInt(e.target.value) })}
              fullWidth
            />
          </div>

          <DateTimeInput
            label={t('admin.contests.form.startTime')}
            value={createForm.start}
            onChange={(e) => setCreateForm({ ...createForm, start: e.target.value })}
          />

          <Input
            label={t('admin.contests.form.durationHours')}
            type="number"
            value={createForm.duration}
            onChange={(e) => setCreateForm({ ...createForm, duration: parseInt(e.target.value) })}
            fullWidth
          />
        </div>
      );
    } else {
      return (
        <DeleteConfirmation
          message={`${t('admin.contests.modal.deletePrompt')} ${selectedContest?.name}?`}
          warning={t('common.modal.deleteWarning')}
        />
      );
    }
  };

  const renderModalFooter = () => {
    return <CRUDModalFooter mode={mode} onCancel={() => setIsModalOpen(false)} onSubmit={handleSubmit} />;
  };

  return (
    <>
      <AdminContests
        contests={contests}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        loading={loading}
        onPageChange={setCurrentPage}
        onCreateContest={handleCreateClick}
        onDeleteContest={handleDeleteClick}
        onContestClick={handleContestClick}
        onPictureUpload={handlePictureUpload}
      />

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={mode === 'create' ? t('admin.contests.modal.createTitle') : t('admin.contests.modal.deleteTitle')}
        size={mode === 'create' ? 'lg' : 'sm'}
        footer={renderModalFooter()}
      >
        {renderModalContent()}
      </Modal>

      <input
        type="file"
        ref={fileInputRef}
        className="hidden"
        accept="image/png,image/jpeg,image/jpg,image/gif"
        onChange={handlePictureChange}
      />
    </>
  );
}

export default ContestsManagement;
