import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getGroupList, createGroup, updateGroup, deleteGroup, getRoleList } from '../../../../api/admin/rbac';
import { toast } from '../../../../utils/toast';
import { useCRUDModal } from '../../../../hooks';
import { Modal } from '../../../common';
import CRUDModalFooter from '../../../common/CRUDModalFooter';
import AdminGroups from './AdminGroups';
import GroupMembersDialog from './GroupMembersDialog';
import RbacForm from './RbacForm';
import { lastAvailablePage } from './payloads';

export default function GroupsManager() {
  const { t } = useTranslation();
  const [groups, setGroups] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [roles, setRoles] = useState([]);
  const [memberGroup, setMemberGroup] = useState(null);
  const [revision, setRevision] = useState(0);
  const pageSize = 20;
  const refreshGroups = () => setRevision((value) => value + 1);

  useEffect(() => {
    let cancelled = false;
    getGroupList({ limit: pageSize, offset: (currentPage - 1) * pageSize })
      .then((response) => {
        if (!cancelled && response.code === 200) {
          const count = response.data.count || 0;
          setGroups(response.data.groups || []);
          setTotalCount(count);
          setCurrentPage(lastAvailablePage(count, pageSize, currentPage));
        }
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message || t('admin.rbac.groups.toast.fetchFailed') });
      });
    return () => {
      cancelled = true;
    };
  }, [currentPage, revision, t]);
  useEffect(() => {
    let cancelled = false;
    getRoleList({ limit: 50, offset: 0 })
      .then((response) => {
        if (!cancelled && response.code === 200) setRoles(response.data.roles || []);
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message || t('admin.rbac.groups.toast.fetchRolesFailed') });
      });
    return () => {
      cancelled = true;
    };
  }, [t]);

  const modal = useCRUDModal({
    defaultForm: { name: '', description: '', role_id: '' },
    createApi: createGroup,
    updateApi: updateGroup,
    deleteApi: deleteGroup,
    onSuccess: refreshGroups,
    itemToForm: (group) => ({ name: group.name, description: group.description || '', role_id: group.role_id || '' }),
    messages: Object.fromEntries(
      ['createSuccess', 'createFailed', 'updateSuccess', 'updateFailed', 'deleteSuccess', 'deleteFailed'].map((key) => [
        key,
        t(`admin.rbac.groups.toast.${key}`),
      ])
    ),
  });

  return (
    <>
      <AdminGroups
        groups={groups}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onCreateGroup={modal.openCreate}
        onEditGroup={modal.openEdit}
        onDeleteGroup={modal.openDelete}
        onManageUsers={setMemberGroup}
        rolesMap={Object.fromEntries(roles.map((role) => [role.id, role.name]))}
      />
      <Modal
        isOpen={modal.isModalOpen}
        onClose={modal.closeModal}
        title={t(`admin.rbac.groups.modal.${modal.mode}Title`)}
        size={modal.mode === 'delete' ? 'sm' : 'lg'}
        footer={<CRUDModalFooter mode={modal.mode} onCancel={modal.closeModal} onSubmit={modal.handleSubmit} />}
      >
        <RbacForm
          domain="groups"
          mode={modal.mode}
          item={modal.selectedItem}
          form={modal.editForm}
          setForm={modal.setEditForm}
          roles={roles}
        />
      </Modal>
      {memberGroup && (
        <GroupMembersDialog
          key={memberGroup.id}
          group={memberGroup}
          onClose={() => setMemberGroup(null)}
          onChanged={refreshGroups}
        />
      )}
    </>
  );
}
