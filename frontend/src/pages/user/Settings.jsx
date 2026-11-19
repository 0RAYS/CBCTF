import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { activateEmail, deleteAccount, updatePassword, updateUserInfo, uploadPicture } from '../../api/user';
import { fetchUserInfo, logoutUser } from '../../store/user';
import { toast } from '../../utils/toast';
import UserSettings from '../../components/features/UserSettings';
import Loading from '../../components/common/Loading';
import { useTranslation } from 'react-i18next';

function Settings() {
  const { user, loading } = useSelector((state) => state.user);
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { t } = useTranslation();

  const handleUpdate = async (data) => {
    try {
      const response = await updateUserInfo(data);
      if (response.code === 200) {
        await dispatch(fetchUserInfo());
      }
    } catch (error) {
      toast.danger({ title: t('toast.common.updateFailed'), description: error.message });
    }
  };

  const handlePasswordChange = async (data) => {
    try {
      const payload = { new: data.newPassword, old: data.currentPassword };
      if (user.has_no_pwd) {
        payload.old = 'never_login_pwd';
      }
      const response = await updatePassword(payload);
      if (response.code === 200) {
        toast.success({ title: t('toast.user.passwordUpdated') });
        await dispatch(fetchUserInfo());
      }
    } catch (error) {
      toast.danger({ title: t('toast.user.passwordUpdateFailed'), description: error.message });
    }
  };

  const handleEmailVerify = async () => {
    try {
      await activateEmail();
    } catch (error) {
      toast.danger({ title: t('toast.common.sendFailed'), description: error.message });
    }
  };

  const handleDeleteAccount = async (password) => {
    try {
      const response = await deleteAccount({ password });
      if (response.code === 200) {
        await dispatch(logoutUser());
        navigate('/login');
      }
    } catch (error) {
      toast.danger({ title: t('toast.common.deleteFailed'), description: error.message });
    }
  };

  const handlePictureChange = async (file) => {
    try {
      await uploadPicture(file);
      await dispatch(fetchUserInfo());
      toast.success({ title: t('toast.user.avatarUpdated') });
    } catch (error) {
      toast.danger({ title: t('toast.common.updateFailed'), description: error.message });
    }
  };

  if (loading || !user) {
    return <Loading />;
  }

  return (
    <UserSettings
      user={{
        name: user.name || '',
        email: user.email || '',
        emailVerified: user.verified || false,
        picture: user.picture || 'https://avatars.githubusercontent.com/u/default',
        description: user.description || '',
        hasNoPwd: user.has_no_pwd || false,
      }}
      onUpdate={handleUpdate}
      onPasswordChange={handlePasswordChange}
      onEmailVerify={handleEmailVerify}
      onDeleteAccount={handleDeleteAccount}
      onPictureChange={handlePictureChange}
      onLogout={() => {
        navigate('/');
        dispatch(logoutUser());
      }}
    />
  );
}

export default Settings;
