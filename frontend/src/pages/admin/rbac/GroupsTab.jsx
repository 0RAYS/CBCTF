import { useEffect, useMemo, useState } from 'react';
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
import { getUserList } from '../../../api/admin/user';
import AdminGroups from '../../../components/features/Admin/AdminGroups';
import { FormField, Input, List, Modal, Pagination, Select, Textarea } from '../../../components/common';
import CRUDModalFooter from '../../../components/common/CRUDModalFooter';
import ModalButton from '../../../components/common/ModalButton';
import { useCRUDModal, useDebounce } from '../../../hooks/index.js';
import { useTranslation } from 'react-i18next';

function GroupsTab() {
  const [groups, setGroups] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;
  const { t } = useTranslation();

  const [roles, setRoles] = useState([]);
  const [rolesMap, setRolesMap] = useState({});

  const [userModalOpen, setUserModalOpen] = useState(false);
  const [selectedGroupForUsers, setSelectedGroupForUsers] = useState(null);
  const [groupUsers, setGroupUsers] = useState([]);
  const [groupUserIds, setGroupUserIds] = useState([]);
  const [loadingUsers, setLoadingUsers] = useState(false);
  const [loadingGroupUserIds, setLoadingGroupUserIds] = useState(false);
  const [userPage, setUserPage] = useState(1);
  const [userTotalCount, setUserTotalCount] = useState(0);
  const userPageSize = 10;

  const [candidateNameQuery, setCandidateNameQuery] = useState('');
  const [candidateEmailQuery, setCandidateEmailQuery] = useState('');
  const [candidateDescQuery, setCandidateDescQuery] = useState('');
  const [candidateUsers, setCandidateUsers] = useState([]);
  const [loadingCandidateUsers, setLoadingCandidateUsers] = useState(false);
  const [candidatePage, setCandidatePage] = useState(1);
  const [candidateTotalCount, setCandidateTotalCount] = useState(0);
  const candidatePageSize = 10;
  const [selectedCandidateIds, setSelectedCandidateIds] = useState([]);
  const debouncedCandidateNameQuery = useDebounce(candidateNameQuery, 300);
  const debouncedCandidateEmailQuery = useDebounce(candidateEmailQuery, 300);
  const debouncedCandidateDescQuery = useDebounce(candidateDescQuery, 300);

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
  }, [t]);

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

  const fetchGroupUsers = async (groupId, page = 1) => {
    setLoadingUsers(true);
    try {
      const response = await getGroupUsers(groupId, {
        limit: userPageSize,
        offset: (page - 1) * userPageSize,
      });
      if (response.code === 200) {
        setGroupUsers(response.data.users || []);
        setUserTotalCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.rbac.groups.toast.fetchUsersFailed') });
    } finally {
      setLoadingUsers(false);
    }
  };

  const fetchAllGroupUserIds = async (groupId) => {
    setLoadingGroupUserIds(true);
    try {
      const ids = new Set();
      const limit = 100;
      let offset = 0;
      let total = 0;

      do {
        const response = await getGroupUsers(groupId, {
          limit,
          offset,
        });

        if (response.code !== 200) {
          break;
        }

        const users = response.data.users || [];
        users.forEach((user) => ids.add(user.id));
        total = response.data.count || users.length;
        offset += limit;
      } while (offset < total);

      setGroupUserIds(Array.from(ids));
    } catch (error) {
      toast.danger({ description: error.message || t('admin.rbac.groups.toast.fetchUsersFailed') });
      setGroupUserIds([]);
    } finally {
      setLoadingGroupUserIds(false);
    }
  };

  const resetUserModalState = () => {
    setSelectedGroupForUsers(null);
    setGroupUsers([]);
    setGroupUserIds([]);
    setUserPage(1);
    setUserTotalCount(0);
    setCandidateNameQuery('');
    setCandidateEmailQuery('');
    setCandidateDescQuery('');
    setCandidateUsers([]);
    setCandidatePage(1);
    setCandidateTotalCount(0);
    setSelectedCandidateIds([]);
  };

  const closeUserModal = () => {
    setUserModalOpen(false);
    resetUserModalState();
  };

  const handleManageUsers = (group) => {
    resetUserModalState();
    setSelectedGroupForUsers(group);
    setUserModalOpen(true);
    fetchGroupUsers(group.id, 1);
    fetchAllGroupUserIds(group.id);
  };

  const handleRemoveUserFromList = async (user) => {
    try {
      const response = await removeUserFromGroup(selectedGroupForUsers.id, { user_id: user.id });
      if (response.code === 200) {
        toast.success({ description: t('admin.rbac.groups.toast.removeUserSuccess') });
        setGroupUserIds((prev) => prev.filter((id) => id !== user.id));
        setCandidateUsers((prev) => prev.filter((candidate) => candidate.id !== user.id));
        fetchGroups();

        const nextPage = groupUsers.length === 1 && userPage > 1 ? userPage - 1 : userPage;
        if (nextPage !== userPage) {
          setUserPage(nextPage);
        }
        fetchGroupUsers(selectedGroupForUsers.id, nextPage);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.rbac.groups.toast.removeUserFailed') });
    }
  };

  const handleUserPageChange = (page) => {
    setUserPage(page);
    if (selectedGroupForUsers) {
      fetchGroupUsers(selectedGroupForUsers.id, page);
    }
  };

  useEffect(() => {
    if (!userModalOpen) {
      return;
    }

    let cancelled = false;

    const fetchCandidateUsers = async () => {
      setLoadingCandidateUsers(true);
      try {
        const params = {
          limit: candidatePageSize,
          offset: (candidatePage - 1) * candidatePageSize,
        };
        if (debouncedCandidateNameQuery.trim()) {
          params.name = debouncedCandidateNameQuery.trim();
        }
        if (debouncedCandidateEmailQuery.trim()) {
          params.email = debouncedCandidateEmailQuery.trim();
        }
        if (debouncedCandidateDescQuery.trim()) {
          params.description = debouncedCandidateDescQuery.trim();
        }

        const response = await getUserList(params);
        if (response.code !== 200) {
          throw new Error(t('admin.rbac.groups.toast.fetchCandidatesFailed'));
        }

        if (!cancelled) {
          setCandidateUsers(response.data.users || []);
          setCandidateTotalCount(response.data.count || 0);
        }
      } catch (error) {
        if (!cancelled) {
          toast.danger({ description: error.message || t('admin.rbac.groups.toast.fetchCandidatesFailed') });
          setCandidateUsers([]);
          setCandidateTotalCount(0);
        }
      } finally {
        if (!cancelled) {
          setLoadingCandidateUsers(false);
        }
      }
    };

    fetchCandidateUsers();

    return () => {
      cancelled = true;
    };
  }, [
    userModalOpen,
    candidatePage,
    candidatePageSize,
    debouncedCandidateNameQuery,
    debouncedCandidateEmailQuery,
    debouncedCandidateDescQuery,
    t,
  ]);

  useEffect(() => {
    if (!userModalOpen) {
      return;
    }
    setCandidatePage(1);
  }, [userModalOpen, debouncedCandidateNameQuery, debouncedCandidateEmailQuery, debouncedCandidateDescQuery]);

  const availableCandidateUsers = useMemo(
    () => candidateUsers.filter((user) => !groupUserIds.includes(user.id)),
    [candidateUsers, groupUserIds]
  );

  const assignableSelectedIds = useMemo(
    () => selectedCandidateIds.filter((id) => !groupUserIds.includes(id)),
    [selectedCandidateIds, groupUserIds]
  );

  const allVisibleCandidatesSelected =
    availableCandidateUsers.length > 0 &&
    availableCandidateUsers.every((user) => selectedCandidateIds.includes(user.id));

  const handleToggleCandidate = (userId) => {
    setSelectedCandidateIds((prev) => (prev.includes(userId) ? prev.filter((id) => id !== userId) : [...prev, userId]));
  };

  const handleToggleAllCandidates = () => {
    setSelectedCandidateIds((prev) => {
      if (allVisibleCandidatesSelected) {
        return prev.filter((id) => !availableCandidateUsers.some((user) => user.id === id));
      }

      const next = new Set(prev);
      availableCandidateUsers.forEach((user) => next.add(user.id));
      return Array.from(next);
    });
  };

  const handleCandidatePageChange = (page) => {
    setCandidatePage(page);
  };

  const handleAssignSelectedUsers = async () => {
    if (!selectedGroupForUsers || assignableSelectedIds.length === 0) {
      return;
    }

    const results = await Promise.allSettled(
      assignableSelectedIds.map((userId) => assignUserToGroup(selectedGroupForUsers.id, { user_id: userId }))
    );

    const successfulIds = [];
    results.forEach((result, index) => {
      if (result.status === 'fulfilled' && result.value.code === 200) {
        successfulIds.push(assignableSelectedIds[index]);
      }
    });

    if (successfulIds.length > 0) {
      setGroupUserIds((prev) => Array.from(new Set([...prev, ...successfulIds])));
      setSelectedCandidateIds((prev) => prev.filter((id) => !successfulIds.includes(id)));
      fetchGroups();
      fetchGroupUsers(selectedGroupForUsers.id, userPage);
    }

    if (successfulIds.length === assignableSelectedIds.length) {
      toast.success({
        description:
          successfulIds.length === 1
            ? t('admin.rbac.groups.toast.assignUserSuccess')
            : t('admin.rbac.groups.toast.assignUsersSuccess', { count: successfulIds.length }),
      });
      return;
    }

    if (successfulIds.length > 0) {
      toast.warning({
        description: t('admin.rbac.groups.toast.assignUsersPartial', {
          success: successfulIds.length,
          total: assignableSelectedIds.length,
        }),
      });
      return;
    }

    toast.danger({ description: t('admin.rbac.groups.toast.assignUserFailed') });
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
        <FormField label={t('admin.rbac.groups.form.name')}>
          <Input
            type="text"
            value={editForm.name}
            onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
            placeholder={t('admin.rbac.groups.form.namePlaceholder')}
            fullWidth
            required={mode === 'create'}
            disabled={mode === 'edit' && selectedGroup?.default}
          />
        </FormField>

        <FormField label={t('admin.rbac.groups.form.description')}>
          <Textarea
            value={editForm.description}
            onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
            placeholder={t('admin.rbac.groups.form.descriptionPlaceholder')}
            rows={3}
            fullWidth
          />
        </FormField>

        <FormField label={t('admin.rbac.groups.form.role')}>
          <Select
            value={editForm.role_id}
            onChange={(e) => setEditForm({ ...editForm, role_id: parseInt(e.target.value) || '' })}
            options={roleOptions}
            placeholder={t('admin.rbac.groups.form.rolePlaceholder')}
            fullWidth
          />
        </FormField>
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

      <Modal
        isOpen={userModalOpen}
        onClose={closeUserModal}
        title={`${t('admin.rbac.groups.modal.usersTitle')} - ${selectedGroupForUsers?.name || ''}`}
        size="xl"
        footer={
          <>
            <ModalButton onClick={closeUserModal}>{t('common.confirm')}</ModalButton>
            <ModalButton
              variant="primary"
              onClick={handleAssignSelectedUsers}
              disabled={assignableSelectedIds.length === 0 || loadingCandidateUsers || loadingGroupUserIds}
            >
              {t('admin.rbac.groups.form.addSelected')}
            </ModalButton>
          </>
        }
      >
        <div className="space-y-6">
          <div className="space-y-3">
            <div className="flex items-center justify-between gap-3">
              <h3 className="text-sm font-mono text-neutral-200">{t('admin.rbac.groups.modal.currentUsersTitle')}</h3>
              <span className="text-xs text-neutral-500">
                {t('admin.rbac.groups.form.currentUsersCount', { count: userTotalCount })}
              </span>
            </div>
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
              emptyContent={t('admin.rbac.groups.empty.currentUsers')}
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
              footer={
                userTotalCount > userPageSize && (
                  <Pagination
                    total={Math.ceil(userTotalCount / userPageSize)}
                    current={userPage}
                    onChange={handleUserPageChange}
                    showTotal
                    totalItems={userTotalCount}
                  />
                )
              }
            />
          </div>

          <div className="border-t border-neutral-300/10 pt-6">
            <div className="space-y-4">
              <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
                <div className="grid flex-1 gap-3 md:grid-cols-3">
                  <FormField label={t('admin.rbac.groups.form.userNameSearch')}>
                    <Input
                      type="text"
                      value={candidateNameQuery}
                      onChange={(e) => setCandidateNameQuery(e.target.value)}
                      placeholder={t('admin.rbac.groups.form.userNameSearchPlaceholder')}
                      fullWidth
                    />
                  </FormField>
                  <FormField label={t('admin.rbac.groups.form.userEmailSearch')}>
                    <Input
                      type="text"
                      value={candidateEmailQuery}
                      onChange={(e) => setCandidateEmailQuery(e.target.value)}
                      placeholder={t('admin.rbac.groups.form.userEmailSearchPlaceholder')}
                      fullWidth
                    />
                  </FormField>
                  <FormField label={t('admin.rbac.groups.form.userDescSearch')}>
                    <Input
                      type="text"
                      value={candidateDescQuery}
                      onChange={(e) => setCandidateDescQuery(e.target.value)}
                      placeholder={t('admin.rbac.groups.form.userDescSearchPlaceholder')}
                      fullWidth
                    />
                  </FormField>
                </div>
                <div className="pb-1 text-sm text-neutral-400">
                  {t('admin.rbac.groups.form.selectedUsers', { count: assignableSelectedIds.length })}
                </div>
              </div>

              <div className="flex flex-col gap-2 text-sm text-neutral-400 md:flex-row md:items-center md:justify-between">
                <label className="inline-flex items-center gap-2">
                  <input
                    type="checkbox"
                    className="h-4 w-4 rounded border border-neutral-300/40 bg-black/20"
                    checked={allVisibleCandidatesSelected}
                    onChange={handleToggleAllCandidates}
                    disabled={availableCandidateUsers.length === 0}
                  />
                  <span>{t('common.selectAll')}</span>
                </label>
              </div>

              <List
                columns={[
                  { key: 'select', label: '', width: '10%' },
                  { key: 'id', label: 'ID', width: '12%' },
                  { key: 'name', label: t('admin.rbac.groups.columns.userName'), width: '24%' },
                  { key: 'email', label: t('admin.rbac.groups.columns.email'), width: '28%' },
                  { key: 'description', label: t('admin.rbac.groups.columns.description'), width: '26%' },
                ]}
                data={availableCandidateUsers}
                loading={loadingCandidateUsers || loadingGroupUserIds}
                empty={!loadingCandidateUsers && !loadingGroupUserIds && availableCandidateUsers.length === 0}
                emptyContent={
                  debouncedCandidateNameQuery.trim() ||
                  debouncedCandidateEmailQuery.trim() ||
                  debouncedCandidateDescQuery.trim()
                    ? t('admin.rbac.groups.empty.searchCandidates')
                    : t('admin.rbac.groups.empty.availableUsers')
                }
                animate={false}
                renderCell={(item, column) => {
                  if (column.key === 'select') {
                    return (
                      <input
                        type="checkbox"
                        className="h-4 w-4 rounded border border-neutral-300/40 bg-black/20"
                        checked={selectedCandidateIds.includes(item.id)}
                        onChange={() => handleToggleCandidate(item.id)}
                      />
                    );
                  }
                  return item[column.key] || '-';
                }}
                footer={
                  candidateTotalCount > candidatePageSize && (
                    <Pagination
                      total={Math.ceil(candidateTotalCount / candidatePageSize)}
                      current={candidatePage}
                      onChange={handleCandidatePageChange}
                      showTotal
                      totalItems={candidateTotalCount}
                    />
                  )
                }
              />
            </div>
          </div>
        </div>
      </Modal>
    </>
  );
}

export default GroupsTab;
