import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  getPermissionList,
  getRolePermissions,
  assignPermissionToRole,
  revokePermissionFromRole,
} from '../../../../api/admin/rbac';
import { toast } from '../../../../utils/toast';
import { Button, Modal, Pagination, Spinner } from '../../../common';
import Checkbox from '../../../common/Checkbox';

export default function RolePermissionsDialog({ role, onClose }) {
  const { t } = useTranslation();
  const [permissions, setPermissions] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [page, setPage] = useState(1);
  const [assignedIds, setAssignedIds] = useState([]);
  const [ready, setReady] = useState(false);
  const [loading, setLoading] = useState(true);
  const [pendingIds, setPendingIds] = useState([]);
  const busyIds = useRef(new Set());
  const pageSize = 50;

  useEffect(() => {
    let cancelled = false;
    getRolePermissions(role.id)
      .then((response) => {
        if (response.code !== 200) throw new Error(t('admin.rbac.roles.toast.fetchPermFailed'));
        if (cancelled) return;
        setAssignedIds((response.data.permissions || []).map((permission) => permission.id));
        setReady(true);
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message || t('admin.rbac.roles.toast.fetchPermFailed') });
      });
    return () => {
      cancelled = true;
    };
  }, [role.id, t]);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    getPermissionList({ limit: pageSize, offset: (page - 1) * pageSize })
      .then((response) => {
        if (response.code !== 200) throw new Error(t('admin.rbac.roles.toast.fetchPermFailed'));
        if (cancelled) return;
        setPermissions(response.data.permissions || []);
        setTotalCount(response.data.count || 0);
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message || t('admin.rbac.roles.toast.fetchPermFailed') });
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [page, t]);

  async function togglePermission(permission) {
    if (!ready || busyIds.current.has(permission.id)) return;
    const assigned = assignedIds.includes(permission.id);
    const action = assigned ? 'revokePerm' : 'assignPerm';
    busyIds.current.add(permission.id);
    setPendingIds([...busyIds.current]);
    try {
      const api = assigned ? revokePermissionFromRole : assignPermissionToRole;
      const response = await api(role.id, { permission_id: permission.id });
      if (response.code !== 200) throw new Error(t(`admin.rbac.roles.toast.${action}Failed`));
      setAssignedIds((prev) =>
        assigned ? prev.filter((id) => id !== permission.id) : [...new Set([...prev, permission.id])]
      );
      toast.success({ description: t(`admin.rbac.roles.toast.${action}Success`) });
    } catch (error) {
      toast.danger({ description: error.message || t(`admin.rbac.roles.toast.${action}Failed`) });
    } finally {
      busyIds.current.delete(permission.id);
      setPendingIds([...busyIds.current]);
    }
  }

  const byResource = permissions.reduce((groups, permission) => {
    (groups[permission.resource] ||= []).push(permission);
    return groups;
  }, {});
  const close = () => {
    if (pendingIds.length === 0) onClose();
  };
  return (
    <Modal
      isOpen
      onClose={close}
      title={`${t('admin.rbac.roles.modal.permissionsTitle')} - ${role.name}`}
      size="lg"
      footer={
        <Button size="sm" variant="ghost" disabled={pendingIds.length > 0} onClick={close}>
          {t('common.confirm')}
        </Button>
      }
    >
      <div className="space-y-6">
        {loading && <Spinner />}
        {Object.entries(byResource).map(([resource, items]) => (
          <div key={resource}>
            <h4 className="text-sm font-medium text-neutral-300 mb-2 capitalize">
              {t(`admin.rbac.roles.permGroup.${resource}`, resource)}
            </h4>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
              {items.map((permission) => (
                <Checkbox
                  key={permission.id}
                  checked={assignedIds.includes(permission.id)}
                  onChange={() => togglePermission(permission)}
                  disabled={!ready || loading || pendingIds.includes(permission.id)}
                  label={
                    <span className="flex min-w-0 flex-col break-words">
                      <span className="text-sm text-neutral-200 font-mono">{permission.name}</span>
                      <span className="text-xs text-neutral-500">{permission.description}</span>
                    </span>
                  }
                />
              ))}
            </div>
          </div>
        ))}
        {totalCount > pageSize && (
          <Pagination
            total={Math.ceil(totalCount / pageSize)}
            current={page}
            onChange={setPage}
            showTotal
            totalItems={totalCount}
            simple
          />
        )}
      </div>
    </Modal>
  );
}
