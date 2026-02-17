import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from '../../utils/toast';
import AdminContests from '../../components/features/Admin/AdminContests';
import { Modal } from '../../components/common';
import CRUDModalFooter from '../../components/common/CRUDModalFooter';
import DeleteConfirmation from '../../components/common/DeleteConfirmation';
import { getContestList, createContest, deleteContest, updateContestPicture } from '../../api/admin/contest';
import Input from '../../components/common/Input';
import Textarea from '../../components/common/Textarea';
import { DateTimeInput } from '../../components/common';
import { useTranslation } from 'react-i18next';

// 获取明天的日期时间
const getTomorrowDateTime = () => {
  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  return tomorrow.toISOString().slice(0, 16);
};

function ContestsManagement() {
  const [contests, setContests] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [mode, setMode] = useState('create'); // 'create' | 'delete'
  const [selectedContest, setSelectedContest] = useState(null);
  const [createForm, setCreateForm] = useState({
    name: '',
    description: '',
    captcha: '',
    prefix: 'CBCTF',
    size: 4,
    start: getTomorrowDateTime(),
    duration: 24,
  });

  const fileInputRef = useRef(null);
  const pageSize = 10;
  const navigate = useNavigate();
  const { t } = useTranslation();

  const fetchContests = async () => {
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
    }
  };

  useEffect(() => {
    fetchContests();
  }, [currentPage]);

  // 创建比赛按钮点击
  const handleCreateClick = () => {
    setMode('create');
    setCreateForm({
      name: '',
      description: '',
      captcha: '',
      prefix: 'CBCTF',
      size: 4,
      start: getTomorrowDateTime(),
      duration: 24,
    });
    setIsModalOpen(true);
  };

  // 删除比赛按钮点击
  const handleDeleteClick = (contest) => {
    setMode('delete');
    setSelectedContest(contest);
    setIsModalOpen(true);
  };

  // 比赛点击（进入详情）
  const handleContestClick = (contest) => {
    navigate(`/admin/contests/${contest.id}`);
  };

  // 上传头像
  const handlePictureUpload = (contest) => {
    setSelectedContest(contest);
    fileInputRef.current?.click();
  };

  // 处理文件上传变更
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

  // 提交表单（创建或删除）
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

  // 渲染模态框内容
  const renderModalContent = () => {
    if (mode === 'create') {
      return (
        <div className="space-y-4">
          <div>
            <label className="text-neutral-400 text-sm block mb-1">{t('admin.contests.form.name')}</label>
            <Input
              type="text"
              value={createForm.name}
              onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
              fullWidth
            />
          </div>

          <div>
            <label className="text-neutral-400 text-sm block mb-1">{t('admin.contests.form.description')}</label>
            <Textarea
              value={createForm.description}
              onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
              rows={4}
              fullWidth
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-neutral-400 text-sm block mb-1">{t('admin.contests.form.flagPrefix')}</label>
              <Input
                type="text"
                value={createForm.prefix}
                onChange={(e) => setCreateForm({ ...createForm, prefix: e.target.value })}
                fullWidth
              />
            </div>
            <div>
              <label className="text-neutral-400 text-sm block mb-1">{t('admin.contests.form.teamSize')}</label>
              <Input
                type="number"
                value={createForm.size}
                onChange={(e) => setCreateForm({ ...createForm, size: parseInt(e.target.value) })}
                fullWidth
              />
            </div>
          </div>

          <div>
            <label className="text-neutral-400 text-sm block mb-1">{t('admin.contests.form.startTime')}</label>
            <DateTimeInput
              value={createForm.start}
              onChange={(e) => setCreateForm({ ...createForm, start: e.target.value })}
            />
          </div>

          <div>
            <label className="text-neutral-400 text-sm block mb-1">{t('admin.contests.form.durationHours')}</label>
            <Input
              type="number"
              value={createForm.duration}
              onChange={(e) => setCreateForm({ ...createForm, duration: parseInt(e.target.value) })}
              fullWidth
            />
          </div>
        </div>
      );
    } else {
      return (
        <DeleteConfirmation
          message={`${t('admin.contests.modal.deletePrompt')} ${selectedContest?.name}?`}
          warning={t('admin.contests.modal.deleteWarning')}
        />
      );
    }
  };

  // 渲染模态框底部按钮
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

      <input type="file" ref={fileInputRef} className="hidden" accept="image/*" onChange={handlePictureChange} />
    </>
  );
}

export default ContestsManagement;
