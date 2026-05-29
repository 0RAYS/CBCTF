import { useState, useEffect, useRef } from 'react';
import { toast } from '../../../utils/toast';
import { getUserList, updateUser, deleteUser, createUser, updateUserPicture } from '../../../api/admin/user';
import AdminUsers from '../../../components/features/Admin/AdminUsers';
import { FormField, FormSwitch, Input, Modal, Textarea } from '../../../components/common';
import CRUDModalFooter from '../../../components/common/CRUDModalFooter';
import ModalFooter from '../../../components/common/ModalFooter';
import { useDebounce } from '../../../hooks';
import { useCRUDModal } from '../../../hooks/index.js';
import { useTranslation } from 'react-i18next';

function UsersTab() {
  const [users, setUsers] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;
  const fileInputRef = useRef(null);
  const { t } = useTranslation();

  const [showDetailDialog, setShowDetailDialog] = useState(false);
  const [detailUser, setDetailUser] = useState(null);

  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [passwordTarget, setPasswordTarget] = useState(null);
  const [passwordForm, setPasswordForm] = useState({ newPassword: '', confirmPassword: '' });
  const [passwordLoading, setPasswordLoading] = useState(false);

  function fetchUsers() {
    getUserList({
      limit: pageSize,
      offset: (currentPage - 1) * pageSize,
    })
      .then((response) => {
        if (response.code === 200) {
          setUsers(response.data.users);
          setTotalCount(response.data.count);
        }
      })
      .catch((error) => {
        toast.danger({ description: error.message || t('admin.users.toast.fetchFailed') });
      });
  }

  const defaultForm = {
    name: '',
    email: '',
    description: '',
    hidden: false,
    verified: true,
    banned: false,
    password: '',
  };

  const {
    isModalOpen,
    mode,
    selectedItem: selectedUser,
    editForm,
    setEditForm,
    openCreate,
    openEdit,
    openDelete,
    closeModal,
    handleSubmit,
  } = useCRUDModal({
    defaultForm,
    createApi: createUser,
    updateApi: updateUser,
    deleteApi: deleteUser,
    onSuccess: fetchUsers,
    itemToForm: (user) => ({
      name: user.name,
      email: user.email,
      description: user.description || '',
      hidden: user.hidden,
      verified: user.verified,
      banned: user.banned,
    }),
    messages: {
      createSuccess: t('admin.users.toast.createSuccess'),
      createFailed: t('admin.users.toast.createFailed'),
      updateSuccess: t('admin.users.toast.updateSuccess'),
      updateFailed: t('admin.users.toast.updateFailed'),
      deleteSuccess: t('admin.users.toast.deleteSuccess'),
      deleteFailed: t('admin.users.toast.deleteFailed'),
    },
  });

  const [nameQuery, setNameQuery] = useState('');
  const [emailQuery, setEmailQuery] = useState('');
  const [descQuery, setDescQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [searchLoading, setSearchLoading] = useState(false);

  const debouncedName = useDebounce(nameQuery, 300);
  const debouncedEmail = useDebounce(emailQuery, 300);
  const debouncedDesc = useDebounce(descQuery, 300);

  const isSearchMode = !!(nameQuery.trim() || emailQuery.trim() || descQuery.trim());

  useEffect(() => {
    let cancelled = false;
    if (!debouncedName.trim() && !debouncedEmail.trim() && !debouncedDesc.trim()) {
      setSearchResults([]);
      return;
    }
    const doSearch = async () => {
      setSearchLoading(true);
      try {
        const params = { limit: 20, offset: 0 };
        if (debouncedName.trim()) params.name = debouncedName.trim();
        if (debouncedEmail.trim()) params.email = debouncedEmail.trim();
        if (debouncedDesc.trim()) params.description = debouncedDesc.trim();
        const response = await getUserList(params);
        if (!cancelled && response.code === 200) {
          setSearchResults(response.data.users || []);
        }
      } catch (error) {
        if (!cancelled) {
          toast.danger({ description: error.message || t('admin.users.toast.searchFailed') });
          setSearchResults([]);
        }
      } finally {
        if (!cancelled) setSearchLoading(false);
      }
    };
    doSearch();
    return () => {
      cancelled = true;
    };
  }, [debouncedName, debouncedEmail, debouncedDesc]);

  useEffect(() => {
    if (!isSearchMode) {
      fetchUsers();
    }
  }, [currentPage, isSearchMode]);

  const handlePictureUpload = (user) => {
    openEdit(user);
    closeModal();
    fileInputRef.current?.click();
  };

  const handleFileChange = async (event) => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file) return;

    try {
      const response = await updateUserPicture(selectedUser.id, file);
      if (response.code === 200) {
        toast.success({ description: t('admin.users.toast.pictureUpdated') });
        fetchUsers();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.users.toast.pictureUpdateFailed') });
    }
  };

  const handleOpenPasswordModal = (user) => {
    setPasswordTarget(user);
    setPasswordForm({ newPassword: '', confirmPassword: '' });
    setShowPasswordModal(true);
  };

  const handleClosePasswordModal = () => {
    setShowPasswordModal(false);
    setPasswordTarget(null);
    setPasswordForm({ newPassword: '', confirmPassword: '' });
  };

  const handleChangePassword = async () => {
    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      toast.danger({ description: t('admin.users.form.passwordMismatch') });
      return;
    }
    setPasswordLoading(true);
    try {
      const response = await updateUser(passwordTarget.id, { password: passwordForm.newPassword });
      if (response.code === 200) {
        toast.success({ description: t('admin.users.toast.passwordChanged') });
        handleClosePasswordModal();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.users.toast.passwordChangeFailed') });
    } finally {
      setPasswordLoading(false);
    }
  };

  const renderModalContent = () => {
    if (mode === 'delete') {
      return (
        <div className="space-y-2 text-neutral-300">
          <p>
            {t('admin.users.modal.deletePrompt')} <span className="text-white font-semibold">{selectedUser?.name}</span>
            ?
          </p>
          <p className="text-red-400 text-sm">{t('admin.users.modal.deleteWarning')}</p>
        </div>
      );
    }

    return (
      <div className="space-y-4">
        <FormField label={t('admin.users.form.username')}>
          <Input
            type="text"
            value={editForm.name}
            onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
            placeholder={t('admin.users.form.usernamePlaceholder')}
            fullWidth
            required={mode === 'create'}
          />
        </FormField>

        <FormField label={t('admin.users.form.email')}>
          <Input
            type="email"
            value={editForm.email}
            onChange={(e) => setEditForm({ ...editForm, email: e.target.value })}
            placeholder={t('admin.users.form.emailPlaceholder')}
            fullWidth
            required={mode === 'create'}
          />
        </FormField>

        {mode === 'create' && (
          <FormField label={t('admin.users.form.password')}>
            <Input
              type="password"
              value={editForm.password}
              onChange={(e) => setEditForm({ ...editForm, password: e.target.value })}
              placeholder={t('admin.users.form.passwordPlaceholder')}
              fullWidth
              required
            />
          </FormField>
        )}

        <FormField label={t('admin.users.form.description')}>
          <Textarea
            value={editForm.description}
            onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
            placeholder={t('admin.users.form.descriptionPlaceholder')}
            rows={3}
            fullWidth
          />
        </FormField>

        <div className="flex flex-col gap-2">
          <FormSwitch
            id="verified"
            checked={editForm.verified}
            onChange={(e) => setEditForm({ ...editForm, verified: e.target.checked })}
            label={t('admin.users.status.verified')}
          />
          <FormSwitch
            id="banned"
            checked={editForm.banned}
            onChange={(e) => setEditForm({ ...editForm, banned: e.target.checked })}
            label={t('admin.users.status.banned')}
          />
          <FormSwitch
            id="hidden"
            checked={editForm.hidden}
            onChange={(e) => setEditForm({ ...editForm, hidden: e.target.checked })}
            label={t('admin.users.status.hidden')}
          />
        </div>
      </div>
    );
  };

  const renderModalFooter = () => {
    return <CRUDModalFooter mode={mode} onCancel={closeModal} onSubmit={handleSubmit} />;
  };

  const displayUsers = isSearchMode ? searchResults : users;
  const displayTotalCount = isSearchMode ? searchResults.length : totalCount;

  const handleRowClick = (user) => {
    setDetailUser(user);
    setShowDetailDialog(true);
  };

  const handleDetailClose = () => {
    setShowDetailDialog(false);
    setDetailUser(null);
  };

  return (
    <>
      <AdminUsers
        users={displayUsers}
        totalCount={displayTotalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        loading={false}
        onPageChange={setCurrentPage}
        onCreateUser={openCreate}
        onEditUser={openEdit}
        onDeleteUser={openDelete}
        onPictureUpload={handlePictureUpload}
        onChangePassword={handleOpenPasswordModal}
        nameQuery={nameQuery}
        emailQuery={emailQuery}
        descQuery={descQuery}
        searchLoading={searchLoading}
        isSearchMode={isSearchMode}
        onNameChange={setNameQuery}
        onEmailChange={setEmailQuery}
        onDescChange={setDescQuery}
        onRowClick={handleRowClick}
        showDetailDialog={showDetailDialog}
        detailUser={detailUser}
        onDetailClose={handleDetailClose}
      />

      <Modal
        isOpen={isModalOpen}
        onClose={closeModal}
        title={
          mode === 'create'
            ? t('admin.users.modal.createTitle')
            : mode === 'edit'
              ? t('admin.users.modal.editTitle')
              : t('admin.users.modal.deleteTitle')
        }
        size={mode !== 'delete' ? 'lg' : 'sm'}
        footer={renderModalFooter()}
      >
        {renderModalContent()}
      </Modal>

      <Modal
        isOpen={showPasswordModal}
        onClose={handleClosePasswordModal}
        title={t('admin.users.modal.changePasswordTitle')}
        size="sm"
        footer={
          <ModalFooter
            onCancel={handleClosePasswordModal}
            onSubmit={handleChangePassword}
            cancelLabel={t('common.cancel')}
            submitLabel={t('common.save')}
            submitDisabled={passwordLoading}
          />
        }
      >
        <div className="space-y-4">
          <FormField label={t('admin.users.form.newPassword')}>
            <Input
              type="password"
              value={passwordForm.newPassword}
              onChange={(e) => setPasswordForm({ ...passwordForm, newPassword: e.target.value })}
              placeholder={t('admin.users.form.newPasswordPlaceholder')}
              fullWidth
              required
            />
          </FormField>
          <FormField label={t('admin.users.form.confirmPassword')}>
            <Input
              type="password"
              value={passwordForm.confirmPassword}
              onChange={(e) => setPasswordForm({ ...passwordForm, confirmPassword: e.target.value })}
              placeholder={t('admin.users.form.confirmPasswordPlaceholder')}
              fullWidth
              required
            />
          </FormField>
        </div>
      </Modal>

      <input
        type="file"
        ref={fileInputRef}
        className="hidden"
        accept="image/png,image/jpeg,image/jpg,image/gif"
        onChange={handleFileChange}
      />
    </>
  );
}

export default UsersTab;
