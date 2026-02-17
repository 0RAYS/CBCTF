import { useState, useEffect } from 'react';
import { toast } from '../../utils/toast';
import { getSmtpList, createSmtp, updateSmtp, deleteSmtp } from '../../api/admin/smtp';
import { getEmailHistory, getAllEmailHistory } from '../../api/admin/email';
import AdminSmtp from '../../components/features/Admin/AdminSmtp';
import AdminEmailHistory from '../../components/features/Admin/AdminEmailHistory';
import { Modal } from '../../components/common';
import ModalButton from '../../components/common/ModalButton';
import Input from '../../components/common/Input';
import { useTranslation } from 'react-i18next';

function SmtpManagement() {
  // 状态管理
  const [smtpConfigs, setSmtpConfigs] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedSmtp, setSelectedSmtp] = useState(null);
  const [mode, setMode] = useState('edit'); // 'edit' | 'create' | 'delete'
  const [editForm, setEditForm] = useState({
    address: '',
    host: '',
    port: 587,
    pwd: '',
    on: false,
  });

  // 邮件历史记录状态
  const [emailHistory, setEmailHistory] = useState([]);
  const [emailTotalCount, setEmailTotalCount] = useState(0);
  const [emailCurrentPage, setEmailCurrentPage] = useState(1);
  const [isEmailModalOpen, setIsEmailModalOpen] = useState(false);
  const [selectedEmail, setSelectedEmail] = useState(null);
  const [activeTab, setActiveTab] = useState('smtp'); // 'smtp' | 'history'
  const { t, i18n } = useTranslation();

  const fetchSmtpConfigs = async () => {
    try {
      const response = await getSmtpList({
        limit: 20,
        offset: (currentPage - 1) * 20,
      });

      if (response.code === 200) {
        setSmtpConfigs(response.data.smtps || response.data);
        setTotalCount(response.data.count || response.data.length);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.smtp.toast.fetchListFailed') });
    }
  };

  const fetchEmailHistory = async () => {
    try {
      let response;
      if (selectedSmtp) {
        response = await getEmailHistory(selectedSmtp.id, {
          limit: 20,
          offset: (emailCurrentPage - 1) * 20,
        });
      } else {
        response = await getAllEmailHistory({
          limit: 20,
          offset: (emailCurrentPage - 1) * 20,
        });
      }

      if (response.code === 200) {
        setEmailHistory(response.data.emails || response.data);
        setEmailTotalCount(response.data.count || response.data.length);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.smtp.toast.fetchHistoryFailed') });
    }
  };

  // 数据获取
  useEffect(() => {
    if (activeTab === 'smtp') {
      fetchSmtpConfigs();
    } else if (activeTab === 'history') {
      fetchEmailHistory();
    }
  }, [currentPage, emailCurrentPage, activeTab]);

  const handleCreateClick = () => {
    setMode('create');
    setSelectedSmtp(null);
    setEditForm({
      address: '',
      host: '',
      port: 587,
      pwd: '',
      on: false,
    });
    setIsModalOpen(true);
  };

  const handleEditClick = (smtp) => {
    setMode('edit');
    setSelectedSmtp(smtp);
    setEditForm({
      address: smtp.address,
      host: smtp.host,
      port: smtp.port,
      pwd: '', // 不显示密码
      on: smtp.on,
    });
    setIsModalOpen(true);
  };

  const handleDeleteClick = (smtp) => {
    setMode('delete');
    setSelectedSmtp(smtp);
    setIsModalOpen(true);
  };

  const handleSmtpClick = (smtp) => {
    handleEditClick(smtp);
  };

  const handleViewEmail = (email) => {
    setSelectedEmail(email);
    setIsEmailModalOpen(true);
  };

  const handleEmailClick = (email) => {
    handleViewEmail(email);
  };

  const handleViewHistory = (smtp) => {
    setSelectedSmtp(smtp);
    setActiveTab('history');
    setEmailCurrentPage(1);
  };

  const handleCreateSmtp = async () => {
    try {
      const response = await createSmtp(editForm);
      if (response.code === 200) {
        toast.success({ description: t('admin.smtp.toast.createSuccess') });
        setIsModalOpen(false);
        fetchSmtpConfigs();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.smtp.toast.createFailed') });
    }
  };

  const handleUpdateSmtp = async () => {
    try {
      // 只发送有值的字段
      const updateData = {};
      if (editForm.address !== selectedSmtp.address) updateData.address = editForm.address;
      if (editForm.host !== selectedSmtp.host) updateData.host = editForm.host;
      if (editForm.port !== selectedSmtp.port) updateData.port = editForm.port;
      if (editForm.pwd) updateData.pwd = editForm.pwd; // 只有输入密码时才更新
      if (editForm.on !== selectedSmtp.on) updateData.on = editForm.on;

      const response = await updateSmtp(selectedSmtp.id, updateData);
      if (response.code === 200) {
        toast.success({ description: t('admin.smtp.toast.updateSuccess') });
        setIsModalOpen(false);
        fetchSmtpConfigs();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.smtp.toast.updateFailed') });
    }
  };

  const handleDeleteSmtp = async () => {
    try {
      const response = await deleteSmtp(selectedSmtp.id);
      if (response.code === 200) {
        toast.success({ description: t('admin.smtp.toast.deleteSuccess') });
        setIsModalOpen(false);
        fetchSmtpConfigs();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.smtp.toast.deleteFailed') });
    }
  };

  const renderModalContent = () => {
    if (mode === 'delete') {
      return (
        <div className="text-center">
          <p className="text-neutral-300 mb-4">
            {t('admin.smtp.modal.deletePrompt')}{' '}
            <span className="font-semibold text-red-400">{selectedSmtp?.address}</span>?
          </p>
          <p className="text-neutral-400 text-sm">{t('admin.smtp.modal.deleteWarning')}</p>
        </div>
      );
    }

    return (
      <div className="space-y-3">
        {/* 基本信息 */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.smtp.form.addressLabel')}
            </label>
            <Input
              type="email"
              value={editForm.address}
              onChange={(e) => setEditForm({ ...editForm, address: e.target.value })}
              placeholder={t('admin.smtp.form.addressPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.form.hostLabel')}</label>
            <Input
              type="text"
              value={editForm.host}
              onChange={(e) => setEditForm({ ...editForm, host: e.target.value })}
              placeholder={t('admin.smtp.form.hostPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
        </div>

        {/* 端口和密码 */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.form.portLabel')}</label>
            <Input
              type="number"
              value={editForm.port}
              onChange={(e) => setEditForm({ ...editForm, port: parseInt(e.target.value) || 587 })}
              placeholder={t('admin.smtp.form.portPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {mode === 'create' ? t('admin.smtp.form.passwordLabelCreate') : t('admin.smtp.form.passwordLabelEdit')}
            </label>
            <Input
              type="password"
              value={editForm.pwd}
              onChange={(e) => setEditForm({ ...editForm, pwd: e.target.value })}
              placeholder={
                mode === 'create'
                  ? t('admin.smtp.form.passwordPlaceholderCreate')
                  : t('admin.smtp.form.passwordPlaceholderEdit')
              }
              fullWidth
              required={mode === 'create'}
            />
          </div>
        </div>

        {/* 状态 */}
        <div className="flex items-center">
          <input
            type="checkbox"
            id="on"
            checked={editForm.on}
            onChange={(e) => setEditForm({ ...editForm, on: e.target.checked })}
            className="mr-2"
          />
          <label htmlFor="on" className="text-neutral-300">
            {t('admin.smtp.form.enable')}
          </label>
        </div>
      </div>
    );
  };

  const renderModalFooter = () => {
    return (
      <>
        <ModalButton onClick={() => setIsModalOpen(false)}>{t('common.cancel')}</ModalButton>
        <ModalButton
          variant={mode === 'delete' ? 'danger' : 'primary'}
          onClick={mode === 'create' ? handleCreateSmtp : mode === 'edit' ? handleUpdateSmtp : handleDeleteSmtp}
        >
          {mode === 'create' ? t('common.create') : mode === 'edit' ? t('common.save') : t('common.delete')}
        </ModalButton>
      </>
    );
  };

  const renderEmailModalContent = () => {
    if (!selectedEmail) return null;

    return (
      <div className="space-y-4">
        {/* 基本信息 */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.id')}</label>
            <div className="text-neutral-50 font-mono">#{selectedEmail.id}</div>
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.status')}</label>
            <div>
              {selectedEmail.success ? (
                <span className="text-green-400">{t('admin.smtp.email.statusSuccess')}</span>
              ) : (
                <span className="text-red-400">{t('admin.smtp.email.statusFailed')}</span>
              )}
            </div>
          </div>
        </div>

        {/* 发件人和收件人 */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.from')}</label>
            <div className="text-neutral-50">{selectedEmail.from}</div>
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.to')}</label>
            <div className="text-neutral-50">{selectedEmail.to}</div>
          </div>
        </div>

        {/* 主题 */}
        <div>
          <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.subject')}</label>
          <div className="text-neutral-50">{selectedEmail.subject || t('admin.smtp.email.noSubject')}</div>
        </div>

        {/* 发送时间 */}
        <div>
          <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.time')}</label>
          <div className="text-neutral-50">{new Date(selectedEmail.time).toLocaleString(i18n.language || 'en-US')}</div>
        </div>

        {/* 邮件内容 */}
        <div>
          <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.smtp.email.content')}</label>
          <div className="bg-neutral-800 border border-neutral-700 rounded-lg p-3 max-h-60 overflow-y-auto">
            <pre className="text-neutral-300 text-sm whitespace-pre-wrap font-mono">
              {selectedEmail.content || t('admin.smtp.email.noContent')}
            </pre>
          </div>
        </div>
      </div>
    );
  };

  return (
    <>
      {/* 标签页切换 */}
      <div className="w-full mx-auto mb-6">
        <div className="flex border-b border-neutral-700">
          <button
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === 'smtp'
                ? 'text-blue-400 border-b-2 border-blue-400'
                : 'text-neutral-400 hover:text-neutral-300'
            }`}
            onClick={() => setActiveTab('smtp')}
          >
            {t('admin.smtp.tabs.config')}
          </button>
          <button
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === 'history'
                ? 'text-blue-400 border-b-2 border-blue-400'
                : 'text-neutral-400 hover:text-neutral-300'
            }`}
            onClick={() => {
              setSelectedSmtp(null);
              setActiveTab('history');
              setEmailCurrentPage(1);
            }}
          >
            {t('admin.smtp.tabs.history')}
          </button>
        </div>
      </div>

      {/* SMTP配置管理 */}
      {activeTab === 'smtp' && (
        <AdminSmtp
          smtpConfigs={smtpConfigs}
          totalCount={totalCount}
          currentPage={currentPage}
          pageSize={10}
          loading={false}
          onPageChange={setCurrentPage}
          onCreateSmtp={handleCreateClick}
          onEditSmtp={handleEditClick}
          onDeleteSmtp={handleDeleteClick}
          onSmtpClick={handleSmtpClick}
          onViewHistory={handleViewHistory}
        />
      )}

      {/* 邮件历史记录 */}
      {activeTab === 'history' && (
        <AdminEmailHistory
          emailHistory={emailHistory}
          totalCount={emailTotalCount}
          currentPage={emailCurrentPage}
          pageSize={20}
          loading={false}
          onPageChange={setEmailCurrentPage}
          onViewEmail={handleViewEmail}
          onEmailClick={handleEmailClick}
          smtpAddress={selectedSmtp?.address}
        />
      )}

      {/* SMTP配置模态框 */}
      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={
          mode === 'create'
            ? t('admin.smtp.modal.createTitle')
            : mode === 'edit'
              ? t('admin.smtp.modal.editTitle')
              : t('admin.smtp.modal.deleteTitle')
        }
        size={mode !== 'delete' ? 'lg' : 'sm'}
        footer={renderModalFooter()}
      >
        {renderModalContent()}
      </Modal>

      {/* 邮件详情模态框 */}
      <Modal
        isOpen={isEmailModalOpen}
        onClose={() => setIsEmailModalOpen(false)}
        title={t('admin.smtp.email.detailTitle')}
        size="lg"
        footer={<ModalButton onClick={() => setIsEmailModalOpen(false)}>{t('admin.smtp.actions.close')}</ModalButton>}
      >
        {renderEmailModalContent()}
      </Modal>
    </>
  );
}

export default SmtpManagement;
