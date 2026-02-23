/**
 * 用户设置页面组件
 * @param {Object} props
 * @param {Object} props.user - 用户信息
 * @param {string} props.user.name - 用户显示名称
 * @param {string} props.user.email - 用户邮箱
 * @param {boolean} props.user.emailVerified - 邮箱是否已验证
 * @param {string} props.user.picture - 用户头像URL
 * @param {string} props.user.description - 用户个人简介
 * @param {Function} props.onUpdate - 更新基本信息的回调函数
 * @param {Function} props.onPasswordChange - 修改密码的回调函数
 * @param {Function} props.onEmailVerify - 发送验证邮件的回调函数
 * @param {Function} props.onDeleteAccount - 删除账户的回调函数，接收 password
 * @param {Function} props.onPictureChange - 更新头像的回调函数
 * @example
 * <UserSettings
 *   user={{
 *     name: "CyberNinja",
 *     email: "ninja@example.com",
 *     emailVerified: false,
 *     picture: "https://avatars.githubusercontent.com/u/1",
 *     description: "Security researcher..."
 *   }}
 *   onUpdate={(data) => console.log('Updating profile:', data)}
 *   onPasswordChange={(data) => console.log('Changing password:', data)}
 *   onEmailVerify={(email) => console.log('Verifying email:', email)}
 *   onDeleteAccount={(password) => console.log('Deleting account', password)}
 *   onPictureChange={(file) => console.log('Updating picture:', file)}
 * />
 */

import { motion } from 'motion/react';
import { useState } from 'react';
import ConfirmModal from '../common/ConfirmModal';
import Avatar from '../common/Avatar';
import { useNavigate } from 'react-router-dom';
import { logoutUser } from '../../store/user.js';
import { useDispatch } from 'react-redux';
import { Button, Card, Input } from '../../components/common';
import { useTranslation } from 'react-i18next';

