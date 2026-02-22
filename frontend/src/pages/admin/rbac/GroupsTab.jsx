import { useState, useEffect } from 'react';
import { toast } from '../../../utils/toast';
import {
  getGroupList,
  createGroup,
  updateGroup,
  deleteGroup,
  getRoleList,
  assignUserToGroup,
  removeUserFromGroup,
  getGroupUsers,
} from '../../../api/admin/rbac';
import AdminGroups from '../../../components/features/Admin/AdminGroups';
import { Modal } from '../../../components/common';
import CRUDModalFooter from '../../../components/common/CRUDModalFooter';
import ModalButton from '../../../components/common/ModalButton';
import Input from '../../../components/common/Input';
import Textarea from '../../../components/common/Textarea';
import Select from '../../../components/common/Select';
import List from '../../../components/common/List';
import { useCRUDModal } from '../../../hooks/index.js';
import { useTranslation } from 'react-i18next';

function GroupsTab() {
  const [groups, setGroups] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;
  const { t } = useTranslation();

  // Roles for dropdown
  const [roles, setRoles] = useState([]);
  const [rolesMap, setRolesMap] = useState({});

  // User management modal state
  const [userModalOpen, setUserModalOpen] = useState(false);
  const [selectedGroupForUsers, setSelectedGroupForUsers] = useState(null);
  const [userId, setUserId] = useState('');
  const [groupUsers, setGroupUsers] = useState([]);
  const [loadingUsers, setLoadingUsers] = useState(false);

  function fetchGroups() {
    getGroupList({
      limit: pageSize,
      offset: (currentPage - 1) * pageSize,
    })
      .then((response) => {
        if (response.code === 200) {
          setGroups(response.data.groups);
          setTotalCount(response.data.count);
        }
      })
      .catch((error) => {
        toast.danger({ description: error.message || t('admin.rbac.groups.toast.fetchFailed') });
      });
  }

  // Fetch roles on mount
  useEffect(() => {
    getRoleList({ limit: 50, offset: 0 })
      .then((response) => {
        if (response.code === 200) {
          const rolesList = response.data.roles || [];
          setRoles(rolesList);
          const map = {};
          rolesList.forEach((role) => {
            map[role.id] = role.name;
          });
          setRolesMap(map);
        }
      })
      .catch((error) => {
        toast.danger({ description: error.message || t('admin.rbac.groups.toast.fetchRolesFailed') });
      });
  }, []);

  const defaultForm = {
    name: '',
    description: '',
    role_id: '',
  };

  const {
    isModalOpen,
    mode,
    selectedItem: selectedGroup,
    editForm,
    setEditForm,
    openCreate,
    openEdit,
    openDelete,
    closeModal,
    handleSubmit,
  } = useCRUDModal({
    defaultForm,
    createApi: createGroup,
    updateApi: updateGroup,
    deleteApi: deleteGroup,
    onSuccess: fetchGroups,
    itemToForm: (group) => ({
      name: group.name,
      description: group.description || '',
      role_id: group.role_id || '',
    }),
    messages: {
      createSuccess: t('admin.rbac.groups.toast.createSuccess'),
      createFailed: t('admin.rbac.groups.toast.createFailed'),
      updateSuccess: t('admin.rbac.groups.toast.updateSuccess'),
      updateFailed: t('admin.rbac.groups.toast.updateFailed'),
      deleteSuccess: t('admin.rbac.groups.toast.deleteSuccess'),
      deleteFailed: t('admin.rbac.groups.toast.deleteFailed'),
    },
  });

  useEffect(() => {
    fetchGroups();
  }, [currentPage]);

  const fetchGroupUsers = async (groupId) => {
    setLoadingUsers(true);
    try {
      const response = await getGroupUsers(groupId);
      if (response.code === 200) {
        setGroupUsers(response.data.users || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.rbac.groups.toast.fetchUsersFailed') });
    } finally {
      setLoadingUsers(false);
    }
  };

  const handleManageUsers = (group) => {
    setSelectedGroupForUsers(group);
    setUserId('');
    setGroupUsers([]);
    setUserModalOpen(true);
    fetchGroupUsers(group.id);
  };

  const handleAssignUser = async () => {
    const id = parseInt(userId, 10);
    if (!id || isNaN(id)) return;

    try {
      const response = await assignUserToGroup(selectedGroupForUsers.id, { user_id: id });
      if (response.code === 200) {
        toast.success({ description: t('admin.rbac.groups.toast.assignUserSuccess') });
        setUserId('');
        fetchGroups();
        fetchGroupUsers(selectedGroupForUsers.id);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.rbac.groups.toast.assignUserFailed') });
    }
  };

  const handleRemoveUserFromList = async (user) => {
    try {
      const response = await removeUserFromGroup(selectedGroupForUsers.id, { user_id: user.id });
      if (response.code === 200) {
        toast.success({ description: t('admin.rbac.groups.toast.removeUserSuccess') });
        fetchGroups();
        fetchGroupUsers(selectedGroupForUsers.id);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.rbac.groups.toast.removeUserFailed') });
    }
  };

  const renderModalContent = () => {
    if (mode === 'delete') {
      return (
        <p className="text-neutral-300">
          {t('admin.rbac.groups.modal.deletePrompt')}{' '}
          <span className="text-white font-semibold">{selectedGroup?.name}</span>?{' '}
          {t('admin.rbac.groups.modal.deleteWarning')}
        </p>
      );
    }

    const roleOptions = roles.map((role) => ({
      value: role.id,
      label: role.name,
    }));

    return (
      <div className="space-y-4">
        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-400 mb-1">{t('admin.rbac.groups.form.name')}</label>
          <Input
            type="text"
            value={editForm.name}
            onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
            placeholder={t('admin.rbac.groups.form.namePlaceholder')}
            fullWidth
            required={mode === 'create'}
            disabled={mode === 'edit' && selectedGroup?.default}
          />
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-400 mb-1">
            {t('admin.rbac.groups.form.description')}
          </label>
          <Textarea
            value={editForm.description}
            onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
            placeholder={t('admin.rbac.groups.form.descriptionPlaceholder')}
            rows={3}
            fullWidth
          />
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-400 mb-1">{t('admin.rbac.groups.form.role')}</label>
          <Select
            value={editForm.role_id}
            onChange={(e) => setEditForm({ ...editForm, role_id: parseInt(e.target.value) || '' })}
            options={roleOptions}
            placeholder={t('admin.rbac.groups.form.rolePlaceholder')}
            fullWidth
          />
        </div>
      </div>
    );
  };

  const renderModalFooter = () => {
    return <CRUDModalFooter mode={mode} onCancel={closeModal} onSubmit={handleSubmit} />;
  };

  return (
    <>
      <AdminGroups
        groups={groups}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        loading={false}
        onPageChange={setCurrentPage}
        onCreateGroup={openCreate}
        onEditGroup={openEdit}
        onDeleteGroup={openDelete}
        onManageUsers={handleManageUsers}
        rolesMap={rolesMap}
      />

      {/* CRUD Modal */}
      <Modal
        isOpen={isModalOpen}
        onClose={closeModal}
        title={
          mode === 'create'
            ? t('admin.rbac.groups.modal.createTitle')
            : mode === 'edit'
              ? t('admin.rbac.groups.modal.editTitle')
              : t('admin.rbac.groups.modal.deleteTitle')
        }
        size={mode !== 'delete' ? 'lg' : 'sm'}
        footer={renderModalFooter()}
      >
        {renderModalContent()}
      </Modal>

      {/* User Management Modal */}
      <Modal
        isOpen={userModalOpen}
        onClose={() => setUserModalOpen(false)}
        title={`${t('admin.rbac.groups.modal.usersTitle')} - ${selectedGroupForUsers?.name || ''}`}
        size="lg"
        footer={<ModalButton onClick={() => setUserModalOpen(false)}>{t('common.confirm')}</ModalButton>}
      >
        <div className="space-y-4">
          <List
            columns={[
              { key: 'id', label: 'ID', width: '15%' },
              { key: 'name', label: t('admin.rbac.groups.columns.userName'), width: '30%' },
              { key: 'email', label: t('admin.rbac.groups.columns.email'), width: '35%' },
              { key: 'actions', label: t('admin.rbac.groups.columns.actions'), width: '20%' },
            ]}
            data={groupUsers}
            loading={loadingUsers}
            empty={!loadingUsers && groupUsers.length === 0}
            animate={false}
            renderCell={(item, column) => {
              if (column.key === 'actions') {
                return (
                  <ModalButton variant="danger" onClick={() => handleRemoveUserFromList(item)}>
                    {t('admin.rbac.groups.form.remove')}
                  </ModalButton>
                );
              }
              return item[column.key] ?? '-';
            }}
          />
          <div className="border-t border-neutral-300/10 pt-4">
            <label className="block text-sm font-medium text-neutral-400 mb-1">
              {t('admin.rbac.groups.form.userId')}
            </label>
            <div className="flex gap-2">
              <Input
                type="number"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                placeholder={t('admin.rbac.groups.form.userIdPlaceholder')}
                fullWidth
              />
              <ModalButton variant="primary" onClick={handleAssignUser}>
                {t('admin.rbac.groups.form.addUser')}
              </ModalButton>
            </div>
          </div>
        </div>
      </Modal>
    </>
  );
}

export default GroupsTab;
