import { useDispatch, useSelector } from 'react-redux';
import { updateAdminPassword, updateAdminInfo, updateAdminPicture } from '../../api/admin/system';
import { fetchUserInfo } from '../../store/user.js';
import { toast } from '../../utils/toast.js';
import AdminSetting from '../../components/features/Admin/AdminSetting';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

function AdminSettings() {
  const { user } = useSelector((state) => state.user);
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { t } = useTranslation();

  // 处理更新基本信息
  const handleUpdateInfo = async (data) => {
    try {
      const response = await updateAdminInfo(data);

      if (response.code === 200) {
        toast.success({ description: t('admin.settings.toast.updateInfoSuccess') });
        await dispatch(fetchUserInfo(true));
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.settings.toast.updateInfoFailed') });
    }
  };

  // 处理修改密码
  const handlePasswordChange = async (data) => {
    try {
      const response = await updateAdminPassword({
        old: data.currentPassword,
        new: data.newPassword,
      });

      if (response.code === 200) {
        toast.success({ description: t('admin.settings.toast.passwordUpdated') });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.settings.toast.passwordUpdateFailed') });
    }
  };

  // 处理头像上传
  const handlePictureChange = async (file) => {
    try {
      const response = await updateAdminPicture(file);

      if (response.code === 200) {
        toast.success({ description: t('admin.settings.toast.avatarUpdated') });
        await dispatch(fetchUserInfo(true));
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.settings.toast.avatarUpdateFailed') });
    }
  };

  // 处理登出
  const handleLogout = () => {
    localStorage.removeItem('token');
    toast.success({ description: t('admin.settings.toast.logoutSuccess') });
    navigate('/login');
  };

  // 提取需要的用户信息
  const adminInfo = {
    name: user?.name || '',
    email: user?.email || '',
    picture: user?.picture || '',
  };

  return (
    <AdminSetting
      admin={adminInfo}
      onUpdate={handleUpdateInfo}
      onPasswordChange={handlePasswordChange}
      onPictureChange={handlePictureChange}
      onLogout={handleLogout}
    />
  );
}

export default AdminSettings;
