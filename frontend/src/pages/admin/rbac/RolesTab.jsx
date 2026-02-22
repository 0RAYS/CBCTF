import { useState, useEffect } from 'react';
import { toast } from '../../../utils/toast';
import {
  getRoleList,
  getRolePermissions,
  createRole,
  updateRole,
  deleteRole,
  getPermissionList,
  assignPermissionToRole,
  revokePermissionFromRole,
} from '../../../api/admin/rbac';
import AdminRoles from '../../../components/features/Admin/AdminRoles';
import { Modal } from '../../../components/common';
import CRUDModalFooter from '../../../components/common/CRUDModalFooter';
import ModalButton from '../../../components/common/ModalButton';
import Input from '../../../components/common/Input';
import Textarea from '../../../components/common/Textarea';
import { useCRUDModal } from '../../../hooks/index.js';
import { useTranslation } from 'react-i18next';

function RolesTab() {
  const [roles, setRoles] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;
  const { t } = useTranslation();

  // Permission management modal state
  const [permModalOpen, setPermModalOpen] = useState(false);
  const [selectedRoleForPerms, setSelectedRoleForPerms] = useState(null);
  const [allPermissions, setAllPermissions] = useState([]);
  const [rolePermissions, setRolePermissions] = useState([]);

  function fetchRoles() {
    getRoleList({
      limit: pageSize,
      offset: (currentPage - 1) * pageSize,
    })
      .then((response) => {
        if (response.code === 200) {
          setRoles(response.data.roles);
          setTotalCount(response.data.count);
        }
      })
      .catch((error) => {
        toast.danger({ description: error.message || t('admin.rbac.roles.toast.fetchFailed') });
      });
  }

  const defaultForm = {
    name: '',
    description: '',
  };

  const {
    isModalOpen,
    mode,
    selectedItem: selectedRole,
    editForm,
    setEditForm,
    openCreate,
    openEdit,
    openDelete,
    closeModal,
    handleSubmit,
  } = useCRUDModal({
    defaultForm,
    createApi: createRole,
    updateApi: updateRole,
    deleteApi: deleteRole,
    onSuccess: fetchRoles,
    itemToForm: (role) => ({
      name: role.name,
      description: role.description || '',
    }),
    messages: {
      createSuccess: t('admin.rbac.roles.toast.createSuccess'),
      createFailed: t('admin.rbac.roles.toast.createFailed'),
      updateSuccess: t('admin.rbac.roles.toast.updateSuccess'),
      updateFailed: t('admin.rbac.roles.toast.updateFailed'),
      deleteSuccess: t('admin.rbac.roles.toast.deleteSuccess'),
      deleteFailed: t('admin.rbac.roles.toast.deleteFailed'),
    },
  });

  useEffect(() => {
    fetchRoles();
  }, [currentPage]);

  // Fetch all permissions once on mount
  useEffect(() => {
    getPermissionList({ limit: 50, offset: 0 })
      .then((response) => {
        if (response.code === 200) {
          setAllPermissions(response.data.permissions || []);
        }
      })
      .catch((error) => {
        toast.danger({ description: error.message || t('admin.rbac.roles.toast.fetchPermFailed') });
      });
  }, []);

  const handleManagePermissions = async (role) => {
    setSelectedRoleForPerms(role);
    try {
      const response = await getRolePermissions(role.id);
      if (response.code === 200) {
        setRolePermissions(response.data.permissions || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.rbac.roles.toast.fetchPermFailed') });
    }
    setPermModalOpen(true);
  };

  const handleTogglePermission = async (permission, isAssigned) => {
    try {
      if (isAssigned) {
        const response = await revokePermissionFromRole(selectedRoleForPerms.id, {
          permission_id: permission.id,
        });
        if (response.code === 200) {
          setRolePermissions((prev) => prev.filter((p) => p.id !== permission.id));
          toast.success({ description: t('admin.rbac.roles.toast.revokePermSuccess') });
        }
      } else {
        const response = await assignPermissionToRole(selectedRoleForPerms.id, {
          permission_id: permission.id,
        });
        if (response.code === 200) {
          setRolePermissions((prev) => [...prev, permission]);
          toast.success({ description: t('admin.rbac.roles.toast.assignPermSuccess') });
        }
      }
    } catch (error) {
      toast.danger({
        description:
          error.message ||
          (isAssigned ? t('admin.rbac.roles.toast.revokePermFailed') : t('admin.rbac.roles.toast.assignPermFailed')),
      });
    }
  };

  // Group permissions by resource
  const permissionsByResource = allPermissions.reduce((acc, perm) => {
    if (!acc[perm.resource]) acc[perm.resource] = [];
    acc[perm.resource].push(perm);
    return acc;
  }, {});

  const renderModalContent = () => {
    if (mode === 'delete') {
      return (
        <p className="text-neutral-300">
          {t('admin.rbac.roles.modal.deletePrompt')}{' '}
          <span className="text-white font-semibold">{selectedRole?.name}</span>?{' '}
          {t('admin.rbac.roles.modal.deleteWarning')}
        </p>
      );
    }

    return (
      <div className="space-y-4">
        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-400 mb-1">{t('admin.rbac.roles.form.name')}</label>
          <Input
            type="text"
            value={editForm.name}
            onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
            placeholder={t('admin.rbac.roles.form.namePlaceholder')}
            fullWidth
            required={mode === 'create'}
            disabled={mode === 'edit' && selectedRole?.default}
          />
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-400 mb-1">
            {t('admin.rbac.roles.form.description')}
          </label>
          <Textarea
            value={editForm.description}
            onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
            placeholder={t('admin.rbac.roles.form.descriptionPlaceholder')}
            rows={3}
            fullWidth
          />
        </div>
      </div>
    );
  };

  const renderModalFooter = () => {
    return <CRUDModalFooter mode={mode} onCancel={closeModal} onSubmit={handleSubmit} />;
  };

  const rolePermissionIds = new Set(rolePermissions.map((p) => p.id));

  return (
    <>
      <AdminRoles
        roles={roles}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        loading={false}
        onPageChange={setCurrentPage}
        onCreateRole={openCreate}
        onEditRole={openEdit}
        onDeleteRole={openDelete}
        onManagePermissions={handleManagePermissions}
      />

      {/* CRUD Modal */}
      <Modal
        isOpen={isModalOpen}
        onClose={closeModal}
        title={
          mode === 'create'
            ? t('admin.rbac.roles.modal.createTitle')
            : mode === 'edit'
              ? t('admin.rbac.roles.modal.editTitle')
              : t('admin.rbac.roles.modal.deleteTitle')
        }
        size={mode !== 'delete' ? 'lg' : 'sm'}
        footer={renderModalFooter()}
      >
        {renderModalContent()}
      </Modal>

      {/* Permissions Management Modal */}
      <Modal
        isOpen={permModalOpen}
        onClose={() => setPermModalOpen(false)}
        title={`${t('admin.rbac.roles.modal.permissionsTitle')} - ${selectedRoleForPerms?.name || ''}`}
        size="lg"
        footer={<ModalButton onClick={() => setPermModalOpen(false)}>{t('common.confirm')}</ModalButton>}
      >
        <div className="space-y-6">
          {Object.entries(permissionsByResource).map(([resource, perms]) => (
            <div key={resource}>
              <h4 className="text-sm font-medium text-neutral-300 mb-2 capitalize">
                {t(`admin.rbac.roles.permGroup.${resource}`, resource)}
              </h4>
              <div className="grid grid-cols-2 gap-2">
                {perms.map((perm) => {
                  const isAssigned = rolePermissionIds.has(perm.id);
                  return (
                    <label
                      key={perm.id}
                      className="flex items-center gap-2 p-2 rounded hover:bg-neutral-800 cursor-pointer"
                    >
                      <input
                        type="checkbox"
                        checked={isAssigned}
                        onChange={() => handleTogglePermission(perm, isAssigned)}
                        className="shrink-0"
                      />
                      <div className="flex flex-col">
                        <span className="text-sm text-neutral-200 font-mono">{perm.name}</span>
                        <span className="text-xs text-neutral-500">{perm.description}</span>
                      </div>
                    </label>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      </Modal>
    </>
  );
}

export default RolesTab;
