import { useState, useEffect, useRef } from 'react';
import { toast } from '../../utils/toast';
import {
  getOAuthProviderList,
  createOAuthProvider,
  updateOAuthProvider,
  deleteOAuthProvider,
  uploadOAuthPicture,
} from '../../api/admin/oauth';
import AdminOAuthProviders from '../../components/features/Admin/AdminOAuthProviders';
import { Modal } from '../../components/common';
import ModalButton from '../../components/common/ModalButton';
import Input from '../../components/common/Input';
import { useTranslation } from 'react-i18next';

function OAuthProvidersManagement() {
  // 状态管理
  const [providers, setProviders] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedProvider, setSelectedProvider] = useState(null);
  const [mode, setMode] = useState('edit'); // 'edit' | 'create' | 'delete'
  const fileInputRef = useRef(null);
  const uploadTargetRef = useRef(null);
  const [editForm, setEditForm] = useState({
    provider: '',
    uri: '',
    auth_url: '',
    token_url: '',
    user_info_url: '',
    callback_url: '',
    client_id: '',
    client_secret: '',
    picture: '',
    picture_field: '',
    name_field: '',
    email_field: '',
    description_field: '',
    id_field: '',
    on: false,
  });
  const { t } = useTranslation();

  const fetchProviders = async () => {
    try {
      const response = await getOAuthProviderList({
        limit: pageSize,
        offset: (currentPage - 1) * pageSize,
      });

      if (response.code === 200) {
        setProviders(response.data.providers);
        setTotalCount(response.data.count);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.oauthProviders.toast.fetchFailed') });
    }
  };

  // 数据获取
  useEffect(() => {
    fetchProviders();
  }, [currentPage]);

  const handleCreateClick = () => {
    setMode('create');
    setSelectedProvider(null);
    setEditForm({
      provider: '',
      uri: '',
      auth_url: '',
      token_url: '',
      user_info_url: '',
      callback_url: '',
      client_id: '',
      client_secret: '',
      picture: '',
      picture_field: '',
      name_field: '',
      email_field: '',
      description_field: '',
      id_field: '',
      on: false,
    });
    setIsModalOpen(true);
  };

  const handleEditClick = (provider) => {
    setMode('edit');
    setSelectedProvider(provider);
    setEditForm({
      provider: provider.provider,
      uri: provider.uri,
      auth_url: provider.auth_url,
      token_url: provider.token_url,
      user_info_url: provider.user_info_url,
      callback_url: provider.callback_url || '',
      client_id: provider.client_id,
      client_secret: provider.client_secret,
      picture: provider.picture,
      picture_field: provider.picture_field,
      name_field: provider.name_field,
      email_field: provider.email_field,
      description_field: provider.description_field,
      id_field: provider.id_field,
      on: provider.on,
    });
    setIsModalOpen(true);
  };

  const handleDeleteClick = (provider) => {
    setMode('delete');
    setSelectedProvider(provider);
    setIsModalOpen(true);
  };

  const handleProviderClick = (provider) => {
    handleEditClick(provider);
  };

  const handleCreateProvider = async () => {
    try {
      const response = await createOAuthProvider(editForm);
      if (response.code === 200) {
        toast.success({ description: t('admin.oauthProviders.toast.createSuccess') });
        setIsModalOpen(false);
        fetchProviders();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.oauthProviders.toast.createFailed') });
    }
  };

  const handleUpdateProvider = async () => {
    try {
      const response = await updateOAuthProvider(selectedProvider.id, editForm);
      if (response.code === 200) {
        toast.success({ description: t('admin.oauthProviders.toast.updateSuccess') });
        setIsModalOpen(false);
        fetchProviders();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.oauthProviders.toast.updateFailed') });
    }
  };

  const handleDeleteProvider = async () => {
    try {
      const response = await deleteOAuthProvider(selectedProvider.id);
      if (response.code === 200) {
        toast.success({ description: t('admin.oauthProviders.toast.deleteSuccess') });
        setIsModalOpen(false);
        fetchProviders();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.oauthProviders.toast.deleteFailed') });
    }
  };

  const handlePictureUpload = (provider) => {
    uploadTargetRef.current = provider;
    fileInputRef.current?.click();
  };

  const handleFileChange = async (event) => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file || !uploadTargetRef.current) return;

    try {
      const response = await uploadOAuthPicture(uploadTargetRef.current.id, file);
      if (response.code === 200) {
        toast.success({ description: t('admin.oauthProviders.toast.pictureUpdated') });
        fetchProviders();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.oauthProviders.toast.pictureUpdateFailed') });
    }
  };

  const renderModalContent = () => {
    if (mode === 'delete') {
      return (
        <div className="text-center">
          <p className="text-neutral-300 mb-4">
            {t('admin.oauthProviders.modal.deletePrompt')}{' '}
            <span className="font-semibold text-red-400">{selectedProvider?.provider}</span>?
          </p>
          <p className="text-neutral-400 text-sm">{t('admin.oauthProviders.modal.deleteWarning')}</p>
        </div>
      );
    }

    return (
      <div className="space-y-3">
        {/* 基本信息 */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.oauthProviders.form.providerLabel')}
            </label>
            <Input
              type="text"
              value={editForm.provider}
              onChange={(e) => setEditForm({ ...editForm, provider: e.target.value })}
              placeholder={t('admin.oauthProviders.form.providerPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.oauthProviders.form.uriLabel')}
            </label>
            <Input
              type="text"
              value={editForm.uri}
              onChange={(e) => setEditForm({ ...editForm, uri: e.target.value })}
              placeholder={t('admin.oauthProviders.form.uriPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
        </div>

        {/* URL配置 */}
        <div className="space-y-3">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.oauthProviders.form.authUrlLabel')}
            </label>
            <Input
              type="text"
              value={editForm.auth_url}
              onChange={(e) => setEditForm({ ...editForm, auth_url: e.target.value })}
              placeholder={t('admin.oauthProviders.form.authUrlPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.oauthProviders.form.tokenUrlLabel')}
            </label>
            <Input
              type="text"
              value={editForm.token_url}
              onChange={(e) => setEditForm({ ...editForm, token_url: e.target.value })}
              placeholder={t('admin.oauthProviders.form.tokenUrlPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.oauthProviders.form.userInfoUrlLabel')}
            </label>
            <Input
              type="text"
              value={editForm.user_info_url}
              onChange={(e) => setEditForm({ ...editForm, user_info_url: e.target.value })}
              placeholder={t('admin.oauthProviders.form.userInfoUrlPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.oauthProviders.form.callbackUrlLabel')}
            </label>
            <Input
              type="text"
              value={editForm.callback_url}
              onChange={(e) => setEditForm({ ...editForm, callback_url: e.target.value })}
              placeholder={t('admin.oauthProviders.form.callbackUrlPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
        </div>

        {/* 客户端配置 */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.oauthProviders.form.clientIdLabel')}
            </label>
            <Input
              type="text"
              value={editForm.client_id}
              onChange={(e) => setEditForm({ ...editForm, client_id: e.target.value })}
              placeholder={t('admin.oauthProviders.form.clientIdPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.oauthProviders.form.clientSecretLabel')}
            </label>
            <Input
              type="password"
              value={editForm.client_secret}
              onChange={(e) => setEditForm({ ...editForm, client_secret: e.target.value })}
              placeholder={t('admin.oauthProviders.form.clientSecretPlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
        </div>

        {/* 字段映射 */}
        <div className="space-y-3">
          <h3 className="text-neutral-200 font-medium">{t('admin.oauthProviders.form.mappingTitle')}</h3>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-neutral-300 text-sm font-medium mb-2">
                {t('admin.oauthProviders.form.pictureFieldLabel')}
              </label>
              <Input
                type="text"
                value={editForm.picture_field}
                onChange={(e) => setEditForm({ ...editForm, picture_field: e.target.value })}
                placeholder="{picture_url}"
                fullWidth
              />
            </div>
            <div>
              <label className="block text-neutral-300 text-sm font-medium mb-2">
                {t('admin.oauthProviders.form.nameFieldLabel')}
              </label>
              <Input
                type="text"
                value={editForm.name_field}
                onChange={(e) => setEditForm({ ...editForm, name_field: e.target.value })}
                placeholder="{name}"
                fullWidth
                required={mode === 'create'}
              />
            </div>
            <div>
              <label className="block text-neutral-300 text-sm font-medium mb-2">
                {t('admin.oauthProviders.form.emailFieldLabel')}
              </label>
              <Input
                type="text"
                value={editForm.email_field}
                onChange={(e) => setEditForm({ ...editForm, email_field: e.target.value })}
                placeholder="{email}"
                fullWidth
                required={mode === 'create'}
              />
            </div>
            <div>
              <label className="block text-neutral-300 text-sm font-medium mb-2">
                {t('admin.oauthProviders.form.descriptionFieldLabel')}
              </label>
              <Input
                type="text"
                value={editForm.description_field}
                onChange={(e) => setEditForm({ ...editForm, description_field: e.target.value })}
                placeholder="{description}"
                fullWidth
              />
            </div>
            <div>
              <label className="block text-neutral-300 text-sm font-medium mb-2">
                {t('admin.oauthProviders.form.idFieldLabel')}
              </label>
              <Input
                type="text"
                value={editForm.id_field}
                onChange={(e) => setEditForm({ ...editForm, id_field: e.target.value })}
                placeholder="{id}"
                fullWidth
                required={mode === 'create'}
              />
            </div>
            <div>
              <label className="block text-neutral-300 text-sm font-medium mb-2">
                {t('admin.oauthProviders.form.logoUrlLabel')}
              </label>
              <Input
                type="text"
                value={editForm.picture}
                onChange={(e) => setEditForm({ ...editForm, picture: e.target.value })}
                placeholder={t('admin.oauthProviders.form.logoUrlPlaceholder')}
                fullWidth
              />
            </div>
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
            {t('admin.oauthProviders.form.enable')}
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
          onClick={
            mode === 'create' ? handleCreateProvider : mode === 'edit' ? handleUpdateProvider : handleDeleteProvider
          }
        >
          {mode === 'create' ? t('common.create') : mode === 'edit' ? t('common.save') : t('common.delete')}
        </ModalButton>
      </>
    );
  };

  return (
    <>
      <AdminOAuthProviders
        providers={providers}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        loading={false}
        onPageChange={setCurrentPage}
        onCreateProvider={handleCreateClick}
        onEditProvider={handleEditClick}
        onDeleteProvider={handleDeleteClick}
        onProviderClick={handleProviderClick}
        onPictureUpload={handlePictureUpload}
      />

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={
          mode === 'create'
            ? t('admin.oauthProviders.modal.createTitle')
            : mode === 'edit'
              ? t('admin.oauthProviders.modal.editTitle')
              : t('admin.oauthProviders.modal.deleteTitle')
        }
        size={mode !== 'delete' ? 'lg' : 'sm'}
        footer={renderModalFooter()}
      >
        {renderModalContent()}
      </Modal>

      <input type="file" ref={fileInputRef} className="hidden" accept="image/*" onChange={handleFileChange} />
    </>
  );
}

export default OAuthProvidersManagement;
