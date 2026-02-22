import { useState, useEffect } from 'react';
import { toast } from '../../../utils/toast';
import { getPermissionList, updatePermission } from '../../../api/admin/rbac';
import AdminPermissions from '../../../components/features/Admin/AdminPermissions';
import { Modal } from '../../../components/common';
import ModalButton from '../../../components/common/ModalButton';
import Textarea from '../../../components/common/Textarea';
import { useTranslation } from 'react-i18next';

function PermissionsTab() {
  const [permissions, setPermissions] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 50;
  const { t } = useTranslation();

  // Edit modal state
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedPermission, setSelectedPermission] = useState(null);
  const [editDescription, setEditDescription] = useState('');

  function fetchPermissions() {
    getPermissionList({
      limit: pageSize,
      offset: (currentPage - 1) * pageSize,
    })
      .then((response) => {
        if (response.code === 200) {
          setPermissions(response.data.permissions || []);
          setTotalCount(response.data.count);
        }
      })
      .catch((error) => {
        toast.danger({ description: error.message || t('admin.rbac.permissions.toast.fetchFailed') });
      });
  }

  useEffect(() => {
    fetchPermissions();
  }, [currentPage]);

  const handleEditPermission = (permission) => {
    setSelectedPermission(permission);
    setEditDescription(permission.description || '');
    setIsModalOpen(true);
  };

  const handleUpdate = async () => {
    try {
      const response = await updatePermission(selectedPermission.id, {
        description: editDescription,
      });
      if (response.code === 200) {
        toast.success({ description: t('admin.rbac.permissions.toast.updateSuccess') });
        setIsModalOpen(false);
        fetchPermissions();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.rbac.permissions.toast.updateFailed') });
    }
  };

  return (
    <>
      <AdminPermissions
        permissions={permissions}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        loading={false}
        onPageChange={setCurrentPage}
        onEditPermission={handleEditPermission}
      />

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={t('admin.rbac.permissions.modal.editTitle')}
        size="md"
        footer={
          <>
            <ModalButton onClick={() => setIsModalOpen(false)}>{t('common.cancel')}</ModalButton>
            <ModalButton variant="primary" onClick={handleUpdate}>
              {t('common.save')}
            </ModalButton>
          </>
        }
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-neutral-400 mb-1">
              {t('admin.rbac.permissions.form.description')}
            </label>
            <Textarea
              value={editDescription}
              onChange={(e) => setEditDescription(e.target.value)}
              placeholder={t('admin.rbac.permissions.form.descriptionPlaceholder')}
              rows={3}
              fullWidth
            />
          </div>
        </div>
      </Modal>
    </>
  );
}

export default PermissionsTab;
