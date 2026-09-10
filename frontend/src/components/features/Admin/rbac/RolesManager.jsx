import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getRoleList, createRole, updateRole, deleteRole } from '../../../../api/admin/rbac';
import { toast } from '../../../../utils/toast';
import { useCRUDModal } from '../../../../hooks';
import { Modal } from '../../../common';
import CRUDModalFooter from '../../../common/CRUDModalFooter';
import AdminRoles from './AdminRoles';
import RolePermissionsDialog from './RolePermissionsDialog';
import RbacForm from './RbacForm';
import { lastAvailablePage } from './payloads';

export default function RolesManager() {
  const { t } = useTranslation();
  const [roles, setRoles] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [permissionRole, setPermissionRole] = useState(null);
  const [revision, setRevision] = useState(0);
  const pageSize = 20;
  const refreshRoles = () => setRevision((value) => value + 1);

  useEffect(() => {
    let cancelled = false;
    getRoleList({ limit: pageSize, offset: (currentPage - 1) * pageSize })
      .then((response) => {
        if (!cancelled && response.code === 200) {
          const count = response.data.count || 0;
          setRoles(response.data.roles || []);
          setTotalCount(count);
          setCurrentPage(lastAvailablePage(count, pageSize, currentPage));
        }
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message || t('admin.rbac.roles.toast.fetchFailed') });
      });
    return () => {
      cancelled = true;
    };
  }, [currentPage, revision, t]);

  const modal = useCRUDModal({
    defaultForm: { name: '', description: '' },
    createApi: createRole,
    updateApi: updateRole,
    deleteApi: deleteRole,
    onSuccess: refreshRoles,
    itemToForm: (role) => ({ name: role.name, description: role.description || '' }),
    messages: Object.fromEntries(
      ['createSuccess', 'createFailed', 'updateSuccess', 'updateFailed', 'deleteSuccess', 'deleteFailed'].map((key) => [
        key,
        t(`admin.rbac.roles.toast.${key}`),
      ])
    ),
  });
  return (
    <>
      <AdminRoles
        roles={roles}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onCreateRole={modal.openCreate}
        onEditRole={modal.openEdit}
        onDeleteRole={modal.openDelete}
        onManagePermissions={setPermissionRole}
      />
      <Modal
        isOpen={modal.isModalOpen}
        onClose={modal.closeModal}
        title={t(`admin.rbac.roles.modal.${modal.mode}Title`)}
        size={modal.mode === 'delete' ? 'sm' : 'lg'}
        footer={<CRUDModalFooter mode={modal.mode} onCancel={modal.closeModal} onSubmit={modal.handleSubmit} />}
      >
        <RbacForm
          domain="roles"
          mode={modal.mode}
          item={modal.selectedItem}
          form={modal.editForm}
          setForm={modal.setEditForm}
        />
      </Modal>
      {permissionRole && (
        <RolePermissionsDialog key={permissionRole.id} role={permissionRole} onClose={() => setPermissionRole(null)} />
      )}
    </>
  );
}
