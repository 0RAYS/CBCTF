import { useDispatch, useSelector } from 'react-redux';
import { updateAdminPassword, updateAdminInfo, updateAdminPicture } from '../../api/admin/system';
import { fetchUserInfo, logoutUser } from '../../store/user.js';
import { toast } from '../../utils/toast.js';
import UserSettings from '../../components/features/UserSettings';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

function AdminSettings() {
  const { user } = useSelector((state) => state.user);
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { t } = useTranslation();

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

  const handleLogout = () => {
    navigate('/');
    dispatch(logoutUser());
  };

  return (
    <UserSettings
      user={{
        name: user?.name || '',
        email: user?.email || '',
        emailVerified: true,
        picture: user?.picture || '',
        description: user?.description || '',
        hasNoPwd: false,
      }}
      onUpdate={handleUpdateInfo}
      onPasswordChange={handlePasswordChange}
      onPictureChange={handlePictureChange}
      onLogout={handleLogout}
    />
  );
}

export default AdminSettings;
