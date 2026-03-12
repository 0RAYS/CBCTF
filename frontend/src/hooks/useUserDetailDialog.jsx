import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getUserInfo } from '../api/admin/user';
import { toast } from '../utils/toast';
import AdminUserDetailDialog from '../components/features/Admin/AdminUserDetailDialog';

export function useUserDetailDialog() {
  const { t } = useTranslation();
  const [show, setShow] = useState(false);
  const [userData, setUserData] = useState(null);

  const openUserDetail = async (userId) => {
    if (!userId) return;
    try {
      const res = await getUserInfo(userId);
      if (res.code === 200) {
        setUserData(res.data);
        setShow(true);
      }
    } catch (err) {
      toast.danger({ description: err.message || t('admin.users.toast.fetchFailed') });
    }
  };

  const renderUserDetailDialog = () => (
    <AdminUserDetailDialog
      isOpen={show}
      onClose={() => {
        setShow(false);
        setUserData(null);
      }}
      user={userData}
    />
  );

  return { openUserDetail, renderUserDetailDialog };
}
