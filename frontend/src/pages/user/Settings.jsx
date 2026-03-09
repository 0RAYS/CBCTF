import { useSelector, useDispatch } from 'react-redux';
import { updatePassword, updateUserInfo, deleteAccount, uploadPicture, activateEmail } from '../../api/user';
import { fetchUserInfo } from '../../store/user';
import { toast } from '../../utils/toast';
import UserSettings from '../../components/features/UserSettings';
import Loading from '../../components/common/Loading';
import { useTranslation } from 'react-i18next';

function Settings() {
  const { user, loading } = useSelector((state) => state.user);
  const dispatch = useDispatch();
  const { t } = useTranslation();

  // 处理更新基本信息
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

  // 处理密码修改
  const handlePasswordChange = async (data) => {
    try {
      const payload = { new: data.newPassword };
      if (user.has_no_pwd) {
        payload.old = 'never_login_pwd';
      }
      await updatePassword(payload);
    } catch (error) {
      toast.danger({ title: t('toast.user.passwordUpdateFailed'), description: error.message });
    }
  };

  // 处理邮箱验证
  const handleEmailVerify = async () => {
    try {
      await activateEmail();
    } catch (error) {
      toast.danger({ title: t('toast.common.sendFailed'), description: error.message });
    }
  };

  // 处理账户删除
  const handleDeleteAccount = async (password) => {
    try {
      await deleteAccount({ password });
      // 这里可以添加登出和跳转逻辑
    } catch (error) {
      toast.danger({ title: t('toast.common.deleteFailed'), description: error.message });
    }
  };

  // 处理头像更新
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
    />
  );
}

export default Settings;