function UserSettings({ user, onUpdate, onPasswordChange, onEmailVerify, onDeleteAccount, onPictureChange }) {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const [activeSection, setActiveSection] = useState('profile');
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showEmailModal, setShowEmailModal] = useState(false);
  const [deletePassword, setDeletePassword] = useState('');
  const [deletePasswordError, setDeletePasswordError] = useState('');
  const [passwordValues, setPasswordValues] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });
  const [passwordErrors, setPasswordErrors] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });

  const sections = [
    { id: 'profile', label: t('user.settings.sections.profile'), icon: '👤' },
    { id: 'security', label: t('user.settings.sections.security'), icon: '🔒' },
    { id: 'divider', type: 'divider' },
    { id: 'logout', label: t('user.settings.sections.logout'), icon: '🖖', color: 'text-red-400' },
  ];

  const handlePasswordInputChange = (e) => {
    const { name, value } = e.target;
    setPasswordValues((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handlePasswordBlur = (e) => {
    const { name, value } = e.target;

    if (name === 'newPassword') {
      if (value && value.length < 6) {
        setPasswordErrors((prev) => ({
          ...prev,
          newPassword: t('auth.validation.passwordMin'),
        }));
      } else {
        setPasswordErrors((prev) => ({
          ...prev,
          newPassword: '',
        }));
      }
    }

    if (name === 'confirmPassword' && passwordValues.newPassword) {
      setPasswordErrors((prev) => ({
        ...prev,
        confirmPassword: value !== passwordValues.newPassword ? t('auth.validation.passwordMismatch') : '',
      }));
    }
  };

  const handlePasswordChange = (e) => {
    e.preventDefault();

    const errors = {
      currentPassword: user.hasNoPwd
        ? ''
        : !passwordValues.currentPassword
          ? t('user.settings.validation.currentPasswordRequired')
          : '',
      newPassword: passwordValues.newPassword.length < 6 ? t('auth.validation.passwordMin') : '',
      confirmPassword:
        passwordValues.newPassword !== passwordValues.confirmPassword ? t('auth.validation.passwordMismatch') : '',
    };

    setPasswordErrors(errors);

    if (!Object.values(errors).some((error) => error)) {
      onPasswordChange(passwordValues);
      setPasswordValues({
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
      });
      setPasswordErrors({
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
      });
    }
  };

  const handlePictureClick = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*';
    input.onchange = (e) => {
      const file = e.target.files[0];
      if (file) {
        onPictureChange(file);
      }
    };
    input.click();
  };

  const handleLogout = () => {
    navigate('/');
    dispatch(logoutUser());
  };

  const openDeleteModal = () => {
    setDeletePassword('');
    setDeletePasswordError('');
    setShowDeleteModal(true);
  };

  const closeDeleteModal = () => {
    setShowDeleteModal(false);
    setDeletePassword('');
    setDeletePasswordError('');
  };

  const handleDeletePasswordChange = (event) => {
    setDeletePassword(event.target.value);
    if (deletePasswordError) {
      setDeletePasswordError('');
    }
  };

  return (
    <div className="w-full max-w-[1200px] mx-auto">
      <Card variant="default" padding="none" animate className="overflow-hidden">
        {/* 标题栏 */}
        <div className="p-6 border-b border-neutral-300/30">
          <h2 className="text-2xl font-mono text-neutral-50">{t('user.settings.title')}</h2>
        </div>

        <div className="flex">
          {/* 左侧导航 */}
          <div className="w-[240px] border-r border-neutral-300/30 p-4">
            <div className="space-y-2">
              {sections.map((section) =>
                section.type === 'divider' ? (
                  <div key={section.id} className="my-4 border-t border-neutral-300/30"></div>
                ) : (
                  <Button
                    key={section.id}
                    variant="ghost"
                    align="icon-left-text-center"
                    icon={<span className="text-xl">{section.icon}</span>}
                    textColor={
                      section.id === 'logout'
                        ? 'text-red-400'
                        : activeSection === section.id
                          ? 'text-geek-400'
                          : 'text-neutral-300'
                    }
                    className={`w-full px-4 py-3 rounded-md min-w-0 h-auto
                      ${
                        section.id === 'logout'
                          ? 'hover:bg-red-400/10'
                          : activeSection === section.id
                            ? 'bg-geek-400/10'
                            : 'hover:bg-neutral-700/10 hover:text-neutral-200'
                      }`}
                    onClick={() => {
                      if (section.id === 'logout') {
                        handleLogout();
                      } else {
                        setActiveSection(section.id);
                      }
                    }}
                  >
                    <span className="font-mono text-sm">{section.label}</span>
                  </Button>
                )
              )}
            </div>
          </div>

          {/* 右侧内容 */}
          <div className="flex-1 p-6">
            {/* 个人资料 */}
            {activeSection === 'profile' && (
              <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-6">
                {/* 头像上传 */}
                <div className="flex items-center gap-6">
                  <div className="relative group" onClick={handlePictureClick}>
                    <Avatar src={user.picture} name={user.name} size="xl" shape="circle" />
                    <div
                      className="absolute inset-0 flex items-center justify-center
                                            bg-black/60 rounded-full opacity-0 group-hover:opacity-100 transition-opacity
                                            cursor-pointer"
                    >
                      <span className="text-neutral-200 text-sm font-mono">{t('common.change')}</span>
                    </div>
                  </div>
                  <div>
                    <h3 className="text-neutral-50 font-mono mb-1">{user.name}</h3>
                    <p className="text-neutral-400 text-sm">
                      {user.emailVerified ? t('user.settings.emailVerified') : t('user.settings.emailNotVerified')}
                    </p>
                  </div>
                </div>

                {/* 基本信息表单 */}
                <form
                  onSubmit={(e) => {
                    e.preventDefault();
                    onUpdate({
                      name: e.target.name.value,
                      email: e.target.email.value,
                      description: e.target.description.value,
                    });
                  }}
                  className="space-y-4"
                >
                  <div>
                    <label className="block text-neutral-400 text-sm mb-2">{t('user.settings.displayName')}</label>
                    <input
                      required
                      name="name"
                      type="text"
                      defaultValue={user.name}
                      className="w-full p-3 bg-neutral-900 border border-neutral-300/30 rounded-md
                                                text-neutral-50 font-mono
                                                focus:outline-none focus:border-geek-400"
                    />
                  </div>
                  <div>
                    <label className="block text-neutral-400 text-sm mb-2">{t('auth.placeholders.email')}</label>
                    <div className="flex gap-2">
                      <input
                        required
                        name="email"
                        type="email"
                        defaultValue={user.email}
                        className="flex-1 p-3 bg-neutral-900 border border-neutral-300/30 rounded-md
                                                    text-neutral-50 font-mono
                                                    focus:outline-none focus:border-geek-400"
                      />
                      {!user.emailVerified && (
                        <Button
                          variant="primary"
                          size="sm"
                          className="border-yellow-400 text-yellow-400 hover:bg-yellow-400/10"
                          onClick={() => setShowEmailModal(true)}
                        >
                          {t('user.settings.verify')}
                        </Button>
                      )}
                    </div>
                  </div>
                  <div>
                    <label className="block text-neutral-400 text-sm mb-2">{t('user.settings.bio')}</label>
                    <textarea
                      name="description"
                      defaultValue={user.description}
                      rows={4}
                      className="w-full p-3 bg-neutral-900 border border-neutral-300/30 rounded-md
                                                text-neutral-50
                                                focus:outline-none focus:border-geek-400"
                    />
                  </div>

                  {/* 保存按钮 */}
                  <div className="flex justify-end pt-4 border-t border-neutral-300/30">
                    <Button type="submit" variant="primary" size="sm">
                      {t('common.saveChanges')}
                    </Button>
                  </div>
                </form>
              </motion.div>
            )}

            {/* 安全设置 */}
            {activeSection === 'security' && (
              <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-6">
                {/* 密码修改 */}
                <form onSubmit={handlePasswordChange} className="p-4 border border-neutral-300/30 rounded-md">
                  <h3 className="text-neutral-50 font-mono mb-4">{t('user.settings.changePassword')}</h3>
                  <div className="space-y-4">
                    {!user.hasNoPwd && (
                      <div>
                        <input
                          name="currentPassword"
                          type="password"
                          placeholder={t('user.settings.placeholders.currentPassword')}
                          value={passwordValues.currentPassword}
                          onChange={handlePasswordInputChange}
                          onBlur={handlePasswordBlur}
                          className={`w-full p-3 bg-neutral-900 border rounded-md
                                                    text-neutral-50 font-mono focus:outline-none transition-all duration-200
                                                    ${
                                                      passwordErrors.currentPassword
                                                        ? 'border-red-400 focus:border-red-400 shadow-[0_0_10px_rgba(248,113,113,0.1)]'
                                                        : 'border-neutral-300/30 focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.1)]'
                                                    }`}
                          required
                        />
                        <motion.div
                          initial={{ opacity: 0, height: 0 }}
                          animate={{
                            opacity: passwordErrors.currentPassword ? 1 : 0,
                            height: passwordErrors.currentPassword ? 'auto' : 0,
                          }}
                          transition={{ duration: 0.2 }}
                        >
                          {passwordErrors.currentPassword && (
                            <p className="mt-1 text-red-400 text-sm">{passwordErrors.currentPassword}</p>
                          )}
                        </motion.div>
                      </div>
                    )}

                    <div>
                      <input
                        name="newPassword"
                        type="password"
                        placeholder={t('user.settings.placeholders.newPassword')}
                        value={passwordValues.newPassword}
                        onChange={handlePasswordInputChange}
                        onBlur={handlePasswordBlur}
                        className={`w-full p-3 bg-neutral-900 border rounded-md
                                                    text-neutral-50 font-mono focus:outline-none transition-all duration-200
                                                    ${
                                                      passwordErrors.newPassword
                                                        ? 'border-red-400 focus:border-red-400 shadow-[0_0_10px_rgba(248,113,113,0.1)]'
                                                        : 'border-neutral-300/30 focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.1)]'
                                                    }`}
                        required
                      />
                      <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{
                          opacity: passwordErrors.newPassword ? 1 : 0,
                          height: passwordErrors.newPassword ? 'auto' : 0,
                        }}
                        transition={{ duration: 0.2 }}
                      >
                        {passwordErrors.newPassword && (
                          <p className="mt-1 text-red-400 text-sm">{passwordErrors.newPassword}</p>
                        )}
                      </motion.div>
                    </div>

                    <div>
                      <input
                        name="confirmPassword"
                        type="password"
                        placeholder={t('user.settings.placeholders.confirmPassword')}
                        value={passwordValues.confirmPassword}
                        onChange={handlePasswordInputChange}
                        onBlur={handlePasswordBlur}
                        className={`w-full p-3 bg-neutral-900 border rounded-md
                                                    text-neutral-50 font-mono focus:outline-none transition-all duration-200
                                                    ${
                                                      passwordErrors.confirmPassword
                                                        ? 'border-red-400 focus:border-red-400 shadow-[0_0_10px_rgba(248,113,113,0.1)]'
                                                        : 'border-neutral-300/30 focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.1)]'
                                                    }`}
                        required
                      />
                      <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{
                          opacity: passwordErrors.confirmPassword ? 1 : 0,
                          height: passwordErrors.confirmPassword ? 'auto' : 0,
                        }}
                        transition={{ duration: 0.2 }}
                      >
                        {passwordErrors.confirmPassword && (
                          <p className="mt-1 text-red-400 text-sm">{passwordErrors.confirmPassword}</p>
                        )}
                      </motion.div>
                    </div>

                    <div className="flex justify-end">
                      <Button type="submit" variant="primary" size="sm">
                        {t('user.settings.updatePassword')}
                      </Button>
                    </div>
                  </div>
                </form>

                {/* 账户注销 */}
                <div className="p-4 border border-red-400/30 rounded-md bg-red-400/5">
                  <h3 className="text-red-400 font-mono mb-2">{t('user.settings.deleteAccount')}</h3>
                  <p className="text-neutral-400 text-sm mb-4">{t('user.settings.deleteAccountHint')}</p>
                  <Button variant="danger" size="sm" onClick={openDeleteModal}>
                    {t('user.settings.deleteAccountAction')}
                  </Button>
                </div>
              </motion.div>
            )}
          </div>
        </div>
      </Card>

      {/* 确认模态框 */}
      <ConfirmModal
        isOpen={showDeleteModal}
        onClose={closeDeleteModal}
        onConfirm={() => {
          if (!deletePassword) {
            setDeletePasswordError(t('auth.validation.passwordRequired'));
            return;
          }
          onDeleteAccount(deletePassword);
          closeDeleteModal();
        }}
        title={t('user.settings.deleteAccount')}
        message={
          <div className="space-y-3">
            <p>{t('user.settings.deleteAccountConfirm')}</p>
            <Input
              type="password"
              value={deletePassword}
              onChange={handleDeletePasswordChange}
              placeholder={t('auth.placeholders.password')}
              error={deletePasswordError}
              autoFocus
              autoComplete="current-password"
            />
          </div>
        }
        confirmText={t('common.delete')}
        type="danger"
      />

      <ConfirmModal
        isOpen={showEmailModal}
        onClose={() => setShowEmailModal(false)}
        onConfirm={() => {
          onEmailVerify(user.email);
          setShowEmailModal(false);
        }}
        title={t('user.settings.verifyEmail')}
        message={t('user.settings.verifyEmailHint')}
        confirmText={t('common.send')}
      />
    </div>
  );
}

export default UserSettings;
