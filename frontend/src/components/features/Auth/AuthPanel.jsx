import { motion, AnimatePresence, useAnimationControls } from 'motion/react';
import { useState, useEffect } from 'react';
import { Button, Input } from '../../common';
import OAuthLogin from './OAuthLogin';
import { useTranslation } from 'react-i18next';

function AuthPanel({ onSubmit, registrationEnabled = false }) {
  const { t } = useTranslation();
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    confirmPassword: '',
    email: '',
  });
  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const formControls = useAnimationControls();
  const registrationClosedText = t('auth.registrationClosed', {
    defaultValue: 'Registration is currently disabled.',
  });

  useEffect(() => {
    if (!registrationEnabled && !isLogin) {
      setIsLogin(true);
      setErrors((prev) => {
        if (!prev.submit) {
          return prev;
        }
        return { ...prev, submit: '' };
      });
    }
  }, [registrationEnabled, isLogin]);

  useEffect(() => {
    if (errors.submit) {
      formControls.start({
        x: [-8, 8, -6, 6, -3, 3, 0],
        transition: { duration: 0.35 },
      });
    }
  }, [errors.submit, formControls]);

  const validateForm = () => {
    const newErrors = {};

    if (formData.username.length < 3) {
      newErrors.username = t('auth.validation.usernameMin');
    }

    if (formData.password.length < 6) {
      newErrors.password = t('auth.validation.passwordMin');
    }

    if (!isLogin) {
      if (formData.confirmPassword !== formData.password) {
        newErrors.confirmPassword = t('auth.validation.passwordMismatch');
      }

      if (!/\S+@\S+\.\S+/.test(formData.email)) {
        newErrors.email = t('auth.validation.emailInvalid');
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!registrationEnabled && !isLogin) {
      setErrors((prev) => ({
        ...prev,
        submit: registrationClosedText,
      }));
      return;
    }

    if (!validateForm()) return;

    setIsSubmitting(true);
    try {
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

    if (errors[name] || errors.submit) {
      setErrors((prev) => ({
        ...prev,
        [name]: '',
        submit: '',
      }));
    }
  };

  return (
    <motion.div
      className="w-full max-w-[400px] bg-neutral-800/80 border border-neutral-600/60 rounded-md p-8"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, ease: [0.25, 1, 0.5, 1] }}
    >
      <div className="relative flex justify-center mb-8">
        <div className="absolute -top-[3px] -right-[3px] w-[10px] h-[10px] border-t border-r border-neutral-300" />
        <div className="absolute -bottom-[3px] -left-[3px] w-[10px] h-[10px] border-b border-l border-neutral-300" />
        <motion.h1
          className="text-neutral-50 text-2xl font-mono tracking-wider"
          key={isLogin ? 'login' : 'register'}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.15 }}
        >
          {isLogin ? t('auth.login') : t('auth.register')}
        </motion.h1>
      </div>

      <div className="flex justify-center gap-4 mb-8">
        <Button
          variant={isLogin ? 'primary' : 'outline'}
          size="sm"
          className={registrationEnabled ? 'min-w-[100px]' : 'min-w-[220px]'}
          onClick={() => setIsLogin(true)}
        >
          {t('auth.login')}
        </Button>
        {registrationEnabled && (
          <Button
            variant={!isLogin ? 'primary' : 'outline'}
            size="sm"
            className="min-w-[100px]"
            onClick={() => setIsLogin(false)}
          >
            {t('auth.register')}
          </Button>
        )}
      </div>

      <motion.form onSubmit={handleSubmit} className="space-y-4" animate={formControls}>
        <AnimatePresence mode="wait">
          <motion.div
            key={isLogin ? 'login-form' : 'register-form'}
            initial={{ opacity: 0, x: isLogin ? -20 : 20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: isLogin ? 20 : -20 }}
            transition={{
              duration: 0.2,
              ease: [0.25, 1, 0.5, 1],
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
                    ease: [0.25, 1, 0.5, 1],
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
                    ease: [0.25, 1, 0.5, 1],
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

        <Button type="submit" variant="primary" fullWidth className="shadow-focus-strong" disabled={isSubmitting}>
          {isSubmitting ? t('common.processing') : isLogin ? t('auth.login') : t('auth.register')}
        </Button>

        {!registrationEnabled && (
          <p className="text-neutral-400 text-sm font-mono text-center" role="status">
            {registrationClosedText}
          </p>
        )}

        {errors.submit && (
          <p className="text-red-400 text-sm font-mono text-center" role="alert">
            {errors.submit}
          </p>
        )}
      </motion.form>

      <OAuthLogin />
    </motion.div>
  );
}

export default AuthPanel;
