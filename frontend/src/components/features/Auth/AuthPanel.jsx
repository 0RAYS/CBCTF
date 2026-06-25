import { motion, AnimatePresence, useAnimationControls } from 'motion/react';
import { useState, useEffect } from 'react';
import { Button, Input, Modal } from '../../common';
import OAuthLogin from './OAuthLogin';
import { useTranslation } from 'react-i18next';
import { forgotPassword, getCaptcha } from '../../../api/auth';

function CaptchaField({ captcha, value, onChange, onRefresh, error, disabled }) {
  const { t } = useTranslation();

  return (
    <div className="space-y-2">
      <div className="flex gap-3 items-end">
        <Input
          type="text"
          name="captcha"
          required
          value={value}
          onChange={onChange}
          label={t('auth.captcha.label')}
          placeholder={t('auth.captcha.placeholder')}
          error={error}
          autoComplete="off"
          disabled={disabled}
        />
        <button
          type="button"
          className="h-10 w-32 shrink-0 overflow-hidden rounded-md border border-neutral-600/60 bg-neutral-900/80 transition-colors hover:border-geek-400 disabled:opacity-50"
          onClick={onRefresh}
          disabled={disabled}
          aria-label={t('auth.captcha.refresh')}
        >
          {captcha?.image ? (
            <img src={captcha.image} alt={t('auth.captcha.imageAlt')} className="h-full w-full object-cover" />
          ) : (
            <span className="text-xs text-neutral-400">{t('common.loading')}</span>
          )}
        </button>
      </div>
      <button
        type="button"
        className="text-xs text-neutral-500 hover:text-[#597ef7]"
        onClick={onRefresh}
        disabled={disabled}
      >
        {t('auth.captcha.refresh')}
      </button>
    </div>
  );
}

