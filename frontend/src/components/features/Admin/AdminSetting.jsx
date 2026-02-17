import { motion } from 'motion/react';
import { useState } from 'react';
import { IconUser, IconLock, IconLogout } from '@tabler/icons-react';
import { Button } from '../../../components/common';
import Avatar from '../../common/Avatar';
import { useTranslation } from 'react-i18next';

/**
 * 管理员设置页面组件
 * @param {Object} props
 * @param {Object} props.admin - 管理员信息
 * @param {Function} props.onUpdate - 更新基本信息的回调函数
 * @param {Function} props.onPasswordChange - 修改密码的回调函数
 * @param {Function} props.onPictureChange - 更新头像的回调函数
 * @param {Function} props.onLogout - 登出的回调函数
 * @example
 * <AdminSetting
 *   admin={{
 *     name: "管理员",
 *     email: "admin@example.com",
 *     picture: "https://avatars.githubusercontent.com/u/1"
 *   }}
 *   onUpdate={(data) => console.log('更新信息:', data)}
 *   onPasswordChange={(data) => console.log('修改密码:', data)}
 *   onPictureChange={(file) => console.log('更新头像:', file)}
 *   onLogout={() => console.log('登出')}
 * />
 */
function AdminSetting({ admin, onUpdate, onPasswordChange, onPictureChange, onLogout }) {
  const { t } = useTranslation();
  const [activeSection, setActiveSection] = useState('profile');
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
    { id: 'profile', label: t('admin.settings.ui.sections.profile'), icon: <IconUser size={18} /> },
    { id: 'security', label: t('admin.settings.ui.sections.security'), icon: <IconLock size={18} /> },
    { id: 'divider', type: 'divider' },
    {
      id: 'logout',
      label: t('admin.settings.ui.sections.logout'),
      icon: <IconLogout size={18} />,
      color: 'text-red-400',
    },
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
          newPassword: t('admin.settings.validation.newMin'),
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
        confirmPassword: value !== passwordValues.newPassword ? t('admin.settings.validation.confirmMismatch') : '',
      }));
    }
  };

  const handlePasswordChange = (e) => {
    e.preventDefault();

    const errors = {
      currentPassword: !passwordValues.currentPassword ? t('admin.settings.validation.currentRequired') : '',
      newPassword: passwordValues.newPassword.length < 6 ? t('admin.settings.validation.newMin') : '',
      confirmPassword:
        passwordValues.newPassword !== passwordValues.confirmPassword
          ? t('admin.settings.validation.confirmMismatch')
          : '',
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

  // 统一输入框样式
  const inputBaseClass =
    'w-full bg-black/20 border border-neutral-300/30 rounded-md p-3 text-neutral-50 font-mono focus:border-geek-400 focus:outline-none transition-colors duration-200';

  return (
    <div className="w-full mx-auto">
      <motion.div
        className="border border-neutral-300/30 rounded-md bg-black/30 backdrop-blur-[2px] overflow-hidden"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
      >
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
                    variant={section.id === 'logout' ? 'danger' : activeSection === section.id ? 'primary' : 'ghost'}
                    size="sm"
                    align="icon-left"
                    icon={section.icon}
                    className={`w-full min-w-0 h-[42px] justify-start ${
                      section.id === 'logout'
                        ? 'hover:!bg-red-400/10'
                        : activeSection === section.id
                          ? '!bg-geek-400/10'
                          : 'hover:!bg-neutral-700/10'
                    }`}
                    onClick={() => {
                      if (section.id === 'logout') {
                        onLogout();
                      } else {
                        setActiveSection(section.id);
                      }
                    }}
                  >
                    {section.label}
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
                    <Avatar src={admin.picture} name={admin.name} size="xl" shape="circle" />
                    <div
                      className="absolute inset-0 flex items-center justify-center
                               bg-black/60 rounded-full opacity-0 group-hover:opacity-100 transition-opacity
                               cursor-pointer"
                    >
                      <span className="text-neutral-200 text-sm font-mono">{t('admin.settings.ui.avatar.change')}</span>
                    </div>
                  </div>
                  <div>
                    <h3 className="text-neutral-50 font-mono mb-1">{admin.name}</h3>
                    <p className="text-neutral-400 text-sm">{t('admin.settings.ui.accountLabel')}</p>
                  </div>
                </div>

                {/* 基本信息表单 */}
                <form
                  onSubmit={(e) => {
                    e.preventDefault();
                    onUpdate({
                      name: e.target.name.value,
                      email: e.target.email.value,
                    });
                  }}
                  className="space-y-4"
                >
                  <div>
                    <label className="block text-neutral-400 text-sm mb-2">
                      {t('admin.settings.ui.form.username')}
                    </label>
                    <input name="name" type="text" defaultValue={admin.name} className={inputBaseClass} />
                  </div>
                  <div>
                    <label className="block text-neutral-400 text-sm mb-2">{t('admin.settings.ui.form.email')}</label>
                    <input name="email" type="email" defaultValue={admin.email} className={inputBaseClass} />
                  </div>

                  {/* 保存按钮 */}
                  <div className="flex justify-end pt-4 border-t border-neutral-300/30">
                    <Button type="submit" variant="primary" size="sm">
                      {t('admin.settings.ui.form.save')}
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
                  <h3 className="text-neutral-50 font-mono mb-4">{t('admin.settings.ui.security.title')}</h3>
                  <div className="space-y-4">
                    <div>
                      <input
                        name="currentPassword"
                        type="password"
                        placeholder={t('admin.settings.ui.security.placeholders.current')}
                        value={passwordValues.currentPassword}
                        onChange={handlePasswordInputChange}
                        onBlur={handlePasswordBlur}
                        className={`${inputBaseClass} ${
                          passwordErrors.currentPassword
                            ? 'border-red-400 focus:border-red-400 shadow-[0_0_10px_rgba(248,113,113,0.1)]'
                            : ''
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

                    <div>
                      <input
                        name="newPassword"
                        type="password"
                        placeholder={t('admin.settings.ui.security.placeholders.new')}
                        value={passwordValues.newPassword}
                        onChange={handlePasswordInputChange}
                        onBlur={handlePasswordBlur}
                        className={`${inputBaseClass} ${
                          passwordErrors.newPassword
                            ? 'border-red-400 focus:border-red-400 shadow-[0_0_10px_rgba(248,113,113,0.1)]'
                            : ''
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
                        placeholder={t('admin.settings.ui.security.placeholders.confirm')}
                        value={passwordValues.confirmPassword}
                        onChange={handlePasswordInputChange}
                        onBlur={handlePasswordBlur}
                        className={`${inputBaseClass} ${
                          passwordErrors.confirmPassword
                            ? 'border-red-400 focus:border-red-400 shadow-[0_0_10px_rgba(248,113,113,0.1)]'
                            : ''
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
                        {t('admin.settings.ui.security.update')}
                      </Button>
                    </div>
                  </div>
                </form>

                {/* 账户安全提示 */}
                <div className="p-4 border border-geek-400/30 rounded-md bg-geek-400/5">
                  <h3 className="text-geek-400 font-mono mb-2">{t('admin.settings.ui.security.hintTitle')}</h3>
                  <p className="text-neutral-400 text-sm mb-4">{t('admin.settings.ui.security.hintBody')}</p>
                </div>
              </motion.div>
            )}
          </div>
        </div>
      </motion.div>
    </div>
  );
}

export default AdminSetting;
