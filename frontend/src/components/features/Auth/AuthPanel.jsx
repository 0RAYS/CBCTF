import { motion, AnimatePresence } from 'motion/react';
import { useState } from 'react';
import { Button, Input } from '../../common';
import OAuthLogin from './OAuthLogin';
import { useTranslation } from 'react-i18next';

function AuthPanel({ onSubmit }) {
  const { t } = useTranslation();
  const [isLogin, setIsLogin] = useState(true); // true为登录模式，false为注册模式
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    confirmPassword: '',
    email: '',
  });
  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  const validateForm = () => {
    const newErrors = {};

    // 用户名验证
    if (formData.username.length < 3) {
      newErrors.username = t('auth.validation.usernameMin');
    }

    // 密码验证
    if (formData.password.length < 6) {
      newErrors.password = t('auth.validation.passwordMin');
    }

    // 注册模式的额外验证
    if (!isLogin) {
      // 确认密码验证
      if (formData.confirmPassword !== formData.password) {
        newErrors.confirmPassword = t('auth.validation.passwordMismatch');
      }

      // 邮箱验证
      if (!/\S+@\S+\.\S+/.test(formData.email)) {
        newErrors.email = t('auth.validation.emailInvalid');
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!validateForm()) return;

    setIsSubmitting(true);
    try {
      // 调用外部传入的处理函数
      await onSubmit({
        type: isLogin ? 'login' : 'register',
        data: {
          username: formData.username,
          password: formData.password,
          ...(isLogin
            ? {}
            : {
                email: formData.email,
                confirmPassword: formData.confirmPassword,
              }),
        },
      });
    } catch (error) {
      // 处理错误，可能是显示错误消息
      setErrors((prev) => ({
        ...prev,
        submit: error.message,
      }));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    // 清除对应字段的错误
    if (errors[name]) {
      setErrors((prev) => ({
        ...prev,
        [name]: '',
      }));
    }
  };

  return (
    <motion.div
      className="w-full max-w-[400px] bg-neutral-900 border border-neutral-600 rounded-md p-8"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.2 }}
    >
      {/* 标题区域 */}
      <div className="relative flex justify-center mb-8">
        <div className="absolute -top-[3px] -right-[3px] w-[10px] h-[10px] border-t border-r border-neutral-300" />
        <div className="absolute -bottom-[3px] -left-[3px] w-[10px] h-[10px] border-b border-l border-neutral-300" />
        <motion.h1
          className="text-neutral-300 text-2xl font-mono tracking-wider"
          key={isLogin ? 'login' : 'register'}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.15 }}
        >
          {isLogin ? t('auth.login') : t('auth.register')}
        </motion.h1>
      </div>

      {/* 切换按钮 */}
      <div className="flex justify-center gap-4 mb-8">
        <Button
          variant={isLogin ? 'primary' : 'outline'}
          size="sm"
          className="min-w-[100px]"
          onClick={() => setIsLogin(true)}
        >
          {t('auth.login')}
        </Button>
        <Button
          variant={!isLogin ? 'primary' : 'outline'}
          size="sm"
          className="min-w-[100px]"
          onClick={() => setIsLogin(false)}
        >
          {t('auth.register')}
        </Button>
      </div>

      {/* 账号密码表单 */}
      <motion.form
        onSubmit={handleSubmit}
        className="space-y-4"
        animate={{ height: isLogin ? 'auto' : 'auto' }}
        transition={{ duration: 0.2 }}
      >
        <AnimatePresence mode="wait">
          <motion.div
            key={isLogin ? 'login-form' : 'register-form'}
            initial={{ opacity: 0, x: isLogin ? -20 : 20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: isLogin ? 20 : -20 }}
            transition={{
              duration: 0.2,
              ease: 'easeOut',
            }}
          >
            <div className="space-y-4">
              <Input
                type="text"
                name="username"
                required
                value={formData.username}
                onChange={handleChange}
                label={t('auth.placeholders.username')}
                placeholder={t('auth.placeholders.username')}
                error={errors.username}
              />

              {!isLogin && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: 'auto' }}
                  exit={{ opacity: 0, height: 0 }}
                  transition={{
                    duration: 0.15,
                    ease: 'easeInOut',
                  }}
                >
                  <Input
                    type="email"
                    name="email"
                    required
                    value={formData.email}
                    onChange={handleChange}
                    label={t('auth.placeholders.email')}
                    placeholder={t('auth.placeholders.email')}
                    error={errors.email}
                  />
                </motion.div>
              )}

              <Input
                type="password"
                name="password"
                required
                value={formData.password}
                onChange={handleChange}
                label={t('auth.placeholders.password')}
                placeholder={t('auth.placeholders.password')}
                error={errors.password}
              />

              {!isLogin && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: 'auto' }}
                  exit={{ opacity: 0, height: 0 }}
                  transition={{
                    duration: 0.15,
                    ease: 'easeInOut',
                  }}
                >
                  <Input
                    type="password"
                    name="confirmPassword"
                    required
                    value={formData.confirmPassword}
                    onChange={handleChange}
                    label={t('auth.placeholders.confirmPassword')}
                    placeholder={t('auth.placeholders.confirmPassword')}
                    error={errors.confirmPassword}
                  />
                </motion.div>
              )}
            </div>
          </motion.div>
        </AnimatePresence>

        <Button
          type="submit"
          variant="primary"
          fullWidth
          className="shadow-focus-strong"
          disabled={isSubmitting}
        >
          {isSubmitting ? t('common.processing') : isLogin ? t('auth.login') : t('auth.register')}
        </Button>
      </motion.form>

      {/* OAuth登录 - 独立于表单 */}
      <OAuthLogin />
    </motion.div>
  );
}

export default AuthPanel;