function ForgotPasswordModal({ isOpen, onClose }) {
  const { t } = useTranslation();
  const [email, setEmail] = useState('');
  const [emailError, setEmailError] = useState('');
  const [captcha, setCaptcha] = useState(null);
  const [captchaValue, setCaptchaValue] = useState('');
  const [captchaError, setCaptchaError] = useState('');
  const [isSending, setIsSending] = useState(false);
  const [sent, setSent] = useState(false);

  const loadCaptcha = async () => {
    const response = await getCaptcha();
    if (response.code === 200) {
      setCaptcha(response.data);
      setCaptchaValue('');
    }
  };

  useEffect(() => {
    if (isOpen && !sent) {
      loadCaptcha();
    }
  }, [isOpen, sent]);

  const handleClose = () => {
    setEmail('');
    setEmailError('');
    setCaptcha(null);
    setCaptchaValue('');
    setCaptchaError('');
    setIsSending(false);
    setSent(false);
    onClose();
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!/\S+@\S+\.\S+/.test(email)) {
      setEmailError(t('auth.validation.emailInvalid'));
      return;
    }
    if (!captcha?.id || !captchaValue.trim()) {
      setCaptchaError(t('auth.validation.captchaRequired'));
      return;
    }
    setIsSending(true);
    try {
      const response = await forgotPassword({ email, captchaId: captcha.id, captcha: captchaValue });
      if (response.code !== 200) {
        throw new Error(response.msg);
      }
      setSent(true);
    } catch {
      setCaptchaError(t('auth.validation.captchaInvalid'));
      await loadCaptcha();
    } finally {
      setIsSending(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title={t('auth.forgotPassword.title')} size="sm">
      {sent ? (
        <div className="space-y-4 text-center">
          <p className="text-neutral-300 text-sm leading-relaxed">{t('auth.forgotPassword.successMessage')}</p>
          <Button variant="primary" fullWidth onClick={handleClose}>
            {t('auth.forgotPassword.backToLogin')}
          </Button>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-4">
          <p className="text-neutral-400 text-sm leading-relaxed">{t('auth.forgotPassword.description')}</p>
          <Input
            type="email"
            name="email"
            required
            value={email}
            onChange={(e) => {
              setEmail(e.target.value);
              if (emailError) setEmailError('');
            }}
            label={t('auth.forgotPassword.emailLabel')}
            placeholder={t('auth.placeholders.email')}
            error={emailError}
          />
          <CaptchaField
            captcha={captcha}
            value={captchaValue}
            onChange={(e) => {
              setCaptchaValue(e.target.value);
              if (captchaError) setCaptchaError('');
            }}
            onRefresh={loadCaptcha}
            error={captchaError}
            disabled={isSending}
          />
          <Button type="submit" variant="primary" fullWidth disabled={isSending}>
            {isSending ? t('auth.forgotPassword.sending') : t('auth.forgotPassword.submit')}
          </Button>
        </form>
      )}
    </Modal>
  );
}

function AuthPanel({ onSubmit, registrationEnabled = false }) {
  const { t } = useTranslation();
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    confirmPassword: '',
    email: '',
    captcha: '',
  });
  const [captcha, setCaptcha] = useState(null);
  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [forgotOpen, setForgotOpen] = useState(false);
  const formControls = useAnimationControls();

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

  const shouldShowCaptcha = true;

  const loadCaptcha = async () => {
    const response = await getCaptcha();
    if (response.code === 200) {
      setCaptcha(response.data);
      setFormData((prev) => ({ ...prev, captcha: '' }));
    }
  };

  useEffect(() => {
    if (shouldShowCaptcha) {
      loadCaptcha();
    } else {
      setCaptcha(null);
      setFormData((prev) => ({ ...prev, captcha: '' }));
    }
  }, [shouldShowCaptcha]);

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

    if (shouldShowCaptcha && (!captcha?.id || !formData.captcha.trim())) {
      newErrors.captcha = t('auth.validation.captchaRequired');
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!registrationEnabled && !isLogin) {
      setIsLogin(true);
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
          captchaId: captcha?.id,
          captcha: formData.captcha,
          ...(isLogin
            ? {}
            : {
                email: formData.email,
                confirmPassword: formData.confirmPassword,
              }),
        },
      });
    } catch (error) {
      await loadCaptcha();
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

  const switchMode = (nextIsLogin) => {
    setIsLogin(nextIsLogin);
    setErrors({});
  };

  return (
    <>
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
            onClick={() => switchMode(true)}
          >
            {t('auth.login')}
          </Button>
          {registrationEnabled && (
            <Button
              variant={!isLogin ? 'primary' : 'outline'}
              size="sm"
              className="min-w-[100px]"
              onClick={() => switchMode(false)}
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

                {shouldShowCaptcha && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    exit={{ opacity: 0, height: 0 }}
                    transition={{ duration: 0.15, ease: [0.25, 1, 0.5, 1] }}
                  >
                    <CaptchaField
                      captcha={captcha}
                      value={formData.captcha}
                      onChange={handleChange}
                      onRefresh={loadCaptcha}
                      error={errors.captcha}
                      disabled={isSubmitting}
                    />
                  </motion.div>
                )}
              </div>
            </motion.div>
          </AnimatePresence>

          {isLogin && (
            <div className="flex justify-end">
              <button
                type="button"
                className="text-xs text-neutral-400 hover:text-[#597ef7] font-mono transition-colors"
                onClick={() => setForgotOpen(true)}
              >
                {t('auth.forgotPassword.link')}
              </button>
            </div>
          )}

          <Button type="submit" variant="primary" fullWidth className="shadow-focus-strong" disabled={isSubmitting}>
            {isSubmitting ? t('common.processing') : isLogin ? t('auth.login') : t('auth.register')}
          </Button>

          {errors.submit && (
            <p className="text-red-400 text-sm font-mono text-center" role="alert">
              {errors.submit}
            </p>
          )}
        </motion.form>

        <OAuthLogin />
      </motion.div>

      <ForgotPasswordModal isOpen={forgotOpen} onClose={() => setForgotOpen(false)} />
    </>
  );
}

export default AuthPanel;
