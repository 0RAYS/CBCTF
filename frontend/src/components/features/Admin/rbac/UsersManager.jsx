import { useState, useEffect, useRef } from 'react';
import { toast } from '../../../../utils/toast';
import { getUserList, updateUser, deleteUser, createUser, updateUserPicture } from '../../../../api/admin/user';
import AdminUsers from './AdminUsers';
import UserForm from './UserForm';
import { lastAvailablePage, userUpdatePayload } from './payloads';
import { Modal } from '../../../common';
import CRUDModalFooter from '../../../common/CRUDModalFooter';
import { useDebounce, useCRUDModal } from '../../../../hooks';
import { useUserDetailDialog } from '../details/useUserDetailDialog.jsx';
import { useTranslation } from 'react-i18next';

function UsersManager() {
  const [users, setUsers] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [revision, setRevision] = useState(0);
  const pageSize = 20;
  const fileInputRef = useRef(null);
  const { t } = useTranslation();

  const pictureUserRef = useRef(null);
  const { openUserDetail, renderUserDetailDialog } = useUserDetailDialog();

  const refreshUsers = () => setRevision((value) => value + 1);

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
    onSuccess: refreshUsers,
    itemToForm: (user) => ({
      name: user.name,
      email: user.email,
      description: user.description || '',
      hidden: user.hidden,
      verified: user.verified,
      banned: user.banned,
      password: '',
    }),
    beforeUpdate: userUpdatePayload,
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
      setSearchLoading(false);
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
  }, [debouncedName, debouncedEmail, debouncedDesc, revision]);

  useEffect(() => {
    if (isSearchMode) return;
    let cancelled = false;
    getUserList({ limit: pageSize, offset: (currentPage - 1) * pageSize })
      .then((response) => {
        if (cancelled || response.code !== 200) return;
        const count = response.data.count || 0;
        setUsers(response.data.users || []);
        setTotalCount(count);
        setCurrentPage(lastAvailablePage(count, pageSize, currentPage));
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message || t('admin.users.toast.fetchFailed') });
      });
    return () => {
      cancelled = true;
    };
  }, [currentPage, isSearchMode, revision, t]);

  const handlePictureUpload = (user) => {
    pictureUserRef.current = user;
    fileInputRef.current?.click();
  };

  const handleFileChange = async (event) => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file || !pictureUserRef.current) return;

    try {
      const response = await updateUserPicture(pictureUserRef.current.id, file);
      if (response.code === 200) {
        toast.success({ description: t('admin.users.toast.pictureUpdated') });
        refreshUsers();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.users.toast.pictureUpdateFailed') });
    }
  };

  const renderModalFooter = () => {
    return <CRUDModalFooter mode={mode} onCancel={closeModal} onSubmit={handleSubmit} />;
  };

  const displayUsers = isSearchMode ? searchResults : users;
  const displayTotalCount = isSearchMode ? searchResults.length : totalCount;

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
        nameQuery={nameQuery}
        emailQuery={emailQuery}
        descQuery={descQuery}
        searchLoading={searchLoading}
        isSearchMode={isSearchMode}
        onNameChange={setNameQuery}
        onEmailChange={setEmailQuery}
        onDescChange={setDescQuery}
        onRowClick={(user) => openUserDetail(user.id)}
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
        <UserForm mode={mode} selectedUser={selectedUser} editForm={editForm} setEditForm={setEditForm} />
      </Modal>

      {renderUserDetailDialog()}

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

export default UsersManager;
