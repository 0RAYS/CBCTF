import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getUserInfo } from '../../../../api/admin/user';
import { toast } from '../../../../utils/toast';
import AdminUserDetailDialog from './UserDetailDialog';

export function useUserDetailDialog() {
  const { t } = useTranslation();
  const [show, setShow] = useState(false);
  const [userData, setUserData] = useState(null);
  const request = useRef(0);
  useEffect(
    () => () => {
      request.current += 1;
    },
    []
  );

  const openUserDetail = async (userId) => {
    if (!userId) return;
    const version = ++request.current;
    try {
      const res = await getUserInfo(userId);
      if (version === request.current && res.code === 200) {
        setUserData(res.data);
        setShow(true);
      }
    } catch (err) {
      if (version === request.current) toast.danger({ description: err.message || t('admin.users.toast.fetchFailed') });
    }
  };

  const renderUserDetailDialog = () => (
    <AdminUserDetailDialog
      isOpen={show}
      onClose={() => {
        request.current += 1;
        setShow(false);
        setUserData(null);
      }}
      user={userData}
    />
  );

  return { openUserDetail, renderUserDetailDialog };
}
